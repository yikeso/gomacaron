package models

import (
	"database/sql"
)

type ResourceCenter struct {
	Id             int64         `db:"Id"`
	Version        sql.NullInt64 `db:"Version"`
	Title          sql.NullString `db:"Title"`
	ShortTitle     sql.NullString `db:"ShortTitle"`
	AuthorId       sql.NullInt64 `db:"AuthorId"`
	CreateTime     sql.NullString `db:"CreateTime"`
	CreateUserId   sql.NullInt64 `db:"CreateUserId"`
	CreateUserName sql.NullString `db:"CreateUserName"`
	ExamineYear    sql.NullInt64 `db:"ExamineYear"`
	Info           sql.NullString `db:"Info"`
	KeyWord        sql.NullString `db:"KeyWord"`
	Readlevel      sql.NullInt64 `db:"Readlevel"`
	SelectState    sql.NullInt64 `db:"SelectState"`
	Type           sql.NullInt64 `db:"Type"`
	Length         sql.NullInt64 `db:"Length"`
	Money          sql.NullFloat64 `db:"Money"`
	OrireSourceUrl sql.NullString `db:"OrireSourceUrl"`
	CutUrl         sql.NullString `db:"CutUrl"`
	CutCount       sql.NullInt64 `db:"CutCount"`
	ResourceUrl    sql.NullString `db:"ResourceUrl"`
	ShareType      sql.NullInt64 `db:"ShareType"`
	Deleted        []byte `db:"Deleted"`
	Difficulty     sql.NullInt64 `db:"Difficulty"`
	ViceTitle      sql.NullString `db:"ViceTitle"`
	PublishDate    sql.NullString `db:"PublishDate"`
	Source         sql.NullString `db:"Source"`
	Provider       sql.NullString `db:"Provider"`
	Producer       sql.NullString `db:"Producer"`
	Times          sql.NullString `db:"Times"`
	Precisions     sql.NullString `db:"Precisions"`
	Dimension      sql.NullString `db:"Dimension"`
	Place          sql.NullString `db:"Place"`
	EventTimes     sql.NullString `db:"EventTimes"`
	FigureId       sql.NullString `db:"FigureId"`
	Uuid           sql.NullString `db:"Uuid"`
}

var (
	resourceCenterBaseSql = "Version,Title,ShortTitle,AuthorId,CreateTime,CreateUserId," +
		"CreateUserName,ExamineYear,Info,KeyWord,Readlevel,SelectState,Type,Length," +
		"Money,OrireSourceUrl,CutUrl,CutCount,ResourceUrl,ShareType,Deleted,Difficulty," +
		"ViceTitle,PublishDate,Source,Provider,Producer,Times,Precisions,Dimension," +
		"Place,EventTimes,FigureId,Uuid"
)
//查出还没有生成txt的资源20条
func GetBookWithOutTxt(readerId int64) (resourceCenters []ResourceCenter, err error) {
	query := "Select Id," + resourceCenterBaseSql +
		" from T_RESOURCECENTER WHERE deleted='0' AND type in (6,7) " +
		"AND id>(SELECT MAXRESOURCEID FROM T_READER_RC rc WHERE rc.ID = ?) LIMIT 0,20"
	err = resourceDb.Select(&resourceCenters, query, readerId)
	return
}

//根据资源id查询资源
func FindResourceCenterById(id int64) (resourceCenter ResourceCenter, err error) {
	query := "SELECT Id," + resourceCenterBaseSql + " FROM T_RESOURCECENTER WHERE id = ?"
	err = resourceDb.Get(&resourceCenter, query, id)
	return
}

//根据资源id获取资源的type属性
func GetResourceCenterTypeById(id int64) (t sql.NullInt64, err error) {
	query := "SELECT t.type FROM T_RESOURCECENTER t WHERE t.id = ?"
	err = resourceDb.Get(&t, query, id)
	return
}
