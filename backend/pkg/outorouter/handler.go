package outorouter

import (
	"context"
	"net/http"
	"sync"
)

// HandlerFunc は AutoRouter 拡張のハンドラ関数の型です。
type HandlerFunc[Req any, Res any] func(ctx context.Context, req *Req) (*Res, error)

// MiddlewareFunc は AutoRouter 拡張のミドルウェア関数の型です。
type MiddlewareFunc func(next http.Handler) http.Handler

// RequestObject はHTTP POSTリクエストオブジェクトのインターフェースです。
type RequestObject interface {
	Validate() error
}

// ResponseObject はHTTP POSTレスポンスオブジェクトのインターフェースです。
type ResponseObject interface {
}

type ErrorResponseObject interface {
	ResponseObject
	StatusCode() int
}

type Tag string

func (t Tag) String() string {
	return string(t)
}

func RegisterTags(tags ...Tag) []Tag {
	return tags
}

// MappingResponseObjectFunc はドメインレスポンスをレスポンスオブジェクトに変換する関数の型です。
type MappingResponseObjectFunc[DomainRes any] func(domainRes *DomainRes) (ResponseObject, error)

type EndpointKind string

func (k EndpointKind) String() string {
	return string(k)
}

const (
	KindUnaryJSON  EndpointKind = "JSON"
	KindFileUpload EndpointKind = "FileUpload"
	KindWebSocket  EndpointKind = "WebSocket"
)

type Endpoint interface {
	GetFullPath() string
	GetDomain() string
	GetVersion() uint8
	GetVersionWithPrefix() string
}

// FieldInfo はGoの構造体フィールドの情報を保持します（コード生成用）
type FieldInfo struct {
	Name       string    `json:"name"`                  // Goのフィールド名
	JSONName   string    `json:"json_name"`             // JSONタグの名前
	Type       string    `json:"type"`                  // Goの型名
	TSType     string    `json:"ts_type"`               // TypeScriptの型名
	Optional   bool      `json:"optional"`              // omitemptyの有無
	NestedType *TypeInfo `json:"nested_type,omitempty"` // ネストされた構造体の型情報（構造体またはスライス/配列の要素が構造体の場合）
}

// TypeInfo は構造体の型情報を保持します
type TypeInfo struct {
	Name   string      `json:"name"`   // 型名
	Fields []FieldInfo `json:"fields"` // フィールド情報
}

type internalEndpoint struct {
	Kind EndpointKind

	Domain     string
	Version    uint8
	MethodName string

	Summary     string
	Description string
	Tags        []Tag

	HTTPMethod string
	handler    http.Handler

	// 型名 (コード生成用)
	RequestType  string
	ResponseType string

	// 型情報 (コード生成用)
	RequestTypeInfo  TypeInfo
	ResponseTypeInfo TypeInfo
}

type Router struct {
	mu sync.RWMutex

	// domain → version → kind → []internalEndpoint
	registry map[string]map[uint8]map[EndpointKind][]internalEndpoint

	// 実行時ルーティング: method → fullPath → handler
	httpRoutes map[string]map[string]http.Handler

	// Content-Type別ルーティング: method → fullPath → Content-Type → handler
	contentTypeRoutes map[string]map[string]map[string]http.Handler

	middlewares []MiddlewareFunc
	// TODO: loggerを組み込む
	logger Logger
}

type Option func(*Router)

func New(opts ...Option) *Router {
	r := &Router{
		registry:          make(map[string]map[uint8]map[EndpointKind][]internalEndpoint),
		httpRoutes:        make(map[string]map[string]http.Handler),
		contentTypeRoutes: make(map[string]map[string]map[string]http.Handler),
		middlewares:       make([]MiddlewareFunc, 0),
		logger:            nil,
	}

	for _, opt := range opts {
		opt(r)
	}

	// デフォルトで利用するミドルウェアを追加

	return r
}

// GetRegistries は登録されているエンドポイントのレジストリを取得します。
func (r *Router) GetRegistries() map[string]map[uint8]map[EndpointKind][]internalEndpoint {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.registry
}

