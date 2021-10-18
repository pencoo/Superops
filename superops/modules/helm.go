package modules

//kubernetes部署模板
type KubeTmps struct {
	Id    int    `gorm:"primory_key;auto_increment" json:"id"`
	Name  string `gorm:"type:varchar(50)" binding:"omitempty,min=3,max=50" json:"Name"`
	Ktmp  string `gorm:"type:longblob" json:"Ktmp"`
	Ctime string `gorm:"type:varchar(20)" json:"ctime"`
	Cuser string `gorm:"type:varchar(50)" json:"cuser"`
	State int    `gorm:"type:tinyint;default:1" binding:"gte=0,lte=10" json:"status"`
}

//部署资源模板限制cpu、内存,例如lage、xlage等
type DeployTemplate struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Dtname  string `gorm:"type:varchar(50);comment:'部署名称'"`
	Cpumin  string ``
	Cpumax  string ``
	Memmin  string ``
	Memmax  string ``
	Context string `` //描述文档，选择此文档显示此模板cpu和内存使用范围
}

//部署健康检查模板，包含启动钩子，停止钩子，健康钩子
type DeployHealthCheck struct {
	Id       int    `gorm:"primory_key;auto_increment" json:"id"`
	HcName   string `` //监控检查模板名称
	Htype    string `` //探针类型，LivenessProbe存活探针，ReadinessProbe就绪探针
	Ctype    string `` //健康检查类型，cmd、http、sock
	Cport    string `` //http类型的请求端口，默认为http
	Cpath    string `` //http类型的请求path路径
	Ccmd     string `` //cmd时的命令
	LinitSec int    `` //存活探针开始检测延迟,默认10秒
	Ltimeout int    `` //存活探针检查超时时间,默认1秒
	Lperiod  int    `` //存活探针检查时间间隔,默认10秒
	LSTF     int    `` //失败状态(failed)连续多少次检查成功后转为成功状态(successful)，默认为1
	LFTS     int    `` //成功状态(successful)连续多少次检查失败后转为失败状态(failed)，默认为3
	Context  string `` //描述文档，当选择此模板显示此模板需要程序怎么实现检查探针
}

//部署目录挂载
type DeployVolume struct {
	Id            int    `gorm:"primory_key;auto_increment" json:"id"`
	Assemblyid    int    `` //组件id
	Penv          int    `` //部署环境
	Vtype         string `` //卷类型，emptyDir、hostPath、nfs
	Volumesname   string `` //卷名称，任何类型都需要
	Volumesmounts string `` //容器挂载路径
	Volumespath   string `` //hostPath、nfs模式下的目录path
	Volumestype   string `` //hostPath模式下的type，可以是"DirectoryOrCreate、Directory、FileOrCreate、File"
	Nfsserver     string `` //nfs类型时nfs服务器地址
}

//部署label
type DeployLabel struct {
	Id         int    `gorm:"primory_key;auto_increment" json:"id"`
	Assemblyid int    `` //组件id
	Penv       int    `` //部署环境
	LabelKey   string `` //标签key
	LabelValue string `` //标签值
}

//暴露端口，默认80端口暴露不需要添加
type DeployPort struct {
	Id         int    `gorm:"primory_key;auto_increment" json:"id"`
	Assemblyid int    `` //组件id
	Ingress    int    `` //是否需要使用ingress暴露
	Cport      int    `` //容器暴露的端口
	Sport      int    `` //server暴露的端口
	Nport      int    `` //kube node节点端口
	Proport    string `` //端口类型，TCP、UDP
}
