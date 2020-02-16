package models

import (
	"bookzone/sysinit"
	"bookzone/util"
	"fmt"
	"strconv"
	"strings"
)

type CollectionData struct {
	BookId      int    `json:"book_id"`
	BookName    string `json:"book_name"`
	Identify    string `json:"identify"`
	Description string `json:"description"`
	DocCount    int    `json:"doc_count"`
	Cover       string `json:"cover"`
	MemberId    int    `json:"member_id"`
	Nickname    string `json:"user_name"`
	Vcnt        int    `json:"vcnt"`
	Collection  int    `json:"star"`
	Score       int    `json:"score"`
	CntComment  int    `json:"cnt_comment"`
	CntScore    int    `json:"cnt_score"`
	ScoreFloat  string `json:"score_float"`
	OrderIndex  int    `json:"order_index"`
}

type Collection struct {
	Id       int 	`xorm:"pk autoincr" json:"id"`
	MemberId int 	`json:"member_id"`
	BookId   int	`json:"book_id"`
}

func NewCollection() *Collection {
	return &Collection{}
}

func (m *Collection) TableName() string {
	return "md_star"
}

func (m *Collection) Collection(uid, bid int) (bool, error) {
	var cancel bool
	star := &Collection{MemberId: uid, BookId: bid}
	_, err := sysinit.DatabaseEngine.Get(star)
	if err != nil {
		return false, err
	}

	if star.Id > 0 {
		if _, err = sysinit.DatabaseEngine.Delete(&Collection{Id: star.Id}); err == nil {
			util.IncOrDec("md_books", "star", fmt.Sprintf("book_id = %v and star > 0", bid), false, 1)
		}
		cancel = true
	} else {
		cancel = false
		if _, err = sysinit.DatabaseEngine.Insert(star); err == nil {
			util.IncOrDec("md_books", "star", fmt.Sprintf("book_id = %v", bid), true, 1)
		}
	}

	return cancel, nil
}

func (m *Collection) DoesCollection(uid, bid interface{}) bool {
	var collection Collection
	collection.MemberId, _ = strconv.Atoi(fmt.Sprintf("%v", uid))
	collection.BookId, _ = strconv.Atoi(fmt.Sprintf("%v", bid))

	_, err := sysinit.DatabaseEngine.Get(&collection)
	if err != nil {
		return false
	}
	if collection.Id > 0 {
		return true
	}
	return false
}

func (m *Collection) List(mid, p, listRows int) (int64, []*CollectionData, error) {
	sql := "select book_id from md_star where member_id = ? order by id desc limit %v offset %v"
	sql = fmt.Sprintf(sql, listRows, (p - 1) * listRows)
	retSlice, err := sysinit.DatabaseEngine.QueryString(sql, mid)
	if err != nil {
		return 0, nil, err
	}

	var stars []*Collection
	for _, data := range retSlice {
		var star Collection
		util.Map2struct(data, &star)
		stars = append(stars, &star)
	}

	bids := []string{}
	for _, v := range stars {
		bids = append(bids, strconv.Itoa(v.BookId))
	}
	var books []*CollectionData
	var totalCount int64
	bidstr := strings.Join(bids, ",")
	sql = fmt.Sprintf("select b.*, m.nickname from md_books b left join md_members m on m.member_id = b.member_id where b.book_id in (%v)", bidstr)
	retSlice, err = sysinit.DatabaseEngine.QueryString(sql)
	if err != nil {
		return 0, nil, err
	}

	for _, data := range retSlice {
		var cd CollectionData
		util.Map2struct(data, &cd)
		books = append(books, &cd)
	}
	totalCount, _ = sysinit.DatabaseEngine.Count(&Collection{MemberId: mid})

	return totalCount, books, nil
}