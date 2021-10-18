package rcache

import (
	"encoding/base64"
	"fmt"
	"superops/libs/config"
	"superops/middlewares/ginzap"
	"time"

	"github.com/go-redis/redis"
)

var Rcli *redis.Client

const (
	//缓存需要有缓存刷新接口，刷新前先删除原有缓存信息，在重新创建缓存
	//project相关键
	Project_enable_list      string = "Project_enable_list"       //存储激活中的项目列表
	Project_disable_list     string = "project_disable_list"      //存储被禁用的项目列表
	Project_app_list_enable  string = "project_app_list_enable_"  //项目激活应用列表
	Project_app_list_disable string = "project_app_list_disable_" //项目禁用应用列表
	Project_user_list_dev    string = "project_user_list_dev_"    //项目开发用户列表
	Project_user_list_test   string = "project_user_list_test_"   //项目测试用户列表
	Project_user_list_ops    string = "project_user_list_ops_"    //项目运维用户列表
	Project_env_list_dev     string = "project_env_list_dev_"     //项目开发环境列表
	Project_env_list_test    string = "project_env_list_test_"    //项目测试环境列表
	Project_env_list_ops     string = "project_env_list_ops_"     //项目运维环境列表

	//harbor相关
	Harbor_enable_list  string = "harbor_enable_list"   //查询激活harbor列表
	Harbor_disable_list string = "harbor_disable_list"  //查询禁用harbor列表
	Harbor_project_list string = "harbor_project_list_" //查询harbor项目列表

	//jenkins相关
	Jenkins_enable_list  string = "jenkins_enable_list"  //查询激活harbor列表
	Jenkins_disable_list string = "jenkins_disable_list" //查询禁用harbor列表

	//用户相关
	User_enable_list   string = "user_enable_list"    //激活状态用户列表
	User_disable_list  string = "user_disable_list"   //未激活状态用户列表
	Userinfo_token     string = "userinfo_token_"     //用户token信息
	Userinfo_api_token string = "userinfo_api_token_" //用户api token信息

	//系统类键
	Superops_list_endpoint string = "superops_list_endpoint" //superops节点集合,客户端启动后注册自己
)

func RedisInit() {
	if config.Redisurl == "" {
		ginzap.Lzap.Info("未启用缓存")
	} else {
		Rcli = redis.NewClient(&redis.Options{
			Addr:     config.Redisurl,
			Password: config.Redispass,
			DB:       config.RedisDb,
		})
		_, err := Rcli.Ping().Result()
		if err == nil {
			ginzap.Lzap.Info("redis连接成功")
		} else {
			ginzap.Lzap.Error("redis连接失败，程序退出...")
			panic("redis连接失败")
		}
	}
}

//公共方法操作
//设置过期
func RedisKeyExpire(k string, d time.Duration) error {
	return Rcli.Expire(k, d).Err()
}

//设置过期时间
func RedisKeyExpireAt(k string, d time.Time) error {
	return Rcli.ExpireAt(k, d).Err()
}

//取消过期设置
func RedisKeyNotExpire(k string) error {
	return Rcli.Expire(k, 0).Err()
}

//判断key是否存在
func RedisKeyExists(k string) bool {
	n, err := Rcli.Exists(k).Result()
	if err == nil && n > 0 {
		return true
	} else {
		return false
	}
}

//删除缓存key
func RedisKeyDel(k string) bool {
	n, err := Rcli.Del(k).Result()
	if err == nil && n > 0 {
		return true
	} else {
		return false
	}
}

//string类型操作
//string类型 增
func RedisCacheStringAdd(k, v string) error {
	return Rcli.Set(k, v, 0).Err()
}

//string类型 删
func RedisCacheStringDel(k string) error {
	return Rcli.Del(k).Err()
}

//string类型 查
func RedisCacheStringGet(k string) (string, error) {
	v, err := Rcli.Get(k).Result()
	if err == nil {
		return v, nil
	} else {
		return v, err
	}
}

//set集合类型操作
//set集合类型 增
func RedisCacheSetAdd(k, v string) error {
	_, err := Rcli.SAdd(k, v).Result()
	return err
}

//set集合类型 删
func RedisCacheSetDel(k, v string) error {
	_, err := Rcli.SRem(k, v).Result()
	return err
}

//set集合类型 判断值存在
func RedisCacheSetIn(k, v string) bool {
	return Rcli.SIsMember(k, v).Val()
}

//set集合类型 查询所有元素
func RedisCacheSetGetAll(k string) ([]string, error) {
	return Rcli.SMembers(k).Result()
}

//set集合类型 查询元素个数
func RedisCacheGetLen(k string) (int, error) {
	l, err := Rcli.SCard(k).Result()
	if err == nil && l > 0 {
		return int(l), err
	} else {
		return 0, err
	}
}

//hash类型操作
//hash类型(map类型) 增
func RedisCacheHashAdd(key, k, v string) bool {
	b, err := Rcli.HSet(key, k, v).Result()
	if b && err == nil {
		return true
	} else {
		return false
	}
}

//hash类型(map类型) 删
func RedisCacheHashDel(key, k string) bool {
	b, err := Rcli.HDel(key, k).Result()
	if b > 0 && err == nil {
		return true
	} else {
		return false
	}
}

//hash类型(map类型) 判断元素是否存在
func RedisCacheHashExist(key, k string) bool {
	b, err := Rcli.HExists(key, k).Result()
	if b && err == nil {
		return true
	} else {
		return false
	}
}

//hash类型(map类型) 查询一个元素
func RedisCacheHashGet(key, k string) (string, error) {
	return Rcli.HGet(key, k).Result()
}

//hash类型(map类型) 查询一个元素
func RedisCacheHashGetAll(k string) (map[string]string, error) {
	return Rcli.HGetAll(k).Result()
}

//hash类型(map类型) 查询hash长度
func RedisCacheHashGetLen(k string) (int, error) {
	l, err := Rcli.HLen(k).Result()
	if err == nil && l > 0 {
		return int(l), nil
	} else {
		return 0, err
	}
}

//发布/订阅
//发布
func RedisCachePublish(ch, data string) bool {
	err := Rcli.Publish(ch, AllPassProcessByEncode(data)).Err()
	if err != nil {
		return false
	} else {
		return true
	}
}

//订阅
func RedisCacheSubscribe(ch string) {
	pubsub := Rcli.Subscribe(ch)
	_, err := pubsub.Receive()
	if err != nil {
		fmt.Print(err)
		return
	}
	c := pubsub.Channel()
	for msg := range c {
		fmt.Println(msg.Channel, msg.Payload, msg.String())
	}
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
