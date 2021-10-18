package modules

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//项目管理表
type Projects struct {
	Id         int    `gorm:"primory_key;auto_increment" json:"id"`
	Project    string `gorm:"type:varchar(50);uniqueIndex:proj;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
	Cuser      string `gorm:"type:varchar(50);comment:'项目创建用户'" binding:"omitempty,min=3,max=50" json:"cuser"`
	Ctime      string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	Jenkins    string `gorm:"type:varchar(50);comment:'项目部署使用的jenkins'" binding:"omitempty,min=6,max=50" json:"jenkins"`
	Repository string `gorm:"type:varchar(50);comment:'保存制品的仓库'" binding:"omitempty,min=6,max=50" json:"repository"`
	Context    string `gorm:"type:varchar(500);comment:'项目描述'" binding:"omitempty,min=3,max=500" json:"context"`
	State      int    `gorm:"type:tinyint;default:1;comment:'项目状态，1激活，2禁用'" binding:"gte=0,lte=10" json:"State"`
}

//项目用户管理表
type ProjectUsers struct {
	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
	Project  string `gorm:"type:varchar(50);uniqueIndex:pus;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
	Username string `gorm:"type:varchar(50);uniqueIndex:pus;comment:'用户名称'" binding:"omitempty,min=3,max=50" json:"username"`
	Urole    int    `gorm:"type:tinyint;uniqueIndex:pus;comment:'用户角色'" binding:"gte=0,lte=10" json:"urole"`
	Ctime    string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
}

//项目应用管理
type ProjectApps struct {
	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
	Project  string `gorm:"type:varchar(50);uniqueIndex:pas;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
	Appname  string `gorm:"type:varchar(50);uniqueIndex:pas;comment:'组件名称'" binding:"omitempty,min=3,max=50" json:"appname"`
	Username string `gorm:"type:varchar(20);comment:'创建组件的用户'" binding:"omitempty,min=3,max=50" json:"username"`
	Langtype string `gorm:"type:varchar(50);comment:'开发语言，用于对应jenkins模板'" binding:"omitempty,min=2,max=50" json:"langtype"`
	Dtmp     string `gorm:"type:varchar(50);comment:'部署模板，用于对应部署工具的模板'" binding:"omitempty,min=2,max=50" json:"dtmp"`
	Giturl   string `gorm:"type:varchar(50);comment:'组件git地址'" binding:"omitempty,url,min=6,max=50" json:"giturl"`
	Gituser  string `gorm:"type:varchar(50);comment:'拉取组件的git用户名'" binding:"omitempty,min=6,max=50" json:"gituser"`
	Gitpass  string `gorm:"type:varchar(50);comment:'拉取组件的git密码'" binding:"omitempty,min=6,max=50" json:"gitpass"`
	Gittoken string `gorm:"type:varchar(50);comment:'用户的gittoken'" binding:"omitempty,min=6,max=50" json:"gittoken"`
	Ctime    string `gorm:"type:varchar(20);comment:'创建时间';comment:'创建时间'" json:"ctime"`
	Route    string `gorm:"type:varchar(20);comment:'路由'" binding:"omitempty,min=1,max=20" json:"route"`
	Context  string `gorm:"type:varchar(500);comment:'组件描述信息'" binding:"omitempty,min=6,max=500" json:"context"`
	State    int    `gorm:"type:tinyint;default:1;comment:'组件状态，1激活，2禁用'" binding:"gte=0,lte=10" json:"state"`
}

//项目环境管理
type ProjectEnvs struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Project string `gorm:"type:varchar(50);uniqueIndex:pes;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
	Domain  string `gorm:"type:varchar(20);uniqueIndex:pes;comment:'项目部署域名'" binding:"omitempty,fqdn,min=3,max=20" json:"domain"`
	Penv    string `gorm:"type:varchar(20);uniqueIndex:pes;comment:'环境名称'" binding:"omitempty,min=3,max=20" json:"penv"`
	Kname   string `gorm:"type:varchar(50);comment:'部署环境的k8s名称'" binding:"omitempty,min=3,max=50" json:"kname"`
	Envrole int    `gorm:"type:tinyint;comment:'环境角色'" binding:"gte=0,lte=10" json:"envrole"`
	Ctime   string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	Cuser   string `gorm:"type:varchar(50);comment:'创建用户'" json:"cuser"`
	Context string `gorm:"type:varchar(500);comment:'环境描述'" json:"context"`
	State   int    `gorm:"type:tinyint;default:1;comment:'环境状态，1激活，2禁用'" binding:"gte=0,lte=10" json:"state"`
}

