package models

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Notification struct {
		ID         string    `json:"id" bson:"_id"`
		UserID     string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
		IsRead     bool      `json:"is_read" bson:"is_read"`
		Content    string    `json:"content" bson:"content"`
		ExerciseID string    `json:"exercise_id" bson:"exercise_id"`
		IsDeleted  bool      `json:"is_deleted" bson:"is_deleted"`
		CreatedAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
		UpdatedAt  time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
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

func ListAllNotifications(db *mongo.Client, filter bson.M, listParam ListQueryParam) (*ListNotification, error) {
	limit := listParam.PageSize
	page := listParam.PageToken
	skip := (page - 1) * limit
	orderBy := "created_at"
	orderType := 1
	if listParam.OrderType == internal.DESC.String() {
		orderType = -1
	}

	if listParam.OrderBy != "" {
		orderBy = listParam.OrderBy
	}

	opts := []*options.FindOptions{}
	opts = append(opts, options.Find().SetSort(bson.D{{orderBy, orderType}}))
	opts = append(opts, options.Find().SetSkip(skip))
	opts = append(opts, options.Find().SetLimit(limit))

	notificationCollection := GetNotificationCollection(db)
	ctx := context.TODO()
	cursor, err := notificationCollection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Get all user] Internal server error",
			Internal: err,
		}
	}
	if cursor == nil {
		return nil, &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Not found notifications",
		}
	}

	defer cursor.Close(ctx)

	notificationArray := []Notification{}
	cursor.All(ctx, &notificationArray)
	totalRecords, err := notificationCollection.CountDocuments(ctx, filter)

	return ConvertNotificationArrayToListNotification(notificationArray, page+1, totalRecords), nil
}
