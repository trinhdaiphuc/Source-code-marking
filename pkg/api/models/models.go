package models

import "go.mongodb.org/mongo-driver/mongo"

type (
	ListQueryParam struct {
		FilterBy    string `query:"filter_by"`
		FilterValue string `query:"filter_value"`
		OrderBy     string `query:"order_by"`
		OrderType   string `query:"order_type"`
		PageToken   int64  `query:"page_token"`
		PageSize    int64  `query:"page_size"`
	}

	TotalRecords struct {
		ID    string `json:"id,omitempty" bson:"_id"`
		Total int64  `json:"total" bson:"total"`
	}
)

func NewModelDB(db *mongo.Client) {
	newUserCollection(db)
	newRoleCollection(db)
}
