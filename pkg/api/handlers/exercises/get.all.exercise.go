package exercises

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *ExerciseHandler) GetAllExercises(c echo.Context) (err error) {
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

	var opts []*options.FindOptions
	opts = append(opts, options.Find().SetSort(bson.D{{orderBy, orderType}}))
	opts = append(opts, options.Find().SetSkip(skip))
	opts = append(opts, options.Find().SetLimit(limit))

	filter := bson.M{}
	ExerciseCollection := models.GetExerciseCollection(h.DB)
	ctx := context.Background()
	cursor, err := ExerciseCollection.Find(ctx, filter, opts...)
	if err != nil {
		h.Logger.Error("Internal error when Find: ", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Get all user] Internal server error",
			Internal: err,
		}
	}

	totalRecords, err := ExerciseCollection.CountDocuments(ctx, filter)

	if cursor == nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Not found exercises",
		}
	}

	// An expression with defer will be called at the end of the function
	defer cursor.Close(ctx)

	ExerciseArray := []models.Exercise{}
	cursor.All(ctx, &ExerciseArray)
	return c.JSON(http.StatusOK, models.ConvertExerciseArrayToListExercise(ExerciseArray, page+1, totalRecords))
}
