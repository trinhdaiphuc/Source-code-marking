package handlers

import (
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/classes"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/comments"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/exercises"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/files"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/notifications"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Handler struct for handle all api logic
	Handler struct {
		DB               *mongo.Client
		Logger           *internal.AppLog
		JWTKey           string
		RedisClient      *redis.Client
		WebsocketClients map[*websocket.Conn]string
		UserHandler      UserInterface
		FileHandler      FileInterface
		CommentHandler   CommentInterface
		ClassHandler     ClassInterface
		ExerciseHandler  ExercisesInterface
		Notification     NotificationInterface
	}

	// UserInterface is a interface for handle all user logic
	UserInterface interface {
		Signin(c echo.Context) (err error)
		Signup(c echo.Context) (err error)
		Profile(c echo.Context) (err error)
		CreateUser(c echo.Context) (err error)
		GetUser(c echo.Context) (err error)
		GetAllUsers(c echo.Context) (err error)
		UpdateUser(c echo.Context) (err error)
		DeleteUser(c echo.Context) (err error)
		ListClasses(c echo.Context) (err error)
		ValidateUser(c echo.Context) (err error)
		ForgetPassword(c echo.Context) (err error)
		ResetPassword(c echo.Context) (err error)
		ChangePassword(c echo.Context) (err error)
	}

	// ClassInterface is a interface for handle all class logic
	ClassInterface interface {
		CreateClass(c echo.Context) (err error)
		UpdateClass(c echo.Context) (err error)
		GetClass(c echo.Context) (err error)
		GetAllClasses(c echo.Context) (err error)
		DeleteClass(c echo.Context) (err error)
		EnrollClass(c echo.Context) (err error)
		UnenrollClass(c echo.Context) (err error)
		ListExercises(c echo.Context) (err error)
		ListClassUsers(c echo.Context) (err error)
	}

	// ExercisesInterface is a interface for handle all class logic
	ExercisesInterface interface {
		CreateExercise(c echo.Context) (err error)
		UpdateExercise(c echo.Context) (err error)
		GetExercise(c echo.Context) (err error)
		GetAllExercises(c echo.Context) (err error)
		DeleteExercise(c echo.Context) (err error)
		ListFiles(c echo.Context) (err error)
	}

	// FileInterface is a interface for handle all file logic
	FileInterface interface {
		CreateFile(c echo.Context) (err error)
		UpdateFile(c echo.Context) (err error)
		MarkFile(c echo.Context) (err error)
		GetFile(c echo.Context) (err error)
		GetAllFiles(c echo.Context) (err error)
		DeleteFile(c echo.Context) (err error)
		ListComments(c echo.Context) (err error)
	}

	// CommentInterface is a interface for handle all comment logic
	CommentInterface interface {
		CreateComment(c echo.Context) (err error)
		UpdateComment(c echo.Context) (err error)
		GetComment(c echo.Context) (err error)
		GetAllComments(c echo.Context) (err error)
		DeleteComment(c echo.Context) (err error)
	}

	// NotificationInterface is a interface for handle all notification logic
	NotificationInterface interface {
		WebsocketNotification(c echo.Context) (err error)
		GetAllNotifications(c echo.Context) (err error)
		MarkReadNotification(c echo.Context) (err error)
	}
)

// NewHandlers create a handler pointer
func NewHandlers(db *mongo.Client, logger *internal.AppLog, redisClient *redis.Client) (h *Handler) {
	h = &Handler{
		DB:          db,
		JWTKey:      os.Getenv("SECRET_KEY"),
		Logger:      logger,
		RedisClient: redisClient,
	}

	h.WebsocketClients = make(map[*websocket.Conn]string)

	h.UserHandler = users.NewUserHandler(logger, h.JWTKey, db, h.RedisClient)
	h.FileHandler = files.NewFileHandler(logger, db, h.RedisClient)
	h.CommentHandler = comments.NewCommentHandler(logger, db)
	h.ClassHandler = classes.NewClassHandler(logger, db)
	h.ExerciseHandler = exercises.NewExerciseHandler(logger, db, h.RedisClient)
	h.Notification = notifications.NewNotificationHandler(logger, db, h.RedisClient, h.WebsocketClients)

	return
}
