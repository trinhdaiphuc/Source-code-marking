package handlers

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/comments"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/files"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	// Handler struct for handle all api logic
	Handler struct {
		DB             *mongo.Client
		Logger         *internal.AppLog
		JWTKey         string
		UserHandler    UserInterface
		FileHandler    FileInterface
		CommentHandler CommentInterface
	}

	// UserInterface is a interface for handle all user logic
	UserInterface interface {
		Signin(c echo.Context) (err error)
		Signup(c echo.Context) (err error)
		Profile(c echo.Context) (err error)
		GetUser(c echo.Context) (err error)
		GetAllUsers(c echo.Context) (err error)
		UpdateUser(c echo.Context) (err error)
		DeleteUser(c echo.Context) (err error)
	}

	// FileInterface is a interface for handle all file logic
	FileInterface interface {
		CreateFile(c echo.Context) (err error)
		UpdateFile(c echo.Context) (err error)
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
)

// NewUserHandlers create a handler pointer
func NewHandlers(db *mongo.Client, logger *internal.AppLog) (h *Handler) {
	h = &Handler{
		DB:     db,
		JWTKey: os.Getenv("SECRET_KEY"),
		Logger: logger,
	}

	h.UserHandler = users.NewUserHandler(logger, h.JWTKey, db)
	h.FileHandler = files.NewFileHandler(logger, db)
	h.CommentHandler = comments.NewCommentHandler(logger, db)

	return
}
