package models

import "go.mongodb.org/mongo-driver/mongo"

type ListQueryParam struct {
	ListByRole string `query:"list_by_role"`
	OrderBy    string `query:"order_by"`
	OrderType  string `query:"order_type"`
	PageToken  int64  `query:"page_token"`
	PageSize   int64  `query:"page_size"`
}

func NewModelDB(db *mongo.Client) {
	newUserCollection(db)
	newRoleCollection(db)
}
