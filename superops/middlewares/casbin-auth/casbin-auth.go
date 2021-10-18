package casbin_auth

import (
	"superops/libs/config"
	"superops/modules"

	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/gin-gonic/gin"
)

var Enforcer *casbin.Enforcer

func Csbinit() {
	admin := config.AdminName
	if admin == "" {
		admin = "pencoo"
	}
	minfo := "(g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act) || r.sub == " + admin
	m := casbin.NewModel()
	m.AddDef("r", "r", "sub, obj, act")
	m.AddDef("p", "p", "sub, obj, act")
	m.AddDef("g", "g", "_, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")
	m.AddDef("m", "m", minfo)
	a := gormadapter.NewAdapterByDB(modules.Db)
	Enforcer = casbin.NewEnforcer(m, a)
}
func CasbinAuth(c *gin.Context) {
	var userName string
	userName = c.GetString("Name")
	if userName == "" {
		c.JSON(200, gin.H{
			"code":    401,
			"message": "未认证",
			"data":    "",
		})
		c.Abort()
		return
	}
	p := c.Request.URL.Path
	m := c.Request.Method
	res := Enforcer.Enforce(userName, p, m)
	if !res {
		c.JSON(200, gin.H{
			"code":    401,
			"message": "Unauthorized",
			"data":    "",
		})
		c.Abort()
		return
	}
	c.Next()
}
