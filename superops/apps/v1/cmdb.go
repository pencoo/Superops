package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"superops/libs/e"
	L "superops/middlewares/ginzap"
	"superops/modules"
	"time"

	nmap "github.com/Ullaakut/nmap/v2"
	"github.com/gin-gonic/gin"
)

// @Summary 增加主机扫描网段
// @Description 必须参数：Name、Ipmask
// @Description     Name：描述网段用途
// @Description     Ipmask：网段信息，例如：192.168.0.0/16
// @Description user和pass为空时需要ssh信任
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Ipscan body modules.Ipscan true "主机扫描网段信息"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/createscanhost [post]
func CmdbScanHostsCreate(c *gin.Context) {
	var S modules.Ipscan
	if err := c.ShouldBindJSON(&S); err != nil {
		L.Errlog(L.ErrorInfo{Info: "新增主机扫描网段失败-数据绑定失败", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err, Data: S})
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, err, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if S.Name != "" && S.Ipmask != "" {
			S.Id = 0
			if err := S.CmdbIpscanCreate(); err != nil {
				//新增失败
				L.Errlog(L.ErrorInfo{Info: "新增主机扫描网段失败-数据库写入失败", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err, Data: S})
				c.JSON(200, Response{e.FAILED_INSERT_DB, err, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				L.Infolog(L.SuccessInfo{Info: "新增主机扫描网段成功", User: c.GetString("Name"), Path: c.Request.URL.Path, Data: S})
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			}
		} else {
			//有必须参数为空
			L.Errlog(L.ErrorInfo{Info: "新增主机扫描网段失败-有必须参数为空", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err, Data: S})
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		}
	}
}

// @Summary 修改主机扫描网段
// @Description 必须字段：name且不能修改。
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param modules.Ipscan body modules.Ipscan true "主机扫描网段信息"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/updatescanhost [post]
func CmdbScanHostsUpdate(c *gin.Context) {
	var S modules.Ipscan
	if err := c.ShouldBindJSON(&S); err != nil {
		L.Errlog(L.ErrorInfo{Info: "更新主机扫描网段失败-数据绑定失败", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err, Data: S})
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, err, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if S.Name != "" {
			S.Id = 0
			if err := S.CmdbIpscanUpdate(); err != nil {
				//更新失败
				L.Errlog(L.ErrorInfo{Info: "更新主机扫描网段失败-数据库写入失败", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err, Data: S})
				c.JSON(200, Response{e.FAILED_INSERT_DB, err, e.GetMsg(e.FAILED_INSERT_DB)})
			} else {
				L.Infolog(L.SuccessInfo{Info: "更新主机扫描网段成功", User: c.GetString("Name"), Path: c.Request.URL.Path, Data: S})
				c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
			}
		} else {
			//有必须参数为空
			L.Errlog(L.ErrorInfo{Info: "更新主机扫描网段失败-有必须参数为空", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err, Data: S})
			c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
		}
	}
}

// @Summary 删除主机扫描网段
// @Description 请求方式：get /api/v1/cmdb/deletescanhost?name=xxx
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/deletescanhost [get]
func CmdbScanHostsDelete(c *gin.Context) {
	var S modules.Ipscan
	if S.Name = c.Query("name"); S.Name == "" {
		L.Errlog(L.ErrorInfo{Info: "删除主机扫描网段失败-参数为空", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: errors.New("请求参数为空"), Data: S})
		c.JSON(200, Response{e.FAILED_PARAMS_ANAL, nil, e.GetMsg(e.FAILED_PARAMS_ANAL)})
	} else {
		if err := S.CmdbIpscanDelete(); err != nil {
			//新增失败
			L.Errlog(L.ErrorInfo{Info: "删除主机扫描网段失败-数据库写入失败", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err, Data: S})
			c.JSON(200, Response{e.FAILED_INSERT_DB, err, e.GetMsg(e.FAILED_INSERT_DB)})
		} else {
			L.Infolog(L.SuccessInfo{Info: "删除主机扫描网段成功", User: c.GetString("Name"), Path: c.Request.URL.Path, Data: S})
			c.JSON(200, Response{e.SUCCESS, nil, e.GetMsg(e.SUCCESS)})
		}
	}
}

// @Summary 查询现有主机网段
// @Description 请求方式：get /api/v1/cmdb/getscanhostlist?active=xxx
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param active query int false "查询网段列表，不带则查询所有，active=1表示查询激活，2表示查询停用"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/getscanhostlist [get]
func CmdbScanHostsGetList(c *gin.Context) {
	var S modules.Ipscan
	i := c.Query("active")
	if i == "" || i == "1" || i == "2" {
		S.Status, _ = strconv.Atoi(i)
		R, err := S.CmdbIpscanGetList()
		if err != nil {
			L.Errlog(L.ErrorInfo{Info: "查询主机网段失败-数据库查询失败", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err})
			c.JSON(200, Response{e.FAILED_SELECT_DB, nil, e.GetMsg(e.FAILED_SELECT_DB)})
		} else {
			L.Infolog(L.SuccessInfo{Info: "查询主机网段成功", User: c.GetString("Name"), Path: c.Request.URL.Path})
			c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
		}
	} else {
		L.Errlog(L.ErrorInfo{Info: "", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: errors.New("请求参数错误")})
		c.JSON(200, Response{e.FAILED_PARAMS_NULL, nil, e.GetMsg(e.FAILED_PARAMS_NULL)})
	}
}

// @Summary 获取网段名称列表
// @Description 请求方式：get /api/v1/cmdb/getscanhostnamelist
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/getscanhostnamelist [get]
func CmdbScanHostsNameList(c *gin.Context) {
	var S modules.Ipscan
	R := S.CmdbIpscanGetNameList()
	L.Infolog(L.SuccessInfo{Info: "查询启动名称列表成功", User: c.GetString("Name"), Path: c.Request.URL.Path})
	c.JSON(200, Response{e.SUCCESS, R, e.GetMsg(e.SUCCESS)})
}

// @Summary 启动网段扫描
// @Description 请求方式：get /api/v1/cmdb/beginscanhost?name=xxx
// @Description 获取名称列表接口：get /api/v1/cmdb/getscanhostnamelist
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/beginscanhost [get]
func CmdbScanHostsStart(c *gin.Context) {
	s := c.Query("name")
	err := BeginScanVhost(s, c.GetString("Name"), c.Request.URL.Path)
	if err != nil {
		L.Errlog(L.ErrorInfo{Info: "启动网络扫描失败", User: c.GetString("Name"), Path: c.Request.URL.Path, Err: err})
		c.JSON(200, Response{e.FAILED_PARAMS_NULL, err, e.GetMsg(e.FAILED_PARAMS_NULL)})
	} else {
		L.Infolog(L.SuccessInfo{Info: "开始网络扫描", User: c.GetString("Name"), Path: c.Request.URL.Path, Data: s})
		c.JSON(200, Response{e.SUCCESS, "开始扫描", e.GetMsg(e.SUCCESS)})
	}
}

//s：扫描的网段名称，名称为空扫描所有记录网段
//u：扫描用户，robot表示机器人
//p：发动扫描的url。机器人扫描时此字段为空
func BeginScanVhost(s string, u string, p string) error {
	var S modules.Ipscan
	if s != "" {
		Do := false
		R := S.CmdbIpscanGetNameList()
		for _, v := range R {
			if s == v {
				Do = true
				break
			}
		}
		if Do {
			S.Name = s
			_ = S.CmdbIpscanGetInfo()
			go ScanVhost(S.Ipmask, u, S.Name)
			return nil
		} else {
			L.Errlog(L.ErrorInfo{Info: "扫描失败-" + s, User: u, Path: p, Err: errors.New("扫描名称不存在")})
			return errors.New("扫描名称不存在")
		}
	} else {
		dolist := S.CmdbIpscanGetIpList()
		for _, v := range dolist {
			go ScanVhost(v.Ipmask, u, v.Name)
		}
		return nil
	}
}

//s：扫描的IP段
//u：开启扫描的用户
//n：扫描网段名称
func ScanVhost(s string, u string, n string) {
	var H modules.ScanHistory
	var ListHostIP []string
	H.Starttime = time.Now().Format("2006-01-02 15:04:05")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(s),
		nmap.WithPingScan(),
		nmap.WithContext(ctx),
	)
	if err != nil {
		fmt.Println("错误：", err)
	} else {
		result, _, err := scanner.Run()
		if err != nil {
			fmt.Println("run nmap scan failed: %v", err)
		} else {
			for _, host := range result.Hosts {
				// 查询出所有在线 IP
				ip := fmt.Sprintf("%s", host.Addresses[0])
				// 返回给数组
				ListHostIP = append(ListHostIP, ip)
			}
		}
	}
	if ListHostIP != nil {
		list, _ := json.Marshal(ListHostIP)
		H.Succlist = string(list)
		H.Status = 1
	} else {
		H.Status = 2
		H.Context = "扫描结果为空"
	}
	H.Name = n + "-" + Hx(H.Starttime)
	H.User = u
	H.Stoptime = time.Now().Format("2006-01-02 15:04:05")
	_ = H.CmdbScanHistoryCreate()
	go VhostManagerByScan(ListHostIP, H.Name, DotProcessMask(s))
}

//点分ip处理函数，可通过输入ip加掩码获取ip掩码字符串，只支持8/16/24掩码,默认掩码是24
func DotProcessMask(mask string) string {
	if mask != "" {
		fen := strings.Split(mask, "/")
		ip, m := fen[0], fen[1]
		ip_dot := strings.Split(ip, ".")
		if m == "8" {
			return ip_dot[0] + "*"
		} else if m == "16" {
			return ip_dot[0] + "." + ip_dot[1] + "*"
		} else {
			return ip_dot[0] + "." + ip_dot[1] + "." + ip_dot[2] + "*"
		}
	} else {
		return ""
	}
}

//通过列表完成cmdb建设,传入扫描记录名称.修改后将结果写入Context
func VhostManagerByScan(iplist []string, ScanhistoryName string, s string) {
	if iplist != nil {
		var h modules.Vhostinfo
		//DisableHostList：未激活服务器被扫描到列表
		//enableHostList：激活状态服务器未被扫描到列表
		//NewHostList：新增服务器列表
		var DisableHostList, enableHostList, NewHostList []string
		//查询s网段下激活服务器列表
		EnableHostList := h.CmdbVhostGetList(1, s)
		//查询s网段下所有服务器列表
		AllHostList := h.CmdbVhostGetList(0, s)
		//分析扫描到的服务器
		for _, v := range iplist {
			f := true
			//判断扫描服务器是否在激活列表中
			if EnableHostList != nil {
				for _, k := range EnableHostList {
					if k.Hostip == v {
						f = false
						break
					}
				}
			}
			//判断什么服务器是否在所有服务器中
			if f == true && AllHostList != nil {
				for _, j := range AllHostList {
					if j.Hostip == v {
						DisableHostList = append(DisableHostList, v)
						f = false
						break
					}
				}
			}
			//判断为新增服务器
			if f == true {
				NewHostList = append(NewHostList, v)
			}
		}
		//获取激活状态服务器扫描失败列表
		for _, v := range EnableHostList {
			f := true
			for _, k := range iplist {
				if v.Hostip == k {
					f = false
					break
				}
			}
			if f == true {
				enableHostList = append(enableHostList, v.Hostip)
			}
		}
		//标记未激活服务器扫描成功
		if DisableHostList != nil {
			for _, v := range DisableHostList {
				h := modules.Vhostinfo{Hostip: v}
				_ = h.CmdbVhostUpdateStatusScan(1)
			}
		}
		//标记激活服务器扫描失败
		if enableHostList != nil {
			for _, v := range enableHostList {
				h := modules.Vhostinfo{Hostip: v}
				_ = h.CmdbVhostUpdateStatusScan(2)
			}
		}
		//处理新增服务器
		if NewHostList != nil {
			//扫描主机信息，并添加主机
			for _, v := range NewHostList {
				go AddNewVHost(v, ScanhistoryName)
			}
		}
	}
}

//将扫描到的新服务器添加到数据库
func AddNewVHost(ip string, scaninfo string) {
	var h modules.Vhostinfo
	h.Hostip, h.Status, h.Scanstat, h.Ctime = ip, 0, 1, time.Now().Format("2006-01-02 15:04:05")
	s := modules.Ipscan{Name: scaninfo}
	if s.CmdbIpscanGetInfo() == nil && s.User != "" && s.Pass != "" {
		if s.Sshport == 0 {
			s.Sshport = 22
		}
		session, err := SSHConnect(s.User, s.Pass, ip, s.Sshport)
		if err != nil {
			L.Errlog(L.ErrorInfo{Info: "robot扫描添加服务器失败ssh连接失败。ip：" + ip, Data: "扫描名称：" + scaninfo})
			_ = h.CmdbVhostCreate()
		} else {
			var stdOut, stdErr bytes.Buffer
			defer session.Close()
			session.Stdout = &stdOut
			session.Stderr = &stdErr
			_ = session.Run("hostname;cat /proc/cpuinfo | grep processor | wc -l;cat /proc/cpuinfo | grep 'model name' | tail -1 | awk -F ': ' '{print $2}';free -m | grep Mem | awk '{print $2}'")
			r := strings.Split(strings.Replace(stdOut.String(), "\n", "!!", -1), "!!")
			h.HostName, h.Cputype = r[0], r[2]
			h.Cpunum, _ = strconv.Atoi(r[1])
			h.Memsize, _ = strconv.Atoi(r[3])
			_ = h.CmdbVhostCreate()
		}
	} else {
		//远程执行命令获取结果
	}
}

// @Summary 获取主机列表
// @Description 请求方法：get /api/v1/cmdb/gethostlist
// @Description num表示查询条数，默认查询20条，且最多查询100条
// @Description page表示查询页码，默认查询第1页
// @Description hostname通过过滤主机名称前缀过滤主机
// @Description ip通过过滤ip前缀过滤主机
// @Description state通过状态查询，1表示激活主机，2表示未激活主机
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Param SplitPage body SplitPage false "输入条数和页码，默认显示第1页，显示20条数据"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/gethostlist [POST]
func CmdbGetHostList(c *gin.Context) {
	var h modules.Vhostinfo
	var p SplitPage
	_ = c.ShouldBindJSON(&p)
	if p.Num < 10 {
		p.Num = 20
	}
	if p.Num > 100 {
		p.Num = 100
	}
	if p.Stat != 1 && p.Stat != 2 {
		p.Stat = 0
	}
	h.HostName, h.Hostip, h.Status = p.Hostname, p.Ip, p.Stat
	if p.Page < 1 || p.Page-1 > h.CmdbVhostGetNum()/p.Num {
		p.Page = 0
	} else {
		p.Page = p.Page - 1
	}
	H, err := h.CmdbVhostGetListSplit(p.Num, p.Page)
	if err != nil {
		c.JSON(200, Response{e.FAILED_SELECT_DB, err, e.GetMsg(e.FAILED_SELECT_DB)})
	} else {
		c.JSON(200, Response{e.SUCCESS, H, e.GetMsg(e.SUCCESS)})
	}
}

// @Summary 获取主机名称列表
// @Description 用户注销：此功能只允许用户自己操作
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/gethostsnamelist [get]
func CmdbGetHostNameList(c *gin.Context) {

}

// @Summary 增加主机
// @Description 用户注销：此功能只允许用户自己操作
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/createhost [get]
func CmdbCreateHosts(c *gin.Context) {

}

// @Summary 修改主机信息
// @Description 用户注销：此功能只允许用户自己操作
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/updatehostinfo [get]
func CmdbHostsInfo(c *gin.Context) {

}

// @Summary 删除主机信息
// @Description 用户注销：此功能只允许用户自己操作
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/deletehost [get]
func CmdbDelHostsInfo(c *gin.Context) {

}

// @Summary 修改主机状态
// @Description 用户注销：此功能只允许用户自己操作
// @Tags cmdb
// @Produce json
// @Param Access-Token header string true "用户token"
// @Success 200 {object} Response
// @Router /api/v1/cmdb/updatehoststatus [get]
func CmdbUpdateHostStatus(c *gin.Context) {

}
