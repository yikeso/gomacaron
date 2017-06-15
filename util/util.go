package util

import (
	"unsafe"
	"reflect"
	"container/list"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
	"github.com/Unknwon/com"
)
//电子书索引页相关属性实体
type UrlEntity struct {
	ChapterList *list.List
	Dir string
	Cover string
	Chapters string
}
//解析电子书索引得到相关属性
func GetBookUrl(fileUrl string)(entity *UrlEntity,err error){
	doc,err := goquery.NewDocument(fileUrl)
	if err != nil {
		return
	}
	xhtml := doc.Find("xhtml").Children()
	entity = &UrlEntity{ChapterList:list.New()}
	xhtml.Each(func(i int,s *goquery.Selection){
		attr,b := s.Attr("id")
		if !b{
			return
		}
		if strings.EqualFold("b_content",attr) {
			entity.Dir = s.Text()
		}else if strings.EqualFold("chapter",attr) {
			entity.ChapterList.PushBack(s.Text())
		}else if strings.EqualFold("cover",attr) {
			entity.Cover = s.Text()
		}else if strings.EqualFold("content",attr) {
			entity.Chapters = s.Text()
		}
	})
	return
}

//字符串转byte数组
func Str2Byte(s *string) *[]byte {
	return (*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(s))))
}
//格式化文件路径
func FormatFilePath(path *string)(s string){
	//正则
	reg := regexp.MustCompile("[\\\\/]+")
	return reg.ReplaceAllString(*path,"/")
}
