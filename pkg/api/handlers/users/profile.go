package users

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Profile handler
func (h *UserHandler) Profile(c echo.Context) (err error) {
	// Bind
	user := &models.User{}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	h.Logger.Debug("User id ", userID)
	userCollection := models.GetUserCollection(h.DB)
	result := userCollection.FindOne(context.Background(), bson.M{"_id": userID})
	if err := result.Decode(&user); err != nil {
		h.Logger.Info("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found user %v",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Profile] Internal server error ",
			Internal: err,
		}
	}
	user.Password = ""

	return c.JSON(http.StatusOK, user)
}
