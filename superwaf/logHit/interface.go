package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
}

// @Summary waf与nginx健康检查
// @Description 格式：get /api/v1/health
// @Accept json
// @Produce json
// @Success 200 {object} Response "ok"
// @Router /api/v1/health [get]
func HealthCheck(c *gin.Context) {
	n, err := http.Get(NgxUrl + "/healthcheck")
	if err == nil && n.StatusCode == 200 {
		_ = n.Body.Close()
		_ = RedisCacheSetAdd(RWafClient, LockIpPort)
		_ = RedisCacheSetAdd(RWafNginx, NgxUrl)
		c.JSON(200, Response{200, "ok"})
	} else {
		_ = RedisCacheSetDel(RWafNginx, NgxUrl)
		c.JSON(444, Response{444, "nginx failed"})
	}
}

func NginxStart(c *gin.Context) {
	go WafSync()
	c.JSON(200, Response{200, "ok"})
}

// @Summary 添加列表元素
// @Description 格式：POST /api/v1/addlistpoint
// @Description 不能为空字段：key,value
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/addlistpoint [POST]
func AddListPoint(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		if H.Key == "" || H.Value == "" {
			c.JSON(200, Response{400, "有必须字段为空"})
		} else {
			for _, j := range RedisKeyList {
				if j == H.Key {
					k := keyInfo{Url: NgxUrl, Keyinfo: H}
					if k.AddListInfo() {
						c.JSON(200, Response{200, "ok"})
					} else {
						c.JSON(200, Response{400, "数据增加失败"})
					}
					break
				}
			}
			c.JSON(200, Response{400, "key不在查询失败"})
		}
	}
}

// @Summary 列表中删除元素
// @Description 格式：POST /api/v1/dellistpoint
// @Description 不能为空字段：key,value
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/dellistpoint [POST]
func DelListPoint(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil || H.Key == "" || H.Value == "" {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		for _, j := range RedisKeyList {
			if j == H.Key {
				k := keyInfo{Url: NgxUrl, Keyinfo: H}
				if k.DelListInfo() {
					c.JSON(200, Response{200, "ok"})
				} else {
					c.JSON(200, Response{400, "数据删除失败"})
				}
				break
			}
		}
		c.JSON(200, Response{400, "key不在查询失败"})
	}
}

// @Summary 获取列表元素
// @Description 格式：POST /api/v1/getlistpoint
// @Description 不能为空字段：key
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/getlistpoint [POST]
func GetList(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil || H.Key == "" {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		for _, j := range RedisKeyList {
			if j == H.Key {
				k := keyInfo{Url: NgxUrl, Keyinfo: H}
				if s, err := k.GetListInfo(); err == nil && s != "" {
					c.JSON(200, Response{200, s})
				} else {
					c.JSON(200, Response{400, "数据查询失败"})
				}
				break
			}
		}
		c.JSON(200, Response{400, "key不在查询失败"})
	}
}

// @Summary 查询锁状态
// @Description 格式：POST /api/v1/getlockstate
// @Description 不能为空字段：Lock
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/getlockstate [POST]
func GetLockState(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil || H.Lock == "" {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		for _, j := range RedisKeyList {
			if j == H.Key {
				k := keyInfo{Url: NgxUrl, Keyinfo: H}
				if s, err := k.GetListInfo(); err == nil && s != "" {
					c.JSON(200, Response{200, s})
				} else {
					c.JSON(200, Response{400, "数据查询失败"})
				}
				break
			}
		}
		c.JSON(200, Response{400, "key不在查询失败"})
	}
}

// @Summary 开锁
// @Description 格式：POST /api/v1/lockenable
// @Description 不能为空字段：Lock
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/lockenable [POST]
func LockEnable(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil || H.Lock == "" {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		for _, j := range RedisKeyList {
			if j == H.Key {
				k := keyInfo{Url: NgxUrl, Keyinfo: H}
				if k.EnableLock() {
					c.JSON(200, Response{200, "ok"})
				} else {
					c.JSON(200, Response{400, "数据查询失败"})
				}
				break
			}
		}
		c.JSON(200, Response{400, "key不在查询失败"})
	}
}

