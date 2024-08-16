package logger

import (
	"fmt"
	"go_gin_advanced/gin_01/gin_zap/config"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// gin 框架中使用 zap 日志库

var logger *zap.Logger

// cfg *config.LogConfig 表示调用 config 包下的 LogConfig 结构体
func InitLogger(cfg *config.LogConfig) (err error) {

	writeSyncer := getLogWriter(cfg.Filename, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge)
	encoder := getEncoder()

	//  l.UnmarshalText([]slice)指针负责将一个文本类型的字符串映射到一个相应的日志级别
	var leve = new(zapcore.Level)
	err = leve.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		fmt.Println("found err in leve.unmarshaltext")
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, leve)

	logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger) //替换全局 logger , 后续调用只需使用 zap.L() 即可
	return
}

func getEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxsize int, maxbackup int, maxage int) zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxsize,
		MaxBackups: maxbackup,
		MaxAge:     maxage,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 自定义 GinLogger()中间件
func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path      // path指的是url中在?之前/xxx/xxx 的部分
		query := c.Request.URL.RawQuery // rawQuery 指的是url中?后面的键值对组
		c.Next()

		// c.Errors.ByType(gin.ErrorTypePrivate) 返回一个 gin.Errors类型的切片
		cost := time.Since(start)
		logger.Info(
			path,
			zap.Int("status", c.Writer.Status()),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
		)
	}
}

// 自定义 GinRecovery 中间件，recover掉可能出现的panic，使用zap记录相关日志
func GinRecovery(stack bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {

			// broken pipe 在连接关闭时试图写入数据,
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						seStr := strings.ToLower(se.Error())
						str1 := "broken pipe"
						str2 := "connection reset by peer"
						if strings.Contains(seStr, str1) || strings.Contains(seStr, str2) {
							brokenPipe = true
						}
					}
				}

				// 下面一行将 http 请求转换成文本表示, false 表示不包含请求体,返回的 httpRequest 为 []byte 类型
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(
						c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)

					c.Error(err.(error))
					c.Abort()
					return
				}

				// 这里 debug.Stack() 将追栈跟踪信息记录到日志文件中，通过查看相关函数调用序列发现问题
				if stack {
					logger.Error(
						"[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error(
						"[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				// 这里表示捕获到了panic 并处理了相关错误之后,向客户端发送一个服务器内部出现故障的相应 500
				c.AbortWithStatus(http.StatusInternalServerError) // 500

			}
		}()

		// 未捕获到panic则继续执行之后的中间件函数
		c.Next()
	}
}
