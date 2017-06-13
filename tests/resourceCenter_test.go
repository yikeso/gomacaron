package tests

import (
	_ "github.com/yikeso/gomacaron/models"
	"github.com/yikeso/gomacaron/models"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFindResourceCenterById(t *testing.T)  {
	Convey("T_RESOURCECENTER 表进行findById查询",t,func(){
		var id int64 = 65068
		r,err := models.FindResourceCenterById(id)
		So(err,ShouldBeNil)
		So(r.Title.String,ShouldEqual,"千秋红岩")
	})
}

func TestGetResourceCenterTypeById(t *testing.T)  {
	Convey("T_RESOURCECENTER 表根据id查询资源类型",t,func(){
		var id int64 = 65068
		r,err := models.GetResourceCenterTypeById(id)
		So(err,ShouldBeNil)
		So(r.Int64,ShouldEqual,6)
	})
}

func TestGetBookWithOutTxt(t *testing.T)  {
	Convey("T_RESOURCECENTER 表查出还没创建txt的资源20条",t,func(){
		var readerId int64 = 2
		r,err := models.GetBookWithOutTxt(readerId)
		So(err,ShouldBeNil)
		So(len(r),ShouldBeBetweenOrEqual,0,20)
	})
}

func TestUpdateReaderRcMaxresourceid(t *testing.T)  {
	Convey("t_Reader_Rc 表更新，电子书创建txt任务进度",t,func(){
		var readerId int64 = 2
		var maxId int64 = 65069
		err := models.UpdateReaderRcMaxresourceid(maxId,readerId)
		So(err,ShouldBeNil)
	})
}