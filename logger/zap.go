package logger

import "fmt"

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	*zap.Logger
}

func (z ZapLogger) Debug(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		fields = append(fields, field.(zapcore.Field))
	}
	z.Logger.Debug(msg, zFields...)
}

func (z ZapLogger) Info(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		fields = append(fields, field.(zapcore.Field))
	}
	z.Logger.Info(msg, zFields...)
}

func (z ZapLogger) Warn(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		fields = append(fields, field.(zapcore.Field))
	}
	z.Logger.Warn(msg, zFields...)
}

func (z ZapLogger) Error(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		fields = append(fields, field.(zapcore.Field))
	}
	z.Logger.Error(msg, zFields...)
}

func (z ZapLogger) Fatal(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		fields = append(fields, field.(zapcore.Field))
	}
	z.Logger.Fatal(msg, zFields...)
}

func (z ZapLogger) Panic(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		fields = append(fields, field.(zapcore.Field))
	}
	z.Logger.Panic(msg, zFields...)
}

// GetProductionLogger creates and returns a new zap.Logger configured for production use.
// The production logger is optimized for performance. It uses a JSON encoder, logs to standard
// error, and writes at InfoLevel and above.
//
// Returns:
//
//	*zap.Logger - The configured zap.Logger for production use.
//	error       - An error if the logger could not be created.
func GetProductionLogger(level zap.AtomicLevel) (Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = level
	logger, err := config.Build()
	return ZapLogger{logger}, err
}

// GetDevelopmentLogger creates and returns a new zap.Logger configured for development use.
// The development logger is more verbose and is intended for use during development. It uses
// a console encoder with colored level output and logs at the specified log level.
//
// Parameters:
//
//	level - The minimum logging level at which logs should be written,
//	        e.g., zapcore.DebugLevel, zapcore.InfoLevel.
//
// Returns:
//
//	*zap.Logger - The configured zap.Logger for development use.
//	error       - An error if the logger could not be created.
func GetDevelopmentLogger(level zap.AtomicLevel) (Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.Level = level
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	return ZapLogger{logger}, err
}

func GetLogger(env string, level string) (Logger, error) {
	configLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid logger level provided: %s - err: %s",
			level, err,
		)
	}

	switch env {
	case "development":
		return GetDevelopmentLogger(configLevel)
	case "production":
		return GetProductionLogger(configLevel)
	default:
		return nil, fmt.Errorf("failure to construct logger for env: %s", env)
	}
}
