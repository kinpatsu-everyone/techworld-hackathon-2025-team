package outorouter

import (
	"context"
	"net/http"
	"time"
)

type ctxKeyNowUTC struct{}

// GetNowUTCFromContext はコンテキストから現在時刻（UTC）を取得します。
func GetNowUTCFromContext(ctx context.Context) time.Time {
	if v, ok := ctx.Value(ctxKeyNowUTC{}).(time.Time); ok {
		return v
	}
	return time.Time{}
}

// NowUTCMiddleware はリクエストごとに現在時刻（UTC）をコンテキストにセットするミドルウェアです。
func NowUTCMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			nowUTC := time.Now().UTC()
			ctx := context.WithValue(r.Context(), ctxKeyNowUTC{}, nowUTC)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
