package logger

import (
	"errors"
	"fmt"
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
)

func GinLog(msg string, level zapcore.Level) {
	// 获取原始日志器
	logger := zap.L()

	// 创建一个新的日志器，禁用调用者信息
	noCallerLogger := logger.WithOptions(zap.WithCaller(false))

	// 根据级别记录日志
	switch level {
	case zapcore.ErrorLevel:
		noCallerLogger.Error(msg)
	case zapcore.WarnLevel:
		noCallerLogger.Warn(msg)
	default:
		noCallerLogger.Info(msg)
	}
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		raw := c.Request.URL.RawPath

		// 处理请求
		c.Next()

		// 收集请求细节
		cost := time.Since(start)
		statusCode := c.Writer.Status()
		if raw == "" {
			raw = path
		}
		if query != "" {
			raw = raw + "?" + query
		}

		// 构造类似 Gin 默认日志格式的消息
		// [GIN] | 200 | 5.7823ms | 127.0.0.1 | GET "/"  | 127.0.0.1 | xx.com
		msg := fmt.Sprintf("[GIN] | %3d | %4dms |  %-15s | %-4s | %-10s",
			statusCode, cost.Milliseconds(), c.ClientIP(), c.Request.Method, raw)

		// msg := fmt.Sprintf("[GIN] | %3d | %4dms |  %-15s | %-4s %-10q  | %-20s | %s",
		// 	statusCode, cost.Milliseconds(), c.ClientIP(), c.Request.Method, raw, c.Request.Host, c.Request.Referer())

		// 根据状态码选择日志级别
		var level zapcore.Level
		switch {
		case statusCode >= 500:
			level = zapcore.ErrorLevel
		case statusCode >= 400:
			level = zapcore.WarnLevel
		default:
			level = zapcore.InfoLevel
		}

		GinLog(msg, level)

		// 如果有错误，额外记录
		if len(c.Errors) > 0 {
			GinLog(fmt.Sprintf("Gin errors: %s", c.Errors[0]), zapcore.ErrorLevel)
		}
	}
}

func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool

				// 转换为字符串检查错误信息
				errStr := fmt.Sprintf("%v", err)
				if strings.Contains(strings.ToLower(errStr), "broken pipe") ||
					strings.Contains(strings.ToLower(errStr), "connection reset by peer") {
					brokenPipe = true
				} else {
					// 尝试使用 errors.As 进行类型断言
					var opErr *net.OpError
					if errors.As(err.(error), &opErr) {
						var sysErr *os.SyscallError
						if errors.As(opErr.Err, &sysErr) {
							sysErrStr := strings.ToLower(sysErr.Error())
							if strings.Contains(sysErrStr, "broken pipe") ||
								strings.Contains(sysErrStr, "connection reset by peer") {
								brokenPipe = true
							}
						}
					}
				}

				// 安全地获取HTTP请求信息
				var httpRequest []byte
				if c.Request != nil {
					httpRequest, _ = httputil.DumpRequest(c.Request, false)
				}

				requestStr := string(httpRequest)

				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", requestStr),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(fmt.Errorf("%v", err))
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", requestStr),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", requestStr),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
