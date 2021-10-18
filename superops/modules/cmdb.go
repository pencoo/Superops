package modules

import "errors"

//数据库表
//机房信息
type SupinfoForm struct {
	Id        string      `gorm:"primory_key;auto_increment" json:"id"`
	Envname   string      `gorm:"type:varchar(50);commit:'环境名称'" json:"envname"`
	Context   string      `gorm:"type:varchar(30);commit:'环境概述，例如阿里云、腾讯云、公司机房等'" json:"context"`
	Url       string      `gorm:"type:varchar(100);commit:'云环境登录地址'" json:"url"`
	Logininfo []Logininfo `gorm:"type:varchar(500);commit:'登录信息'" json:"logininfo"`
	Support   []Supinfo   `gorm:"type:varchar(500);commit:'服务支持方联系信息'" json:"support"`
	Subnet    []Subnet    `gorm:"type:varchar(500);commit:'子网信息，通常是[127.0.0.1/24]'" json:"subnet"`
}
type Supinfo struct {
	Id        string `gorm:"primory_key;auto_increment" json:"id"`
	Envname   string `gorm:"type:varchar(50);commit:'环境名称'" json:"envname"`
	Context   string `gorm:"type:varchar(30);commit:'环境概述，例如阿里云、腾讯云、公司机房等'" json:"context"`
	Url       string `gorm:"type:varchar(100);commit:'云环境登录地址'" json:"url"`
	Logininfo string `gorm:"type:varchar(500);commit:'登录信息'" json:"logininfo"`
	Support   string `gorm:"type:varchar(500);commit:'服务支持方联系信息'" json:"support"`
	Subnet    string `gorm:"type:varchar(500);commit:'子网信息，通常是[127.0.0.1/24]'" json:"subnet"`
}

//云账号登录信息,可能会有多个，建议列表使用
type Logininfo struct {
	Id      string `commit:"登录id，无则可忽略" json:"id"`
	User    string `commit:"登录用户名" json:"user"`
	Pass    string `commit:"登录密码" json:"pass"`
	Context string `commit:"登录描述" json:"context"`
	Ctime   string `commit:"账号创建时间" json:"ctime"`
}

//服务支持方信息,可能会有多个建议列表使用
type Support struct {
	Name  string `commit:"支持方名称" json:"name"`
	Phone string `commit:"支持方电话" json:"phone"`
	Email string `commit:"支持方邮件" json:"email"`
	Qq    string `commit:"支持方qq信息" json:"qq"`
	Ctime string `commit:"支持方创建时间" json:"ctime"`
}

//子网信息，通常是list类型用于cmdb扫描
type Subnet struct {
	Name string `commit:"子网信息" json:"name"`
	Net  string `commit:"子网信息,用于cmdb自动扫描，例如：'192.168.0.0/24'" json:"net"`
	Neti string `commit:"跨网段扫描代理地址" json:"neti"`
}

//数据库表
//物理服务器
type HostInfo struct {
	Id            int    `gorm:"primory_key;auto_increment" json:"id"`
	Hostname      string `gorm:"type:varchar(70);commit:'主机名称'" json:"hostname"`
	Roomid        int    `gorm:"type:int;commit:'机房id'" json:"roomid"`
	Interfaceinfo string `gorm:"type:varchar(300);commit:'网络接口名称'" json:"interfaceinfo"`
	Cpuinfo       string `gorm:"type:varchar(200);commit:'cpu信息'" json:"cpuinfo"`
	Cpunum        int    `gorm:"type:tinyint;commit:'cpu数量'" json:"cpunum"`
	Cpuall        int    `gorm:"type:smallint;commit:'cpu核心数量'" json:"cpuall"`
	Meminfo       string `gorm:"type:varchar(100);commit:'内存型号'" json:"meminfo"`
	Memnum        int    `gorm:"type:tinyint;commit:'内存数量'" json:"memnum"`
	Memall        int    `gorm:"type:smallint;commit:'内存总数'" json:"memall"`
	Diskinfo      string `gorm:"type:varchar(700);commit:'磁盘信息'" json:"diskinfo"`
	Masterip      string `gorm:"type:varchar(15);commit:'服务器主IP'" json:"masterip"`
	Manageip      string `gorm:"type:varchar(15);commit:'管理ip'" json:"manageip"`
	Ctime         string `gorm:"type:varchar(20);commit:'服务器购买时间'" json:"ctime"`
	Status        string `gorm:"type:tinyint;commit:'状态，1运行，2停机，3暂停'" json:"status"`
}

