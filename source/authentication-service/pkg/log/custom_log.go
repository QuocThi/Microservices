package log

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

func newLoggerConfig(serviceName string) zapcore.EncoderConfig {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	if os.Getenv("LOG_ENV") == "prod" {
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "Timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.MessageKey = serviceName
		encoderConfig.CallerKey = "Caller"
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	return encoderConfig
}

// Log with normal info, no addition info like log name.
func NewCustomLog(serviceName string) *zap.Logger {
	var file zapcore.WriteSyncer
	var core zapcore.Core

	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	encoderConfig := newLoggerConfig(serviceName)

	if os.Getenv("LOG_ENV") == "prod" {
		file = zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./app.log",
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     7, // days
		})

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		core = zapcore.NewCore(fileEncoder, file, level)
	} else {
		stdout := zapcore.AddSync(os.Stdout)
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewCore(consoleEncoder, stdout, level)
	}

	// To combine production and development Core into one Core use NewTee.
	// https://betterstack.com/community/guides/logging/go/zap/
	// core := zapcore.NewTee(
	// 	zapcore.NewCore(consoleEncoder, stdout, level),
	// 	zapcore.NewCore(fileEncoder, file, level),
	// )

	return zap.New(core, zap.AddCaller())
}

func NewCustomLogger() logr.Logger {
	var file zapcore.WriteSyncer
	var core zapcore.Core
	logConfig := zap.NewDevelopmentConfig()
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	if os.Getenv("LOG_ENV") == "prod" {
		logConfig := zap.NewProductionConfig()
		file = zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./app.log",
			MaxSize:    10, // megabytes
			MaxBackups: 3,
			MaxAge:     7, // days
		})

		fileEncoder := zapcore.NewJSONEncoder(logConfig.EncoderConfig)

		core = zapcore.NewCore(fileEncoder, file, level)
	} else {
		stdout := zapcore.AddSync(os.Stdout)
		consoleEncoder := zapcore.NewConsoleEncoder(logConfig.EncoderConfig)
		core = zapcore.NewCore(consoleEncoder, stdout, level)
	}

	logger := zap.New(core, zap.AddCaller())
	log := zapr.NewLogger(logger)
	log = log.WithName("AuthenticationSerivce")

	return log
}
