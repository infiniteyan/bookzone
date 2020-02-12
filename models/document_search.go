package models

import (
	"bookzone/sysinit"
	"bookzone/util"
	"bookzone/util/html2text"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DocumentData struct {
	DocumentId   int       `json:"doc_id"`
	DocumentName string    `json:"doc_name"`
	Identify     string    `json:"identify"`
	Release      string    `json:"release"`
	Vcnt         int       `json:"vcnt"`
	CreateTime   time.Time `json:"create_time"`
	BookId       int       `json:"book_id"`
	BookIdentify string    `json:"book_identify"`
	BookName     string    `json:"book_name"`
}

type DocumentSearch struct {
	DocumentId   int       `json:"doc_id"`
	BookId       int       `json:"book_id"`
	DocumentName string    `json:"doc_name"`
	Identify     string    `json:"identify"`
	Description  string    `json:"description"`
	Author       string    `json:"author"`
	BookName     string    `json:"book_name"`
	BookIdentify string    `json:"book_identify"`
	ModifyTime   time.Time `json:"modify_time"`
	CreateTime   time.Time `json:"create_time"`
}

func NewDocumentSearch() *DocumentSearch {
	return &DocumentSearch{}
}

func (m *DocumentSearch) SearchDocument(keyword string, bookId int, page, size int) ([]*DocumentSearch, int, error) {
	fields := []string{"document_id", "document_name", "identify", "book_id"}

	var sql, sqlCount string
	if bookId == 0 {
		sql = "select %v from md_documents d left join md_books b on d.book_id = b.book_id where b.privately_owned = 0 and (d.document_name like ? or d.release like ? )"
		sqlCount = fmt.Sprintf(sql, "count(d.document_id) cnt")
		sql = fmt.Sprintf(sql, "d." + strings.Join(fields, ",d.")) + " order by d.vcnt desc"
	} else {
		sql = "select %v from md_documents where book_id = " + strconv.Itoa(bookId) + " and (document_name like ? or release like ?) "
		sqlCount = fmt.Sprintf(sql, "count(document_id) cnt")
		sql = fmt.Sprintf(sql, strings.Join(fields, ",")) + " order by vcnt desc"
	}

	var count int
	like := "%" + keyword + "%"

	var err error
	var docs []*DocumentSearch
	retSlice, err := sysinit.DatabaseEngine.QueryString(sqlCount, like, like)
	if err != nil {
		return nil, 0, err
	}

	count = len(retSlice)
	limit := fmt.Sprintf(" limit %v offset %v", size, (page - 1) * size)
	if count > 0 {
		retSlice, err = sysinit.DatabaseEngine.QueryString(sql + limit, like, like)
		if err != nil {
			return nil, 0, err
		}

		for _, data := range retSlice {
			var doc DocumentSearch
			util.Map2struct(data, &doc)
			docs = append(docs, &doc)
		}
	}
	return docs, count, nil
}

func (m *DocumentSearch) GetDocsById(id []int, withoutCont ...bool) ([]*DocumentData, error) {
	if len(id) == 0 {
		return nil, errors.New("no id")
	}

	var idArr []string
	for _, i := range id {
		idArr = append(idArr, fmt.Sprint(i))
	}

	fields := []string{
		"d.document_id", "d.document_name", "d.identify", "d.vcnt", "d.create_time", "b.book_id",
	}

	if len(withoutCont) == 0 || !withoutCont[0] {
		fields = append(fields, "b.identify book_identify", "d.release", "b.book_name")
	}

	sqlFmt := "select " + strings.Join(fields, ",") + " from md_documents d left join md_books b on d.book_id = b.book_id where d.document_id in(%v)"
	sql := fmt.Sprintf(sqlFmt, strings.Join(idArr, ","))

	var docs []*DocumentData
	var err error

	retSlice, err := sysinit.DatabaseEngine.QueryString(sql)
	if err != nil {
		return nil, err
	}

	if len(retSlice) > 0 {
		docMap := make(map[int]*DocumentData)
		for _, data := range retSlice {
			var docData DocumentData
			util.Map2struct(data, &docData)
			docMap[docData.DocumentId] = &docData
		}

		for _, i := range id {
			if doc, ok := docMap[i]; ok {
				doc.Release = html2text.Html2Text(doc.Release)
				docs = append(docs, doc)
			}
		}
	}
	return docs, nil
}