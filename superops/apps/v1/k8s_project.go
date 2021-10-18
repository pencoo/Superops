package v1

import (
	"strconv"
	"superops/libs/config"
	"superops/libs/e"
	"superops/libs/rcache"
	L "superops/middlewares/ginzap"
	"superops/modules"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary 新增k8s项目
// @Description 格式：POST /api/v1/Project/createk8sproject
// @Description 不可为空字段：project、jenkins、repository
// @Description jenkins获取接口：/api/v1/service/getjenkinsnamelist
// @Description repository获取接口：/api/v1/service/getharbornamelist
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Projects body modules.Projects true "Projects参数"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/createk8sproject [post]
func ProjectManageCreateK8sProject(c *gin.Context) {
	var p modules.Projects
	if err := c.BindJSON(&p); err != nil {
		L.Lzap.Error("新增k8s项目数据解析失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/createk8sproject"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
	} else {
		if p.Project == "" || p.Jenkins == "" || p.Repository == "" {
			L.Lzap.Error("新增k8s项目有必须字段为空，添加失败", zap.Reflect("data", p), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/createk8sproject"))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		} else {
			j := modules.Jenkins{Name: p.Jenkins}
			h := modules.Harbors{Harbor: p.Repository}
			if err := j.ServiceJenkinsInfo(); err != nil || j.State != 1 {
				L.Lzap.Error("新增k8s项目jenkins不存在或已停用，添加失败", zap.Error(err), zap.String("data", p.Jenkins), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/createk8sproject"))
				c.JSON(200, Response{e.ERROR, nil, "jenkins不存在或已停用"})
			} else if !h.HarborInfo() || h.State != 1 {
				L.Lzap.Error("新增k8s项目repository不存在或已停用，添加失败", zap.String("data", p.Repository), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/createk8sproject"))
				c.JSON(200, Response{e.ERROR, nil, "repository不存在或已停用"})
			} else {
				p.Cuser = c.GetString("Name")
				if !rcache.RedisCacheSetIn(rcache.Project_enable_list, p.Project) && !rcache.RedisCacheSetIn(rcache.Project_disable_list, p.Project) {
					L.Lzap.Error("新增k8s项目失败,项目已存在", zap.String("data", p.Project), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/createk8sproject"))
					c.JSON(200, Response{e.ERROR, nil, "项目已存在"})
				} else {
					if err := p.ProjectManageProjectCreate(); err != nil {
						L.Lzap.Error("新增k8s项目数据写入失败", zap.String("data", p.Project), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/createk8sproject"))
						c.JSON(200, Response{e.ERROR, nil, "数据库写入失败"})
					} else {
						_ = rcache.RedisCacheSetAdd(rcache.Project_enable_list, p.Project)
						L.Lzap.Info("新增k8s项目成功", zap.String("data", p.Project), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/createk8sproject"))
						c.JSON(200, Response{e.SUCCESS, nil, "项目添加成功"})
					}
				}
			}
		}
	}
}

// @Summary 更新k8s项目
// @Description 格式：POST /api/v1/Project/updatek8sproject
// @Description 不可为空字段：project
// @Description jenkins获取接口：/api/v1/service/getjenkinsnamelist
// @Description repository获取接口：/api/v1/service/getharbornamelist
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Projects body modules.Projects true "Projects参数"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/updatek8sproject [post]
func ProjectManageUpdateK8sProject(c *gin.Context) {
	var p modules.Projects
	if err := c.BindJSON(&p); err != nil {
		L.Lzap.Error("更新k8s项目数据解析失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/updatek8sproject"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
	} else {
		if p.Project == "" || p.Jenkins == "" || p.Repository == "" {
			L.Lzap.Error("更新k8s项目有必须字段为空，添加失败", zap.Reflect("data", p), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/updatek8sproject"))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		} else {
			j := modules.Jenkins{Name: p.Jenkins}
			h := modules.Harbors{Harbor: p.Repository}
			if err := j.ServiceJenkinsInfo(); err != nil || j.State != 1 {
				L.Lzap.Error("更新k8s项目jenkins不存在或已停用，添加失败", zap.Error(err), zap.String("data", p.Jenkins), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/updatek8sproject"))
				c.JSON(200, Response{e.ERROR, nil, "jenkins不存在或已停用"})
			} else if !h.HarborInfo() || h.State != 1 {
				L.Lzap.Error("更新k8s项目repository不存在或已停用，添加失败", zap.String("data", p.Repository), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/updatek8sproject"))
				c.JSON(200, Response{e.ERROR, nil, "repository不存在或已停用"})
			} else {
				p.Cuser = c.GetString("Name")
				if err := p.ProjectManageProjectUpdate(); err != nil {
					L.Lzap.Error("更新k8s项目数据写入失败", zap.String("data", p.Project), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/updatek8sproject"))
					c.JSON(200, Response{e.ERROR, nil, "数据库写入失败"})
				} else {
					L.Lzap.Info("更新k8s项目成功", zap.String("data", p.Project), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/updatek8sproject"))
					c.JSON(200, Response{e.SUCCESS, nil, "项目更新成功"})
				}
			}
		}
	}
}

// @Summary 激活、禁用k8s项目
// @Description 格式：GET /api/v1/Project/controlk8sproject?project=xxx&do=1|2
// @Description project：项目名称
// @Description do：1表示激活项目，2表示禁用项目
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param project query string true "Projects名称"
// @Param do query string true "操作"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/controlk8sproject [get]
func ProjectManageControlK8sProject(c *gin.Context) {
	p := c.Query("project")
	b := c.Query("do")
	if p == "" || b == "" || !(b == "1" || b == "2") {
		L.Lzap.Error("激活、禁用k8s项目失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/controlk8sproject"))
		c.JSON(200, Response{e.ERROR, nil, "请求参数错误"})
	} else {
		P := modules.Projects{Project: p}
		if (rcache.RedisCacheSetIn(rcache.Project_enable_list, p) || rcache.RedisCacheSetIn(rcache.Project_disable_list, p)) || (P.ProjectManageProjectGetInfo() != nil || P.Id < 0) {
			L.Lzap.Error("激活、禁用k8s项目失败，项目不存在", zap.Reflect("data", P), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/controlk8sproject"))
			c.JSON(200, Response{e.ERROR, nil, "项目不存在"})
		} else {
			B, _ := strconv.Atoi(b)
			if err := P.ProjectManageProjectControl(B); err != nil {
				L.Lzap.Error("激活、禁用k8s项目数据库写入失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/controlk8sproject"))
				c.JSON(200, Response{e.ERROR, nil, "数据写入失败"})
			} else {
				if B == 1 {
					_ = rcache.RedisCacheSetDel(rcache.Project_disable_list, P.Project)
					_ = rcache.RedisCacheSetAdd(rcache.Project_enable_list, P.Project)
				} else if B == 2 {
					_ = rcache.RedisCacheSetDel(rcache.Project_enable_list, P.Project)
					_ = rcache.RedisCacheSetAdd(rcache.Project_disable_list, P.Project)
				}
				L.Lzap.Info("激活、禁用k8s项目成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/controlk8sproject"))
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			}
		}
	}
}

// @Summary 查询k8s项目列表(管理员表格)（此接口会刷新项目列表缓存）
// @Description 格式：GET /api/v1/Project/getk8sprojectlist
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/getk8sprojectlist [get]
func ProjectManageGetK8sProjectList(c *gin.Context) {
	var p modules.Projects
	list, err := p.ProjectManageProjectGetList()
	if err != nil {
		L.Lzap.Error("查询k8s项目列表失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/getk8sprojectlist"))
		c.JSON(200, Response{e.ERROR, err, "数据查询错误"})
	} else {
		for _, v := range list {
			if v.State == 0 {
				_ = rcache.RedisCacheSetDel(rcache.Project_disable_list, v.Project)
				_ = rcache.RedisCacheSetAdd(rcache.Project_enable_list, v.Project)
			} else if v.State == 1 {
				_ = rcache.RedisCacheSetDel(rcache.Project_enable_list, v.Project)
				_ = rcache.RedisCacheSetAdd(rcache.Project_disable_list, v.Project)
			}
		}
		L.Lzap.Info("查询k8s项目列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/getk8sprojectlist"))
		c.JSON(200, Response{e.SUCCESS, list, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 管理员查询k8s项目名称列表（此接口会刷新项目列表缓存）
// @Description 格式：GET /api/v1/Project/getk8sprojectnamelist
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/getk8sprojectnamelist [get]
func ProjectManageGetK8sProjectNameList(c *gin.Context) {
	EnableList, err := rcache.RedisCacheSetGetAll(rcache.Project_enable_list)
	if err != nil {
		var p modules.Projects
		pl1, err := p.ProjectManageProjectGetNameList(1)
		pl2, err := p.ProjectManageProjectGetNameList(2)
		if err != nil || pl1 == nil || pl2 == nil {
			L.Lzap.Error("管理员查询k8s项目名称列表失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/getk8sprojectnamelist"))
			c.JSON(200, Response{e.ERROR, err, "缓存查询失败"})
		} else {
			var R []string
			for _, v := range pl1 {
				R = append(R, v)
				_ = rcache.RedisCacheSetAdd(rcache.Project_enable_list, v)
			}
			for _, v := range pl2 {
				R = append(R, v)
				_ = rcache.RedisCacheSetAdd(rcache.Project_disable_list, v)
			}
			L.Lzap.Info("管理员查询k8s项目名称列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/getk8sprojectnamelist"))
			c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
		}
	} else {
		L.Lzap.Info("管理员查询k8s项目名称列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/getk8sprojectnamelist"))
		c.JSON(200, Response{e.SUCCESS, EnableList, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 用户查询k8s项目名称列表
// @Description 格式：GET /api/v1/Project/usergetk8sprojectnamelist
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/usergetk8sprojectnamelist [get]
func ProjectManageGetK8sProjectNameListByUser(c *gin.Context) {
	EnableList, err := rcache.RedisCacheSetGetAll(rcache.Project_enable_list)
	if err != nil {
		pu := modules.ProjectUsers{Username: c.GetString("Name")}
		pl, err := pu.ProjectManageProjectUsersGetProjectList()
		if err != nil || pl == nil {
			L.Lzap.Error("用户查询k8s项目名称列表失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/usergetk8sprojectnamelist"))
			c.JSON(200, Response{e.ERROR, "", "数据查询失败"})
		} else {
			var R []string
			for _, v := range pl {
				R = append(R, v.Project)
				if v.Urole == 1 {
					_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_dev+v.Project, v.Username)
				} else if v.Urole == 2 {
					_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_test+v.Project, v.Username)
				} else if v.Urole == 3 {
					_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_ops+v.Project, v.Username)
				}
			}
			L.Lzap.Info("用户查询k8s项目名称列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/usergetk8sprojectnamelist"))
			c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
		}
	} else {
		var R []string
		for _, v := range EnableList {
			if rcache.RedisCacheSetIn(rcache.Project_user_list_dev+v, c.GetString("Name")) {
				R = append(R, v)
			} else if rcache.RedisCacheSetIn(rcache.Project_user_list_test+v, c.GetString("Name")) {
				R = append(R, v)
			} else if rcache.RedisCacheSetIn(rcache.Project_user_list_ops+v, c.GetString("Name")) {
				R = append(R, v)
			}
		}
		L.Lzap.Info("用户查询k8s项目名称列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/usergetk8sprojectnamelist"))
		c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 查询k8s项目信息
// @Description 格式：GET /api/v1/Project/selectk8sprojectinfo?project=xxx
// @Description project：项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param project query string true "Projects名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/selectk8sprojectinfo [get]
func ProjectManageGetK8sProjectInfo(c *gin.Context) {
	p := c.Query("project")
	if rcache.RedisCacheSetIn(rcache.Project_enable_list, p) || rcache.RedisCacheSetIn(rcache.Project_disable_list, p) {
		P := modules.Projects{Project: p}
		if err := P.ProjectManageProjectGetInfo(); err != nil {
			L.Lzap.Error("查询k8s项目信息失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/selectk8sprojectinfo"))
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			L.Lzap.Info("查询k8s项目信息成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/selectk8sprojectinfo"))
			c.JSON(200, Response{e.SUCCESS, P, e.GetMsg(e.SUCCESS)})
		}
	} else {
		P := modules.Projects{Project: p}
		if err := P.ProjectManageProjectGetInfo(); err != nil {
			L.Lzap.Error("查询k8s项目信息失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/selectk8sprojectinfo"))
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			if P.State == 1 {
				_ = rcache.RedisCacheSetAdd(rcache.Project_enable_list, P.Project)
			} else {
				_ = rcache.RedisCacheSetAdd(rcache.Project_disable_list, P.Project)
			}
			L.Lzap.Info("查询k8s项目信息成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/selectk8sprojectinfo"))
			c.JSON(200, Response{e.SUCCESS, P, e.GetMsg(e.SUCCESS)})
		}
	}
}

//
type Pcu struct {
	Project string `json:"project"`
	modules.Paul
}

type SFLS struct {
	Success []string `json:"success"`
	Failed  []string `json:"failed"`
}

// @Summary 给k8s项目新增用户
// @Description 格式：POST /api/v1/Project/k8sprojectcreateuser
// @Description 获取用户列表接口：
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param Pcu body Pcu true "Projects新增用户信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectcreateuser [post]
func ProjectManageCreateUserForProject(c *gin.Context) {
	var p Pcu
	if err := c.BindJSON(&p); err != nil || !(p.Project != "" && (p.Ops != nil || p.Dev != nil || p.Test != nil)) {
		L.Lzap.Error("给k8s项目新增用户失败,参数解析失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateuser"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if rcache.RedisCacheSetIn(rcache.Project_enable_list, p.Project) {
			R := SFLS{}
			tmp := modules.ProjectUsers{Project: p.Project}
			if p.Dev != nil {
				tmp.Urole = 1
				for _, v := range p.Dev {
					tmp.Username = v
					err = tmp.ProjectManageProjectUsersCreate()
					if err != nil {
						R.Failed = append(R.Failed, v)
					} else {
						R.Success = append(R.Success, v)
						_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_dev+p.Project, v)
					}
				}
			}
			if p.Test != nil {
				tmp.Urole = 2
				for _, v := range p.Test {
					tmp.Username = v
					err = tmp.ProjectManageProjectUsersCreate()
					if err != nil {
						R.Failed = append(R.Failed, v)
					} else {
						R.Success = append(R.Success, v)
						_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_test+p.Project, v)
					}
				}
			}
			if p.Ops != nil {
				tmp.Urole = 3
				for _, v := range p.Ops {
					tmp.Username = v
					err = tmp.ProjectManageProjectUsersCreate()
					if err != nil {
						R.Failed = append(R.Failed, v)
					} else {
						R.Success = append(R.Success, v)
						_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_ops+p.Project, v)
					}
				}
			}
			L.Lzap.Info("给k8s项目新增用户成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateuser"))
			c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Error("给k8s项目新增用户失败,项目不存在或被禁用", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateuser"))
			c.JSON(200, Response{e.ERROR, nil, "项目不存在或项目已被禁用"})
		}
	}
}

// @Summary 给k8s项目删除用户
// @Description 格式：POST /api/v1/Project/k8sprojectdeleteuser
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param Pcu body Pcu true "Projects删除用户信息名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectdeleteuser [post]
func ProjectManageDeleteUserForProject(c *gin.Context) {
	var p Pcu
	if err := c.BindJSON(&p); err != nil || !(p.Project != "" && (p.Ops != nil || p.Dev != nil || p.Test != nil)) {
		L.Lzap.Error("给k8s项目删除用户失败,参数解析失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectdeleteuser"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if rcache.RedisCacheSetIn(rcache.Project_enable_list, p.Project) {
			R := SFLS{}
			tmp := modules.ProjectUsers{Project: p.Project}
			if p.Dev != nil {
				tmp.Urole = 1
				for _, v := range p.Dev {
					tmp.Username = v
					err = tmp.ProjectManageProjectUsersDelete()
					if err != nil {
						R.Failed = append(R.Failed, v)
					} else {
						R.Success = append(R.Success, v)
						_ = rcache.RedisCacheSetDel(rcache.Project_user_list_dev+p.Project, v)
					}
				}
			}
			if p.Test != nil {
				tmp.Urole = 2
				for _, v := range p.Test {
					tmp.Username = v
					err = tmp.ProjectManageProjectUsersDelete()
					if err != nil {
						R.Failed = append(R.Failed, v)
					} else {
						R.Success = append(R.Success, v)
						_ = rcache.RedisCacheSetDel(rcache.Project_user_list_test+p.Project, v)
					}
				}
			}
			if p.Ops != nil {
				tmp.Urole = 3
				for _, v := range p.Ops {
					tmp.Username = v
					err = tmp.ProjectManageProjectUsersDelete()
					if err != nil {
						R.Failed = append(R.Failed, v)
					} else {
						R.Success = append(R.Success, v)
						_ = rcache.RedisCacheSetDel(rcache.Project_user_list_ops+p.Project, v)
					}
				}
			}
			L.Lzap.Info("给k8s项目删除用户成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectdeleteuser"))
			c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
		} else {
			L.Lzap.Error("给k8s项目删除用户失败,项目不存在或被禁用", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectdeleteuser"))
			c.JSON(200, Response{e.ERROR, nil, "项目不存在或项目已被禁用"})
		}
	}
}

// @Summary 查询k8s项目用户列表
// @Description 格式：GET /api/v1/Project/k8sprojectgetuserslist?project=xxx
// @Description project：项目名称，必须参数
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param Pcu body Pcu true "Projects名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetuserslist [get]
func ProjectManageGetUsersListForProject(c *gin.Context) {
	p := Pcu{Project: c.Query("project")}
	if p.Project == "" {
		//请求项目为空
		L.Lzap.Error("查询项目用户列表失败，项目名称为空", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetuserslist"))
		c.JSON(200, Response{e.ERROR, nil, "项目名称不能为空"})
	} else {
		if rcache.RedisCacheSetIn(rcache.Project_enable_list, p.Project) {
			//项目存在缓存中
			p.Dev, _ = rcache.RedisCacheSetGetAll(rcache.Project_user_list_dev + p.Project)
			p.Test, _ = rcache.RedisCacheSetGetAll(rcache.Project_user_list_test + p.Project)
			p.Ops, _ = rcache.RedisCacheSetGetAll(rcache.Project_user_list_ops + p.Project)
			L.Lzap.Info("查询项目用户列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetuserslist"))
			c.JSON(200, Response{e.SUCCESS, p, e.GetMsg(e.SUCCESS)})
		} else {
			//项目不在缓存中
			err := Rebuildprojectuserlist(p.Project)
			if err == nil {
				p.Dev, _ = rcache.RedisCacheSetGetAll(rcache.Project_user_list_dev + p.Project)
				p.Test, _ = rcache.RedisCacheSetGetAll(rcache.Project_user_list_test + p.Project)
				p.Ops, _ = rcache.RedisCacheSetGetAll(rcache.Project_user_list_ops + p.Project)
				L.Lzap.Info("查询项目用户列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetuserslist"))
				c.JSON(200, Response{e.SUCCESS, p, e.GetMsg(e.SUCCESS)})
			} else {
				L.Lzap.Error("查询项目用户列表失败，缓存更新失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetuserslist"))
				c.JSON(200, Response{e.ERROR, nil, "项目用户信息查询失败"})
			}
		}
	}
}

// @Summary 查询k8s项目用户列表外的用户列表
// @Description 格式：GET /api/v1/Project/k8sprojectoutuserslist?project=xxx
// @Description project：项目名称，必须参数
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param Pcu body Pcu true "Projects名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectoutuserslist [get]
func ProjectManageGetUsersListOurOfProject(c *gin.Context) {
	p := c.Query("project")
	outUser := rcache.Rcli.SDiff(rcache.User_enable_list, rcache.Project_user_list_dev+p, rcache.Project_user_list_test+p, rcache.Project_user_list_ops+p)
	if outUser != nil {
		L.Lzap.Info("查询k8s项目用户列表外的用户列表成功", zap.String("from", "cache"), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetuserslist"))
		c.JSON(200, Response{e.SUCCESS, outUser, e.GetMsg(e.SUCCESS)})
	} else {
		//主动刷新缓存后再次尝试
		P := modules.ProjectUsers{Project: p}
		ul, err := P.ProjectManageGetOutOfProjectUserList()
		if err != nil || ul == nil {
			L.Lzap.Error("查询k8s项目用户列表外的用户列表失败", zap.Error(err), zap.Reflect("data", ul), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetuserslist"))
			c.JSON(200, Response{e.ERROR, err, e.GetMsg(e.ERROR)})
		} else {
			L.Lzap.Info("查询k8s项目用户列表外的用户列表成功", zap.String("from", "db"), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetuserslist"))
			c.JSON(200, Response{e.SUCCESS, ul, e.GetMsg(e.SUCCESS)})
		}
	}
}

// @Summary 新增k8s应用
// @Description 格式：POST /api/v1/Project/k8sprojectcreateapp
// @Description 必须字段：project、appname、langtype、dtmp、giturl、gituser、gitpass、gittoken、route
// @Description langtype获取url：
// @Description dtmp获取url：
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.ProjectApps body modules.ProjectApps true "项目应用信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectcreateapp [post]
func ProjectManageCreatek8sapp(c *gin.Context) {
	var P modules.ProjectApps
	if err := c.BindJSON(&P); err != nil {
		L.Lzap.Error("新增k8s应用失败，数据绑定失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if P.Project == "" || P.Appname == "" || P.Langtype == "" || P.Dtmp == "" || P.Route == "" || P.Gituser == "" || P.Gitpass == "" || P.Giturl == "" || P.Gittoken == "" {
			L.Lzap.Error("新增k8s应用失败，有必须参数为空", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
			c.JSON(200, Response{e.ERROR, nil, "新增k8s应用失败，有必须参数为空"})
		} else {
			if rcache.RedisCacheSetIn(rcache.Project_enable_list, P.Project) {
				if rcache.RedisCacheSetIn(rcache.Project_app_list_enable+P.Project, P.Appname) {
					L.Lzap.Error("新增k8s应用失败，应用已存在", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
					c.JSON(200, Response{e.ERROR, nil, "新增k8s应用失败，应用已存在"})
				} else if rcache.RedisCacheSetIn(rcache.Project_app_list_disable+P.Project, P.Appname) {
					L.Lzap.Error("新增k8s应用失败，应用已存在并被禁用", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
					c.JSON(200, Response{e.ERROR, nil, "新增k8s应用失败，应用已存在并被禁用"})
				} else {
					P.Id, P.Ctime, P.State, P.Username = 0, "", 1, c.GetString("Name")
					if err := P.ProjectManageProjectAppsCreate(); err != nil {
						L.Lzap.Error("新增k8s应用失败，数据库写入失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
						c.JSON(200, Response{e.ERROR, nil, "新增k8s应用失败，数据库写入失败"})
					} else {
						_ = rcache.RedisCacheSetAdd(rcache.Project_app_list_enable+P.Project, P.Appname)
						L.Lzap.Info("新增k8s应用成功", zap.String("app", P.Appname), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
						c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
					}
				}
			} else if rcache.RedisCacheSetIn(rcache.Project_disable_list, P.Project) {
				L.Lzap.Error("新增k8s应用失败，项目已禁用", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
				c.JSON(200, Response{e.ERROR, nil, "新增k8s应用失败，项目已禁用"})
			} else {
				L.Lzap.Error("新增k8s应用失败，项目不存在", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcreateapp"))
				c.JSON(200, Response{e.ERROR, nil, "新增k8s应用失败，项目不存在"})
			}
		}
	}
}

// @Summary 更新k8s应用
// @Description 格式：POST /api/v1/Project/k8sprojectupdateapp
// @Description 必须字段：project、appname
// @Description 不可修改字段：Id、state、ctime、username
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.ProjectApps body modules.ProjectApps true "项目应用信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectupdateapp [post]
func ProjectManageUpdatek8sapp(c *gin.Context) {
	var P modules.ProjectApps
	if err := c.BindJSON(&P); err != nil {
		L.Lzap.Error("更新k8s应用失败，请求参数解析失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectupdateapp"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if rcache.RedisCacheSetIn(rcache.Project_app_list_enable+P.Project, P.Appname) && rcache.RedisCacheSetIn(rcache.Project_enable_list, P.Project) {
			P.Id, P.State, P.Ctime, P.Username = 0, 0, "", ""
			if err := P.ProjectManageProjectAppsUpdate(); err != nil {
				L.Lzap.Error("更新k8s应用失败，数据库写入失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectupdateapp"))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				L.Lzap.Info("更新k8s应用成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectupdateapp"))
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			}
		} else {
			L.Lzap.Error("更新k8s应用失败，项目或应用已被禁用", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectupdateapp"))
			c.JSON(200, Response{e.ERROR, nil, "项目或应用被禁用无法更新"})
		}
	}
}

// @Summary 启用或禁用k8s应用
// @Description 格式：GET /api/v1/Project/k8sprojectcontrolapp?project=xxx&appname=xxx&type=1|2
// @Description 必须字段：project项目名称、appname应用名称、type：1：启用，2：禁用
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectcontrolapp [get]
func ProjectManageControlk8sapp(c *gin.Context) {
	p := c.Query("project")
	a := c.Query("appname")
	do := c.Query("type")
	P := modules.ProjectApps{Project: p, Appname: a}
	if P.Project == "" || P.Appname == "" || !(do == "1" || do == "2") {
		L.Lzap.Error("启用或禁用k8s应用失败，请求参数解析失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcontrolapp"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if do == "1" {
			if rcache.RedisCacheSetIn(rcache.Project_enable_list, P.Project) && rcache.RedisCacheSetIn(rcache.Project_app_list_disable+P.Project, P.Appname) {
				if err := P.ProjectManageProjectAppsControl(1); err != nil {
					L.Lzap.Error("启用或禁用k8s应用失败，数据库写入失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcontrolapp"))
					c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
				} else {
					_ = rcache.RedisCacheSetAdd(rcache.Project_app_list_enable+P.Project, P.Appname)
					_ = rcache.RedisCacheSetDel(rcache.Project_app_list_disable+P.Project, P.Appname)
					L.Lzap.Info("启用或禁用k8s应用成功", zap.String("appname", P.Appname), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcontrolapp"))
					c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
				}
			} else {
				L.Lzap.Error("启用或禁用k8s应用失败，项目已禁用或应用已是激活状态", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcontrolapp"))
				c.JSON(200, Response{e.ERROR, nil, "项目已禁用或应用已是激活状态"})
			}
		} else {
			if rcache.RedisCacheSetIn(rcache.Project_enable_list, P.Project) && rcache.RedisCacheSetIn(rcache.Project_app_list_enable+P.Project, P.Appname) {
				if err := P.ProjectManageProjectAppsControl(2); err != nil {
					L.Lzap.Error("启用或禁用k8s应用失败，数据库写入失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcontrolapp"))
					c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
				} else {
					_ = rcache.RedisCacheSetAdd(rcache.Project_app_list_disable+P.Project, P.Appname)
					_ = rcache.RedisCacheSetDel(rcache.Project_app_list_enable+P.Project, P.Appname)
					L.Lzap.Info("启用或禁用k8s应用成功", zap.String("appname", P.Appname), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcontrolapp"))
					c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
				}
			} else {
				L.Lzap.Error("启用或禁用k8s应用失败，项目已禁用或应用已是禁用状态", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectcontrolapp"))
				c.JSON(200, Response{e.ERROR, nil, "项目已禁用或应用已是禁用状态"})
			}
		}
	}
}

// @Summary 用户获取k8s项目应用列表(项目开发角色可访问)
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]
func ProjectManageGetAppListsByUser(c *gin.Context) {
	P := modules.ProjectApps{Project: c.Query("project")}
	if P.Project == "" {
		L.Lzap.Error("用户获取k8s应用列表失败，请求参数解析失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetapplist"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		u := c.GetString("Name")
		dev, _ := rcache.RedisCacheSetGetAll(rcache.Project_user_list_dev + P.Project)
		if dev != nil || u == config.AdminName {
			applist, err := rcache.RedisCacheSetGetAll(rcache.Project_app_list_enable)
			if err != nil && applist == nil {
				applist, err = P.ProjectManageProjectAppsGetNameList()
				if err != nil {
					L.Lzap.Error("用户获取k8s应用列表失败，数据库查询失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetapplist"))
					c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
				} else {
					for _, v := range applist {
						_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_ops+P.Project, v)
					}
					L.Lzap.Info("用户获取k8s应用列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetapplist"))
					c.JSON(200, Response{e.SUCCESS, applist, e.GetMsg(e.SUCCESS)})
				}
			} else {
				L.Lzap.Info("用户获取k8s应用列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetapplist"))
				c.JSON(200, Response{e.SUCCESS, applist, e.GetMsg(e.SUCCESS)})
			}
		} else {
			L.Lzap.Error("用户获取k8s应用列表失败，用户无查询项目权限", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/Project/k8sprojectgetapplist"))
			c.JSON(200, Response{e.ERROR, nil, "用户无权限查询此项目"})
		}
	}
}

// @Summary 管理员获取k8s项目应用(未完成待续。。。)
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]
func ProjectManageGetAppListsByAdmin(c *gin.Context) {
	//
}

// @Summary 增加k8s项目环境
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 获取k8s项目环境信息
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 更新k8s项目环境信息
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 获取k8s项目环境名称列表
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 管理员获取k8s项目环境名称列表
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 删除k8s项目环境
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 增加k8s项目应用配置
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 删除k8s项目应用配置
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 查看k8s项目应用配置列表
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 拷贝k8s项目应用配置列表(项目配置拷贝)
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 配置k8s项目应用配置键值(修改和增加)
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 查看k8s项目应用配置键值列表
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 部署应用包到项目指定环境
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 获取项目应用环境部署列表
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]

// @Summary 管理员查询包部署环境列表
// @Description 格式：GET /api/v1/Project/k8sprojectgetapplist?project=xxx
// @Description 必须字段：project项目名称
// @Tags Project
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/Project/k8sprojectgetapplist [get]
