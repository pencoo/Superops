package main

import (
	"io/ioutil"
	_ "logHit/docs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	sf "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

// @title Superwaf-logHit
// @version 1.0
// @description Superwaf project api document
// @contact.name API Support
// @BasePath
// @author pengbilong email:919182546@qq.com tel:13008811488

func main() {
	InitConfig()
	InitEs()
	RedisInit()
	go NginxActiveCheck(NgxUrl)
	go ElasticDo()
	for _, v := range LogFileList {
		if v != "" {
			go TailFile(v)
		}
	}
	r := gin.Default()
	v1 := r.Group("/api/v1")
	{
		v1.GET("/swagger/*any", gs.WrapHandler(sf.Handler))
		v1.GET("/health", HealthCheck)
		v1.GET("/start", NginxStart)
		v1.POST("/addlistpoint", AddListPoint)
		v1.POST("/dellistpoint", DelListPoint)
		v1.POST("/getlistpoint", GetList)
		v1.POST("/getlockstate", GetLockState)
		v1.POST("/lockenable", LockEnable)
		v1.POST("/lockdisable", LockDisable)
		v1.POST("/getmeminfo", GetMemInfo)
		v1.POST("/setwafurl", SetWafServerUrl)
		v1.GET("/getalllock", GetAllLock)
		v1.GET("/getalllist", GetAllList)
		v1.GET("/wafrulesync", RedisToWafSync)
	}
	r.Run(LockIpPort)
}

func NginxActiveCheck(url string) {
	_ = RedisCacheSetAdd(RWafClient, LockIpPort)
	var Count = 0
	for {
		resp, err := http.Get(url + "/lockstate?lock=NgStat")
		if err != nil {
			if Count <= 2 {
				loginfo := WafLog{DTime: time.Now().Format("2006-01-02 15:04:05"), Client: LockIp, Type: "Error", Message: "nginx stopping..."}
				loginfo.WafDoLog()
				_ = RedisCacheSetDel(RWafNginx, NgxUrl)
				time.Sleep(time.Second * 60)
			} else {
				Count++
				time.Sleep(time.Second * 5)
			}
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			if string(body) == "off\n" {
				l := keyInfo{Url: NgxUrl, Path: "lockenable"}
				l.Lock = "NgStat"
				if l.EnableLock() {
					loginfo := WafLog{DTime: time.Now().Format("2006-01-02 15:04:05"), Client: LockIp, Type: "Error", Message: "nginx start failed"}
					loginfo.WafDoLog()
				} else {
					//go NginxRestartInItWafConfig()
					go WafSync()
					_ = RedisCacheSetAdd(RWafClient, LockIpPort)
					_ = RedisCacheSetAdd(RWafNginx, NgxUrl)
					time.Sleep(time.Second * 5)
				}
			} else {
				_ = resp.Body.Close()
				Count = 0
				time.Sleep(time.Second * 60)
			}
		}
	}
}
