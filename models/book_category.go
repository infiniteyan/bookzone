package models

import (
	"bookzone/sysinit"
	"bookzone/util"
	"bookzone/util/log"
	"fmt"
	"strconv"
	"strings"
)

type BookCategory struct {
	Id         int		`json:"id"`
	BookId     int		`json:"book_id"`
	CategoryId int 		`json:"category_id"`
}

func (m *BookCategory) TableName() string {
	return "md_book_category"
}

func (m *BookCategory) SelectByBookId(book_id int) ([]*Category, int64, error) {
	var ret []*Category
	var totalCount int64
	var err error
	sql := "select c.* from md_category c left join md_book_category bc on c.id = bc.category_id where bc.book_id = ?"
	retSlice, err := sysinit.DatabaseEngine.QueryString(sql, book_id)
	if err != nil {
		return nil, 0, err
	}
	totalCount = int64(len(ret))

	for _, data := range retSlice {
		var cat Category
		util.Map2struct(data, &cat)
		ret = append(ret, &cat)
	}

	return ret, totalCount, nil
}

func (m *BookCategory) SetBookCates(bookId int, cids []string) {
	if len(cids) == 0 {
		return
	}
	var err error
	var cates []*Category
	sql := "select id, pid from md_category where id in (%s)"
	sql = fmt.Sprintf(sql, strings.Join(cids, ","))
	retSlice, err := sysinit.DatabaseEngine.QueryString(sql)
	for _, data := range retSlice {
		var cate Category
		util.Map2struct(data, &cate)
		cates = append(cates, &cate)
	}
	cidMap := make(map[string]bool)
	for _, cate := range cates {
		cidMap[strconv.Itoa(cate.Pid)] = true
		cidMap[strconv.Itoa(cate.Id)] = true
	}
	cids = []string{}
	for cid, _ := range cidMap {
		cids = append(cids, cid)
	}

	sysinit.DatabaseEngine.Delete(&BookCategory{BookId: bookId})
	var bookCates []*BookCategory
	for _, cid := range cids {
		cidNum, _ := strconv.Atoi(cid)
		bookCate := &BookCategory{
			CategoryId: cidNum,
			BookId: bookId,
		}
		bookCates = append(bookCates, bookCate)
	}

	for _, bookCate := range bookCates {
		_, err = sysinit.DatabaseEngine.Insert(bookCate)
		if err != nil {
			log.Errorf(err.Error())
		}
	}

	go CountCategory()
}
