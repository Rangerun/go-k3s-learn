package main

import (
	"os"
	"github.com/gin-gonic/gin"
)

// 标记是否故障，用于测试存活探针
var isBroken bool = false

func main() {
	r := gin.Default()

	appVer := os.Getenv("APP_VER")
	appMsg := os.Getenv("APP_MSG")
	dbPass := os.Getenv("DB_PASSWORD")

	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ver":     appVer,
			"msg":     appMsg,
			"db_pass": dbPass,
		})
	})

	// 健康检查接口
	r.GET("/health", func(c *gin.Context) {
		if isBroken {
			c.Status(500)
			return
		}
		c.Status(200)
	})

	// 制造故障的接口，访问这个接口把isBroken置为true，模拟程序卡死
	r.GET("/break", func(c *gin.Context) {
		isBroken = true
		c.JSON(200, gin.H{"info": "now service broken, health will return 500"})
	})

	r.Run(":8080")
}
