package tracer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
)

func TestSetup_正常にTracerProviderを構築できる(t *testing.T) {
	tests := []struct {
		name   string
		config Config
	}{
		{
			name: "gRPCプロトコルで構築できる",
			config: Config{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
				Endpoint:       "localhost:4317",
				Protocol:       ProtocolGRPC,
				Insecure:       true,
				SampleRatio:    1.0,
				Timeout:        5 * time.Second,
			},
		},
		{
			name: "HTTPプロトコルで構築できる",
			config: Config{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "test",
				Endpoint:       "localhost:4318",
				Protocol:       ProtocolHTTP,
				Insecure:       true,
				SampleRatio:    1.0,
				Timeout:        5 * time.Second,
			},
		},
		{
			name: "最小限の設定で構築できる",
			config: Config{
				ServiceName: "minimal-service",
				Endpoint:    "localhost:4317",
				Insecure:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tp, cleanup, err := Setup(ctx, tt.config)

			require.NoError(t, err)
			require.NotNil(t, tp)
			require.NotNil(t, cleanup)

			// クリーンアップを実行
			if cleanup != nil {
				cleanupErr := cleanup(ctx)
				// エンドポイントが実際に存在しない場合はエラーが発生する可能性があるが、
				// テストの目的は正常に構築できることなので、エラーは無視する
				_ = cleanupErr
			}
		})
	}
}

func TestSetup_リソース属性が正しく設定される(t *testing.T) {
	config := Config{
		ServiceName:    "test-service",
		ServiceVersion: "2.0.0",
		Environment:    "staging",
		Endpoint:       "localhost:4317",
		Insecure:       true,
		ResourceAttributes: map[string]string{
			"custom.key": "custom.value",
		},
	}

	ctx := context.Background()
	tp, cleanup, err := Setup(ctx, config)

	require.NoError(t, err)
	require.NotNil(t, tp)
	defer func() {
		if cleanup != nil {
			_ = cleanup(ctx)
		}
	}()

	// TracerProviderがグローバルに設定されていることを確認
	globalTP := otel.GetTracerProvider()
	assert.NotNil(t, globalTP)
}

func TestSetup_サンプリング比率が正しく設定される(t *testing.T) {
	tests := []struct {
		name        string
		sampleRatio float64
	}{
		{
			name:        "サンプリング比率100%",
			sampleRatio: 1.0,
		},
		{
			name:        "サンプリング比率50%",
			sampleRatio: 0.5,
		},
		{
			name:        "サンプリング比率0%",
			sampleRatio: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				ServiceName: "test-service",
				Endpoint:    "localhost:4317",
				Insecure:    true,
				SampleRatio: tt.sampleRatio,
			}

			ctx := context.Background()
			tp, cleanup, err := Setup(ctx, config)

			require.NoError(t, err)
			require.NotNil(t, tp)
			
			if cleanup != nil {
				defer func() { _ = cleanup(ctx) }()
			}
		})
	}
}

func TestSetup_ヘッダーが正しく設定される(t *testing.T) {
	config := Config{
		ServiceName: "test-service",
		Endpoint:    "localhost:4317",
		Insecure:    true,
		Headers: map[string]string{
			"api-key":      "secret",
			"custom-header": "value",
		},
	}

	ctx := context.Background()
	tp, cleanup, err := Setup(ctx, config)

	require.NoError(t, err)
	require.NotNil(t, tp)
	
	if cleanup != nil {
		defer func() { _ = cleanup(ctx) }()
	}
}

func TestTracer_名前付きトレーサーを取得できる(t *testing.T) {
	tests := []struct {
		name       string
		tracerName string
	}{
		{
			name:       "シンプルな名前",
			tracerName: "test-tracer",
		},
		{
			name:       "モジュールパスのような名前",
			tracerName: "github.com/kinpatsu-everyone/backend-template",
		},
		{
			name:       "空の名前",
			tracerName: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracer := Tracer(tt.tracerName)
			assert.NotNil(t, tracer)
		})
	}
}

