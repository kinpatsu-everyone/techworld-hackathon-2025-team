package outorouter

import (
	"context"
	"net/http"
	"time"
)

// Logger はログ出力のためのインターフェースです
// さまざまなログレベルでログを出力するメソッドを提供します
type Logger interface {
	// Debug レベルのログを出力します
	Debug(ctx context.Context, msg string, keyAndValues map[string]any)
	// Info レベルのログを出力します
	Info(ctx context.Context, msg string, keyAndValues map[string]any)
	// Error レベルのログを出力します
	Error(ctx context.Context, msg string, keyAndValues map[string]any)
}

func WithLogger(l Logger) Option {
	return func(o *Router) {
		o.logger = l
	}
}

// statusRecorder はレスポンスのステータスコードとサイズを記録するための構造体です
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.statusCode == 0 {
		r.statusCode = code
		r.ResponseWriter.WriteHeader(code)
	}
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if r.statusCode == 0 {
		// Writeが呼ばれたからステータスコードが設定されていない場合は200 OKと見なす
		r.statusCode = http.StatusOK
	}
	size, err := r.ResponseWriter.Write(b)
	r.size += size
	return size, err
}

// LoggingMiddleware はHTTPリクエストとレスポンスの情報をログに出力するミドルウェアを生成します
func LoggingMiddleware(logger Logger) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			start := GetNowUTCFromContext(ctx)
			requestID := GetRequestIDFromContext(ctx)

			recorder := &statusRecorder{
				ResponseWriter: w,
				statusCode:     0,
				size:           0,
			}

			// Panicをキャッチしてログに出力する
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error(ctx, "panic が発生しました", map[string]any{
						"request_id": requestID,
						"panic":      rec,
						"method":     r.Method,
						"path":       r.URL.Path,
					})
					// Only write error response if headers haven't been sent
					if recorder.statusCode == 0 {
						http.Error(recorder, "サーバー内部で予期しないエラーが発生しました", http.StatusInternalServerError)
					}
				}
			}()

			next.ServeHTTP(recorder, r)

			latency := time.Since(start)
			logger.Info(ctx, "HTTP リクエストが処理されました", map[string]any{
				"request_id":  requestID,
				"method":      r.Method,
				"path":        r.URL.Path,
				"status_code": recorder.statusCode,
				"size":        recorder.size,
				"duration_ms": latency.Milliseconds(),
				"remote_addr": r.RemoteAddr,
				"user_agent":  r.UserAgent(),
			})
		})
	}
}
