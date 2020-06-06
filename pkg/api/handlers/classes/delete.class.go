package classes

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *ClassHandler) DeleteClass(c echo.Context) (err error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	classID := c.Param("id")
	h.Logger.Debug("UserID ", userID, ", ClassID ", classID)
	ctx := context.Background()
	classCollection := models.GetClassCollection(h.DB)
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now().UTC(),
		},
	}

	filter := bson.M{"_id": classID, "teachers._id": userID}
	_, err = classCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[DeleteClass] Internal server error",
			Internal: err,
		}
	}
	return c.NoContent(http.StatusNoContent)
}
