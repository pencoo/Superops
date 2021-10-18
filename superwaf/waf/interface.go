package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code    int
	Data    interface{}
	Message string
}

//base64编码处理
func AllPassProcessByEncode(pass string) string {
	return base64.StdEncoding.EncodeToString([]byte(pass))
}

//base64解码处理
func AllPassProcessBydecode(pass string) string {
	re, _ := base64.StdEncoding.DecodeString(pass)
	return string(re)
}

// @Summary 客户端启动通知接口
// @Description 格式：get /api/v1/wafclientstart?host=x.x.x.x:port
// @Description host参数不能为空，用于标识本机
// @Accept json
// @Produce json
// @Success 200 {object} Resp "ok"
// @Router /api/v1/wafclientstart [get]
func WafClientInit(c *gin.Context) {
	go RulewafPush(c.Query("host"))
	c.JSON(200, Resp{200, nil, "ok"})
}

// 日志写入es中
// @Summary 设置waf地址
// @Description 格式：POST /api/v1/waflog
// @Description 不能为空字段：key
// @Accept json
// @Produce json
// @Param WafLog body WafLog true "相关信息"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/waflog [POST]
func WafLogReceive(c *gin.Context) {
	var H WafLog
	if c.BindJSON(&H) == nil {
		EsChan <- H
		c.JSON(200, Resp{200, nil, "ok"})
	} else {
		c.JSON(200, Resp{400, nil, "数据解析失败"})
	}
}

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

//nginx日志接口用于接收client发送过来的日志控制台打印(并将之写入redis并启动对此url的分析过程)
func NginxLogReceive(c *gin.Context) {
	var H LogHit
	if c.BindJSON(&H) == nil {
		RsChan <- H
		c.String(200, "")
	} else {
		c.String(400, "数据解析失败")
	}
}

// @Summary 开启状态锁列表
// @Description 格式：GET /api/v1/enablelocklist
// @Accept json
// @Produce json
// @Success 200 {object} Resp "ok"
// @Router /api/v1/enablelocklist [get]
func GetEnableLockList(c *gin.Context) {
	var l WafLock
	L, err := l.FindLockEnableList()
	if err != nil {
		c.JSON(200, Resp{400, err, "数据查询失败"})
	} else {
		c.JSON(200, Resp{200, L, ""})
	}
}

// @Summary 关闭状态锁列表
// @Description 格式：get /api/v1/disablelocklist
// @Accept json
// @Produce json
// @Success 200 {object} Resp "ok"
// @Router /api/v1/disablelocklist [get]
func GetDisableLockList(c *gin.Context) {
	var l WafLock
	L, err := l.FindLockDisableList()
	if err != nil {
		c.JSON(200, Resp{400, err, "数据查询失败"})
	} else {
		c.JSON(200, Resp{200, L, ""})
	}
}

// @Summary 开锁
// @Description 格式：get /api/v1/enablelock?key=
// @Description 不能为空字段：key
// @Accept json
// @Produce json
// @Param key query string true "锁名称"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/enablelock [get]
func EnableLock(c *gin.Context) {
	l := WafLock{Lockname: c.Query("key")}
	if l.Lockname != "" {
		if err := l.EnableLock(); err != nil {
			c.JSON(200, Resp{400, err, ""})
		} else {
			c.JSON(200, Resp{200, nil, "ok"})
		}
	} else {
		c.JSON(200, Resp{400, nil, "请求key为空"})
	}
}

// @Summary 关锁
// @Description 格式：get /api/v1/disablelock?key=
// @Description 不能为空字段：key
// @Accept json
// @Produce json
// @Param key query string true "锁名称"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/disablelock [get]
func DisableLock(c *gin.Context) {
	l := WafLock{Lockname: c.Query("key")}
	if l.Lockname != "" {
		if err := l.DisableLock(); err != nil {
			c.JSON(200, Resp{400, err, ""})
		} else {
			c.JSON(200, Resp{200, nil, "ok"})
		}
	} else {
		c.JSON(200, Resp{400, nil, "请求key为空"})
	}
}

