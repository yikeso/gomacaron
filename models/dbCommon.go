package models

import (
	_ "github.com/go-sql-driver/mysql"//初始化驱动
	"github.com/yikeso/gomacaron/config"
	"github.com/jmoiron/sqlx"
	"github.com/alecthomas/log4go"
)

const TIMESTAMP_FORMATE = "2006-01-02 03:04:05"

var resourceDb *sqlx.DB
var errorLogDb *sqlx.DB

func init(){
	initDb()
	/*node,_ := config.Read("common","runmodel")
	task := cron.New()
	spec,_ := config.Read(node,"reloadDB")
	task.AddFunc(spec,initDb)
	task.Start()*/
}

func initDb(){
	log4go.Info("初始化数据源")
	node,_ := config.Read("common","runmodel")
	driver,_ := config.Read(node,"drivername1")
	datasource,_ := config.Read(node,"datasourcename1")
	var err error
	resourceDb,err = sqlx.Connect(driver,datasource)
	resourceDb.SetMaxOpenConns(2)
	if err != nil {
		panic(err)
	}
	driver,_ = config.Read(node,"drivername2")
	datasource,_ = config.Read(node,"datasourcename2")
	errorLogDb,err = sqlx.Connect(driver,datasource)
	errorLogDb.SetMaxOpenConns(2)
	if err != nil {
		panic(err)
	}
}
//获取错误日志的事务
func GetErrorLogTx()(tx *sqlx.Tx,err error){
	return errorLogDb.Beginx()
}
//获取电子书资源的事务
func GetResourceTx()(tx *sqlx.Tx,err error){
	return resourceDb.Beginx()
}