//项目配置管理
type ProjectConfs struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Project string `gorm:"type:varchar(50);uniqueIndex:pcs;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
	Appname string `gorm:"type:varchar(50);uniqueIndex:pcs;comment:'组件名称'" binding:"omitempty,min=3,max=50" json:"appname"`
	Confkey string `gorm:"type:varchar(50);uniqueIndex:pcs;comment:'配置键'" binding:"omitempty,min=1,max=50" json:"confkey"`
	Context string `gorm:"type:varchar(200);comment:'配置键描述'" binding:"omitempty,min=1,max=200" json:"context"`
	Ctime   string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	Cuser   string `gorm:"type:varchar(50);comment:'创建用户'" json:"cuser"`
	State   int    `gorm:"type:tinyint;default:1;comment:'配置状态，1激活，2禁用'" binding:"gte=0,lte=3" json:"state"`
}

//项目配置键值
type AppEnvConfig struct {
	Id        int    `gorm:"primory_key;auto_increment" json:"id"`
	Pcid      int    `gorm:"type:int;uniqueIndex:aec;comment:'项目配置(ProjectConfs)ID'" json:"pcid"`
	Peid      int    `gorm:"type:int;uniqueIndex:aec;comment:'项目环境ID'" json:"peid"`
	Cuser     string `gorm:"type:varchar(50);comment:'创建用户'" json:"cuser"`
	Confvalue string `gorm:"type:varchar(500);comment:'配置值'" binding:"omitempty,min=1,max=500" json:"confvalue"`
}

//部署管理
type ProjectDeploy struct {
	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
	Project  string `gorm:"type:varchar(50);comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
	Appname  string `gorm:"type:varchar(50);comment:'组件名称'" binding:"omitempty,min=3,max=50" json:"appname"`
	Duser    string `gorm:"type:varchar(50);comment:'部署用户'" binding:"omitempty,min=3,max=50" json:"duser"`
	Package  string `gorm:"type:varchar(50);comment:'部署包'" binding:"omitempty,min=3,max=50" json:"package"`
	Denv     string `gorm:"type:varchar(20);comment:'部署环境'" binding:"omitempty,min=3,max=20" json:"denv"`
	State    int    `gorm:"type:tinyint;default:0;comment:'部署状态'" binding:"gte=0,lte=10" json:"state"`
	Manifest string `gorm:"type:text;comment:'部署细节'" json:"manifest"`
	Ctime    string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
}

func (j *Projects) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *ProjectUsers) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *ProjectApps) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *ProjectEnvs) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *ProjectConfs) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *AppEnvConfig) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *ProjectDeploy) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}

//项目管理表
//type Projects struct {
//	Id         int    `gorm:"primory_key;auto_increment" json:"id"`
//	Project    string `gorm:"type:varchar(50);uniqueIndex:proj;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
//	Cuser      string `gorm:"type:varchar(50);comment:'项目创建用户'" binding:"omitempty,min=3,max=50" json:"cuser"`
//	Ctime      string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
//	Jenkins    string `gorm:"type:varchar(50);comment:'项目部署使用的jenkins'" binding:"omitempty,min=6,max=50" json:"jenkins"`
//	Repository string `gorm:"type:varchar(50);comment:'保存制品的仓库'" binding:"omitempty,min=6,max=50" json:"repository"`
//	Context    string `gorm:"type:varchar(500);comment:'项目描述'" binding:"omitempty,min=3,max=500" json:"context"`
//	State      int    `gorm:"type:tinyint;default:1;comment:'项目状态，1激活，2禁用'" binding:"gte=0,lte=10" json:"State"`
//}

//新增项目
func (p *Projects) ProjectManageProjectCreate() error {
	return Db.Model(&p).Create(&p).Error
}

//修改项目
func (p *Projects) ProjectManageProjectUpdate() error {
	return Db.Model(&p).Where("`project` = ?", p.Project).Update(&p).Error
}

//激活、禁用项目,n=1:激活，n=2:禁用
func (p *Projects) ProjectManageProjectControl(n int) error {
	return Db.Model(&p).Where("`project` = ?", p.Project).Update("state", n).Error
}

//查询项目列表
func (p *Projects) ProjectManageProjectGetList() ([]Projects, error) {
	var H []Projects
	err := Db.Model(&p).Scan(&H).Error
	return H, err
}

