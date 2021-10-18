package main

import (
	"strings"

	"github.com/go-redis/redis"
)

var Rcli *redis.Client

func RedisInit() {
	Rinfo := strings.Split(RedisInfo, ":")
	Rcli = redis.NewClient(&redis.Options{
		Addr:     Rinfo[0] + ":" + Rinfo[1],
		Password: "",
		DB:       0,
	})
	_, err := Rcli.Ping().Result()
	if err == nil {
	} else {
		panic("redis连接失败")
	}
	GetWafUrl()
}

func GetWafUrl() {
	WafUrl, _ = RedisCacheStringGet(RWafUrl)
}

var RedisKeyList = map[string]string{
	"waf_lock_enable_list":  "",
	"waf_lock_disable_list": "",
	"white_url_list":        "WhiteUrl",
	"white_ip_list":         "WhiteIp",
	"bad_method_list":       "BadMethod",
	"bad_agent_list":        "BadAgent",
	"bad_domain_list":       "BadDomain",
	"bad_ip_list":           "BadIp",
	"bad_ip_5min_list":      "BadIp_5",
	"bad_ip_10min_list":     "BadIp_10",
	"bad_ip_60min_list":     "BadIp_60",
	"bad_ip_1_1_list":       "BadIp1_1",
	"bad_ip_1_10_list":      "BadIp1_10",
	"bad_ip_5_10_list":      "BadIp5_10",
	"bad_url_list":          "BadUrl",
	"bad_url_5min_list":     "BadUrl_5",
	"bad_url_10min_list":    "BadUrl_10",
	"bad_url_60min_list":    "BadUrl_60",
	"bad_url_1_1_list":      "BadUrl1_1",
	"bad_url_1_10_list":     "BadUrl1_10",
	"bad_url_5_10_list":     "BadUrl5_10"}

const (
	//redis规则维护
	RWafClient = "waf_client_list"       //waf客户端列表，元素为ip
	RWafNginx  = "waf_enable_nginx_list" //waf活跃nginx列表，元素为ip
	RWafUrl    = "waf_server_url"        //waf地址
	//REnableLock  = "waf_lock_enable_list"  //lock_enable 开锁状态列表，列表中是锁名称
	//RDisableLock = "waf_lock_disable_list" //lock_disable 关锁状态列表，列表中是锁名称
	//RWhiteUrl    = "white_url_list"        //白名单url列表
	//RWhiteIp     = "white_ip_list"         //白名单ip列表
	//RBadMethod   = "bad_method_list"       //method黑名单列表
	//RBadAgent    = "bad_agent_list"        //agent黑名单列表
	//RBadDomain   = "bad_domain_list"       //域名黑名单
	//RBadIp       = "bad_ip_list"           //黑名单ip列表
	//RBadIp5      = "bad_ip_5min_list"      //ip禁用5分钟列表
	//RBadIp10     = "bad_ip_10min_list"     //ip禁用10分钟列表
	//RBadIp60     = "bad_ip_60min_list"     //ip禁用60分钟列表
	//RBadIp1_1    = "bad_ip_1_1_list"       //ip限制每秒1次请求
	//RBadIp1_10   = "bad_ip_1_10_list"      //ip限制10秒1次请求
	//RBadIp5_10   = "bad_ip_5_10_list"      //ip限制10秒5次请求
	//RBadUrl      = "bad_url_list"          //黑名单url列表
	//RBadUrl5     = "bad_url_5min_list"     //url禁用5分钟列表
	//RBadUrl10    = "bad_url_10min_list"    //url禁用10分钟列表
	//RBadUrl60    = "bad_url_60min_list"    //url禁用60分钟列表
	//RBadUrl1_1   = "bad_url_1_1_list"      //url限制每秒1次请求
	//RBadUrl1_10  = "bad_url_1_10_list"     //url限制10秒1次请求
	//RBadUrl5_10  = "bad_url_5_10_list"     //url限制10秒5次请求
)

//nginx重启时初始化nginx规则
func NginxRestartInItWafConfig() {
	//_ = RedisCacheSetAdd(RWafClient, LockIpPort)
	//_ = RedisCacheSetAdd(RWafNginx, LockIp)
	for k, v := range RedisKeyList {
		if v == "" {
			MoreLockControl(k)
		} else {
			MoreListControl(k, v)
		}
	}
}

func MoreLockControl(s string) {
	locklist, _ := RedisCacheSetGetAll(s)
	if locklist != nil {
		var k keyInfo
		k = keyInfo{Url: NgxUrl}
		for _, v := range locklist {
			k.Lock = v
			if strings.Contains(s, "enable") {
				go k.EnableLock()
			} else {
				go k.DisableLock()
			}
		}
	}
}

func MoreListControl(s string, c string) {
	badmethodlist, err := RedisCacheSetGetAll(s)
	if err == nil && badmethodlist != nil {
		k := keyInfo{Url: NgxUrl}
		k.Key = c
		for _, v := range badmethodlist {
			k.Value = v
			go k.AddListInfo()
		}
	}
}

//string类型 查
func RedisCacheStringGet(k string) (string, error) {
	v, err := Rcli.Get(k).Result()
	if err == nil {
		return v, nil
	} else {
		return v, err
	}
}

//set集合类型 增
func RedisCacheSetAdd(k, v string) error {
	_, err := Rcli.SAdd(k, v).Result()
	return err
}

//set集合类型 删
func RedisCacheSetDel(k, v string) error {
	_, err := Rcli.SRem(k, v).Result()
	return err
}

//set集合类型 查询所有元素
func RedisCacheSetGetAll(k string) ([]string, error) {
	return Rcli.SMembers(k).Result()
}
