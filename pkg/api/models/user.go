package models

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	User struct {
		ID         string    `json:"id,omitempty" bson:"_id,omitempty"`
		Email      string    `json:"email" bson:"email" validate:"required,email"`
		Password   string    `json:"password,omitempty" bson:"password,omitempty"`
		Name       string    `json:"name,omitempty" bson:"name,omitempty"`
		Role       string    `json:"role,omitempty" bson:"role,omitempty" validate:"required"`
		IsVerified bool      `json:"is_verified" bson:"is_verified"`
		Service    string    `json:"service,omitempty" bson:"service,omitempty"`
		IsDeleted  bool      `json:"is_deleted" bson:"is_deleted"`
		CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	}

	ListUser struct {
		Users         []User `json:"users"`
		NextPageToken int64  `json:"next_page_token"`
		TotalRecords  int64  `json:"total_records"`
	}

	ResetPassword struct {
		Password string `json:"password" validate:"required,min=6"`
	}
)

func newUserCollection(db *mongo.Client) {
	// Create indexs
	mod := mongo.IndexModel{
		Keys: bson.M{
			"email": -1, // index in ascending order
		},
		// create UniqueIndex option
		Options: options.Index().SetUnique(true),
	}
	ctx := context.Background()
	userCollection := getDatabase(db).Collection("users")
	userCollection.Indexes().CreateOne(ctx, mod)
	password, _ := internal.HashPassword("phucdeptrai")
	admin := &User{
		ID:         uuid.NewV4().String(),
		Email:      "admin@gmail.com",
		Password:   password,
		Name:       "admin",
		Role:       "ADMIN",
		IsVerified: true,
		IsDeleted:  false,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
	userCollection.InsertOne(ctx, admin)
}

func GetUserCollection(db *mongo.Client) *mongo.Collection {
	userCollection := getDatabase(db).Collection("users")
	return userCollection
}

func ConvertUserArrayToListUser(users []User, nextPageToken, totalRecords int64) *ListUser {
	listUser := &ListUser{
		Users:         users,
		NextPageToken: nextPageToken,
		TotalRecords:  totalRecords,
	}
	for i := range listUser.Users {
		listUser.Users[i].Password = ""
	}
	return listUser
}
