package sysinit

import (
	"bookzone/conf"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DatabaseEngine *xorm.Engine

func dbinit() {
	registerDatabase()
}

func registerDatabase() {
	var err error
	dbNameKey := "db_database"
	dbUserKey := "db_username"
	dbPwdKey := "db_password"
	dbHostKey := "db_host"
	dbPortKey := "db_port"

	dbName := conf.GlobalCfg.Section("mysql").Key(dbNameKey).String()
	dbUser := conf.GlobalCfg.Section("mysql").Key(dbUserKey).String()
	dbPwd := conf.GlobalCfg.Section("mysql").Key(dbPwdKey).String()
	dbHost := conf.GlobalCfg.Section("mysql").Key(dbHostKey).String()
	dbPort := conf.GlobalCfg.Section("mysql").Key(dbPortKey).String()
	dataSourceName := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8", dbUser, dbPwd, dbHost, dbPort, dbName)

	DatabaseEngine, err = xorm.NewEngine("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	DatabaseEngine = DatabaseEngine
	fmt.Println("register database success.")
}