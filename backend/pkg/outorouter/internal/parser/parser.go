package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
)

// ParseFile は metadata.json のパスを受け取り、中間表現を返す。
func ParseFile(filePath string) (*Metadata, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open metadata file: %w", err)
	}
	defer f.Close()

	return Parse(f)
}

// Parse は metadata.json の内容を解析し、中間表現を返す。
func Parse(r io.Reader) (*Metadata, error) {
	var raw rawMetadata
	dec := json.NewDecoder(r)
	if err := dec.Decode(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode metadata: %w", err)
	}

	meta := &Metadata{
		Domains: make(map[string]Domain),
		All:     make([]Endpoint, 0),
	}

	for domain, versions := range raw {
		if domain == "" {
			return nil, fmt.Errorf("metadata contains empty domain key")
		}

		if _, exists := meta.Domains[domain]; !exists {
			meta.Domains[domain] = Domain{Versions: make(map[uint8][]Endpoint)}
		}

		for versionKey, endpoints := range versions {
			version, err := parseVersion(versionKey)
			if err != nil {
				return nil, fmt.Errorf("domain %q: %w", domain, err)
			}

			for idx, ep := range endpoints {
				converted, err := ep.toEndpoint(domain, version)
				if err != nil {
					return nil, fmt.Errorf("domain %q version %d index %d: %w", domain, version, idx, err)
				}

				meta.Domains[domain].Versions[version] = append(meta.Domains[domain].Versions[version], converted)
				meta.All = append(meta.All, converted)
			}
		}
	}

	return meta, nil
}

func parseVersion(key string) (uint8, error) {
	v, err := strconv.ParseUint(key, 10, 8)
	if err != nil {
		return 0, fmt.Errorf("invalid version key %q: %w", key, err)
	}
	return uint8(v), nil
}

type rawMetadata map[string]map[string][]rawEndpoint

type rawEndpoint struct {
	Kind         string   `json:"kind"`
	MethodName   string   `json:"method_name"`
	HTTPMethod   string   `json:"http_method"`
	RequestType  string   `json:"request_type"`
	ResponseType string   `json:"response_type"`
	Summary      string   `json:"summary"`
	Description  string   `json:"description"`
	Tags         []Tag    `json:"tags"`
	Domain       string   `json:"domain"`
	Version      uint8    `json:"version"`

	// 型情報
	RequestTypeInfo  rawTypeInfo `json:"request_type_info"`
	ResponseTypeInfo rawTypeInfo `json:"response_type_info"`
}

type rawTypeInfo struct {
	Name   string         `json:"name"`
	Fields []rawFieldInfo `json:"fields"`
}

type rawFieldInfo struct {
	Name       string       `json:"name"`
	JSONName   string       `json:"json_name"`
	Type       string       `json:"type"`
	TSType     string       `json:"ts_type"`
	Optional   bool         `json:"optional"`
	NestedType *rawTypeInfo `json:"nested_type,omitempty"`
}

func (r rawEndpoint) toEndpoint(domain string, version uint8) (Endpoint, error) {
	kind, err := toKind(r.Kind)
	if err != nil {
		return Endpoint{}, err
	}

	if r.MethodName == "" {
		return Endpoint{}, fmt.Errorf("method_name is empty")
	}

	if r.HTTPMethod == "" {
		return Endpoint{}, fmt.Errorf("http_method is empty")
	}

	return Endpoint{
		Kind:             kind,
		Domain:           domain,
		Version:          version,
		MethodName:       r.MethodName,
		HTTPMethod:       r.HTTPMethod,
		RequestType:      r.RequestType,
		ResponseType:     r.ResponseType,
		Summary:          r.Summary,
		Description:      r.Description,
		Tags:             r.Tags,
		RequestTypeInfo:  convertTypeInfo(r.RequestTypeInfo),
		ResponseTypeInfo: convertTypeInfo(r.ResponseTypeInfo),
	}, nil
}

func convertTypeInfo(raw rawTypeInfo) TypeInfo {
	fields := make([]FieldInfo, len(raw.Fields))
	for i, f := range raw.Fields {
		fieldInfo := FieldInfo{
			Name:     f.Name,
			JSONName: f.JSONName,
			Type:     f.Type,
			TSType:   f.TSType,
			Optional: f.Optional,
		}
		// ネストされた型情報があれば再帰的に変換
		if f.NestedType != nil {
			nestedTypeInfo := convertTypeInfo(*f.NestedType)
			fieldInfo.NestedType = &nestedTypeInfo
		}
		fields[i] = fieldInfo
	}
	return TypeInfo{
		Name:   raw.Name,
		Fields: fields,
	}
}

func toKind(k string) (EndpointKind, error) {
	switch EndpointKind(k) {
	case KindUnaryJSON, KindFileUpload, KindWebSocket:
		return EndpointKind(k), nil
	default:
		return "", fmt.Errorf("unknown kind %q", k)
	}
}
