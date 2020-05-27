package users

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Claims struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	jwt.StandardClaims
}

type UserHandler struct {
	Logger *internal.AppLog
	DB     *mongo.Client
	JWTKey string
}

func NewUserHandler(logger *internal.AppLog, jwtKey string, db *mongo.Client) (u *UserHandler) {
	u = &UserHandler{
		Logger: logger,
		JWTKey: jwtKey,
		DB:     db,
	}
	return
}

func createTokenWithUser(userID, role, JwtKey string) (string, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 24 hours
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID:   userID,
		Role: role,
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
