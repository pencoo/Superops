package main

import "github.com/gin-gonic/gin"

func main() {
	RedisInit()
	DBinit()
	r := gin.Default()
	r.Run()
}
