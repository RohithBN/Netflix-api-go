package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Netflix represents a movie document in the database
type Netflix struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Movie   string             `json:"movie" bson:"movie"`
	Year    int                `json:"year" bson:"year"`
	Stars   int                `json:"stars" bson:"stars"`
	Watched bool               `json:"watched" bson:"watched"`
}
