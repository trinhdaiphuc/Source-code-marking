package files

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

	exerciseCollection := models.GetExerciseCollection(h.DB)
	resultFind := exerciseCollection.FindOne(context.Background(), bson.M{"_id": fileItem.ExerciseID})

	Exercise := &models.Exercise{}
	if err := resultFind.Decode(&Exercise); err != nil {
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

	fileItem.UserID = userID
	fileItem.CreatedAt = time.Now().UTC()
	fileItem.UpdatedAt = time.Now().UTC()

	fileCollection := models.GetFileCollection(h.DB)
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
