package models

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Role struct {
	ID   int    `json:"id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}

func newRoleCollection(db *mongo.Client) {
	roleCollection := getDatabase(db).Collection("roles")

	roleCollection.InsertMany(context.Background(), []interface{}{
		Role{ID: 1, Name: "ADMIN"},
		Role{ID: 2, Name: "TEACHER"},
		Role{ID: 3, Name: "STUDENT"},
	})
}

func GetRoleCollection(db *mongo.Client) *mongo.Collection {
	roleCollection := getDatabase(db).Collection("roles")
	return roleCollection
}

func GetARole(db *mongo.Client, filter bson.M) (*Role, error) {
	roleCollection := GetRoleCollection(db)
	resultFind := roleCollection.FindOne(context.TODO(), filter)

	data := &Role{}
	if err := resultFind.Decode(&data); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "Invalid role",
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
