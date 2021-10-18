package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hpcloud/tail"
	"github.com/tidwall/gjson"
)

type LogHit struct {
	Timespace string  `json:"timespace"` //访问时间
	ServerIp  string  `json:"serverip"`  //服务器ip
	ClientIp  string  `json:"clientip"`  //客户端ip，先检查xxf，如果xxf不为空，则取client_ip
	Domain    string  `json:"domain"`    //访问域名
	Referer   int     `json:"referer"`   //判断是否为空，1有值，2为空，referer
	AppTime   float64 `json:"apptime"`   //应用处理时间
	NgTime    float64 `json:"ngtime"`    //nginx处理时间
	AllTime   float64 `json:"alltime"`   //请求总耗时
	ReqSize   int     `json:"reqsize"`   //请求大小
	Size      int     `json:"size"`      //请求响应大小
	AppIp     string  `json:"appip"`     //应用ip地址，如果ip包含127.0.0.1则记录ServerIp，否则取IP
	Path      string  `json:"path"`      //请求url
	Agent     string  `json:"agent"`     //客户端agent
	Status    int     `json:"status"`    //请求状态
	Method    string  `json:"method"`    //请求method方法
}

//f：文件名
func TailFile(f string) {
	tailFile, err := tail.TailFile(f, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		fmt.Println("tail file err:", err)
		return
	}
	for {
		msg, ok := <-tailFile.Lines
		if !ok {
			fmt.Printf("tail file close reopen, filename: %s\n", tailFile.Filename)
			time.Sleep(1 * time.Second)
			continue
		}
		message := LogProcess(msg.Text)
		if message.Method != "" {
			EsChan <- message
			go message.NginxDoLog()
		}
	}
}

func LogProcess(log string) LogHit {
	defer func() {
		recover()
	}()
	var Log LogHit
	Log.Timespace = gjson.Get(log, "@timestamp").Str
	Log.ServerIp = gjson.Get(log, "server_ip").Str
	Log.ClientIp = gjson.Get(log, "xff").Str
	if Log.ClientIp == "" {
		Log.ClientIp = gjson.Get(log, "client_ip").Str
	}
	Log.Domain = gjson.Get(log, "domain").Str
	if gjson.Get(log, "referer").Str != "" {
		Log.Referer = 1
	} else {
		Log.Referer = 2
	}
	Log.AppTime, _ = strconv.ParseFloat(gjson.Get(log, "upstreamtime").Str, 64)
	Log.AllTime, _ = strconv.ParseFloat(gjson.Get(log, "responsetime").Str, 64)
	Log.NgTime = Log.AllTime - Log.AppTime
	Log.Size, _ = strconv.Atoi(gjson.Get(log, "size").Str)
	Log.ReqSize, _ = strconv.Atoi(gjson.Get(log, "request_length").Str)
	if strings.Contains(gjson.Get(log, "upstreamhost").Str, "127.0.0.1") {
		Log.AppIp = Log.ServerIp
	} else {
		Log.AppIp = strings.Split(gjson.Get(log, "upstreamhost").Str, ":")[0]
	}
	if gjson.Get(log, "request").Str == "" {
		Log.Path = "/unknown"
	} else {
		Log.Path = strings.Split(gjson.Get(log, "request").Str, " ")[1]
	}
	Log.Agent = gjson.Get(log, "http_user_agent").Str
	Log.Status, _ = strconv.Atoi(gjson.Get(log, "status").Str)
	Log.Method = gjson.Get(log, "request_method").Str
	return Log
}

func (l *LogHit) NginxDoLog() {
	client := &http.Client{}
	req, err := json.Marshal(l)
	if err != nil {
		fmt.Println("log json to byte failed, error:", err)
	}
	request, _ := http.NewRequest("POST", WafUrl+"/api/v1/nginxlog", bytes.NewBuffer(req))
	_, err = client.Do(request)
	if err != nil {
		fmt.Println("nginx do log failed, err: ", err)
	}
}
