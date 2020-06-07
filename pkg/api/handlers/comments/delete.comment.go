package comments

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *CommentHandler) DeleteComment(c echo.Context) (err error) {
	commentID := c.Param("id")

	fileItem := &models.File{}
	ctx := context.Background()
	filter := bson.M{"comments._id": commentID}
	update := bson.M{
		"$pull": bson.M{
			"comments": bson.M{
				"_id": commentID,
			},
		},
	}

	fileCollection := models.GetFileCollection(h.DB)
	result := fileCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	err = result.Decode(&fileItem)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found comment",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "Internal server error",
			Internal: err,
		}
	}
	return c.NoContent(http.StatusNoContent)
}
