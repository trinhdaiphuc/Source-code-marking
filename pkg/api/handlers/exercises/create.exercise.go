package exercises

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *ExerciseHandler) CreateExercise(c echo.Context) (err error) {
	exerciseItem := &models.Exercise{}
	if err = c.Bind(exerciseItem); err != nil {
		return
	}

	h.Logger.Debug("Create Exercise parameters ", exerciseItem)

	classItem := &models.Class{}
	classCollection := models.GetClassCollection(h.DB)
	result := classCollection.FindOne(context.Background(), bson.M{"_id": exerciseItem.ClassID})
	if err := result.Decode(&classItem); err != nil {
		h.Logger.Info("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found class ",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Profile] Internal server error ",
			Internal: err,
		}
	}

	exerciseItem.ID = uuid.NewV4().String()
	exerciseItem.IsOpen = false
	exerciseItem.CreatedAt = time.Now().UTC()
	exerciseItem.UpdatedAt = time.Now().UTC()

	ExerciseCollection := models.GetExerciseCollection(h.DB)
	_, err = ExerciseCollection.InsertOne(context.Background(), exerciseItem)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "MongoDB is not avalable.",
			Internal: err,
		}
	}

	return c.JSON(http.StatusCreated, exerciseItem)
}
