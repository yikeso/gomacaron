package models

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/yikeso/gomacaron/config"
	"database/sql"
	"github.com/yikeso/gomacaron/config"
)

var resourceDb *sql.DB
var errorLogDb *sql.DB

func init(){
	initDb()
}

func initDb(){
	node := config.Read("common","runmodel")
	var err error
	resourceDb,err = sql.Open(config.Read(node,"drivername1"), config.Read(node,"datasourcename1"))
	if err != nil {
		panic(err)
	}
	errorLogDb,err = sql.Open(config.Read(node,"drivername1"), config.Read(node,"datasourcename1"))
	if err != nil {
		panic(err)
	}
}