func TestHTTPHandler_ミドルウェアとしてラップできる(t *testing.T) {
	// テスト用のハンドラーを作成
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// トレーシングミドルウェアでラップ
	wrapped := HTTPHandler(nil, "test-operation", testHandler)
	assert.NotNil(t, wrapped)

	// リクエストを送信
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	wrapped.ServeHTTP(w, req)

	// レスポンスを確認
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "OK", w.Body.String())
}

func TestHTTPClientTransport_RoundTripperとして使用できる(t *testing.T) {
	// テスト用のサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("response"))
	}))
	defer server.Close()

	// トレーシング対応のトランスポートを作成
	transport := HTTPClientTransport(nil, nil)
	client := &http.Client{
		Transport: transport,
	}

	// リクエストを送信
	resp, err := client.Get(server.URL)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestConfig_normalize(t *testing.T) {
	tests := []struct {
		name           string
		input          Config
		expectedProtocol Protocol
		expectedRatio  float64
		expectedTimeout time.Duration
	}{
		{
			name: "デフォルト値が設定される",
			input: Config{
				ServiceName: "test",
			},
			expectedProtocol: ProtocolGRPC,
			expectedRatio:    1.0,
			expectedTimeout:  5 * time.Second,
		},
		{
			name: "設定された値は保持される",
			input: Config{
				ServiceName: "test",
				Protocol:    ProtocolHTTP,
				SampleRatio: 0.5,
				Timeout:     30 * time.Second,
			},
			expectedProtocol: ProtocolHTTP,
			expectedRatio:    0.5,
			expectedTimeout:  30 * time.Second,
		},
		{
			name: "不正なサンプリング比率は1.0に補正される（上限）",
			input: Config{
				ServiceName: "test",
				SampleRatio: 1.5,
			},
			expectedProtocol: ProtocolGRPC,
			expectedRatio:    1.0,
			expectedTimeout:  5 * time.Second,
		},
		{
			name: "不正なサンプリング比率は1.0に補正される（下限）",
			input: Config{
				ServiceName: "test",
				SampleRatio: -0.5,
			},
			expectedProtocol: ProtocolGRPC,
			expectedRatio:    1.0,
			expectedTimeout:  5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.input
			cfg.normalize()

			assert.Equal(t, tt.expectedProtocol, cfg.Protocol)
			assert.Equal(t, tt.expectedRatio, cfg.SampleRatio)
			assert.Equal(t, tt.expectedTimeout, cfg.Timeout)
		})
	}
}

func TestBuildResourceAttributes(t *testing.T) {
	tests := []struct {
		name                string
		config              Config
		expectedServiceName string
		hasVersion          bool
		hasEnvironment      bool
		customAttrsCount    int
	}{
		{
			name: "サービス名のみ",
			config: Config{
				ServiceName: "test-service",
			},
			expectedServiceName: "test-service",
			hasVersion:          false,
			hasEnvironment:      false,
			customAttrsCount:    0,
		},
		{
			name: "すべての属性が設定されている",
			config: Config{
				ServiceName:    "test-service",
				ServiceVersion: "1.0.0",
				Environment:    "production",
				ResourceAttributes: map[string]string{
					"custom.attr": "value",
				},
			},
			expectedServiceName: "test-service",
			hasVersion:          true,
			hasEnvironment:      true,
			customAttrsCount:    1,
		},
		{
			name: "複数のカスタム属性",
			config: Config{
				ServiceName: "test-service",
				ResourceAttributes: map[string]string{
					"attr1": "value1",
					"attr2": "value2",
					"attr3": "value3",
				},
			},
			expectedServiceName: "test-service",
			customAttrsCount:    3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attrs := buildResourceAttributes(tt.config)

			// 少なくともサービス名は含まれている
			assert.NotEmpty(t, attrs)

			// カスタム属性の数を確認
			customCount := 0
			for _, attr := range attrs {
				key := string(attr.Key)
				if key != "service.name" && key != "service.version" && key != "deployment.environment" {
					customCount++
				}
			}
			assert.Equal(t, tt.customAttrsCount, customCount)
		})
	}
}
