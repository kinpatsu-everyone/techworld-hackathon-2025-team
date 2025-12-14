package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseBool(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		def      bool
		expected bool
	}{
		{
			name:     "空文字列の場合はデフォルト値を返す",
			value:    "",
			def:      true,
			expected: true,
		},
		{
			name:     "trueの文字列はtrueを返す",
			value:    "true",
			def:      false,
			expected: true,
		},
		{
			name:     "falseの文字列はfalseを返す",
			value:    "false",
			def:      true,
			expected: false,
		},
		{
			name:     "1はtrueを返す",
			value:    "1",
			def:      false,
			expected: true,
		},
		{
			name:     "0はfalseを返す",
			value:    "0",
			def:      true,
			expected: false,
		},
		{
			name:     "不正な値の場合はデフォルト値を返す",
			value:    "invalid",
			def:      true,
			expected: true,
		},
		{
			name:     "前後の空白は無視される",
			value:    "  true  ",
			def:      false,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseBool(tt.value, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		def      int
		expected int
	}{
		{
			name:     "空文字列の場合はデフォルト値を返す",
			value:    "",
			def:      42,
			expected: 42,
		},
		{
			name:     "正の整数をパースできる",
			value:    "123",
			def:      0,
			expected: 123,
		},
		{
			name:     "負の整数をパースできる",
			value:    "-456",
			def:      0,
			expected: -456,
		},
		{
			name:     "0をパースできる",
			value:    "0",
			def:      42,
			expected: 0,
		},
		{
			name:     "不正な値の場合はデフォルト値を返す",
			value:    "invalid",
			def:      42,
			expected: 42,
		},
		{
			name:     "前後の空白は無視される",
			value:    "  789  ",
			def:      0,
			expected: 789,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseInt(tt.value, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		def      time.Duration
		expected time.Duration
	}{
		{
			name:     "空文字列の場合はデフォルト値を返す",
			value:    "",
			def:      5 * time.Second,
			expected: 5 * time.Second,
		},
		{
			name:     "秒単位の時間をパースできる",
			value:    "10s",
			def:      0,
			expected: 10 * time.Second,
		},
		{
			name:     "分単位の時間をパースできる",
			value:    "5m",
			def:      0,
			expected: 5 * time.Minute,
		},
		{
			name:     "時間単位の時間をパースできる",
			value:    "2h",
			def:      0,
			expected: 2 * time.Hour,
		},
		{
			name:     "ミリ秒単位の時間をパースできる",
			value:    "500ms",
			def:      0,
			expected: 500 * time.Millisecond,
		},
		{
			name:     "複合的な時間をパースできる",
			value:    "1h30m",
			def:      0,
			expected: 90 * time.Minute,
		},
		{
			name:     "不正な値の場合はデフォルト値を返す",
			value:    "invalid",
			def:      5 * time.Second,
			expected: 5 * time.Second,
		},
		{
			name:     "前後の空白は無視される",
			value:    "  3s  ",
			def:      0,
			expected: 3 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDuration(tt.value, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDefaultString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		def      string
		expected string
	}{
		{
			name:     "空文字列の場合はデフォルト値を返す",
			value:    "",
			def:      "default",
			expected: "default",
		},
		{
			name:     "値がある場合はその値を返す",
			value:    "value",
			def:      "default",
			expected: "value",
		},
		{
			name:     "前後の空白のみの場合はデフォルト値を返す",
			value:    "   ",
			def:      "default",
			expected: "default",
		},
		{
			name:     "前後の空白は削除される",
			value:    "  value  ",
			def:      "default",
			expected: "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := defaultString(tt.value, tt.def)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadEnv_必須環境変数が設定されている場合(t *testing.T) {
	// テスト用の環境変数を設定
	envVars := map[string]string{
		"ENV":            "test",
		"MYSQL_USER":     "testuser",
		"MYSQL_PASSWORD": "testpass",
		"MYSQL_DATABASE": "testdb",
		"MYSQL_HOST":     "localhost",
		"MYSQL_PORT":     "3306",
		"PORT":           "8080",
		"GEMINI_API_KEY": "test-api-key",
	}

	// 環境変数を設定し、テスト後に元に戻す
	for key, value := range envVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	ctx := context.Background()

	// LoadEnvがパニックしないことを確認
	require.NotPanics(t, func() {
		LoadEnv(ctx)
	})

	// 設定された値を確認
	assert.Equal(t, "test", ENV)
	assert.Equal(t, "testuser", MySQLUser)
	assert.Equal(t, "testpass", MySQLPassword)
	assert.Equal(t, "testdb", MySQLDatabase)
	assert.Equal(t, "localhost", MySQLHost)
	assert.Equal(t, "3306", MySQLPort)
	assert.Equal(t, "8080", ApiPort)
}

func TestLoadEnv_必須環境変数が不足している場合はパニックする(t *testing.T) {
	// 環境変数をクリア
	envVars := []string{
		"ENV", "MYSQL_USER", "MYSQL_PASSWORD", "MYSQL_DATABASE",
		"MYSQL_HOST", "MYSQL_PORT", "PORT", "GEMINI_API_KEY",
	}
	for _, key := range envVars {
		os.Unsetenv(key)
	}

	ctx := context.Background()

	// LoadEnvがパニックすることを確認
	assert.Panics(t, func() {
		LoadEnv(ctx)
	})
}
