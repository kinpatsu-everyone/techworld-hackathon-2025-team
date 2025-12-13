package handler

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthzRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request HealthzRequest
	}{
		{
			name:    "空のリクエストは常に有効",
			request: HealthzRequest{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			assert.NoError(t, err)
		})
	}
}

func TestHealthz(t *testing.T) {
	tests := []struct {
		name             string
		request          *HealthzRequest
		expectedVersion  string
		expectedStatus   string
		expectedMessage  string
		shouldReturnError bool
	}{
		{
			name:             "正常なヘルスチェックレスポンスを返す",
			request:          &HealthzRequest{},
			expectedVersion:  "1.0.0",
			expectedStatus:   "ok",
			expectedMessage:  "Service is healthy",
			shouldReturnError: false,
		},
		{
			name:             "nilリクエストでも正常に動作する",
			request:          nil,
			expectedVersion:  "1.0.0",
			expectedStatus:   "ok",
			expectedMessage:  "Service is healthy",
			shouldReturnError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			
			response, err := Healthz(ctx, tt.request)

			if tt.shouldReturnError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				require.NoError(t, err)
				require.NotNil(t, response)
				assert.Equal(t, tt.expectedVersion, response.Version)
				assert.Equal(t, tt.expectedStatus, response.Status)
				assert.Equal(t, tt.expectedMessage, response.Message)
			}
		})
	}
}

func TestHealthzResponse_構造体フィールド(t *testing.T) {
	tests := []struct {
		name     string
		response HealthzResponse
	}{
		{
			name: "すべてのフィールドが設定されている",
			response: HealthzResponse{
				Version: "1.0.0",
				Status:  "ok",
				Message: "Service is healthy",
			},
		},
		{
			name: "空のレスポンスも作成できる",
			response: HealthzResponse{},
		},
		{
			name: "部分的に設定されたレスポンスも作成できる",
			response: HealthzResponse{
				Status: "ok",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 構造体が正しく作成できることを確認
			assert.NotNil(t, tt.response)
		})
	}
}

func TestHealthz_コンテキストキャンセル時の動作(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // 即座にキャンセル

	request := &HealthzRequest{}
	
	// キャンセルされたコンテキストでもHealthzは正常に動作する
	// （実装が単純でコンテキストを使用していないため）
	response, err := Healthz(ctx, request)

	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "ok", response.Status)
}
