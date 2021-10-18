package main

import (
	"fmt"
	"net/http"
	"time"
	_ "waf/docs"

	"github.com/gin-gonic/gin"
	sf "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

var WafLockTag string

func main() {
	WafLockTag = "off"
	DBinit()
	RedisInit()
	InitEs()
	go ElasticDo()
	go wafClientCheck()
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]string{"Esurl": Esurl, "EsIndex": EsIndex, "RedisInfo": RedisInfo, "Ip": Ip})
	})
	v1 := r.Group("/api/v1")
	{
		v1.GET("/swagger/*any", gs.WrapHandler(sf.Handler))
		v1.POST("/waflog", WafLogReceive)              //waf客户端日志接口
		v1.POST("/nginxlog", NginxLogReceive)          //nginx日志接收接口
		v1.GET("/enablelocklist", GetEnableLockList)   //获取当前开启锁列表
		v1.GET("/disablelocklist", GetDisableLockList) //获取当前关闭锁列表
		v1.GET("/enablelock", EnableLock)              //开锁
		v1.GET("/disablelock", DisableLock)            //关锁
		v1.GET("/getlistinfo", GetListInfo)            //获取单个名单内容
		v1.GET("/getalllistinfo", GetAllListInfo)      //获取所有名单列表
		v1.GET("/getlistnamelists", GetListNameLists)  //获取所有名单名称列表
		v1.GET("/updatelistinfo", UpdateListInfo)      //更新名单(会同步刷新缓存和推送waf)
		v1.POST("/addlistinfo", AddListInfo)           //增加名单元素
		v1.POST("/dellistinfo", DelListInfo)           //删除名单元素
		v1.GET("/cachesync", RedisCacheSync)           //名单同步到redis并推送
		v1.GET("/wafsync", WafSync)                    //名单推送
		v1.GET("/wafclientstart", WafClientInit)       //waf客户端启动通知接口
		v1.GET("/wafclientlist", GetWafClientList)     //获取在线waf client列表
		v1.GET("/waflock", WafLockControl)             //waf开关
	}
	r.Run(Ip)
}

func wafClientCheck() {
	for {
		waf, errwaf := Rcli.SMembers(RWafClient).Result()
		if errwaf != nil {
			fmt.Println("waf客户端列表获取失败")
			time.Sleep(time.Minute * 1)
			continue
		}
		if waf != nil {
			for _, v := range waf {
				CheckClient(v)
			}
		}
		time.Sleep(time.Minute * 1)
	}
}

func CheckClient(waf string) {
	if waf == "" {
		fmt.Println("waf client check failed！host is null")
	} else {
		w, errwaf := http.Get("http://" + waf + "/api/v1/health")
		if errwaf == nil {
			if w.StatusCode == 200 {
				_ = w.Body.Close()
			} else if w.StatusCode == 444 {
				_ = w.Body.Close()
				_ = RedisCacheSetDel(RWafClient, waf)
			}
		} else {
			_ = RedisCacheSetDel(RWafClient, waf)
		}
	}
}
