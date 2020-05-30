package classes

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *ClassHandler) CreateClass(c echo.Context) (err error) {
	classItem := &models.Class{}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	userRole := claims["role"].(string)

	if userRole != "TEACHER" {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "Invalid role",
		}
	}

	if err := c.Bind(classItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

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
	classItem.ID = uuid.NewV4().String()
	classItem.Teachers = []models.User{user}
	classItem.CreatedAt = time.Now().UTC()
	classItem.UpdatedAt = time.Now().UTC()

	classCollection := models.GetClassCollection(h.DB)

	ctx := context.Background()
	_, err = classCollection.InsertOne(ctx, classItem)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[CreateClass] Internal server error",
			Internal: err,
		}
	}
	return c.JSON(http.StatusCreated, classItem)
}
