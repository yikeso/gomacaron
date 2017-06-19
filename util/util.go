package util

import (
	"unsafe"
	"reflect"
	"container/list"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
	"github.com/yikeso/gomacaron/jsonobj"
	"fmt"
	"strconv"
	"net/http"
	"bytes"
	"errors"
	"image"
	"github.com/yikeso/gomacaron/models"
)

var  blockLevelElements []string

func init(){
	initBlockLevelElements()
}

func initBlockLevelElements(){
	blockLevelElements = make([]string,0,30)
	blockLevelElements = append(blockLevelElements,"address") //定义地址
	blockLevelElements = append(blockLevelElements,"caption");//定义表格标题
	blockLevelElements = append(blockLevelElements,"dd");//定义列表中定义条目
	blockLevelElements = append(blockLevelElements,"div");//定义文档中的分区或节
	blockLevelElements = append(blockLevelElements,"dl");//定义列表
	blockLevelElements = append(blockLevelElements,"dt");//定义列表中的项目
	blockLevelElements = append(blockLevelElements,"fieldset");//定义一个框架集
	blockLevelElements = append(blockLevelElements,"form");//创建 HTML 表单
	blockLevelElements = append(blockLevelElements,"h1");//定义最大的标题
	blockLevelElements = append(blockLevelElements,"h2");//定义副标题
	blockLevelElements = append(blockLevelElements,"h3");//定义副标题
	blockLevelElements = append(blockLevelElements,"h4");//定义副标题
	blockLevelElements = append(blockLevelElements,"h5");//定义副标题
	blockLevelElements = append(blockLevelElements,"h6");//定义副标题
	blockLevelElements = append(blockLevelElements,"hr");//定义水平线
	blockLevelElements = append(blockLevelElements,"fieldset");//元素定义标题
	blockLevelElements = append(blockLevelElements,"noframes");//为那些不支持框架的浏览器显示文本，于 frameset 元素内部
	blockLevelElements = append(blockLevelElements,"noscript");//定义在脚本未被执行时的替代内容
	blockLevelElements = append(blockLevelElements,"ol");//定义有序列表
	blockLevelElements = append(blockLevelElements,"ul");//定义无序列表
	blockLevelElements = append(blockLevelElements,"li");
	blockLevelElements = append(blockLevelElements,"pre");//定义预格式化的文本
	blockLevelElements = append(blockLevelElements,"table");//标签定义 HTML 表格
	blockLevelElements = append(blockLevelElements,"thead");
	blockLevelElements = append(blockLevelElements,"tbody");
	blockLevelElements = append(blockLevelElements,"tr");
	blockLevelElements = append(blockLevelElements,"td");
}

