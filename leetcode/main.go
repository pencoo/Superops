package main

import (
	"fmt"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "47.114.168.49:33307",
		Password: "Pbl123!@#",
		DB:       0,
	})
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}
