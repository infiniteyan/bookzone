package models

import (
	"bookzone/sysinit"
	"bookzone/util"
	"errors"
	"fmt"
	"bookzone/util/log"
	"strings"
)

type Category struct {
	Id 		int
	Pid 	int
	Title 	string
	Intro 	string
	Icon 	string
	Cnt 	int
	Sort 	int
	Status 	bool
}

func (this *Category) TableName() string {
	return "md_category"
}

func (this *Category) GetCates(pid int, status int) ([]Category, error) {
	if pid == -1 {
		every := make([]Category, 0)
		err := sysinit.DatabaseEngine.Find(&every)
		return every, err
	} else {
		category := &Category{Pid: pid}
		_, err := sysinit.DatabaseEngine.Get(category)
		if err != nil {
			return nil, err
		} else {
			return []Category{*category}, nil
		}
	}
}

func (this *Category) Find(id int) (*Category, error) {
	category := &Category{Id: id}
	_, err := sysinit.DatabaseEngine.Get(category)
	if err != nil {
		return nil, err
	} else {
		return category, nil
	}
}

func (this *Category) InsertMulti(pid int, cates string) error {
	var err error
	slice := strings.Split(cates, "\n")
	if len(slice) == 0 {
		return errors.New("invalid parameter")
	}

	for _, item := range slice {
		if item = strings.TrimSpace(item); item != "" {
			cate := Category{
				Pid: pid,
				Title: item,
				Status: true,
			}
			_, err = sysinit.DatabaseEngine.Insert(&cate)
			if err != nil {
				log.Errorf(err.Error())
			}
		}
	}

	return nil
}

func (this *Category) Delete(id int) error {
	cate := Category{Id: id}
	_, err := sysinit.DatabaseEngine.Get(&cate)
	if err != nil {
		return err
	}
	if cate.Cnt > 0 {
		return errors.New("fail to delete, Category.Cnt > 0")
	}

	_, err = sysinit.DatabaseEngine.Delete(&Category{Id: id})
	if err != nil {
		return err
	}

	_, err = sysinit.DatabaseEngine.Exec("delete from md_category where pid = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (this *Category) UpdateField(id int, field, val string) error {
	_, err := sysinit.DatabaseEngine.Exec(fmt.Sprintf("update md_category set %v = %v where id = ?", field, val), id)
	return err
}

var counting = false

type Count struct {
	Cnt 		int
	CategoryId 	int
}

func CountCategory() {
	if counting {
		return
	}

	counting = true
	defer func() {
		counting = false
	}()

	var counts []*Count
	var err error
	sql := "select count(bc.id) cnd, bc.category_id from md_book_category bc left join md_books b on b.book_id = bc.book_id where b.privately_owned = 0 group by bc.category_id"
	retSlice, err := sysinit.DatabaseEngine.QueryString(sql)
	if err != nil {
		return
	}

	for _, data := range retSlice {
		var cnt Count
		util.Map2struct(data, &cnt)
		counts = append(counts, &cnt)
	}
	if len(counts) == 0 {
		return
	}

}