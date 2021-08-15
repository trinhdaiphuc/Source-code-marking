package notifications

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func sendNotification(db *mongo.Client, redisClient *redis.Client, userID, userEmail string) {
	filter := bson.M{"user_id": userID, "is_deleted": false}

	listParam := models.ListQueryParam{
		PageSize:  5,
		PageToken: 1,
		OrderBy:   "created_at",
		OrderType: internal.DESC.String(),
	}

	listNotification, _ := models.ListAllNotifications(db, filter, listParam)
	notificationCollection := models.GetNotificationCollection(db)
	filter["is_read"] = false
	totalUnread, _ := notificationCollection.CountDocuments(context.TODO(), filter)
	listNotificationWebsocket := models.ListNotificationWebsocket{
		Notifications: listNotification.Notifications,
		TotalUnread:   totalUnread,
	}
	message, _ := json.Marshal(listNotificationWebsocket)
	redisClient.Publish(context.Background(), userEmail, message).Err()
}

func (h *NotificationHandler) MarkReadNotification(c echo.Context) (err error) {
	id := c.Param("id")
	filter := bson.M{"_id": id, "is_deleted": false}
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
		},
	}
	notificationItem := &models.Notification{}
	notificationCollection := models.GetNotificationCollection(h.DB)
	result := notificationCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	if err := result.Decode(&notificationItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found notification",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update notification] Internal server error",
			Internal: err,
		}
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	userEmail := claims["email"].(string)

	go sendNotification(h.DB, h.RedisClient, userID, userEmail)

	return c.NoContent(http.StatusNoContent)
}