//查询单条项目
func (p *Projects) ProjectManageProjectGetInfo() error {
	return Db.Model(&p).Where("`project` = ?", p.Project).Find(&p).Error
}

//查询项目名称列表
func (p *Projects) ProjectManageProjectGetNameList(i int) ([]string, error) {
	var H []Projects
	var R []string
	var err error
	err = Db.Model(&p).Where("state = ?", i).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Project)
		}
	}
	return R, err
}

//根据harbor获取项目列表
func (p *Projects) ProjectManageProjectGetListByHarbor() ([]Projects, error) {
	var H []Projects
	err := Db.Model(&p).Where("repository", p.Repository).Scan(&H).Error
	return H, err
}

//根据Jenkins获取项目列表
func (p *Projects) ProjectManageProjectGetListByJenkins() ([]Projects, error) {
	var H []Projects
	err := Db.Model(&p).Where("jenkins", p.Jenkins).Scan(&H).Error
	return H, err
}

//项目用户管理表
//type ProjectUsers struct {
//	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
//	Project  string `gorm:"type:varchar(50);uniqueIndex:pus;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
//	Username string `gorm:"type:varchar(50);uniqueIndex:pus;comment:'用户名称'" binding:"omitempty,min=3,max=50" json:"username"`
//	Urole    int    `gorm:"type:tinyint;uniqueIndex:pus;comment:'用户角色'" binding:"gte=0,lte=10" json:"urole"`
//	Ctime    string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
//}

//新增项目用户
func (p *ProjectUsers) ProjectManageProjectUsersCreate() error {
	return Db.Model(&p).Create(&p).Error
}

//删除项目用户
func (p *ProjectUsers) ProjectManageProjectUsersDelete() error {
	return Db.Model(&p).Where("`Urole` = ? and `Username` = ? and `project` = ?", p.Urole, p.Username, p.Project).Delete(&p).Error
}

//用户查询项目列表
func (p *ProjectUsers) ProjectManageProjectUsersGetProjectList() ([]ProjectUsers, error) {
	var H []ProjectUsers
	err := Db.Model(&p).Where("`username` = ?", p.Username).Scan(&H).Error
	return H, err
}

//查询项目开发、测试、运维用户列表，需要传入project和urole
func (p *ProjectUsers) ProjectManageProjectUsersGetRoleUserList() ([]string, error) {
	var H []ProjectUsers
	var R []string
	err := Db.Model(&p).Where("`Urole` = ? and `project` = ?", p.Urole, p.Project).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Username)
		}
	}
	return R, err
}

type Paul struct {
	Dev  []string
	Test []string
	Ops  []string
}

//查询项目所有用户
func (p *ProjectUsers) ProjectManageProjectUsersGetList() (Paul, error) {
	var H Paul
	var u ProjectUsers
	var err1, err2, err3 error
	u.Project, u.Urole = p.Project, 1
	H.Dev, err1 = u.ProjectManageProjectUsersGetRoleUserList()
	u.Urole = 2
	H.Test, err2 = u.ProjectManageProjectUsersGetRoleUserList()
	u.Urole = 3
	H.Ops, err3 = u.ProjectManageProjectUsersGetRoleUserList()
	if err1 == nil && err2 == nil && err3 == nil {
		return H, nil
	} else {
		return H, errors.New("查询出错")
	}
}

//查询项目角色
func (p *ProjectUsers) ProjectManageProjectUsersGetInfo() (Paul, error) {
	var H []ProjectUsers
	var R Paul
	err := Db.Model(&p).Where("`username` = ?", p.Username).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			if v.Urole == 1 {
				R.Dev = append(R.Dev, v.Project)
			} else if v.Urole == 2 {
				R.Test = append(R.Test, v.Project)
			} else if v.Urole == 3 {
				R.Ops = append(R.Ops, v.Project)
			}
		}
		return R, err
	} else {
		return R, err
	}
}

//查询不在项目中的用户列表
func (p *ProjectUsers) ProjectManageGetOutOfProjectUserList() ([]string, error) {
	var H []ProjectUsers
	var R []string
	err := Db.Raw("select username from users where state = 1 and username not in (select username from project_users where project=?)", p.Project).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Username)
		}
	}
	return R, err
}

