package handlers

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/handlers/users"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Handler struct {
		DB          *mongo.Client
		Logger      *internal.AppLog
		JWTKey      string
		UserHandler UserInterface
	}

	UserInterface interface {
		Signin(c echo.Context) (err error)
		Signup(c echo.Context) (err error)
		Profile(c echo.Context) (err error)
	}
)

func NewUserHandlers(db *mongo.Client, logger *internal.AppLog) (h *Handler) {
	h = &Handler{}
	h.DB = db
	h.JWTKey = os.Getenv("SECRET_KEY")
	h.Logger = logger
	h.UserHandler = users.NewUserHandler(logger, h.JWTKey, db)
	return
}

func (h *Handler) Signup(c echo.Context) (err error) {
	return h.UserHandler.Signup(c)
}

func (h *Handler) Signin(c echo.Context) (err error) {
	return h.UserHandler.Signin(c)
}

func (h *Handler) Profile(c echo.Context) (err error) {
	return h.UserHandler.Profile(c)
}
