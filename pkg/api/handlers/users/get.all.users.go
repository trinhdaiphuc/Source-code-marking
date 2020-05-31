package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *UserHandler) GetAllUsers(c echo.Context) (err error) {
	listParam := &models.ListQueryParam{}
	if err = c.Bind(listParam); err != nil {
		return
	}
	h.Logger.Debug(fmt.Sprintf("List query parameters: %v", listParam))
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

	opts := []*options.FindOptions{}
	opts = append(opts, options.Find().SetSort(bson.D{{orderBy, orderType}}))
	opts = append(opts, options.Find().SetSkip(skip))
	opts = append(opts, options.Find().SetLimit(limit))

	filter := bson.M{}

	if listParam.FilterBy != "" && listParam.FilterValue != "" {
		filter = bson.M{listParam.FilterBy: listParam.FilterValue}
	}

	userCollection := models.GetUserCollection(h.DB)
	ctx := context.Background()
	cursor, err := userCollection.Find(ctx, filter, opts...)
	if err != nil {
		h.Logger.Error("Internal error when Find: ", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Get all user] Internal server error",
			Internal: err,
		}
	}

	totalRecords, err := userCollection.CountDocuments(ctx, filter)

	if cursor == nil {
		status.New(codes.FailedPrecondition, "No books have been created")
	}

	// An expression with defer will be called at the end of the function
	defer cursor.Close(ctx)

	userArray := []models.User{}
	cursor.All(ctx, &userArray)
	return c.JSON(http.StatusOK, models.ConvertUserArrayToListUser(userArray, page+1, totalRecords))
}