//项目应用管理
//type ProjectApps struct {
//	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
//	Project  string `gorm:"type:varchar(50);uniqueIndex:pas;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
//	Appname  string `gorm:"type:varchar(50);uniqueIndex:pas;comment:'组件名称'" binding:"omitempty,min=3,max=50" json:"appname"`
//	Username string `gorm:"type:varchar(20);comment:'创建组件的用户'" binding:"omitempty,min=3,max=50" json:"username"`
//	Langtype string `gorm:"type:varchar(50);comment:'开发语言，用于对应jenkins模板'" binding:"omitempty,min=2,max=50" json:"langtype"`
//	Dtmp     string `gorm:"type:varchar(50);comment:'部署模板，用于对应部署工具的模板'" binding:"omitempty,min=2,max=50" json:"dtmp"`
//	Giturl   string `gorm:"type:varchar(50);comment:'组件git地址'" binding:"omitempty,url,min=6,max=50" json:"giturl"`
//	Gituser  string `gorm:"type:varchar(50);comment:'拉取组件的git用户名'" binding:"omitempty,min=6,max=50" json:"gituser"`
//	Gitpass  string `gorm:"type:varchar(50);comment:'拉取组件的git密码'" binding:"omitempty,min=6,max=50" json:"gitpass"`
//	Gittoken string `gorm:"type:varchar(50);comment:'用户的gittoken'" binding:"omitempty,min=6,max=50" json:"gittoken"`
//	Ctime    string `gorm:"type:varchar(20);comment:'创建时间';comment:'创建时间'" json:"ctime"`
//	Route    string `gorm:"type:varchar(20);comment:'路由'" binding:"omitempty,min=1,max=20" json:"route"`
//	Context  string `gorm:"type:varchar(500);comment:'组件描述信息'" binding:"omitempty,min=6,max=500" json:"context"`
//	State    int    `gorm:"type:tinyint;default:1;comment:'组件状态，1激活，2禁用'" binding:"gte=0,lte=10" json:"state"`
//}

//新增项目应用
func (p *ProjectApps) ProjectManageProjectAppsCreate() error {
	return Db.Model(&p).Create(&p).Error
}

//修改项目应用
func (p *ProjectApps) ProjectManageProjectAppsUpdate() error {
	return Db.Model(&p).Where("`appname` = ? and ", p.Appname, p.Project).Update(&p).Error
}

//激活、禁用项目应用,n=1:激活，n=2:禁用
func (p *ProjectApps) ProjectManageProjectAppsControl(n int) error {
	return Db.Model(&p).Where("`appname` = ? and `project` = ?", p.Appname, p.Project).Update("state", n).Error
}

//获取项目应用列表
func (p *ProjectApps) ProjectManageProjectAppsGetList() ([]ProjectApps, error) {
	var H []ProjectApps
	err := Db.Model(&p).Where("`project` = ?", p.Project).Scan(&H).Error
	return H, err
}

//获取项目应用名称列表
func (p *ProjectApps) ProjectManageProjectAppsGetNameList() ([]string, error) {
	var H []ProjectApps
	var R []string
	err := Db.Model(&p).Where("`project` = ? and `state` = 1", p.Project).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Appname)
		}
	}
	return R, err
}

//项目环境管理
//type ProjectEnvs struct {
//	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
//	Project string `gorm:"type:varchar(50);uniqueIndex:pes;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
//	Domain  string `gorm:"type:varchar(20);uniqueIndex:pes;comment:'项目部署域名'" binding:"omitempty,fqdn,min=3,max=20" json:"domain"`
//	Penv    string `gorm:"type:varchar(20);uniqueIndex:pes;comment:'环境名称'" binding:"omitempty,min=3,max=20" json:"penv"`
//	Kname   string `gorm:"type:varchar(50);comment:'部署环境的k8s名称'" binding:"omitempty,min=3,max=50" json:"kname"`
//	Envrole int    `gorm:"type:tinyint;comment:'环境角色'" binding:"gte=0,lte=10" json:"envrole"`
//	Ctime   string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
//	Cuser   string `gorm:"type:varchar(50);comment:'创建用户'" json:"cuser"`
//	Context string `gorm:"type:varchar(500);comment:'环境描述'" json:"context"`
//	State   int    `gorm:"type:tinyint;default:1;comment:'环境状态，1激活，2禁用'" binding:"gte=0,lte=10" json:"state"`
//}

//新增项目环境
func (p *ProjectEnvs) ProjectManageProjectEnvsCreate() error {
	return Db.Model(&p).Create(&p).Error
}

//修改项目环境
func (p *ProjectEnvs) ProjectManageProjectEnvsUpdate() error {
	return Db.Model(&p).Where("`envrole` = ? and `penv` = ? and `project` = ?", p.Envrole, p.Penv, p.Project).Update(&p).Error
}