//网卡信息,可能会有多个，建议使用列表
type Interfaceinfo struct {
	Ifname  string  `commit:"网络卡名称，通常是pci插槽位_网卡类型" json:"ifname"`
	Ifspeed string  `commit:"接口速率，通常是速率类型，例如e100，f1000,l1000,单位mbps" json:"ifspeed"`
	Buytime string  `commit:"购买时间" json:"buytime"`
	Ifinfo  []Ifuse `commit:"接口使用信息使用Ifuse结构" json:"ifinfo"`
}

//接口使用
type Ifuse struct {
	Id  int    `commit:"接口号" json:"id"`
	Ip  string `commit:"ip地址" json:"ip"`
	Mac string `commit:"mac地址" json:"mac"`
}

//cpu信息，可能会有多个建议使用列表
type Cpuinfo struct {
	Soltid  int    `commit:"cpu插槽ID" json:"soltid"`
	Type    string `commit:"cpu型号" json:"type"`
	Hz      int    `commit:"cpu主频，单位MHz" json:"hz"`
	Core    int    `commit:"cpu核心数" json:"core"`
	L1      int    `commit:"cpu一级缓存单位KB" json:"l1"`
	L2      int    `commit:"cpu二级缓存单位KB" json:"l2"`
	L3      int    `commit:"cpu三级缓存单位KB" json:"l3"`
	Buytime string `commit:"cpu购买时间" json:"buytime"`
}

//磁盘信息，可能会有多个建议使用列表
type Diskinfo struct {
	Type      string     `commit:"磁盘类型，mechain：机械，SSD：ssd磁盘，PCIE：PCIE磁盘等" json:"type"`
	Raid      []int      `commit:"Raid支持列表" json:"raid"`
	Raidtype  string     `commit:"raid类型，当前使用的raid类型" json:"raidtype"`
	Disktype  string     `commit:"磁盘型号" json:"disktype"`
	Disksize  int        `commit:"磁盘大小，单位MB" json:"disksize"`
	Spacesize int        `commit:"总容量，单位GB" json:"spacesize"`
	Mountpath string     `commit:"挂载路径" json:"mountpath"`
	Disklist  []Disklist `commit:"磁盘具体信息" json:"disklist"`
}

type Disklist struct {
	Sole    int    `commit:"槽位" json:"sole"`
	Buytime string `commit:"购买时间" json:"buytime"`
}

//数据库表
type Ipscan struct {
	Id      int    `gorm:"primory_key;auto_increment" json:"id"`
	Name    string `gorm:"type:varchar(70);commit:'ip段名称'" binding:"required,min=4,max=70" json:"name"`
	Ipmask  string `gorm:"type:varchar(20);commit:'ip段'" binding:"omitempty,min=10,max=16" json:"ipmask"`
	Context string `gorm:"type:varchar(100);commit:'ip段描述'" json:"context"`
	User    string `gorm:"type:varchar(20);commit:'登录服务器账号'" json:"user"`
	Pass    string `gorm:"type:varchar(100);commit:'登录服务器密码'" json:"pass"`
	Sshport int    `gorm:"type:int;commit:'ssh登录端口'" json:"sshport"`
	Status  int    `gorm:"type:tinyint;default:1;commit:'状态,1启用扫描，2停止扫描'" binding:"omitempty,1|2" json:"status"`
}

//扫描记录
type ScanHistory struct {
	Id        int    `gorm:"primory_key;auto_increment" json:"id"`
	Name      string `gorm:"type:varchar(70);commit:'ip段名称'" binding:"required,min=4,max=70" json:"name"`
	User      string `gorm:"type:varchar(50);commit:'扫描用户，robot是自动扫描'" json:"user"`
	Starttime string `gorm:"type:varchar(30);commit:'开始扫描时间'" json:"starttime"`
	Stoptime  string `gorm:"type:varchar(30);commit:'扫描完成时间'" json:"stoptime"`
	Succlist  string `gorm:"type:varchar(200);commit:'扫描到的机器'" json:"succlist"`
	Context   string `gorm:"type:varchar(100);commit:'扫描信息'" json:"context"`
	Status    int    `gorm:"type:tinyint;commit:'扫描状态，1表示成功，2表示失败'" json:"status"`
}

