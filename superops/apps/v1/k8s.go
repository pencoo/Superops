package v1

import (
	"fmt"
	"os"
	"superops/libs/config"
	"superops/libs/e"
	L "superops/middlewares/ginzap"
	"superops/modules"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// @Summary 添加kubeconfig文件
// @Description 格式：GET /api/v1/cmdb/createkubeconfig
// @Description 不可为空字段：kname、kenvrole、kconfig、comname
// @Description hcomname获取接口：/api/v1/commands/getcommandsnamelist
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.KubernetesConfigs body modules.KubernetesConfigs true "kubeconfig参数"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/createkubeconfig [post]
func CmdbCreateKubeconfig(c *gin.Context) {
	var kc modules.KubernetesConfigs
	if err := c.ShouldBindJSON(&kc); err != nil {
		L.Lzap.Error("数据绑定失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/createkubeconfig"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if kc.Kname == "" || kc.KenvRole == 0 || kc.Kconfig == "" || kc.Comname == "" {
			L.Lzap.Error("有必须参数为空", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/createkubeconfig"), zap.Reflect("data", kc))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		} else {
			kc.Confpath = config.FilePath + "/config/" + kc.Kname
			if err := kc.KubeconfigCreate(); err != nil {
				L.Lzap.Error("新增kubeconfig数据库写入失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/createkubeconfig"), zap.Error(err))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				NotExistCreatePath(kc.Confpath)
				f, _ := os.Open(kc.Confpath)
				defer f.Close()
				_, err := f.WriteString(kc.Kconfig)
				if err != nil {
					L.Lzap.Error("数据库写入成功，文件写入失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/createkubeconfig"), zap.Error(err))
					c.JSON(200, Response{e.ERROR, nil, "文件写入失败"})
				} else {
					L.Lzap.Info("配置添加成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/createkubeconfig"))
					c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
				}
			}
		}
	}
}

