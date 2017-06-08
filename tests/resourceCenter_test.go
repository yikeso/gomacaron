package tests

import (
	_ "github.com/yikeso/gomacaron/models"
	"github.com/yikeso/gomacaron/models"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"log"
)

func TestFindResourceCenterById(t *testing.T)  {
	Convey("T_RESOURCECENTER 表进行findById查询",t,func(){
		var id int64 = 218
		r,err := models.FindResourceCenterById(id)
		log.Println("deleted:",r.Deleted.String)
		So(err,ShouldBeNil)
	})
}