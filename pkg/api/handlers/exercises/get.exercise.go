package exercises

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
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
	filter := bson.M{"_id": ExerciseID}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}
	resultFind := exerciseCollection.FindOne(context.Background(), filter)

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
	return c.JSON(http.StatusOK, Exercise)
}
