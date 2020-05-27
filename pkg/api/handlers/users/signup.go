package users

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Signup handler
func (h *UserHandler) Signup(c echo.Context) (err error) {
	h.Logger.Info("Sign-up handler")

	// Bind
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	h.Logger.Debug("Sign-up parameters: ", *u)
	// Validate
	if u.Email == "" || len(u.Password) < 6 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid email or password",
		}
	}

	ctx := context.Background()
	// Check email had created or not.
	userCollection := models.GetUserCollection(h.DB)
	resultFind := userCollection.FindOne(ctx, bson.M{"email": u.Email})

	user := models.User{}
	if err := resultFind.Decode(&user); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err != mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:    http.StatusInternalServerError,
				Message: "MongoDB is not avalable.",
			}
		}
	}

	if user.Email != "" {
		return &echo.HTTPError{
			Code:    http.StatusConflict,
			Message: "This email have already existed.",
		}
	}

	// Hash password
	u.Password, err = internal.HashPassword(u.Password)
	if err != nil {
		h.Logger.Error("Error when hashpassword ", err)
		return
	}

	u.ID = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	roleCollection := models.GetRoleCollection(h.DB)
	result := roleCollection.FindOne(ctx, bson.M{"name": u.Role})

	role := &models.Role{}

	if err := result.Decode(&role); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: "Invalid role",
			}
		}
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: "MongoDB is not avalable.",
		}
	}

	h.Logger.Debug("UUID ", u)

	// Save user
	_, err = userCollection.InsertOne(context.Background(), u)
	if err != nil {
		h.Logger.Debug("Error when sign-up ", err.Error())
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "MongoDB is not avalable.",
			Internal: err,
		}
	}
	u.Password = ""
	return c.JSON(http.StatusCreated, u)
}
