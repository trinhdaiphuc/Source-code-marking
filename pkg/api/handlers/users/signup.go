package users

import (
	"context"
	"net/http"
	"os"
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
	u.IsVerified = false
	u.CreatedAt = time.Now().UTC()
	u.UpdatedAt = time.Now().UTC()

	roleCollection := models.GetRoleCollection(h.DB)
	result := roleCollection.FindOne(ctx, bson.M{"name": u.Role})

	role := &models.Role{}

	if err := result.Decode(&role); err != nil {
		h.Logger.Debug("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{
				Code:     http.StatusBadRequest,
				Message:  "Invalid role",
				Internal: err,
			}
		}
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Signup] Internal server error",
			Internal: err,
		}
	}

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
	go func() {
		token, _ := createTokenWithUser(u.ID, u.Role, h.JWTKey, 24)
		validationLink := os.Getenv("FRONT_END_SERVER_HOST") + "/confirmation?confirmation_token=" + token
		h.Logger.Info("Validation link ", validationLink)
		content := "Please click this link to verify your email: " + validationLink
		subject := "Welcome to Source code marking"
		id, err := internal.SendMail(os.Getenv("EMAIL_USERNAME"), u.Email, subject, content)
		if err != nil {
			h.Logger.Error("Error when send mail ", err)
		} else {
			h.Logger.Info("Send mail success with id ", id)
		}
	}()
	return c.JSON(http.StatusCreated, u)
}
