package task

import (
	_ "github.com/yikeso/gomacaron/models"
	"github.com/yikeso/gomacaron/models"
	"errors"
	"fmt"
	"github.com/yikeso/gomacaron/config"
	"github.com/yikeso/gomacaron/util"
	"time"
	"os"
	"container/list"
	"github.com/yikeso/gomacaron/jsonobj"
	"strconv"
	"strings"
)

func init(){

}


func createTxtByResourceCenter(rs *models.ResourceCenter,isDev bool)(err error){
	//捕获异常返回错误
	defer func(){
		if r := recover();r != nil{
			err = errors.New(fmt.Sprint(r))
		}
	}()
	resourceUrl := config.GetProp("","resourceUrl","http://resource.gbxx123.com/")
	resourceAbsolutePath := rs.ResourceUrl
	urlEntity,err := util.GetBookUrl(fmt.Sprint(resourceUrl,resourceAbsolutePath))
	if err != nil {
		err = errors.New(fmt.Sprint("id为：",rs.Id," 的电子书解析索引index失败\n",err.Error()))
		return
	}
	var bookTxtDir,bookHtmlDir string;
	if isDev {
		bookTxtDir = "e:/bookTxtDir/"
		bookHtmlDir = "e:/bookHtmlDir/"
	}else {
		bookTxtDir = "./bookTxtDir/"
		bookHtmlDir = "./bookHtmlDir/"
	}
	datePath := fmt.Sprint(time.Unix(rs.Id*config.ID_TO_TIME,0).Format("2006/01/02"),
		"/",rs.Id,"/")
	switch rs.Type {
	case 7:
		bookTxtDir = fmt.Sprint(bookTxtDir,"7article/",datePath)
		bookHtmlDir = fmt.Sprint(bookHtmlDir,"7article/",datePath)
	default :
		bookTxtDir = fmt.Sprint(bookTxtDir,"6book/",datePath)
		bookHtmlDir = fmt.Sprint(bookHtmlDir,"6book/",datePath)
	}
	os.MkdirAll(bookTxtDir,0777)
	os.MkdirAll(bookHtmlDir,0777)
	chapterList,err := getChapterListByCharpterUrlList(urlEntity.ChapterList)
	if err != nil {
		err = errors.New(fmt.Sprint("id为：",rs.Id," 的电子书解析目录失败\n",err.Error()))
		return
	}
}

func getChapterListByCharpterUrlList(urlList *list.List)(chapterList *list.List,err error){
	if urlList == nil || urlList.Len() == 0{
		return
	}
	chapterList = list.New()
	var dir jsonobj.Directory
	id := 0
	for e := urlList.Front();e != nil;e.Next() {
		id++
		dir = jsonobj.Directory{SubDirectory:list.New()}
		dir.Id = strconv.Itoa(id)
		dir.Level = 1
		s := e.Value.(string)
		if i := strings.Index(s,"#");i > 0 {
			dir.Url = util.FormatFilePath(&s[0:i-1])
			dir.Anchor = s[i+1:]
		}else if i == 0{
			err = errors.New("章节路径为空")
		}else{
			dir.Url = s
		}
		chapterList.PushBack(dir)
	}
	return
}