package models

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Exercise struct {
		ID          string    `json:"id,omitempty" bson:"_id"`
		ClassID     string    `json:"class_id" bson:"class_id" validate:"required"`
		Name        string    `json:"name" bson:"name"`
		Description string    `json:"description" bson:"description"`
		Deadline    time.Time `json:"deadline" bson:"deadline" validate:"required"`
		IsOpen      bool      `json:"is_open" bson:"is_open"`
		IsDeleted   bool      `json:"is_deleted" bson:"is_deleted"`
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
	ExerciseCollection := getDatabase(db).Collection("exercises")
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

func GetAExercise(db *mongo.Client, filter bson.M) (*Exercise, error) {
	exerciseCollection := GetExerciseCollection(db)
	result := exerciseCollection.FindOne(context.Background(), filter)

	exerciseItem := &Exercise{}

	if err := result.Decode(&exerciseItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found exercise",
				Internal: err,
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Delete exercise] Internal server error",
			Internal: err,
		}
	}
	return exerciseItem, nil
}
