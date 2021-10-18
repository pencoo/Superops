package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	U = append(U, Student{1, "小明", 18, "男"})
	U = append(U, Student{2, "小红", 16, "女"})
	U = append(U, Student{3, "小刚", 17, "男"})
	U = append(U, Student{4, "小丽", 17, "女"})
	r.Use(Cors())
	r.GET("/", FunGet)
	r.POST("/", FunPost)
	r.GET("/d", FunDel)
	r.Run(":8080")
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		// 放行所有OPTIONS方法，因为有的模板是要请求两次的
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

var U []Student

type Student struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

func FunGet(c *gin.Context) {
	u := c.Query("name")
	if u != "" {
		r := 0
		for _, v := range U {
			if v.Name == u {
				r = 1
				c.JSON(200, v)
			}
		}
		if r == 0 {
			c.String(400, "用户不存在")
		}
	} else {
		c.JSON(200, U)
	}
}

func FunPost(c *gin.Context) {
	var u Student
	if err := c.BindJSON(&u); err == nil {
		fmt.Println(u)
		U = append(U, u)
		c.JSON(200, u)
	} else {
		if err != nil {
			c.JSON(400, err)
		} else {
			c.JSON(400, "有必须参数为空")
		}
	}
}

func FunDel(c *gin.Context) {
	u := c.Query("name")
	if u != "" {
		r := 0
		var Uu []Student
		for _, v := range U {
			if v.Name != u {
				Uu = append(Uu, v)
			} else {
				r += 1
			}
		}
		if r == 1 {
			U = Uu
			Uu = nil
			c.JSON(200, "数据删除成功")
		} else {
			c.JSON(200, "数据删除失败")
		}
	} else {
		c.JSON(400, "用户名为空")
	}
}