func (r *Router) Handler() http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		r.mu.RLock()

		// Content-Type別ルーティングを優先的にチェック
		contentType := req.Header.Get("Content-Type")

		// multipart/form-dataの場合、boundaryパラメータを除去して判定
		normalizedContentType := contentType
		if len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
			normalizedContentType = "multipart/form-data"
		}

		// デバッグ: ログが有効な場合はリクエスト情報を出力
		if r.logger != nil {
			r.logger.Info(req.Context(), "request routing", map[string]any{
				"method":       req.Method,
				"path":         req.URL.Path,
				"content_type": contentType,
				"normalized":   normalizedContentType,
			})
		}

		if contentTypeRoutes, methodExists := r.contentTypeRoutes[req.Method]; methodExists {
			if pathRoutes, pathExists := contentTypeRoutes[req.URL.Path]; pathExists {
				// Content-Typeに基づいてハンドラーを選択
				if h, exists := pathRoutes[normalizedContentType]; exists {
					r.mu.RUnlock()
					h.ServeHTTP(res, req)
					return
				}
				// Content-Typeが空の場合はapplication/jsonを試す
				if contentType == "" {
					if h, exists := pathRoutes["application/json"]; exists {
						r.mu.RUnlock()
						h.ServeHTTP(res, req)
						return
					}
				}
			}
		}

		// フォールバック: 通常のルーティング
		// 通常のルーティングでも、Content-Typeに応じて適切なハンドラーを選択
		methodRoutes, methodExists := r.httpRoutes[req.Method]
		if !methodExists {
			r.mu.RUnlock()
			// デバッグ: ログが有効な場合は404の原因を出力
			if r.logger != nil {
				r.logger.Info(req.Context(), "404: method not found", map[string]any{
					"method": req.Method,
					"path":   req.URL.Path,
				})
			}
			http.NotFound(res, req)
			return
		}

		h, pathExists := methodRoutes[req.URL.Path]
		r.mu.RUnlock()

		if !pathExists {
			// デバッグ: ログが有効な場合は404の原因を出力
			if r.logger != nil {
				r.logger.Info(req.Context(), "404: path not found", map[string]any{
					"method": req.Method,
					"path":   req.URL.Path,
					"available_paths": func() []string {
						r.mu.RLock()
						defer r.mu.RUnlock()
						paths := make([]string, 0)
						if routes, ok := r.httpRoutes[req.Method]; ok {
							for p := range routes {
								paths = append(paths, p)
							}
						}
						return paths
					}(),
				})
			}
			http.NotFound(res, req)
			return
		}

		h.ServeHTTP(res, req)
	})
}

// Use はミドルウェアを追加します。
func (r *Router) Use(mws ...MiddlewareFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.middlewares = append(r.middlewares, mws...)
}

// RegisterCustomHandler はカスタムHTTPハンドラーを登録します
func (r *Router) RegisterCustomHandler(method, path string, h http.Handler) {
	r.addHTTPRoute(method, path, r.applyMiddlewares(h))
}

func (r *Router) applyMiddlewares(h http.Handler) http.Handler {
	r.mu.RLock()
	mws := make([]MiddlewareFunc, len(r.middlewares))
	copy(mws, r.middlewares)
	r.mu.RUnlock()

	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

func (r *Router) addHTTPRoute(method, fullPath string, h http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.httpRoutes[method]; !exists {
		r.httpRoutes[method] = make(map[string]http.Handler)
	}

	r.httpRoutes[method][fullPath] = h
}

func (r *Router) addContentTypeRoute(method, fullPath, contentType string, h http.Handler) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.contentTypeRoutes[method]; !exists {
		r.contentTypeRoutes[method] = make(map[string]map[string]http.Handler)
	}

	if _, exists := r.contentTypeRoutes[method][fullPath]; !exists {
		r.contentTypeRoutes[method][fullPath] = make(map[string]http.Handler)
	}

	r.contentTypeRoutes[method][fullPath][contentType] = h
}

func (r *Router) addToRegistry(ep internalEndpoint) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// ドメインが存在しない場合は新たに追加する
	if _, ok := r.registry[ep.Domain]; !ok {
		r.registry[ep.Domain] = make(map[uint8]map[EndpointKind][]internalEndpoint)
	}

	// そのドメインの指定のバージョンが存在しない場合は新たに追加する
	if _, ok := r.registry[ep.Domain][ep.Version]; !ok {
		r.registry[ep.Domain][ep.Version] = make(map[EndpointKind][]internalEndpoint)
	}

	r.registry[ep.Domain][ep.Version][ep.Kind] = append(r.registry[ep.Domain][ep.Version][ep.Kind], ep)
}
