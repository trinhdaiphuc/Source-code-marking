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
	userRole := claims["role"].(string)

	user, err := models.GetAUser(h.DB, bson.M{"_id": userID}, userRole)
	if err != nil {
		return err
	}

	class, err := models.GetAClass(h.DB, bson.M{"_id": classID, "students._id": userID})

	if err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if code == http.StatusInternalServerError {
			return err
		}
	}

	if class != nil && class.ID != "" {
		return echo.NewHTTPError(http.StatusConflict, "You already enroll")
	}

	user.Password = ""

	filter := bson.M{"_id": classID, "is_deleted": false}
	data := &models.Class{}
	update := bson.M{
		"$addToSet": bson.M{
			"students": bson.M{
				"$each": []models.User{*user},
			},
		},
	}

	classCollection := models.GetClassCollection(h.DB)
	result := classCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

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
