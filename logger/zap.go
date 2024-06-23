package logger

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger is an implementation of the Logger interface using the zap logging library.
type ZapLogger struct {
	*zap.Logger
}

// Debug logs a message at the debug level.
func (z ZapLogger) Debug(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		zFields = append(zFields, field.(zapcore.Field))
	}
	z.Logger.Debug(msg, zFields...)
}

// Info logs a message at the info level.
func (z ZapLogger) Info(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		zFields = append(zFields, field.(zapcore.Field))
	}
	z.Logger.Info(msg, zFields...)
}

// Warn logs a message at the warn level.
func (z ZapLogger) Warn(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		zFields = append(zFields, field.(zapcore.Field))
	}
	z.Logger.Warn(msg, zFields...)
}

// Error logs a message at the error level.
func (z ZapLogger) Error(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		zFields = append(zFields, field.(zapcore.Field))
	}
	z.Logger.Error(msg, zFields...)
}

// Fatal logs a message at the fatal level and then calls os.Exit(1).
func (z ZapLogger) Fatal(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		zFields = append(zFields, field.(zapcore.Field))
	}
	z.Logger.Fatal(msg, zFields...)
}

// Panic logs a message at the panic level and then panics.
func (z ZapLogger) Panic(msg string, fields ...interface{}) {
	var zFields []zapcore.Field
	for _, field := range fields {
		zFields = append(zFields, field.(zapcore.Field))
	}
	z.Logger.Panic(msg, zFields...)
}

// GetZapProductionLogger creates and returns a new zap.Logger configured for production use.
// The production logger is optimized for performance. It uses a JSON encoder, logs to standard
// error, and writes at InfoLevel and above.
//
// Parameters:
//   - level: The minimum logging level at which logs should be written,
//     e.g., zapcore.DebugLevel, zapcore.InfoLevel.
//
// Returns:
//   - Logger: The configured zap.Logger for production use.
//   - error: An error if the logger could not be created.
func GetZapProductionLogger(level zap.AtomicLevel) (Logger, error) {
	config := zap.NewProductionConfig()
	config.Level = level
	logger, err := config.Build()
	return ZapLogger{logger}, err
}

// GetZapDevelopmentLogger creates and returns a new zap.Logger configured for development use.
// The development logger is more verbose and is intended for use during development. It uses
// a console encoder with colored level output and logs at the specified log level.
//
// Parameters:
//   - level: The minimum logging level at which logs should be written,
//     e.g., zapcore.DebugLevel, zapcore.InfoLevel.
//
// Returns:
//   - Logger: The configured zap.Logger for development use.
//   - error: An error if the logger could not be created.
func GetZapDevelopmentLogger(level zap.AtomicLevel) (Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.Level = level
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, err := config.Build()
	return ZapLogger{logger}, err
}

// GetZapLogger creates and returns a new zap.Logger based on the environment and log level.
// It supports "development" and "production" environments.
//
// Parameters:
//   - env: The environment for the logger, e.g., "development", "production".
//   - level: The logging level, e.g., "debug", "info".
//
// Returns:
//   - Logger: The configured zap.Logger based on the environment and log level.
//   - error: An error if the logger could not be created.
func GetZapLogger(env string, level string) (Logger, error) {
	configLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, errors.Wrapf(err,
			"invalid logger level provided: %s - err: %s",
			level,
		)
	}

	switch env {
	case "development":
		return GetZapDevelopmentLogger(configLevel)
	case "production":
		return GetZapProductionLogger(configLevel)
	default:
		return nil, fmt.Errorf("failure to construct logger for env: %s", env)
	}
}
