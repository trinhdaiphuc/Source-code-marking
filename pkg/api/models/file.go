package models

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type File struct {
	ID         string    `json:"id,omitempty" bson:"_id"`
	UserID     string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ExerciseID string    `json:"exercise_id,omitempty" bson:"exercise_id,omitempty" validate:"required"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Data       string    `json:"data" bson:"data" validate:"required"`
	Mark       float32   `json:"mark,omitempty" bson:"mark,omitempty" validate:"gte=0,lte=10"`
	Comments   []Comment `json:"comments,omitempty" bson:"comments,omitempty"`
	IsDeleted  bool      `json:"is_deleted" bson:"is_deleted"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type ListFile struct {
	Files         []File `json:"files"`
	NextPageToken int64  `json:"next_page_token"`
	TotalRecords  int64  `json:"total_records"`
}

func GetFileCollection(db *mongo.Client) *mongo.Collection {
	fileCollection := getDatabase(db).Collection("files")
	return fileCollection
}

func ConvertFileArrayToListFile(Files []File, nextPageToken, totalRecords int64) *ListFile {
	listFile := &ListFile{
		Files:         Files,
		NextPageToken: nextPageToken,
		TotalRecords:  totalRecords,
	}
	for i := range listFile.Files {
		listFile.Files[i].Comments = nil
	}
	return listFile
}

func GetAFile(db *mongo.Client, filter bson.M) (*File, error) {
	fileCollection := GetFileCollection(db)
	result := fileCollection.FindOne(context.TODO(), filter)

	data := &File{}
	if err := result.Decode(&data); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found file",
				Internal: err,
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal server error",
			Internal: err,
		}
	}

	return data, nil
}
