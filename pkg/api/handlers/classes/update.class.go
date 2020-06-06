package classes

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

func (h *ClassHandler) UpdateClass(c echo.Context) (err error) {
	h.Logger.Debug("Get user handler")
	// Get param
	class := &models.Class{}
	classID := c.Param("id")
	if err := c.Bind(class); err != nil {
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
	classCollection := models.GetClassCollection(h.DB)
	resultFind := classCollection.FindOne(ctx, bson.M{"_id": classID})

	data := models.Class{}
	if err := resultFind.Decode(&data); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "Not found class",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[UpdateClass] Internal server error",
			Internal: err,
		}
	}
	userRole := claims["role"].(string)

	if userRole == "ADMIN" {
		goto UPDATECLASS
	}

	for _, v := range data.Teachers {
		if v.ID == userID {
			goto UPDATECLASS
		}
	}

	return &echo.HTTPError{
		Code:    http.StatusForbidden,
		Message: "User cannot update this class.",
	}

UPDATECLASS:

	update := bson.M{
		"$set": bson.M{
			"name":        class.Name,
			"description": class.Description,
			"updated_at":  time.Now().UTC(),
		},
	}
	filter := bson.M{"_id": classID}

	if userRole != "ADMIN" {
		filter["is_deleted"] = false
	}

	resultUpdate := classCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
	err = resultUpdate.Decode(&data)
	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Update user] Internal server error",
			Internal: err,
		}
	}

	if class.Teachers != nil {
		update = bson.M{
			"$addToSet": bson.M{
				"teachers": bson.M{
					"$each": class.Teachers,
				},
			},
		}
		filter = bson.M{"_id": classID}

		resultUpdate = classCollection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
		err = resultUpdate.Decode(&data)
		if err != nil {
			return &echo.HTTPError{
				Code:     http.StatusInternalServerError,
				Message:  "[Update user] Internal server error",
				Internal: err,
			}
		}
	}

	return c.JSON(http.StatusOK, data)
}
