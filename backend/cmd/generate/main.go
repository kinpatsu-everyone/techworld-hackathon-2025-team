package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/kinpatsu-everyone/backend-template/pkg/outologger"
	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter"
	"github.com/kinpatsu-everyone/backend-template/router"
)

func main() {
	// フラグの定義
	outputPath := flag.String("output", ".api/client.ts", "Output path for the TypeScript client file")
	baseURL := flag.String("base-url", "http://localhost:8080", "Base URL for the API client")
	metadataPath := flag.String("metadata", ".api/metadata.json", "Output path for the metadata JSON file")
	flag.Parse()

	// Logger の設定（quietモード）
	logger := outologger.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelWarn,
	})))
	outologger.SetLogger(logger)

	// ルーターの設定
	r := outorouter.New(
		outorouter.WithLogger(logger),
	)

	// router/router.go からエンドポイントを登録
	if _, err := router.Build(r); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to build router: %v\n", err)
		os.Exit(1)
	}

	// DevConfig でコード生成
	dev := outorouter.NewDevConfig(
		outorouter.WithMetadataFilePath(*metadataPath),
		outorouter.WithTypeScriptClient(*outputPath, *baseURL),
	)

	if err := dev.Run(r); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate client: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated TypeScript client: %s\n", *outputPath)
	fmt.Printf("Generated metadata: %s\n", *metadataPath)
}
