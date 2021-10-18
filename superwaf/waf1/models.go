package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func DBinit() {
	var err error
	Db, err = gorm.Open("mysql", "pencoo:309745197@tcp(106.75.169.216:33306)/pencoo?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic("数据库连接失败！error: " + fmt.Sprint(err))
	}
	Db.Debug()
	Db.SingularTable(true)
	Db.DB().SetMaxIdleConns(3)
	Db.DB().SetMaxOpenConns(10)
	Db.AutoMigrate(&NginxCluster{}, &NginxServerInfo{}, &BlackWhiteList{}, &DropRole{}, &DropList{}, &LimitRole{}, &LimitList{})
}

//nginx集群表
type NginxCluster struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Name    string `gorm:"type:varchar(50);not null;unique_index;commit:'nginx集群名称'" json:"name"`
	Context string `gorm:"type:varchar(100);commit:'集群描述信息'" json:"context"`
	Status  int    `gorm:"type:tinyint;default:1;commit:'集群状态，1表示集群激活，2表示集群禁用'" json:"status"`
}

//nginx集群服务器列表
type NginxServerInfo struct {
	Id     int    `gorm:"primory_key;auto_increment" json:"id"`
	Nid    int    `gorm:"commit:'nginx集群id'" json:"nid"`
	Url    string `gorm:"type:varchar(30);commit:'nginx waf接口地址，例如：http://192.168.1.1:9527'" json:"url"`
	Status int    `gorm:"type:tinyint;default:1;commit:'nginx状态，1表示集群激活，2表示集群禁用，3表示检查失败'" json:"status"`
}

//名单规则表
type BlackWhiteList struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Nid     int    `gorm:"commit:'nginx集群id'" json:"nid"`
	Name    string `gorm:"type:varchar(11);commit:'名单名称，可以是ipwhite、ipblack、urlwhite、urlblack、domainblack、domainwhite'" json:"name"`
	Lvalue  string `gorm:"type:varchar(50);commit:'名单值，单个值'" json:"lvalue"`
	Ctime   string `gorm:"type:varchar(30);commit:'创建时间'" json:"ctime"`
	Context string `gorm:"type:varchar(100);commit:'规则说明'" json:"context"`
	Status  int    `gorm:"type:tinyint;default:1;commit:'名单状态，1表示激活，2表示禁用'" json:"status"`
}

//drop名单规则
type DropRole struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Name    string `gorm:"type:varchar(30);commit:'drop名单名称,推荐名称：Dtype_Dtime'" json:"name"`
	Dtype   string `gorm:"type:varchar(11);commit:'drop类型，可以是ipdrop、domaindrop、urldrop'" json:"dtype"`
	Dtime   int    `gorm:"commit:'drop时间，单位为秒'" json:"dtime"`
	Context string `gorm:"type:varchar(100);commit:'规则描述'" json:"context"`
	Status  int    `gorm:"type:tinyint;default:1;commit:'名单状态，1表示激活，2表示禁用'" json:"status"`
}

//限制规则
type DropList struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Nid     int    `gorm:"commit:'nginx集群id'" json:"nid"`
	Did     int    `gorm:"commit:'drop规则id'" json:"did"`
	Name    string `gorm:"type:varchar(30);commit:'限制规则名称，推荐与Dvalue相关'" json:"name"`
	Dvalue  string `gorm:"type:varchar(50);commit:'规则值，单条'" json:"dvalue"`
	Ctime   string `gorm:"type:varchar(30);commit:'规则创建时间'" json:"ctime"`
	Context string `gorm:"type:varchar(100);commit:'规则描述'" json:"context"`
	Status  int    `gorm:"type:tinyint;default:1;commit:'规则状态，1表示激活，2表示禁用'" json:"status"`
}

