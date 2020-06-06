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
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ValidateUser handler
func (h *UserHandler) ValidateUser(c echo.Context) (err error) {
	h.Logger.Debug("Get user handler")
	// Get param
	token := c.QueryParam("confirmation_token")
	h.Logger.Debug("Token ", token)
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(h.JWTKey), nil
		})

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid validate token",
			Internal: err,
		}
	}
	if !tkn.Valid {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid validate token",
			Internal: err,
		}
	}

	userCollection := models.GetUserCollection(h.DB)
	ctx := context.Background()
	update := bson.M{
		"$set": bson.M{
			"is_verified": true,
			"updated_at":  time.Now().UTC(),
		},
	}
	resultFind := userCollection.FindOneAndUpdate(ctx, bson.M{"_id": claims.ID}, update, options.FindOneAndUpdate().SetReturnDocument(1))

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

	if user.IsDeleted {
		return &echo.HTTPError{
			Code:    http.StatusGone,
			Message: "User has been deleted.",
		}
	}

	return c.JSON(http.StatusOK, user)
}
