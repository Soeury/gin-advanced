package logger

import (
	"fmt"
	"go_gin_advanced/gin_02_struct/settings"
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

func LogInit(cfg *settings.LogConfig) (err error) {

	filename := cfg.Filename
	maxsize := cfg.Max_size
	maxbackups := cfg.Max_backups
	maxage := cfg.Max_age
	writeSyncer := GetLogSyncer(filename, maxsize, maxbackups, maxage)
	encoder := GetEncoder()

	leve := new(zapcore.Level)
	err = leve.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		fmt.Printf("found err in leve.unmarshal : %v\n", err)
		return err
	}
	core := zapcore.NewCore(encoder, writeSyncer, leve)

	log := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(log) // 全局替换之后使用zap.L().method() 即可写入日志
	return err
}

func GetEncoder() zapcore.Encoder {

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func GetLogSyncer(filename string, maxsize int, maxbackups int, maxage int) zapcore.WriteSyncer {

	lumberjackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxsize,
		MaxBackups: maxbackups,
		MaxAge:     maxage,
	}
	return zapcore.AddSync(lumberjackLogger)
}

func GinLogger() gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()
		cost := time.Since(start)

		// 记录日志
		zap.L().Info(
			path,
			zap.Int("status", c.Writer.Status()),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost), // duration 表示以易于阅读的方式记录时间
		)
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {

	return func(c *gin.Context) {
		defer func() {

			// broken : 连接关闭时试图写入数据
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

				// 将请求转换为 []byte 类型
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(
						c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)

					c.Error(err.(error))
					c.Abort()
					return
				}

				// stack(bool) 表示日志是否记录追栈跟踪的信息，
				if stack {
					zap.L().Error(
						"[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error(
						"[recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				// 这里表示捕获到了panic 然后像客户端发送一个服务器出现问题的 状态码 500
				c.AbortWithStatus(http.StatusInternalServerError) // 500
			}
		}()

		// 未捕获到panic 则继续执行之后的中间件
		c.Next()
	}
}
