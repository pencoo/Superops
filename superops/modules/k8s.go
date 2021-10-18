package modules

import (
	"time"

	"github.com/jinzhu/gorm"
)

//kubeconfig配置管理
//kenvRole：开发可选择的集群应该小于等于1，测试可选择的集群应该小于等于2，运维可选择的集群应该大于2小于等于4
type KubernetesConfigs struct {
	Id        int    `gorm:"primory_key;auto_increment" json:"id"`
	Kname     string `gorm:"type:varchar(50);uniqueIndex:kname;comment:'配置名称'" binding:"omitempty,min=3,max=50" json:"kname"`
	Comname   string `gorm:"comment:'kubectl命令名称，对应kubecmds中的id'" binding:"omitempty,min=3,max=50" json:"comname"`
	Kcontexts string `gorm:"type:varchar(100);comment:'集群介绍说明'" binding:"omitempty,min=3,max=100" json:"kcontexts"`
	Kconfig   string `gorm:"type:longblob;comment:'kubernetes上下文配置'" json:"kconfig"`
	Ctime     string `gorm:"type:varchar(20);comment:'配置创建时间'" json:"ctime"`
	KenvRole  int    `gorm:"type:tinyint;comment:'环境角色，1表示开发测试环境，2表示压测环境，3表示灰度环境，4表示生产环境'" binding:"gte=0,lte=10" json:"kenvrole"`
	Confpath  string `gorm:"type:varchar(200);comment:'配置文件存放路径，此参数不需要传递'" binding:"omitempty,min=10,max=200" json:"confpath"`
	State     int    `gorm:"type:tinyint;default:1;comment:'配置状态，1表示启用，2表示禁用'" binding:"gte=0,lte=10" json:"state"`
}

//标签和污点管理
type TagAndTaint struct {
	Id int `gorm:"primory_key;auto_increment" json:"id"`
}

//服务器管理
type MethineList struct {
	Id int `gorm:"primory_key;auto_increment" json:"id"`
}

func (*Kubecmds) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}

//新增kubeconfig配置
func (k *KubernetesConfigs) KubeconfigCreate() error {
	return Db.Model(&k).Create(&k).Error
}

//判断kubeconfig是否存在
func (k *KubernetesConfigs) KubeconfigExist() bool {
	Db.Model(&k).Where("`kname` = ?", k.Kname).Find(&k)
	if k.Id > 0 {
		return true
	} else {
		return false
	}
}

//删除kubeconfig
func (k *KubernetesConfigs) KubeconfigDelete() error {
	return Db.Model(&k).Where("`kname` = ?", k.Kname).Delete(&k).Error
}

//更新kubeconfig
func (k *KubernetesConfigs) KubeconfigUpdate() error {
	return Db.Model(&k).Where("`kname` = ?", k.Kname).Update(&k).Error
}

//启用/禁用kubeconfig配置,i=1启用，i=2禁用
func (k *KubernetesConfigs) KubeconfigControl(i int) error {
	return Db.Model(&k).Where("`kname` = ?", k.Kname).Update("state", i).Error
}

//获取kubeconfig 名称列表
type kcn struct {
	Kname string
}

//开发查询输入1，测试查询输入2，运维查询输入3，流程自动化输入3
func KubeconfigGetNameList(r int) ([]string, error) {
	var H []kcn
	var R []string
	var err error
	if r < 3 {
		err = Db.Model(KubernetesConfigs{}).Where("`kenvrole` <= ? and `state` = 1", r).Scan(&H).Error
	} else {
		err = Db.Model(KubernetesConfigs{}).Where("`kenvrole` <= 4 and `kenvrole` >= 3 and `state` = 1").Scan(&H).Error
	}

	if H != nil && err == nil {
		for _, v := range H {
			R = append(R, v.Kname)
		}
		return R, err
	} else {
		return nil, err
	}
}

//查询kubeconfig列表(表格展示)
func KubeconfigGetList() ([]KubernetesConfigs, error) {
	var H []KubernetesConfigs
	return H, Db.Model(KubernetesConfigs{}).Scan(&H).Error
}

//查询kubeconfig单条
func (k *KubernetesConfigs) KubeconfigGetInfo() error {
	return Db.Model(&k).Where("`Kname` = ?", k.Kname).Find(&k).Error
}

//获取kubeconfig列表,i=1查询激活的配置，i=2查询禁用的配置
func KubeconfigList(i int) ([]KubernetesConfigs, error) {
	var H []KubernetesConfigs
	return H, Db.Model(KubernetesConfigs{}).Where("`state` = ?", i).Scan(&H).Error
}