// @Summary 关锁
// @Description 格式：POST /api/v1/lockdisable
// @Description 不能为空字段：Lock
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/lockdisable [POST]
func LockDisable(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil || H.Lock == "" {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		for _, j := range RedisKeyList {
			if j == H.Key {
				k := keyInfo{Url: NgxUrl, Keyinfo: H}
				if k.DisableLock() {
					c.JSON(200, Response{200, "ok"})
				} else {
					c.JSON(200, Response{400, "数据查询失败"})
				}
				break
			}
		}
		c.JSON(200, Response{400, "key不在查询失败"})
	}
}

// @Summary 获取内存元素
// @Description 格式：POST /api/v1/getmeminfo
// @Description 不能为空字段：key
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/getmeminfo [POST]
func GetMemInfo(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil || H.Key == "" {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		k := keyInfo{Url: NgxUrl, Keyinfo: H}
		if s, err := k.GetMemInfo(); err != nil || s == "" {
			c.JSON(200, Response{400, err})
		} else {
			c.JSON(200, Response{200, s})
		}
	}
}

// @Summary 设置waf地址
// @Description 格式：POST /api/v1/setwafurl
// @Description 不能为空字段：key
// @Accept json
// @Produce json
// @Param Keyinfo body Keyinfo true "相关信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/setwafurl [POST]
func SetWafServerUrl(c *gin.Context) {
	var H Keyinfo
	if err := c.BindJSON(&H); err != nil || H.Key == "" {
		c.JSON(200, Response{Code: 400, Message: err})
	} else {
		WafUrl = H.Key
		c.JSON(200, Response{200, "ok"})
	}
}

// @Summary 获取锁名称列表
// @Description 格式：get /api/v1/getalllock
// @Accept json
// @Produce json
// @Success 200 {object} Response "ok"
// @Router /api/v1/getalllock [get]
func GetAllLock(c *gin.Context) {
	var list []string
	list = getalllock()
	if list != nil {
		c.JSON(200, Response{200, list})
	} else {
		c.JSON(200, Response{400, "数据查询失败"})
	}
}
func getalllock() []string {
	var list []string
	en, _ := RedisCacheSetGetAll("waf_lock_enable_list")
	dis, _ := RedisCacheSetGetAll("waf_lock_disable_list")
	if en != nil || dis != nil {
		list = append(list, en...)
		list = append(list, dis...)
	}
	return list
}

// @Summary 获取名单名称列表
// @Description 格式：get /api/v1/getalllist
// @Accept json
// @Produce json
// @Success 200 {object} Response "ok"
// @Router /api/v1/getalllist [get]
func GetAllList(c *gin.Context) {
	list := getalllist()
	if list != nil {
		c.JSON(200, Response{200, list})
	} else {
		c.JSON(200, Response{400, ""})
	}
}

func getalllist() []string {
	var list []string
	for _, v := range RedisKeyList {
		if v != "" {
			list = append(list, v)
		}
	}
	return list
}

// @Summary 从redis同步规则到waf
// @Description 格式：get /api/v1/wafrulesync
// @Accept json
// @Produce json
// @Success 200 {object} Response "ok"
// @Router /api/v1/wafrulesync [get]
func RedisToWafSync(c *gin.Context) {
	go WafSync()
	c.JSON(200, Response{200, "ok"})
}

func WafSync() {
	Rinfo := strings.Split(RedisInfo, ":")
	url := NgxUrl + "/redis?ip=" + Rinfo[0] + "&port=" + Rinfo[1]
	_, _ = http.Get(url)
	_, _ = http.Get(NgxUrl + "/cacherebuild")
}
