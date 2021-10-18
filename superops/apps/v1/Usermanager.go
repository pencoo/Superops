package v1

import (
	"errors"
	"fmt"
	"superops/libs/config"
	"superops/libs/e"
	"superops/libs/rcache"
	casbin_auth "superops/middlewares/casbin-auth"
	L "superops/middlewares/ginzap"
	jwt_auth "superops/middlewares/jwt-auth"
	"superops/modules"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary 用户登录接口
// @Description 可以使用格式:
// @Description    用户名+密码
// @Description    手机号+密码
// @Description    邮箱+密码
// @Description 例如：{"username":"","password":""}
// @Tags Users
// @Accept json
// @Produce json
// @Param modules.Users body modules.Users true "username、password用于登录验证"
// @Success 200 {object} Response
// @Router /api/v1/login [post]
func UserLogin(c *gin.Context) {
	var uinfo modules.Users
	if err := c.ShouldBindJSON(&uinfo); err != nil {
		L.Lzap.Error("数据绑定失败", zap.String("path", c.Request.URL.Path), zap.Reflect("data", uinfo), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if uinfo.Password != "" && (uinfo.Username != "" || uinfo.Phone != "" || uinfo.Email != "") {
			uinfo.Password = Hx(uinfo.Password)
			if uinfo.LoginAuth() {
				t, err := jwt_auth.CreateToken(jwt_auth.CustomClaims{Id: uinfo.Id, Name: uinfo.Username})
				if err != nil {
					L.Lzap.Error("用户信息验证成功，token生成失败", zap.String("path", "/api/v1/login"), zap.Reflect("data", uinfo))
					c.JSON(200, Response{e.FAILED_AUTH_TOKEN, nil, e.GetMsg(e.FAILED_AUTH_TOKEN)})
				} else {
					var u modules.UsersInfo
					u.Username, u.Phone, u.Email, u.Autograph, u.Context = uinfo.Username, uinfo.Phone, uinfo.Email, uinfo.Autograph, uinfo.Context
					_ = rcache.RedisCacheStringAdd(rcache.Userinfo_token+uinfo.Username, t)
					if config.Tokentimeout == 0 {
						_ = rcache.RedisKeyExpireAt(rcache.Userinfo_token+uinfo.Username, time.Unix(jwt_auth.WeekEnd(), 0))
					} else {
						fmt.Println("超时过期时间：", config.Tokentimeout)
						_ = rcache.RedisKeyExpire(rcache.Userinfo_token+uinfo.Username, config.Tokentimeout)
					}
					L.Lzap.Info("用户登录成功", zap.String("path", "/api/v1/login"), zap.String("Access-Token", t), zap.Reflect("data", u))
					c.JSON(200, Response{e.SUCCESS, map[string]interface{}{"Access-Token": t, "userinfo": u}, e.GetMsg(e.SUCCESS)})
				}
			} else {
				L.Lzap.Error("登录验证失败", zap.String("path", "/api/v1/login"), zap.Reflect("data", uinfo))
				//L.Errlog(L.ErrorInfo{Info: "数据绑定失败", Do: "新增云服务商或机房信息", Err: err, User: c.GetString("Name"), Path: c.Request.URL.Path, Data: uinfo})
				c.JSON(200, Response{e.FAILED_AUTH, nil, e.GetMsg(e.FAILED_AUTH)})
			}
		} else {
			L.Lzap.Error("用户登录有必须字段为空", zap.String("path", "/api/v1/login"), zap.Reflect("data", uinfo))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		}
	}
}

// @Summary 用户注销登录接口
// @Description 用户注销：此功能只允许用户自己操作
// @Tags Users
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/users/logout [get]
func UsersLogout(c *gin.Context) {
	u := c.GetString("Name")
	if u != "" {
		_ = rcache.RedisCacheStringDel(rcache.Userinfo_token + u)
		c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 用户注册接口
// @Description  必须参数：用户名、密码、手机、邮箱
// @Tags Users
// @Produce json
// @Param modules.Users body modules.Users true "用户注册"
// @Success 200 {object} Response
// @Router /api/v1/register [post]
func UsersRegister(c *gin.Context) {
	var uinfo modules.Users
	if err := c.ShouldBindJSON(&uinfo); err != nil {
		L.Lzap.Error("用户登录数据绑定失败", zap.String("path", "/api/v1/register"), zap.Reflect("data", uinfo), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if uinfo.Username == "" || uinfo.Password == "" || uinfo.Phone == "" || uinfo.Email == "" {
			L.Lzap.Error("有必须字段为空", zap.String("path", "/api/v1/register"), zap.Reflect("data", uinfo))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		} else {
			uinfo.Id, uinfo.Ustat = 0, 0
			uinfo.Password = Hx(uinfo.Password)
			if err := uinfo.UsersRegister(); err != nil {
				L.Lzap.Error("数据库数据插入失败", zap.String("path", "/api/v1/register"), zap.Reflect("data", uinfo), zap.Error(err))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				L.Lzap.Info("用户注册成功", zap.String("path", "/api/v1/register"), zap.Reflect("data", uinfo))
				c.JSON(200, Response{e.SUCCESS, "", "用户注册成功，请通知管理员审核账号"})
			}
		}
	}
}

// @Summary 用户信息修改
// @Description 修改用户信息，可修改：签名、手机号、邮箱、描述信息
// @Tags Users
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.UsersInfo body modules.UsersInfo true "修改用户信息"
// @Success 200 {object} Response "ok" "返回用户信息"
// @Router /api/v1/users/chageinfo [post]
func UsersUpdateInfo(c *gin.Context) {
	var uinfo modules.UsersInfo
	if err := c.ShouldBindJSON(&uinfo); err != nil {
		L.Lzap.Error("修改用户信息失败，用户信息绑定失败", zap.String("path", "/api/v1/users/chageinfo"), zap.String("user", c.GetString("Name")), zap.Error(err), zap.Reflect("data", uinfo))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if uinfo.Username == "" && uinfo.Phone == "" && uinfo.Email == "" && uinfo.Context == "" {
			L.Lzap.Error("提交信息不能为空", zap.String("path", "/api/v1/users/chageinfo"), zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
		} else {
			var u modules.Users
			u.Id, u.Autograph, u.Email, u.Phone, u.Context = c.GetInt("Id"), uinfo.Autograph, uinfo.Email, uinfo.Phone, uinfo.Context
			if err := u.UsersInfoUpdate(); err != nil {
				L.Lzap.Error("修改用户信息失败", zap.String("path", "/api/v1/users/chageinfo"), zap.String("user", c.GetString("Namw")), zap.Reflect("data", uinfo))
				c.JSON(200, Response{e.FAILED_UPDATE_DB, nil, e.GetMsg(e.FAILED_UPDATE_DB)})
			} else {
				L.Lzap.Info("用户信息修改成功", zap.String("path", "/api/v1/users/chageinfo"), zap.String("user", c.GetString("Namw")), zap.Reflect("data", uinfo))
				c.JSON(200, Response{e.SUCCESS, uinfo, e.GetMsg(e.SUCCESS)})
			}
		}
	}
}

// @Summary 用户修改密码
// @Description 修改用户密码，直接在url参数中提交密码即可，例如：password=xxxxx
// @Description 用户修改密码后需要重新登录
// @Tags Users
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param password query string true "password参数"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/chagespass [post]
func UsersChangePassword(c *gin.Context) {
	pass := c.Query("password")
	if pass == "" {
		L.Lzap.Error("输入密码为空，密码修改失败", zap.String("path", "/api/v1/users/changepass"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		var uinfo modules.Users
		uinfo.Id, uinfo.Password = c.GetInt("Id"), Hx(pass)
		if err := uinfo.UsersInfoUpdate(); err != nil {
			L.Lzap.Error("密码修改失败", zap.String("path", "/api/v1/users/changepass"), zap.String("user", c.GetString("Name")), zap.Error(err))
			c.JSON(200, Response{e.FAILED_UPDATE_DB, err, e.GetMsg(e.FAILED_UPDATE_DB)})
		} else {
			_ = rcache.RedisCacheStringDel(rcache.Userinfo_token + c.GetString("Name"))
			L.Lzap.Info("用户密码修改成功", zap.String("path", "/api/v1/users/changepass"), zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.SUCCESS, "密码修改成功，请使用新密码登录", e.GetMsg(e.SUCCESS)})
		}
	}
}

// @Summary 获取用户信息
// @Description 获取用户信息接口: get /api/v1/users/userinfo
// @Tags Users
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "返回用户信息"
// @Router /api/v1/users/userinfo [get]
func UsersGetUserInfo(c *gin.Context) {
	var uinfo modules.Users
	uinfo.Id = c.GetInt("Id")
	u, err := uinfo.UsersInfos()
	if err != nil {
		L.Lzap.Error("用户信息获取失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/userinfo"), zap.Error(err))
		c.JSON(200, Response{e.ERROR, nil, fmt.Sprint(err)})
	} else {
		L.Lzap.Info("用户信息获取成功", zap.String("path", "/api/v1/users/userinfo"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.SUCCESS, u, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags Users
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/getlists [get]
func GetUsersLists(c *gin.Context) {
	var uinfo modules.Users
	ulist := uinfo.GetUsersLists(1)
	if ulist != nil {
		L.Lzap.Info("用户信息列表获取成功", zap.String("path", "/api/v1/users/getlists"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.SUCCESS, ulist, e.GetMsg(e.SUCCESS)})
	} else {
		L.Lzap.Error("用户信息列表获取失败", zap.String("path", "/api/v1/users/getlists"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.ERROR, nil, e.GetMsg(e.ERROR)})
	}
}

// @Summary 用户创建API Token
// @Description 请求格式：GET /api/v1/users/createapitoken
// @Tags Users
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/createapitoken [get]
func UsersCreateApiToken(c *gin.Context) {
	t, err := jwt_auth.CreateTokenForFuture(jwt_auth.CustomClaims{Id: c.GetInt("Id"), Name: c.GetString("Name")})
	if err != nil {
		L.Lzap.Error("Api Token申请失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/createapitoken"), zap.Error(err))
		c.JSON(200, Response{e.ERROR, nil, "Token生成失败，请重新申请"})
	} else {
		var uak modules.Userapitoken
		uak.Username, uak.Token = c.GetString("Name"), t
		if err := uak.UsersAddApiToken(); err == nil {
			_ = rcache.RedisCacheStringAdd(rcache.Userinfo_api_token+uak.Username, uak.Token)
			L.Lzap.Info("Api Token申请成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/createapitoken"))
			c.JSON(200, Response{e.SUCCESS, t, e.GetMsg(e.SUCCESS)})
		} else if err == errors.New("exist token") {
			L.Lzap.Error("用户api token已经存在无法创建", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/createapitoken"))
			c.JSON(200, Response{e.ERROR, nil, "用户api token已经存在无法创建"})
		} else {
			L.Lzap.Error("新增token写入数据库失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/createapitoken"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
		}
	}
}

// @Summary 用户删除API Token
// @Description 请求格式：GET /api/v1/users/deleteapitoken
// @Tags Users
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/deleteapitoken [get]
func UsersDeleteApiToken(c *gin.Context) {
	var uak modules.Userapitoken
	uak.Username = c.GetString("Name")
	if err := uak.UsersDeleteApiToken(); err == nil {
		_ = rcache.RedisCacheStringDel(rcache.Userinfo_api_token + uak.Username)
		L.Lzap.Info("用户删除API Token成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/deleteapitoken"))
		c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
	} else {
		L.Lzap.Error("api token删除失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/deleteapitoken"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_INSERT_DB, err, e.GetMsg(e.FAILED_INSERT_DB)})
	}
}

// @Summary 用户开启API Token
// @Description 请求格式：GET /api/v1/users/enableapitoken
// @Tags Users
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/enableapitoken [get]
func UsersEnableApiToken(c *gin.Context) {
	var ut modules.Userapitoken
	ut.Username = c.GetString("Name")
	if ut.UserApiTokenExist() {
		if err := ut.UsersEnableApiToken(); err == nil {
			_ = rcache.RedisCacheStringAdd(rcache.Userinfo_api_token+ut.Username, ut.Token)
			L.Lzap.Info("用户API Token成功启用", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/enableapitoken"))
			c.JSON(200, Response{e.SUCCESS, ut.Token, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Error("用户API Token启动失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/enableapitoken"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_INSERT_DB, err, e.GetMsg(e.FAILED_INSERT_DB)})
		}
	} else {
		L.Lzap.Error("用户无API Token", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/enableapitoken"))
		c.JSON(200, Response{e.ERROR, nil, e.GetMsg(e.ERROR)})
	}
}

// @Summary 用户停止API Token
// @Description 请求格式：GET /api/v1/users/disableapitoken
// @Tags Users
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/disableapitoken [get]
func UsersDisableApiToken(c *gin.Context) {
	var ut modules.Userapitoken
	ut.Username = c.GetString("Name")
	if ut.UserApiTokenExist() {
		if err := ut.UsersDisableApiToken(); err == nil {
			_ = rcache.RedisCacheStringDel(rcache.Userinfo_api_token + ut.Username)
			L.Lzap.Info("用户API Token成功禁用", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/disableapitoken"))
			c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Error("用户API Token禁用失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/disableapitoken"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_INSERT_DB, err, e.GetMsg(e.FAILED_INSERT_DB)})
		}
	} else {
		L.Lzap.Error("用户无API Token", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/disableapitoken"))
		c.JSON(200, Response{e.ERROR, nil, e.GetMsg(e.ERROR)})
	}
}

// @Summary 用户获取API Token
// @Description 请求格式：GET /api/v1/users/getapitoken
// @Tags Users
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param file formData file true "file"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/getapitoken [get]
func UsersGetApiToken(c *gin.Context) {
	var u modules.Userapitoken
	u.Username = c.GetString("Name")
	if err := u.UsersGetApiToken(); err == nil {
		if u.Token != "" {
			L.Lzap.Info("用户token获取成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/getapitoken"))
			c.JSON(200, Response{e.SUCCESS, u.Token, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Error("用户token为空", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/getapitoken"))
			c.JSON(200, Response{e.SUCCESS, "用户无Api Token", e.GetMsg(e.SUCCESS)})
		}
	} else {
		L.Lzap.Error("用户获取API token失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/getapitoken"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_SELECT_DB, err, e.GetMsg(e.FAILED_SELECT_DB)})
	}
}

//管理员相关接口

// @Summary 获取url列表
// @Description 获取url列表
// @Tags UsersAdmin
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "返回url列表"
// @Router /api/v1/users/geturllist [get]
func GetUrlListInfo(c *gin.Context) {
	c.JSON(200, Response{e.SUCCESS, Urllist, e.GetMsg(e.SUCCESS)})
}

// @Summary 获取url - method列表
// @Description 返回格式：[{url,method},{url,method},{url,method}...]
// @Tags UsersAdmin
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "返回url列表"
// @Router /api/v1/users/geturlmethodlist [get]
func GetUrlMethodListInfo(c *gin.Context) {
	c.JSON(200, Response{e.SUCCESS, UrlMethodlist, e.GetMsg(e.SUCCESS)})
}

type UserlistSelectType struct {
	Stype   int    `binding:"omitempty,gte=0,lte=3" json:"stype"`
	Lines   int    `binding:"omitempty,gte=0,lte=100" json:"lines"`
	Offset  int    `binding:"omitempty,gt=0" json:"offset"`
	OrderBy string `binding:"omitempty,min=1,max=30" json:"orderby"`
}

// @Summary 获取用户列表
// @Description 请求格式：{"stype":0,"lines":20,"offset":1,"orderby":"id"}
// @Description stype为查询类型：0、查询新注册用户,1、查询激活用户，2、查询禁用用户，3、查询所有用户
// @Description lines为查询条数，默认20条，最多不能查询100条
// @Description offset查询页码，默认查询第1页
// @Description orderby查询排序列，默认为id
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param UserlistSelectType body UserlistSelectType true "查询用户列表参数"
// @Success 200 {object} Response "ok" "返回用户信息"
// @Router /api/v1/users/admingetlist [post]
func UsersGetUserlistByAdmin(c *gin.Context) {
	var Ut UserlistSelectType
	if err := c.ShouldBindJSON(&Ut); err != nil {
		L.Lzap.Error("数据查询绑定失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/admingetlist"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		var u modules.Users
		ul, err := u.GetUsersListsByAdmin(Ut.Stype, Ut.Lines, Ut.Offset, Ut.OrderBy)
		if err != nil {
			L.Lzap.Error("数据查询失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/admingetlist"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_SELECT_DB, err, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			L.Lzap.Info("数据查询成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/admingetlist"), zap.Reflect("data", ul))
			c.JSON(200, Response{e.SUCCESS, ul, e.GetMsg(e.SUCCESS)})
		}
	}
}

// @Summary 单用户启用/禁用
// @Description 启用和禁用单个用户
// @Description 格式：GET /api/v1/users/controls?user=xxx&do=1
// @Description user：输入用户名
// @Description do：操作，1表示启用，2表示禁用
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param user query string true "启用/禁用用户名"
// @Param do query string true "启用：1，禁用：2"
// @Success 200 {object} Response "ok" "返回用户信息"
// @Router /api/v1/users/controls [get]
func UsersOpenOrOff(c *gin.Context) {
	u, d := c.Query("user"), c.Query("do")
	if u != "" && (d == "1" || d == "2") {
		var uinfo modules.Users
		uinfo.Username = u
		var ok string
		if d == "1" {
			_ = rcache.RedisCacheSetAdd(rcache.User_enable_list, uinfo.Username)
			_ = rcache.RedisCacheSetDel(rcache.User_disable_list, uinfo.Username)
			if i := uinfo.UserDoActivateOrDisable(0); i == "" {
				ok = ""
			} else {
				ok = i
			}
		} else {
			_ = rcache.RedisCacheSetAdd(rcache.User_disable_list, uinfo.Username)
			_ = rcache.RedisCacheSetDel(rcache.User_enable_list, uinfo.Username)
			if i := uinfo.UserDoActivateOrDisable(1); i == "" {
				ok = ""
			} else {
				ok = i
			}
		}
		if ok == "" {
			L.Lzap.Info("用户状态操作成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/controls"), zap.String("douser", u))
			c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Error("用户状态操作失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/controls"), zap.String("douser", u), zap.String("error", ok))
			c.JSON(200, Response{e.ERROR, ok, e.GetMsg(e.ERROR)})
		}
	} else {
		L.Lzap.Error("数据绑定失败", zap.String("user", c.GetString("Name")), zap.String("douser", u), zap.String("do", d), zap.String("path", "/api/v1/users/controls"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	}
}

type UsersStatusChange struct {
	DoEnable  []string `json:"doenable"`
	DoDisable []string `json:"dodisable"`
}

// @Summary 批量用户启用/禁用
// @Description 用户批量启用和禁用
// @Description 格式：POST /api/v1/users/controls {"doenable":["x","x","x"],"dodisable":["x","x","x"]}
// @Description doenable：激活用户列表
// @Description dodisable：禁用用户列表
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param UsersStatusChange query UsersStatusChange true "启用/禁用列表"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/controls [post]
func UsersGroupsOpenOrOff(c *gin.Context) {
	var ul UsersStatusChange
	if err := c.ShouldBindJSON(&ul); err != nil {
		L.Lzap.Error("数据绑定失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/controls"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		var fu []string
		var um modules.Users
		for _, u := range ul.DoEnable {
			um.Username = u
			_ = rcache.RedisCacheSetAdd(rcache.User_enable_list, u)
			_ = rcache.RedisCacheSetDel(rcache.User_disable_list, u)
			_ = rcache.RedisCacheStringDel(rcache.Userinfo_token + u)
			if um.UserDoActivateOrDisable(0) != "" {
				fu = append(fu, um.Username)
			}
		}
		for _, u := range ul.DoDisable {
			um.Username = u
			_ = rcache.RedisCacheSetAdd(rcache.User_disable_list, u)
			_ = rcache.RedisCacheSetDel(rcache.User_enable_list, u)
			_ = rcache.RedisCacheStringDel(rcache.Userinfo_token + u)
			if um.UserDoActivateOrDisable(1) != "" {
				fu = append(fu, um.Username)
			}
		}
		L.Lzap.Info("用户批量启用/禁用操作", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/users/controls"), zap.Reflect("data", ul), zap.Reflect("failed", fu))
		c.JSON(200, Response{e.SUCCESS, map[string]interface{}{"failed": fu}, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 添加角色权限
// @Description 添加角色或给角色添加权限
// @Description 格式：{"role":"dev","path":"/api/v1/roleadd","method":"POST"}
// @Description role：角色
// @Description path：url接口
// @Description method：调用方法
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param RM body RM true "body参数"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/roleadd [post]
func UsersRoleAdd(c *gin.Context) {
	var rm RM
	if err := c.ShouldBindJSON(&rm); err == nil {
		if rm.Role != "" && rm.Path != "" && CheckMethod(rm.Method) {
			if casbin_auth.Enforcer.AddPolicy(rm.Role, rm.Path, rm.Method) {
				L.Lzap.Info("角色权限添加成功", zap.String("user", c.GetString("Name")), zap.Reflect("data", rm))
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			} else {
				L.Lzap.Error("角色权限添加失败", zap.String("user", c.GetString("Name")), zap.Reflect("data", rm))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			}
		} else {
			L.Lzap.Error("角色信息不完整，添加失败", zap.String("user", c.GetString("Name")), zap.Reflect("data", rm))
			c.JSON(200, Response{e.FAILED_PARAMS_NUM, nil, e.GetMsg(e.FAILED_PARAMS_NUM)})
		}
	} else {
		L.Lzap.Error("角色权限添加失败，数据绑定失败", zap.String("user", c.GetString("Name")), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	}
}

// @Summary 删除角色权限
// @Description 格式：post /api/v1/users/roleDelPower {"role":"dev","path":"/api/v1/userlogin","method":"POST"}
// @Description role：角色
// @Description path：url接口
// @Description method：调用方法
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param RM body RM true "body参数"
// @Success 200 {object} Response "ok" "返回用户信息"
// @Router /api/v1/users/roleDelPower [post]
func UsersRoleDelPower(c *gin.Context) {
	//删除角色的一条权限规则
	var rm RM
	if err := c.ShouldBindJSON(&rm); err == nil {
		if rm.Role != "" && rm.Path != "" && CheckMethod(rm.Method) {
			if casbin_auth.Enforcer.RemovePolicy(rm.Role, rm.Path, rm.Method) {
				L.Lzap.Info("角色权限删除成功", zap.String("user", c.GetString("Name")), zap.Reflect("data", rm))
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			} else {
				L.Lzap.Info("角色权限删除失败", zap.String("user", c.GetString("Name")), zap.Reflect("data", rm))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			}
		} else {
			L.Lzap.Error("角色信息不完整，删除失败", zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.FAILED_PARAMS_NUM, nil, e.GetMsg(e.FAILED_PARAMS_NUM)})
		}
	} else {
		L.Lzap.Error("角色权限删除失败，数据绑定失败", zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	}
}

// @Summary 获取角色列表
// @Description 获取权限角色列表
// @Tags UsersAdmin
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/rolelist [post]
func UsersGetRoleList(c *gin.Context) {
	//查询系统所有角色
	r := casbin_auth.Enforcer.GetAllRoles()
	c.JSON(200, Response{200, r, ""})
}

// @Summary 删除角色
// @Description 格式：GET /api/v1/users/roledel?role=xxx
// @Tags UsersAdmin
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param role query string true "角色名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/roledel [get]
func UsersRoleDel(c *gin.Context) {
	//删除角色和用户拥有的此角色
	r := c.Query("role")
	casbin_auth.Enforcer.DeleteRole(r)
	c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
}

// @Summary 获取角色权限列表
// @Description 格式：get /api/v1/users/rolePower?role=xxx
// @Tags Users
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param role query string true "角色名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/rolePower [post]
func UsersRolePower(c *gin.Context) {
	//查询一个角色拥有哪些权限
	var RolePowers []RM
	r := c.Query("role")
	rs := casbin_auth.Enforcer.GetPolicy()
	for _, v := range rs {
		if v[0] == r {
			a := RM{Role: v[0], Path: v[1], Method: v[2]}
			RolePowers = append(RolePowers, a)
		}
	}
	c.JSON(200, Response{e.SUCCESS, RolePowers, e.GetMsg(e.SUCCESS)})
}

// @Summary 角色添加用户
// @Description 格式：GET /api/v1/users/roleadduser?role=xxx&user=xxx
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param role query string true "角色"
// @Param user query string true "用户"
// @Success 200 {object} Response "ok" "返回用户信息"
// @Router /api/v1/users/roleadduser [get]
func UsersAddToRole(c *gin.Context) {
	//给用户添加一个角色
	r := c.Query("role")
	u := c.Query("user")
	casbin_auth.Enforcer.AddRoleForUser(u, r)
	c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
}

// @Summary 角色删除用户
// @Description 格式：GET /api/v1/users/roledeluser?role=xxx&user=xxx
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param role query string true "角色"
// @Param user query string true "用户"
// @Success 200 {object} Response "ok" "返回用户信息"
// @Router /api/v1/users/roledeluser [get]
func UsersDelToRole(c *gin.Context) {
	//给用hu删除一个角色
	r := c.Query("role")
	u := c.Query("user")
	casbin_auth.Enforcer.DeleteRoleForUser(u, r)
	c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
}

// @Summary 获取用户角色列表
// @Description 格式：GET /api/v1/users/rolelistbyuser?user=xxx
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param user query string true "用户"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/rolelistbyuser [get]
func UsersGetRole(c *gin.Context) {
	//获取用户拥有的角色
	u := c.Query("user")
	rl, _ := casbin_auth.Enforcer.GetRolesForUser(u)
	c.JSON(200, Response{e.SUCCESS, rl, e.GetMsg(e.SUCCESS)})
}

// @Summary 获取用户权限列表
// @Description 格式：GET /api/v1/users/powerbyuser?user=xxx
// @Tags UsersAdmin
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param user query string true "用户"
// @Success 200 {object} Response "ok"
// @Router /api/v1/users/powerbyuser [get]
func UsersGetUserPower(c *gin.Context) {
	//获取用户拥有的权限
	u := c.Query("user")
	if u == "" {
		u = c.GetString("Name")
	}
	rl, _ := casbin_auth.Enforcer.GetRolesForUser(u)
	var RolePowers []RM
	if rl != nil {
		rs := casbin_auth.Enforcer.GetPolicy()
		for _, r := range rl {
			for _, v := range rs {
				if v[0] == r {
					a := RM{Role: v[0], Path: v[1], Method: v[2]}
					RolePowers = append(RolePowers, a)
				}
			}
		}
		c.JSON(200, Response{e.SUCCESS, RolePowers, e.GetMsg(e.SUCCESS)})
	}
}
