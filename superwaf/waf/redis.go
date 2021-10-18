package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

var Rcli *redis.Client

func RedisInit() {
	Rinfo := strings.Split(RedisInfo, ":")
	if len(Rinfo) < 3 {
		Rcli = redis.NewClient(&redis.Options{
			Addr:     Rinfo[0] + ":" + Rinfo[1],
			Password: "",
			DB:       0,
		})
	} else {
		Rcli = redis.NewClient(&redis.Options{
			Addr:     Rinfo[0] + ":" + Rinfo[1],
			Password: Rinfo[2],
			DB:       0,
		})
	}
	_, err := Rcli.Ping().Result()
	if err == nil {
	} else {
		panic("redis连接失败")
	}
	waf, _ := RedisCacheStringGet(RWafUrl)
	if waf == "" || !strings.Contains(waf, Ip) {
		_ = RedisCacheStringAdd(RWafUrl, "http://"+Ip)
	}
	go DelRedisExpireKey()
	go RedisDo()
	go InterfaceTimePolymerization()
	go DelRoleListDoClean()
}

var RedisKeylist = map[string]string{
	"EnableLock":  "waf_lock_enable",    //lock_enable 开锁状态列表，列表中是锁名称
	"DisableLock": "waf_lock_disable",   //lock_disable 关锁状态列表，列表中是锁名称
	"WhiteUrl":    "white_url_list",     //白名单url列表
	"WhiteIp":     "white_ip_list",      //白名单ip列表
	"BadMethod":   "bad_method_list",    //method黑名单列表
	"BadAgent":    "bad_agent_list",     //agent黑名单列表
	"BadDomain":   "bad_domain_list",    //域名黑名单
	"BadIp":       "bad_ip_list",        //黑名单ip列表
	"BadIp_5":     "bad_ip_5min_list",   //ip禁用5分钟列表
	"BadIp_10":    "bad_ip_10min_list",  //ip禁用10分钟列表
	"BadIp_60":    "bad_ip_60min_list",  //ip禁用60分钟列表
	"BadIp1_1":    "bad_ip_1_1_list",    //ip限制每秒1次请求
	"BadIp1_10":   "bad_ip_1_10_list",   //ip限制10秒1次请求
	"BadIp5_10":   "bad_ip_5_10_list",   //ip限制10秒5次请求
	"BadUrl":      "bad_url_list",       //黑名单url列表
	"BadUrl_5":    "bad_url_5min_list",  //url禁用5分钟列表
	"BadUrl_10":   "bad_url_10min_list", //url禁用10分钟列表
	"BadUrl_60":   "bad_url_60min_list", //url禁用60分钟列表
	"BadUrl1_1":   "bad_url_1_1_list",   //url限制每秒1次请求
	"BadUrl1_10":  "bad_url_1_10_list",  //url限制10秒1次请求
	"BadUrl5_10":  "bad_url_5_10_list",  //url限制10秒5次请求
}

const (
	//redis规则维护
	RWafClient  = "waf_client_list"       //waf客户端列表，元素为ip
	RWafNginx   = "waf_enable_nginx_list" //waf活跃nginx列表，元素为ip
	RWafUrl     = "waf_server_url"        //waf地址
	RAUS        = "access_url_0_"         //url记数访问成功,ddd是时间延时到分钟，xxx是url，过期时间为5分钟，ddd_xxx
	RAUF        = "access_url_1_"         //url记数访问失败,ddd是时间延时到分钟，xxx是url，过期时间为5分钟
	RAIS        = "access_ip_0_"          //Ip访问成功记数,ddd是时间延时到分钟，xxx是url，过期时间位5分钟
	RAIF        = "access_ip_1_"          //Ip访问失败记数,ddd是时间延时到分钟，xxx是url，过期时间位5分钟
	RASS        = "access_status_0"       //正确状态码统计，过期时间为5分钟
	RASF        = "access_status_1"       //错误状态码统计，过期时间为5分钟
	RADS        = "access_domain_0_"      //域名访问正确统计，ddd是时间延时到分钟，xxx是url，过期时间为5分钟
	RADF        = "access_domain_1_"      //域名访问错误统计，ddd是时间延时到分钟，xxx是url，过期时间为5分钟
	RAUIF       = "access_err_"           //单ip访问错误统计，过期时间5分钟,ip_url_ddd
	RAUT        = "access_url_time_"      //url访问耗时，ddd是时间延时到分钟，xxx是url，过期时间为30天
	RAUTOK      = "ok_access_url_time_"   //url访问耗时app，ddd是时间延时到分钟，xxx是url，过期时间为30天
	DelRoleList = "delete_role_list"      //规则删除队列
)

