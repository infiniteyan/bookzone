package models

import "time"

type BookData struct {
	BookId         int       `json:"book_id"`
	BookName       string    `json:"book_name"`
	Identify       string    `json:"identify"`
	OrderIndex     int       `json:"order_index"`
	Description    string    `json:"description"`
	PrivatelyOwned int       `json:"privately_owned"`
	PrivateToken   string    `json:"private_token"`
	DocCount       int       `json:"doc_count"`
	CommentCount   int       `json:"comment_count"`
	CreateTime     time.Time `json:"create_time"`
	CreateName     string    `json:"create_name"`
	ModifyTime     time.Time `json:"modify_time"`
	Cover          string    `json:"cover"`
	MemberId       int       `json:"member_id"`
	Username       int       `json:"user_name"`
	Editor         string    `json:"editor"`
	RelationshipId int       `json:"relationship_id"`
	RoleId         int       `json:"role_id"`
	RoleName       string    `json:"role_name"`
	Status         int
	Vcnt           int    `json:"vcnt"`
	Collection     int    `json:"star"`
	Score          int    `json:"score"`
	CntComment     int    `json:"cnt_comment"`
	CntScore       int    `json:"cnt_score"`
	ScoreFloat     string `json:"score_float"`
	LastModifyText string `json:"last_modify_text"`
	Author         string `json:"author"`
	AuthorURL      string `json:"author_url"`
}

func NewBookData() *BookData {
	return &BookData{}
}