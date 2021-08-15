package files

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *FileHandler) UpdateFile(c echo.Context) (err error) {
	h.Logger.Debug("Get user handler")
	// Get param
	file := &models.File{}
	fileID := c.Param("id")
	if err := c.Bind(file); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid parameters",
			Internal: err,
		}
	}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	fileItem, err := models.GetAFile(h.DB, bson.M{"_id": fileID, "is_deleted": false})
	if err != nil {
		return err
	}

	if userID != fileItem.UserID {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "User cannot update this file.",
		}
	}

	exercise, err := models.GetAExercise(h.DB, bson.M{"_id": fileItem.ExerciseID})
	if err != nil {
		return err
	}

	if exercise.Deadline.Sub(time.Now()) < 0 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Over deadline.",
		}
	}

	update := bson.M{
		"$set": bson.M{
			"name":       file.Name,
			"data":       file.Data,
			"updated_at": time.Now().UTC(),
		},
	}
	filter := bson.M{"_id": fileID, "is_deleted": false}
	fileCollection := models.GetFileCollection(h.DB)
	resultUpdate := fileCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	err = resultUpdate.Decode(&fileItem)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}

	return c.JSON(http.StatusOK, fileItem)
}
