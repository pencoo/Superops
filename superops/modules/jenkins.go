package modules

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

//jenkins服务表
type Jenkins struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Name    string `gorm:"type:varchar(50);uniqueIndex;comment:'jenkins名称'" binding:"omitempty,min=3,max=50" json:"name"`
	Jurl    string `gorm:"type:varchar(100);comment:'jenkins url地址'" binding:"omitempty,min=3,max=100,url" json:"jurl"`
	Juser   string `gorm:"type:varchar(50);comment:'jenkins用户名'" binding:"omitempty,min=3,max=50" json:"juser"`
	Jpass   string `gorm:"type:varchar(50);comment:'jenkins密码'" binding:"omitempty,min=3,max=50" json:"jpass"`
	Kname   string `gorm:"type:varchar(50);comment:'部署jenkins的k8s名称'" binding:"omitempty,min=3,max=50" json:"kname"`
	Ctime   string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	Context string `gorm:"type:varchar(500);comment:'说明'" binding:"omitempty,min=3,max=500" json:"context"`
	State   int    `gorm:"type:tinyint;default:1;comment:'状态，1表示启用，2表示禁用，3表示已迁移'" binding:"gte=0,lte=10" json:"state"`
}

//jenkins job管理表
type JenkinsJobs struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Jenkins string `gorm:"type:varchar(50);comment:'jenkins名称name'" binding:"omitempty,min=3,max=50" json:"jenkins"`
	Jfolder string `gorm:"type:varchar(50);comment:'jenkins文件夹'" binding:"omitempty,min=3,max=50" json:"jfolder"`
	Jobname string `gorm:"type:varchar(50);uniqueIndex;comment:'job名称'" binding:"omitempty,min=3,max=50" json:"jobname"`
	Jtmp    string `gorm:"type:varchar(50);comment:'job模板名称'" binding:"omitempty,min=3,max=50" json:"jtmp"`
	Gurl    string `gorm:"type:varchar(50);comment:'组件git地址'" binding:"omitempty,min=3,max=50" json:"gurl"`
	Guser   string `gorm:"type:varchar(50);comment:'组件git用户'" binding:"omitempty,min=3,max=50" json:"Guser"`
	Gpass   string `gorm:"type:varchar(50);comment:'组件git密码'" binding:"omitempty,min=3,max=50" json:"Gpass"`
	State   int    `gorm:"type:tinyint;default:1;comment:'job状态，1表示启用，2表示禁用，3表示已迁移'" binding:"gte=0,lte=10" json:"state"`
	Ctime   string `gorm:"type:varchar(20);comment:'job创建时间'" json:"ctime"`
}

//文件夹表
type JenkinsFolders struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Folder  string `gorm:"type:varchar(50);comment:'文件夹名称，等于项目名称'" binding:"omitempty,min=3,max=50" json:"folder"`
	Jenkins string `gorm:"type:varchar(50);comment:'jenkins名称，对应jenkins表里面的名称'" binding:"omitempty,min=3,max=50" json:"jenkins"`
	Ctime   string `gorm:"type:varchar(20);comment:'文件夹配置创建时间'" json:"ctime"`
}

//jenkins构建信息表
type JenkinsBuilds struct {
	Id        int    `gorm:"primory_key;auto_increment" json:"id"`
	Project   string `gorm:"type:varchar(50);index:project_index;comment:'项目名称'" binding:"omitempty,min=3,max=50" json:"project"`
	Appname   string `gorm:"type:varchar(50);index:project_index;comment:'应用名称'" binding:"omitempty,min=3,max=50" json:"appname"`
	Buildstat int    `gorm:"type:tinyint;default:0;comment:'构建状态，0表示构建中，1表示构建成功，2表示构建失败'" binding:"gte=0,lte=10" json:"buildstat"`
	Jenkins   string `gorm:"type:varchar(50);comment:'构建任务的jenkins名称'" json:"jenkins"`
	Buildinfo string `gorm:"type:longblob;comment:'构建信息'" json:"buildinfo"`
	Gittag    string `gorm:"type:varchar(50);comment:'构建任务的tag标签'" json:"gittag"`
	Gitcommit string `gorm:"type:varchar(50);comment:'构建任务的gitcommit位置'" json:"gitcommit"`
	Package   string `gorm:"type:varchar(50);uniqueIndex;comment:'包名称'" binding:"omitempty,min=3,max=50" json:"package"`
	Ctime     string `gorm:"type:varchar(20);comment:'构建时间'" json:"ctime"`
	Cuser     string `gorm:"type:varchar(50);comment:'构建用户'" binding:"omitempty,min=3,max=50" json:"cuser"`
	Pushtag   int    `gorm:"type:tinyint;default:1;comment:'推送标签，1、开发环境，2、测试环境，3、运维环境'" binding:"gte=0,lte=10" json:"pushtag"`
	Savetag   int    `gorm:"type:tinyint;default:0;comment:'保留标签，0、可删除，1、部署重要版本不可删除，2、已从harbor中删除'" binding:"gte=0,lte=10" json:"savetag"`
}

