package files

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func publishMarkingNotification(
	redisClient *redis.Client, logger *internal.AppLog, db *mongo.Client,
	dataFile models.File, dataExercise models.Exercise) {
	ctx := context.Background()
	userCollection := models.GetUserCollection(db)
	user := &models.User{}
	result := userCollection.FindOne(ctx, bson.M{"_id": dataFile.UserID})

	if err := result.Decode(&user); err != nil {
		logger.Error("Error when get user ", err)
		return
	}

	class, err := models.GetAClass(db, bson.M{"_id": dataExercise.ClassID, "is_deleted": false})
	if err != nil {
		logger.Error("Error when get a class ", err)
	}

	notification := &models.Notification{
		ID:         uuid.NewV4().String(),
		Content:    "Bài tập " + dataExercise.Name + " của lớp " + class.Name + " đã được chấm",
		IsRead:     false,
		IsDeleted:  false,
		ExerciseID: dataFile.ExerciseID,
		UserID:     dataFile.UserID,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	notificationCollection := models.GetNotificationCollection(db)
	notificationCollection.InsertOne(ctx, notification)

	listParam := models.ListQueryParam{
		PageSize:  5,
		PageToken: 1,
		OrderBy:   "created_at",
		OrderType: internal.DESC.String(),
	}

	filter := bson.M{"user_id": dataFile.UserID, "is_deleted": false}

	listNotification, err := models.ListAllNotifications(db, filter, listParam)
	if err != nil {
		logger.Error("[Mark file] Error when find ", err)
		return
	}

	totalUnread, err := notificationCollection.CountDocuments(context.TODO(), bson.M{"is_read": false})
	listNotificationWebsocket := models.ListNotificationWebsocket{
		Notifications: listNotification.Notifications,
		TotalUnread:   totalUnread,
	}
	message, _ := json.Marshal(listNotificationWebsocket)
	err = redisClient.Publish(ctx, user.Email, message).Err()

	logger.Debug("User Email ", user.Email)

	if err != nil {
		logger.Error("Error when publish ", err)
	}
	return
}

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

	fileItem, err := models.GetAFile(h.DB, bson.M{"_id": fileID, "is_deleted": false})
	if err != nil {
		return err
	}

	exercise, err := models.GetAExercise(h.DB, bson.M{"_id": fileItem.ExerciseID, "is_deleted": false})
	if err != nil {
		return err
	}

	if exercise.Deadline.Sub(time.Now()) > 0 {
		return &echo.HTTPError{
			Code:     http.StatusTooEarly,
			Message:  "It has not been over deadline",
			Internal: err,
		}
	}

	update := bson.M{
		"$set": bson.M{
			"mark":       file.Mark,
			"updated_at": time.Now().UTC(),
		},
	}
	filter := bson.M{"_id": fileID}

	fileCollection := models.GetFileCollection(h.DB)
	resultUpdate := fileCollection.FindOneAndUpdate(context.TODO(), filter, update, options.FindOneAndUpdate().SetReturnDocument(1))
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
	go publishMarkingNotification(h.RedisClient, h.Logger, h.DB, *file, *exercise)
	return c.NoContent(http.StatusNoContent)
}
