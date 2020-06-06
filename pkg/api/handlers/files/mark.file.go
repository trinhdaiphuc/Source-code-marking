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

	classCollection := models.GetClassCollection(db)
	class, err := models.GetAClass(classCollection, bson.M{"_id": dataExercise.ClassID, "is_deleted": false})

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

	opts := []*options.FindOptions{}
	opts = append(opts, options.Find().SetSort(bson.D{{"created_at", 1}}))
	opts = append(opts, options.Find().SetSkip(0))
	opts = append(opts, options.Find().SetLimit(5))
	opts = append(opts, options.Find().SetProjection(bson.D{{"content", 1}, {"exercise_id", 1}, {"is_read", 1}}))

	filter := bson.M{"user_id": dataFile.UserID, "is_deleted": false}

	cursor, err := notificationCollection.Find(ctx, filter, opts...)
	if err != nil {
		logger.Error("Error when find ", err)
		return
	}
	notificationArray := []models.Notification{}
	cursor.All(ctx, &notificationArray)

	logger.Debug("Notifications ", notificationArray)

	message, _ := json.Marshal(notificationArray)
	err = redisClient.Publish(ctx, user.Email, message).Err()
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

	ctx := context.Background()
	fileCollection := models.GetFileCollection(h.DB)

	resultFind := fileCollection.FindOne(context.Background(), bson.M{"_id": fileID, "is_deleted": false})

	data := &models.File{}
	if err := resultFind.Decode(&data); err != nil {
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

	exerciseCollection := models.GetExerciseCollection(h.DB)
	resultFind = exerciseCollection.FindOne(context.Background(), bson.M{"_id": data.ExerciseID, "is_deleted": false})

	exercise := &models.Exercise{}
	if err := resultFind.Decode(&exercise); err != nil {
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusNotFound,
				Message:  "Not found Exercise",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[GetExercise] Internal server error",
			Internal: err,
		}
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
	go publishMarkingNotification(h.RedisClient, h.Logger, h.DB, *file, *exercise)
	return c.NoContent(http.StatusNoContent)
}
