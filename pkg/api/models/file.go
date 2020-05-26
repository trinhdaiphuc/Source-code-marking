package models

import "go.mongodb.org/mongo-driver/mongo"

type File struct {
	ID         string    `json:"id,omitempty" bson:"_id"`
	UserID     string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ExerciseID string    `json:"exercise_id,omitempty" bson:"exercise_id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Data       string    `json:"data" bson:"data"`
	Comments   []Comment `json:"comments,omitempty" bson:"comments,omitempty"`
}

func GetFileCollection(db *mongo.Client) *mongo.Collection {
	userCollection := db.Database("Source-code-marking").Collection("files")
	return userCollection
}
