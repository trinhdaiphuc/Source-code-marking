package classes

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *ClassHandler) GetClass(c echo.Context) (err error) {
	classID := c.Param("id")

	classCollection := models.GetClassCollection(h.DB)
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	filter := bson.M{"_id": classID}

	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}

	classItem, err := models.GetAClass(classCollection, filter)

	if err != nil {
		h.Logger.Error("Error when get a class ", err)
		return err
	}

	return c.JSON(http.StatusOK, classItem)
}
