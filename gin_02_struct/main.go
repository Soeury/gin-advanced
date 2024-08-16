package main

import (
	"context"
	"fmt"
	"go_gin_advanced/gin_02_struct/dao/mysql"
	"go_gin_advanced/gin_02_struct/logger"
	"go_gin_advanced/gin_02_struct/routes"
	"go_gin_advanced/gin_02_struct/settings"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// web_staging 脚手架 - 通过结构体保存配置文件中的信息
//    当其他人接管了我们的项目时，可能不知道我们的配置文件中有哪些内容，
//    这时候通过结构体保存配置文件中的信息就显得非常友好

func main() {

	// 1. 加载配置
	if err := settings.ConfigInit(); err != nil {
		fmt.Printf("found err in settings.ConfigInit : %v\n", err)
		return
	}
	// 2. 初始化日志
	if err := logger.LogInit(settings.Conf.LogConfig); err != nil {
		fmt.Printf("found err in logger.Loginit : %v\n", err)
		return
	}

	defer zap.L().Sync() // 缓冲区的日志刷到磁盘中
	// 3. 初始化mysql连接
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("found err in mysql.MysqlInit : %v\n", err)
		return
	}

	defer mysql.Close()

	/*
		// redis 这里有点问题 : 连接不上
		// 4. 初始化redis连接
		if err := redis.Init(settings.Conf.RedisConfig); err != nil {
			fmt.Printf("found err in redis.RedisInit : %v\n", err)
			return
		}

		defer redis.Close()

	*/

	// 5. 注册路由
	r := routes.SetUp()

	// 6. 优雅关机
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.StagingConfig.Port),
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
