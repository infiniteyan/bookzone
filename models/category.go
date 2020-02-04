package models

import (
	"bookzone/sysinit"
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