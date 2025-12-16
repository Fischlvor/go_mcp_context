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

// InitLogger 初始化日志，使用 Gin 原生格式
func InitLogger() *zap.Logger {
	zapCfg := global.Config.Zap

	// 创建 lumberjack 日志文件（支持轮转）
	lumberJackLogger := &lumberjack.Logger{
		Filename:   zapCfg.Filename,
		MaxSize:    zapCfg.MaxSize,
		MaxBackups: zapCfg.MaxBackups,
		MaxAge:     zapCfg.MaxAge,
	}

	// 设置全局日志写入器（文件）
	global.LogWriter = lumberJackLogger

	// 重定向标准库 log 到文件
	log.SetOutput(lumberJackLogger)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Gin 日志：控制台 + 文件
	if zapCfg.IsConsolePrint {
		// 同时输出到文件和控制台
		gin.DefaultWriter = io.MultiWriter(lumberJackLogger, os.Stdout)
		gin.DefaultErrorWriter = io.MultiWriter(lumberJackLogger, os.Stderr)
		gin.SetMode(gin.DebugMode)
	} else {
		// 只输出到文件
		gin.DefaultWriter = lumberJackLogger
		gin.DefaultErrorWriter = lumberJackLogger
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 zap logger（用于结构化日志场景）
	writeSyncer := zapcore.AddSync(global.LogWriter)
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		CallerKey:    "caller",
		MessageKey:   "msg",
		LineEnding:   zapcore.DefaultLineEnding,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.TimeEncoderOfLayout("2006/01/02 15:04:05"),
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	var logLevel zapcore.Level
	if err := logLevel.UnmarshalText([]byte(zapCfg.Level)); err != nil {
		log.Fatalf("Failed to parse log level: %v", err)
	}

	core := zapcore.NewCore(encoder, writeSyncer, logLevel)
	return zap.New(core, zap.AddCaller())
}
