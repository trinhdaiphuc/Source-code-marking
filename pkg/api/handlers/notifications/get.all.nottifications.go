package notifications

import (
	"context"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *NotificationHandler) GetAllNotifications(c echo.Context) (err error) {
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

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	filter := bson.M{"user_id": userID}

	notificationCollection := models.GetNotificationCollection(h.DB)
	ctx := context.Background()
	cursor, err := notificationCollection.Find(ctx, filter, opts...)
	if err != nil {
		h.Logger.Error("Internal error when Find: ", err)
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Get all user] Internal server error",
			Internal: err,
		}
	}

	totalRecords, err := notificationCollection.CountDocuments(ctx, filter)

	if cursor == nil {
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: "Not found notifications",
		}
	}

	// An expression with defer will be called at the end of the function
	defer cursor.Close(ctx)

	notificationArray := []models.Notification{}
	cursor.All(ctx, &notificationArray)
	return c.JSON(http.StatusOK, models.ConvertNotificationArrayToListNotification(notificationArray, page+1, totalRecords))
}