//数据库表
//虚拟服务器、云服务器
type Vhostinfo struct {
	Id       string `gorm:"primory_key;auto_increment" json:"id"`
	HostName string `gorm:"type:varchar(70);commit:'服务器名称'" binding:"" json:"host_name"`
	Hostip   string `gorm:"type:varchar(15);commit:'服务器ip'" json:"hostip"`
	Cputype  string `gorm:"type:varchar(30);commit:'cpu型号'" json:"cputype"`
	Cpunum   int    `gorm:"type:tinyint;commit:'cpu核心数量'" json:"cpunum"`
	Memsize  int    `gorm:"type:int;commit:'内存大小，单位MB'" json:"memsize"`
	Disksys  int    `gorm:"type:int;commit:'系统盘大小，单位GB'" json:"disksys"`
	//Disklists string `gorm:"type:varchar(500);commit:'数据盘信息'" json:"disklists"`
	//Mhostid   int    `gorm:"type:int;commit:'宿主服务器Id'" json:"mhostid"`
	//Cloudid  int    `gorm:"type:int;commit:'云服务器时为云机房id'" json:"cloudid"`
	Netspeed int    `gorm:"type:smallint;commit:'网卡速率，单位Mbps'" json:"netspeed"`
	Ctime    string `gorm:"type:varchar(50);commit:'服务器创建时间'" json:"ctime"`
	Hostenv  string `gorm:"type:varchar(10);commit:'服务器环境，例如dev、test、pro等'" json:"hostenv"`
	//Stime    string `gorm:"type:varchar(50);commit:'服务器停机时间，对应停机或暂停状态'" json:"stime"`
	Status   int `gorm:"type:tinyint;default:1;commit:'服务器状态，0扫描到待确认，1运行，2停机，3暂停'" json:"status"`
	Scanstat int `gorm:"type:tinyint;default:1;commit:'扫描状态，1运行，2停机'" json:"scanstat"`
}

//虚拟服务器数据磁盘列表
type Disklists struct {
	Size    int    `commit:"磁盘大小" json:"size"`
	Type    int    `commit:"磁盘类型，mechain：机械，SSD：ssd磁盘，PCIE：PCIE磁盘" json:"type"`
	Path    string `commit:"磁盘挂载路径" json:"path"`
	Context string `commit:"磁盘描述信息" json:"context"`
}

//数据库表
//ip扫描信息
//type Ipscan struct {
//	Id      string `gorm:"primory_key;auto_increment" json:"id"`
//	Name    string `gorm:"type:varchar(70);commit:'ip段名称'" json:"name"`
//	Ipmask  string `gorm:"type:varchar(20);commit:'ip段'" json:"ipmask"`
//	Context string `gorm:"type:varchar(100);commit:'ip段描述'" json:"context"`
//	Status  string `gorm:"type:tinyint;commit:'状态,1启用扫描，2停止扫描'" json:"status"`
//}

//增加ip扫描信息
func (ip *Ipscan) CmdbIpscanCreate() error {
	return Db.Model(&ip).Create(&ip).Error
}

//修改ip扫描信息(根据name字段修改)
func (ip *Ipscan) CmdbIpscanUpdate() error {
	return Db.Model(&ip).Where("`name` = ?", ip.Name).Update(&ip).Error
}

//启用和禁用扫描
func (ip *Ipscan) CmdbIpscanControl(i int) error {
	return Db.Model(&ip).Where("`name` = ?", ip.Name).Update("state", i).Error
}

//删除ip扫描信息
func (ip *Ipscan) CmdbIpscanDelete() error {
	return Db.Where("`name` = ?", ip.Name).Delete(&Ipscan{}).Error
}

//查询ip扫描列表,ip.Status=1查询激活中的扫描列表，ip.Status=2查询未激活列表，其它值表示查询所有
func (ip *Ipscan) CmdbIpscanGetList() ([]Ipscan, error) {
	var H []Ipscan
	var err error
	if ip.Status == 1 || ip.Status == 2 {
		err = Db.Model(&ip).Where("`status` = ?", ip.Status).Scan(&H).Error
	} else {
		err = Db.Model(&ip).Scan(&H).Error
	}
	return H, err
}