//激活、禁用项目环境，1：激活，2禁用
func (p *ProjectEnvs) ProjectManageProjectEnvsControl(n int) error {
	return Db.Model(&p).Where("`envrole` = ? and `penv` = ? and `project` = ?", p.Envrole, p.Penv, p.Project).Update("state", n).Error
}

//删除项目环境
func (p *ProjectEnvs) ProjectManageProjectEnvsDelete() error {
	return Db.Model(&p).Where("`envrole` = ? and `penv` = ? and `project` = ?", p.Envrole, p.Penv, p.Project).Delete(&p).Error
}

//用户查询有权限的项目环境列表
func (p *ProjectEnvs) ProjectManageProjectEnvsGetList() ([]ProjectEnvs, error) {
	var H []ProjectEnvs
	err := Db.Model(&p).Where("`envrole` = ? and `project` = ?", p.Envrole, p.Project).Scan(&H).Error
	return H, err
}

//用户查询有权限的项目环境名称列表
func (p *ProjectEnvs) ProjectManageProjectEnvsGetNameList() ([]string, error) {
	var H []ProjectEnvs
	var R []string
	err := Db.Model(&p).Where("`state` = ? and `envrole` = ? and `project` = ?", p.State, p.Envrole, p.Project).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Penv)
		}
	}
	return R, err
}

//项目配置管理
//type ProjectConfs struct {
//	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
//	Project string `gorm:"type:varchar(50);uniqueIndex:pcs;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
//	Appname string `gorm:"type:varchar(50);uniqueIndex:pcs;comment:'组件名称'" binding:"omitempty,min=3,max=50" json:"appname"`
//	Confkey string `gorm:"type:varchar(50);uniqueIndex:pcs;comment:'配置键'" binding:"omitempty,min=1,max=50" json:"confkey"`
//	Context string `gorm:"type:varchar(200);comment:'配置键描述'" binding:"omitempty,min=1,max=200" json:"context"`
//	Ctime   string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
//	Cuser   string `gorm:"type:varchar(50);comment:'创建用户'" json:"cuser"`
//	State   int    `gorm:"type:tinyint;default:1;comment:'配置状态，1激活，2禁用'" binding:"gte=0,lte=3" json:"state"`
//}
//新增项目配置项
func (p *ProjectConfs) ProjectManageProjectConfsCreate() error {
	return Db.Model(&p).Create(&p).Error
}

//删除项目配置项
func (p *ProjectConfs) ProjectManageProjectConfsDelete() error {
	return Db.Model(&p).Where("`project` = ? and `appname` = ? and `confkey` = ?", p.Project, p.Appname, p.Confkey).Delete(&p).Error
}

//查询项目应用配置列表
func (p *ProjectConfs) ProjectManageProjectConfsAppServerGetList() ([]ProjectConfs, error) {
	var H []ProjectConfs
	err := Db.Model(&p).Where("`project` = ? and `appname` = ?", p.Project, p.Appname).Scan(&H).Error
	return H, err
}

//项目应用配置拷贝,Destination:拷贝目标,User:拷贝用户，p是源，需要提供project和appname
func (p *ProjectConfs) ProjectManageProjectConfsGetList(Destination string, User string) error {
	L, err := p.ProjectManageProjectConfsAppServerGetList()
	if err == nil {
		for _, v := range L {
			v.Appname, v.Ctime, v.Cuser = Destination, time.Now().In(TL).Format("2006-01-02 15:04:05"), User
			_ = v.ProjectManageProjectConfsCreate()
		}
		return nil
	} else {
		return err
	}
}

//项目配置键值
//type AppEnvConfig struct {
//	Id        int    `gorm:"primory_key;auto_increment" json:"id"`
//	Pcid      int    `gorm:"type:int;uniqueIndex:aec;comment:'项目配置(ProjectConfs)ID'" json:"pcid"`
//	Peid      int    `gorm:"type:int;uniqueIndex:aec;comment:'项目环境ID'" json:"peid"`
//  Cuser     string `gorm:"type:varchar(50);comment:'创建用户'" json:"cuser"`
//	Confvalue string `gorm:"type:varchar(500);comment:'配置值'" binding:"omitempty,min=1,max=500" json:"confvalue"`
//}

//新增配置值
func (c *AppEnvConfig) ProjectManageProjectConfCreateValues() error {
	return Db.Model(&c).Create(&c).Error
}

