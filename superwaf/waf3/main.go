package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

var (
	Dbinfo  string
	Logpath string
)

func init() {
	flag.StringVar(&Dbinfo, "mysql", "", "mysql连接信息")
	flag.StringVar(&Logpath, "log", "", "日志文件路径")
}

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.String(200, "hello world")
	})
	r.Run()
}
