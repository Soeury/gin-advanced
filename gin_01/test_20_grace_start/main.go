package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

// 优雅的重启 : 被用于需要高可用性、频繁更新且对中断敏感的生产环境中
// 这里使用 fvbock/endless 来替换默认的 listenandserver 启动服务来实现
// endless 允许服务器在不停止处理当前连接的情况下进行代码或配置的更新

func main() {

	// 创建自定义日志文件:
	path := "./server.log"
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)
	if err != nil {
		log.Fatalf("found err in os.openfile : %s", err)
	}
	defer logFile.Close()

	//  log.New(输入目标文件 ， 前缀 ， 日志格式)
	logger := log.New(logFile, "shutdown: ", log.LstdFlags)
	logger.Println("this log entry will be written to server.log")

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(time.Second * 10)
		c.String(http.StatusOK, "welcome to graceful start ...")
	})

	// 默认 endless 服务器会监听以下信号:
	// syscall.SIGHUP , syscall.sigusr1 , syscall.SIGUSR2 , syscall.SIGINT , syscall.SIGTERM , syscall.SIGTSTP
	// 接收到  syscall.SIGHUP 信号将触发 'fork/restart' 实现优雅重启
	// 接收到  syscall.SIGINT , syscall.SIGTERM 信号将触发优雅关机

	// 不停止处理当前连接的情况下重启服务器 (fork 了一个子进程)
	if err := endless.ListenAndServe(":8080", r); err != nil {
		logger.Fatalf("listen: %s\n", err)
	}

	logger.Println("server existing ...")
}
