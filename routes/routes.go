package routes

import (
	controller "github.com/RohithBN/netflix-api/controllers"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router:=mux.NewRouter()
	router.HandleFunc("/api/movies",controller.GetAllMovies).Methods("GET")
	router.HandleFunc("/api/delete-movies",controller.DeleteAllMovies).Methods("DELETE")
	router.HandleFunc("/api/movie/{id}",controller.MarkMovieWatched).Methods("PUT")
	router.HandleFunc("/api/movie/{id}",controller.DeleteOneMovie).Methods("DELETE")
	router.HandleFunc("/api/addmovie",controller.CreateMovie).Methods("POST")

	return router

}