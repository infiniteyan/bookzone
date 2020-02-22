package models

import (
	"bookzone/sysinit"
	"bookzone/util"
	"bytes"
	"errors"
	"fmt"
	"bookzone/util/log"
	"strings"
	"time"
)

type Document struct {
	DocumentId   int           `xorm:"pk autoincr" json:"document_id"`
	DocumentName string        `json:"document_name"`
	Identify     string        `json:"identify"`
	BookId       int           `json:"book_id"`
	ParentId     int           `json:"parent_id"`
	OrderSort    int           `json:"order_sort"`
	Release      string        `json:"release"`
	CreateTime   time.Time     `json:"create_time"`
	MemberId     int           `json:"member_id"`
	ModifyTime   time.Time     `json:"modify_time"`
	ModifyAt     int           `json:"-"`
	Version      int64         `json:"version"`
	AttachList   []*Attachment `xorm:"-" json:"attach"`
	Vcnt         int           `json:"vcnt"`
	Markdown     string        `xorm:"-" json:"markdown"`
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
	log.Infof("execute sql:%s", sql)
	retSlice, err := sysinit.DatabaseEngine.QueryString(sql, bookId)
	if err != nil {
		return nil ,err
	}

	var docs []*Document
	for _, data := range retSlice {
		var doc Document
		util.Map2struct(data, &doc)
		docs = append(docs, &doc)
	}

	return docs, nil
}

func (this *Document) ReleaseContent(bookId int) {
	util.BooksRelease.Set(bookId)
	defer util.BooksRelease.Delete(bookId)

	sql := "select document_id from md_documents where book_id = ? limit 5000"
	retSlice, err := sysinit.DatabaseEngine.QueryString(sql, bookId)
	if err != nil {
		log.Infof(err.Error())
		return
	}

	documents := []*Document{}
	for _, data := range retSlice {
		var doc Document
		util.Map2struct(data, &doc)
		documents = append(documents, &doc)
	}

	documentStore := NewDocumentStore()
	for _, element := range documents {
		content := strings.TrimSpace(documentStore.SelectField(element.DocumentId, "content"))
		element.Release = content
		attachList, err := NewAttachment().SelectByDocumentId(element.DocumentId)

		if err == nil && len(attachList) > 0 {
			content := bytes.NewBufferString("<div class=\"attach-list\"><strong>附件</strong><ul>")
			for _, attach := range attachList {
				li := fmt.Sprintf("<li><a href=\"%s\" target=\"_blank\" title=\"%s\">%s</a></li>", attach.HttpPath, attach.Name, attach.Name)
				content.WriteString(li)
			}
			content.WriteString("</ul></div>")
			element.Release += content.String()
		}
		_, err = sysinit.DatabaseEngine.Update(element, &Document{DocumentId: element.DocumentId})
		if err != nil {
			log.Infof(err.Error())
			continue
		}
	}
}