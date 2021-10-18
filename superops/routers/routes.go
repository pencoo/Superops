package routers

import (
	"net/http"
	v1 "superops/apps/v1"
	_ "superops/docs"
	"superops/libs/config"
	jwt_auth "superops/middlewares/jwt-auth"

	"github.com/gin-gonic/gin"
	sf "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
)

func RouterInit(r *gin.Engine) {
	r.Use(Cors())
	if config.Apidoc == "true" {
		r.GET("/api/v1/swagger/*any", gs.WrapHandler(sf.Handler))
	}
	r.POST("/api/v1/login", v1.UserLogin)        //用户登录
	r.POST("/api/v1/register", v1.UsersRegister) //用户注册
	u := r.Group("/api/v1/users", jwt_auth.JWTAuth())
	{
		//用户接口
		u.POST("/chageinfo", v1.UsersUpdateInfo)           //修改用户信息
		u.GET("/logout", v1.UsersLogout)                   //用户注销登录
		u.POST("/chagespass", v1.UsersChangePassword)      //修改用户密码
		u.GET("/userinfo", v1.UsersGetUserInfo)            //获取用户信息
		u.GET("/getlists", v1.GetUsersLists)               //获取用户列表
		u.GET("/createapitoken", v1.UsersCreateApiToken)   //用户创建API Token
		u.GET("/deleteapitoken", v1.UsersDeleteApiToken)   //用户删除API Token
		u.GET("/enableapitoken", v1.UsersEnableApiToken)   //用户启用API Token
		u.GET("/disableapitoken", v1.UsersDisableApiToken) //用户禁用API Token
		u.GET("/getapitoken", v1.UsersGetApiToken)         //获取用户API Token
		//管理员接口
		u.POST("/admingetlist", v1.UsersGetUserlistByAdmin) //查询用户信息列表
		u.GET("/geturllist", v1.GetUrlListInfo)             //获取url列表
		u.GET("/geturlmethodlist", v1.GetUrlMethodListInfo) //获取url-method列表
		u.GET("/controls", v1.UsersOpenOrOff)               //单用户启用或禁用
		u.POST("/controls", v1.UsersGroupsOpenOrOff)        //多用户启用或禁用
		u.POST("/rolelist", v1.UsersGetRoleList)            //获取角色列表
		u.POST("/roleadd", v1.UsersRoleAdd)                 //添加角色或给角色添加权限
		u.POST("/roleDelPower", v1.UsersRoleDelPower)       //给角色删除权限
		u.GET("/roledel", v1.UsersRoleDel)                  //删除角色
		u.GET("/rolePower", v1.UsersRolePower)              //获取角色权限列表
		u.GET("/roleadduser", v1.UsersAddToRole)            //将用户加入一个角色
		u.GET("/roledeluser", v1.UsersDelToRole)            //将用户从一个角色中删除
		u.GET("/rolelistbyuser", v1.UsersGetRole)           //获取用户角色列表
		u.GET("/powerbyuser", v1.UsersGetUserPower)         //获取用户权限列表
	}
	cmd := r.Group("/api/v1/commands", jwt_auth.JWTAuth()) //管理员接口
	{
		cmd.POST("/createcommand", v1.ComsCreate)                   //新增命令
		cmd.POST("/createfile", v1.ComsUploadFile)                  //上传命令执行文件
		cmd.GET("/getcommandslist", v1.ComsGetCommandsList)         //获取命令信息列表
		cmd.GET("/getcommandsnamelist", v1.ComsGetCommandsNameList) //获取命令名称列表
		cmd.GET("/controlcommand", v1.ComsControlCommand)           //启用或禁用命令
		cmd.GET("/delcommands", v1.ComsDelCommands)                 //删除命令
	}
	cmdb := r.Group("/api/v1/cmdb", jwt_auth.JWTAuth()) //kubeconfig
	{
		cmdb.POST("/createkubeconfig", v1.CmdbCreateKubeconfig)   //添加kubeconfig文件
		cmdb.POST("/getkubeconfiglist", v1.CmdbGetKubeConfigList) //获取kubeconfig列表
		cmdb.GET("/getkubeconfig", v1.CmdbGetKubeConfig)          //获取kubeconfig文件详情(单条详情)
		cmdb.GET("/deletekubeconfig", v1.CmdbDeleteKubeConfig)    //删除kubeconfig文件
		cmdb.POST("/updatekubeconfig", v1.CmdbUpdateKubeConfig)   //修改kubeconfig文件
		cmdb.GET("/controlkubeconfig", v1.CmdbControlKubeConfig)  //启用或禁用kubeconfig文件(未完待续)
	}
	service := r.Group("/api/v1/service", jwt_auth.JWTAuth()) //service
	{
		service.POST("/createharbor", v1.ServiceCreateHarbor)            //添加harbor
		service.GET("/deleteharbor", v1.ServiceDeleteHarbor)             //删除harbor
		service.POST("/Updateharbor", v1.ServiceUpdateHarbor)            //修改harbor
		service.GET("/getharbornamelist", v1.ServiceHarborGetNameList)   //获取harbor名称列表
		service.GET("/getharborinfo", v1.ServiceHarborGetInfo)           //获取harbor信息
		service.GET("/getharborlist", v1.ServiceHarborGetList)           //获取harbor列表
		service.POST("/createjenkins", v1.ServiceCreateJenkins)          //添加jenkins
		service.POST("/updatejenkins", v1.ServiceUpdateJenkins)          //修改jenkins
		service.POST("/controljenkins", v1.ServiceJenkinsControls)       //修改jenkins状态
		service.GET("/jenkinsinfo", v1.ServiceGetJenkinsInfo)            //获取jenkins信息
		service.GET("/getjenkinsnamelist", v1.ServiceGetJenkinsNameList) //获取jenkins名称列表
		service.GET("/getjenkinslist", v1.ServiceGetJenkinsList)         //获取jenkins列表
		service.GET("/getjenkinsjobslist", v1.ServiceGetJenkinsJobsList) //获取jenkins job列表
	}
	project := r.Group("/api/v1/project", jwt_auth.JWTAuth())
	{
		project.POST("/createk8sproject", v1.ProjectManageCreateK8sProject)                    //添加k8s项目
		project.POST("/updatek8sproject", v1.ProjectManageUpdateK8sProject)                    //更新k8s项目
		project.GET("/controlk8sproject", v1.ProjectManageControlK8sProject)                   //启用、禁用k8s项目
		project.GET("/getk8sprojectlist", v1.ProjectManageGetK8sProjectList)                   //查询k8s项目列表
		project.GET("/getk8sprojectnamelist", v1.ProjectManageGetK8sProjectNameList)           //管理员查询k8s项目名称列表
		project.GET("/usergetk8sprojectnamelist", v1.ProjectManageGetK8sProjectNameListByUser) //用户查询k8s项目名称列表
		project.GET("/selectk8sprojectinfo", v1.ProjectManageGetK8sProjectInfo)                //查询k8s项目信息
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Access-Token,Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
