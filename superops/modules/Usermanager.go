package modules

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

type Users struct {
	Id        int    `gorm:"primory_key;auto_increment;comment:'用户id'" json:"id"`
	Username  string `gorm:"type:varchar(50);uniqueIndex:username;comment:'姓名'" binding:"omitempty,min=3,max=50" json:"username"`
	Autograph string `gorm:"type:varchar(50);comment:'用户签名'" binding:"omitempty,min=3,max=50" json:"autograph"`
	Email     string `gorm:"type:varchar(50);uniqueIndex:email;comment:'邮箱'" binding:"omitempty,email" json:"email"`
	Phone     string `gorm:"type:varchar(11);uniqueIndex:login;comment:'手机'" binding:"omitempty,len=11" json:"phone"`
	Context   string `gorm:"type:varchar(255);comment:'用户描述信息'" binding:"omitempty,min=1,max=100" json:"context"`
	Password  string `gorm:"type:varchar(100);comment:'密码'" binding:"omitempty,min=6,max=100" json:"password"`
	Ctime     string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	Utime     string `gorm:"type:varchar(20);comment:'修改时间'" json:"utime"`
	Ltime     string `gorm:"type:varchar(20);comment:'禁用/启用时间'" json:"ltime"`
	Ustat     int    `gorm:"type:tinyint;default:0;comment:'用户状态，0、注册 1、激活 2、禁用 3、删除'" binding:"gte=0,lte=10" json:"ustat"`
}

type UsersInfo struct {
	Username  string `gorm:"type:varchar(50);unique:username;comment:'姓名'" binding:"omitempty,min=3,max=50" json:"username"`
	Autograph string `gorm:"type:varchar(50);comment:'用户签名'" binding:"omitempty,min=3,max=50" json:"autograph"`
	Email     string `gorm:"type:varchar(50);unique:email;comment:'邮箱'" binding:"omitempty,min=6,max=50" json:"email"`
	Phone     string `gorm:"type:varchar(11);unique:login;comment:'手机'" binding:"omitempty,len=11" json:"phone"`
	Context   string `gorm:"type:varchar(255);comment:'用户描述信息'" binding:"omitempty,min=1,max=100" json:"context"`
}

type UStat struct {
	Ctime string `gorm:"type:varchar(20);comment:'创建时间'" json:"ctime"`
	Utime string `gorm:"type:varchar(20);comment:'修改时间'" json:"utime"`
	Ltime string `gorm:"type:varchar(20);comment:'禁用/启用时间'" json:"ltime"`
	Ustat int    `gorm:"type:tinyint;default:0;comment:'用户状态，0、注册 1、激活 2、禁用 3、删除'" binding:"gte=0,lte=10" json:"ustat"`
}

type UName struct {
	Username string `gorm:"type:varchar(50);unique:username;comment:'姓名'" binding:"omitempty,min=3,max=50" json:"username"`
}
type ULists struct {
	UsersInfo
	UStat
}

//userapi token存储
type Userapitoken struct {
	Id       int    `gorm:"primory_key;auto_increment;comment:'Token id'" json:"id"`
	Username string `gorm:"type:varchar(50);comment:'用户名'" json:"uname"`
	Token    string `gorm:"type:varchar(255);comment:'token信息'" json:"token"`
	Tstat    int    `gorm:"type:tinyint;default:1;comment:'token状态，1启用，2禁用'" json:"tstat"`
}

type UserTokenInfo struct {
	Username string `gorm:"type:varchar(50);comment:'用户名'" json:"uname"`
	Token    string `gorm:"type:varchar(255);comment:'token信息'" json:"token"`
}

func (u *Users) BeforeCreate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Ctime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}
func (u *Users) BeforeUpdate(scope *gorm.Scope) error {
	_ = scope.SetColumn("Utime", time.Now().In(TL).Format("2006-01-02 15:04:05"))
	return nil
}

//查询用户是否存在且为激活状态
func (u *Users) UserExist() bool {
	if Db.Model(&u).Where("`Username` = ? and `ustat` = 1", u.Username).Find(&u).Error != nil {
		return false
	} else {
		if u.Id != 0 {
			return true
		} else {
			return false
		}
	}
}

//用户登录验证,支持：用户名+密码、手机号+密码、邮箱+密码
func (u *Users) LoginAuth() bool {
	if u.Username != "" {
		Db.Where("`username` = ? and `password` = ? and `ustat` = 1", u.Username, u.Password).Find(&u)
	} else if u.Phone != "" {
		Db.Where("`phone` = ? and `password` = ? and `ustat` = 1", u.Phone, u.Password).Find(&u)
	} else if u.Email != "" {
		Db.Where("`email` = ? and `password` = ? and `ustat` = 1", u.Email, u.Password).Find(&u)
	} else {
		return false
	}
	if u.Id != 0 {
		return true
	} else {
		return false
	}
}

//用户注册
func (u *Users) UsersRegister() error {
	return Db.Create(&u).Error
}

//根据id查询用户名
func (u *Users) SelectUserNameById() {
	Db.Model(&u).Where("`id` = ?", u.Id).Find(&u)
}

//更新用户信息
func (u *Users) UsersInfoUpdate() error {
	return Db.Model(&u).Where("`id` = ?", u.Id).Update(&u).Error
}

//密码修改
func (u *Users) UsersChangePass() error {
	return Db.Model(&u).Where("`username` = ?", u.Username).Update("password", u.Password).Error
}

//获取用户信息,需要携带uid
func (u *Users) UsersInfos() (UsersInfo, error) {
	var H UsersInfo
	var err error
	err = Db.Model(&u).Select("username,autograph,email,phone,context,ctime,utime,ltime,ustat").Where("`id` = ?", u.Id).Scan(&H).Error
	return H, err
}

