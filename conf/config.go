package conf

import (
	"gopkg.in/ini.v1"
)

var GlobalCfg *ini.File

func init() {
	var err error
	GlobalCfg, err = ini.Load("./config.ini")
	if err != nil {
		panic(err)
	}
}