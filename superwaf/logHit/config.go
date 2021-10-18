package main

import (
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	WafUrl     string
	NgxUrl     string
	RedisInfo  string
	LockIp     string
	LockIpPort string
)

var LogFileList []string
var EsChan = make(chan LogHit, 3000)
var E EsConf

type EsConf struct {
	EsUrl   string
	EsIndex string
}

func InitConfig() {
	configfile, err := ini.Load("./conf.ini")
	if err != nil {
		panic("配置文件获取错误，error：" + fmt.Sprint(err))
	}
	WafUrl = configfile.Section("").Key("wafurl").String()
	lfl := configfile.Section("").Key("logfile").String()
	LogFileList = strings.Split(lfl, ",")
	E.EsUrl = configfile.Section("").Key("esurl").String()
	E.EsIndex = configfile.Section("").Key("esindex").String()
	NgxUrl = configfile.Section("").Key("ngxurl").String()
	RedisInfo = configfile.Section("").Key("redisinfo").String()
	lockinfo := strings.Split(configfile.Section("").Key("ip").String(), ":")
	LockIp = lockinfo[0]
	LockIpPort = lockinfo[0] + ":" + lockinfo[1]
}
