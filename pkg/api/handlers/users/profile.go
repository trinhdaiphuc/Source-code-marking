package users

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

// Profile handler
func (h *UserHandler) Profile(c echo.Context) (err error) {
	h.Logger.Debug("Profile user handler")
	user := &models.User{}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	key := "user:" + userID
	cached, err := internal.RedisGetCachedWithHash(key, h.RedisClient)

	if err != nil {
		h.Logger.Error("Error when get cache ", err)
		goto FIND_DB
	}

	mapstructure.Decode(cached, &user)

	if user.ID != "" {
		user.CreatedAt, err = time.Parse(time.RFC3339, cached["CreatedAt"])
		user.UpdatedAt, err = time.Parse(time.RFC3339, cached["UpdatedAt"])
		h.Logger.Debug("User get from cached ", user)
		return c.JSON(http.StatusOK, user)
	}

FIND_DB:
	userRole := claims["role"].(string)
	user, err = models.GetAUser(h.DB, bson.M{"_id": userID}, userRole)

	if err != nil {
		return err
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