//消息写入redis
func RedisDo() {
	for {
		select {
		case ri := <-RsChan:
			RedisInfoInto(ri)
			InfoProc(ri)
		}
	}
}

//redis写入格式处理
func RedisInfoInto(j LogHit) {
	t := time.Now().Format("01021504")
	to := time.Minute * time.Duration(Dtime)
	pipe := Rcli.Pipeline()
	url := t + "_" + j.Path
	domain := t + "_" + j.Domain
	ip := t + "_" + j.ClientIp
	bip := j.ClientIp + "_" + j.Path + "_" + t
	if j.Status < 400 {
		//成功请求
		if RedisKeyExists(RAUS + url) {
			pipe.Incr(RAUS + url)
		} else {
			pipe.Set(RAUS+url, 0, to)
			pipe.Incr(RAUS + url)
		}
		if RedisKeyExists(RADS + domain) {
			pipe.Incr(RADS + domain)
		} else {
			pipe.Set(RADS+domain, 0, to)
			pipe.Incr(RADS + domain)
		}
		if RedisKeyExists(RAIS + ip) {
			pipe.Incr(RAIS + ip)
		} else {
			pipe.Set(RAIS+ip, 0, to)
			pipe.Incr(RAIS + ip)
		}
		if RedisKeyExists(RASS) {
			pipe.Incr(RASS)
		} else {
			pipe.Set(RASS, 0, to)
			pipe.Incr(RASS)
		}
	} else {
		//失败请求
		if RedisKeyExists(RAUF + url) {
			pipe.Incr(RAUF + url)
		} else {
			pipe.Set(RAUF+url, 0, to)
			pipe.Incr(RAUF + url)
		}
		if RedisKeyExists(RADF + domain) {
			pipe.Incr(RADF + domain)
		} else {
			pipe.Set(RADF+domain, 0, to)
			pipe.Incr(RADF + domain)
		}
		if RedisKeyExists(RAIF + ip) {
			pipe.Incr(RAIF + ip)
		} else {
			pipe.Set(RAIF+ip, 0, to)
			pipe.Incr(RAIF + ip)
		}
		if RedisKeyExists(RASF) {
			pipe.Incr(RASF)
		} else {
			pipe.Set(RASF, 0, to)
			pipe.Incr(RASF)
		}
		if RedisKeyExists(RAUIF + bip) {
			pipe.Incr(RAUIF + bip)
		} else {
			pipe.Set(RAUIF+bip, 0, to)
			pipe.Incr(RAUIF + bip)
		}
	}
	if RedisKeyExists(RAUT + url) {
		pipe.RPush(RAUT+url, j.AllTime)
	} else {
		pipe.RPush(RAUT+url, j.AllTime)
		pipe.Expire(RAUT+url, to)
	}
	_, _ = pipe.Exec()
}

