package models

import (
	"bookzone/sysinit"
	"bookzone/util"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Book struct {
	BookId 			int					`json:"book_id"`
	BookName 		string				`json:"book_name"`
	Identify       	string    			`json:"identify"`
	OrderIndex     	int       			`json:"order_index"`
	Description    	string    			`json:"description"`
	Cover          	string    			`json:"cover"`
	Editor         	string    			`json:"editor"`
	Status        	int       			`json:"status"`
	PrivatelyOwned 	int       			`json:"privately_owned"`
	PrivateToken   	string    			`json:"private_token"`
	MemberId       	int      			`json:"member_id"`
	CreateTime     	time.Time 			`son:"create_time"`
	ModifyTime     	time.Time 			`json:"modify_time"`
	ReleaseTime    	time.Time 			`json:"release_time"`
	DocCount       	int       			`json:"doc_count"`
	CommentCount   	int       			`json:"comment_count"`
	Vcnt           	int       			`json:"vcnt"`
	Collection     	int       			`xorm:"star" json:"star"`
	Score          	int       			`json:"score"`
	CntScore      	int 				`json:"cnt_score"`
	CntComment     	int 				`json:"cnt_comment"`
	Author         	string    			`json:"author"`
	AuthorUrl      	string    			`json:"author_url"`
}

func (this *Book) TableName() string {
	return "md_books"
}

func NewBook() *Book {
	return &Book{}
}

func (this *Book) HomeData(pageIndex, pageSize int, cid int, fields ...string) ([]Book, int, error) {
	if len(fields) == 0 {
		fields = append(fields, "book_id", "book_name", "identify", "cover", "order_index")
	}

	fieldStr := "b." + strings.Join(fields, ",b.")
	sqlFmt := "select %s from md_books b left join md_book_category c on b.book_id = c.book_id where c.category_id = " + strconv.Itoa(cid)
	sql := fmt.Sprintf(sqlFmt, fieldStr)
	fmt.Println(sql)

	resultSlice, err := sysinit.DatabaseEngine.Query(sql)
	if err != nil {
		return nil, 0, err
	}
	var books []Book
	for _, data := range resultSlice {
		var book Book
		id , err := strconv.Atoi(string(data["book_id"]))
		if err != nil {
			continue
		}
		index, err := strconv.Atoi(string(data["order_index"]))
		if err != nil {
			continue
		}
		book.BookId = id
		book.OrderIndex = index
		book.BookName = string(data["book_name"])
		book.Identify = string(data["identify"])
		book.Cover = string(data["cover"])

		books = append(books, book)
	}

	return books, len(books), nil
}

func (this *Book) Insert() error {
	var err error
	if _, err =sysinit.DatabaseEngine.Insert(this); err != nil {
		return err
	}

	relationship := RelationShip{BookId: this.BookId, MemberId: this.MemberId, RoleId: 0}
	if err = relationship.Insert(); err != nil {
		return err
	}

	document := Document{BookId: this.BookId, DocumentName: "空白文档", Identify: "blank", MemberId: this.MemberId}
	var id int64
	if id, err = document.InsertOrUpdate(); err == nil {
		docStore := DocumentStore{DocumentId: int(id), Markdown: ""}
		err = docStore.InsertOrUpdate()
	}
	return err
}

func (this *Book) Update(cols ...string) error {
	bk := &Book{BookId: this.BookId}
	has, err := sysinit.DatabaseEngine.Get(bk)

	if err != nil {
		return err
	}
	if !has {
		log.Println("please insert first")
		return err
	}

	_, err = sysinit.DatabaseEngine.Cols(cols...).Update(this, bk)
	return err
}

func (this *Book) SelectByIdentify(value string) (*Book, error) {
	tmp := &Book{}
	tmp.Identify = value
	_, err := sysinit.DatabaseEngine.Get(tmp)
	if err != nil {
		return nil, err
	} else {
		this = tmp
		return this, nil
	}
}

func (this *Book) ToBookData() *BookData {
	m := &BookData{}
	m.BookId = this.BookId
	m.BookName = this.BookName
	m.Identify = this.Identify
	m.OrderIndex = this.OrderIndex
	m.Description = strings.Replace(this.Description, "\r\n", "<br/>", -1)
	m.PrivatelyOwned = this.PrivatelyOwned
	m.PrivateToken = this.PrivateToken
	m.DocCount = this.DocCount
	m.CommentCount = this.CommentCount
	m.CreateTime = this.CreateTime
	m.ModifyTime = this.ModifyTime
	m.Cover = this.Cover
	m.MemberId = this.MemberId
	m.Status = this.Status
	m.Editor = this.Editor
	m.Vcnt = this.Vcnt
	m.Collection = this.Collection
	m.Score = this.Score
	m.ScoreFloat = util.ScoreFloat(this.Score)
	m.CntScore = this.CntScore
	m.CntComment = this.CntComment
	m.Author = this.Author
	m.AuthorURL = this.AuthorUrl
	if this.Editor == "" {
		m.Editor = "markdown"
	}
	return m
}

func (this *Book) SelectPage(pageIndex, pageSize, memberId int, PrivatelyOwned int) ([]*BookData, int, error) {
	var totalCount int
	var err error
	var books []*BookData
	sql1 := "select count(b.book_id) as total_count from md_books as b left join " +
		"md_relationship as r on b.book_id=r.book_id and r.member_id = ? where r.relationship_id > 0  and b.privately_owned = " + strconv.Itoa(PrivatelyOwned)

	retSlice, err := sysinit.DatabaseEngine.Query(sql1, memberId)
	if err != nil {
		return nil, 0, err
	}
	totalCount = len(retSlice)

	offset := (pageIndex - 1) * pageSize
	sql2 := "select book.*,rel.member_id,rel.role_id,m.account as create_name from md_books as book" +
		" left join md_relationship as rel on book.book_id=rel.book_id and rel.member_id = ?" +
		" left join md_relationship as rel1 on book.book_id=rel1.book_id and rel1.role_id=0" +
		" left join md_members as m on rel1.member_id=m.member_id " +
		" where rel.relationship_id > 0 %v order by book.book_id desc limit " + fmt.Sprintf("%d,%d", offset, pageSize)
	sql2 = fmt.Sprintf(sql2, " and book.privately_owned="+strconv.Itoa(PrivatelyOwned))

	retSlice, err = sysinit.DatabaseEngine.Query(sql2, memberId)
	if err != nil {
		return nil, totalCount, err
	}

	for _, data := range retSlice {
		var book Book
		byteContent, err := json.Marshal(data)
		if err != nil {
			continue
		}
		err = json.Unmarshal(byteContent, &book)
		if err != nil {
			continue
		}
		books = append(books, book.ToBookData())
	}

	return books, totalCount, nil
}

func (this *Book) RefreshDocumentCount(bookId int) {
	doc := &Document{BookId: bookId}
	count, err := sysinit.DatabaseEngine.Count(doc)
	if err != nil {
		log.Println(err)
	} else {
		bean := NewBook()
		bean.DocCount = int(count)
		condiBean := NewBook()
		condiBean.BookId = bookId
		sysinit.DatabaseEngine.Update(bean, condiBean)
	}
}