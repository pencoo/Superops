package main

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func DBinit() {
	var err error
	Db, err = gorm.Open("mysql", Mysqlinfo)
	if err != nil {
		panic("数据库连接失败！error: " + fmt.Sprint(err))
	}
	Db.Debug()
	Db.SingularTable(true)
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
	Db.AutoMigrate(&WafLock{}, &WafList{})
}

type WafRole struct {
	Id        int    `gorm:"primory_key;auto_increment" json:"id"`
	Rolename  string `gorm:"type:varchar(50);uniqueIndex;comment:'规则名称'" json:"rolename"`
	Key       string `gorm:"type:varchar(50);uniqueIndex;comment:'监听字段，ip、url、domain'" json:"key"`
	Type      string `gorm:"type:varchar(50);uniqueIndex;comment:'持续时间'" json:"type"`
	Condition string `gorm:"type:varchar(50);uniqueIndex;comment:'条件：大于、小于、等于'" json:"condition"`
	Action    string `gorm:"type:varchar(50);uniqueIndex;comment:'加入黑名单、白名单、禁用ip、禁用url、限速ip、限速url'" json:"action"`
}

type WafLock struct {
	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
	Lockname string `gorm:"type:varchar(50);uniqueIndex;comment:'锁名称'" json:"lockname"`
	Lockstat int    `gorm:"type:varchar(4);comment:'锁状态，1开启，2关闭'" json:"lockstat"`
}

type WafList struct {
	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
	Listname string `gorm:"type:varchar(50);uniqueIndex;comment:'列表名称'" json:"listname"`
	Listinfo string `gorm:"type:varchar(5000);comment:'列表内存'" json:"listinfo"`
}

//根据名称查询单个锁状态
func (w *WafLock) FindLockStatus() error {
	return Db.Model(&w).Where("`lockname` = ?", w.Lockname).Find(&w).Error
}

type Locknames struct {
	Lockname string
}

//查询开启状态锁列表
func (w *WafLock) FindLockEnableList() ([]string, error) {
	var H []Locknames
	var R []string
	err := Db.Model(&w).Select("lockname").Where("`lockstat` = 1").Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Lockname)
		}
		return R, err
	} else {
		return nil, err
	}
}

//查询关闭状态锁列表
func (w *WafLock) FindLockDisableList() ([]string, error) {
	var H []Locknames
	var R []string
	err := Db.Model(&w).Select("lockname").Where("`lockstat` = 2").Scan(&H).Error
	if H != nil {
		for _, v := range H {
			R = append(R, v.Lockname)
		}
		return R, err
	} else {
		return nil, err
	}
}

//根据名称开锁
func (w *WafLock) EnableLock() error {
	return Db.Model(&w).Where("`lockname` = ?", w.Lockname).Update("lockstat", 1).Error
}

//根据名称关锁
func (w *WafLock) DisableLock() error {
	return Db.Model(&w).Where("`lockname` = ?", w.Lockname).Update("lockstat", 2).Error
}

//根据名称查询单个列表
func (w *WafList) FindWafListInfo() error {
	return Db.Model(&w).Where("`listname` = ?", w.Listname).Find(&w).Error
}

//查询所有列表
func (w *WafList) FindWafListAllInfo() ([]WafList, error) {
	var H []WafList
	err := Db.Model(&w).Scan(&H).Error
	return H, err
}

//根据名称修改单个列表
func (w *WafList) UpdateWafListInfo() error {
	return Db.Model(&w).Where("`listname` = ?", w.Listname).Update("listinfo", w.Listinfo).Error
}

//获取白名单列表
func (w *WafList) DataIsInWhiteList(s string) bool {
	var R, b []string
	var H []WafList
	Db.Model(&w).Where("`listname` = ? or `listname` = ?", "WhiteIp", "WhiteUrl").Scan(&H)
	if H != nil {
		for _, v := range H {
			a := AllPassProcessBydecode(v.Listinfo)
			_ = json.Unmarshal([]byte(a), &b)
			if b != nil {
				R = append(R, b...)
			}
		}
	}
	for _, k := range R {
		if k == s {
			return true
		}
	}
	return false
}