//jenkins build tips
type JenkinsBuildsTips struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Buildid int    `gorm:"type:int;index:tips;comment:'jenkinsbuild id'" json:"buildid"`
	Tips    string `gorm:"type:varchar(500);comment:'tips提示内容'" json:"tips"`
	Cuser   string `gorm:"type:varchar(50);comment:'tips书写用户'" binding:"omitempty,min=3,max=50" json:"cuser"`
	Atuser  string `gorm:"type:varchar(50);comment:'需要处理此tips的用户'" binding:"omitempty,min=3,max=50" json:"atuser"`
	Ctime   string `gorm:"type:varchar(20);comment:'tips书写时间'" json:"ctime"`
}

//jenkins模板表
type JenkinsTmps struct {
	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
	Tmpname  string `gorm:"type:varchar(50);uniqueIndex;comment:'模板名称'" json:"tmpname"`
	Context  string `gorm:"type:varchar(500);comment:'说明'" binding:"omitempty,min=3,max=500" json:"context"`
	Langtype string `gorm:"type:varchar(50);uniqueIndex;comment:'模板语言'" json:"langtype"`
	Tmp      string `gorm:"type:longblob;comment:'jenkins模板'" json:"tmp"`
	Ctime    string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	Cuser    string `gorm:"type:varchar(50);comment:'创建用户'" json:"cuser"`
	State    int    `gorm:"type:tinyint;comment:'模板状态，1激活，2禁用'" binding:"gte=0,lte=10" json:"state"`
}

func (j *Jenkins) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *JenkinsJobs) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *JenkinsFolders) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *JenkinsBuilds) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (j *JenkinsTmps) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}

//jenkins服务表

//增加jenkins服务
func (j *Jenkins) ServiceCreateJenkins() error {
	return Db.Model(&j).Create(&j).Error
}

//修改jenkins服务
func (j *Jenkins) ServiceUpdateJenkins() error {
	return Db.Model(&j).Where("`name` = ?", j.Name).Update(&j).Error
}

//查询单条jenkins信息/判断jenkins是否存在
func (j *Jenkins) ServiceJenkinsInfo() error {
	if err := Db.Model(&j).Where("`name` = ?", j.Name).Find(&j).Error; err != nil {
		return err
	} else {
		return nil
	}
}

//查询jenkins名称列表
type JenkinsList struct {
	Name string
}

func (j *Jenkins) ServiceJenkinsNameLists() ([]string, error) {
	var H []JenkinsList
	var R []string
	err := Db.Model(&j).Select("name").Where("`state` = 1").Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Name)
		}
	}
	return R, err
}

//查询jenkins列表
func (j *Jenkins) ServiceJenkinsLists() ([]Jenkins, error) {
	var H []Jenkins
	err := Db.Model(&j).Scan(&H).Error
	return H, err
}

//启用、禁用、迁移jenkins
//d=2：启用，d=2：禁用，d=3：迁移
func (j *Jenkins) ServiceJenkinsControls(d int) error {
	return Db.Model(&j).Where("`name` = ", j.Name).Update("state", d).Error
}

//jenkins job管理表
//添加jenkins job
func (j *JenkinsJobs) ServiceJenkinsCreateJob() error {
	return Db.Model(&j).Create(&j).Error
}

//删除jenkins job
func (j *JenkinsJobs) ServiceJenkinsDeleteJob() error {
	return Db.Model(&j).Where("`jenkins` = ?", j.Jenkins).Where("`jobname` = ?", j.Jobname).Delete(&j).Error
}

//修改jenkins job
func (j *JenkinsJobs) ServiceJenkinsUpdateJob() error {
	return Db.Model(&j).Where("`jenkins` = ?", j.Jenkins).Where("`jobname` = ?", j.Jobname).Update(&j).Error
}

//获取jenkins job信息
func (j *JenkinsJobs) ServiceJenkinsGetJobsInfo() error {
	return Db.Model(&j).Where("`jenkins` = ?", j.Jenkins).Where("`jobname` = ?", j.Jobname).Find(&j).Error
}

//获取jenkins job状态
func (j *JenkinsJobs) ServiceJenkinsGetJobsState() bool {
	err := Db.Model(&j).Where("`jenkins` = ?", j.Jenkins).Where("`jobname` = ?", j.Jobname).Find(&j).Error
	if err != nil || j.State != 1 {
		return false
	} else {
		return true
	}
}

