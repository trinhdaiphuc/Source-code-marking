package users

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (h *UserHandler) signinByEmail(u models.User, db *mongo.Client) (*models.User, error) {
	// Validate
	if len(u.Password) < 6 {
		return nil, &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Password length must be greater than 6",
		}
	}

	user, err := models.GetAUser(db, bson.M{"email": u.Email}, u.Role)

	if err != nil {
		return nil, err
	}

	if !user.IsVerified {
		return nil, &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Email has not verified.",
		}
	}

	if len(user.Password) == 0 {
		return nil, &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Sign in with wrong method.",
		}
	}

	if ok := internal.CheckPasswordHash(u.Password, user.Password); !ok {
		return nil, &echo.HTTPError{
			Code:    http.StatusUnauthorized,
			Message: "Password is invalid.",
		}
	}

	return user, nil
}

func (h *UserHandler) signinByThirdparty(u *models.User, db *mongo.Client) (*models.User, error) {
	if u.Password != "" {
		return nil, &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Please use right sign in method",
		}
	}

	if !(u.Role == "STUDENT" || u.Role == "TEACHER") {
		return nil, &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid arguments: role",
		}
	}

	data, err := models.GetAUser(db, bson.M{"email": u.Email}, u.Role)
	if err != nil {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}

		if code == http.StatusNotFound {
			data = &models.User{
				ID:         uuid.NewV4().String(),
				Email:      u.Email,
				Name:       u.Name,
				IsVerified: true,
				Role:       u.Role,
				IsDeleted:  false,
				CreatedAt:  time.Now().UTC(),
				UpdatedAt:  time.Now().UTC(),
			}
			err := models.CreateAUser(db, data)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
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

	user := &models.User{}
	switch u.Service {
	case "EMAIL":
		if user, err = h.signinByEmail(*u, h.DB); err != nil {
			return err
		}
	case "GOOGLE":
		if user, err = h.signinByThirdparty(u, h.DB); err != nil {
			return err
		}
	case "FACEBOOK":
		if user, err = h.signinByThirdparty(u, h.DB); err != nil {
			return err
		}
	default:
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "Invalid services method."}
	}

	// Generate encoded token and send it as response
	tokenString, err := createTokenWithUser(*user, h.JWTKey, 24)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Access-Token", tokenString)

	user.Password = "" // Don't send password

	go func() {
		key := "user:" + user.ID
		err = internal.RedisSetCachedWithHash(key, h.RedisClient, user)
		if err != nil {
			h.Logger.Error("Error when cached user ", err)
		}
	}()

	return c.JSON(http.StatusOK, user)
}
