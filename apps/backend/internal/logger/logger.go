package logger

import (
	"log/slog"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap  *zap.Logger
	slog *slog.Logger
}

func New(environment string) (*Logger, error) {
	var zapLogger *zap.Logger
	var err error

	// Configure based on environment
	if environment == "production" {
		config := zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "level"
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.StacktraceKey = "stacktrace"
		
		zapLogger, err = config.Build(
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	} else {
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		
		zapLogger, err = config.Build(
			zap.AddCallerSkip(1),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	}

	if err != nil {
		return nil, err
	}

	// Create slog logger for middleware
	var slogHandler slog.Handler
	if environment == "production" {
		slogHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
			AddSource: true,
		})
	} else {
		slogHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
			AddSource: true,
		})
	}

	slogLogger := slog.New(slogHandler)

	return &Logger{
		zap:  zapLogger,
		slog: slogLogger,
	}, nil
}

// Zap returns the underlying zap logger
func (l *Logger) Zap() *zap.Logger {
	return l.zap
}

// Slog returns the underlying slog logger for middleware
func (l *Logger) Slog() *slog.Logger {
	return l.slog
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.zap.Sync()
}

// With adds fields to the logger
func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		zap:  l.zap.With(fields...),
		slog: l.slog,
	}
}

// WithGroup creates a new logger with a group
func (l *Logger) WithGroup(group string) *Logger {
	return &Logger{
		zap:  l.zap.Named(group),
		slog: l.slog.WithGroup(group),
	}
}

// Logger methods
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.zap.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.zap.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.zap.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.zap.Error(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...zap.Field) {
	l.zap.Fatal(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...zap.Field) {
	l.zap.Panic(msg, fields...)
}

// Structured logging helpers
func (l *Logger) WithRequest(requestID string) *Logger {
	return l.With(zap.String("request_id", requestID))
}

func (l *Logger) WithUser(userID string) *Logger {
	return l.With(zap.String("user_id", userID))
}

func (l *Logger) WithTenant(tenantID string) *Logger {
	return l.With(zap.String("tenant_id", tenantID))
}

func (l *Logger) WithDuration(duration time.Duration) *Logger {
	return l.With(zap.Duration("duration", duration))
}

func (l *Logger) WithError(err error) *Logger {
	return l.With(zap.Error(err))
}