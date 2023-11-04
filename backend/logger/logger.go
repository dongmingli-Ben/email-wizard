package logger

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)


var logger *zap.Logger

var LOGGING_LEVEL_MAP = map[string]zapcore.Level{
	"DEBUG": zap.DebugLevel,
	"INFO": zap.InfoLevel,
	"WARN": zap.WarnLevel,
	"ERROR": zap.ErrorLevel,
	"FATAL": zap.FatalLevel,
}

func InitLogger(log_dir string, name string, rotation_period int, backup int, level string) {
	encoder_config := zap.NewProductionEncoderConfig()
	encoder_config.EncodeTime = zapcore.ISO8601TimeEncoder
	stacktrace_encoder := zapcore.NewJSONEncoder(encoder_config)

	if _, ok := LOGGING_LEVEL_MAP[level]; !ok {
		fmt.Printf("log level %v not recognized, setting level to INFO ...", level)
	}
	core := zapcore.NewCore(
		stacktrace_encoder,
		zapcore.AddSync(&lumberjack.Logger{
			Filename: filepath.Join(log_dir, fmt.Sprintf("%v.log", name)),
			MaxBackups: backup,
			MaxAge: rotation_period,
			Compress: false,
		}),
		LOGGING_LEVEL_MAP[level],
	)
	logger = zap.New(core)
}

func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zapcore.Field) {
	logger.Fatal(msg, fields...)
}

func LogErrorStackTrace() {
	if r := recover(); r != nil {
		logger.Error("Panic occurred", zap.Any("panic", r), zap.Stack("stack"))
		panic(r)
	}
}

func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
		defer LogErrorStackTrace()
        start := time.Now()
		
		// Create a custom io.ReadCloser to capture and buffer the request body
        var buf bytes.Buffer
        tee := io.TeeReader(c.Request.Body, &buf)
        c.Request.Body = io.NopCloser(tee)
        c.Next()

        end := time.Now()
        latency := end.Sub(start)

        logger.Info("HTTP request",
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Int("status", c.Writer.Status()),
            zap.Any("params", c.Request.URL.Query()),
			zap.String("payload", buf.String()),
            zap.Duration("latency", latency),
            zap.String("userAgent", c.GetHeader("User-Agent")),
        )
    }
}