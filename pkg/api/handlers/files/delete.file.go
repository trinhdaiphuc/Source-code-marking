package files

import (
	"context"
	"net/http"

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
	h.Logger.Debug("UserID ", userID, ", FileID ", fileID)
	ctx := context.Background()
	fileCollection := models.GetFileCollection(h.DB)
	_, err = fileCollection.DeleteOne(ctx, bson.M{"_id": fileID, "user_id": userID})
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[DeleteFile] Internal server error",
			Internal: err,
		}
	}
	return c.NoContent(204)
}
