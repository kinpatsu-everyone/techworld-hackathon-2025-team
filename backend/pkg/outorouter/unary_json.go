package outorouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

type UnaryJSONHandlerFunc[Req RequestObject, Res ResponseObject] func(ctx context.Context, req *Req) (*Res, error)

// UnaryJSONEndpoint はHTTP 1.1 のPOSTメソッドでやり取りするためのエンドポイントです
type UnaryJSONEndpoint[Req RequestObject, Res any] struct {
	Domain     string
	Version    uint8
	MethodName string

	Summary     string
	Description string
	Tags        []Tag

	Handler UnaryJSONHandlerFunc[Req, Res]
}

func (u UnaryJSONEndpoint[Req, Res]) GetFullPath() string {
	return fmt.Sprintf("/%s/%s/%s", u.Domain, u.GetVersionWithPrefix(), u.MethodName)
}

func (u UnaryJSONEndpoint[Req, Res]) GetDomain() string {
	return u.Domain
}

func (u UnaryJSONEndpoint[Req, Res]) GetVersion() uint8 {
	return u.Version
}

func (u UnaryJSONEndpoint[Req, Res]) GetVersionWithPrefix() string {
	return fmt.Sprintf("v%d", u.Version)
}

func RegisterUnaryJSONEndpoint[Req RequestObject, Res any](
	r *Router,
	ep UnaryJSONEndpoint[Req, Res],
) {
	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// POSTメソッドではない場合は不正と見做す
		if req.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "指定した形式のリクエストではありません", http.StatusMethodNotAllowed)
			return
		}

		// Content-Type validation
		contentType := req.Header.Get("Content-Type")
		if contentType != "" && contentType != "application/json" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		ctx := req.Context()

		var request Req
		// リクエストのパースとバリデーションを行う
		// リクエスト型にフィールドがない場合はJSONパースをスキップ
		reqType := reflect.TypeOf(request)
		hasFields := reqType.Kind() == reflect.Struct && reqType.NumField() > 0

		if hasFields && req.Body != nil {
			defer req.Body.Close()
			// Set maximum request body size to 10MB
			req.Body = http.MaxBytesReader(w, req.Body, 10*1024*1024)

			// ここでリクエストボディをパースして request にセットする
			if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, "リクエストのJSON形式が不正です", http.StatusBadRequest)
				return
			}
		}

		// Validate request if Validate method exists
		if err := request.Validate(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, fmt.Sprintf("リクエストのバリデーションに失敗しました: %v", err), http.StatusBadRequest)
			return
		}

		response, err := ep.Handler(ctx, &request)
		if err != nil {
			var httpErr HTTPError
			if errors.As(err, &httpErr) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(httpErr.StatusCode())

				_ = json.NewEncoder(w).Encode(map[string]any{
					"error": map[string]any{
						"error":   httpErr.Error(),
						"message": httpErr.Message(),
					},
				})
				return
			}

			// その他のエラーは500として扱う
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)

			_ = json.NewEncoder(w).Encode(map[string]any{
				"error": map[string]any{
					"error":   "UNKNOWN_INTERNAL_ERROR",
					"message": "サーバー内部で予期しないエラーが発生しました",
				},
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "レスポンスの返却に失敗しました", http.StatusInternalServerError)
			return
		}
	})

	r.addHTTPRoute("POST", ep.GetFullPath(), r.applyMiddlewares(h))
	r.addContentTypeRoute("POST", ep.GetFullPath(), "application/json", r.applyMiddlewares(h))

	// リクエスト・レスポンスモデルのメタデータ
	var reqZero Req
	var resZero Res

	internalEp := internalEndpoint{
		Kind:             KindUnaryJSON,
		Domain:           ep.Domain,
		Version:          ep.Version,
		MethodName:       ep.MethodName,
		Summary:          ep.Summary,
		Description:      ep.Description,
		Tags:             ep.Tags,
		HTTPMethod:       http.MethodPost,
		handler:          h,
		RequestType:      reflect.TypeOf(reqZero).String(),
		ResponseType:     reflect.TypeOf(resZero).String(),
		RequestTypeInfo:  extractTypeInfo(reflect.TypeOf(reqZero)),
		ResponseTypeInfo: extractTypeInfo(reflect.TypeOf(resZero)),
	}

	r.addToRegistry(internalEp)
}

