package users

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *UserHandler) signinByEmail(u *models.User, userColection *mongo.Collection) (err error) {
	// Validate
	if len(u.Password) < 6 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Password length must be greater than 6",
		}
	}

	result := userColection.FindOne(context.Background(), bson.M{"email": u.Email})

	user := models.User{}
	if err := result.Decode(&user); err != nil {
		h.Logger.Info("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Invalid email or password."}
		}
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message: "MongoDB is not avalable.",
		}
	}

	if len(user.Password) == 0 {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "You have already signed in with other services. Please use another sign-in method.",
		}
	}

	if ok := internal.CheckPasswordHash(u.Password, user.Password); !ok {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Password is invalid."}
	}

	u.ID = user.ID
	if err != nil {
		return &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Could not convert to ObjectId",
		}
	}
	u.Name = user.Name
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	return nil
}

// Signin handler
func (h *UserHandler) Signin(c echo.Context) (err error) {
	h.Logger.Debug("Sign-in handler")
	// Bind
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	h.Logger.Debug("Sign-in parameters: ", *u)

	// Validate
	if u.Email == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing parameters",
		}
	}

	userCollection := models.GetUserCollection(h.DB)
	if err := h.signinByEmail(u, userCollection); err != nil {
		return err
	}

	// Generate encoded token and send it as response
	tokenString, err := createTokenWithUserID(u.ID, h.JWTKey)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Access-Token", tokenString)

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}
