package models

import (
	"bookzone/sysinit"
	"log"
)

func init() {
	sysinit.DatabaseEngine.Sync2(new(Category))
	log.Println("register Category...")
}