package models

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Role struct {
	ID   int    `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

func newRoleCollection(db *mongo.Client) {
	roleCollection := db.Database("Source-code-marking").Collection("roles")

	roleCollection.InsertMany(context.Background(), []interface{}{
		Role{ID: 1, Name: "ADMIN"},
		Role{ID: 2, Name: "TEACHER"},
		Role{ID: 3, Name: "STUDENT"},
	})
}

func GetRoleCollection(db *mongo.Client) *mongo.Collection {
	roleCollection := db.Database("Source-code-marking").Collection("roles")
	return roleCollection
}
