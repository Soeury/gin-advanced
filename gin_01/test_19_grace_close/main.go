package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

// 优雅的关机和重启
// 1. 优雅的关机 : 服务器关闭后将正在处理的请求处理完毕后再退出程序
//    直接关闭服务器端会强制结束进程导致正在处理的请求出现问题

// 这里使用的是 context 包结合 http.server 和 Shutdown 方法
// 类似处理'优雅的关机'的方法还有很多，gin.engine 和 第三方库 graceful 和 grace ...

func main() {

	// 创建自定义日志文件
	path := "./server.log"
	logFile, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0744)
	if err != nil {
		log.Fatalf("found err in os.openfile : %s", err)
	}
	defer logFile.Close()

	//  log.New(输入目标文件 ， 前缀 ， 日志格式)
	logger := log.New(logFile, "shutdown: ", log.LstdFlags)
	logger.Println("this log entry will be written to server.log")

	// 注册路由
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		time.Sleep(time.Second * 10)
		c.String(http.StatusOK, "welcome to server")
	})

	// 创建一个 http.server 实例
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	// 运行一个长时间的任务会阻塞主goroutine , 所以这里通过开启一个goroutine来启动http服务器
	// 再这个 goroutine 里面，ListenAndServe 会阻塞，直到关闭或遇到错误
	// log.Fatalf("msg") 打印消息后会立即终止程序，使用前提是找到了错误
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待终端信号来优雅的关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道

	// kill 默认发送一个 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，像常用的 ctrl + c 就是触发系统的 SIGINT 信号
	// kill -9 发送 syscall.SIGKILL 信号， 但是不能被捕获，所以不用添加

	// signal.Notify 把收到的信号转发给 quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 阻塞在这里，只有通道收到了信号之后(才能把信号拿出来)才会往下执行
	log.Println("Shutdowm Server ...")

	// 创建一个5秒超时的context
	// defer cancel() 应该是在5秒内没有处理完请求，然后取消所有的操作，释放资源
	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 5秒内关闭服务(将未处理完的请求处理完成后再关闭) ， 超过5秒就退出
	if err := srv.Shutdown(context); err != nil {
		logger.Println("server shutdown : ", err)
	}

	logger.Println("Server existing ...")
}
