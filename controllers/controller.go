package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RohithBN/netflix-api/model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ConnectionString string

const dbName = "netflix"
const colName = "watchlist"

var collection *mongo.Collection

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to load .env file:", err)
	}

	ConnectionString = os.Getenv("MONGO_URI")
	if ConnectionString == "" {
		log.Fatal("MONGO_URI is not set in the environment variables")
	}

	log.Println("MongoDB connection string loaded successfully")

	clientOption := options.Client().ApplyURI(ConnectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDb connection ssuccesssful")
	collection = client.Database(dbName).Collection(colName)
}

func addMovie(movie model.Netflix) {
	result, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The movie is inserted with id:", result.InsertedID)
}

func UpdateMovie(movieId string) bool {
	id, err := primitive.ObjectIDFromHex(movieId)
	if err != nil {
		log.Printf("Error converting movieId to ObjectID: %v", err)
		return false
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Printf("Error updating movie: %v", err)
		return false
	}

	if result.MatchedCount == 0 {
		log.Println("No movie found with the given ID")
		return false
	}

	log.Printf("Movie updated successfully with ID: %v", movieId)
	return true
}

func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The movie is deleted with id:", result.DeletedCount)
}

func deleteAllMovies() float64 {
	result, err := collection.DeleteMany(context.Background(), bson.M{}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("The movies are deleted with count:", result.DeletedCount)
	return float64(result.DeletedCount)
}

func getAllMovies() []primitive.M {
	findOptions := options.Find()

	cursor, err := collection.Find(context.Background(), bson.D{}, findOptions)
	if err != nil {
		log.Println("Error during collection.Find:", err)

	}
	defer cursor.Close(context.Background())

	var results []primitive.M

	for cursor.Next(context.TODO()) {
		var elem bson.M
		if err := cursor.Decode(&elem); err != nil {
			log.Println("Error during cursor.Decode:", err)

		}

		results = append(results, elem)
	}

	if err := cursor.Err(); err != nil {
		log.Println("Error during cursor iteration:", err)

	}

	return results
}

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	movies := getAllMovies()
	err := json.NewEncoder(w).Encode(movies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func MarkMovieWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	id := params["id"]

	success := UpdateMovie(id)
	if success {
		log.Println("Movie marked as watched successfully") // Debugging log
		err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Movie marked as watched successfully",
			"id":      id,
		})
		if err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	} else {
		http.Error(w, "Movie not found", http.StatusNotFound)
	}
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	if r.Body == nil {
		http.Error(w, "No request body provided", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var movie model.Netflix
	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	addMovie(movie)
	err = json.NewEncoder(w).Encode(movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	params := mux.Vars(r)
	id := params["id"]
	deleteOneMovie(id)
	json.NewEncoder(w).Encode(id)
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	count := deleteAllMovies()
	json.NewEncoder(w).Encode(count)
}
