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
	userCollection := models.GetUserCollection(h.DB)
	result := userCollection.FindOne(context.Background(), bson.M{"_id": userID})
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

	classCollection := models.GetClassCollection(h.DB)

	update := bson.M{
		"$addToSet": bson.M{
			"students": bson.M{
				"$each": []models.User{user},
			},
		},
	}
	filter := bson.M{"_id": classID, "students._id": user.ID}
	data := &models.Class{}
	ctx := context.Background()

	result = classCollection.FindOne(ctx, filter)
	err = result.Decode(&data)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found class",
				Internal: err,
			}
		}
	}
	if data.ID != "" {
		return &echo.HTTPError{
			Code:     http.StatusConflict,
			Message:  "Already enroll",
			Internal: err,
		}
	}

	filter = bson.M{"_id": classID}
	result = classCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	err = result.Decode(&data)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found class",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}
	return c.JSON(http.StatusOK, data)
}
