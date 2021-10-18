package main

import (
	"fmt"
	"strconv"

	"gopkg.in/ini.v1"
)

var (
	Esurl     string
	EsIndex   string
	RedisInfo string
	Ip        string
	Dtime     int
	Mysqlinfo string
)
var RsChan = make(chan LogHit, 3000)
var EsChan = make(chan WafLog, 3000)

func init() {
	configfile, err := ini.Load("./conf.ini")
	if err != nil {
		panic("配置文件获取错误，error：" + fmt.Sprint(err))
	}
	Esurl = configfile.Section("").Key("esurl").String()
	EsIndex = configfile.Section("").Key("esindex").String()
	RedisInfo = configfile.Section("").Key("redis").String()
	Ip = configfile.Section("").Key("listen").String()
	cachetime := configfile.Section("").Key("cachetimeout").String()
	Dtime, _ = strconv.Atoi(cachetime)
	Mysqlinfo = configfile.Section("").Key("mysql").String()
}
