package files

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *FileHandler) GetFile(c echo.Context) (err error) {
	h.Logger.Debug("Get file handler")
	// Get param
	fileID := c.Param("id")

	userCollection := models.GetFileCollection(h.DB)
	resultFind := userCollection.FindOne(context.Background(), bson.M{"_id": fileID})

	file := &models.File{}
	if err := resultFind.Decode(&file); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found file",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[GetFile] Internal server error",
			Internal: err,
		}
	}
	return c.JSON(http.StatusOK, file)
}
