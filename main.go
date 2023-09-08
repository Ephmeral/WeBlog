package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个Gin实例
	r := gin.New()

	// 注册中间件
	r.Use(gin.Logger(), gin.Recovery())

	// 注册一个路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Hello": "World!",
		})
	})

	// 处理404请求
	r.NoRoute(func(c *gin.Context) {
		// 获取标头信息的Accept信息
		acceptString := c.Request.Header.Get("Accept")
		if strings.Contains(acceptString, "text/html") {
			// 如果是HTML的话
			c.String(http.StatusNotFound, "页面返回404")
		} else {
			// 默认返回JSON
			c.JSON(http.StatusNotFound, gin.H{
				"error_code":    404,
				"error_message": "路由未定义，请确认url和请求方法是否正确",
			})
		}
	})

	// 启动服务，默认为8080，这里改成9090
	r.Run(":9090")
}
