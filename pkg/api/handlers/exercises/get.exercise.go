package exercises

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *ExerciseHandler) GetExercise(c echo.Context) (err error) {
	h.Logger.Debug("Get Exercise handler")
	// Get param
	exerciseID := c.Param("id")
	filter := bson.M{"_id": exerciseID}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}

	exercise, err := models.GetAExercise(h.DB, filter)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, exercise)
}
