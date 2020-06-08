package users

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func deleteUserClasss(db *mongo.Client, logger *internal.AppLog, userID string) {
	classCollection := models.GetClassCollection(db)

	userItem, err := models.GetAUser(db, bson.M{"_id": userID}, "")
	if err != nil {
		logger.Error("Error when get user ")
		return
	}

	update := bson.M{}
	filter := bson.M{}

	switch userItem.Role {
	case "TEACHER":
		filter["teachers._id"] = userID
		update = bson.M{
			"$pull": bson.M{
				"teachers": bson.M{
					"_id": userID,
				},
			},
		}
	case "STUDENT":
		filter["students._id"] = userID
		update = bson.M{
			"$pull": bson.M{
				"students": bson.M{
					"_id": userID,
				},
			},
		}
	default:
		filter["_id"] = "12"
	}

	ctx := context.Background()
	_, err = classCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		logger.Error("Error when delete user in class ", err)
	}
}

func (h *UserHandler) DeleteUser(c echo.Context) (err error) {
	id := c.Param("id")

	filter := bson.M{"_id": id}
	userCollection := models.GetUserCollection(h.DB)
	update := bson.M{
		"$set": bson.M{
			"is_deleted": true,
			"updated_at": time.Now().UTC(),
		},
	}
	_, err = userCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Delete user] Internal server error",
			Internal: err,
		}
	}
	go deleteUserClasss(h.DB, h.Logger, id)
	return c.NoContent(http.StatusNoContent)
}
