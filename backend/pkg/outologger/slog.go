package outologger

import (
	"context"
	"log/slog"
)

type SlogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(logger *slog.Logger) *SlogLogger {
	return &SlogLogger{logger: logger}
}

func (l *SlogLogger) Debug(ctx context.Context, msg string, keyAndValues map[string]any) {
	l.logger.LogAttrs(
		ctx,
		slog.LevelDebug,
		msg,
		mapToSlogAttrs(keyAndValues)...,
	)
}

func (l *SlogLogger) Info(ctx context.Context, msg string, keyAndValues map[string]any) {
	l.logger.LogAttrs(
		ctx,
		slog.LevelInfo,
		msg,
		mapToSlogAttrs(keyAndValues)...,
	)
}

func (l *SlogLogger) Warn(ctx context.Context, msg string, keyAndValues map[string]any) {
	l.logger.LogAttrs(
		ctx,
		slog.LevelWarn,
		msg,
		mapToSlogAttrs(keyAndValues)...,
	)
}

func (l *SlogLogger) Error(ctx context.Context, msg string, keyAndValues map[string]any) {
	l.logger.LogAttrs(
		ctx,
		slog.LevelError,
		msg,
		mapToSlogAttrs(keyAndValues)...,
	)
}

// mapToSlogAttrs は map[string]any → []slog.Attr に変換
func mapToSlogAttrs(m map[string]any) []slog.Attr {
	if len(m) == 0 {
		return nil
	}
	attrs := make([]slog.Attr, 0, len(m))
	for k, v := range m {
		attrs = append(attrs, slog.Any(k, v))
	}
	return attrs
}
