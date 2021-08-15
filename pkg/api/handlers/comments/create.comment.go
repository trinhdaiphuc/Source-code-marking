package comments

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
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

	if err := c.Validate(commentItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	commentItem.ID = uuid.NewV4().String()
	commentItem.IsDeleted = false
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	filter := bson.M{"_id": commentItem.FileID, "is_deleted": false}
	fileItem, err := models.GetAFile(h.DB, filter)
	if err != nil {
		return err
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
	fileCollection := models.GetFileCollection(h.DB)
	resultUpdate := fileCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	if err := resultUpdate.Decode(&fileItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[CreateComment] Internal server error",
			Internal: err,
		}
	}

	return c.JSON(http.StatusCreated, commentItem)
}
