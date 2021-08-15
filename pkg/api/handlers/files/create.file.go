package files

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

func (h *FileHandler) CreateFile(c echo.Context) (err error) {
	fileItem := &models.File{}
	if err := c.Bind(fileItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	if err := c.Validate(fileItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	h.Logger.Debug("Create file parameters ", fileItem)

	fileItem.ID = uuid.NewV4().String()
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	fileCollection := models.GetFileCollection(h.DB)

	filter := bson.M{"_id": fileItem.ExerciseID}
	exercise, err := models.GetAExercise(h.DB, filter)
	if err != nil {
		return err
	}

	if exercise.Deadline.Sub(time.Now()) < 0 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Over deadline.",
		}
	}

	filter = bson.M{"user_id": userID, "exercise_id": fileItem.ExerciseID}
	data, err := models.GetAFile(h.DB, filter)
	if err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if code != http.StatusNotFound {
			return err
		}
	}

	if data != nil && data.ID != "" {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "File has already existed.",
		}
	}

	fileItem.UserID = userID
	fileItem.IsDeleted = false
	fileItem.CreatedAt = time.Now().UTC()
	fileItem.UpdatedAt = time.Now().UTC()

	_, err = fileCollection.InsertOne(context.Background(), fileItem)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "MongoDB is not avalable.",
			Internal: err,
		}
	}

	return c.JSON(http.StatusCreated, fileItem)
}