//查询扫描地址名称列表
func (ip *Ipscan) CmdbIpscanGetNameList() []string {
	var H []Ipscan
	var R []string
	Db.Model(&ip).Select("name").Where("`state` = ?", 1).Scan(&H)
	for _, v := range H {
		R = append(R, v.Name)
	}
	return R
}

//查询扫描地址列表
func (ip *Ipscan) CmdbIpscanGetIpList() []Ipscan {
	var H []Ipscan
	Db.Model(&ip).Where("`state` = ?", 1).Scan(&H)
	return H
}

//查询单条扫描信息
func (ip *Ipscan) CmdbIpscanGetInfo() error {
	return Db.Model(&ip).Where("`state` = ?", 1).Where("`name` = ?", ip.Name).Find(&ip).Error
}

//扫描记录
//type ScanHistory struct {
//	Id        int    `gorm:"primory_key;auto_increment" json:"id"`
//	Name      string `gorm:"type:varchar(70);commit:'ip段名称'" binding:"required,min=4,max=70" json:"name"`
//	User      string `gorm:"type:varchar(50);commit:'扫描用户，robot是自动扫描'" json:"user"`
//	Starttime string `gorm:"type:varchar(30);commit:'开始扫描时间'" json:"starttime"`
//	Stoptime  string `gorm:"type:varchar(30);commit:'扫描完成时间'" json:"stoptime"`
//	Succlist  string `gorm:"type:varchar(200);commit:'扫描到的机器'" json:"succlist"`
//	Context   string `gorm:"type:varchar(100);commit:'扫描信息'" json:"context"`
//	Status    int    `gorm:"type:tinyint;commit:'扫描状态，1表示成功，2表示失败'" json:"status"`
//}
//增加扫描记录
func (s *ScanHistory) CmdbScanHistoryCreate() error {
	return Db.Model(&s).Create(&s).Error
}

//查询扫描成功记录

//查询扫描失败记录

//查询扫描记录，user不空可查询指定用户的扫描记录，默认显示最近10行

//数据库表
//虚拟服务器、云服务器
//type Vhostinfo struct {
//	Id       string `gorm:"primory_key;auto_increment" json:"id"`
//	HostName string `gorm:"type:varchar(70);commit:'服务器名称'" json:"host_name"`
//	Hostip   string `gorm:"type:varchar(15);commit:'服务器ip'" json:"hostip"`
//	Cputype  string `gorm:"type:varchar(30);commit:'cpu型号'" json:"cputype"`
//	Cpunum   int    `gorm:"type:tinyint;commit:'cpu核心数量'" json:"cpunum"`
//	Memsize  int    `gorm:"type:int;commit:'内存大小，单位MB'" json:"memsize"`
//	Disksys  int    `gorm:"type:int;commit:'系统盘大小，单位GB'" json:"disksys"`
//	Netspeed int    `gorm:"type:smallint;commit:'网卡速率，单位Mbps'" json:"netspeed"`
//	Ctime    string `gorm:"type:varchar(50);commit:'服务器创建时间'" json:"ctime"`
//	Status int `gorm:"type:tinyint;commit:'服务器状态，1运行，2停机，3暂停'" json:"status"`
//  Scanstat int `gorm:"type:tinyint;commit:'扫描状态，1运行，2停机'" json:"scanstat"`
//}

//增加服务器
func (h *Vhostinfo) CmdbVhostCreate() error {
	return Db.Model(&h).Create(&h).Error
}

//修改服务器信息
func (h *Vhostinfo) CmdbVhostUpdate() error {
	return Db.Model(&h).Where("`hostname` = ?", h.HostName).Update(&h).Error
}

//删除服务器信息
func (h *Vhostinfo) CmdbVhostDelete() error {
	return Db.Where("`hostname` = ?", h.HostName).Delete(&h).Error
}

//web修改服务器状态，i=1表示运行服务器，i=2表示服务器停机，i=3表示暂停
func (h *Vhostinfo) CmdbVhostUpdateStatus(i int) error {
	if i == 1 || i == 2 {
		return Db.Model(&h).Where("`hostname` = ?", h.HostName).Update("status", i).Error
	} else {
		return errors.New("param is need 1 or 2 or 3 by int")
	}
}

//扫描修改服务器状态，i=1表示运行服务器，i=2表示服务器停机
func (h *Vhostinfo) CmdbVhostUpdateStatusScan(i int) error {
	if i == 1 || i == 2 {
		return Db.Model(&h).Where("`hostip` = ?", h.Hostip).Update("scanstat", i).Error
	} else {
		return errors.New("param is need 1 or 2 by int")
	}
}

