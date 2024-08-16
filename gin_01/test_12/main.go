package main

import (
	"net/http"

	"go.uber.org/zap"
)

//  ----  go - zap日志库: 提供了两种日志记录器
// -1  logger
//   - 调用 zap.NewProduction()  /  zap.NewDevelopment  /  zap.Example()
//   - 上面每个函数都将创建一个 logger , 唯一的区别是记录的信息不同，所以打印出的日志格式和内容会有一些差异
//   - 通过 logger 调用一些方法  info/error/panic/warning ...

// -2  sugared logger
//   - 与 logger 非常相似, 定义一个全局变量 : var sugarLogger *zap.SugaredLogger
//   - 在初始化 logger 的时候在后面加上 sugaredlogger := logger.Sugar() 得到的就是 sugaredlogger 日志记录器
//   - 之后使用 sugaredlogger 就可以了

// 定义一个全局变量
var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

// 初始化一个日志记录器
// 要更改打印的格式可以更换 zap.xxxxx()
func initLogger() {

	logger, _ = zap.NewProduction()
	sugarLogger = logger.Sugar()
}

// 日志演示
// 首先创建了一个 logger , 然后使用  logger.info 或者 logger.error 等方法记录消息
// 为什么 logger 可以 .info 和 .error 呢?  func (log *logger) methodxxx(msg string , field...)
// zap背后封装了 info/error/debug/panic... 的方法, 每个方法都接收一个消息字符串和 任意数量的 field(键值对参数)
func loggerDemo(url string) {

	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"found err in http.get...",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info(
			"success...",
			zap.String("statusCode", resp.Status))
		resp.Body.Close()
	}
}

func main() {

	initLogger()
	defer sugarLogger.Sync() // 程序结束时把缓冲区中的日志刷到磁盘上
	loggerDemo("https://www.baidu.com")
	loggerDemo("www.google.com")
}
