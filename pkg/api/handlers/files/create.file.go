package files

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
)

func (h *FileHandler) CreateFile(c echo.Context) (err error) {
	fileItem := &models.File{}
	if err = c.Bind(fileItem); err != nil {
		return
	}

	h.Logger.Debug("Create file parameters ", fileItem)

	fileItem.ID = uuid.NewV4().String()
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	fileItem.UserID = userID
	fileItem.CreatedAt = time.Now().UTC()
	fileItem.UpdatedAt = time.Now().UTC()

	h.Logger.Debug("Sign-in parameters: ", *fileItem)
	fileCollection := models.GetFileCollection(h.DB)
	_, err = fileCollection.InsertOne(context.Background(), fileItem)

	if err != nil {
		h.Logger.Debug("Error when sign-up ", err.Error())
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "MongoDB is not avalable.",
			Internal: err,
		}
	}

	return c.JSON(http.StatusCreated, fileItem)
}
