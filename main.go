package main

import (
	"runtime"
	_ "github.com/yikeso/gomacaron/models"
	"log"
	"net/http"
	"github.com/yikeso/gomacaron/routers"
)

func main(){
	//设置并发使用的核数 为cup核数
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Println("service is running...")
	m := routers.GetRouters()
	log.Println(http.ListenAndServe("0.0.0.0:8080",m))
}