//修改配置值
func (c *AppEnvConfig) ProjectManageProjectConfUpdateValues() error {
	return Db.Raw("update app_env_config set cuser = ?,confvalue = ? where pcid = ? and peid = ?", c.Cuser, c.Confvalue, c.Pcid, c.Peid).Error
}

//删除配置值
func (c *AppEnvConfig) ProjectManageProjectConfDeleteValues() error {
	return Db.Model(&c).Where("pcid = ? and peid = ?", c.Pcid, c.Peid).Delete(&c).Error
}

type AppConfValue struct {
	Project   string `json:"project"`
	Appname   string `json:"appname"`
	Confkey   string `json:"confkey"`
	Context   string `json:"context"`
	Confvalue string `json:"confvalue"`
}

//查询项目应用配置列表,需要c具备：项目名称、应用名称。n：环境ID
func (c *ProjectConfs) ProjectManageProjectConfValueLists(n int) ([]AppConfValue, error) {
	var H []AppConfValue
	err := Db.Raw("select c.project as project,c.appname as appname,c.confkey as confkey,c.context as context,v.confvalue as confvalue from project_confs as c,app_env_config as v where c.project=? and c.appname=? and v.peid=? and c.id=v.pcid", c.Project, c.Appname, n).Error
	return H, err
}

//部署管理
//type ProjectDeploy struct {
//	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
//	Project  string `gorm:"type:varchar(50);index:pdl;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
//	Appname  string `gorm:"type:varchar(50);index:pdl;comment:'组件名称'" binding:"omitempty,min=3,max=50" json:"appname"`
//	Duser    string `gorm:"type:varchar(50);comment:'部署用户'" binding:"omitempty,min=3,max=50" json:"duser"`
//	Package  string `gorm:"type:varchar(50);index:pdl;comment:'部署包'" binding:"omitempty,min=3,max=50" json:"package"`
//	Denv     string `gorm:"type:varchar(20);index:pdl;comment:'部署环境'" binding:"omitempty,min=3,max=20" json:"denv"`
//	State    int    `gorm:"type:tinyint;default:0;comment:'部署状态'" binding:"gte=0,lte=10" json:"state"`
//	Manifest string `gorm:"type:varchar(65535);comment:'部署细节'" json:"hmanifest"`
//	Ctime    string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
//}

//新增部署
func (p *ProjectDeploy) ProjectManageProjectDeployCreate() error {
	return Db.Model(&p).Create(&p).Error
}

//更新部署
func (p *ProjectDeploy) ProjectManageProjectDeployUpdate() error {
	return Db.Model(&p).Where("`project` = ? and `appname` = ? and `package` = ? and `denv` = ?", p.Project, p.Appname, p.Package, p.Denv).Update(&p).Error
}

type LimitOffset struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

//查询部署列表(管理员查询)

//查询项目部署列表(管理员查询)

//查询项目应用部署列表(管理员查询)

//查询环境部署列表(管理员查询)

//查询项目环境部署列表(管理员查询)

//查询项目应用环境部署次数
func (p *ProjectDeploy) ProjectManageProjectDeployGetNumByPAE() (int, error) {
	var H []ProjectDeploy
	err := Db.Model(&p).Where("`project` = ? and `appname` = ? and `denv` = ?", p.Project, p.Appname, p.Denv).Scan(&H).Error
	if H != nil {
		i := 0
		for _, v := range H {
			if v.Project != "" {
				i++
			}
		}
		return i, err
	} else {
		return 0, err
	}
}

//查询项目应用环境部署列表
func (p *ProjectDeploy) ProjectManageProjectDeployGetListByPAE(l LimitOffset) ([]ProjectDeploy, error) {
	var H []ProjectDeploy
	if l.Limit <= 0 {
		l.Limit = 10
	}
	if l.Offset <= 0 {
		l.Offset = 1
	}
	err := Db.Model(&p).Where("`project` = ? and `appname` = ? and `denv` = ?", p.Project, p.Appname, p.Denv).Order("id desc").Limit(l.Limit).Offset((l.Offset - 1) * l.Limit).Scan(&H).Error
	return H, err
}

//查询指定包部署环境列表
func (p *ProjectDeploy) ProjectManageProjectDeployGetListByPackage() ([]string, error) {
	var H []ProjectDeploy
	var R []string
	err := Db.Model(&p).Where("`project` = ? and `appname` = ? and `package` = ?", p.Project, p.Appname, p.Package).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Denv)
		}
		return R, err
	} else {
		return nil, err
	}
}
