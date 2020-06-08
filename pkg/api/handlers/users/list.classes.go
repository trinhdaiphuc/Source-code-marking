package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) ListClasses(c echo.Context) (err error) {
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

	userID := c.Param("id")
	h.Logger.Debug(fmt.Sprintf("List query parameters: %v", *listParam))
	limit := listParam.PageSize
	page := listParam.PageToken
	skip := (page - 1) * limit
	orderBy := "created_at"
	orderType := 1
	if listParam.OrderType == internal.DESC.String() {
		orderType = -1
	}

	if listParam.OrderBy != "" {
		orderBy = listParam.OrderBy
	}

	ctx := context.Background()

	user, err := models.GetAUser(h.DB, bson.M{"_id": userID}, "")

	if err != nil {
		return err
	}

	opts := []*options.FindOptions{}
	opts = append(opts, options.Find().SetSort(bson.D{primitive.E{Key: orderBy, Value: orderType}}))
	opts = append(opts, options.Find().SetSkip(skip))
	opts = append(opts, options.Find().SetLimit(limit))

	filterBy := ""
	switch user.Role {
	case "STUDENT":
		filterBy = "students._id"
	case "TEACHER":
		filterBy = "teachers._id"
	}

	filter := bson.M{filterBy: userID, "is_deleted": false}
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

	totalRecords, err := classCollection.CountDocuments(ctx, filter)

	if cursor == nil {
		status.New(codes.FailedPrecondition, "No books have been created")
	}

	// An expression with defer will be called at the end of the function
	defer cursor.Close(ctx)

	classArray := []models.Class{}
	cursor.All(ctx, &classArray)
	return c.JSON(http.StatusOK, models.ConvertClassArrayToListClass(classArray, page+1, totalRecords))
}
