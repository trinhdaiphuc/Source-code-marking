package files

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *FileHandler) DeleteFile(c echo.Context) (err error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	fileID := c.Param("id")
	ctx := context.Background()

	fileItem, err := models.GetAFile(h.DB, bson.M{"_id": fileID, "is_deleted": false})

	if err != nil {
		return
	}

	filter := bson.M{"_id": fileItem.ExerciseID, "is_deleted": false}
	exercise, err := models.GetAExercise(h.DB, filter)
	if err != nil {
		return err
	}

	if exercise.Deadline.Sub(time.Now()) < 0 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Over deadline",
		}
	}

	fileCollection := models.GetFileCollection(h.DB)
	_, err = fileCollection.DeleteOne(ctx, bson.M{"_id": fileID, "user_id": userID})
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[DeleteFile] Internal server error",
			Internal: err,
		}
	}
	return c.NoContent(http.StatusNoContent)
}
