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
	Deleted        []byte
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
//查出还没有生成txt的资源20条
func GetBookWithOutTxt(readerId int64)(resourceCenters []*ResourceCenter,err error) {
	resourceCenters = make([]*ResourceCenter, 0, 10)
	query := "Select * from T_RESOURCECENTER WHERE deleted='0' AND type in (6,7) " +
		"AND id>(SELECT MAXRESOURCEID FROM T_READER_RC rc WHERE rc.ID = ?) LIMIT 0,20"
	rows, err := resourceDb.Query(query, readerId)
	if err != nil {
		return
	}
	for rows.Next() {
		resourceCenter := new(ResourceCenter)
		err = fillInResourceCenterByRows(rows,resourceCenter)
		if err != nil{
			return
		}
		resourceCenters = append(resourceCenters, resourceCenter)
	}
	return
}
//根据资源id查询资源
func FindResourceCenterById(id int64)(resourceCenter ResourceCenter,err error){
	query := "SELECT Id,"+resourceCenterBaseSql+" FROM T_RESOURCECENTER WHERE id = ?"
	row := resourceDb.QueryRow(query,id)
	err = fillInResourceCenterByRow(row,&resourceCenter)
	return
}
//根据资源id获取资源的type属性
func GetResourceCenterTypeById(id int64)(t sql.NullInt64,err error){
	query := "SELECT t.type FROM T_RESOURCECENTER t WHERE t.id = ?"
	row := resourceDb.QueryRow(query,id)
	err = row.Scan(&t)
	return
}

/**
 * 将查询的结果转换为ResourceCenter结构体
 */
func fillInResourceCenterByRow(row *sql.Row,resourceCenter *ResourceCenter)error{
	return row.Scan(&resourceCenter.Id,&resourceCenter.Version,&resourceCenter.Title,&resourceCenter.ShortTitle,
		&resourceCenter.AuthorId,&resourceCenter.CreateTime,&resourceCenter.CreateUserId,&resourceCenter.CreateUserName,
		&resourceCenter.ExamineYear,&resourceCenter.Info,&resourceCenter.KeyWord,&resourceCenter.Readlevel,
		&resourceCenter.SelectState,&resourceCenter.Type,&resourceCenter.Length,&resourceCenter.Money,
		&resourceCenter.OrireSourceUrl,&resourceCenter.CutUrl,&resourceCenter.CutCount,&resourceCenter.ResourceUrl,
		&resourceCenter.ShareType,&resourceCenter.Deleted,&resourceCenter.Difficulty,&resourceCenter.ViceTitle,
		&resourceCenter.PublishDate,&resourceCenter.Source,&resourceCenter.Provider,&resourceCenter.Producer,
		&resourceCenter.Times,&resourceCenter.Precisions,&resourceCenter.Dimension,&resourceCenter.Place,
		&resourceCenter.EventTimes,&resourceCenter.FigureId,&resourceCenter.Uuid)
}

/**
 * 将查询的结果转换为ResourceCenter结构体
 */
func fillInResourceCenterByRows(rows *sql.Rows,resourceCenter *ResourceCenter)error{
	return rows.Scan(&resourceCenter.Id,&resourceCenter.Version,&resourceCenter.Title,&resourceCenter.ShortTitle,
		&resourceCenter.AuthorId,&resourceCenter.CreateTime,&resourceCenter.CreateUserId,&resourceCenter.CreateUserName,
		&resourceCenter.ExamineYear,&resourceCenter.Info,&resourceCenter.KeyWord,&resourceCenter.Readlevel,
		&resourceCenter.SelectState,&resourceCenter.Type,&resourceCenter.Length,&resourceCenter.Money,
		&resourceCenter.OrireSourceUrl,&resourceCenter.CutUrl,&resourceCenter.CutCount,&resourceCenter.ResourceUrl,
		&resourceCenter.ShareType,&resourceCenter.Deleted,&resourceCenter.Difficulty,&resourceCenter.ViceTitle,
		&resourceCenter.PublishDate,&resourceCenter.Source,&resourceCenter.Provider,&resourceCenter.Producer,
		&resourceCenter.Times,&resourceCenter.Precisions,&resourceCenter.Dimension,&resourceCenter.Place,
		&resourceCenter.EventTimes,&resourceCenter.FigureId,&resourceCenter.Uuid)
}