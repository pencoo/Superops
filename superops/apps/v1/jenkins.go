package v1

import (
	"strconv"
	"superops/libs/e"
	L "superops/middlewares/ginzap"
	"superops/modules"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// @Summary 添加jenkins
// @Description 格式：GET /api/v1/service/createjenkins
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Jenkins body modules.Jenkins true "jenkins信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/createjenkins [post]
func ServiceCreateJenkins(c *gin.Context) {
	var j modules.Jenkins
	if err := c.BindJSON(&j); err != nil {
		L.Lzap.Error("新增jenkins数据绑定失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/createjenkins"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if j.Name == "" || j.Jurl == "" || j.Juser == "" || j.Jpass == "" {
			L.Lzap.Error("新增jenkins失败，有必须字段为空", zap.Reflect("data", j), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/createjenkins"))
			c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, "又必须字段为空"})
		} else {
			J := j
			err := J.ServiceJenkinsInfo()
			if err != nil || J.Id > 0 {
				L.Lzap.Error("新增jenkins失败，jenkins已存在", zap.Reflect("data", j), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/createjenkins"))
				c.JSON(200, Response{e.ERROR, nil, "jenkins已存在"})
			} else {
				if err := j.ServiceCreateJenkins(); err != nil {
					L.Lzap.Error("新增jenkins失败，数据库写入失败", zap.Reflect("data", j), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/createjenkins"))
					c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
				} else {
					L.Lzap.Info("新增jenkins成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/createjenkins"))
					c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
				}
			}
		}
	}
}

// @Summary 修改jenkins
// @Description 格式：GET /api/v1/service/updatejenkins
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Jenkins body modules.Jenkins true "jenkins信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/updatejenkins [post]
func ServiceUpdateJenkins(c *gin.Context) {
	var j modules.Jenkins
	if err := c.BindJSON(&j); err != nil {
		L.Lzap.Error("修改jenkins数据绑定失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkins"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if j.Name == "" {
			L.Lzap.Error("修改jenkins失败，有必须字段为空", zap.Reflect("data", j), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkins"))
			c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, "又必须字段为空"})
		} else {
			J := j
			err := J.ServiceJenkinsInfo()
			if err != nil || J.Id < 0 {
				L.Lzap.Error("新增jenkins失败，jenkins不存在", zap.Reflect("data", j), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkins"))
				c.JSON(200, Response{e.ERROR, nil, "jenkins不存在"})
			} else {
				if err := j.ServiceUpdateJenkins(); err != nil {
					L.Lzap.Error("修改jenkins失败，数据库写入失败", zap.Reflect("data", j), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkins"))
					c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
				} else {
					L.Lzap.Info("修改jenkins成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkins"))
					c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
				}
			}
		}
	}
}

// @Summary jenkins状态控制(未完待续)
// @Description 格式：GET /api/v1/service/controljenkins?jenkins=xxx&type=1|2|3&to
// @Description 必须字段：
// @Description     jenkins：指定需要修改的jenkins名称
// @Description     type：1表示启用，2表示禁用，3表示迁移(to=xxx必须，用于指定迁移的目标jenkins)
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param jenkins query string true "jenkins名称"
// @Param type query string true "操作"
// @Param to query string true "目标jenkins"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/updatejenkins [post]
func ServiceJenkinsControls(c *gin.Context) {
	j := c.Query("jenkins")
	d := c.Query("type")
	t := c.Query("to")
	if j == "" || !(d == "1" || d == "2" || (d == "3" && t != "")) {
		//条件不足
	} else {
		D, _ := strconv.Atoi(d)
		if D == 3 {
			//迁移jenkins，需要先判断源jenkins和目标jenkins是否存在且目标是否激活可用
			//J := modules.Jenkins{Name: j}
			//T := modules.Jenkins{Name: t}
		} else {
			//启用或禁用，禁用前需要先检查是否有项目在jenkins中
		}
	}
}

// @Summary 查询单条jenkins详情
// @Description 格式：GET /api/v1/service/jenkinsinfo?jenkins=xxx
// @Description jenkins：查询jenkins的名称
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param jenkins query string true "查询的jenkins名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/jenkinsinfo [get]
func ServiceGetJenkinsInfo(c *gin.Context) {
	j := c.Query("jenkins")
	if j == "" {
		L.Lzap.Error("查询jenkins信息失败", zap.String("error", "jenkins name获取失败"), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/jenkinsinfo"))
		c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
	} else {
		J := modules.Jenkins{Name: j}
		err := J.ServiceJenkinsInfo()
		if err != nil {
			L.Lzap.Error("查询jenkins信息失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/jenkinsinfo"))
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			L.Lzap.Info("查询jenkins信息成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/jenkinsinfo"))
			c.JSON(200, Response{e.SUCCESS, J, e.GetMsg(e.SUCCESS)})
		}
	}
}

// @Summary 查询jenkins名称列表
// @Description 格式：GET /api/v1/service/getjenkinsnamelist
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/getjenkinsnamelist [get]
func ServiceGetJenkinsNameList(c *gin.Context) {
	var j modules.Jenkins
	r, err := j.ServiceJenkinsNameLists()
	if err != nil || r == nil {
		L.Lzap.Error("查询jenkins名称信息失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinsnamelist"))
		c.JSON(200, Response{e.FAILED_SELECT_DB, r, e.GetMsg(e.FAILED_SELECT_DB)})
	} else {
		L.Lzap.Info("查询jenkins信息成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinsnamelist"))
		c.JSON(200, Response{e.SUCCESS, r, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 查询jenkins列表
// @Description 格式：GET /api/v1/service/getjenkinslist
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/getjenkinslist [get]
func ServiceGetJenkinsList(c *gin.Context) {
	var j modules.Jenkins
	r, err := j.ServiceJenkinsLists()
	if err != nil || r == nil {
		L.Lzap.Error("查询jenkins名称信息失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinslist"))
		c.JSON(200, Response{e.FAILED_SELECT_DB, r, e.GetMsg(e.FAILED_SELECT_DB)})
	} else {
		L.Lzap.Info("查询jenkins信息成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinslist"))
		c.JSON(200, Response{e.SUCCESS, r, e.GetMsg(e.SUCCESS)})
	}
}

//jenkins job

// @Summary 查询jenkins job列表
// @Description 格式：GET /api/v1/service/getjenkinsjobslist?jenkins=xxx
// @Description jenkins:jenkins名称
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param jenkins query string true "jenkins名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/getjenkinsjobslist [get]
func ServiceGetJenkinsJobsList(c *gin.Context) {
	j := c.Query("jenkins")
	J := modules.JenkinsJobs{Jenkins: j}
	r, err := J.ServiceJenkinsJobsList()
	if err != nil {
		L.Lzap.Error("查询jenkins job列表失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinsjobslist"))
		c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
	} else {
		L.Lzap.Info("查询jenkins job列表成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinsjobslist"))
		c.JSON(200, Response{e.SUCCESS, r, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 查询jenkins job详情
// @Description 格式：GET /api/v1/service/getjenkinsjobsinfo?jenkins=xxx&jobname=xxx
// @Description jenkins:jenkins名称
// @Description jobname:job名称
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param jenkins query string true "jenkins名称"
// @Param jobname query string true "job名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/getjenkinsjobsinfo [get]
func ServiceGetJenkinsJobsInfo(c *gin.Context) {
	j := c.Query("jenkins")
	b := c.Query("jobname")
	if j == "" || b == "" {
		L.Lzap.Error("查询jenkins job信息失败", zap.String("info", "有必须参数为空"), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinsjobsinfo"))
		c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
	} else {
		J := modules.JenkinsJobs{Jenkins: j, Jobname: b}
		if err := J.ServiceJenkinsGetJobsInfo(); err != nil {
			L.Lzap.Error("查询jenkins job信息失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinsjobsinfo"))
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			L.Lzap.Info("查询jenkins job信息成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/getjenkinsjobsinfo"))
			c.JSON(200, Response{e.SUCCESS, J, e.GetMsg(e.SUCCESS)})
		}
	}
}

// @Summary 修改jenkins job详情(更新jenkins pipeline未完待续)
// @Description 格式：GET /api/v1/service/updatejenkinsjobs
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.JenkinsJobs body modules.JenkinsJobs true "jenkins job信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/updatejenkinsjobs [get]
func ServiceUpdateJenkinsJobs(c *gin.Context) {
	var J modules.JenkinsJobs
	if err := c.BindJSON(&J); err != nil {
		L.Lzap.Error("修改jenkins job信息失败，数据绑定失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkinsjobs"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if J.Jenkins == "" && J.Jobname == "" {
			L.Lzap.Error("修改jenkins job信息失败，必须字段为空", zap.Reflect("data", J), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkinsjobs"))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		} else {
			J.Id, J.State, J.Ctime, J.Jfolder = 0, 0, "", ""
			if err := J.ServiceJenkinsUpdateJob(); err != nil {
				L.Lzap.Error("修改jenkins job信息失败，数据库写入失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkinsjobs"))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				//jenkins job发生变更，需要重新生成jenkins pipeline
				if err := JenkinsPipelineInit(J); err != nil {
					L.Lzap.Error("修改jenkins job信息数据库修改成功，jenkins pipeline修改失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkinsjobs"))
					_ = J.ServiceJenkinsControl(2)
					c.JSON(200, Response{e.SUCCESS, nil, "jenkins pipeline修改失败，请联系管理员"})
				} else {
					L.Lzap.Info("修改jenkins job信息成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/updatejenkinsjobs"))
					c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
				}
			}
		}
	}
}

//更新jenkins pipeline（大任务，要仔细）
func JenkinsPipelineInit(j modules.JenkinsJobs) error {
	return nil
}

type JenkinsJobDo struct {
	Jenkins string `json:"jenkins"`
	Type    string `json:"type"`
	Jobname string `json:"jobname"`
}

// @Summary 启用、禁用jenkins job
// @Description 格式：GET /api/v1/service/controljenkinsjobs
// @Description jenkins：jenkins名称
// @Description jobname：job名称(应用名称)
// @Description type: 1：启用，2：禁用
// @Description 启用或禁用jenkins：jenkins、jobname、type不能为空
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param JenkinsJobDo body JenkinsJobDo true "jenkins job操作信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/controljenkinsjobs [get]
func ServiceJenkinsJobControl(c *gin.Context) {
	var Do JenkinsJobDo
	if err := c.BindJSON(&Do); err != nil {
		L.Lzap.Error("启用或禁用jenkins job失败，数据绑定失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/controljenkinsjobs"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if Do.Jobname != "" && Do.Jenkins != "" && (Do.Type == "1" || Do.Type == "2") {
			d, _ := strconv.Atoi(Do.Type)
			J := modules.JenkinsJobs{Jenkins: Do.Jenkins, Jobname: Do.Jobname, State: d}
			if err := J.ServiceJenkinsControl(J.State); err != nil {
				L.Lzap.Error("启用或禁用jenkins job失败，数据库更新失败", zap.Error(err), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/controljenkinsjobs"))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				L.Lzap.Info("启用或禁用jenkins job成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/controljenkinsjobs"))
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			}
		} else {
			L.Lzap.Error("启用或禁用jenkins job失败,有必须字段为空", zap.Reflect("data", Do), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/service/controljenkinsjobs"))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		}
	}
}

type JenkinsJobMove struct {
	Tojenkins string `json:"tojenkins"`
	JenkinsJobDo
}

// @Summary 迁移jenkins job(未完待续)
// @Description 格式：GET /api/v1/service/movejenkinsjobs
// @Description jenkins：源jenkins名称
// @Description tojenkins: 目标jenkins名称
// @Description jobname：job名称(应用名称)
// @Description 迁移jenkins：jenkins、jobname、tojenkins均不能为空
// @Tags service
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param JenkinsJobMove body JenkinsJobMove true "jenkins job信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/service/movejenkinsjobs [get]
func ServiceJenkinsJobMove(c *gin.Context) {
	//迁移流程：
	//	1、生成迁移记录
	//	2、
	//
	//
}

//jenkins构建信息
//查询项目应用构建列表
//查询
