package files

import (
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileHandler struct {
	Logger *internal.AppLog
	DB     *mongo.Client
}

func NewFileHandler(logger *internal.AppLog, db *mongo.Client) (u *FileHandler) {
	u = &FileHandler{
		Logger: logger,
		DB:     db,
	}
	return
}
