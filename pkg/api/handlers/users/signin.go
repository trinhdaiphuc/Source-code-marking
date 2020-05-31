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

func (h *UserHandler) signinByEmail(u models.User, userCollection *mongo.Collection) (*models.User, error) {
	// Validate
	if len(u.Password) < 6 {
		return nil, &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Password length must be greater than 6",
		}
	}

	result := userCollection.FindOne(context.Background(), bson.M{"email": u.Email})

	user := &models.User{}
	if err := result.Decode(&user); err != nil {
		h.Logger.Info("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return nil, &echo.HTTPError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password.",
			}
		}
		return nil, &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "[Signin] Internal server error",
			Internal: err,
		}
	}

	if len(user.Password) == 0 {
		return nil, &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "You have already signed in with other services. Please use another sign-in method.",
		}
	}

	if ok := internal.CheckPasswordHash(u.Password, user.Password); !ok {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "Password is invalid."}
	}

	return user, nil
}

func (h *UserHandler) signinByThirdparty(u *models.User, userCollection *mongo.Collection) (*models.User, error) {
	if u.Password != "" {
		return nil, &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Please use right sign in method",
		}
	}

	result := userCollection.FindOne(context.Background(), bson.M{"email": u.Email})

	data := &models.User{}
	if err := result.Decode(&data); err != nil {
		if err == mongo.ErrNoDocuments {
			data = &models.User{
				ID:         uuid.NewV4().String(),
				Email:      u.Email,
				Name:       u.Name,
				IsVerified: true,
				Role:       u.Role,
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			}
			_, err = userCollection.InsertOne(context.Background(), data)
			if err != nil {
				// return internal gRPC error to be handled later
				return nil, &echo.HTTPError{
					Code:     http.StatusBadRequest,
					Message:  "[Signin] Internal server error",
					Internal: err,
				}
			}
			return data, nil
		}
		h.Logger.Error("Error when sign in by third-party ", err)
		return nil, &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "[Signin] Internal server error",
			Internal: err,
		}
	}

	return data, nil
}

// Signin handler
func (h *UserHandler) Signin(c echo.Context) (err error) {
	h.Logger.Debug("Sign-in handler")
	// Bind
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return &echo.HTTPError{
			Code:     http.StatusBadRequest,
			Message:  "Invalid arguments error",
			Internal: err,
		}
	}

	h.Logger.Debug("Sign-in parameters: ", *u)

	// Validate
	if u.Email == "" {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Missing parameters email",
		}
	}

	userCollection := models.GetUserCollection(h.DB)
	user := &models.User{}
	switch u.Service {
	case "EMAIL":
		if user, err = h.signinByEmail(*u, userCollection); err != nil {
			return err
		}
	case "GOOGLE":
		if user, err = h.signinByThirdparty(u, userCollection); err != nil {
			return err
		}
	case "FACEBOOK":
		if user, err = h.signinByThirdparty(u, userCollection); err != nil {
			return err
		}
	default:
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid services method."}
	}

	if err != nil {
		return err
	}

	// Generate encoded token and send it as response
	tokenString, err := createTokenWithUser(user.ID, user.Role, h.JWTKey)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Access-Token", tokenString)

	user.Password = "" // Don't send password

	return c.JSON(http.StatusOK, user)
}
