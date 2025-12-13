package parser

import "fmt"

// Tag はAPIタグを表す文字列型です。
type Tag string

// EndpointKind はエンドポイントの種類を表す文字列型です。
type EndpointKind string

const (
	KindUnaryJSON  EndpointKind = "JSON"
	KindFileUpload EndpointKind = "FileUpload"
	KindWebSocket  EndpointKind = "WebSocket"
)

// FieldInfo はGoの構造体フィールドの情報を保持します
type FieldInfo struct {
	Name       string    `json:"name"`
	JSONName   string    `json:"json_name"`
	Type       string    `json:"type"`
	TSType     string    `json:"ts_type"`
	Optional   bool      `json:"optional"`
	NestedType *TypeInfo `json:"nested_type,omitempty"` // ネストされた構造体の型情報
}

// TypeInfo は構造体の型情報を保持します
type TypeInfo struct {
	Name   string      `json:"name"`
	Fields []FieldInfo `json:"fields"`
}

// Endpoint は metadata.json の1エントリを中間表現として保持する。
type Endpoint struct {
	Kind         EndpointKind
	Domain       string
	Version      uint8
	MethodName   string
	HTTPMethod   string
	RequestType  string
	ResponseType string
	Summary      string
	Description  string
	Tags         []Tag

	// 型情報
	RequestTypeInfo  TypeInfo
	ResponseTypeInfo TypeInfo
}

// Metadata はドメイン別・バージョン別のエンドポイント集合を表す。
type Metadata struct {
	Domains map[string]Domain
	All     []Endpoint
}

// Domain はバージョンごとのエンドポイント群を持つ。
type Domain struct {
	Versions map[uint8][]Endpoint
}

// Path はエンドポイントのパス (domain/v{version}/method) を返す。
func (e Endpoint) Path() string {
	return fmt.Sprintf("%s/v%d/%s", e.Domain, e.Version, e.MethodName)
}
