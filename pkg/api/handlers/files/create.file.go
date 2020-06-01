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

func fileExisted(fileCollection *mongo.Collection, userID, exerciseID string, c chan *models.File) {

}

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

	exerciseCollection := models.GetExerciseCollection(h.DB)
	resultFind := exerciseCollection.FindOne(context.Background(), bson.M{"_id": fileItem.ExerciseID})

	exercise := &models.Exercise{}
	if err := resultFind.Decode(&exercise); err != nil {
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

	h.Logger.Debug("Time deadline ", exercise.Deadline.Sub(time.Now()))

	if exercise.Deadline.Sub(time.Now()) < 0 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Over deadline.",
		}
	}

	filter := bson.M{"user_id": userID, "exercise_id": fileItem.ExerciseID}
	result := fileCollection.FindOne(context.Background(), filter)

	data := &models.File{}
	if err := result.Decode(&data); err != nil {
		if err != mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "[GetExercise] Internal server error",
				Internal: err,
			}
		}
	}

	if data != nil {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "File has already existed.",
		}
	}

	fileItem.UserID = userID
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
