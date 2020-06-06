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
)

func (h *UserHandler) ChangePassword(c echo.Context) (err error) {
	password := &models.ChangePassword{}
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

	filter := bson.M{"_id": userID}
	userCollection := models.GetUserCollection(h.DB)
	resultFind := userCollection.FindOne(context.Background(), bson.M{"_id": userID})

	ctx := context.Background()
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

	if user.IsDeleted {
		return &echo.HTTPError{
			Code:    http.StatusGone,
			Message: "User has been deleted.",
		}
	}

	if ok := internal.CheckPasswordHash(password.OldPassword, user.Password); !ok {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid old password",
		}
	}

	newPassword, err := internal.HashPassword(password.NewPassword)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[ChangePassword] Internal server error",
			Internal: err,
		}
	}

	update := bson.M{
		"$set": bson.M{
			"password": newPassword,
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[ChangePassword] Internal server error",
			Internal: err,
		}
	}

	return c.NoContent(http.StatusNoContent)
}
