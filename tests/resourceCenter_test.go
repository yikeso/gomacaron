package tests

import (
	_"github.com/yikeso/gomacaron/config"
	_ "github.com/yikeso/gomacaron/models"
	_ "github.com/yikeso/gomacaron/util"
	"github.com/yikeso/gomacaron/models"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"database/sql"
	"time"
	"github.com/yikeso/gomacaron/task"
	"fmt"
)
//不间断创建txt任务测试
func TestCreateTxtTask(t *testing.T)  {
	task.CreateTxtTask()
}

func TestFindResourceCenterById(t *testing.T)  {
	Convey("T_RESOURCECENTER 表进行findById查询",t,func(){
		var id int64 = 65068
		r,err := models.FindResourceCenterById(id)
		log.Println(r)
		So(err,ShouldBeNil)
		So(r.Title.String,ShouldEqual,"千秋红岩")
	})
}

func TestGetResourceCenterTypeById(t *testing.T)  {
	Convey("T_RESOURCECENTER 表根据id查询资源类型",t,func(){
		var id int64 = 65068
		r,err := models.GetResourceCenterTypeById(id)
		log.Println(r)
		So(err,ShouldBeNil)
		So(r.Int64,ShouldEqual,6)
	})
}

func TestGetBookWithOutTxt(t *testing.T)  {
	Convey("T_RESOURCECENTER 表查出还没创建txt的资源20条",t,func(){
		var readerId int64 = 2
		r,err := models.GetBookWithOutTxt(readerId)
		log.Println(r)
		So(err,ShouldBeNil)
		So(len(r),ShouldBeBetweenOrEqual,0,20)
	})
}

func TestUpdateReaderRcMaxresourceid(t *testing.T) {
	Convey("t_Reader_Rc 表更新，电子书创建txt任务进度",t,func(){
		var readerId int64 = 2
		var maxId int64 = 65069
		err := models.UpdateReaderRcMaxresourceid(maxId,readerId)
		So(err,ShouldBeNil)
	})
}

func TestFindErroeLogByStatus(t *testing.T) {
	Convey("根据错误状态分页查找错误",t, func() {
		status := []int{0}
		errlogs,err := models.FindErroeLogByStatus(&status,1,10)
		log.Println(errlogs)
		So(err,ShouldBeNil)
		So(len(errlogs),ShouldBeBetweenOrEqual,0,10)
	})
}

func TestInsertCreatTxtError(t *testing.T) {

	Convey("插入一个新的创建txt错误",t, func() {
		errorLog := models.BookCreateTxtErrorLog{Bookid:1234,Booktitle:sql.NullString{String:"test",Valid:true},
			Errormessage:sql.NullString{String:"panic",Valid:true},
			Createtime:sql.NullString{String:time.Now().Format(models.TIMESTAMP_FORMATE),Valid:true},
			Lastmodify:sql.NullString{String:time.Now().Format(models.TIMESTAMP_FORMATE),Valid:true}}
		tx,err := models.GetErrorLogTx()
		So(err,ShouldBeNil)
		if err != nil {
			return
		}
		err = models.InsertCreatTxtError(tx,&errorLog)
		log.Println(errorLog)
		tx.Commit()
		So(err,ShouldBeNil)
		So(errorLog.Id,ShouldNotEqual,0)
	})
}

func TestCreateTxtByResourceCenter(t *testing.T) {
	Convey("测试根据resourceCenter生成txt文件",t, func() {
		rs,_ := models.FindResourceCenterById(85398)
		err := task.CreateTxtByResourceCenter(&rs,true)
		So(err,ShouldBeNil)
	})
}

func TestString(t *testing.T) {
	Convey("测试字符串是指针传递还是对象复制",t, func() {
		old := "asdfghjkl"
		fmt.Println(old)
		test1 := old[:4]
		fmt.Println(test1)
		runeArray := []rune(old)
		runeArray[1] = 'a'
		fmt.Println(old)
		fmt.Println(string(runeArray))
	})
}