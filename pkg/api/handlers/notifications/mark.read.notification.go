package notifications

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	return c.NoContent(http.StatusNoContent)
}
