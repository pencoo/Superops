package main

import (
	"strings"
	v1 "superops/apps/v1"
	"superops/libs/config"
	"superops/libs/rcache"
	casbin_auth "superops/middlewares/casbin-auth"
	"superops/middlewares/ginzap"
	"superops/modules"
	"superops/routers"
	"time"

	"github.com/gin-gonic/gin"
)

// @title Superops-system
// @version 1.0
// @description Superops project api document
// @contact.name API Support
// @contact.url http://127.0.0.1:8080/api/v1/swagger/index.html
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8080
// @BasePath
// @author pengbilong email:919182546@qq.com tel:13008811488

func main() {
	BeforeInit()
	defer rcache.Rcli.Close() //关闭redis句柄
	defer modules.Db.Close()  //关闭mysql句柄
	defer ginzap.Lzap.Sync()
	defer rcache.RedisCacheSetDel("superops_list_endpoint", config.Endpoint)
	r := gin.Default()
	r.Use(ginzap.Logger(3*time.Second, ginzap.Lzap))
	routers.RouterInit(r)
	rl := r.Routes()
	for _, u := range rl {
		if !strings.Contains(u.Path, "swagger") {
			v1.Urllist = append(v1.Urllist, u.Path)
			v1.UrlMethodlist = append(v1.UrlMethodlist, v1.UrlMethodList{UrlPath: u.Path, Method: u.Method})
		}
	}
	r.Run()
}

func BeforeInit() {
	config.Confinit("./conf.ini")
	ginzap.GinzapInit()
	modules.DBinit()
	rcache.RedisInit()
	casbin_auth.Csbinit()
	if !rcache.RedisCacheSetIn("superops_list_endpoint", config.Endpoint) {
		_ = rcache.RedisCacheSetAdd("superops_list_endpoint", config.Endpoint)
	}
}