// extractTypeInfo はGoの型からTypeInfoを抽出します
func extractTypeInfo(t reflect.Type) TypeInfo {
	visited := make(map[reflect.Type]bool)
	return extractTypeInfoRecursive(t, visited)
}

// extractTypeInfoRecursive はGoの型からTypeInfoを再帰的に抽出します
// visitedは循環参照を防ぐために使用します
func extractTypeInfoRecursive(t reflect.Type, visited map[reflect.Type]bool) TypeInfo {
	// ポインタの場合は要素型を取得
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	info := TypeInfo{
		Name:   t.Name(),
		Fields: make([]FieldInfo, 0),
	}

	// 構造体でない場合は空のフィールドリストを返す
	if t.Kind() != reflect.Struct {
		return info
	}

	// 循環参照チェック
	if visited[t] {
		return info
	}
	visited[t] = true

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 非公開フィールドはスキップ
		if !field.IsExported() {
			continue
		}

		// JSONタグを解析
		jsonTag := field.Tag.Get("json")
		jsonName, optional := parseJSONTag(jsonTag, field.Name)

		// "-"の場合はスキップ（JSONで無視されるフィールド）
		if jsonName == "-" {
			continue
		}

		fieldInfo := FieldInfo{
			Name:     field.Name,
			JSONName: jsonName,
			Type:     field.Type.String(),
			TSType:   goTypeToTSType(field.Type),
			Optional: optional,
		}

		// ネストされた構造体の型情報を抽出
		nestedType := extractNestedTypeInfo(field.Type, visited)
		if nestedType != nil {
			fieldInfo.NestedType = nestedType
		}

		info.Fields = append(info.Fields, fieldInfo)
	}

	return info
}

// extractNestedTypeInfo はフィールドの型からネストされた構造体の型情報を抽出します
func extractNestedTypeInfo(t reflect.Type, visited map[reflect.Type]bool) *TypeInfo {
	// ポインタの場合は要素型を取得
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// スライスまたは配列の場合は要素型を取得
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		t = t.Elem()
		// 要素がポインタの場合はさらに要素型を取得
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
	}

	// 構造体の場合のみ型情報を抽出
	if t.Kind() == reflect.Struct {
		// time.Time などの標準ライブラリの型はスキップ
		if t.PkgPath() == "time" {
			return nil
		}
		// multipart.FileHeader などのmime/multipartパッケージの型はスキップ
		if t.PkgPath() == "mime/multipart" {
			return nil
		}

		// 循環参照チェック
		if visited[t] {
			return nil
		}

		typeInfo := extractTypeInfoRecursive(t, visited)
		return &typeInfo
	}

	return nil
}

// parseJSONTag はJSONタグを解析してフィールド名とomitemptyの有無を返します
func parseJSONTag(tag string, defaultName string) (name string, optional bool) {
	if tag == "" {
		return defaultName, false
	}

	parts := strings.Split(tag, ",")
	name = parts[0]
	if name == "" {
		name = defaultName
	}

	for _, part := range parts[1:] {
		if part == "omitempty" {
			optional = true
			break
		}
	}

	return name, optional
}

// goTypeToTSType はGoの型をTypeScriptの型に変換します
func goTypeToTSType(t reflect.Type) string {
	// ポインタの場合は要素型を取得
	if t.Kind() == reflect.Ptr {
		return goTypeToTSType(t.Elem())
	}

	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		elemType := goTypeToTSType(t.Elem())
		return elemType + "[]"
	case reflect.Map:
		keyType := goTypeToTSType(t.Key())
		valType := goTypeToTSType(t.Elem())
		return fmt.Sprintf("Record<%s, %s>", keyType, valType)
	case reflect.Struct:
		// time.Time は特別扱い
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return "string" // ISO 8601形式の文字列として扱う
		}
		// その他の構造体はインライン型として展開するか、anyとして扱う
		return t.Name()
	case reflect.Interface:
		// any型
		if t.NumMethod() == 0 {
			return "unknown"
		}
		return "unknown"
	default:
		return "unknown"
	}
}