//获取用户名列表,n=1：激活用户列表，n=2：禁用用户列表
func (u *Users) GetUsersLists(n int) []string {
	var H []UName
	var R []string
	Db.Debug().Model(&u).Select("username").Where("`ustat` = ?", n).Scan(&H)
	for _, v := range H {
		R = append(R, v.Username)
	}
	return R
}

//管理员获取用户列表
//t为查询类型，0、查询新注册用户,1、查询激活用户，2、查询禁用用户，3、查询所有用户
//l为查询条数，默认20条
//o为查询分页，默认第1页
//od为查询排序列
func (u *Users) GetUsersListsByAdmin(t, l, o int, od string) ([]ULists, error) {
	var H []ULists
	var err error
	if l == 0 {
		l = 20
	}
	if o == 0 {
		o = 1
	}
	if t >= 0 && t <= 2 {
		if od != "" {
			err = Db.Model(&u).Select("username,autograph,email,phone,context,ctime,utime,ltime,ustat").Where("`ustat` = ?", t).Order(od).Limit(l).Offset(o * (l - 1)).Scan(&H).Error
		} else {
			err = Db.Model(&u).Select("username,autograph,email,phone,context,ctime,utime,ltime,ustat").Where("`ustat` = ?", t).Order("id").Limit(l).Offset(o * (l - 1)).Scan(&H).Error
		}
	} else if t == 3 {
		if od != "" {
			err = Db.Model(&u).Select("username,autograph,email,phone,context,ctime,utime,ltime,ustat").Order(od).Limit(l).Offset(o * (l - 1)).Scan(&H).Error
		} else {
			err = Db.Model(&u).Select("username,autograph,email,phone,context,ctime,utime,ltime,ustat").Order("id").Limit(l).Offset(o * (l - 1)).Scan(&H).Error
		}
	}
	return H, err
}

//用户锁定
func (u *Users) UserDoActivateOrDisable(i int) string {
	u.Ltime = time.Now().In(TL).Format("2006-01-02 15:04:05")
	if i == 0 {
		u.Ustat = 1 //激活
	} else {
		u.Ustat = 2 //禁用
	}
	return u.UserStatusUpdate()
}

//用户解锁/激活
//func (u *Users) UserUnlock() string {
//	u.Ltime = time.Now().Format("2006-01-02 15:04:05")
//	u.Ustat = 1
//	return u.UserStatusUpdate()
//}

//用户状态更新
func (u *Users) UserStatusUpdate() string {
	if err := Db.Model(&u).Where("`username` = ?", u.Username).Update("ustat", u.Ustat).Error; err == nil {
		return ""
	} else {
		return fmt.Sprint(err)
	}
}

//查询角色列表
//type CasbinRoleList struct {
//	Role string
//}
//
//func SelectRoleList() ([]string, error) {
//	var H []CasbinRoleList
//	var R []string
//	err := Db.Debug().Raw("select v0 from casbin_rule where p_type='p' group by v0").Scan(&H).Error
//	for _, v := range H {
//		R = append(R, v.Role)
//	}
//	return R, err
//}

//判断角色权限是否存在
//func SelectRolePowerExist(csb Casbinrule) (Casbinrule, error) {
//	var H Casbinrule
//	err := Db.Debug().Raw("select * from casbin_rule where p_type='p' and v0 = ? and v1 = ? and v2 = ?", csb.V0, csb.V1, csb.V2).Scan(&H).Error
//	return H, err
//}

//删除一条角色权限
//func DeleteRolePowerExist(csb Casbinrule) bool {
//	err := Db.Debug().Raw("delete from casbin_rule where p_type='p' and v0 = ? and v1 = ? and v2 = ?", csb.V0, csb.V1, csb.V2).Error
//	if err != nil {
//		return false
//	} else {
//		return true
//	}
//}

//查询一个用户API Token是否有效
func (u *Userapitoken) UsersApiTokenIsActive() bool {
	err := Db.Debug().Model(&u).Where("`username` = ? and `token` = ? and `tstat` = 1", u.Username, u.Token).Find(&u).Error
	if u.Id == 0 || err != nil {
		return false
	} else {
		return true
	}
}

//新增一个用户API Token
func (u *Userapitoken) UsersAddApiToken() error {
	if !u.UserApiTokenExist() {
		return Db.Debug().Create(&u).Error
	} else {
		return errors.New("exist token")
	}
}

func (u Userapitoken) UserApiTokenExist() bool {
	Db.Debug().Model(Userapitoken{}).Where("`username` = ?", u.Username).Find(&u)
	if u.Id == 0 {
		return false
	} else {
		return true
	}
}

//禁用一个用户API Token
func (u *Userapitoken) UsersDisableApiToken() error {
	return Db.Debug().Model(&u).Where("`username` = ?", u.Username).Update("tstat", 2).Error
}

//启用一个用户API Token
func (u *Userapitoken) UsersEnableApiToken() error {
	return Db.Debug().Model(&u).Where("`username` = ?", u.Username).Update("tstat", 1).Error
}

//删除一个用户API Token
func (u *Userapitoken) UsersDeleteApiToken() error {
	return Db.Debug().Model(&u).Where("`username` = ?", u.Username).Delete(&u).Error
}

//获取api token
func (u *Userapitoken) UsersGetApiToken() error {
	return Db.Debug().Model(&u).Where("`username` = ?", u.Username).Find(&u).Error
}

//查询激活的API Token列表
func UsersApiTokenLists() []UserTokenInfo {
	var H []UserTokenInfo
	Db.Debug().Model(Userapitoken{}).Select("username,token").Where("`tstat` = 1").Scan(&H)
	return H
}
