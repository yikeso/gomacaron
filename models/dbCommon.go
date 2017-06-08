package models

import (
	"database/sql"
	"github.com/go-gorp/gorp"
	_ "github.com/go-sql-driver/mysql"
	"github.com/yikeso/gomacaron/config"
	"log"
)

var resourceDb *gorp.DbMap
var errorLogDb *gorp.DbMap

func init(){
	initDb()
}

func initDb(){
	config.ReadConfig()
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	node := config.Read("common","runmodel")
	db, err := sql.Open(config.Read(node,"drivername1"), config.Read(node,"datasourcename1"))
	logError(err, "sql.Open failed")
	db.SetMaxOpenConns(10)
	// construct a gorp DbMap
	resourceDb = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	resourceDb.AddTableWithName(ResourceCenter{}, "T_RESOURCECENTER").SetKeys(true, "Id")
	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	//err = resourceDb.CreateTablesIfNotExists()

	db1, err := sql.Open(config.Read(node,"drivername1"), config.Read(node,"datasourcename1"))
	logError(err, "sql.Open failed")
	db1.SetMaxOpenConns(10)
	errorLogDb = &gorp.DbMap{Db: db1, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	//errorLogDb.AddTableWithName(ResourceCenter{}, "posts").SetKeys(true, "Id")
	//err = resourceDb.CreateTablesIfNotExists()
	//logError(err, "Create tables failed")
}

func logError(err error, msg string){
	if err != nil {
		log.Fatalln(msg, err)
	}
}