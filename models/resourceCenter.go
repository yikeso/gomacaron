package models

import (
	"database/sql"
)

type ResourceCenter struct {
	Id             int64
	Version        sql.NullInt64
	Title          sql.NullString
	ShortTitle     sql.NullString
	AuthorId       sql.NullInt64
	CreateTime     sql.NullString
	CreateUserId   sql.NullInt64
	CreateUserName sql.NullString
	ExamineYear    sql.NullInt64
	Info           sql.NullString
	KeyWord        sql.NullString
	Readlevel      sql.NullInt64
	SelectState    sql.NullInt64
	Type           sql.NullInt64
	Length         sql.NullInt64
	Money          sql.NullFloat64
	OrireSourceUrl sql.NullString
	CutUrl         sql.NullString
	CutCount       sql.NullInt64
	ResourceUrl    sql.NullString
	ShareType      sql.NullInt64
	Deleted        sql.NullString
	Difficulty     sql.NullInt64
	ViceTitle      sql.NullString
	PublishDate    sql.NullString
	Source         sql.NullString
	Provider       sql.NullString
	Producer       sql.NullString
	Times          sql.NullString
	Precisions     sql.NullString
	Dimension      sql.NullString
	Place          sql.NullString
	EventTimes     sql.NullString
	FigureId       sql.NullString
	Uuid           sql.NullString
}

var (
	resourceCenterBaseSql = "Version,Title,ShortTitle,AuthorId,CreateTime,CreateUserId," +
		"CreateUserName,ExamineYear,Info,KeyWord,Readlevel,SelectState,Type,Length," +
		"Money,OrireSourceUrl,CutUrl,CutCount,ResourceUrl,ShareType,Deleted,Difficulty," +
		"ViceTitle,PublishDate,Source,Provider,Producer,Times,Precisions,Dimension," +
		"Place,EventTimes,FigureId,Uuid"
)

func FindResourceCenterById(id int64)(resourceCenter ResourceCenter,err error){
	query := "SELECT Id,"+resourceCenterBaseSql+" FROM T_RESOURCECENTER WHERE id = ?"
	resourceCenter = ResourceCenter{}
	err = resourceDb.SelectOne(&resourceCenter, query,id)
	logError(err, "SelectOne failed")
	return
}
