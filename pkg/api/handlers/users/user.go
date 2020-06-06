package users

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"github.com/trinhdaiphuc/Source-code-marking/pkg/api/models"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Claims struct {
	ID    string `json:"id"`
	Role  string `json:"role"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type UserHandler struct {
	Logger      *internal.AppLog
	DB          *mongo.Client
	JWTKey      string
	RedisClient *redis.Client
}

func NewUserHandler(logger *internal.AppLog, jwtKey string, db *mongo.Client, redisClient *redis.Client) (u *UserHandler) {
	u = &UserHandler{
		Logger:      logger,
		JWTKey:      jwtKey,
		DB:          db,
		RedisClient: redisClient,
	}
	return
}

func createTokenWithUser(user models.User, jwtKey string, expireTime int) (string, error) {
	// Declare the expiration time of the token
	// here, we have kept it as 24 hours
	expirationTime := time.Now().Add(time.Duration(expireTime) * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		ID:    user.ID,
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	stringToken, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error when generate token: %v", err),
		)
	}
	return stringToken, nil
}
