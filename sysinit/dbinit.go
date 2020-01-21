package sysinit

import (
	"bookzone/conf"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var GlobalDB *sql.DB

func dbinit() {
	registerDatabase("w")
}

func registerDatabase(alias string) {
	if len(alias) == 0 {
		return
	}

	if alias == "w" || alias == "default" {
		alias = "w"
	}

	dbNameKey := "db_" + alias + "_database"
	dbUserKey := "db_" + alias + "_username"
	dbPwdKey := "db_" + alias + "_password"
	dbHostKey := "db_" + alias + "_host"
	dbPortKey := "db_" + alias + "_port"

	dbName := conf.GlobalCfg.Section("mysql").Key(dbNameKey).String()
	dbUser := conf.GlobalCfg.Section("mysql").Key(dbUserKey).String()
	dbPwd := conf.GlobalCfg.Section("mysql").Key(dbPwdKey).String()
	dbHost := conf.GlobalCfg.Section("mysql").Key(dbHostKey).String()
	dbPort := conf.GlobalCfg.Section("mysql").Key(dbPortKey).String()

	fmt.Println(dbName, dbUser)

	GlobalDB, err := sql.Open("mysql", dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", GlobalDB.Stats())
}