package log

import (
	"context"

	"github.com/go-park-mail-ru/2025_2_Avrora/internal/delivery/http/middleware"
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}

func New(logger *zap.Logger) *Logger {
	return &Logger{Logger: logger}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

func (l *Logger) withReqID(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return l.Logger
	}
	if rid, ok := ctx.Value(middleware.RequestIDKey).(string); ok && rid != "" {
		return l.Logger.With(zap.String("request_id", rid))
	}
	return l.Logger
}

func (l *Logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	l.withReqID(ctx).Info(msg, fields...)
}

func (l *Logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	l.withReqID(ctx).Error(msg, fields...)
}

func (l *Logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	l.withReqID(ctx).Warn(msg, fields...)
}