package main

import (
	"github.com/go-redis/redis"
)

var Rcli *redis.Client

func RedisInit() {
	Rcli = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1",
		Password: "",
		DB:       0,
	})
	_, err := Rcli.Ping().Result()
	if err == nil {
	} else {
		panic("redis连接失败")
	}
}
