package files

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h *FileHandler) UpdateFile(c echo.Context) (err error) {
	h.Logger.Debug("Get user handler")
	// Get param
	file := &models.File{}
	fileID := c.Param("id")
	if err := c.Bind(file); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid parameters",
			Internal: err,
		}
	}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["id"].(string)

	ctx := context.Background()
	fileCollection := models.GetFileCollection(h.DB)
	resultFind := fileCollection.FindOne(ctx, bson.M{"_id": fileID})

	data := models.File{}
	if err := resultFind.Decode(&data); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "Not found file",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[UpdateFile] Internal server error",
			Internal: err,
		}
	}

	if userID != data.UserID {
		return &echo.HTTPError{
			Code:    http.StatusForbidden,
			Message: "User cannot update this file.",
		}
	}

	update := bson.M{
		"$set": bson.M{
			"name":       file.Name,
			"data":       file.Data,
			"updated_at": time.Now().UTC(),
		},
	}
	filter := bson.M{"_id": fileID}

	resultUpdate := fileCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	err = resultUpdate.Decode(&data)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}

	return c.JSON(http.StatusOK, data)
}
