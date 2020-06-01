package users

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ValidateUser handler
func (h *UserHandler) ValidateUser(c echo.Context) (err error) {
	h.Logger.Debug("Get user handler")
	// Get param
	userID := c.Param("confirmation_token")

	userCollection := models.GetUserCollection(h.DB)
	resultFind := userCollection.FindOne(context.Background(), bson.M{"_id": userID})

		// // Initialize a new instance of `Claims`
		// claims := &Claims{}
		// tkn, err := jwt.ParseWithClaims(refreshToken, claims,
		// 	func(token *jwt.Token) (interface{}, error) {
		// 		return []byte(os.Getenv("SECRET_KEY")), nil
		// 	})
	user := models.User{}
	if err := resultFind.Decode(&user); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found user ",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Get user] Internal server error",
			Internal: err,
		}
	}
	user.Password = ""
	return c.JSON(http.StatusOK, user)
}
