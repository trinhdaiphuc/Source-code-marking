package users

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Forbidden update user account",
		}
	}

	userCollection := models.GetUserCollection(h.DB)
	ctx := context.Background()
	// Get the old data of user
	result := userCollection.FindOne(ctx, bson.M{"_id": id})
	data := &models.User{}
	if err = result.Decode(&data); err != nil {
		h.Logger.Error("Error when get a user: ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found user",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}

	if data.IsDeleted {
		return &echo.HTTPError{
			Code:    http.StatusGone,
			Message: "User has been deleted.",
		}
	}

	update := bson.M{
		"$set": bson.M{
			"name":       u.Name,
			"updated_at": time.Now().UTC(),
		},
	}
	filter := bson.M{"_id": id}

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
		return err
	}

	c.Response().Header().Set("Access-Token", tokenString)
	data.Password = ""
	return c.JSON(http.StatusOK, data)
}
