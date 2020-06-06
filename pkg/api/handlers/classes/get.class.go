package classes

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *ClassHandler) GetClass(c echo.Context) (err error) {
	classID := c.Param("id")
	classItem := &models.Class{}
	classCollection := models.GetClassCollection(h.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	filter := bson.M{"_id": classID}

	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}

	result := classCollection.FindOne(context.Background(), filter)
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
	return c.JSON(http.StatusOK, classItem)
}
