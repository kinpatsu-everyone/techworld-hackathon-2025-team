package outorouter

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKeyRequestID struct{}

func GetRequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxKeyRequestID{}).(string); ok {
		return v
	}
	return ""
}

// RequestIDMiddleware はリクエストごとに一意な RequestID を発行し、
// context とレスポンスヘッダにセットする。
func RequestIDMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = uuid.NewString()
			}

			ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, requestID)
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
