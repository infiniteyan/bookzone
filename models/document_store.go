package models

import (
	"bookzone/sysinit"
	"bookzone/util/log"
)

type DocumentStore struct {
	DocumentId 		int 	`orm:"pk;auto;column(document_id)"`
	Markdown 		string	`orm:"type(text);"`
	Content 		string	`orm:"type(text);"`
}

func (this *DocumentStore) TableName() string {
	return "md_document_store"
}

func (this *DocumentStore) SelectField(docId int, field string) string {
	docStore := &DocumentStore{DocumentId: docId}

	sysinit.DatabaseEngine.Get(docStore)
	if field == "content" {
		return docStore.Content
	} else {
		return docStore.Markdown
	}
}

func (this *DocumentStore) InsertOrUpdate(fields ...string) error {
	docStore := &DocumentStore{DocumentId: this.DocumentId}
	has, err := sysinit.DatabaseEngine.Get(docStore)

	if err != nil {
		log.Errorf(err.Error())
		return err
	}
	if has {
		_, err = sysinit.DatabaseEngine.Update(this, &DocumentStore{DocumentId: this.DocumentId})
	} else {
		_, err = sysinit.DatabaseEngine.Insert(this)
	}
	return err
}

func (this *DocumentStore) Delete(docId int) {
	docStore := &DocumentStore{DocumentId: docId}
	sysinit.DatabaseEngine.Delete(docStore)
}