//请求判断处理
func InfoProc(i LogHit) {
	//响应时间大于平均5分钟内平均响应时间的2倍时进行报警
	avg := Resposeavgtime(i.Path, 5)
	if avg > 0 && avg*10 < i.AllTime {
		if WafLockTag == "on" {
			var wl WafList
			wl.Listname = "BadUrl_5"
			wl.Listinfo = i.Path
			fmt.Println("平均值：", avg, " 当前响应时间：", i.AllTime, "  ,path drop 5 min, path: ", i.Path)
			go HttpPostDoTest("http://"+Ip+"/api/v1/addlistinfo", wl)
		} else {
			fmt.Println("waf lock未开启", "平均值：", avg, " 当前响应时间：", i.AllTime, "  ,path drop 5 min, path: ", i.Path)
		}
	}
	//1分钟内连续错误5次封禁ip5分钟
	_, f := ResponseSingleIpCountByTime(i.ClientIp, 1)
	if f > 5 {
		if WafLockTag == "on" {
			var wl WafList
			wl.Listname = "BadIp_5"
			wl.Listinfo = i.ClientIp
			fmt.Println("连续访问错误次数：", f, " IP drop 5 min, ip: ", i.ClientIp)
			go HttpPostDoTest("http://"+Ip+"/api/v1/addlistinfo", wl)
		} else {
			fmt.Println("waf lock未开启", "连续访问错误次数：", f, " IP drop 5 min, ip: ", i.ClientIp)
		}
	}
	//rt := Resposeavgtime(i.Path, 5)
	//fmt.Println("url: ", i.Path, ",5分钟内平均响应时间：")
	//s, f := ResponseSingleIpCountByTime(i.ClientIp, 5)
	//fmt.Println("ip: ", i.ClientIp, ",5分钟内访问成功数量：", s, ",失败数量：", f)
	//s, f = ResponseSingleDomainCountByTime(i.Domain, 5)
	//fmt.Println("domain: ", i.Domain, ",5分钟内访问成功数量：", s, ",失败数量：", f)
	//if i.Status >= 400 {
	//	fmt.Println("ip: ", i.ClientIp, "访问url：", i.Path, "最近五分钟平均错误次数: ", SingleIpAccessSingleUrlError(i.ClientIp, i.Path))
	//}
	//if WafLockTag == "on" {
	//	//
	//} else {
	//	if Resposeavgtime(i.Path, 5) > i.AllTime*5 {
	//	}
	//}
}

func HttpPostDoTest(url string, l WafList) {
	client := &http.Client{}
	z := WafList{Listname: l.Listname, Listinfo: l.Listinfo}
	i, _ := json.Marshal(z)
	req_new := bytes.NewBuffer(i)
	request, _ := http.NewRequest("POST", url, req_new)
	request.Header.Set("Content-type", "application/json")
	response, _ := client.Do(request)
	_ = response.Body.Close()
}

//func AccessRoleRead() {
//	var R []WafRole
//	for _, R := range R {
//
//	}
//}

//请一定时间内请求响应平均值,单位分钟
func Resposeavgtime(url string, t int) float64 {
	l := GetTimeDot(t)
	var j float64
	for _, v := range l {
		i, _ := Rcli.Get(RAUTOK + v + "_" + url).Result()
		m, err := strconv.ParseFloat(i, 64)
		if err != nil {
			m = 0.0
			t = t - 1
		}
		j = j + m
	}
	if t > 0 {
		return j / float64(t)
	} else {
		return j
	}
}

//单ip一段时间内请求总数，返回成功数和失败数
func ResponseSingleIpCountByTime(ip string, t int) (int, int) {
	l := GetTimeDot(t)
	var s, f int
	for _, v := range l {
		i, _ := Rcli.Get(RAIS + v + "_" + ip).Result()
		m, _ := strconv.Atoi(i)
		s = s + m
	}
	for _, v := range l {
		i, _ := Rcli.Get(RAIF + v + "_" + ip).Result()
		m, _ := strconv.Atoi(i)
		f = f + m
	}
	return s, f
}

//单ip一段时间内请求总数，返回成功数和失败数
func ResponseSingleDomainCountByTime(domain string, t int) (int, int) {
	l := GetTimeDot(t)
	var s, f int
	for _, v := range l {
		i, _ := Rcli.Get(RADS + v + "_" + domain).Result()
		m, _ := strconv.Atoi(i)
		s = s + m
	}
	for _, v := range l {
		i, _ := Rcli.Get(RADF + v + "_" + domain).Result()
		m, _ := strconv.Atoi(i)
		f = f + m
	}
	return s, f
}

//单ip访问单url错误，五分钟内的每分钟平均值
func SingleIpAccessSingleUrlError(ip string, url string) int {
	i, _ := Rcli.Keys(RAUIF + ip + "_" + url).Result()
	if len(i) > 0 {
		num := len(i)
		count := 0
		for _, v := range i {
			j, _ := Rcli.Get(v).Result()
			m, _ := strconv.Atoi(j)
			if m == 0 {
				num = num - 1
			} else {
				count = count + num
			}
		}
		if num != 0 {
			return count / num
		} else {
			return count
		}
	} else {
		return 0
	}
}

