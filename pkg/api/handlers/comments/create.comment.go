package comments

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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *CommentHandler) CreateComment(c echo.Context) (err error) {
	commentItem := &models.Comment{}
	if err = c.Bind(commentItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments",
			Internal: err,
		}
	}

	h.Logger.Debug("Create comments parameters ", commentItem)

	commentItem.ID = uuid.NewV4().String()
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)
	fileCollection := models.GetFileCollection(h.DB)

	fileItem := &models.File{}
	ctx := context.Background()
	filter := bson.M{"_id": commentItem.FileID}
	result := fileCollection.FindOne(ctx, filter)

	if err := result.Decode(&fileItem); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found file.",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[CreateComment] Internal server error",
			Internal: err,
		}
	}

	commentItem.UserID = userID
	commentItem.CreatedAt = time.Now().UTC()
	commentItem.UpdatedAt = time.Now().UTC()

	update := bson.M{
		"$addToSet": bson.M{
			"comments": bson.M{
				"$each": []interface{}{commentItem},
			},
		},
	}
	resultUpdate := fileCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	if err := resultUpdate.Decode(&fileItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[CreateComment] Internal server error",
			Internal: err,
		}
	}

	h.Logger.Debug("File Item ", fileItem)

	return c.JSON(http.StatusCreated, commentItem)
}
