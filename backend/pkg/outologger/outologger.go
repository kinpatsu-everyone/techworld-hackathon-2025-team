package outologger

import "context"

var globalLogger Logger

// Logger はログ出力のためのインターフェースです
// さまざまなログレベルでログを出力するメソッドを提供します
type Logger interface {
	// Debug レベルのログを出力します
	Debug(ctx context.Context, msg string, keyAndValues map[string]any)
	// Info レベルのログを出力します
	Info(ctx context.Context, msg string, keyAndValues map[string]any)
	// Warn レベルのログを出力します
	Warn(ctx context.Context, msg string, keyAndValues map[string]any)
	// Error レベルのログを出力します
	Error(ctx context.Context, msg string, keyAndValues map[string]any)
}

// SetLogger はグローバルなロガーを設定します
func SetLogger(logger Logger) {
	globalLogger = logger
}

// GetLogger はグローバルなロガーを取得します
func GetLogger() Logger {
	return globalLogger
}
