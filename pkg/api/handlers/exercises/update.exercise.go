package exercises

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	ctx := context.Background()
	exerciseCollection := models.GetExerciseCollection(h.DB)
	resultFind := exerciseCollection.FindOne(ctx, bson.M{"_id": ExerciseID})

	data := models.Exercise{}
	if err := resultFind.Decode(&data); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "Not found Exercise",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[UpdateExercise] Internal server error",
			Internal: err,
		}
	}

	update := bson.M{
		"$set": bson.M{
			"name":        Exercise.Name,
			"description": Exercise.Description,
			"deadline":    Exercise.Deadline,
			"updated_at":  time.Now().UTC(),
		},
	}
	filter := bson.M{"_id": ExerciseID}

	resultUpdate := exerciseCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
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