//电子书索引页相关属性实体
type UrlEntity struct {
	ChapterArray []string
	Dir string
	Cover string
	Chapters string
}
//解析电子书索引得到相关属性
func GetBookUrl(fileUrl string)(entity *UrlEntity,err error){
	fmt.Println(fileUrl)
	doc,err := goquery.NewDocument(fileUrl)
	if err != nil {
		return
	}
	xhtml := doc.Find("xhtml").Children()
	entity = new(UrlEntity)
	xhtml.Each(func(i int,s *goquery.Selection){
		_,b := s.Attr("id")
		if !b{
			return
		}
		t := s.Text()
		if strings.Index(t,"b_content") > 0 {
			entity.Dir = t
		}else if strings.Index(t,"chapter") > 0 {
			entity.ChapterArray = append(entity.ChapterArray,t)
		}else if strings.Index(t,"cover") > 0 {
			entity.Cover = s.Text()
		}else if strings.Index(t,"content") > 0 {
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
func FormatFilePath(path string)(s string){
	//正则
	reg := regexp.MustCompile("[\\\\/]+")
	return reg.ReplaceAllString(path,"/")
}
//解析封面url获取封面图片的url
func GetCoverImageUrlByCoverUrl(coverUrl string)(imageUrl string){
	doc,err := goquery.NewDocument(coverUrl)
	if err != nil {
		return
	}
	s := doc.Find("img")
	if s.Length() > 0{
		imageUrl,_ = s.Eq(0).Attr("src")
		imageUrl = fmt.Sprint(coverUrl[:strings.LastIndex(coverUrl,"/")],imageUrl)
	}
	return
}

//解析bcontent目录
func GetChapterArrayByBcontent(resourceUrl string,bcontentPath string,ty int)(chapterArray []*jsonobj.Directory,err error){
	var doc *goquery.Document
	switch ty {
	case 7:
		doc,err = goquery.NewDocument(fmt.Sprint(resourceUrl,bcontentPath))
	default:
		doc,err = goquery.NewDocument(fmt.Sprint(resourceUrl,bcontentPath))
		if err != nil {
			panic(err)
		}
	}
	allA := doc.Find("a")
	if allA.Length() < 1{
		return
	}
	dirPath := bcontentPath[:strings.LastIndex(bcontentPath,"/")+1]
	rootArray := getSubDir(nil,allA,dirPath)
	dirSlice := make([]*jsonobj.Directory,0,allA.Length())
	for _,dir := range rootArray {
		dirSlice = append(dirSlice,dir)
		setSubDir(dir,allA,dirPath,dirSlice)
	}
	latelyUrl := ""
	for _,dir := range dirSlice {
		if strings.EqualFold(latelyUrl,dir.Url){
			dir.NewPage = false
		}else {
			dir.NewPage = true
		}
		if dir.NewPage || dir.Level == 1 {
			chapterArray = append(chapterArray,dir)
		}
	}
	return
}
//给标题设置子标题
func setSubDir(parent *jsonobj.Directory,sel *goquery.Selection,dirPath string,dirSlice []*jsonobj.Directory){
	if parent == nil || len(parent.Id) < 1{
		return
	}
	subArray := getSubDir(parent,sel,dirPath)
	if len(subArray) < 1{
		return
	}
	parent.SubDirectory = subArray
	for _,dir := range subArray{
		dirSlice = append(dirSlice,dir)
		setSubDir(dir,sel,dirPath,dirSlice)
	}
}

//获取子标题
func getSubDir(parent *jsonobj.Directory,sel *goquery.Selection,dirPath string)(subArray []*jsonobj.Directory){
	n := 1
	subArray = make([]*jsonobj.Directory,0,10)
	switch {
		case parent == nil || len(parent.Id) == 0 :
			sel.Each(func(i int,s *goquery.Selection){
				title := strings.TrimSpace(s.Text())
				if len(title) == 0 {
					return
				}
				dir := jsonobj.Directory{}
				dir.Title = title
				dir.Level = 1
				href,found := s.Attr("href")
				if !found {
					return
				}
				href = strings.TrimSpace(href)
				if len(href) < 1 {
					return
				}
				anchorIndex := strings.Index(href,"#")
				if anchorIndex > 0 {
					anchor := href[anchorIndex+1:]
					if strings.Index(anchor,"-") != -1 {
						return
					}else{
						dir.Url = fmt.Sprint(dirPath,href[:anchorIndex-1])
						dir.Anchor = anchor
					}
				}else {
					dir.Url = fmt.Sprint(dirPath,href)
				}
				dir.Id = strconv.Itoa(n)
				subArray = append(subArray,&dir)
				n++
			})
	case len(strings.TrimSpace(parent.Anchor)) > 0:
		sel.Each(func(i int,s *goquery.Selection){
			pal := len(parent.Anchor)
			if pal == 0 {
				return
			}
			title := strings.TrimSpace(s.Text())
			if len(title) < 1 {
				return
			}
			dir := jsonobj.Directory{}
			dir.Title = title
			dir.Level = parent.Level + 1
			href,found := s.Attr("href")
			if !found {
				return
			}
			href = strings.TrimSpace(href)
			if len(href) < 1 {
				return
			}
			anchorIndex := strings.Index(href,"#")
			if anchorIndex == -1{
				return
			}
			anchor := href[anchorIndex+1:]
			if len(anchor) <= pal {
				return
			}
			if !strings.HasPrefix(anchor,parent.Anchor) {
				return
			}
			anchorRight := anchor[pal:]
			if strings.Count(anchorRight,"-") != 1 {
				return
			}
			dir.Id = fmt.Sprint(parent.Id,"-",n)
			dir.Url = fmt.Sprint(dirPath,href[:anchorIndex-1])
			dir.Anchor = anchor
			subArray = append(subArray,&dir)
			n++
		})
	}
	return
}
//处理对应的章节内容
func GetChapterEntity(url,anchor1,anchor2 string,ty int,createHtml bool)(chatperEntity *models.ChapterEntity,err error){
	switch ty {
	case 7:
		chatperEntity,err = getArticleChapterEntity(url,anchor1,anchor2,true)
	default:
		chatperEntity,err = getBookChapterEntity(url,anchor1,anchor2,true)
	}
	return
}
//处理电子书章节内容
func getBookChapterEntity(url,anchor1,anchor2 string,createHtml bool)(chatperEntity *models.ChapterEntity,err error){
	return makeContentByUrl(url,anchor1,anchor2,true)
}
//处理文章章节内容
func getArticleChapterEntity(url,anchor1,anchor2 string,createHtml bool)(chatperEntity *models.ChapterEntity,err error){
	return makeContentByUrl(url,anchor1,anchor2,true)
}

func makeContentByUrl(url,anchor1,anchor2 string,createHtml bool)(chatperEntity *models.ChapterEntity,err error){
	body,err := GetHtmlBodyByUrl(url)
	if err != nil{
		return
	}
	chapter,err := cutBodyByAnchor(body,anchor1,anchor2)
	if err != nil{
		return
	}
	p := splitChapterByTag(chapter)
	chatperEntity.Paragraph = len(p)
	buffer := bytes.NewBuffer([]byte(body))
	buffer.Reset()
	paragraphArrayToString(p,buffer,url)
	chatperEntity.Content = buffer
	if createHtml {
		htmlBuffer := bytes.NewBuffer(make([]byte, 10000))
		err = paragraphArrayAddTagToString(p, htmlBuffer, url)
		if err != nil {
			return
		}
		chatperEntity.HtmlContent = htmlBuffer
	}
	return
}

func paragraphArrayAddTagToString(p []string,buffer *bytes.Buffer,url string) (err error){
	dir := url[:strings.LastIndex(url,"/")]
	pBuffer := bytes.NewBuffer(make([]byte,1000))
	for i,s := range p{
		if strings.HasPrefix(s,"</img") ||
			strings.HasPrefix(s,"</video") ||
			strings.HasPrefix(s,"</audio") ||
			strings.HasPrefix(s,"</embed"){
			doc,innerErr := goquery.NewDocumentFromReader(bytes.NewBufferString(
				fmt.Sprint("<div><span>",s,"</span></div>")))
			if innerErr != nil {
				err = innerErr
				return
			}
			e := doc.Find("div").Eq(0)
			e.SetAttr("id",fmt.Sprint(i+1))
			e.SetAttr("pid",fmt.Sprint(i+1))
			e = e.Find("span").Eq(0)
			e.SetAttr("parentid",fmt.Sprint(i+1))
			e.SetAttr("wordoffset",fmt.Sprint(1))
			e.AddClass("wchar")
			e.SetAttr("sytle","display:inline-block;text-indent:0em;")
			e = e.Children().Eq(0)
			src,found := e.Attr("src")
			if found && !strings.HasPrefix(src,"http") {
				e.SetAttr("src",fmt.Sprint(dir,src))
			}
			e.SetAttr("parentid",fmt.Sprint(i+1))
			if strings.HasPrefix(s,"</img") {
				imageStyle, err := getImageStyle(src)
				if err == nil {
					e.SetAttr("style", imageStyle)
				}
			}
			s,innerErr = doc.Find("body").Eq(0).Html()
			if innerErr != nil {
				err = innerErr
				return
			}
			buffer.WriteString(s)
		}else {
			err = addSpan(s,i+1,pBuffer)
			if err != nil {
				return
			}
			buffer.WriteString("<div id = \"")
			buffer.WriteString(string(i+1))
			buffer.WriteString("\" pid = \"")
			buffer.WriteString(string(i+1))
			buffer.WriteString("\">")
			buffer.WriteString(pBuffer.String())
			buffer.WriteString("</div>")
		}
	}
	return
}

//给段落添加span标签
func addSpan(cont string,p int,buffer *bytes.Buffer)(err error){
	reg := regexp.MustCompile("<!--.*-->")
	cont = reg.ReplaceAllString(cont,"")
	reg = regexp.MustCompile("[\\n\\t\\r]+")
	cont = reg.ReplaceAllString(cont,"")
	doc,err := goquery.NewDocumentFromReader(bytes.NewBufferString(cont))
	if err != nil {
		return
	}
	//给所有元素添加parentid属性
	doc.Find("body").Find("").Each(func(i int,s *goquery.Selection){
		s.SetAttr("parentid",fmt.Sprint(p))
	})
	buffer.Reset()
	var c,n rune
	var isAddTage,isEscape bool
	for _,c = range cont{
		if c == '<' {
			isAddTage = false
		}
		if isAddTage && c == '&'{
			isEscape = true
		}
		if isAddTage{
			if isEscape{
				if c == '&'{
					buffer.WriteString("<span ")
					buffer.WriteString("parentid=\"")
					buffer.WriteString(fmt.Sprint(p))
					buffer.WriteString("\" class=\"wchar\" wordoffset=\"")
					n++
					buffer.WriteString(fmt.Sprint(n))
					buffer.WriteString("\">")
					buffer.WriteString(string(c))
				}else if c == ';' {
					buffer.WriteString(string(c))
					buffer.WriteString("</span>")
				}else {
					buffer.WriteString(string(c))
				}
			}else{
				buffer.WriteString("<span parentid=\"")
				buffer.WriteString(fmt.Sprint(p))
				buffer.WriteString("\" class=\"wchar\" wordoffset=\"")
				n++
				buffer.WriteString(fmt.Sprint(n))
				buffer.WriteString("\">")
				buffer.WriteString(string(c))
				buffer.WriteString("</span>")
			}
		}else{
			buffer.WriteString(string(c))
		}
		if c == '>'{
			isAddTage = true
		}
		if isAddTage {
			if c == ';' {
				isEscape = false
			}
		}
	}
	return
}

//读取图片宽高
func getImageStyle(imageUrl string)(imageStyle string,err error){
	resp,err := http.Get(imageUrl)
	if err != nil {
		return
	}
	img,_,err := image.Decode(resp.Body)
	if err != nil{
		return
	}
	b := img.Bounds()
	imageStyle = fmt.Sprint("width:",b.Dy(),"px;height:",b.Dy(),"px;")
	return
}

//段落拼接字符串
func paragraphArrayToString(p []string,buffer *bytes.Buffer,url string){
	for _,s := range p{
		buffer.WriteString(s)
		buffer.WriteString("~~")
	}
}

//根据p标签，div标签，将段落分段
func splitChapterByTag(chapter string)(r []string){
	p := make([]string,0,100)
	var ltIndex,gtIndex int
	var str string
	ltIndex = strings.Index(chapter,"<")
	for ltIndex != -1 {
		gtIndex = strings.Index(chapter,">")
		if gtIndex < 1 {
			break
		}
		str = chapter[ltIndex:gtIndex]
		if strings.HasPrefix(str,"<div") || strings.HasPrefix(str,"</div") ||
			strings.HasPrefix(str,"<p") || strings.HasPrefix(str,"</p") ||
			strings.HasPrefix(str,"</img") ||
			strings.HasPrefix(str,"</video") ||
			strings.HasPrefix(str,"</audio") ||
			strings.HasPrefix(str,"</embed"){
			if ltIndex > 0 {
				str = chapter[:ltIndex-1]
				p = append(p, str)
			}
		}else if strings.HasPrefix(str,"<img") ||
			strings.HasPrefix(str,"<video") ||
			strings.HasPrefix(str,"<audio") ||
			strings.HasPrefix(str,"<embed"){
			if ltIndex > 0 {
				str = chapter[:ltIndex-1]
				p = append(p, str)
			}
			p = append(p,chapter[ltIndex:gtIndex])
		}
		chapter = chapter[gtIndex+1:]
		ltIndex = strings.Index(chapter,"<")
	}
	if len(chapter) > 0{
		p = append(p,chapter)
	}
	buffer := bytes.NewBuffer(make([]byte,500))
	deque := list.New()
	var s string
	for _,str = range p{
		str = removeSpan(str)
		str = strings.TrimSpace(str)
		if len(str) < 1 {
			continue
		}
		if strings.HasPrefix(str,"<") && !strings.HasPrefix(str,"</"){
			s = str[1:strings.Index(str,">")]
			if !strings.HasSuffix(s,"/>"){
				s = str[:strings.Index(s,"\\s")]
				if matcheBlockLevelElements(s){
					deque.PushBack(s)
				}
			}
		}
		if deque.Len() < 1 {
			s = buffer.String()
			if len(s) > 0{
				r = append(r,s)
			}
			buffer.Reset()
			r = append(r,str)
			continue
		}else {
			buffer.WriteString(str)
		}
		if strings.HasSuffix(str,">") && !strings.HasSuffix(str,"/>") {
			s = str[strings.LastIndex(str,"<"):]
			if strings.HasPrefix(s,"</"){
				s = strings.TrimSpace(str[2:strings.Index(s,">")])
				e := deque.Back()
				if strings.EqualFold(s,e.Value.(string)) {
					deque.Remove(e)
				}
			}
		}
	}
	return
}

func matcheBlockLevelElements(targetName string)bool{
	for _,s := range blockLevelElements{
		if strings.EqualFold(s,targetName){
			return true
		}
	}
	return false
}

//去掉span标签
func removeSpan(partHtml string) string {
	i := strings.Index(partHtml,"<span")
	reg := regexp.MustCompile("[\\n\\t\\r]+")
	if i == -1 {
		return reg.ReplaceAllString(partHtml,"")
	}
	buffer := bytes.NewBuffer(make([]byte,500))
	for i != -1 {
		buffer.WriteString(partHtml[:i-1])
		partHtml = partHtml[i+1:]
		i = strings.Index(partHtml,">")
		partHtml = partHtml[i+1:]
		i = strings.Index(partHtml,"</span")
		if i != -1 {
			buffer.WriteString(partHtml[:i-1])
			partHtml = partHtml[i+1:]
			i = strings.Index(partHtml,">")
			partHtml = partHtml[i+1:]
		}
		i = strings.Index(partHtml,"<span")
	}
	buffer.WriteString(partHtml)
	return reg.ReplaceAllString(buffer.String(),"")
}

//根据章节锚点，截取章节内容
func cutBodyByAnchor(body,anchor1,anchor2 string)(subBody string,err error){
	if len(body) == 0 {
		err = errors.New("章节内容为空")
		return
	}
	if anchor1 == anchor2 {
		err = errors.New("章节内容截取，两锚点相同")
		return
	}
	if len(anchor1) == 0{
		subBody = body
	}else if len(anchor2) == 0{
		i := strings.Index(body,anchor1)
		pre := body[:i]
		i = strings.LastIndex(pre,"<")
		subBody = body[i:]
	}else{
		i := strings.Index(body,anchor1)
		pre := body[:i]
		i = strings.LastIndex(pre,"<")
		e := strings.LastIndex(body,anchor2)
		subBody = body[i:e]
		e = strings.LastIndex(subBody,"<")
		subBody = body[:e-1]
	}
	//将两个以上英文空格，换成一个
	reg := regexp.MustCompile("\\s{2,}")
	subBody = reg.ReplaceAllString(subBody,"\\s")
	//将标签内容为一个空格的内容，去掉
	reg = regexp.MustCompile(">\\s<")
	subBody = reg.ReplaceAllString(subBody,"><")
	if len(subBody) == 0 {
		err = errors.New("截取后章节内容为空")
		return
	}
	return
}

func GetHtmlByUrl(url string)(html string,err error){
	resp,err := http.Get(url)
	if err != nil{
		return
	}
	buffer := bytes.NewBuffer(make([]byte,resp.ContentLength))
	_,err = buffer.ReadFrom(resp.Body)
	if err != nil {
		return
	}
	html = buffer.String()
	return
}

func GetHtmlBodyByUrl(url string)(body string,err error){
	html,err := GetHtmlByUrl(url)
	if err != nil {
		return
	}
	start := strings.Index(html,"<body")
	end := strings.LastIndex(html,"</body")
	body = html[start+1:end-1]
	start = strings.Index(body,">")
	body = body[start+1:]
	return
}