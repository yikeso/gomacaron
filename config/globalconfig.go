package config

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/robfig/cron"
	"log"
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
	fmt.Println("开始读取配置文件")
	task := cron.New()
	spec := "*/8 * * * * ?"
	task.AddFunc(spec, ReadConfig)
	task.Start()
}
/**
 * 读取配置文件
 */
func ReadConfig(){
	c.Mymap = make(map[string]string)
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
func Read(node, key string) (str string) {
	if len(node) > 0{
		key = node + middle + key
	}
	str, found := c.Mymap[key]
	if !found {
		log.Println("该配置属性:"+key+" 不存在")
		return
	}
	return
}