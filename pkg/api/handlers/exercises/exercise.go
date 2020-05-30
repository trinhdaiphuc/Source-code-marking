package exercises

import (
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type ExerciseHandler struct {
	Logger *internal.AppLog
	DB     *mongo.Client
}

func NewExerciseHandler(logger *internal.AppLog, db *mongo.Client) (u *ExerciseHandler) {
	u = &ExerciseHandler{
		Logger: logger,
		DB:     db,
	}
	return
}
