package exercises

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *ExerciseHandler) UpdateExercise(c echo.Context) (err error) {
	h.Logger.Debug("Get user handler")
	// Get param
	Exercise := &models.Exercise{}
	ExerciseID := c.Param("id")
	if err := c.Bind(Exercise); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid parameters",
			Internal: err,
		}
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	filter := bson.M{"_id": ExerciseID}
	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}

	data, err := models.GetAExercise(h.DB, filter)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":        Exercise.Name,
			"description": Exercise.Description,
			"deadline":    Exercise.Deadline,
			"updated_at":  time.Now().UTC(),
		},
	}

	filter = bson.M{"_id": ExerciseID}
	exerciseCollection := models.GetExerciseCollection(h.DB)
	resultUpdate := exerciseCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	err = resultUpdate.Decode(&data)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}

	return c.JSON(http.StatusOK, data)
}
