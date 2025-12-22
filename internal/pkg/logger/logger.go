package logger

import (
	"fmt"

	"github.com/spksupakorn/go-restful-authentication/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger based on configuration
func NewLogger(cfg *config.Config) (*zap.Logger, error) {
	var zapConfig zap.Config

	// Determine log level
	level := zapcore.InfoLevel
	switch cfg.Logging.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure based on format and environment
	if cfg.Logging.Format == "json" || cfg.Server.Env == "production" {
		// Production config with JSON format
		zapConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(level),
			Development: false,
			Encoding:    "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "timestamp",
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	} else {
		// Development config with console format
		zapConfig = zap.Config{
			Level:       zap.NewAtomicLevelAt(level),
			Development: true,
			Encoding:    "console",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "T",
				LevelKey:       "L",
				NameKey:        "N",
				CallerKey:      "C",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "M",
				StacktraceKey:  "S",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.CapitalColorLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
	}

	logger, err := zapConfig.Build(zap.AddCallerSkip(0))
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	// Add common fields
	logger = logger.With(
		zap.String("service", "user-management-api"),
		zap.String("environment", cfg.Server.Env),
	)

	return logger, nil
}
