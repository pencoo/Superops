package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type keyInfo struct {
	Url  string //接口地址
	Path string //url地址
	Keyinfo
}
type Keyinfo struct {
	Key   string //键
	Value string //值
	Lock  string //锁
}

//获取列表信息接口
func (k *keyInfo) GetListInfo() (string, error) {
	k.Path = "/getlist"
	resp, err := http.Get(k.Url + k.Path + "?key=" + k.Key)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "{}\n" {
			return "", nil
		} else {
			return strings.Replace(string(body), "\n", "", -1), nil
		}
	}
}

//添加列表元素接口
func (k *keyInfo) AddListInfo() bool {
	k.Path = "/addlist"
	resp, err := http.Get(k.Url + k.Path + "?key=" + k.Key + "&value=" + k.Value)
	defer resp.Body.Close()
	if err != nil {
		return false
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "yes\n" {
			return true
		} else {
			return false
		}
	}
}

//删除列表元素接口
func (k *keyInfo) DelListInfo() bool {
	k.Path = "/dellist"
	resp, err := http.Get(k.Url + k.Path + "?key=" + k.Key + "&value=" + k.Value)
	defer resp.Body.Close()
	if err != nil {
		return false
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "yes\n" {
			return true
		} else {
			return false
		}
	}
}

//查看锁开关状态
func (k *keyInfo) GetLockStatus() bool {
	k.Path = "/lockstate"
	resp, err := http.Get(k.Url + k.Path + "?lock=" + k.Lock)
	defer resp.Body.Close()
	if err != nil {
		return false
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "on\n" {
			return true
		} else {
			return false
		}
	}
}

//开启锁开关
func (k *keyInfo) EnableLock() bool {
	k.Path = "/lockenable"
	resp, err := http.Get(k.Url + k.Path + "?lock=" + k.Lock)
	defer resp.Body.Close()
	if err != nil {
		return false
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "on\n" {
			return true
		} else {
			return false
		}
	}
}

//关闭锁开关
func (k *keyInfo) DisableLock() bool {
	k.Path = "/lockdisable"
	resp, err := http.Get(k.Url + k.Path + "?lock=" + k.Lock)
	if err != nil {
		return false
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) == "off\n" {
			return true
		} else {
			return false
		}
	}
}

//获取内存键信息
func (k *keyInfo) GetMemInfo() (string, error) {
	k.Path = "/getmeminfo"
	resp, err := http.Get(k.Url + k.Path + "?key=" + k.Key)
	if err != nil {
		return "", err
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		if string(body) != "nil\n" {
			return strings.Replace(string(body), "\n", "", -1), nil
		} else {
			return "", nil
		}
	}
}

type WafLog struct {
	DTime   string `json:"dtime"`
	Client  string `json:"client"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (wlog *WafLog) WafDoLog() {
	client := &http.Client{}
	req, _ := json.Marshal(wlog)
	reqinfo := bytes.NewBuffer([]byte(req))
	request, _ := http.NewRequest("POST", WafUrl+"/api/v1/waflog", reqinfo)
	_, _ = client.Do(request)
}