type WafLists struct {
	Listname string
	Listinfo []string
}

// @Summary 获取单个名单列表
// @Description 格式：get /api/v1/getlistinfo?key=
// @Description 不能为空字段：key
// @Accept json
// @Produce json
// @Param key query string true "锁名称"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/getlistinfo [get]
func GetListInfo(c *gin.Context) {
	l := WafList{Listname: c.Query("key")}
	if l.Listname == "" {
		c.JSON(200, Resp{400, nil, "key为空无法查询"})
	} else {
		if err := l.FindWafListInfo(); err != nil {
			c.JSON(200, Resp{400, err, "数据查询失败"})
		} else {
			a := AllPassProcessBydecode(l.Listinfo)
			L := WafLists{Listname: l.Listname}
			_ = json.Unmarshal([]byte(a), &L.Listinfo)
			c.JSON(200, Resp{200, L, "ok"})
		}
	}
}

// @Summary 获取名单名称列表
// @Description 格式：get /api/v1/getlistnamelists
// @Accept json
// @Produce json
// @Success 200 {object} Resp "ok"
// @Router /api/v1/getlistnamelists [get]
func GetListNameLists(c *gin.Context) {
	var l WafList
	L, err := l.FindWafListAllInfo()
	if err != nil {
		c.JSON(200, Resp{400, err, "数据查询错误"})
	} else {
		var R []string
		for _, k := range L {
			if k.Listname != "" {
				R = append(R, k.Listname)
			}
		}
		c.JSON(200, Resp{200, R, "ok"})
	}
}

// @Summary 获取名单列表
// @Description 格式：get /api/v1/getalllistinfo
// @Accept json
// @Produce json
// @Success 200 {object} Resp "ok"
// @Router /api/v1/getalllistinfo [get]
func GetAllListInfo(c *gin.Context) {
	var l WafList
	L, err := l.FindWafListAllInfo()
	if err != nil {
		c.JSON(200, Resp{400, err, "数据查询错误"})
	} else {
		a := []WafLists{}
		for _, v := range L {
			d := AllPassProcessBydecode(v.Listinfo)
			b := WafLists{Listname: v.Listname}
			_ = json.Unmarshal([]byte(d), &b.Listinfo)
			a = append(a, b)
		}
		c.JSON(200, Resp{200, a, "ok"})
	}
}

// @Summary 修改名单列表(会同步刷新缓存和推送waf)
// @Description 格式：POST /api/v1/Updatelistinfo
// @Accept json
// @Produce json
// @Param WafLists body WafLists true "参数"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/Updatelistinfo [post]
func UpdateListInfo(c *gin.Context) {
	var w WafLists
	if err := c.BindJSON(&w); err != nil || w.Listname == "" || w.Listinfo == nil {
		c.JSON(200, Resp{400, nil, "数据解析失败"})
	} else {
		s := WafList{Listname: w.Listname}
		d, err := json.Marshal(w.Listinfo)
		if err != nil {
			c.JSON(200, Resp{400, err, "数据编码失败"})
		} else {
			if string(d) != "" {
				s.Listinfo = AllPassProcessByEncode(string(d))
				if err := s.UpdateWafListInfo(); err != nil {
					c.JSON(200, Resp{400, err, "数据库写入失败"})
				} else {
					go rediscachesync() //同步缓存和推送waf
					c.JSON(200, Resp{200, nil, "ok"})
				}
			} else {
				c.JSON(200, Resp{400, nil, "listinfo 解析失败"})
			}
		}
	}
}

