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
	"github.com/yikeso/gomacaron/jsonobj"
	"strconv"
	"strings"
	"github.com/alecthomas/log4go"
	"encoding/json"
	"io/ioutil"
	"database/sql"
)

var resourceUrl string

//不间断制作txt
func CreateTxtTask(){
	for {
		var dev bool = false
		if strings.EqualFold("development",
			config.GetProp("", "runmodel", "development")) {
			dev = true
		}
		var resourceCenters []models.ResourceCenter
		if dev {
			resourceCenters, _ = models.GetBookWithOutTxt(2)
		} else {
			resourceCenters, _ = models.GetBookWithOutTxt(3)
		}
		if len(resourceCenters) < 1 {
			time.Sleep(10 * time.Minute)
			continue
		}
		var err error
		for _,rsc := range resourceCenters {
			err = CreateTxtByResourceCenter(&rsc,dev)
			if err != nil {
				log4go.Error(err.Error())
				errLog := new(models.BookCreateTxtErrorLog)
				errLog.Bookid = rsc.Id
				errLog.Booktitle = rsc.Title
				errLog.Errormessage = sql.NullString{String:err.Error(),
					Valid:true}
				errLog.Createtime = sql.NullString{String:time.Now().Format(models.TIMESTAMP_FORMATE),
					Valid:true}
				errLog.Lastmodify = errLog.Createtime
				tx,subErr := models.GetErrorLogTx()
				if subErr != nil {
					log4go.Error(subErr.Error())
				}else {
					models.InsertCreatTxtError(tx, errLog)
					tx.Commit()
				}
			}
			if dev {
				models.UpdateReaderRcMaxresourceid(rsc.Id,2)
			}else {
				models.UpdateReaderRcMaxresourceid(rsc.Id,3)
			}
		}
	}
}

func CreateTxtByResourceCenter(rs *models.ResourceCenter,isDev bool)(err error){
	//捕获异常返回错误
	defer func(){
		if r := recover();r != nil{
			err = errors.New(fmt.Sprint("电子书id为",rs.Id," 的资源制作txt文件失败\n",r))
		}
	}()
	resourceUrl = config.GetProp("","resourceUrl","http://resource.gbxx123.com/")
	resourceAbsolutePath := rs.ResourceUrl.String
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
	switch rs.Type.Int64 {
	case 7:
		bookTxtDir = fmt.Sprint(bookTxtDir,"7article/",datePath)
		bookHtmlDir = fmt.Sprint(bookHtmlDir,"7article/",datePath)
	default :
		bookTxtDir = fmt.Sprint(bookTxtDir,"6book/",datePath)
		bookHtmlDir = fmt.Sprint(bookHtmlDir,"6book/",datePath)
	}
	os.MkdirAll(bookTxtDir,0777)
	os.MkdirAll(bookHtmlDir,0777)
	chapterArray,_ := getChapterArrayByCharpterUrlList(resourceUrl,urlEntity.ChapterArray)
	coverUrl := urlEntity.Cover
	chapter0 := new(jsonobj.Chapter0)
	chapter0.BookTitle = rs.Title.String
	if coverUrl = strings.TrimSpace(coverUrl);len(coverUrl) > 0{
		coverUrl = fmt.Sprint(resourceUrl,coverUrl)
		chapter0.ImageUrl = util.GetCoverImageUrlByCoverUrl(coverUrl)
	}
	if len(strings.TrimSpace(urlEntity.Dir)) > 0 {
		ca,_ := util.GetChapterArrayByBcontent(resourceUrl,urlEntity.Dir,int(rs.Type.Int64))
		if len(ca) > 0 {
			chapterArray = ca
		}
	}
	if len(chapterArray) == 0 {
		panic("电子书无任何章节")
	}
	var creatHtml bool = false
	if strings.EqualFold("true",config.GetProp("","creatHtml","true")){
		creatHtml = true
	}
	var txtPath string
	var result *models.ChapterEntity
	for i,dir := range chapterArray {
		log4go.Debug(fmt.Sprint("当前制作id为 ",rs.Id," 的电子书\nURL为 ",dir.Url," 的章节"))
		if len(dir.Anchor) == 0 {
			result,err = util.GetChapterEntity(dir.Url,"","",int(rs.Type.Int64),creatHtml)
		}else if i < len(chapterArray) - 1{
			if len(chapterArray[i+1].Anchor) == 0 {
				result,err = util.GetChapterEntity(dir.Url,dir.Anchor,"",int(rs.Type.Int64),creatHtml)
			}else{
				result,err = util.GetChapterEntity(dir.Url,dir.Anchor,chapterArray[i+1].Anchor,int(rs.Type.Int64),creatHtml)
			}
		}else{
			result,err = util.GetChapterEntity(dir.Url,dir.Anchor,"",int(rs.Type.Int64),creatHtml)
		}
		if err != nil {
			err = errors.New(fmt.Sprint("章节id为 ",i+1," 的章节制作失败",err.Error()))
			return
		}
		dir.ParagraphNum = result.Paragraph
		txtPath = fmt.Sprint(bookTxtDir,i+1,".txt")
		err = ioutil.WriteFile(txtPath,result.Content.Bytes(),0666)
		if err != nil {
			err = errors.New(fmt.Sprint("章节id为 ",i+1," 的章节制作失败",err.Error()))
			return
		}
		txtPath = fmt.Sprint(bookHtmlDir,i+1,".txt")
		err = ioutil.WriteFile(txtPath,result.HtmlContent.Bytes(),0666)
		if err != nil {
			err = errors.New(fmt.Sprint("章节id为 ",i+1," 的章节制作失败",err.Error()))
			return
		}
	}
	for _,dir := range chapterArray {
		dir.ChapterParagraphNum = countChapterParagraphNum(dir)
		if dir.Level == 1 {
			chapter0.Directories = append(chapter0.Directories,dir)
		}
	}
	txtPath = fmt.Sprint(bookTxtDir,0,".txt")
	jsonStr,err := json.Marshal(chapter0)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(txtPath,jsonStr,0666)
	if creatHtml {
		txtPath = fmt.Sprint(bookHtmlDir, 0, ".txt")
		err = ioutil.WriteFile(txtPath, jsonStr, 0666)
	}
	return
}

//统计章节总段落数，包括其子标题
func countChapterParagraphNum(dir *jsonobj.Directory)(p int){
	if dir == nil || len(dir.Id) == 0 {
		return
	}
	p += dir.ParagraphNum
	if len(dir.SubDirectory) > 0 {
		for _,sub := range dir.SubDirectory{
			p += countChapterParagraphNum(sub)
		}
	}
	return
}

func getChapterArrayByCharpterUrlList(resourceUrl string,urlList []string)(chapterArray []*jsonobj.Directory,err error){
	if urlList == nil || len(urlList) == 0{
		return
	}
	var dir *jsonobj.Directory
	for i,s := range urlList {
		dir = new(jsonobj.Directory)
		dir.Id = strconv.Itoa(i+1)
		dir.Level = 1
		if i := strings.Index(s,"#");i > 0 {
			dir.Url = fmt.Sprint(resourceUrl,util.FormatFilePath(s[0:i-1]))
			dir.Anchor = s[i+1:]
		}else if i == 0{
			err = errors.New("章节路径为空")
		}else{
			dir.Url = s
		}
		chapterArray = append(chapterArray,dir)
	}
	return
}