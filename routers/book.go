package routers

import (
	"github.com/Unknwon/macaron"
	"strings"
	"github.com/yikeso/gomacaron/config"
	"time"
	"github.com/alecthomas/log4go"
	"fmt"
	"net/http"
)

func GetRouters() (m *macaron.Macaron){
	//读取配置文件中的运行模式
	runmodel,_ := config.Read("common", "runmodel")
	if strings.EqualFold("development", runmodel){
		macaron.Env = macaron.DEV
	}else if strings.EqualFold("production", runmodel) {
		macaron.Env = macaron.PROD
	}
	m = macaron.New()
	//日志
	m.Use(logger())
	//服务器异常捕获
	m.Use(macaron.Recovery())
	//500错误处理
	m.InternalServerError()
	//404错误处理
	m.NotFound()
	m.Get("/", myHandleer)
	return m
}

func myHandleer(ctx *macaron.Context) (string){
	return "the request path is:" + ctx.Req.RequestURI
}

func logger() macaron.Handler{
	return func (ctx *macaron.Context){
		start := time.Now()
		log4go.Debug(fmt.Sprintf("Started %s %s for %s", ctx.Req.Method, ctx.Req.RequestURI, ctx.RemoteAddr()))
		rw := ctx.Resp.(macaron.ResponseWriter)
		ctx.Next()
		content := fmt.Sprintf("Completed %s %v %s in %v", ctx.Req.RequestURI, rw.Status(), http.StatusText(rw.Status()), time.Since(start))
		if macaron.ColorLog {
			switch rw.Status() {
			case 200, 201, 202:
				content = fmt.Sprintf("\033[1;32m%s\033[0m", content)
			case 301, 302:
				content = fmt.Sprintf("\033[1;37m%s\033[0m", content)
			case 304:
				content = fmt.Sprintf("\033[1;33m%s\033[0m", content)
			case 401, 403:
				content = fmt.Sprintf("\033[4;31m%s\033[0m", content)
			case 404:
				content = fmt.Sprintf("\033[1;31m%s\033[0m", content)
			case 500:
				content = fmt.Sprintf("\033[1;36m%s\033[0m", content)
			}
		}
		log4go.Debug(content)
	}
}