// @Summary 增加名单元素(会同步刷新缓存和推送waf)
// @Description 格式：POST /api/v1/addlistinfo
// @Accept json
// @Produce json
// @Param WafList body WafList true "参数"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/addlistinfo [post]
func AddListInfo(c *gin.Context) {
	var i WafList
	dl := DelListRole{}
	if err := c.BindJSON(&i); err != nil {
		c.JSON(200, Resp{400, err, "数据解析失败"})
	} else {
		if i.Listname != "" && i.Listinfo != "" {
			if i.DataIsInWhiteList(i.Listinfo) {
				c.JSON(200, Resp{400, nil, "数据在白名单中"})
			} else {
				j := i.Listinfo
				i.Listinfo = ""
				if i.FindWafListInfo() != nil {
					c.JSON(200, Resp{400, nil, "查询失败，无此列表"})
				} else {
					if i.Listinfo == "null" {
						i.Listinfo = ""
					}
					if i.Listinfo == "" {
						d := WafLists{Listname: i.Listname}
						d.Listinfo = append(d.Listinfo, j)
						e, _ := json.Marshal(d.Listinfo)
						i.Listinfo = AllPassProcessByEncode(string(e))
						if i.UpdateWafListInfo() == nil {
							dl.Listname = i.Listname
							dl.ListKey = j
							go dl.DataIntoIntoDeleteList()
							f := WafList{Listname: i.Listname, Listinfo: j}
							go WafClientAddRole(f)
							c.JSON(200, Resp{200, nil, "ok"})
						}
					} else {
						jode := AllPassProcessBydecode(i.Listinfo)
						var jo []string
						_ = json.Unmarshal([]byte(jode), &jo)
						if jo != nil {
							jj := 0
							for _, v := range jo {
								if j == v {
									jj = 2
									break
								}
							}
							if jj < 1 {
								jo = append(jo, j)
								d, _ := json.Marshal(jo)
								joen := AllPassProcessByEncode(string(d))
								i.Listinfo = joen
								if i.UpdateWafListInfo() == nil {
									dl.Listname = i.Listname
									dl.ListKey = j
									go dl.DataIntoIntoDeleteList()
									l := WafList{Listname: i.Listname, Listinfo: j}
									go WafClientAddRole(l)
									c.JSON(200, Resp{200, nil, "ok"})
								} else {
									c.JSON(200, Resp{400, nil, "数据更新写入失败"})
								}
							} else {
								c.JSON(200, Resp{400, nil, "数据重复，此数据已存在"})
							}
						}
					}
				}
			}
		} else {
			c.JSON(200, Resp{400, i, "有必须参数为空新增失败"})
		}
	}
}

type DelListRole struct {
	Dtime    int64  //过期时间
	Listname string //过期列表名称
	ListKey  string //过期内容
}

func (dr *DelListRole) DataIntoIntoDeleteList() {
	if dr.Listname == "BadIp_10" || dr.Listname == "BadIp1_1" || dr.Listname == "BadIp1_10" || dr.Listname == "BadIp5_10" || dr.Listname == "BadUrl_10" || dr.Listname == "BadUrl1_1" || dr.Listname == "BadUrl1_10" || dr.Listname == "BadUrl5_10" {
		dr.Dtime = time.Now().Add(time.Duration(10) * time.Minute).Unix()
		i, _ := json.Marshal(dr)
		Rcli.SAdd(DelRoleList, string(i))
	} else if dr.Listname == "BadIp_5" || dr.Listname == "BadUrl_5" {
		dr.Dtime = time.Now().Add(time.Duration(5) * time.Minute).Unix()
		i, _ := json.Marshal(dr)
		Rcli.SAdd(DelRoleList, string(i))
	} else if dr.Listname == "BadIp_60" || dr.Listname == "BadUrl_60" {
		dr.Dtime = time.Now().Add(time.Duration(60) * time.Minute).Unix()
		i, _ := json.Marshal(dr)
		Rcli.SAdd(DelRoleList, string(i))
	}
}

