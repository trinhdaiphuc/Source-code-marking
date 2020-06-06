package users

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Signin handler
func (h *UserHandler) GetUser(c echo.Context) (err error) {
	h.Logger.Debug("Get user handler")
	// Get param
	userID := c.Param("id")

	key := "user:" + userID
	cached, err := internal.RedisGetCachedWithHash(key, h.RedisClient)

	if err != nil {
		h.Logger.Error("Error when get cache ", err)
	}
	h.Logger.Debug("Cached ", cached)
	user := models.User{}
	mapstructure.Decode(cached, &user)

	if user.ID != "" {
		user.CreatedAt, err = time.Parse(time.RFC3339, cached["CreatedAt"])
		user.UpdatedAt, err = time.Parse(time.RFC3339, cached["UpdatedAt"])
		h.Logger.Debug("User get from cached ", user)
		return c.JSON(http.StatusOK, user)
	}

	userCollection := models.GetUserCollection(h.DB)
	resultFind := userCollection.FindOne(context.Background(), bson.M{"_id": userID})

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

	user.Password = ""
	go func() {
		key := "user:" + user.ID
		err = internal.RedisSetCachedWithHash(key, h.RedisClient, user)
		if err != nil {
			h.Logger.Error("Error when cached user ", err)
		}
	}()

	return c.JSON(http.StatusOK, user)
}
