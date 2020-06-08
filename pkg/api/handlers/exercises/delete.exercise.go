package exercises

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *ExerciseHandler) DeleteExercise(c echo.Context) (err error) {
	exerciseID := c.Param("id")

	ctx := context.Background()

	exerciseItem, err := models.GetAExercise(h.DB, bson.M{"_id": exerciseID})
	if err != nil {
		return err
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	userRole := claims["role"].(string)

	filter := bson.M{"_id": exerciseItem.ClassID}
	if userRole != "ADMIN" {
		filter["teachers._id"] = userID
	}

	_, err = models.GetAClass(h.DB, filter)

	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now().UTC(),
		},
	}

	exercise := &models.Exercise{}
	exerciseCollection := models.GetExerciseCollection(h.DB)
	result := exerciseCollection.FindOneAndUpdate(ctx, bson.M{"_id": exerciseID}, update, options.FindOneAndUpdate().SetReturnDocument(1))
	if err = result.Decode(&exercise); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found exercise",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[DeleteExercise] Internal server error",
			Internal: err,
		}
	}
	return c.NoContent(http.StatusNoContent)
}
