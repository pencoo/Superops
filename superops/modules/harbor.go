package modules

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Harbors struct {
	Id          int    `gorm:"primory_key;auto_increment" json:"id"`
	Harbor      string `gorm:"type:varchar(50);uniqueIndex;comment:'harbor名称'" binding:"omitempty,min=3,max=50" json:"harbor"`
	Kname       string `gorm:"type:varchar(50);comment:'部署harbor'" binding:"omitempty,min=3,max=50" json:"kname"`
	Hurl        string `gorm:"type:varchar(100);comment:'harbor url地址'" binding:"omitempty,min=3,max=100,url" json:"hurl"`
	Huser       string `gorm:"type:varchar(50);comment:'harbor用户名'" binding:"omitempty,min=3,max=50" json:"huser"`
	Hpass       string `gorm:"type:varchar(50);comment:'harbor密码'" binding:"omitempty,min=3,max=50" json:"hpass"`
	Dimgtimeout int    `gorm:"type:int;default:7;comment:'开发(镜像标记为1)镜像过期时间，单位为天，默认7天'" binding:"gte=0,lte=200" json:"dimgtimeout"`
	Timgtimeout int    `gorm:"type:int;default:30;comment:'测试(镜像标记为2)镜像过期时间，单位为天，默认30天'" binding:"gte=0,lte=200" json:"timgtimeout"`
	Oimgtimeout int    `gorm:"type:int;default:90;comment:'运维(镜像标记为3)镜像过期时间，单位为天，默认90天'" binding:"gte=0,lte=200" json:"oimgtimeout"`
	Ctime       string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	State       int    `gorm:"type:tinyint;default:1;comment:'harbor状态，1表示启用，2表示禁用'" binding:"gte=0,lte=10" json:"state"`
}

func (h *Harbors) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}

//新增harbor
func (h *Harbors) HarborCreate() bool {
	if Db.Debug().Create(&h).Error != nil {
		return false
	} else {
		return true
	}
}

//修改harbor
func (h *Harbors) HarborUpdate() bool {
	if Db.Debug().Model(&h).Where("`harbor` = ?", h.Harbor).Update(&h).Error != nil {
		return false
	} else {
		return true
	}
}

//删除harbor
func (h *Harbors) HarborDelete() bool {
	if Db.Debug().Where("`harbor` = ?", h.Harbor).Delete(&h).Error != nil {
		return false
	} else {
		return true
	}
}

//查询harbor列表
type Hl struct {
	Harbor string
}

func (h *Harbors) HarborNameList() ([]string, error) {
	var H []Hl
	var R []string
	err := Db.Model(&h).Where("`state` = 1", h.State).Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Harbor)
		}
	}
	return R, err
}

//harbor列表查询
func (h *Harbors) HarborListAll() ([]Harbors, error) {
	var H []Harbors
	err := Db.Debug().Raw("select * from harbors").Scan(&H).Error
	return H, err
}

//查询单条harbor信息，或者判断项目是否存在
func (h *Harbors) HarborInfo() bool {
	if Db.Debug().Where("`harbor` = ?", h.Harbor).Find(&h).Error != nil {
		return false
	} else {
		if h.Id > 0 {
			return true
		} else {
			return false
		}
	}
}
