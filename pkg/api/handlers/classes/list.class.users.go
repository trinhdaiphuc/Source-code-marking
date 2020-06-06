package classes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *ClassHandler) ListClassUsers(c echo.Context) (err error) {
	listParam := &models.ListQueryParam{}

	if err := c.Bind(listParam); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	if err := c.Validate(listParam); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	classID := c.Param("id")
	h.Logger.Debug(fmt.Sprintf("List query parameters: %v", listParam))
	limit := listParam.PageSize
	page := listParam.PageToken
	skip := (page - 1) * limit

	opts := []*options.FindOptions{}

	ctx := context.Background()
	filter := bson.M{
		"_id": classID,
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userRole := claims["role"].(string)

	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}

	switch listParam.FilterValue {
	case "STUDENT":
		opts = append(opts, options.Find().SetProjection(bson.D{
			{"students", bson.D{
				{"$slice", []interface{}{skip, limit}},
			}},
		}))
	case "TEACHER":
		opts = append(opts, options.Find().SetProjection(bson.D{
			{"teachers", bson.D{
				{"$slice", []interface{}{skip, limit}},
			}},
		}))
	}

	classCollection := models.GetClassCollection(h.DB)
	cursor, err := classCollection.Find(ctx, filter, opts...)

	if err != nil {
		h.Logger.Error("Internal error when Find: ", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Get all class] Internal server error",
			Internal: err,
		}
	}

	classArray := []models.Class{}
	cursor.All(ctx, &classArray)

	userList := classArray[0].Teachers

	return c.JSON(http.StatusOK, userList)
}
