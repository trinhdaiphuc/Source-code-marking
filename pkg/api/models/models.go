package models

import (
	"os"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	ListQueryParam struct {
		FilterBy    string `query:"filter_by"`
		FilterValue string `query:"filter_value"`
		OrderBy     string `query:"order_by"`
		OrderType   string `query:"order_type"`
		PageToken   int64  `query:"page_token" validate:"gte=1"`
		PageSize    int64  `query:"page_size" validate:"gte=10,lte=50"`
	}

	TotalRecords struct {
		ID    string `json:"id,omitempty" bson:"_id"`
		Total int64  `json:"total" bson:"total"`
	}

	CustomValidator struct {
		validator *validator.Validate
	}
)

func NewModelDB(db *mongo.Client) {
	newUserCollection(db)
	newRoleCollection(db)
}

func getDatabase(db *mongo.Client) (database *mongo.Database) {
	return db.Database(os.Getenv("DATABASE_NAME"))
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}
