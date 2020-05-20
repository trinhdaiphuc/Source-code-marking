package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string    `json:"email" form:"email" bson:"email"`
	Password  string    `json:"password,omitempty" form:"password" bson:"password"`
	Name      string    `json:"name,omitempty" form:"name" bson:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

func NewUserCollection(db *mongo.Client) {
	// Create indexs
	mod := mongo.IndexModel{
		Keys: bson.M{
			"email": -1, // index in ascending order
		},
		// create UniqueIndex option
		Options: options.Index().SetUnique(true),
	}
	userCollection := db.Database("Source-code-marking").Collection("users")
	userCollection.Indexes().CreateOne(context.Background(), mod)
}

func GetUserCollection(db *mongo.Client) *mongo.Collection {
	userCollection := db.Database("Source-code-marking").Collection("users")
	return userCollection
}
