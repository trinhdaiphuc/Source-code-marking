package users

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
)

func (h *UserHandler) CreateUser(c echo.Context) (err error) {
	// Bind
	u := &models.User{}

	if err := c.Bind(u); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	if err := c.Validate(u); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	h.Logger.Debug("Sign-up parameters: ", *u)
	// Validate
	if u.Email == "" || len(u.Password) < 6 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid email or password",
		}
	}

	// Check email had created or not.
	user, err := models.GetAUser(h.DB, bson.M{"email": u.Email}, u.Role)

	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	if code == http.StatusInternalServerError {
		return err
	}

	if user.Email != "" {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "This email have already existed.",
		}
	}

	if !(u.Role == "STUDENT" || u.Role == "TEACHER") {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid arguments: role",
		}
	}

	// Hash password
	u.Password, err = internal.HashPassword(u.Password)
	if err != nil {
		h.Logger.Error("Error when hashpassword ", err)
		return
	}

	u.ID = uuid.NewV4().String()
	u.IsVerified = false
	u.IsDeleted = false
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()

	_, err = models.GetARole(h.DB, bson.M{"name": u.Role})

	if err != nil {
		return err
	}

	// Save user
	if err = models.CreateAUser(h.DB, u); err != nil {
		return err
	}

	u.Password = ""
	go sendValidationMail(*u, h.JWTKey, h.Logger)
	return c.JSON(http.StatusCreated, u)
}
