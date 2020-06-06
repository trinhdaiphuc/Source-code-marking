package comments

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *CommentHandler) GetComment(c echo.Context) (err error) {
	h.Logger.Debug("Get comment handler")
	// Get param
	commentID := c.Param("id")

	fileCollection := models.GetFileCollection(h.DB)
	resultFind := fileCollection.FindOne(context.Background(), bson.M{"comments._id": commentID})

	comment := models.Comment{}
	if err := resultFind.Decode(&comment); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found comment",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[GetComment] Internal server error",
			Internal: err,
		}
	}
	return c.JSON(http.StatusOK, comment)
}
