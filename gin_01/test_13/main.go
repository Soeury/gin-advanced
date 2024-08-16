package main

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// zap - 定制 logger
//  -1  将日志写入文件而不是在终端打印
//  -2  将日志进行切割归档 (当日志文件大小过大的时候(通常采用JSON格式)，操作起来很不方便)
//       zap本身不支持切割归档，使用第三方库 Lumberjack

var logger *zap.Logger
var sugarLogger *zap.SugaredLogger

func initLogger() {

	// core := zapcore.NewCore() 打造一个核心，通过这个核心去创建一个 logger
	// NewCore里面要传入三个参数 :
	//  -1 encoder : 编码器(如何写入日志)
	//  -2 writeSyncer : 指定日志写到哪里去
	//  -3 LogLevel : 哪种级别的日志将被写入
	encoder := getEncoder()
	writerSyncer := getLogWriter()
	level := zapcore.DebugLevel
	core := zapcore.NewCore(encoder, writerSyncer, level)

	logger = zap.New(core, zap.AddCaller())
	sugarLogger = logger.Sugar()
}

func getEncoder() zapcore.Encoder {

	// 以下三条代码将日志内部时间的格式改成我们可读的形式
	// 不要忘记了在 zap.New() 中加上 zap.AddCaller()
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
	// return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {

	/*

		// 这里 os.Create(path) 表示每次都创建一个新的文件
		// os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744) 表示每次将日志追加到文件中而不是覆盖
		// 0744：这是文件的权限设置。文件权限通常以三位数字表示，分别代表文件所有者、所属组和其他用户的权限
		// 0744表示文件所有者可以读写执行（7），所属组和其他用户可以读（4）
		path := "D:\\M_GO\\GO_gin_advanced\\gin_01\\test_13\\log.txt"
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0744)
		if err != nil {
			panic(err)
		}

		// 利用 io.MultiWriter 支持文件和终端两个输出目标
		ws := io.MultiWriter(file, os.Stdout)
		return zapcore.AddSync(ws)

	*/

	// 日志切割归档使用如下代码:
	LumberjackLogger := &lumberjack.Logger{
		Filename:   "./log.txt",
		MaxSize:    1,     // 按照大小切割
		MaxBackups: 5,     // 最大备份数量
		MaxAge:     30,    // 最大备份天数
		Compress:   false, // 是否压缩
	}
	return zapcore.AddSync(LumberjackLogger)
}

// logger - demo
func Demo(url string) {

	resp, err := http.Get(url)
	if err != nil {
		sugarLogger.Error(
			"found err in http.get...",
			zap.String("url", url),
			zap.Error(err))
	} else {
		sugarLogger.Info(
			"success...",
			zap.String("statuscode", resp.Status))
		resp.Body.Close()
	}
}

func main() {

	initLogger()
	defer sugarLogger.Sync()

	/*  测试日志切割:
	for i := 0; i < 10000; i++ {
	    sugarLogger.Info("test for rotate...")
	}
	*/

	Demo("www.google.com")
	Demo("https://www.baidu.com")
}
