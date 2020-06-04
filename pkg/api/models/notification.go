package models

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Notification struct {
		ID        string    `json:"id,omitempty" bson:"_id"`
		UserID    string    `json:"user_id" bson:"user_id"`
		IsRead    string    `json:"is_read" bson:"is_read"`
		Content   string    `json:"content" bson:"content"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	}

	ListNotification struct {
		Notifications []Notification `json:"notifications"`
		NextPageToken int64          `json:"next_page_token"`
		TotalRecords  int64          `json:"total_records"`
	}
)

func GetNotificationCollection(db *mongo.Client) *mongo.Collection {
	NotificationCollection := getDatabase(db).Collection("notifications")
	return NotificationCollection
}

func ConvertNotificationArrayToListNotification(Notifications []Notification, nextPageToken, totalRecords int64) *ListNotification {
	listNotification := &ListNotification{
		Notifications: Notifications,
		NextPageToken: nextPageToken,
		TotalRecords:  totalRecords,
	}

	return listNotification
}
