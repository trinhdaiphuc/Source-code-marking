package exercises

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
)

func (h *ExerciseHandler) CreateExercise(c echo.Context) (err error) {
	ExerciseItem := &models.Exercise{}
	if err = c.Bind(ExerciseItem); err != nil {
		return
	}

	h.Logger.Debug("Create Exercise parameters ", ExerciseItem)

	ExerciseItem.ID = uuid.NewV4().String()
	ExerciseItem.IsOpen = false
	ExerciseItem.CreatedAt = time.Now().UTC()
	ExerciseItem.UpdatedAt = time.Now().UTC()

	ExerciseCollection := models.GetExerciseCollection(h.DB)
	_, err = ExerciseCollection.InsertOne(context.Background(), ExerciseItem)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "MongoDB is not avalable.",
			Internal: err,
		}
	}

	return c.JSON(http.StatusCreated, ExerciseItem)
}
