package comments

import (
	"github.com/trinhdaiphuc/Source-code-marking/internal"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentHandler struct {
	Logger *internal.AppLog
	DB     *mongo.Client
}

func NewCommentHandler(logger *internal.AppLog, db *mongo.Client) (c *CommentHandler) {
	c = &CommentHandler{
		Logger: logger,
		DB:     db,
	}
	return
}
