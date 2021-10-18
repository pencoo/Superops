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

// @Summary 缓存重建
// @Description 格式：GET /api/v1/cache/rebuild
// @Tags cache
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param user query string true "用户"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cache/rebuild [get]
func RedisCacheRebuild() {
	//
}

// @Summary 用户项目(所有激活项目)列表缓存重建
// @Description 格式：GET /api/v1/cache/rebuildprojectuserlists
// @Tags cache
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cache/rebuildprojectuserlists [get]
func RedisCacheRebuildprojectuserlists(c *gin.Context) {
	var p modules.Projects
	pls, err := p.ProjectManageProjectGetNameList(1)
	if err != nil || pls == nil {
		L.Lzap.Info("项目缓存刷新失败，数据库查询失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cache/rebuildprojectuserlists"))
		c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
	} else {
		var R SFLS
		for _, v := range pls {
			err := Rebuildprojectuserlist(v)
			if err != nil {
				R.Failed = append(R.Failed, v)
			} else {
				R.Success = append(R.Success, v)
			}
		}
		L.Lzap.Info("项目缓存刷新成功", zap.Reflect("project", R), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cache/rebuildprojectuserlists"))
		c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 用户项目(指定)列表缓存重建
// @Description 格式：GET /api/v1/cache/rebuildprojectuserlist?project=xxx
// @Description project：指定重建缓存的项目名称
// @Tags cache
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param project query string true "项目名称"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cache/rebuildprojectuserlist [get]
func RedisCacheRebuildprojectuserlist(c *gin.Context) {
	p := c.Query("project")
	if p == "" {
		L.Lzap.Error("项目缓存重建失败", zap.String("error", "项目名称为空"), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cache/rebuildprojectuserlist"))
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if err := Rebuildprojectuserlist(p); err != nil {
			L.Lzap.Error("项目缓存重建失败", zap.Error(err), zap.String("project", p), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cache/rebuildprojectuserlist"))
			c.JSON(200, Response{e.ERROR, err, e.GetMsg(e.ERROR)})
		} else {
			L.Lzap.Info("项目缓存重建成功", zap.String("project", p), zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cache/rebuildprojectuserlist"))
			c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
		}
	}
}

func Rebuildprojectuserlist(p string) error {
	P := modules.ProjectUsers{Project: p}
	pp := modules.Projects{Project: p}
	if pp.ProjectManageProjectGetInfo() != nil || pp.State != 1 {
		return errors.New("项目不存在")
	} else {
		pual, err := P.ProjectManageProjectUsersGetList()
		if err != nil || !(pual.Ops != nil || pual.Dev != nil || pual.Test != nil) {
			return err
		} else {
			if pual.Dev != nil {
				_ = rcache.RedisKeyDel(rcache.Project_user_list_dev + P.Project)
				for _, v := range pual.Dev {
					_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_dev+P.Project, v)
				}
			}
			if pual.Test != nil {
				_ = rcache.RedisKeyDel(rcache.Project_user_list_test + P.Project)
				for _, v := range pual.Test {
					_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_test+P.Project, v)
				}
			}
			if pual.Ops != nil {
				_ = rcache.RedisKeyDel(rcache.Project_user_list_ops + P.Project)
				for _, v := range pual.Ops {
					_ = rcache.RedisCacheSetAdd(rcache.Project_user_list_ops+P.Project, v)
				}
			}
		}
		return nil
	}
}

// @Summary 用户列表缓存重建
// @Description 格式：GET /api/v1/cache/rebuilduserlist
// @Tags cache
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cache/rebuilduserlist [get]
func RedisCacheRebuildenableUserList(c *gin.Context) {
	var u modules.Users
	enableuserlist := u.GetUsersLists(1)
	disableuserlist := u.GetUsersLists(2)
	if enableuserlist != nil {
		_ = rcache.RedisKeyDel(rcache.User_enable_list)
		for _, v := range enableuserlist {
			_ = rcache.RedisCacheSetAdd(rcache.User_enable_list, v)
		}
	}
	if disableuserlist != nil {
		_ = rcache.RedisKeyDel(rcache.User_disable_list)
		for _, v := range disableuserlist {
			_ = rcache.RedisCacheSetAdd(rcache.User_disable_list, v)
		}
	}
}

// @Summary 刷新所有用户API token
// @Description 格式：GET /api/v1/cache/rebuilduserapitoken
// @Tags cache
// @Accept json
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response "ok"
// @Router /api/v1/cache/rebuilduserapitoken [get]
func RedisCacheRebuildUserApiToken(c *gin.Context) {
	apitoken := modules.UsersApiTokenLists()
	if apitoken != nil {
		for _, v := range apitoken {
			t, _ := rcache.RedisCacheStringGet(rcache.Userinfo_api_token + v.Username)
			if t == "" {
				_ = rcache.RedisCacheStringAdd(rcache.Userinfo_api_token+v.Username, v.Token)
			}
		}
		L.Lzap.Info("刷新用户api token成功", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cache/rebuilduserapitoken"))
		c.JSON(200, Response{e.SUCCESS, apitoken, e.GetMsg(e.SUCCESS)})
	} else {
		L.Lzap.Error("刷新用户api token数据查询失败", zap.String("user", c.GetString("Name")), zap.String("path", "/api/v1/cache/rebuilduserapitoken"))
		c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
	}
}