//查询服务器列表,i=1表示查询运行服务器，i=2表示查询停机服务器，i=3表示查询暂停服务器，其它值表示查询所有
func (h *Vhostinfo) CmdbVhostGetList(i int, ip_prefix string) []Vhostinfo {
	var H []Vhostinfo
	if ip_prefix != "" {
		if i == 1 || i == 2 || i == 3 {
			Db.Model(&h).Where("`status` = ?", i).Where("`hostip` = ?", ip_prefix).Scan(&H)
		} else {
			Db.Model(&h).Where("`hostip` = ?", ip_prefix).Scan(&H)
		}
	} else {
		if i == 1 || i == 2 || i == 3 {
			Db.Model(&h).Where("`status` = ?", i).Scan(&H)
		} else {
			Db.Model(&h).Scan(&H)
		}
	}
	return H
}

//分页查询服务器列表
func (h *Vhostinfo) CmdbVhostGetListSplit(num, page int) ([]Vhostinfo, error) {
	var H []Vhostinfo
	var err error
	if h.HostName == "" && h.Hostip == "" {
		err = Db.Model(&h).Limit(num).Offset(page * num).Order("id desc").Scan(&H).Error
	}
	if h.HostName != "" && h.Hostip != "" {
		err = Db.Model(&h).Where("`hostname` like ?, `hostip` like ?", h.HostName+"%", h.Hostip+"%").Limit(num).Offset(page * num).Order("id desc").Scan(&H).Error
	}
	if h.HostName != "" && h.Hostip == "" {
		err = Db.Model(&h).Where("`hostname` like ?", h.HostName+"%").Limit(num).Offset(page * num).Order("id desc").Scan(&H).Error
	}
	if h.HostName == "" && h.Hostip != "" {
		err = Db.Model(&h).Where("`hostip` like ?", h.Hostip+"%").Limit(num).Offset(page * num).Order("id desc").Scan(&H).Error
	}
	return H, err
}

//查询指定IP前缀服务器列表
func (h *Vhostinfo) CmdbVhostGetListByIpPrefix(prefix string) {
}

//查询指定主机名前缀服务器列表
func (h *Vhostinfo) CmdbVhostGetListByHostnamePrefix(prefix string) {
}

//查询所有服务器Ip列表
func (h *Vhostinfo) CmdbVhostGetIpList() {
}

//查询指定IP前缀服务器Ip列表
func (h *Vhostinfo) CmdbVhostGetIpListByIpPrefix() {
}

//查询服务器记录数量hostname,ip用于判断hostname或ip前缀，s用于判断查询激活还是非激活
func (h *Vhostinfo) CmdbVhostGetNum() int {
	var i int
	//主机名和ip都为空判断
	if h.HostName == "" && h.Hostip == "" {
		if h.Status == 1 || h.Status == 2 {
			Db.Model(&h).Where("`status` = ?", h.Status).Count(&i)
			return i
		} else {
			Db.Model(&h).Count(&i)
			return i
		}
	} else {
		if h.HostName != "" && h.Hostip != "" {
			if h.Status == 1 || h.Status == 2 {
				Db.Model(&h).Where("`status` = ?, `hostname` like ?, `hostip` like ?", h.Status, h.HostName+"%", h.Hostip+"%").Count(&i)
				return i
			} else {
				Db.Model(&h).Where("`hostname` like ?, `hostip` like ?", h.HostName+"%", h.Hostip+"%").Count(&i)
				return i
			}
		} else {
			if h.HostName != "" && h.Hostip == "" {
				if h.Status == 1 || h.Status == 2 {
					Db.Model(&h).Where("`status` = ?, `hostname` like ?", h.Status, h.HostName+"%").Count(&i)
					return i
				} else {
					Db.Model(&h).Where("`hostname` like ?", h.HostName+"%").Count(&i)
					return i
				}
			}
			if h.HostName == "" && h.Hostip != "" {
				if h.Status == 1 || h.Status == 2 {
					Db.Model(&h).Where("`status` = ?, `hostip` like ?", h.Status, h.Hostip+"%").Count(&i)
					return i
				} else {
					Db.Model(&h).Where("`hostip` like ?", h.Hostip+"%").Count(&i)
					return i
				}
			}
		}
	}
	return 0
}
