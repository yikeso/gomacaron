package routers

import (
	"github.com/Unknwon/macaron"
	"strings"
	"github.com/yikeso/gomacaron/config"
)

func GetRouters() (m *macaron.Macaron){
	//读取配置文件中的运行模式
	runmodel := config.Read("common", "runmodel")
	if strings.EqualFold("development", runmodel){
		macaron.Env = macaron.DEV
	}else if strings.EqualFold("production", runmodel) {
		macaron.Env = macaron.PROD
	}else {
		macaron.Env = macaron.TEST
	}
	m = macaron.Classic()
	m.Get("/", myHandleer)
	return m
}

func myHandleer(ctx *macaron.Context) (string){
	return "the request path is:" + ctx.Req.RequestURI
}