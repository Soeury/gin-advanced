package main

import (
	"fmt"
	"go_gin_advanced/gin_01/gin_zap/config"
	"go_gin_advanced/gin_01/gin_zap/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 处理main函数
func main() {

	// 初始化config配置
	err := config.Init()
	if err != nil {
		fmt.Println("panic in config.init")
		panic(err)
	}

	// 初始化 logger
	err = logger.InitLogger(config.Conf.LogConfig)
	if err != nil {
		fmt.Printf("found err in initLogger : %v\n", err)
		return
	}

	// 设置gin框架的运行模式  常见的有 debug , release (production) ...
	gin.SetMode(config.Conf.Mode)

	// 配置路由
	r := gin.Default()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/zap", func(c *gin.Context) {
		// 数据记录到日志
		name := "zap_in_gin"
		msg := 200

		zap.L().Info(
			"this is zap func",
			zap.String("user", name),
			zap.Int("msg", msg),
		)

		c.String(http.StatusOK, "hello zap")
	})

	addr := fmt.Sprintf(":%v", config.Conf.Port)
	r.Run(addr)
}
