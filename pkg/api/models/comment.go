package models

import (
	"time"
)

type (
	Comment struct {
		ID        string    `json:"id,omitempty" bson:"_id"`
		FileID    string    `json:"file_id" bson:"file_id" validate:"required"`
		UserID    string    `json:"user_id" bson:"user_id"`
		Content   string    `json:"content" bson:"content" validate:"required"`
		StartLine Position  `json:"start_line,omitempty" bson:"start_line,omitempty"`
		EndLine   Position  `json:"end_line,omitempty" bson:"end_line,omitempty"`
		IsDeleted bool      `json:"is_deleted" bson:"is_deleted"`
		CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at"`
		UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at"`
	}

	Position struct {
		Row    int `json:"row" bson:"row"`
		Column int `json:"column" bson:"column"`
	}

	ListComment struct {
		Comments      []Comment `json:"comments"`
		NextPageToken int64     `json:"next_page_token"`
		TotalRecords  int64     `json:"total_records"`
	}
)

func ConvertCommentArrayToListComment(Comments []Comment, nextPageToken, totalRecords int64) *ListComment {
	listComment := &ListComment{
		Comments:      Comments,
		NextPageToken: nextPageToken,
		TotalRecords:  totalRecords,
	}

	return listComment
}
