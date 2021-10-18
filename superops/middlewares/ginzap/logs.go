package ginzap

//https://github.com/akath19/gin-zap/blob/master/gin-zap.go
import (
	"os"
	"superops/libs/config"
	"time"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Lzap *zap.Logger

func GinzapInit() {
	if config.LogPath == "" {
		Lzap, _ = zap.NewProduction()
	} else {
		writeSyncer := getLogWriter()
		encoder := getEncoder()
		core := zapcore.NewCore(encoder, writeSyncer, config.Loglevel)
		Lzap = zap.New(core)
	}
}
func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter() zapcore.WriteSyncer {
	file, _ := os.Create(config.LogPath)
	return zapcore.AddSync(file)
}

//Logging colors, unused until zap implements colored logging -> https://github.com/uber-go/zap/issues/489
var (
	green   = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	white   = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	yellow  = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	red     = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	blue    = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	magenta = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	cyan    = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	reset   = string([]byte{27, 91, 48, 109})
)

//setupLogging setups the logger to use zap
func setupLogging(duration time.Duration, zap *zap.Logger) {
	go func() {
		for range time.Tick(duration) {
			zap.Sync()
		}
	}()
}

//ErrorLogger returns a gin handler func for errors
func ErrorLogger() gin.HandlerFunc {
	return ErrorLoggerT(gin.ErrorTypeAny)
}

// ErrorLoggerT returns a gin handler middleware with the given
// type gin.ErrorType.
func ErrorLoggerT(t gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.Writer.Written() {
			json := c.Errors.ByType(t).JSON()
			if json != nil {
				c.JSON(-1, json)
			}
		}
	}
}

//Logger returns a gin handler func for all logging
func Logger(duration time.Duration, logger *zap.Logger) gin.HandlerFunc {
	setupLogging(duration, logger)

	return func(c *gin.Context) {
		t := time.Now()

		c.Next()

		latency := time.Since(t)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		//statusColor := colorForStatus(statusCode)
		//methodColor := colorForMethod(method)
		path := c.Request.URL.Path

		switch {
		case statusCode >= 400 && statusCode <= 499:
			{
				logger.Warn("[GIN]",
					//zap.String("statusColor", statusColor),
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency.String()),
					zap.String("clientIP", clientIP),
					//zap.String("methodColor", methodColor),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
				)
			}
		case statusCode >= 500:
			{
				logger.Error("[GIN]",
					//zap.String("statusColor", statusColor),
					zap.Int("statusCode", statusCode),
					zap.String("latency", latency.String()),
					zap.String("clientIP", clientIP),
					//zap.String("methodColor", methodColor),
					zap.String("method", method),
					zap.String("path", path),
					zap.String("error", c.Errors.String()),
				)
			}
		default:
			logger.Info("[GIN]",
				//zap.String("statusColor", statusColor),
				zap.Int("statusCode", statusCode),
				zap.String("latency", latency.String()),
				zap.String("clientIP", clientIP),
				//zap.String("methodColor", methodColor),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("error", c.Errors.String()),
			)
		}
	}
}

//coorForStatus returns a color based on the status code of the response
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code <= 299:
		return green
	case code >= 300 && code <= 399:
		return white
	case code >= 400 && code <= 499:
		return yellow
	default:
		return red
	}
}

//colorForMethod returns a color based on the HTTP method of the request
func colorForMethod(method string) string {
	switch {
	case method == "GET":
		return blue
	case method == "POST":
		return cyan
	case method == "PUT":
		return yellow
	case method == "DELETE":
		return red
	case method == "PATCH":
		return green
	case method == "HEAD":
		return magenta
	case method == "OPTIONS":
		return white
	default:
		return reset
	}
}

//GormLogger struct
type GormLogger struct{}

//Print - Log Formatter
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		Lzap.Debug(
			"sql",
			zap.String("module", "gorm"),
			zap.String("type", "sql"),
			zap.Any("src", v[1]),
			zap.Any("duration", v[2]),
			zap.Any("sql", v[3]),
			zap.Any("values", v[4]),
			zap.Any("rows_returned", v[5]),
		)
	case "log":
		Lzap.Debug("log", zap.Any("gorm", v[2]))
	}
}

//日志操作封装
type ErrorInfo struct {
	Info string      //错误描述信息
	User string      //操作用户
	Path string      //url路径
	Err  error       //错误信息
	Data interface{} //操作原始数据
}

func Errlog(e ErrorInfo) {
	Lzap.Error(e.Info, zap.String("user", e.User), zap.String("path", e.Path), zap.Reflect("data", e.Data), zap.Error(e.Err))
}

type SuccessInfo struct {
	Info string      //描述信息
	User string      //操作用户
	Path string      //url路径
	Data interface{} //操作原始数据
}

func Infolog(s SuccessInfo) {
	Lzap.Info(s.Info, zap.String("user", s.User), zap.String("path", s.Path), zap.Reflect("data", s.Data))
}

func Warnlog(s SuccessInfo) {
	Lzap.Warn(s.Info, zap.String("user", s.User), zap.String("path", s.Path), zap.Reflect("data", s.Data))
}
