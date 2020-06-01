package files

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *FileHandler) MarkFile(c echo.Context) (err error) {
	file := &models.File{}
	fileID := c.Param("id")
	if err := c.Bind(file); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid parameters",
			Internal: err,
		}
	}

	if !(file.Mark >= 0 && file.Mark <= 10) {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid parameters, mark must in [0,10]",
		}
	}

	ctx := context.Background()
	fileCollection := models.GetFileCollection(h.DB)

	h.Logger.Debug("Mark ", file.Mark)

	update := bson.M{
		"$set": bson.M{
			"mark":       file.Mark,
			"updated_at": time.Now().UTC(),
		},
	}
	filter := bson.M{"_id": fileID}

	resultUpdate := fileCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	err = resultUpdate.Decode(&file)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "Not found file",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}

	return c.NoContent(http.StatusNoContent)
}