//limit名单规则
type LimitRole struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Name    string `gorm:"type:varchar(30);commit:'drop名单名称,推荐格式：Ltype_Lnum_Ltime'" json:"name"`
	Ltype   string `gorm:"type:varchar(11);commit:'limit类型，可以是iplimit、domainlimit、urllimit'" json:"ltype"`
	Lnum    int    `gorm:"commit:'限制访问量'" json:"lnum"`
	Ltime   int    `gorm:"commit:'限制时间，单位为秒'" json:"ltime"`
	Context string `gorm:"type:varchar(100);type:varchar(100);commit:'规则描述'" json:"context"`
	Status  int    `gorm:"type:tinyint;default:1;commit:'名单状态，1表示激活，2表示禁用'" json:"status"`
}

//限流规则
type LimitList struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Nid     int    `gorm:"commit:'nginx集群id'" json:"nid"`
	Lid     int    `gorm:"commit:'limit规则id'" json:"lid"`
	Name    string `gorm:"type:varchar(30);commit:'限流名称'" json:"name"`
	Lvalue  string `gorm:"type:varchar(50);commit:'规则值'" json:"lvalue"`
	Ctime   string `gorm:"type:varchar(30);commit:'规则创建时间'" json:"ctime"`
	Context string `gorm:"type:varchar(100);commit:'规则描述'" json:"context"`
	Status  int    `gorm:"type:tinyint;default:1;commit:'规则状态，1表示激活，2表示禁用'" json:"status"`
}

//NginxCluster表
//创建nginx集群
func (n *NginxCluster) NginxClusterCreate() error {
	return Db.Model(&n).Create(&n).Error
}

//修改集群状态,i=1表示激活，i=2表示禁用
func (n *NginxCluster) NginxClusterUpdateStatus(i int) error {
	return Db.Model(&n).Where("`name` = ?", n.Name).Update("status", i).Error
}

//获取Nginx集群列表,i用于判断状态，无表示查询所有，有且为1表示查询激活状态列表
func (n *NginxCluster) NginxClusterGetList(i int) ([]NginxCluster, error) {
	var H []NginxCluster
	var err error
	if i == 1 {
		err = Db.Model(&n).Where("`status` = 1").Scan(&H).Error
	} else {
		err = Db.Model(&n).Scan(&H).Error
	}
	return H, err
}

//添加nginx机器
func (n *NginxServerInfo) NginxServerInfoCreate() error {
	return Db.Model(&n).Create(&n).Error
}

//删除nginx机器
func (n *NginxServerInfo) NginxServerInfoDelete() error {
	return Db.Raw("delete from nginx_server_info where url = ?", n.Url).Error
}

//修改nginx集群状态,i=1表示激活，i=2表示禁用
func (n *NginxServerInfo) NginxServerInfoUpdateStatus(i int) error {
	return Db.Model(&n).Where("`url` = ?", n.Url).Update("status", i).Error
}

//获取nginx集群下nginx服务器列表,i=1表示激活，i=2表示禁用
func (n *NginxServerInfo) NginxServerInfoGetList(i int) ([]NginxServerInfo, error) {
	var H []NginxServerInfo
	var err error
	if i == 1 {
		err = Db.Model(&n).Where("`status` = ? and `nid` = ?", i, n.Nid).Scan(&H).Error
	} else {
		err = Db.Model(&n).Where("`nid` = ?", n.Nid).Scan(&H).Error
	}
	return H, err
}

//名单规则表
//type BlackWhiteList struct {
//	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
//	Nid     int    `gorm:"commit:'nginx集群id'" json:"nid"`
//	Name    string `gorm:"type:varchar(11);commit:'名单名称，可以是ipwhite、ipblack、urlwhite、urlblack、domainblack、domainwhite'" json:"name"`
//	Lvalue  string `gorm:"type:varchar(50);commit:'名单值，单个值'" json:"lvalue"`
//	Ctime   string `gorm:"type:varchar(30);commit:'创建时间'" json:"ctime"`
//	Context string `gorm:"type:varchar(100);commit:'规则说明'" json:"context"`
//	Status  int    `gorm:"type:tinyint;default:1;commit:'名单状态，1表示激活，2表示禁用'" json:"status"`
//}
