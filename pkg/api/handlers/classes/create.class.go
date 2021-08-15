package classes

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/middlewares"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *ClassHandler) CreateClass(c echo.Context) (err error) {
	classItem := &models.Class{}
	userContext := middlewares.GetUser(c)

	if err := c.Bind(classItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	if err := c.Validate(classItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	filter := bson.M{"_id": userContext.ID, "is_deleted": false}
	user, err := models.GetAUser(h.DB, filter, userContext.Role)
	if err != nil {
		return err
	}

	user.Password = ""
	classItem.ID = uuid.NewV4().String()
	classItem.Teachers = []models.User{*user}
	classItem.IsDeleted = false
	classItem.CreatedAt = time.Now().UTC()
	classItem.UpdatedAt = time.Now().UTC()

	classCollection := models.GetClassCollection(h.DB)
	_, err = classCollection.InsertOne(context.TODO(), classItem)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[CreateClass] Internal server error",
			Internal: err,
		}
	}
	return c.JSON(http.StatusCreated, classItem)
}
