package users

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
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

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	userRole := claims["role"].(string)

	filter := bson.M{"_id": userID}
	user, err := models.GetAUser(h.DB, filter, userRole)

	if err != nil {
		return err
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

	userCollection := models.GetUserCollection(h.DB)
	_, err = userCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[ChangePassword] Internal server error",
			Internal: err,
		}
	}

	return c.NoContent(http.StatusNoContent)
}
