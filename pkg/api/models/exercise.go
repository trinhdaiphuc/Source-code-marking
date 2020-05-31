package models

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Exercise struct {
		ID          string    `json:"id,omitempty" bson:"_id"`
		ClassID     string    `json:"class_id" bson:"class_id"`
		Name        string    `json:"name" bson:"name"`
		Description string    `json:"description" bson:"description"`
		Deadline    time.Time `json:"deadline" bson:"deadline"`
		IsOpen      bool      `json:"is_open" bson:"is_open"`
		CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	ListExercise struct {
		Exercises     []Exercise `json:"exercisees"`
		NextPageToken int64      `json:"next_page_token"`
		TotalRecords  int64      `json:"total_records"`
	}
)

func GetExerciseCollection(db *mongo.Client) *mongo.Collection {
	ExerciseCollection := getDatabase(db).Collection("exercisees")
	return ExerciseCollection
}

func ConvertExerciseArrayToListExercise(Exercises []Exercise, nextPageToken, totalRecords int64) *ListExercise {
	listExercise := &ListExercise{
		Exercises:     Exercises,
		NextPageToken: nextPageToken,
		TotalRecords:  totalRecords,
	}

	return listExercise
}
