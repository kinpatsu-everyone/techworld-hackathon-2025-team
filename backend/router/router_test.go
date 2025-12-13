package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuild(t *testing.T) {
	tests := []struct {
		name          string
		setupRouter   func() *outorouter.Router
		shouldError   bool
	}{
		{
			name: "正常にルーターをビルドできる",
			setupRouter: func() *outorouter.Router {
				return outorouter.New()
			},
			shouldError: false,
		},
		{
			name: "ミドルウェアが設定されたルーターをビルドできる",
			setupRouter: func() *outorouter.Router {
				r := outorouter.New()
				r.Use(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
						next.ServeHTTP(w, req)
					})
				})
				return r
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := tt.setupRouter()
			
			handler, err := Build(router)

			if tt.shouldError {
				assert.Error(t, err)
				assert.Nil(t, handler)
			} else {
				require.NoError(t, err)
				require.NotNil(t, handler)
			}
		})
	}
}

func TestBuild_ヘルスチェックエンドポイントが登録されている(t *testing.T) {
	router := outorouter.New()
	handler, err := Build(router)
	
	require.NoError(t, err)
	require.NotNil(t, handler)

	// ヘルスチェックエンドポイントにリクエストを送信
	req := httptest.NewRequest(http.MethodPost, "/healthz/v1/Healthz", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// レスポンスを確認
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")
	
	body := w.Body.String()
	assert.Contains(t, body, "version")
	assert.Contains(t, body, "status")
	assert.Contains(t, body, "message")
	assert.Contains(t, body, "ok")
}

func TestBuild_POSTメソッドのみ受け付ける(t *testing.T) {
	router := outorouter.New()
	handler, err := Build(router)
	
	require.NoError(t, err)
	require.NotNil(t, handler)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
	}{
		{
			name:           "POSTメソッドは成功する",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GETメソッドは404を返す（ルートが登録されていない）",
			method:         http.MethodGet,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "PUTメソッドは404を返す（ルートが登録されていない）",
			method:         http.MethodPut,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "DELETEメソッドは404を返す（ルートが登録されていない）",
			method:         http.MethodDelete,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/healthz/v1/Healthz", strings.NewReader("{}"))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBuild_存在しないパスへのリクエストは404を返す(t *testing.T) {
	router := outorouter.New()
	handler, err := Build(router)
	
	require.NoError(t, err)
	require.NotNil(t, handler)

	tests := []struct {
		name string
		path string
	}{
		{
			name: "存在しないドメイン",
			path: "/notfound/v1/Test",
		},
		{
			name: "存在しないバージョン",
			path: "/healthz/v99/Healthz",
		},
		{
			name: "存在しないメソッド",
			path: "/healthz/v1/NotExists",
		},
		{
			name: "ルートパス",
			path: "/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader("{}"))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, http.StatusNotFound, w.Code)
		})
	}
}

func TestBuild_正しいJSONレスポンスを返す(t *testing.T) {
	router := outorouter.New()
	handler, err := Build(router)
	
	require.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/healthz/v1/Healthz", strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	body := w.Body.String()
	// JSONとして正しい形式であることを確認
	assert.True(t, strings.HasPrefix(strings.TrimSpace(body), "{"))
	assert.True(t, strings.HasSuffix(strings.TrimSpace(body), "}"))
	
	// 必須フィールドが含まれていることを確認
	assert.Contains(t, body, `"version"`)
	assert.Contains(t, body, `"status"`)
	assert.Contains(t, body, `"message"`)
}
