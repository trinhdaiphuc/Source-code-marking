package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	uuid "github.com/satori/go.uuid"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gopkg.in/mgo.v2/bson"
)

type Claims struct {
	ID string `json:"id"`
	jwt.StandardClaims
}

func createTokenWithUserID(userID string, JwtKey string) (string, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	stringToken, err := token.SignedString([]byte(JwtKey))
	if err != nil {
		return "", status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error when generate token: %v", err),
		)
	}
	return stringToken, nil
}

func (h *Handler) signinByEmail(u *models.User, userColection *mongo.Collection) (err error) {
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
		h.AppLog.Info("Error when sign in by email ", err)
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

// Signup handler
func (h *Handler) Signup(c echo.Context) (err error) {
	h.AppLog.Info("Sign-up handler")

	// Bind
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	h.AppLog.Debug("Sign-up parameters: ", *u)
	// Validate
	if u.Email == "" || len(u.Password) < 6 {
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: "Invalid email or password",
		}
	}

	// Check email had created or not.
	userCollection := models.GetUserCollection(h.DB)
	resultFind := userCollection.FindOne(context.Background(), bson.M{"email": u.Email})

	user := models.User{}
	if err := resultFind.Decode(&user); err != nil {
		h.AppLog.Debug("Error when sign in by email ", err)
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
		h.AppLog.Error("Error when hashpassword ", err)
		return
	}

	u.ID = uuid.NewV4().String()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	h.AppLog.Debug("UUID ", u)

	// Save user
	_, err = userCollection.InsertOne(context.Background(), u)
	if err != nil {
		h.AppLog.Debug("Error when sign-up ", err.Error())
		return &echo.HTTPError{
			Code:     http.StatusInternalServerError,
			Message:  "MongoDB is not avalable.",
			Internal: err,
		}
	}
	u.Password = ""
	return c.JSON(http.StatusCreated, u)
}

// Signin handler
func (h *Handler) Signin(c echo.Context) (err error) {
	h.AppLog.Debug("Sign-in handler")
	// Bind
	u := &models.User{}
	if err = c.Bind(u); err != nil {
		return
	}
	h.AppLog.Debug("Sign-in parameters: ", *u)

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
	tokenString, err := createTokenWithUserID(u.ID, Key)
	if err != nil {
		return err
	}

	c.Response().Header().Set("token", tokenString)

	u.Password = "" // Don't send password
	return c.JSON(http.StatusOK, u)
}

// Profile handler
func (h *Handler) Profile(c echo.Context) (err error) {
	// Bind
	user := &models.User{}
	authHeader := strings.Split(c.Request().Header.Get("Authorization"), "Bearer ")
	jwtToken := authHeader[1]
	claims := &Claims{}
	_, err = jwt.ParseWithClaims(jwtToken, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(Key), nil
		})
	userCollection := models.GetUserCollection(h.DB)
	result := userCollection.FindOne(context.Background(), bson.M{"_id": claims.ID})
	if err := result.Decode(&user); err != nil {
		h.AppLog.Info("Error when sign in by email ", err)
		if err == mongo.ErrNoDocuments {
			return &echo.HTTPError{Code: http.StatusNotFound, Message: "Not found user"}
		}
		return &echo.HTTPError{Code: http.StatusInternalServerError,
			Message: "MongoDB is not avalable.",
		}
	}
	user.Password = ""
	c.JSON(http.StatusOK, user)
	return nil
}
