package handlers

import (
	"os"

	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	Handler struct {
		DB     *mongo.Client
		AppLog *internal.AppLog
	}
)

var Key string = os.Getenv("SECRET_KEY")
