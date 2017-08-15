package context

import (
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is a interface to log object
type Logger interface {
	// Debug write a debug log level
	Debug(action, result, message string)

	// Info write a info log level
	Info(action, result, message string)

	// Warn write a warning log level
	Warn(action, result, message string)

	// Error write a error log level
	Error(action, result, message string)
}

// APILog is the API logger
type APILog struct {
	logger *zap.Logger
}

// NewAPILog returns a pointer of the APILog
func NewAPILog() *APILog {
	loggerInstance := getNewLogInstance()
	return &APILog{
		logger: loggerInstance,
	}
}

// Debug write a debug log level
func (l APILog) Debug(action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(action, result)
	l.logger.Debug(message, fields...)
}

// Info write a info log level
func (l APILog) Info(action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(action, result)
	l.logger.Info(message, fields...)
}

// Warn write a warning log level
func (l APILog) Warn(action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(action, result)
	l.logger.Warn(message, fields...)
}

// Error write a error log level
func (l APILog) Error(action, result, message string) {
	defer l.logger.Sync()

	fields := getLogFields(action, result)
	l.logger.Error(message, fields...)
}

// getLogFields return a logrus fields instance
func getLogFields(action, result string) []zapcore.Field {
	return []zapcore.Field{
		zap.String("action", action),
		zap.String("result", result),
	}
}

var apiLogger *APILog
var onceLog sync.Once

// GetLogger return a new instance of the log
func GetLogger() Logger {
	onceLog.Do(func() {
		apiLogger = NewAPILog()
	})

	return apiLogger
}

// getNewLogInstance return a new instance of Logrus log
func getNewLogInstance() *zap.Logger {
	logLevel := getLogLevel()

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(logLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		// EncoderConfig: zap.NewProductionEncoderConfig(),
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger
}

func getLogLevel() zapcore.Level {
	apiConfig := GetAPIConfig()
	logLevelConfig := strings.ToLower(apiConfig.LogConfig.LogLevel)

	if logLevelConfig == "debug" {
		return zap.DebugLevel
	}

	if logLevelConfig == "info" {
		return zap.InfoLevel
	}

	if logLevelConfig == "warn" {
		return zap.WarnLevel
	}

	return zap.ErrorLevel
}
