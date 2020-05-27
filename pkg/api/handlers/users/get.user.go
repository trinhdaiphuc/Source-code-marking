package users

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Signin handler
func (h *UserHandler) GetUser(c echo.Context) (err error) {
	h.Logger.Debug("Sign-in handler")
	// Get param
	userID := c.Param("id")

	userCollection := models.GetUserCollection(h.DB)
	resultFind := userCollection.FindOne(context.Background(), bson.M{"_id": userID})

	user := models.User{}
	if err := resultFind.Decode(&user); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err != mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "MongoDB is not avalable.",
			}
		}
	}
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}
