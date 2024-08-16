package main

import (
	"context"
	"fmt"
	"go_gin_advanced/gin_02/dao/mysql"
	"go_gin_advanced/gin_02/logger"
	"go_gin_advanced/gin_02/routes"
	"go_gin_advanced/gin_02/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// web_staging 脚手架

func main() {

	// 1. 加载配置
	if err := settings.ConfigInit(); err != nil {
		fmt.Printf("found err in settings.ConfigInit : %v\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.LogInit(); err != nil {
		fmt.Printf("found err in logger.Loginit : %v\n", err)
		return
	}

	defer zap.L().Sync() // 缓冲区的日志刷到磁盘中
	// 3. 初始化mysql连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("found err in mysql.MysqlInit : %v\n", err)
		return
	}

	defer mysql.Close()

	/*
		// redis 这里有点问题 : 连接不上
		// 4. 初始化redis连接
		if err := redis.Init(); err != nil {
			fmt.Printf("found err in redis.RedisInit : %v\n", err)
			return
		}

		defer redis.Close()

	*/

	// 5. 注册路由
	r := routes.SetUp()

	// 6. 优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("staging.port")),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdowm Server ...")
	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := srv.Shutdown(context); err != nil {
		zap.L().Fatal("err in shutdown server : ", zap.Error(err))
		return
	}

	zap.L().Info("server existing ...")
}
