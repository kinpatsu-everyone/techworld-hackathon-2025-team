package outorouter

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ExportedEndpoint struct {
	Kind         string `json:"kind"`
	Domain       string `json:"domain"`
	Version      uint8  `json:"version"`
	MethodName   string `json:"method_name"`
	HTTPMethod   string `json:"http_method"`
	RequestType  string `json:"request_type"`
	ResponseType string `json:"response_type"`
	Summary      string `json:"summary"`
	Description  string `json:"description"`
	Tags         []Tag  `json:"tags"`

	// 型情報（コード生成用）
	RequestTypeInfo  TypeInfo `json:"request_type_info"`
	ResponseTypeInfo TypeInfo `json:"response_type_info"`
}

// ExportMetadataJSON はルーターのメタデータを JSON ファイルとしてエクスポートします。
// オプションでファイルパスを指定できます。デフォルトは ".api/metadata.json" です。
func ExportMetadataJSON(router *Router, opts ...ExportMetadataJSONOptions) error {
	options := defaultExportMetadataJSONOptions
	for _, opt := range opts {
		opt(&options)
	}

	if err := os.MkdirAll(filepath.Dir(options.filePath), 0o755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	f, err := os.Create(options.filePath)
	if err != nil {
		return fmt.Errorf("failed to create metadata file: %w", err)
	}
	defer f.Close()

	data, err := ExportMetadata(router)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func ExportMetadata(router *Router) (map[string]map[uint8][]ExportedEndpoint, error) {
	registries := router.GetRegistries()

	result := make(map[string]map[uint8][]ExportedEndpoint)

	for domain, versions := range registries {
		// ドメインが存在しなかったら初期化
		if _, exists := result[domain]; !exists {
			result[domain] = make(map[uint8][]ExportedEndpoint)
		}

		for version, kinds := range versions {
			var list []ExportedEndpoint

			for _, eps := range kinds {
				for _, ep := range eps {
					list = append(list, ExportedEndpoint{
						Kind:             kindToString(ep.Kind),
						Domain:           ep.Domain,
						Version:          ep.Version,
						MethodName:       ep.MethodName,
						HTTPMethod:       ep.HTTPMethod,
						RequestType:      ep.RequestType,
						ResponseType:     ep.ResponseType,
						Summary:          ep.Summary,
						Description:      ep.Description,
						Tags:             ep.Tags,
						RequestTypeInfo:  ep.RequestTypeInfo,
						ResponseTypeInfo: ep.ResponseTypeInfo,
					})
				}
			}
			result[domain][version] = list
		}
	}

	return result, nil
}

type ExportMetadataJSONOptions func(*exportMetadataJSONOptions)

type exportMetadataJSONOptions struct {
	filePath string
}

var defaultExportMetadataJSONOptions = exportMetadataJSONOptions{
	filePath: ".api/metadata.json",
}

func WithFilePath(filePath string) ExportMetadataJSONOptions {
	return func(opts *exportMetadataJSONOptions) {
		opts.filePath = filePath
	}
}

func kindToString(kind EndpointKind) string {
	return kind.String()
}
