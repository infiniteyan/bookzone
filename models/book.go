package models

import (
	"bookzone/sysinit"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Book struct {
	BookId 			int					`orm:"pk;auto" json:"book_id"`
	BookName 		string				`orm:"size(500)" json:"book_name"`
	Identify       	string    			`orm:"size(100);unique" json:"identify"`
	OrderIndex     	int       			`orm:"default(0)" json:"order_index"`
	Description    	string    			`orm:"size(1000)" json:"description"`
	Cover          	string    			`orm:"size(1000)" json:"cover"`
	Editor         	string    			`orm:"size(50)" json:"editor"`
	Status        	int       			`orm:"default(0)" json:"status"`
	PrivatelyOwned 	int       			`orm:"default(0)" json:"privately_owned"`
	PrivateToken   	string    			`orm:"size(500);null" json:"private_token"`
	MemberId       	int      			`orm:"size(100)" json:"member_id"`
	CreateTime     	time.Time 			`orm:"type(datetime);auto_now_add" json:"create_time"`
	ModifyTime     	time.Time 			`orm:"type(datetime);auto_now_add" json:"modify_time"`
	ReleaseTime    	time.Time 			`orm:"type(datetime);" json:"release_time"`
	DocCount       	int       			`json:"doc_count"`
	CommentCount   	int       			`orm:"type(int)" json:"comment_count"`
	Vcnt           	int       			`orm:"default(0)" json:"vcnt"`
	Collection     	int       			`orm:"column(star);default(0)" json:"star"`
	Score          	int       			`orm:"default(40)" json:"score"`
	CntScore      	int
	CntComment     	int
	Author         	string    			`orm:"size(50)"`
	AuthorURL      	string    			`orm:"column(author_url);size(1000)"`
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