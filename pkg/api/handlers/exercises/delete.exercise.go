package exercises

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *ExerciseHandler) DeleteExercise(c echo.Context) (err error) {
	ExerciseID := c.Param("id")

	ctx := context.Background()
	ExerciseCollection := models.GetExerciseCollection(h.DB)
	_, err = ExerciseCollection.DeleteOne(ctx, bson.M{"_id": ExerciseID})
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[DeleteExercise] Internal server error",
			Internal: err,
		}
	}
	return c.NoContent(204)
}
