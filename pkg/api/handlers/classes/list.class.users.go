package classes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *ClassHandler) ListClassUsers(c echo.Context) (err error) {
	listParam := &models.ListQueryParam{}
	if err = c.Bind(listParam); err != nil {
		return
	}
	classID := c.Param("id")
	h.Logger.Debug(fmt.Sprintf("List query parameters: %v", listParam))

	ctx := context.Background()
	filter := bson.M{"_id": classID}
	data := &models.Class{}
	classCollection := models.GetClassCollection(h.DB)
	result := classCollection.FindOne(ctx, filter)

	if err = result.Decode(&data); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found class",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[ListClassUser] Internal server error",
			Internal: err,
		}
	}

	userList := models.ListUser{}

	switch listParam.FilterValue {
	case "STUDENT":
		userList.Users = data.Students
		userList.TotalRecords = int64(len(data.Students))
	case "TEACHER":
		userList.Users = data.Teachers
		userList.TotalRecords = int64(len(data.Teachers))
	}

	return c.JSON(http.StatusOK, userList)
}
