package controllers

import (
	"github.com/Unknwon/macaron"
	"fmt"
	"github.com/yikeso/gomacaron/jsonobj"
	"github.com/yikeso/gomacaron/util"
	"strings"
	"github.com/yikeso/gomacaron/config"
	"strconv"
	"github.com/yikeso/gomacaron/models"
	"io/ioutil"
	"github.com/alecthomas/log4go"
	"encoding/json"
)

func CoverImageHandler(ctx *macaron.Context) (result string){
	bookIdStr := ctx.Req.Form.Get("bookId")
	if len(bookIdStr) <1 {
		br := &jsonobj.BaseRespone{Status:3,Message:"电子书id不得为空"}
		result = util.Obj2String(br)
		return
	}
	bookId,err := strconv.ParseInt(bookIdStr, 10, 64)
	if err != nil {
		br := &jsonobj.BaseRespone{Status:4,Message:"电子书id不合法"}
		result = util.Obj2String(br)
		return
	}
	var dev bool = false
	if strings.EqualFold("development",
		config.GetProp("", "runmodel", "development")) {
		dev = true
	}
	ty,err := models.GetResourceCenterTypeById(bookId)
	if err != nil {
		log4go.Error(err.Error())
		panic(err)
	}
	txtDir,_ := util.GetResouceCenterDirByBookIdAndBookType(bookId,int(ty.Int64),dev)
	chapter0Content,err := ioutil.ReadFile(fmt.Sprint(txtDir,"0.txt"))
	if err != nil {
		log4go.Error(err.Error())
		panic(err)
	}
	chapter0 := &jsonobj.Chapter0{}
	err = json.Unmarshal(chapter0Content,chapter0)
	if err != nil {
		log4go.Error(err.Error())
		panic(err)
	}
	or := &jsonobj.OneRespone{One:chapter0.ImageUrl}
	return util.Obj2String(or)
}