// @Summary 获取kubeconfig表格(用于管理页面表格展示)
// @Description 格式：GET /api/v1/cmdb/getkubeconfiglist
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/getkubeconfiglist [get]
func CmdbGetKubeConfigList(c *gin.Context) {
	list, err := modules.KubeconfigGetList()
	if err != nil {
		L.Lzap.Error("数据库查询失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/getkubeconfiglist"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
	} else {
		if list != nil {
			L.Lzap.Info("kubeconfig列表获取成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/getkubeconfiglist"))
			c.JSON(200, Response{e.SUCCESS, list, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Info("kubeconfig列表为空", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/getkubeconfiglist"))
			c.JSON(200, Response{e.SUCCESS, nil, "kubeconfig列表为空"})
		}
	}
}

// @Summary 获取kubeconfig名称列表(用于新增/修改项目选择k8s集群)功能未完成，待续...
// @Description 格式：GET /api/v1/cmdb/getkubeconfiglistname?project=xxx
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param project query string true "项目名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/getkubeconfiglistname [get]
func GetKubeConfigListName(c *gin.Context) {
	//查询用户项目角色r,查询可部署
	//u := c.GetInt("Id")
	//p := c.Query("project")
}

// @Summary 获取kubeconfig文件详情(单条详情)
// @Description 格式：GET /api/v1/cmdb/getkubeconfig?kubeconfig=xxx
// @Description kubeconfig: 配置名称
// @Description 名称列表获取接口：GET /api/v1/cmdb/getkubeconfiglist
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param kubeconfig query string true "kubeconfig名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/getkubeconfig [get]
func CmdbGetKubeConfig(c *gin.Context) {
	k := c.Query("kubeconfig")
	if k == "" {
		L.Lzap.Error("请求数据为空", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/getkubeconfig"))
		c.JSON(200, Response{e.ERROR, nil, "请求数据为空"})
	} else {
		kc := modules.KubernetesConfigs{Kname: k}
		if err := kc.KubeconfigGetInfo(); err != nil {
			L.Lzap.Error("数据查询失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/getkubeconfig"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			L.Lzap.Info("查询kubeconfig数据成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/getkubeconfig"))
			c.JSON(200, Response{e.SUCCESS, kc, e.GetMsg(e.SUCCESS)})
		}
	}
}

// @Summary 删除kubeconfig文件
// @Description 格式：GET /api/v1/cmdb/deletekubeconfig?kubeconfig=xxx
// @Description kubeconfig: 配置名称
// @Description 获取kubeconfig列表接口：/api/v1/cmdb/getkubeconfiglist
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param kubeconfig query string true "kubeconfig名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/deletekubeconfig [get]
func CmdbDeleteKubeConfig(c *gin.Context) {
	k := c.Query("kubeconfig")
	if k == "" {
		L.Lzap.Error("删除kubeconfig请求数据为空", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/deletekubeconfig"))
		c.JSON(200, Response{e.ERROR, nil, "请求数据为空"})
	} else {
		kc := modules.KubernetesConfigs{Kname: k}
		_ = kc.KubeconfigGetInfo()
		dpath := kc.Confpath
		if err := kc.KubeconfigDelete(); err != nil {
			L.Lzap.Error("删除kubeconfig失败", zap.String("data", k), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/deletekubeconfig"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			err := os.Remove(dpath)
			if err != nil {
				L.Lzap.Info("删除kubeconfig数据成功,命令删除失败", zap.String("command", config.FilePath+"/config/"+kc.Kname), zap.String("data", k), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/deletekubeconfig"))
				c.JSON(200, Response{e.ERROR, nil, "数据删除成功，命令删除失败(请手动删除)。命令：" + config.FilePath + "/config/" + kc.Kname})
			} else {
				L.Lzap.Info("删除kubeconfig数据成功", zap.String("data", k), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/deletekubeconfig"))
				c.JSON(200, Response{e.SUCCESS, kc, e.GetMsg(e.SUCCESS)})
			}
		}
	}
}

// @Summary 修改kubeconfig文件
// @Description 格式：GET /api/v1/cmdb/updatekubeconfig
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.KubernetesConfigs body modules.KubernetesConfigs true "kubeconfig配置信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/updatekubeconfig [post]
func CmdbUpdateKubeConfig(c *gin.Context) {
	var kc modules.KubernetesConfigs
	if err := c.ShouldBindJSON(&kc); err != nil {
		L.Lzap.Error("数据解析失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/updatekubeconfig"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if kc.KubeconfigExist() {
			kc.Id, kc.KenvRole = 0, 0
			cc := 0
			if kc.Kconfig != "" {
				cc = 1
			}
			if err := kc.KubeconfigUpdate(); err != nil {
				L.Lzap.Error("kubeconfig数据修改失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/updatekubeconfig"))
				c.JSON(200, Response{e.ERROR, nil, fmt.Sprint(err)})
			} else {
				if cc == 1 {
					f, _ := os.Open(kc.Confpath)
					defer f.Close()
					_, _ = f.WriteString(kc.Kconfig)
				}
				_ = kc.KubeconfigGetInfo()
				L.Lzap.Info("kubeconfig数据修改成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/updatekubeconfig"))
				c.JSON(200, Response{e.SUCCESS, kc, e.GetMsg(e.SUCCESS)})
			}
		} else {
			L.Lzap.Error("未查询到此kubeconfig配置", zap.String("kubeconfig", kc.Kname), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/updatekubeconfig"))
			c.JSON(200, Response{e.ERROR, nil, "无此kubeconfig配置"})
		}
	}
}

// @Summary 启用或禁用kubeconfig文件(未完待续)
// @Description 格式：GET /api/v1/cmdb/controlkubeconfig?kubeconfig=xxx&do=enable|disable
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param kubeconfig query string true "kubeconfig名称"
// @Param do query string true "使用enable或disable来启停"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/controlkubeconfig [get]
func CmdbControlKubeConfig(c *gin.Context) {
	kc := c.GetString("kubeconfig")
	d := c.GetString("do")
	if kc == "" || (d == "enable" || d == "disable") {
		//禁用kubeconfig前需要判断是否有项目部署在此集群中
		var do int
		if d == "enable" {
			do = 1
		} else {
			do = 2
		}
		k := modules.KubernetesConfigs{Kname: kc}
		if err := k.KubeconfigControl(do); err != nil {
			L.Lzap.Error("启用或禁用kubeconfig失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/controlkubeconfig"))
			c.JSON(200, Response{e.ERROR, nil, fmt.Sprint(err)})
		} else {
			L.Lzap.Info("启用或禁用kubeconfig成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/controlkubeconfig"))
			c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
		}
	} else {
		L.Lzap.Error("参数解析失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cmdb/controlkubeconfig"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	}
}

// @Summary 后续功能展望
// @Description 需要添加给node加/减tag，node污点、驱逐。项目迁移。(整集群迁移)
// @Tags cmdb
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cmdb/controlkubeconfig [get]
