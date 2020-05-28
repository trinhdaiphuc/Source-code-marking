package comments

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *CommentHandler) UpdateComment(c echo.Context) (err error) {
	commentItem := &models.Comment{}
	if err = c.Bind(commentItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments",
			Internal: err,
		}
	}

	commentID := c.Param("id")

	h.Logger.Debug("Create comments parameters ", commentItem)

	fileCollection := models.GetFileCollection(h.DB)

	fileItem := &models.File{}
	ctx := context.Background()
	filter := bson.M{"_id": commentItem.FileID, "comments._id": commentID}

	update := bson.M{
		"$set": bson.M{
			"comments.$.content":    commentItem.Content,
			"comments.$.start_line": commentItem.StartLine,
			"comments.$.end_line":   commentItem.EndLine,
			"comments.$.updated_at": time.Now().UTC(),
		},
	}
	result := fileCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	if err := result.Decode(&fileItem); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[CreateComment] Internal server error",
			Internal: err,
		}
	}

	h.Logger.Debug("File Item ", fileItem)
	for _, v := range fileItem.Comments {
		if v.ID == commentID {
			commentItem = &v
			break
		}
	}
	return c.JSON(http.StatusCreated, commentItem)
}
