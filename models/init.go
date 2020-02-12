package models

import (
	"bookzone/sysinit"
	"bookzone/util/log"
)

func init() {
	sysinit.DatabaseEngine.Sync2(new(Category))
	log.Infof("register Category...")
}