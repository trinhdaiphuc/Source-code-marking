package exercises

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *ExerciseHandler) GetExercise(c echo.Context) (err error) {
	h.Logger.Debug("Get Exercise handler")
	// Get param
	ExerciseID := c.Param("id")

	exerciseCollection := models.GetExerciseCollection(h.DB)
	resultFind := exerciseCollection.FindOne(context.Background(), bson.M{"_id": ExerciseID})

	Exercise := &models.Exercise{}
	if err := resultFind.Decode(&Exercise); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found Exercise",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[GetExercise] Internal server error",
			Internal: err,
		}
	}
	return c.JSON(http.StatusOK, Exercise)
}
