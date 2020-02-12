package models

import (
	"bookzone/sysinit"
	"bookzone/common"
	"time"
	"errors"
)

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

func (m *BookData) SelectByIdentify(identify string, memberId int) (*BookData, error) {
	var bookData *BookData
	var err error
	if identify == "" || memberId <= 0 {
		return nil, errors.New("invalid parameter")
	}

	book := &Book{Identify: identify}
	_, err = sysinit.DatabaseEngine.Get(book)
	if err != nil {
		return nil, err
	}

	relationship := NewRelationship()
	relationship.BookId = book.BookId
	relationship.RoleId = 0
	_, err = sysinit.DatabaseEngine.Get(relationship)
	if err != nil {
		return nil, errors.New("permission denied")
	}

	member, err := NewMember().Find(relationship.MemberId)
	if err != nil {
		return nil, err
	}

	relationship = NewRelationship()
	relationship.BookId = book.BookId
	relationship.MemberId = memberId
	_, err = sysinit.DatabaseEngine.Get(relationship)
	if err != nil {
		return nil, err
	}

	bookData = book.ToBookData()
	bookData.CreateName = member.Account
	bookData.MemberId = relationship.MemberId
	bookData.RoleId = relationship.RoleId
	bookData.RoleName = common.BookRole(bookData.RoleId)
	bookData.RelationshipId = relationship.RelationshipId
	return bookData, nil
}