//启用、禁用jenkins job
//参数： d=1:启用，d=2：禁用
func (j *JenkinsJobs) ServiceJenkinsControl(d int) error {
	return Db.Model(&j).Where("`jenkins` = ?", j.Jenkins).Where("`jobname` = ?", j.Jobname).Update("state", d).Error
}

//获取jenkins job列表
func (j *JenkinsJobs) ServiceJenkinsJobsList() ([]JenkinsJobs, error) {
	var H []JenkinsJobs
	err := Db.Model(&j).Where("`jenkins` = ?", j.Jenkins).Where("`state` = ?", 1).Scan(&H).Error
	return H, err
}

//jenkins文件夹表
//新增jenkins文件夹
func (j *JenkinsFolders) ServiceJenkinsFolderCreate() error {
	return Db.Model(&j).Create(&j).Error
}

//删除jenkins文件夹
func (j *JenkinsFolders) ServiceJenkinsFolderDelete() error {
	return Db.Model(&j).Where("`folder` = ?", j.Folder).Delete(&j).Error
}

//查询jenkins文件夹列表,j.Jenkins不能为空
type jfs struct {
	Folder string
}

func (j *JenkinsFolders) ServiceJenkinsFolderList() ([]string, error) {
	var H []jfs
	var R []string
	err := Db.Model(&j).Select("folder").Where("`jenkins` = ?", j.Jenkins).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Folder)
		}
	}
	return R, err
}

//jenkins构建信息表
//添加jenkins构建
func (j *JenkinsBuilds) ServiceJenkinsJobBuilding() error {
	return Db.Model(&j).Create(&j).Error
}

//更新jenkins构建状态，更新构建信息和构建状态
func (j *JenkinsBuilds) ServiceJenkinsJobBuildUpdateStatus() error {
	return Db.Raw("update jenkins_builds set buildstat = ?,buildinfo=? where package=?", j.Buildstat, j.Buildinfo, j.Package).Error
}

//更新jenkins构建推送标签,推送标签，1:开发环境，2:测试环境，3:运维环境
func (j *JenkinsBuilds) ServiceJenkinsJobBuildUpdatePushtag(d int) error {
	return Db.Model(&j).Where("`package` = ?", j.Package).Update("pushtag", d).Error
}

//更新jenkins构建保存标签，1：可删除、2：永久保留、3：已删除
func (j *JenkinsBuilds) ServiceJenkinsJobBuildUpdateSavetag(d int) error {
	return Db.Model(&j).Where("`package` = ?", j.Package).Update("savetag", d).Error
}

//查询单条构建信息,包名查询
func (j *JenkinsBuilds) ServiceJenkinsJobBuildGetInfo() error {
	return Db.Model(&j).Where("`package` = ?", j.Package).Find(&j).Error
}

//查询构建包列表(用于表格展示)，需要传入项目名称、应用名称,查询id逆序，默认查询最近30条
func (j *JenkinsBuilds) ServiceJenkinsJobBuildLists(l, o int) ([]JenkinsBuilds, error) {
	var H []JenkinsBuilds
	if l == 0 {
		l = 30
	}
	if o == 0 {
		o = 1
	}
	err := Db.Model(&j).Where("`project` = ? and `appname` = ?", j.Project, j.Appname).Order("id desc").Limit(l).Offset((o - 1) * l).Scan(&H).Error
	return H, err
}

//构建成功/失败查询(用于表格展示)，l表示需要查询的条数，j需要传入项目名称、组件名称、构建状态
func (j *JenkinsBuilds) ServiceJenkinsJobBuildGetSuccess(l int) ([]JenkinsBuilds, error) {
	var H []JenkinsBuilds
	if l == 0 {
		l = 30
	}
	err := Db.Model(&j).Where("`project` = ? and `appname` = ? and `buildstat` = ?", j.Project, j.Appname, j.Buildstat).Order("id desc").Limit(l).Scan(&H).Error
	return H, err
}

//按用户查询构建信息(用于表格展示)，l表示需要查询的条数,j需要传入项目名称、组件名称、构建用户
func (j *JenkinsBuilds) ServiceJenkinsJobBuildByUser(l int) ([]JenkinsBuilds, error) {
	var H []JenkinsBuilds
	if l == 0 {
		l = 30
	}
	err := Db.Model(&j).Where("`project` = ? and `appname` = ? and `cuser` = ?", j.Project, j.Appname, j.Cuser).Order("id desc").Limit(l).Scan(&H).Error
	return H, err
}

//查询可部署包,j需要传入项目名称、组件名称、推送标签值
func (j *JenkinsBuilds) ServiceJenkinsJobBuildPackagesList() ([]string, error) {
	var J []JenkinsBuilds
	var R []string
	var err error
	err = Db.Model(&j).Where("`project` = ? and `appname` = ? and `buildstat` = 1 and `pushtag` = ?", j.Project, j.Appname, j.Pushtag).Scan(&J).Error
	if J != nil {
		for _, v := range J {
			R = append(R, v.Package)
		}
	}
	return R, err
}

