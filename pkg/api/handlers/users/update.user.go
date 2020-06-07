package users

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *UserHandler) UpdateUser(c echo.Context) (err error) {
	h.Logger.Info("Sign-up handler")

	// Bind
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}

	id := c.Param("id")

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	userRole := claims["role"].(string)

	if userRole != "ADMIN" && id != userID {
		return echo.NewHTTPError(http.StatusForbidden, "Forbidden update user account")
	}

	userCollection := models.GetUserCollection(h.DB)
	ctx := context.Background()
	filter := bson.M{"_id": id}
	// Get the old data of user
	data, err := models.GetAUser(h.DB, filter, u.Role)

	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":       u.Name,
			"updated_at": time.Now().UTC(),
		},
	}

	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}

	// Generate encoded token and send it as response
	tokenString, err := createTokenWithUser(*data, h.JWTKey, 24)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	c.Response().Header().Set("Access-Token", tokenString)
	data.Password = ""
	return c.JSON(http.StatusOK, data)
}
