package models

import (
	"bookzone/sysinit"
	"time"
)

type AttachmentData struct {
	Attachment
	IsExist       bool
	BookName      string
	DocumentName  string
	FileShortSize string
	Account       string
	LocalHttpPath string
}

type Attachment struct {
	AttachmentId int `orm:"pk;auto" json:"attachment_id"`
	BookId       int ` json:"book_id"`
	DocumentId   int ` json:"doc_id"`
	Name         string
	Path         string    `orm:"size(2000)" json:"file_path"`
	Size         float64   `orm:"type(float)" json:"file_size"`
	Ext          string    `orm:"size(50)" json:"file_ext"`
	HttpPath     string    `orm:"size(2000)" json:"http_path"`
	CreateTime   time.Time `orm:"type(datetime);auto_now_add" json:"create_time"`
	CreateAt     int       `orm:"type(int)" json:"create_at"`
}

func (this *Attachment) TableName() string {
	return "md_attachment"
}

func NewAttachment() *Attachment {
	return &Attachment{}
}

func (this *Attachment) Insert() error {
	_, err := sysinit.DatabaseEngine.Insert(this)
	return err
}

func (this *Attachment) Update() error {
	_, err := sysinit.DatabaseEngine.Update(this)
	return err
}

func (this *Attachment) SelectByDocumentId(docId int) ([]*Attachment, error) {
	var attachs []*Attachment

	err := sysinit.DatabaseEngine.Where("doc_id = ?", docId).Find(&attachs)
	if err != nil {
		return nil, err
	} else {
		return attachs, nil
	}
}