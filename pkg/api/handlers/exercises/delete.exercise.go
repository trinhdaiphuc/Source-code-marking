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
	exerciseCollection := models.GetExerciseCollection(h.DB)
	result := exerciseCollection.FindOne(ctx, bson.M{"_id": exerciseID})

	exerciseItem := &models.Exercise{}

	if err = result.Decode(&exerciseItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found exercise",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Delete exercise] Internal server error",
			Internal: err,
		}
	}

	classItem := &models.Class{}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	classCollection := models.GetClassCollection(h.DB)
	result = classCollection.FindOne(ctx, bson.M{"_id": exerciseItem.ClassID, "teachers._id": userID})

	if err = result.Decode(&classItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found class with teacher id",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Delete exercise] Internal server error",
			Internal: err,
		}
	}

	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now().UTC(),
		},
	}

	exercise := &models.Exercise{}

	result = exerciseCollection.FindOneAndUpdate(ctx, bson.M{"_id": exerciseID}, update, options.FindOneAndUpdate().SetReturnDocument(1))
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