//返回时间切片
func GetTimeDot(t int) []string {
	var r []string
	for i := t; i > 0; i-- {
		r = append(r, time.Now().Add(-time.Duration(i)*time.Minute).Format("01021504"))
	}
	if r == nil {
		r = append(r, time.Now().Format("01021504"))
	}
	return r
}

//处理过期数据
func DelRedisExpireKey() {
	for {
		var Exlist []string
		i := TimeSlotProcess("clean")
		for _, v := range i {
			exlist, _ := Rcli.Keys("*" + v + "*").Result()
			Exlist = append(Exlist, exlist...)
		}
		if Exlist != nil {
			pipe := Rcli.Pipeline()
			for _, v := range Exlist {
				if !strings.Contains(v, RAUTOK) && v != "" {
					pipe.Del(v)
				}
			}
			pipe.Exec()
		}
		Exlist = []string{}
		time.Sleep(time.Duration(Dtime) * time.Minute)
	}
}

//耗时聚合处理
func InterfaceTimePolymerization() {
	for {
		Exlist := []string{}
		exlist, _ := Rcli.Keys(RAUT + "*").Result()
		if exlist != nil {
			t := time.Now().Format("01021504")
			for _, v := range exlist {
				if !strings.Contains(v, t) && v != "" {
					Exlist = append(Exlist, v)
				}
			}
			go Polymerization(Exlist)
		}
		time.Sleep(20 * time.Second)
	}
}

//聚合函数
func Polymerization(keys []string) {
	for _, v := range keys {
		r, _ := Rcli.LLen(v).Result()
		if r > 0 {
			R, _ := Rcli.LRange(v, 0, r-1).Result()
			var L float64
			for _, j := range R {
				l, _ := strconv.ParseFloat(j, 64)
				L = L + l
			}
			Rcli.Set("ok_"+v, L/float64(r), time.Hour*24)
		}
		Rcli.Del(v)
	}
}

//时间段处理函数
func TimeSlotProcess(c string) []string {
	var ctime int
	if c == "clean" {
		ctime = Dtime
	} else {
		ctime = 1
	}
	var R []string
	for i := 0; i < Dtime; i++ {
		R = append(R, time.Now().Add(-time.Minute*time.Duration(ctime+i)).Format("01021504"))
	}
	return R
}

//判断key是否存在
func RedisKeyExists(k string) bool {
	n, err := Rcli.Exists(k).Result()
	if err == nil && n > 0 {
		return true
	} else {
		return false
	}
}

//删除缓存key
func RedisKeyDel(k string) bool {
	n, err := Rcli.Del(k).Result()
	if err == nil && n > 0 {
		return true
	} else {
		return false
	}
}

//string类型操作
//string类型 增
func RedisCacheStringAdd(k, v string) error {
	return Rcli.Set(k, v, 0).Err()
}

//string类型 删
func RedisCacheStringDel(k string) error {
	return Rcli.Del(k).Err()
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

//定期检查过期key
func DelRoleListDoClean() {
	for {
		l, _ := Rcli.SMembers(DelRoleList).Result()
		if l != nil {
			var info DelListRole
			var dolist []DelListRole
			t := time.Now().Unix()
			for _, v := range l {
				_ = json.Unmarshal([]byte(v), &info)
				if info.Dtime <= t {
					dolist = append(dolist, info)
				}
			}
			if dolist != nil {
				go DeleteExpireData(dolist)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

//redis和mysql过期删除
func DeleteExpireData(d []DelListRole) {
	pipe := Rcli.Pipeline()
	for _, v := range d {
		e, _ := json.Marshal(v)
		pipe.SRem(DelRoleList, string(e))
	}
	_, _ = pipe.Exec()
	c := WafList{}
	for _, v := range d {
		c.Listname = v.Listname
		if c.FindWafListInfo() == nil {
			var lts, j []string
			_ = json.Unmarshal([]byte(c.Listinfo), &lts)
			for _, k := range lts {
				if k != v.ListKey {
					j = append(j, k)
				}
			}
			d, _ := json.Marshal(j)
			c.Listinfo = string(d)
			_ = c.UpdateWafListInfo()
		}
	}
}
