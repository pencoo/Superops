package modules

import (
	"fmt"
	"superops/libs/config"
	L "superops/middlewares/ginzap"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"go.uber.org/zap"
)

var Db *gorm.DB
var TL *time.Location

func DBinit() {
	dbString := config.Dbinfo
	var err error
	Db, err = gorm.Open("mysql", dbString)
	if err != nil {
		L.Lzap.Error("数据库连接失败", zap.String("error", fmt.Sprint(err)))
		panic("数据库连接失败！error: " + fmt.Sprint(err))
	} else {
		L.Lzap.Info("数据库连接成功", zap.String("Dbinfo", config.Dbinfo))
	}
	TL, _ = time.LoadLocation(config.Timezone)
	Db.Debug()
	Db.LogMode(true)
	Db.SetLogger(&L.GormLogger{})
	Db.SingularTable(true)
	Db.DB().SetMaxIdleConns(10)
	Db.DB().SetMaxOpenConns(100)
	tableInit()
}

func tableInit() {
	Db.AutoMigrate(&Users{}, &Userapitoken{}, &Harbors{})
	Db.AutoMigrate(&Kubecmds{}, &KubernetesConfigs{}, &Projects{}, &ProjectUsers{}, &ProjectApps{}, &ProjectEnvs{}, &ProjectConfs{}, &AppEnvConfig{}, &ProjectDeploy{})
	Db.AutoMigrate(&Jenkins{}, &JenkinsJobs{}, &JenkinsFolders{}, &JenkinsBuilds{}, &JenkinsBuildsTips{}, &JenkinsTmps{})
}

//索引管理
// 用户表：
// alter table users add unique key(`username`);
// alter table users add unique key(`phone`);
// alter table users add unique key(`email`);
