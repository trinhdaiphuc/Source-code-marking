package models

type (
	Comment struct {
		ID        string   `json:"id,omitempty" bson:"_id"`
		FileID    string   `json:"file_id" bson:"file_id"`
		Content   string   `json:"content" bson:"content"`
		StartLine Position `json:"start_line,omitempty" bson:"start_line,omitempty"`
		EndLine   Position `json:"end_line,omitempty" bson:"end_line,omitempty"`
	}

	Position struct {
		Row    int `json:"row" bson:"row"`
		Column int `json:"column" bson:"column"`
	}
)
