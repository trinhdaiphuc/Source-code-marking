package users

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *UserHandler) DeleteUser(c echo.Context) (err error) {
	id := c.Param("id")

	filter := bson.M{"_id": id}
	userCollection := models.GetUserCollection(h.DB)
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now().UTC(),
		},
	}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Delete user] Internal server error",
			Internal: err,
		}
	}
	return c.NoContent(http.StatusNoContent)
}
