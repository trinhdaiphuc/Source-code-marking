package files

import (
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *FileHandler) GetFile(c echo.Context) (err error) {
	h.Logger.Debug("Get file handler")
	// Get param
	fileID := c.Param("id")
	filter := bson.M{"_id": fileID}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)
	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}

	file, err := models.GetAFile(h.DB, filter)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, file)
}