// @Summary 删除名单元素(会同步刷新缓存和推送waf)
// @Description 格式：POST /api/v1/dellistinfo
// @Accept json
// @Produce json
// @Param WafList body WafList true "参数"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/dellistinfo [post]
func DelListInfo(c *gin.Context) {
	var i WafList
	if err := c.BindJSON(&i); err != nil {
		c.JSON(200, Resp{400, err, "数据解析失败"})
	} else {
		if i.Listname != "" && i.Listinfo != "" {
			j := i.Listinfo
			if i.FindWafListInfo() != nil {
				c.JSON(200, Resp{400, nil, "查询失败，无此列表"})
			} else {
				if i.Listinfo == "" {
					c.JSON(200, Resp{400, nil, "此列表为空，无法删除元素"})
				} else {
					jode := AllPassProcessBydecode(i.Listinfo)
					var jo []string
					_ = json.Unmarshal([]byte(jode), &jo)
					if jo != nil {
						cz := 0
						for _, v := range jo {
							if j == v {
								cz = 1
							}
						}
						if cz == 1 {
							var newlist []string
							for _, v := range jo {
								if j != v {
									newlist = append(newlist, v)
								}
							}
							if newlist != nil {
								d, _ := json.Marshal(newlist)
								joen := AllPassProcessByEncode(string(d))
								i.Listinfo = joen
								if i.UpdateWafListInfo() == nil {
									l := WafList{Listname: i.Listname, Listinfo: j}
									go WafClientDelRole(l)
									c.JSON(200, Resp{200, nil, "ok"})
								} else {
									c.JSON(200, Resp{400, nil, "数据更新写入失败"})
								}
							} else {
								i.Listinfo = ""
								if i.UpdateWafListInfo() == nil {
									l := WafList{Listname: i.Listname, Listinfo: j}
									go WafClientDelRole(l)
									c.JSON(200, Resp{200, nil, "ok"})
								} else {
									c.JSON(200, Resp{400, nil, "数据更新写入失败"})
								}
							}
						} else {
							c.JSON(200, Resp{400, nil, "信息不存在，无法删除"})
						}
					} else {
						c.JSON(200, Resp{400, nil, "数据解码失败或数据为空"})
					}
				}
			}
		} else {
			c.JSON(200, Resp{400, nil, "有必须参数为空删除失败"})
		}
	}
}

// @Summary waf缓存同步
// @Description 格式：get /api/v1/cachesync?pushwaf=yes|no
// @Description 参数pushwaf可选，默认为yes，即缓存同步redis后推送至waf。no为不推送
// @Accept json
// @Produce json
// @param pushwaf query string false "可选参数"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/cachesync [get]
func RedisCacheSync(c *gin.Context) {
	d := c.Query("pushwaf")
	if d == "no" || d == "NO" || d == "No" {
		go redisCacheRebuild("no")
	} else {
		go redisCacheRebuild("yes")
	}
	c.JSON(200, Resp{200, nil, "ok"})
}

func redisCacheRebuild(do string) {
	if do == "no" {
		rediscachesync()
	} else {
		rediscachesync()
		RulewafPush("")
	}
}

//数据库读取所有规则到redis缓存
func rediscachesync() {
	var wl WafLock
	lockenablelist, err := wl.FindLockEnableList()
	if err == nil && lockenablelist != nil {
		pipe := Rcli.Pipeline()
		pipe.Del(RedisKeylist["REnableLock"])
		for _, v := range lockenablelist {
			pipe.SAdd(RedisKeylist["REnableLock"], v)
		}
		_, _ = pipe.Exec()
	}
	lockDisablelist, err := wl.FindLockDisableList()
	if err == nil && lockDisablelist != nil {
		pipe := Rcli.Pipeline()
		pipe.Del(RedisKeylist["RDisableLock"])
		for _, v := range lockDisablelist {
			pipe.SAdd(RedisKeylist["RDisableLock"], v)
		}
		_, _ = pipe.Exec()
	}
	var wls WafList
	L, err := wls.FindWafListAllInfo()
	a := []WafLists{}
	if err == nil && L != nil {
		for _, v := range L {
			d := AllPassProcessBydecode(v.Listinfo)
			b := WafLists{Listname: v.Listname}
			_ = json.Unmarshal([]byte(d), &b.Listinfo)
			a = append(a, b)
		}
	}
	if a != nil {
		fmt.Println(a)
		for _, v := range a {
			Rcli.Del(RedisKeylist[v.Listname])
			for _, k := range v.Listinfo {
				Rcli.SAdd(RedisKeylist[v.Listname], k)
			}
		}
	}
}

