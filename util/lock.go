package util

import "sync"

var BooksRelease = BooksLock{Books: make(map[int]bool)}

type BooksLock struct {
	Books 		map[int]bool
	Lock 		sync.RWMutex
}

func (this BooksLock) Exist(bookId int) bool {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	_, exist := this.Books[bookId]
	return exist
}

func (this BooksLock) Set(bookId int) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	this.Books[bookId] = true
}

func (this BooksLock) Delete(bookId int) {
	this.Lock.RLock()
	defer this.Lock.RUnlock()
	delete(this.Books, bookId)
}