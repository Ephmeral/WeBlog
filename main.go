package main

import (
	"fmt"

	"gohub/bootstrap"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建一个Gin实例
	router := gin.New()

	// 初始化路由绑定
	bootstrap.SetupRoute(router)

	// 启动服务，默认为8080，这里改成9090
	err := router.Run(":9090")
	if err != nil {
		// 错误处理，端口被占用或其他错误
		fmt.Println(err.Error())
	}
}
