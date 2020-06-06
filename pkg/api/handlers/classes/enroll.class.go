package classes

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *ClassHandler) EnrollClass(c echo.Context) (err error) {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	classID := c.Param("id")

	user := models.User{}
	ctx := context.Background()

	userCollection := models.GetUserCollection(h.DB)
	result := userCollection.FindOne(ctx, bson.M{"_id": userID})
	if err := result.Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found user",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Profile] Internal server error ",
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

	classCollection := models.GetClassCollection(h.DB)

	filter := bson.M{"_id": classID, "is_deleted": false}
	data := &models.Class{}
	update := bson.M{
		"$addToSet": bson.M{
			"students": bson.M{
				"$each": []models.User{user},
			},
		},
	}

	result = classCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	if err = result.Decode(&data); err != nil {
		if err != mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found class",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Enroll class] Internal server error",
			Internal: err,
		}
	}

	return c.JSON(http.StatusOK, data)
}
