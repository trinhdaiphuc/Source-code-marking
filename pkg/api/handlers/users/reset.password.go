package users

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *UserHandler) ResetPassword(c echo.Context) (err error) {
	h.Logger.Debug("Reset password handler")
	password := &models.ResetPassword{}
	if err := c.Bind(password); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	if err := c.Validate(password); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	user := &models.User{}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	h.Logger.Debug("Password: ", password.Password)
	filter := bson.M{"_id": userID, "is_deleted": false}
	ctx := context.TODO()
	hashPassword, err := internal.HashPassword(password.Password)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[ResetPassword] Internal server error",
			Internal: err,
		}
	}

	update := bson.M{
		"$set": bson.M{
			"password": hashPassword,
		},
	}

	userCollection := models.GetUserCollection(h.DB)
	result := userCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	if err = result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found User",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[ResetPassword] Internal server error",
			Internal: err,
		}
	}

	if user.IsDeleted {
		return &echo.HTTPError{
			Code:    http.StatusGone,
			Message: "User has been deleted.",
		}
	}

	return c.NoContent(http.StatusNoContent)
}
