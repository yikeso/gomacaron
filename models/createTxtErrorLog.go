package models

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"fmt"
	"bytes"
)

type BookCreateTxtErrorLog struct {
	Id           int64			`db:"Id"`
	Bookid       int64			`db:"Bookid"`
	Booktitle    sql.NullString	`db:"Booktitle"`
	Createtime   sql.NullString	`db:"Createtime"`
	Lastmodify   sql.NullString	`db:"Lastmodify"`
	Errormessage sql.NullString	`db:"Errormessage"`
	Status       int			`db:"Status"`
	Chapternum   sql.NullInt64	`db:"Chapternum"`
}

var (
	errorLogBaseSql = "Bookid,Booktitle,Createtime,Lastmodify,Errormessage,Status,Chapternum"
)
//标签错误为重试，仍失败
func MarkBookCreateTxtFaildAnginByBookId(tx *sqlx.Tx,bookid int64)(err error){
	return UpdateErrorLogStatusByBookId(tx,1,bookid)
}

//标签错误为正在重试
func MarkBookCreateTxtTryAnginNowByBookId(tx *sqlx.Tx,bookid int64)(err error){
	return UpdateErrorLogStatusByBookId(tx,3,bookid)
}
//update错误状态
func UpdateErrorLogStatusByBookId(tx *sqlx.Tx,status int,bookid int64)(err error){
	query := "UPDATE book_create_txt_error_log SET status = ?,lastmodify = CURRENT_TIMESTAMP WHERE bookid = ?"
	r,err := tx.Exec(query,status,bookid)
	if err != nil {
		return
	}
	_,err = r.RowsAffected()
	return
}
//传入一个错误状态数组的指针，根据错误状态查找错误记录
//page页码，pageSize页面大小
func FindErroeLogByStatus(status *[]int,page int,pageSize int)(errorLogs []BookCreateTxtErrorLog,err error){
	query := bytes.Buffer{}
	query.WriteString("SELECT Id,")
	query.WriteString(errorLogBaseSql)
	query.WriteString(" FROM book_create_txt_error_log WHERE status IN (")
	l := len(*status) - 1
	for i,s := range *status {
		query.WriteString(fmt.Sprint(s))
		if i == l{
			break
		}
		query.WriteString(",")
	}
	query.WriteString(fmt.Sprint(") LIMIT ",(page -1)*pageSize,",",page*pageSize))
	err = errorLogDb.Select(&errorLogs,query.String())
	return
}
//根据bookid将错误标记为处理成功
func MarkErrorDealSuccess(tx *sqlx.Tx,bookid int64)(err error){
	return UpdateErrorLogStatusByBookId(tx,-1,bookid)
}
//根据错误id批量忽略错误
func IngorCreateTxtErrorByIds(tx *sqlx.Tx,ids *[]int64)(err error){
	query := bytes.Buffer{}
	query.WriteString("UPDATE book_create_txt_error_log SET status = 2 WHERE id IN (")
	l := len(*ids) - 1
	for i,s := range *ids {
		query.WriteString(fmt.Sprint(s))
		if i == l{
			break
		}
		query.WriteString(",")
	}
	query.WriteString(")")
	r,err := tx.Exec(query.String())
	if err != nil{
		return
	}
	_,err = r.RowsAffected()
	return
}
//根据电子统计电子书的错误记录个数
func CountCreatTxtErrorByBookId(booid int64)(i int,err error){
	query := "SELECT COUNT(id) FROM book_create_txt_error_log WHERE bookid = ?"
	err = errorLogDb.Select(&i,query,booid)
	return
}
//根据错误id查找错误
func FindBookCreateTxtErrorLogById(id int64)(errorLog BookCreateTxtErrorLog,err error){
	query := fmt.Sprint("SELECT Id,",errorLogBaseSql," FROM book_create_txt_error_log WHERE id = ?")
	err = errorLogDb.Select(&errorLog,query,id)
	return
}
//根据错误状态统计电子书的错误记录个数
func CountCreatTxtErrorByStatus(status *[]int)(i int,err error){
	query := bytes.Buffer{}
	query.WriteString("SELECT COUNT(id) FROM book_create_txt_error_log WHERE status IN (")
	l := len(*status) - 1
	for i,s := range *status {
		query.WriteString(fmt.Sprint(s))
		if i == l{
			break
		}
		query.WriteString(",")
	}
	query.WriteString(")")
	err = errorLogDb.Select(&i,query.String(),status)
	return
}
//插入错误日志
func InsertCreatTxtError(tx *sqlx.Tx,errorLog *BookCreateTxtErrorLog)(err error){
	query := bytes.Buffer{}
	query.WriteString("INSERT INTO book_create_txt_error_log (")
	query.WriteString(errorLogBaseSql)
	query.WriteString(") VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE ")
	query.WriteString("Bookid=?,Booktitle=?,Createtime=?,Lastmodify=?,")
	query.WriteString("Errormessage=?,Status=?,Chapternum=?")
	r,err := tx.Exec(query.String(),errorLog.Bookid,errorLog.Booktitle,errorLog.Createtime,errorLog.Lastmodify,
		errorLog.Errormessage,errorLog.Status,errorLog.Chapternum,
		errorLog.Bookid,errorLog.Booktitle,errorLog.Createtime,errorLog.Lastmodify,
		errorLog.Errormessage,errorLog.Status,errorLog.Chapternum)
	if err != nil {
		return
	}
	id,e := r.LastInsertId()
	if e == nil {
		errorLog.Id = id
	}
	return
}