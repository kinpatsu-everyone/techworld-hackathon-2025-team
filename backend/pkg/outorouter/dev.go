package outorouter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter/internal/generator"
	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter/internal/parser"
)

// TypeScriptClientConfig はTypeScriptクライアント生成の設定です。
type TypeScriptClientConfig struct {
	// Enabled はTypeScriptクライアント生成を有効にするかどうか
	Enabled bool
	// OutputPath は生成されるTypeScriptファイルの出力パス
	OutputPath string
	// BaseURL はAPIのベースURL
	BaseURL string
}

type DevConfig struct {
	metadataFilePath       string
	typeScriptClientConfig TypeScriptClientConfig
}

func DefaultDevConfig() *DevConfig {
	return &DevConfig{
		metadataFilePath: ".api/metadata.json",
		typeScriptClientConfig: TypeScriptClientConfig{
			Enabled:    false,
			OutputPath: ".api/client.ts",
			BaseURL:    "http://localhost:8080",
		},
	}
}

type DevConfigOption func(*DevConfig)

func WithMetadataFilePath(filePath string) DevConfigOption {
	return func(cfg *DevConfig) {
		cfg.metadataFilePath = filePath
	}
}

// WithTypeScriptClient はTypeScriptクライアント生成を有効にします。
func WithTypeScriptClient(outputPath string, baseURL string) DevConfigOption {
	return func(cfg *DevConfig) {
		cfg.typeScriptClientConfig.Enabled = true
		if outputPath != "" {
			cfg.typeScriptClientConfig.OutputPath = outputPath
		}
		if baseURL != "" {
			cfg.typeScriptClientConfig.BaseURL = baseURL
		}
	}
}

func NewDevConfig(opts ...DevConfigOption) *DevConfig {
	cfg := DefaultDevConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func (cfg *DevConfig) Run(router *Router) error {
	// メタデータJSONをエクスポート
	if err := ExportMetadataJSON(router, WithFilePath(cfg.metadataFilePath)); err != nil {
		return fmt.Errorf("failed to export metadata: %w", err)
	}

	// TypeScriptクライアント生成
	if cfg.typeScriptClientConfig.Enabled {
		if err := cfg.generateTypeScriptClient(); err != nil {
			return fmt.Errorf("failed to generate TypeScript client: %w", err)
		}
	}

	return nil
}

func (cfg *DevConfig) generateTypeScriptClient() error {
	// メタデータをパース
	meta, err := parser.ParseFile(cfg.metadataFilePath)
	if err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	// TypeScriptクライアントストラテジーでコード生成
	strategy := generator.TypeScriptClientStrategy{
		BaseURL: cfg.typeScriptClientConfig.BaseURL,
	}
	gen := generator.New(strategy)

	code, err := gen.Generate(meta)
	if err != nil {
		return fmt.Errorf("failed to generate code: %w", err)
	}

	// 出力ディレクトリを作成
	outputDir := filepath.Dir(cfg.typeScriptClientConfig.OutputPath)
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// ファイルに書き出し
	if err := os.WriteFile(cfg.typeScriptClientConfig.OutputPath, []byte(code), 0o644); err != nil {
		return fmt.Errorf("failed to write TypeScript client: %w", err)
	}

	return nil
}
