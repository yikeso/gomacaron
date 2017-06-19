package config

import (
	"os"
	"bufio"
	"strings"
	"log"
	l4g "github.com/alecthomas/log4go"
	"fmt"
)

const middle = "========="


type config struct {
	Mymap  map[string]string
	strcet string
}

/**
 * 局部变量保存配置文件属性
 */
var c config

/**
 * 初始化配置，每两分钟读取一次配置文件，更新配置
 */
func init(){
	readConfig()
	/*task := cron.New()
	spec := "* *//*5 * * * ?"
	task.AddFunc(spec, readConfig)
	task.Start()*/
}
/**
 * 读取配置文件
 */
func readConfig(){
	log.Println("加载log配置文件")
	l4g.LoadConfiguration("D:/go/goDevelopDemo/src/github.com/yikeso/gomacaron/config/log.xml")
	l4g.Debug("开始读取配置文件")
	if c.Mymap == nil {
		c.Mymap = make(map[string]string)
	}
	f,err := os.Open("D:/go/goDevelopDemo/src/github.com/yikeso/gomacaron/config/app.conf")
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer f.Close()
	r := bufio.NewReader(f)
	for {
		b,err := r.ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			return
		}
		s := strings.TrimSpace(b)
		if strings.Index(s, "#") == 0 {
			continue
		}
		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}
		key := c.strcet + middle + frist
		c.Mymap[key] = strings.TrimSpace(second)
	}
}
/**
 * 传入属性，获取属性值
 */
func Read(node, key string) (str string,found bool) {
	if len(node) > 0{
		key = fmt.Sprint(node,middle ,key)
	}else{
		key = fmt.Sprint("common" ,middle ,key)
	}
	str, found = c.Mymap[key]
	if !found {
		log.Println("该配置属性:"+key+" 不存在")
		return
	}
	return
}
//传入节点属性，获取属性值
//node节点，key属性，def如果没有该属性，则使用传入的默认值
func GetProp(node,key,def string) (str string){
	str,found := Read(node,key)
	if !found{
		str = def
	}
	return
}