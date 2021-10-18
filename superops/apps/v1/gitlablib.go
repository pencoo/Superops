package v1

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//git := "http://gitlab.jhongnet.com:8888/pencoo/app.git"
//token := "oGwVFu6S5WSyPQGRxhva"
//fmt.Println("分支列表：", BranchLists(git,token))
//fmt.Println("项目id：", GitProjectId(git,token))

//获取分支列表
//输入url是gitlab仓库地址:
//    http://gitlab.jhongnet.com:8888/ops/sailing-pipeline.git
//    git@gitlab.jhongnet.com:ops/sailing-pipeline.git
//token:项目的gitlab token
func BranchLists(url string, token string) []string {
	apiurl := GitUrlGetServiceUrl(url) + "/api/v4/projects/" + GitProjectId(url, token) + "/repository/branches?private_token=" + token
	resp, _ := http.Get(apiurl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var Re []map[string]interface{}
	var Return []string
	_ = json.Unmarshal(body, &Re)
	for _, v := range Re {
		Return = append(Return, v["name"].(string))
	}
	return Return
}

//获取项目id
//输入url是gitlab仓库地址:
//    http://gitlab.jhongnet.com:8888/ops/sailing-pipeline.git
//    git@gitlab.jhongnet.com:ops/sailing-pipeline.git
//token:项目的gitlab token
func GitProjectId(url string, token string) string {
	apiurl := GitUrlGetServiceUrl(url) + "/api/v4/projects/?search=" + GitUrlGetProjectName(url) + "&private_token=" + token
	resp, _ := http.Get(apiurl)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var Re []map[string]interface{}
	_ = json.Unmarshal(body, &Re)
	if Re != nil {
		for _, v := range Re {
			if v["ssh_url_to_repo"].(string) == url || v["http_url_to_repo"].(string) == url || v["web_url"].(string) == url {
				return fmt.Sprint(v["id"].(float64))
			}
		}
	}
	return ""
}

//获取git地址的项目名称
//提取如下：
//    http://gitlab.jhongnet.com:8888/ops/sailing-pipeline.git ==> sailing-pipeline
//    git@gitlab.jhongnet.com:ops/sailing-pipeline.git         ==> sailing-pipeline
//传入参数为git地址，输出为项目名称
func GitUrlGetProjectName(g string) string {
	var r []string
	if strings.HasPrefix(g, "http") {
		r = strings.Split(g, "/")
	} else if strings.HasPrefix(g, "git@") {
		r = strings.Split(g, ":")
		if strings.Contains(r[len(r)-1], "/") {
			r = strings.Split(r[len(r)-1], "/")
		}
	}
	if len(r) > 0 {
		r1 := r[len(r)-1]
		r2 := strings.Split(r1, ".")
		return r2[0]
	} else {
		return ""
	}
}

//获取gitlab url
//提取如下：
//    http://gitlab.jhongnet.com:8888/ops/sailing-pipeline.git ==> http://gitlab.jhongnet.com:8888
//    git@gitlab.jhongnet.com:ops/sailing-pipeline.git         ==> http://gitlab.jhongnet.com:8888
//传入参数为git地址
func GitUrlGetServiceUrl(g string) string {
	if strings.HasPrefix(g, "http") {
		r := strings.Split(g, "/")
		return r[0] + "//" + r[2]
	} else {
		//当方式为git@的方式时从数据库获取
		return ""
	}
}
