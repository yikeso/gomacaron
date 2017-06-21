package main

import (
	_ "github.com/yikeso/gomacaron/config"//加载配置文件
	_ "github.com/yikeso/gomacaron/models"//初始化数据源
	_ "github.com/yikeso/gomacaron/util"
	"runtime"
	l4g "github.com/alecthomas/log4go"
	"net/http"
	"github.com/yikeso/gomacaron/routers"
)

func main(){
	//设置并发使用的核数 为cup核数
	runtime.GOMAXPROCS(runtime.NumCPU())
	l4g.Info("service is running...")
	m := routers.GetRouters()
	l4g.Error(http.ListenAndServe("0.0.0.0:8080",m))
}
