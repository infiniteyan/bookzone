package models

import (
	"bookzone/sysinit"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

type Document struct {
	DocumentId   int           `orm:"pk;auto;column(document_id)" json:"doc_id"`
	DocumentName string        `orm:"column(document_name);size(500)" json:"doc_name"`
	Identify     string        `orm:"column(identify);size(100);index;null;default(null)" json:"identify"`
	BookId       int           `orm:"column(book_id);type(int)" json:"book_id"`
	ParentId     int           `orm:"column(parent_id);type(int);default(0)" json:"parent_id"`
	OrderSort    int           `orm:"column(order_sort);default(0);type(int)" json:"order_sort"`
	Release      string        `orm:"column(release);type(text);null" json:"release"`
	CreateTime   time.Time     `orm:"column(create_time);type(datetime);auto_now_add" json:"create_time"`
	MemberId     int           `orm:"column(member_id);type(int)" json:"member_id"`
	ModifyTime   time.Time     `orm:"column(modify_time);type(datetime);default(null);auto_now" json:"modify_time"`
	ModifyAt     int           `orm:"column(modify_at);type(int)" json:"-"`
	Version      int64         `orm:"type(bigint);column(version)" json:"version"`
	AttachList   []*Attachment `orm:"-" json:"attach"`
	Vcnt         int           `orm:"column(vcnt);default(0)" json:"vcnt"`
	Markdown     string        `orm:"-" json:"markdown"`
}

func (this *Document) TableName() string {
	return "md_documents"
}

func NewDocument() *Document {
	return &Document{Version:time.Now().Unix()}
}

func (this *Document) SelectByDocId(id int) (*Document, error) {
	if id <= 0 {
		return nil, errors.New("Invalid parameter")
	}

	doc := &Document{DocumentId: id}
	_, err := sysinit.DatabaseEngine.Get(doc)
	if err != nil {
		return nil, err
	} else {
		this = doc
		return this, nil
	}
}

func (this *Document) SelectByIdentify(bookId int, identify string) (*Document, error) {
	doc := &Document{BookId: bookId, Identify: identify}
	_, err := sysinit.DatabaseEngine.Get(doc)
	if err != nil {
		return nil, err
	} else {
		this = doc
		return this, nil
	}
}

func (this *Document) InsertOrUpdate() (int64, error) {
	var id int64
	var err error

	id = int64(this.DocumentId)
	this.ModifyTime = time.Now()
	this.DocumentName = strings.TrimSpace(this.DocumentName)
	if this.DocumentId > 0 {
		condiDoc := &Document{DocumentId: this.DocumentId}
		_, err = sysinit.DatabaseEngine.Update(this, condiDoc)
		return id, err
	}

	selectedDocument := &Document{Identify: this.Identify, BookId: this.BookId}
	_, err = sysinit.DatabaseEngine.Get(selectedDocument)
	if err != nil {
		return -1, errors.New("fail to query database")
	}

	if selectedDocument.DocumentId == 0 {
		this.CreateTime = time.Now()
		id, err = sysinit.DatabaseEngine.Insert(this)
		NewBook().RefreshDocumentCount(this.BookId)
		return id, err
	} else {
		_, err = sysinit.DatabaseEngine.Update(this)
		id = int64(selectedDocument.DocumentId)
		return id, err
	}
}

func (this *Document) GetMenuTop(bookId int) ([]*Document, error) {
	cols := []string{"document_id", "document_name", "member_id", "parent_id", "book_id", "identify"}
	sql := "select %v from md_documents where book_id = ? and parent_id = 0 order by order_sort, document_id limit 5000"
	sql = fmt.Sprintf(sql, strings.Join(cols, ","))
	log.Printf("execute sql:%s\n", sql)
	retSlice, err := sysinit.DatabaseEngine.Query(sql, bookId)
	if err != nil {
		return nil ,err
	}

	var docs []*Document
	for _, data := range retSlice {
		var doc Document
		byteContent, err := json.Marshal(data)
		if err != nil {
			continue
		}
		err = json.Unmarshal(byteContent, &doc)
		if err != nil {
			continue
		}
		docs = append(docs, &doc)
	}

	return docs, nil
}