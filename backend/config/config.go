package config

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ENV = "local"
	// MySQLUser はMySQLのユーザー名です
	MySQLUser = ""
	// MySQLPassword はMySQLのパスワードです
	MySQLPassword = ""
	// MySQLDatabase はMySQLのデータベース名です
	MySQLDatabase = ""
	// MySQLHost はMySQLのホスト名です
	MySQLHost = ""
	// MySQLPort はMySQLのポート番号です
	MySQLPort = ""

	// ApiPort はAPIのポート番号です
	ApiPort = ""

	// CORSAllowedOrigins はCORSで許可するオリジンのリスト（カンマ区切り）
	CORSAllowedOrigins = []string{}

	RedisConfig = CacheConfig{}

	// GeminiAPIKey はGoogle Gemini APIの認証キーです
	GeminiAPIKey = ""

	// GeminiBaseURL はGoogle Gemini APIのベースURLです
	GeminiBaseURL = ""
)

type CacheConfig struct {
	Addr        string
	Username    string
	Password    string
	DB          int
	TLSEnabled  bool
	TLSInsecure bool
	KeyPrefix   string
	DefaultTTL  time.Duration
}

func loadEnv(ctx context.Context, key string, isSecret bool) string {
	result := os.Getenv(key)
	// TODO: isSecretがtrueの場合はSecretManagerから取得するようにする

	if result == "" {
		panic("config: required environment variable " + key + " is missing. Please set " + key + " in your environment.")
	}
	return result
}

// LoadEnv loads environment variables into global config variables.
// It panics if any required configuration is missing or invalid.
func LoadEnv(ctx context.Context) {
	ENV = loadEnv(ctx, "ENV", false)
	MySQLUser = loadEnv(ctx, "MYSQL_USER", true)
	MySQLPassword = loadEnv(ctx, "MYSQL_PASSWORD", true)
	MySQLDatabase = loadEnv(ctx, "MYSQL_DATABASE", true)
	MySQLHost = loadEnv(ctx, "MYSQL_HOST", true)
	MySQLPort = loadEnv(ctx, "MYSQL_PORT", true)
	ApiPort = loadEnv(ctx, "PORT", false)
	GeminiAPIKey = loadEnv(ctx, "GEMINI_API_KEY", true)
	GeminiBaseURL = defaultString(os.Getenv("GEMINI_BASE_URL"), "https://generativelanguage.googleapis.com")

	// CORS設定（オプション、カンマ区切りで複数指定可能）
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if corsOrigins != "" {
		CORSAllowedOrigins = parseCSV(corsOrigins)
	}

	RedisConfig = CacheConfig{
		Addr:        defaultString(os.Getenv("REDIS_ADDR"), "127.0.0.1:6379"),
		Username:    os.Getenv("REDIS_USERNAME"),
		Password:    os.Getenv("REDIS_PASSWORD"),
		DB:          parseInt(os.Getenv("REDIS_DB"), 0),
		TLSEnabled:  parseBool(os.Getenv("REDIS_TLS_ENABLED"), false),
		TLSInsecure: parseBool(os.Getenv("REDIS_TLS_INSECURE"), false),
		KeyPrefix:   defaultString(os.Getenv("REDIS_KEY_PREFIX"), "app"),
		DefaultTTL:  parseDuration(os.Getenv("REDIS_DEFAULT_TTL"), 5*time.Minute),
	}
	if RedisConfig.DefaultTTL <= 0 {
		RedisConfig.DefaultTTL = 5 * time.Minute
	}
	if RedisConfig.KeyPrefix == "" {
		RedisConfig.KeyPrefix = "app"
	}
}

func defaultString(value, def string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return def
	}
	return value
}

func parseBool(value string, def bool) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return def
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return def
	}
	return b
}

func parseInt(value string, def int) int {
	value = strings.TrimSpace(value)
	if value == "" {
		return def
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return def
	}
	return i
}

func parseDuration(value string, def time.Duration) time.Duration {
	value = strings.TrimSpace(value)
	if value == "" {
		return def
	}
	dur, err := time.ParseDuration(value)
	if err != nil {
		return def
	}
	return dur
}

func parseCSV(value string) []string {
	var result []string
	for _, s := range strings.Split(value, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			result = append(result, s)
		}
	}
	return result
}
