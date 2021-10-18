package v1

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

var OkMethod = []string{"POST", "GET"}
var Urllist []string              //系统接口列表
var UrlMethodlist []UrlMethodList //系统接口-method方法列表

type UrlMethodList struct {
	UrlPath string `json:"urlpath"`
	Method  string `json:"method"`
}

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

//分页结构体
type SplitPage struct {
	SplitAll
	Stat     int    `json:"stat"`
	Ip       string `json:"Ip"`
	Hostname string `json:"hostname"`
}

type SplitAll struct {
	Num  int `json:"num"`
	Page int `json:"page"`
}

//md5哈希
func Hx(a string) string {
	hash := md5.New()
	hash.Write([]byte(a))
	return hex.EncodeToString(hash.Sum(nil))
}

//base64编码处理
func AllPassProcessByEncode(pass string) string {
	return base64.StdEncoding.EncodeToString([]byte(pass))
}

//base64解码处理
func AllPassProcessBydecode(pass string) string {
	re, _ := base64.StdEncoding.DecodeString(pass)
	return string(re)
}

//角色权限管理
type RM struct {
	Role   string `json:"role"`
	Path   string `json:"path"`
	Method string `json:"method"`
}

//检查method方法是否合法
func CheckMethod(m string) bool {
	for _, f := range OkMethod {
		if f == m {
			return true
			break
		}
	}
	return false
}

//创建目录，如果目录不存在则创建
func NotExistCreatePath(File string) {
	p := path.Dir(File)
	_, err := os.Stat(p)
	if err != nil {
		_ = os.MkdirAll(p, 0777)
	}
}

//ssh通过账号密码连接
func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}
