package config

import (
	"crypto/md5"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
)

var (
	Dbinfo       string //数据库连接信息
	TokenSecret  string //jwt加密信息
	TokenTimeout string //token超时时间设置,默认2小时,单位：s(秒)，m(分)，h(时)，d(日)，M(月),y(年)
	Tokentimeout time.Duration
	LogPath      string //日志文件路径，不配置默认输出到控制台
	LogLevel     string //日志输出级别
	Loglevel     = zapcore.DebugLevel
	AdminName    string //管理员账号
	Apidoc       string //是否启用API文档
	Redisurl     string //redis地址，格式：x.x.x.x:xxx，多个地址以逗号分隔
	Redispass    string //redis密码，默认是空
	RedisDb      int    //redis数据库，默认是0
	FilePath     string //命令、配置等文件存放路径
	Timezone     string //系统运行时区
	Endpoint     string //节点名称
)

func Confinit(cf string) {
	configfile, err := ini.Load(cf)
	if err != nil {
		panic("配置文件获取错误，error：" + fmt.Sprint(err))
	}
	Dbinfo = configfile.Section("").Key("Dbinfo").String()
	TokenSecret = configfile.Section("").Key("TokenSecret").String()
	TokenTimeout = configfile.Section("").Key("TokenTimeout").String()
	Tokentimeout = TokentimeoutProcess()
	LogPath = configfile.Section("").Key("LogPath").String()
	LogLevel = configfile.Section("").Key("LogLevel").String()
	AdminName = configfile.Section("").Key("AdminName").String()
	Apidoc = configfile.Section("").Key("Apidoc").String()
	Redisurl = configfile.Section("").Key("Redisurl").String()
	Redispass = configfile.Section("").Key("Redispass").String()
	RedisDb, _ = configfile.Section("").Key("RedisDb").Int()
	FilePath = configfile.Section("").Key("FilePath").String()
	Timezone = configfile.Section("").Key("Timezone").String()
	Endpoint = configfile.Section("").Key("Endpoint").String()

	//处理tokensecret为空的情况
	if TokenSecret == "" {
		TokenSecret = fmt.Sprint(md5.Sum([]byte("pengbilong19891029")))
	}
	//命令和配置文件工作路径
	if FilePath == "" {
		FilePath = "/data"
	}
	//默认时区
	if Timezone == "" {
		Timezone = "Asia/Chongqing"
	}
	//日志基本处理
	if LogLevel == "debug" {
		Loglevel = zapcore.DebugLevel
	} else if LogLevel == "info" {
		Loglevel = zapcore.InfoLevel
	} else if LogLevel == "warn" {
		Loglevel = zapcore.WarnLevel
	} else if LogLevel == "error" {
		Loglevel = zapcore.ErrorLevel
	} else if LogLevel == "dpanic" {
		Loglevel = zapcore.DPanicLevel
	} else if LogLevel == "panic" {
		Loglevel = zapcore.PanicLevel
	} else if LogLevel == "fatal" {
		Loglevel = zapcore.FatalLevel
	}
	//endpoint处理
	if Endpoint == "" {
		var err error
		Endpoint, err = os.Hostname()
		if err != nil {
			Endpoint = fmt.Sprint(md5.Sum([]byte("pengbilong19891029" + time.Now().String())))
		}
	}
}

func TokentimeoutProcess() time.Duration {
	if strings.Contains(TokenTimeout, "s") {
		tmp := strings.Trim(TokenTimeout, "s")
		t, err := strconv.Atoi(tmp)
		if err != nil {
			t = 180
		}
		return time.Second * time.Duration(t)
	} else if strings.Contains(TokenTimeout, "m") {
		tmp := strings.Trim(TokenTimeout, "m")
		t, err := strconv.Atoi(tmp)
		if err != nil {
			t = 30
		}
		return time.Minute * time.Duration(t)
	} else if strings.Contains(TokenTimeout, "h") {
		tmp := strings.Trim(TokenTimeout, "h")
		t, err := strconv.Atoi(tmp)
		if err != nil {
			t = 8
		}
		return time.Hour * time.Duration(t)
	} else if strings.Contains(TokenTimeout, "d") {
		tmp := strings.Trim(TokenTimeout, "d")
		t, err := strconv.Atoi(tmp)
		if err != nil {
			t = 2
		}
		return time.Hour * time.Duration(t*24)
	} else if strings.Contains(TokenTimeout, "M") {
		tmp := strings.Trim(TokenTimeout, "M")
		t, err := strconv.Atoi(tmp)
		if err != nil {
			t = 1
		}
		return time.Hour * time.Duration(t*24*30)
	} else if strings.Contains(TokenTimeout, "y") {
		tmp := strings.Trim(TokenTimeout, "y")
		t, err := strconv.Atoi(tmp)
		if err != nil {
			t = 1
		}
		return time.Hour * time.Duration(t*24*30*365)
	} else {
		return time.Duration(0)
	}
}
