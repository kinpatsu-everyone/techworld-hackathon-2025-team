package outorouter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"reflect"
	"strconv"
)

type MultipartHandlerFunc[Req RequestObject, Res ResponseObject] func(ctx context.Context, req *Req) (*Res, error)

// MultipartEndpoint はmultipart/form-dataでファイルアップロードを受け付けるエンドポイントです
type MultipartEndpoint[Req RequestObject, Res any] struct {
	Domain     string
	Version    uint8
	MethodName string

	Summary     string
	Description string
	Tags        []Tag

	Handler MultipartHandlerFunc[Req, Res]

	// MaxMemory はmultipart/form-dataのパース時に使用する最大メモリサイズ（バイト）です
	// このサイズを超える場合は一時ファイルに保存されます
	MaxMemory int64
}

func (m MultipartEndpoint[Req, Res]) GetFullPath() string {
	return fmt.Sprintf("/%s/%s/%s", m.Domain, m.GetVersionWithPrefix(), m.MethodName)
}

func (m MultipartEndpoint[Req, Res]) GetDomain() string {
	return m.Domain
}

func (m MultipartEndpoint[Req, Res]) GetVersion() uint8 {
	return m.Version
}

func (m MultipartEndpoint[Req, Res]) GetVersionWithPrefix() string {
	return fmt.Sprintf("v%d", m.Version)
}

// RegisterMultipartEndpoint はmultipart/form-dataを受け付けるエンドポイントを登録します
func RegisterMultipartEndpoint[Req RequestObject, Res any](
	r *Router,
	ep MultipartEndpoint[Req, Res],
) {
	maxMemory := ep.MaxMemory
	if maxMemory == 0 {
		// デフォルトは32MB
		maxMemory = 32 * 1024 * 1024
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		// POSTメソッドではない場合は不正と見做す
		if req.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "指定した形式のリクエストではありません", http.StatusMethodNotAllowed)
			return
		}

		// Content-Type validation
		contentType := req.Header.Get("Content-Type")
		if contentType == "" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "Content-Type must be multipart/form-data", http.StatusUnsupportedMediaType)
			return
		}

		// multipart/form-dataであることを確認
		// boundaryが含まれているかチェック
		if len(contentType) < 19 || contentType[:19] != "multipart/form-data" {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, "Content-Type must be multipart/form-data", http.StatusUnsupportedMediaType)
			return
		}

		ctx := req.Context()

		// Set maximum request body size to 100MB (multipartは大きいファイルを扱うため)
		req.Body = http.MaxBytesReader(w, req.Body, 100*1024*1024)

		// multipart/form-dataをパース
		err := req.ParseMultipartForm(maxMemory)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, fmt.Sprintf("multipart/form-dataのパースに失敗しました: %v", err), http.StatusBadRequest)
			return
		}
		defer req.MultipartForm.RemoveAll()

		// デバッグログ: フォームデータの内容を出力
		slog.Info("multipart form data received",
			"values", req.MultipartForm.Value,
			"file_keys", func() []string {
				keys := make([]string, 0, len(req.MultipartForm.File))
				for k := range req.MultipartForm.File {
					keys = append(keys, k)
				}
				return keys
			}(),
		)

		// リクエスト構造体を作成
		var request Req
		reqType := reflect.TypeOf(request)
		hasFields := reqType.Kind() == reflect.Struct && reqType.NumField() > 0

		if hasFields {
			// リクエスト構造体のフィールドをmultipart/form-dataから埋める
			if err := populateRequestFromMultipart(&request, req.MultipartForm); err != nil {
				slog.Error("populateRequestFromMultipart failed", "error", err)
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, fmt.Sprintf("リクエストのパースに失敗しました: %v", err), http.StatusBadRequest)
				return
			}
			slog.Info("request populated", "request", fmt.Sprintf("%+v", request))
		}

		// Validate request if Validate method exists
		if err := request.Validate(); err != nil {
			slog.Error("validation failed", "error", err)
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
						"logerror": err,
						"error":    httpErr.Error(),
						"message":  httpErr.Message(),
					},
				})
				return
			}

			// その他のエラーは500として扱う
			// ロガーが設定されている場合はエラーをログに出力
			if r.logger != nil {
				r.logger.Error(ctx, "ハンドラーでエラーが発生しました", map[string]any{
					"error":  err.Error(),
					"path":   req.URL.Path,
					"method": req.Method,
				})
			}

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
			// ロガーが設定されている場合はエラーをログに出力
			if r.logger != nil {
				r.logger.Error(ctx, "レスポンスのエンコードに失敗しました", map[string]any{
					"error":  err.Error(),
					"path":   req.URL.Path,
					"method": req.Method,
				})
			}
			// 既にヘッダーが送信されている可能性があるため、エラーレスポンスは送信しない
			return
		}
	})

	r.addHTTPRoute("POST", ep.GetFullPath(), r.applyMiddlewares(h))
	r.addContentTypeRoute("POST", ep.GetFullPath(), "multipart/form-data", r.applyMiddlewares(h))

	// リクエスト・レスポンスモデルのメタデータ
	var reqZero Req
	var resZero Res

	internalEp := internalEndpoint{
		Kind:             KindFileUpload,
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
		RequestTypeInfo:  extractMultipartTypeInfo(reflect.TypeOf(reqZero)),
		ResponseTypeInfo: extractTypeInfo(reflect.TypeOf(resZero)),
	}

	r.addToRegistry(internalEp)
}

