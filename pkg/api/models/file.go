package models

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type File struct {
	ID         string    `json:"id,omitempty" bson:"_id"`
	UserID     string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ExerciseID string    `json:"exercise_id,omitempty" bson:"exercise_id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Data       string    `json:"data" bson:"data"`
	Mark       float32   `json:"mark,omitempty" bson:"mark,omitempty"`
	Comments   []Comment `json:"comments,omitempty" bson:"comments,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type ListFile struct {
	Files         []File `json:"files"`
	NextPageToken int64  `json:"next_page_token"`
	TotalRecords  int64  `json:"total_records"`
}

func GetFileCollection(db *mongo.Client) *mongo.Collection {
	fileCollection := db.Database("Source-code-marking").Collection("files")
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
