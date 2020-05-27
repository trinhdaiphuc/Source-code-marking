package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Email     string    `json:"email" form:"email" bson:"email"`
	Password  string    `json:"password,omitempty" form:"password" bson:"password"`
	Name      string    `json:"name,omitempty" form:"name" bson:"name,omitempty"`
	Role      string    `json:"role,omitempty" form:"role" bson:"role,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

type ListUser struct {
	Users         []User `json:"users"`
	NextPageToken int64  `json:"next_page_token"`
	TotalRecords  int64  `json:"total_records"`
}

func newUserCollection(db *mongo.Client) {
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

func ConvertUserArrayToListUser(users []User, nextPageToken, totalRecords int64) *ListUser {
	listUser := &ListUser{
		Users:         users,
		NextPageToken: nextPageToken,
		TotalRecords:  totalRecords,
	}
	for i, _ := range listUser.Users {
		listUser.Users[i].Password = ""
	}
	return listUser
}
