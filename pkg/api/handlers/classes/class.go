package classes

import (
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type ClassHandler struct {
	Logger *internal.AppLog
	DB     *mongo.Client
	JWTKey string
}

func NewClassHandler(logger *internal.AppLog, db *mongo.Client) (u *ClassHandler) {
	u = &ClassHandler{
		Logger: logger,
		DB:     db,
	}
	return
}
