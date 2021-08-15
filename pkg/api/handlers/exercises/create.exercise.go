package exercises

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *ExerciseHandler) CreateExercise(c echo.Context) (err error) {
	exerciseItem := &models.Exercise{}

	if err := c.Bind(exerciseItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	if err := c.Validate(exerciseItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	filter := bson.M{"_id": exerciseItem.ClassID, "is_deleted": false, "teachers._id": userID}
	_, err = models.GetAClass(h.DB, filter)
	if err != nil {
		return err
	}

	exerciseItem.ID = uuid.NewV4().String()
	exerciseItem.IsOpen = false
	exerciseItem.IsDeleted = false
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