// populateRequestFromMultipart はmultipart/form-dataからリクエスト構造体を埋めます
func populateRequestFromMultipart(req interface{}, form *multipart.Form) error {
	reqValue := reflect.ValueOf(req).Elem()
	reqType := reqValue.Type()

	for i := 0; i < reqType.NumField(); i++ {
		field := reqType.Field(i)
		fieldValue := reqValue.Field(i)

		// multipartタグを優先、なければjsonタグを使用
		fieldName := ""
		multipartTag := field.Tag.Get("multipart")
		if multipartTag != "" && multipartTag != "-" {
			fieldName = multipartTag
		} else {
			// JSONタグからフィールド名を取得
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" || jsonTag == "-" {
				continue
			}
			// JSONタグからフィールド名を抽出（例: "image_data,omitempty" -> "image_data"）
			fieldName = jsonTag
			for i, r := range jsonTag {
				if r == ',' {
					fieldName = jsonTag[:i]
					break
				}
			}
		}

		// ファイルフィールドかどうかを判定
		// *multipart.FileHeader 型のフィールドはファイルとして扱う
		fileHeaderPtrType := reflect.TypeOf((*multipart.FileHeader)(nil))
		fileHeaderSliceType := reflect.TypeOf(([]*multipart.FileHeader)(nil))
		if fieldValue.Type() == fileHeaderPtrType {
			// ファイルフィールド（ポインタ型）
			if files, ok := form.File[fieldName]; ok && len(files) > 0 {
				fieldValue.Set(reflect.ValueOf(files[0]))
			}
		} else if fieldValue.Type() == fileHeaderSliceType {
			// ファイル配列フィールド
			if files, ok := form.File[fieldName]; ok {
				fieldValue.Set(reflect.ValueOf(files))
			}
		} else {
			// 通常のテキストフィールド
			if values, ok := form.Value[fieldName]; ok && len(values) > 0 {
				value := values[0]
				switch fieldValue.Kind() {
				case reflect.String:
					fieldValue.SetString(value)
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if intVal, err := strconv.ParseInt(value, 10, 64); err == nil {
						fieldValue.SetInt(intVal)
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
						fieldValue.SetUint(uintVal)
					}
				case reflect.Float32, reflect.Float64:
					if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
						fieldValue.SetFloat(floatVal)
					}
				case reflect.Bool:
					if boolVal, err := strconv.ParseBool(value); err == nil {
						fieldValue.SetBool(boolVal)
					}
				case reflect.Slice:
					if fieldValue.Type().Elem().Kind() == reflect.String {
						fieldValue.Set(reflect.ValueOf(values))
					}
				default:
					// JSON文字列としてパースを試みる
					if len(value) > 0 && (value[0] == '{' || value[0] == '[') {
						if err := json.Unmarshal([]byte(value), fieldValue.Addr().Interface()); err != nil {
							// JSONパースに失敗した場合は文字列として扱う
							if fieldValue.Kind() == reflect.String {
								fieldValue.SetString(value)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// extractMultipartTypeInfo はmultipart/form-dataリクエスト型からTypeInfoを抽出します
// multipartタグとjsonタグの両方をサポートします
func extractMultipartTypeInfo(t reflect.Type) TypeInfo {
	visited := make(map[reflect.Type]bool)
	return extractMultipartTypeInfoRecursive(t, visited)
}

// extractMultipartTypeInfoRecursive はGoの型からTypeInfoを再帰的に抽出します（multipart用）
func extractMultipartTypeInfoRecursive(t reflect.Type, visited map[reflect.Type]bool) TypeInfo {
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

		// multipartタグを優先し、なければJSONタグを使用
		fieldName := ""
		optional := false

		multipartTag := field.Tag.Get("multipart")
		if multipartTag != "" && multipartTag != "-" {
			fieldName = multipartTag
		} else {
			jsonTag := field.Tag.Get("json")
			fieldName, optional = parseJSONTag(jsonTag, field.Name)
		}

		// "-"の場合はスキップ
		if fieldName == "-" {
			continue
		}

		fieldInfo := FieldInfo{
			Name:     field.Name,
			JSONName: fieldName,
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
