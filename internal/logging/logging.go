package logging

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	DebugLevel   = "debug"
	InfoLevel    = "info"
	WarningLevel = "warning"
	ErrorLevel   = "error"
)

var AllowedLogLevels = []string{DebugLevel, InfoLevel, WarningLevel, ErrorLevel}

var Logger = zap.Must(zap.NewDevelopmentConfig().Build()).Sugar()

// InitializeLogging initializes the zap logging library with the specified log level.
func InitializeLogging(logLevel string) error {
	config := zap.NewProductionConfig()
	config.Encoding = "console"

	// Configure time encoder to use ISO8601 format
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	level, err := parseLogLevel(logLevel)
	if err != nil {
		return err // Return the error from parseLogLevel
	}
	config.Level.SetLevel(level)

	// Build the logger
	l, err := config.Build()
	if err != nil {
		return fmt.Errorf("failed to initialize zap logger: %w", err)
	}
	Logger = l.Sugar()

	return nil
}

func parseLogLevel(logLevel string) (zapcore.Level, error) {
	switch logLevel {
	case DebugLevel:
		return zapcore.DebugLevel, nil
	case InfoLevel:
		return zapcore.InfoLevel, nil
	case WarningLevel:
		return zapcore.WarnLevel, nil
	case ErrorLevel:
		return zapcore.ErrorLevel, nil
	default:
		return zapcore.InfoLevel, fmt.Errorf("invalid log level: %s. Allowed levels are %v. Defaulting to info", logLevel, AllowedLogLevels)
	}
}

func CloseLogger() {
	if Logger != nil {
		Logger.Sync() // flushes buffer, if any
	}
}
