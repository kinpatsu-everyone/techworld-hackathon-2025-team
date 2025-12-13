package outorouter

import (
	"net/http"
	"strings"
)

// CORSConfig はCORSミドルウェアの設定です
type CORSConfig struct {
	// AllowedOrigins は許可するオリジンのリスト
	// "*" を含む場合はすべてのオリジンを許可
	AllowedOrigins []string
	// AllowedMethods は許可するHTTPメソッドのリスト
	AllowedMethods []string
	// AllowedHeaders は許可するヘッダーのリスト
	AllowedHeaders []string
	// AllowCredentials はクレデンシャルを許可するかどうか
	AllowCredentials bool
	// MaxAge はプリフライトリクエストのキャッシュ時間（秒）
	MaxAge int
}

// DefaultCORSConfig はデフォルトのCORS設定を返します
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           86400, // 24時間
	}
}

// CORSMiddleware はCORSを処理するミドルウェアを生成します
func CORSMiddleware(config CORSConfig) MiddlewareFunc {
	allowedOrigins := make(map[string]bool)
	allowAll := false
	for _, origin := range config.AllowedOrigins {
		if origin == "*" {
			allowAll = true
			break
		}
		allowedOrigins[origin] = true
	}

	methods := strings.Join(config.AllowedMethods, ", ")
	headers := strings.Join(config.AllowedHeaders, ", ")
	maxAge := "86400"
	if config.MaxAge > 0 {
		maxAge = strings.TrimSpace(strings.Join(strings.Fields(strings.Repeat(" ", config.MaxAge)), ""))
		// 単純に整数を文字列に変換
		maxAge = intToString(config.MaxAge)
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// オリジンが許可リストにあるか確認
			if origin != "" {
				if allowAll {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				} else if allowedOrigins[origin] {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Vary", "Origin")
				}
			}

			// プリフライトリクエスト (OPTIONS) の処理
			if r.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Methods", methods)
				w.Header().Set("Access-Control-Allow-Headers", headers)
				w.Header().Set("Access-Control-Max-Age", maxAge)
				if config.AllowCredentials {
					w.Header().Set("Access-Control-Allow-Credentials", "true")
				}
				w.WriteHeader(http.StatusNoContent)
				return
			}

			// 通常のリクエストにもCORSヘッダーを設定
			if config.AllowCredentials && origin != "" && !allowAll {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			next.ServeHTTP(w, r)
		})
	}
}

func intToString(n int) string {
	if n == 0 {
		return "0"
	}
	if n < 0 {
		return "-" + intToString(-n)
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}
