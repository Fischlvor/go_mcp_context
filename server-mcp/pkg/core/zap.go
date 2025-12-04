package core

import (
	"io"
	"log"
	"os"

	"go-mcp-context/pkg/global"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化并返回一个基于配置设置的新 zap.Logger 实例
func InitLogger() *zap.Logger {
	zapCfg := global.Config.Zap

	// 创建 lumberjack 日志文件
	lumberJackLogger := &lumberjack.Logger{
		Filename:   zapCfg.Filename,
		MaxSize:    zapCfg.MaxSize,
		MaxBackups: zapCfg.MaxBackups,
		MaxAge:     zapCfg.MaxAge,
	}

	// 创建一个用于日志输出的 writeSyncer
	writeSyncer := zapcore.AddSync(lumberJackLogger)

	// 如果配置了控制台输出，则添加控制台输出
	if zapCfg.IsConsolePrint {
		writeSyncer = zapcore.NewMultiWriteSyncer(writeSyncer, zapcore.AddSync(os.Stdout))
		// 同时重定向标准库 log 到文件和控制台
		global.LogWriter = io.MultiWriter(lumberJackLogger, os.Stdout)
	} else {
		// 只重定向到文件
		global.LogWriter = lumberJackLogger
	}

	// 重定向标准库 log
	log.SetOutput(global.LogWriter)

	// 重定向 Gin 的默认输出
	gin.DefaultWriter = global.LogWriter
	gin.DefaultErrorWriter = global.LogWriter

	// 创建日志格式化的编码器
	encoder := getEncoder()

	// 根据配置确定日志级别
	var logLevel zapcore.Level

	if err := logLevel.UnmarshalText([]byte(zapCfg.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	// 创建核心和日志实例
	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	logger := zap.New(core, zap.AddCaller())
	return logger
}

// getEncoder 返回一个为生产日志配置的 JSON 编码器
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}
