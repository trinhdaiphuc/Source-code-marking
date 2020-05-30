package models

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Class struct {
		ID          string    `json:"id,omitempty" bson:"_id"`
		Name        string    `json:"name" bson:"name"`
		Description string    `json:"description" bson:"description"`
		Teachers    []User    `json:"teachers" bson:"teachers"`
		Students    []User    `json:"students,omitempty" bson:"students,omitempty"`
		CreatedAt   time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt   time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	}

	ListClass struct {
		Classes       []Class `json:"classes"`
		NextPageToken int64   `json:"next_page_token"`
		TotalRecords  int64   `json:"total_records"`
	}
)

func GetClassCollection(db *mongo.Client) *mongo.Collection {
	classCollection := db.Database("Source-code-marking").Collection("classes")
	return classCollection
}

func ConvertClassArrayToListClass(Classs []Class, nextPageToken, totalRecords int64) *ListClass {
	listClass := &ListClass{
		Classes:       Classs,
		NextPageToken: nextPageToken,
		TotalRecords:  totalRecords,
	}

	for i := range listClass.Classes {
		listClass.Classes[i].Students = nil
	}

	return listClass
}
