package v1

import (
	"errors"
	"superops/libs/e"
	"superops/libs/rcache"
	L "superops/middlewares/ginzap"
	"superops/modules"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// @Summary 添加harbor
// @Description 格式：GET /api/v1/service/createharbor
// @Description 必须字段：Harbor(harbor名称)、Hurl(harbor地址)、Huser(harbor用户)、Hpass(harbor密码)
// @Description 可选字段：Dimgtimeout(开发镜像保留的时间默认7天，0天表示用户过期需要手动清理harbor，单位是天)
// @Description 可选字段：Timgtimeout(测试镜像保留的时间默认30天，0天表示用户过期需要手动清理harbor，单位是天)
// @Description 可选字段：Timgtimeout(生产镜像保留的时间默认90天，0天表示用户过期需要手动清理harbor，单位是天)
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Harbors body modules.Harbors true "Harbors参数"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/createharbor [post]
func ServiceCreateHarbor(c *gin.Context) {
	var j modules.Harbors
	if err := c.BindJSON(&j); err != nil {
		//数据解析失败
		L.Lzap.Error("添加harbor数据解析失败", zap.Error(err), zap.String("path", "/api/v1/service/createharbor"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		jt := j
		if (j.Harbor == "" || j.Hurl == "" || j.Huser == "" || j.Hpass == "") || jt.HarborInfo() {
			//必须字段有空值
			L.Lzap.Error("添加harbor失败,有必须字段为空或harbor已存在", zap.Reflect("data", j), zap.Error(errors.New("又必须字段为空或harbor已存在")), zap.String("path", "/api/v1/service/createharbor"), zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		} else {
			//缓存中查询是否harbor已存在
			if rcache.RedisCacheSetIn(rcache.Harbor_enable_list, j.Harbor) || rcache.RedisCacheSetIn("harbor_disable_list", j.Harbor) {
				//harbor已存在
				L.Lzap.Error("添加harbor失败，项目在缓存中", zap.Reflect("data", j), zap.String("path", "/api/v1/service/createharbor"), zap.String("user", c.GetString("Name")))
				c.JSON(200, Response{e.ERROR, nil, "harbor已存在，请确认harbor名称及harbor是否已禁用"})
			} else {
				if j.HarborCreate() {
					//数据写入成功
					//加入缓存
					_ = rcache.RedisCacheSetAdd(rcache.Harbor_enable_list, j.Harbor)
					L.Lzap.Info("harbor添加成功", zap.String("path", "/api/v1/service/createharbor"), zap.String("user", c.GetString("Name")))
					c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
				} else {
					//数据写入失败
					L.Lzap.Error("添加harbor失败，数据插入失败", zap.Reflect("data", j), zap.String("path", "/api/v1/service/createharbor"), zap.String("user", c.GetString("Name")))
					c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
				}
			}
		}
	}
}

// @Summary 删除harbor(未完成，删除前需要查询所有项目都没有使用此harbor)
// @Description 格式：GET /api/v1/service/deleteharbor?harbor=xxx
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param harbor query string true "Harbors名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/deleteharbor [get]
func ServiceDeleteHarbor(c *gin.Context) {
	h := c.Query("harbor")
	if h == "" {
		//必须参数harbor为空
		L.Lzap.Error("删除harbor失败，未获取到需要删除的harbor名称", zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
	} else {
		hb := modules.Harbors{Harbor: h}
		//缓存中查询是否存在harbor
		if rcache.RedisCacheSetIn(rcache.Harbor_enable_list, h) {
			//缓存中存在，需要先判断是否有项目在harbor中
			if rcache.RedisKeyExists(rcache.Harbor_project_list + h) {
				//有项目在缓存中
				pl, _ := rcache.RedisCacheSetGetAll(rcache.Harbor_project_list + h)
				L.Lzap.Error("删除harbor失败，尚有项目在仓库中", zap.Reflect("project", pl), zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
				c.JSON(200, Response{e.ERROR, pl, "仍有项目在harbor中，请先迁移项目在进行删除"})
			} else {
				//无项目在缓存中
				if hb.HarborInfo() {
					if hb.HarborDelete() {
						//缓存处理harbor
						rcache.RedisKeyDel(rcache.Harbor_project_list + h)
						_ = rcache.RedisCacheSetDel(rcache.Harbor_enable_list, h)
						_ = rcache.RedisCacheSetDel(rcache.Harbor_disable_list, h)
						L.Lzap.Info("删除harbor成功", zap.String("harbor", h), zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
						c.JSON(200, Response{e.SUCCESS, h, e.GetMsg(e.SUCCESS)})
					} else {
						L.Lzap.Error("从数据库删除harbor失败", zap.String("harbor", h), zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
						c.JSON(200, Response{e.ERROR, nil, "harbor从数据库删除失败"})
					}
				} else {
					L.Lzap.Error("删除harbor失败，harbor不存在", zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
					c.JSON(200, Response{e.ERROR, nil, "harbor不存在"})
				}
			}
		} else {
			//缓存中没有，需要在数据库中查询并操作，然后更新缓存
			p := modules.Projects{Repository: h}
			pl, err := p.ProjectManageProjectGetListByHarbor()
			if err != nil || pl != nil {
				//有项目在harbor中
				//补全harbor缓存
				_ = rcache.RedisCacheSetAdd(rcache.Harbor_enable_list, h)
				var pls []string
				for _, v := range pl {
					pls = append(pls, v.Project)
					_ = rcache.RedisCacheSetAdd(rcache.Harbor_project_list+h, v.Project)
				}
				L.Lzap.Error("删除harbor失败，有项目在缓存仓库中", zap.Reflect("data", pls), zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
				c.JSON(200, Response{e.ERROR, pls, "有项目在harbor中，请先进行harbor迁移或项目迁移"})
			} else {
				if hb.HarborInfo() {
					if hb.HarborDelete() {
						//缓存处理harbor
						rcache.RedisKeyDel(rcache.Harbor_project_list + h)
						_ = rcache.RedisCacheSetDel(rcache.Harbor_disable_list, h)
						L.Lzap.Info("删除harbor成功", zap.String("harbor", h), zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
						c.JSON(200, Response{e.SUCCESS, h, e.GetMsg(e.SUCCESS)})
					} else {
						L.Lzap.Error("从数据库删除harbor失败", zap.String("harbor", h), zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
						c.JSON(200, Response{e.ERROR, nil, "harbor从数据库删除失败"})
					}
				} else {
					L.Lzap.Error("删除harbor失败，harbor不存在", zap.String("path", "/api/v1/service/deleteharbor"), zap.String("user", c.GetString("Name")))
					c.JSON(200, Response{e.ERROR, nil, "harbor不存在"})
				}
			}
		}
	}
}

// @Summary 修改harbor
// @Description 格式：GET /api/v1/service/Updateharbor
// @Description 必须字段：Harbor(harbor名称)
// @Description 可选字段：Hurl(harbor地址)、Huser(harbor用户)、Hpass(harbor密码)
// @Description 可选字段：Dimgtimeout(开发镜像保留的时间默认7天，0天表示不过期需要手动清理harbor，单位是天)
// @Description 可选字段：Timgtimeout(测试镜像保留的时间默认30天，0天表示不过期需要手动清理harbor，单位是天)
// @Description 可选字段：Timgtimeout(生产镜像保留的时间默认90天，0天表示不过期需要手动清理harbor，单位是天)
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Harbors body modules.Harbors true "Harbors参数"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/Updateharbor [post]
func ServiceUpdateHarbor(c *gin.Context) {
	var h modules.Harbors
	if err := c.BindJSON(&h); err != nil || h.Harbor == "" {
		L.Lzap.Error("数据解析失败", zap.Error(err), zap.Reflect("data", h), zap.String("path", "/api/v1/service/Updateharbor"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		h.Id, h.Ctime = 0, ""
		if h.HarborUpdate() {
			L.Lzap.Info("数据更新成功", zap.String("path", "/api/v1/service/Updateharbor"), zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Error("数据更新失败", zap.Reflect("data", h), zap.String("path", "/api/v1/service/Updateharbor"), zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.FAILED_INSERT_DB, nil, "数据写入失败"})
		}
	}
}

// @Summary 查询harbor名称列表
// @Description 格式：GET /api/v1/service/getharbornamelist
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param harbor query string true "Harbors名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/getharbornamelist [get]
func ServiceHarborGetNameList(c *gin.Context) {
	var h modules.Harbors
	//查询缓存如果结果为空在查询数据库，然后更新缓存
	eh, err1 := rcache.RedisCacheSetGetAll(rcache.Harbor_enable_list)
	dh, err2 := rcache.RedisCacheSetGetAll(rcache.Harbor_disable_list)
	if err1 != nil && err2 != nil {
		hl, err := h.HarborNameList()
		//添加缓存
		for _, v := range hl {
			_ = rcache.RedisCacheSetAdd(rcache.Harbor_enable_list, v)
		}
		if err != nil || hl == nil {
			L.Lzap.Error("harbor列表获取失败", zap.Error(err), zap.String("path", "/api/v1/service/getharbornamelist"), zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.ERROR, hl, e.GetMsg(e.ERROR)})
		} else {
			L.Lzap.Info("harbor列表获取成功", zap.String("path", "/api/v1/service/getharbornamelist"), zap.String("user", c.GetString("Name")))
			c.JSON(200, Response{e.SUCCESS, hl, e.GetMsg(e.SUCCESS)})
		}
	} else {
		if eh != nil {
			if dh != nil {
				eh = append(eh, dh...)
			}
		}
		L.Lzap.Info("harbor列表获取成功", zap.String("path", "/api/v1/service/getharbornamelist"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.SUCCESS, eh, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 查询harbor单条详情
// @Description 格式：GET /api/v1/service/getharborinfo?harbor=xxx
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param harbor query string true "Harbors名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/getharborinfo [get]
func ServiceHarborGetInfo(c *gin.Context) {
	var h modules.Harbors
	h.Harbor = c.Query("harbor")
	//先查询缓存，如果缓存为空在查询数据库，然后更新缓存
	if h.Harbor != "" || h.HarborInfo() {
		L.Lzap.Info("harbor信息获取成功", zap.String("path", "/api/v1/service/getharborinfo"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.SUCCESS, h, e.GetMsg(e.SUCCESS)})
	} else {
		L.Lzap.Error("harbor信息获取失败", zap.Error(errors.New("数据查询失败")), zap.String("path", "/api/v1/service/getharborinfo"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.ERROR, nil, "harbor不存在或者结果获取失败"})
	}
}

// @Summary 查询harbor列表(管理员页面展示表格)
// @Description 格式：GET /api/v1/service/getharborlist
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param harbor query string true "Harbors名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/getharborlist [get]
func ServiceHarborGetList(c *gin.Context) {
	var h modules.Harbors
	//先查询环境如果缓存为空则查询数据库，然后更新缓存
	hl, err := h.HarborListAll()
	if err != nil || hl == nil {
		L.Lzap.Error("harbor列表获取失败", zap.Error(err), zap.String("path", "/api/v1/service/getharborlist"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.ERROR, hl, e.GetMsg(e.ERROR)})
	} else {
		L.Lzap.Info("harbor列表获取成功", zap.String("path", "/api/v1/service/getharborlist"), zap.String("user", c.GetString("Name")))
		c.JSON(200, Response{e.SUCCESS, hl, e.GetMsg(e.SUCCESS)})
	}
}