type PackageId struct {
	Package string
	Id      int
}

//查询两个package之间的所有包
func (j *JenkinsBuilds) ServiceJenkinsBuildGetBetweenPackage(a, b int) ([]PackageId, error) {
	var J []JenkinsBuilds
	var R []PackageId
	var err error
	err = Db.Model(&j).Where("`project` = ? and `appname` = ? and `buildstat` = 1 and `pushtag` = ? and id between (?,?)", j.Project, j.Appname, j.Pushtag, a, b).Scan(&J).Error
	if J != nil {
		for _, v := range J {
			R = append(R, PackageId{Package: v.Package, Id: v.Id})
		}
	}
	return R, err
}

//构建提示表
//添加构建提示
func (j *JenkinsBuildsTips) ServiceJenkinsTipsCreate() error {
	return Db.Model(&j).Create(&j).Error
}

//查询tips,查询一个package的所有tips
func (j *JenkinsBuildsTips) ServiceJenkinsTipsSelect() []JenkinsBuildsTips {
	var J []JenkinsBuildsTips
	Db.Model(&j).Where("`buildid` = ?", j.Buildid).Scan(&J)
	return J
}

//查询两个包之间的tips,bp表示上一个包，ep表示下一个包
func ServiceJenkinsTipsBetweenMapSelect(bp, ep string) (map[string][]map[string]string, error) {
	B := JenkinsBuilds{Package: bp}
	E := JenkinsBuilds{Package: ep}
	if B.ServiceJenkinsJobBuildGetInfo() == nil && E.ServiceJenkinsJobBuildGetInfo() == nil {
		if B.Id < 0 || E.Id < 0 {
			return nil, errors.New("package不存在")
		} else {
			if B.Pushtag != E.Pushtag {
				return nil, errors.New("package pushtag标志不匹配")
			} else {
				if B.Pushtag == E.Pushtag && B.Appname == E.Appname {
					L, err := B.ServiceJenkinsBuildGetBetweenPackage(B.Id, E.Id)
					if err != nil {
						return nil, errors.New("package 列表查询失败")
					} else {
						var R map[string][]map[string]string
						if L != nil {
							for _, v := range L {
								j := JenkinsBuildsTips{Buildid: v.Id}
								if tl := j.ServiceJenkinsTipsSelect(); tl != nil {
									for _, k := range tl {
										R[v.Package] = append(R[v.Package], map[string]string{"user": k.Cuser, "time": k.Ctime, "tips": k.Tips})
									}
								}
							}
						}
						return R, nil
					}
				} else {
					return nil, errors.New("package项目不匹配")
				}
			}
		}
	} else {
		return nil, errors.New("package查询失败")
	}
}

//jenkins模板表

//添加jenkins模板
func (j *JenkinsTmps) ServiceJenkinsTmpsCreate() error {
	return Db.Model(&j).Create(&j).Error
}

//更新jenkins模板
func (j *JenkinsTmps) ServiceJenkinsTmpsUpdate() error {
	return Db.Model(&j).Where("`tmpname` = ?", j.Tmpname).Update(&j).Error
}

//启用和禁用jenkins模板,d=1:启用，d=2:禁用(禁用后使用此模板的项目奖无法构建)
func (j *JenkinsTmps) ServiceJenkinsTmpsControl(d int) error {
	return Db.Model(&j).Where("`tmpname` = ?", j.Tmpname).Update("state", d).Error
}

//删除jenkins模板(功能未完成待续，需要判断是否有项目在使用此模板，如果有项目在使用则不能删除)
func (j *JenkinsTmps) ServiceJenkinsTmpsDelete() error {
	return nil
}

//查询单条模板信息
func (j *JenkinsTmps) ServiceJenkinsTmpsGetInfo() error {
	if j.Tmpname != "" {
		return Db.Model(&j).Where("`tmpname` = ?", j.Tmpname).Find(&j).Error
	} else {
		return Db.Model(&j).Where("`id` = ?", j.Id).Find(&j).Error
	}
}

type Jtsinfo struct {
	Name    string `json:"name"`
	Devlang string `json:"devlang"`
	Context string `json:"context"`
}

//获取jenkins模板名称列表
func (j *JenkinsTmps) ServiceJenkinsTmpNameList() ([]Jtsinfo, error) {
	var H []JenkinsTmps
	var h []Jtsinfo
	err := Db.Model(&j).Where("`state` = 1").Scan(&H).Error
	if H != nil {
		for _, v := range H {
			h = append(h, Jtsinfo{Name: v.Tmpname, Devlang: v.Langtype, Context: v.Context})
		}
	}
	return h, err
}