// @Summary 获取在线waf client列表
// @Description 格式：get /api/v1/wafclientlist
// @Accept json
// @Produce json
// @Success 200 {object} Resp "ok"
// @Router /api/v1/wafclientlist [get]
func GetWafClientList(c *gin.Context) {
	l, err := Rcli.SMembers(RWafClient).Result()
	if err == nil && l != nil {
		var L []string
		for _, v := range l {
			L = append(L, v)
		}
		c.JSON(200, Resp{200, L, ""})
	} else {
		c.JSON(200, Resp{200, nil, "数据为空"})
	}
}

// @Summary waf客户端同步
// @Description 格式：get /api/v1/wafsync?host=xxx
// @Description host非必须参数，有参数时同步指定节点，无参数时同步所有节点
// @Description host列表获取接口：get /api/v1/wafclientlist
// @Accept json
// @Produce json
// @param host query string false "主机"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/wafsync [get]
func WafSync(c *gin.Context) {
	h := c.Query("host")
	if h == "" {
		go RulewafPush("")
		c.JSON(200, Resp{200, nil, "ok"})
	} else {
		l, err := Rcli.SMembers(RWafClient).Result()
		if err == nil && l != nil {
			a := 0
			for _, k := range l {
				if h == k {
					go RulewafPush(h)
					a = a + 1
					c.JSON(200, Resp{200, nil, "ok"})
				}
			}
			if a == 0 {
				c.JSON(200, Resp{200, nil, "未找到此主机"})
			}
		} else {
			c.JSON(200, Resp{200, err, "未找到此主机"})
		}
	}
}

//节点刷新全部缓存
func RulewafPush(host string) {
	if host == "" {
		l, err := Rcli.SMembers(RWafClient).Result()
		if err == nil && l != nil {
			for _, v := range l {
				_, _ = http.Get("http://" + v + "/wafrulesync")
			}
		}
	} else {
		_, _ = http.Get("http://" + host + "/wafrulesync")
	}
}

type Keyinfo struct {
	Key   string //键
	Value string //值
	Lock  string //锁
}

//节点指定列表增加元素
func WafClientAddRole(list WafList) {
	_ = RedisCacheSetAdd(RedisKeylist[list.Listname], list.Listinfo)
	l, err := Rcli.SMembers(RWafClient).Result()
	if err == nil && l != nil {
		for _, v := range l {
			go HttpPostDo("http://"+v+"/api/v1/addlistpoint", list)
		}
	}
}

//节点指定列表删除元素
func WafClientDelRole(list WafList) {
	_ = RedisCacheSetDel(RedisKeylist[list.Listname], list.Listinfo)
	l, err := Rcli.SMembers(RWafClient).Result()
	if err == nil && l != nil {
		for _, v := range l {
			go HttpPostDo("http://"+v+"/api/v1/dellistpoint", list)
		}
	}
}

//http调用
func HttpPostDo(url string, l WafList) {
	client := &http.Client{}
	z := Keyinfo{Key: l.Listname, Value: l.Listinfo}
	i, _ := json.Marshal(z)
	req_new := bytes.NewBuffer(i)
	request, _ := http.NewRequest("POST", url, req_new)
	request.Header.Set("Content-type", "application/json")
	response, _ := client.Do(request)
	_ = response.Body.Close()
}

// @Summary WafLock
// @Description 格式：get /api/v1/waflock?do=enable|disable
// @Description do为空，表示查询waflock状态，do=enable表示开启waf，do=disable表示关闭waf
// @Accept json
// @Produce json
// @param do query string false "控制锁"
// @Success 200 {object} Resp "ok"
// @Router /api/v1/waflock [get]
func WafLockControl(c *gin.Context) {
	do := c.Query("do")
	if do == "" {
		c.JSON(200, Resp{200, WafLockTag, ""})
	} else if do == "enable" {
		WafLockTag = "on"
		c.JSON(200, Resp{200, WafLockTag, ""})
	} else if do == "disable" {
		WafLockTag = "off"
		c.JSON(200, Resp{200, WafLockTag, ""})
	}
}
