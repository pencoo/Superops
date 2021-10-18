package v1

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"superops/libs/config"
	"superops/libs/e"
	L "superops/middlewares/ginzap"
	"superops/modules"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Summary 新增命令
// @Description 格式：GET /api/v1/commands/createcommand
// @Tags Commands
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Kubecmds body modules.Kubecmds true "命令信息"
// @Success 200 {object} Response "ok"
// @Router /api/v1/commands/createcommand [post]
func ComsCreate(c *gin.Context) {
	var cm modules.Kubecmds
	if err := c.ShouldBindJSON(&cm); err != nil {
		//参数错误
		L.Lzap.Error("数据解析失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createcommand"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if cm.Comname == "" || cm.Comvers == "" {
			//数据格式不对
			L.Lzap.Error("有必须字段为空", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createcommand"), zap.Reflect("data", cm))
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		} else {
			cm.State = 0
			cm.Compath = config.FilePath + "/bin/" + cm.Comname
			if err := cm.ComsCreate(); err != nil {
				//数据写入失败
				L.Lzap.Error("数据库写入失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createcommand"), zap.Error(err))
				c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				//命令新增成功
				L.Lzap.Info("命令信息新增成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createcommand"))
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			}
		}
	}
}

// @Summary 上传命令执行文件
// @Description 格式：GET /api/v1/commands/createfile?comname=xxx file=文件
// @Description 注意：需要先成功创建命令后，才能上传命令执行文件(命令必须是激活状态才能上传可执行文件)
// @Description comname是命令名称，获取接口：/api/v1/commands/getcommandsnamelist
// @Tags Commands
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param comname query string true "命令名称"
// @Param file formData file true "命令文件"
// @Success 200 {object} Response "ok"
// @Router /api/v1/commands/createfile [post]
func ComsUploadFile(c *gin.Context) {
	var cm modules.Kubecmds
	cc := c.Query("comname")
	cm.Comname = cc
	if cm.Comname == "" || cm.ComsSelectByName() != nil {
		//命令为空或命令不存在
		L.Lzap.Error("命令名称为空或命令不存在", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createfile"))
		c.JSON(200, Response{e.ERROR, nil, "命令名称为空或者命令不存在"})
	} else {
		f, err := c.FormFile("file")
		if err != nil {
			//请求中获取文件失败
			L.Lzap.Error("请求中获取文件失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createfile"), zap.Error(err))
			c.JSON(200, Response{e.ERROR, "", fmt.Sprint(err)})
		} else {
			NotExistCreatePath(cm.Compath)
			if c.SaveUploadedFile(f, cm.Compath) != nil {
				//文件写入磁盘失败
				L.Lzap.Error("命令文件写入磁盘失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createfile"))
				c.JSON(200, Response{e.ERROR, nil, "文件写入磁盘失败"})
			} else {
				//文件写入成功
				os.Chmod(cm.Compath, 0777)
				L.Lzap.Info("命令文件写入成功", zap.String("command", cm.Compath), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/createfile"))
				c.JSON(200, Response{e.SUCCESS, "", e.GetMsg(e.SUCCESS)})
			}
		}
	}
}

// @Summary 获取命令列表(详情)
// @Description 格式：GET /api/v1/commands/getcommandslist?type=0|1|2|3
// @Description type: 0、查询无命令文件的命令，1、获取激活中的命令，2、获取禁用的命令，3、获取所有命令
// @Tags Commands
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param type query string true "查询类型"
// @Success 200 {object} Response "ok"
// @Router /api/v1/commands/getcommandslist [get]
func ComsGetCommandsList(c *gin.Context) {
	t := c.Query("type")
	if t == "0" || t == "1" || t == "2" || t == "3" {
		ct, _ := strconv.Atoi(t)
		cl, err := modules.ComsSelectList(ct)
		if err != nil {
			L.Lzap.Error("命令列表获取失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/getcommandslist"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			L.Lzap.Info("命令列表获取成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/getcommandslist"))
			c.JSON(200, Response{e.SUCCESS, cl, e.GetMsg(e.SUCCESS)})
		}
	} else {
		L.Lzap.Error("查询类型解析失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/getcommandslist"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	}
}

// @Summary 获取激活命令名称列表
// @Description 格式：GET /api/v1/commands/getcommandsnamelist
// @Tags Commands
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/commands/getcommandsnamelist [get]
func ComsGetCommandsNameList(c *gin.Context) {
	li, err := modules.ComsSelectNameList()
	if err != nil {
		//数据查询失败
		L.Lzap.Error("数据查询失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/getcommandsnamelist"), zap.Error(err))
		c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
	} else {
		L.Lzap.Info("命令名称列表获取成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/getcommandsnamelist"))
		c.JSON(200, Response{e.SUCCESS, li, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 命令启用或禁用
// @Description 格式：GET /api/v1/commands/controlcommand?command=xxx&type=0|1|2|3
// @Description command: 命令名称
// @Description type: 0、查询无命令文件的命令，1、获取激活中的命令，2、获取禁用的命令，3、获取所有命令
// @Tags Commands
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param type query string true "查询类型。0、查询无命令文件的命令，1、获取激活中的命令，2、获取禁用的命令，3、获取所有命令"
// @Param command query string true "需要禁用的命令名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/commands/controlcommand [get]
func ComsControlCommand(c *gin.Context) {
	cm := c.Query("command")
	t := c.Query("type")
	if cm != "" && (t == "0" || t == "1" || t == "2" || t == "3") {
		ct, _ := strconv.Atoi(t)
		cms := modules.Kubecmds{Comname: cm}
		if err := cms.ComsEnableOrDisable(ct); err != nil {
			//命令操作成功
			L.Lzap.Info("命令启用/禁用成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/controlcommand"))
			c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
		} else {
			//命令操作失败
			L.Lzap.Error("命令启用/禁用失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/controlcommand"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_INSERT_DB, nil, e.GetMsg(e.FAILED_INSERT_DB)})
		}
	} else {
		//参数错误
		L.Lzap.Error("请求参数解析错误", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/controlcommand"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	}
}

// @Summary 删除命令
// @Description 格式：GET /api/v1/commands/delcommands?command=xxx
// @Description command: 命令名称
// @Description command获取接口：/api/v1/commands/getcommandsnamelist
// @Tags Commands
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param command query string true "删除命令名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/commands/delcommands [get]
func ComsDelCommands(c *gin.Context) {
	cm := c.Query("command")
	if cm == "" {
		c.JSON(200, Response{e.ERROR, nil, "命令为空"})
	} else {
		cms := modules.Kubecmds{Comname: cm}
		cms.ComsSelectByName()
		if err := cms.ComsDelete(); err != nil {
			L.Lzap.Error("数据库删除命令失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/delcommands"), zap.Error(err))
			c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
		} else {
			cmd := exec.Command("sh", "-c", "rm -rf "+cms.Compath)
			info, err := cmd.Output()
			if err != nil {
				L.Lzap.Error("数据删除成功，命令文件删除失败", zap.String("command_strout", string(info)), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/delcommands"), zap.Error(err), zap.String("file", cms.Compath))
				c.JSON(200, Response{e.ERROR, err, "数据删除成功，命令文件删除失败：" + cms.Compath})
			} else {
				L.Lzap.Info("命令删除成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/commands/delcommands"))
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			}
		}
	}
}
