package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kinpatsu-everyone/backend-template/config"
	"github.com/kinpatsu-everyone/backend-template/internal/mysql"
	"github.com/kinpatsu-everyone/backend-template/pkg/outologger"
	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter"
	"github.com/kinpatsu-everyone/backend-template/router"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	run(ctx)
}

func run(ctx context.Context) {
	config.LoadEnv(ctx)

	// --- Outorouter と Logger の設定 ---
	// Development モード用の設定
	dev := outorouter.NewDevConfig(
		outorouter.WithMetadataFilePath(".api/metadata.json"),
	)

	// Logger の設定
	logger := outologger.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	outologger.SetLogger(logger)

	// MySQLの設定
	if err := mysql.InitDB(ctx, mysql.DefaultPoolConfig()); err != nil {
		logger.Error(ctx, "❌MySQLの起動に失敗しました", map[string]any{
			"error": err,
		})
	}

	logger.Info(ctx, "✅MySQLの起動に成功しました", map[string]any{
		"stat": mysql.GetDB().Stats(),
	})

	// ルーターの設定
	r := outorouter.New(
		outorouter.WithLogger(outologger.GetLogger()),
	)

	// CORS設定
	corsConfig := outorouter.DefaultCORSConfig()
	corsConfig.AllowedOrigins = config.CORSAllowedOrigins

	// ミドルウェアの登録（適用順序が重要）
	r.Use(
		outorouter.CORSMiddleware(corsConfig), // 1. CORS処理（最初に実行）
		outorouter.NowUTCMiddleware(),         // 2. リクエスト時刻を記録
		outorouter.RequestIDMiddleware(),      // 3. リクエストIDを生成
		outorouter.LoggingMiddleware(logger),  // 4. アクセスログとパニックリカバリー
	)

	// 起動ログ
	logger.Info(ctx, "Starting server", map[string]any{
		"environment": config.ENV,
		"port":        config.ApiPort,
		"development": config.IsDevelopment(),
	})

	mux, err := router.Build(r)
	if err != nil {
		logger.Error(ctx, "failed to build router", map[string]any{
			"error": err,
		})
		return
	}

	// Development モードの場合、メタデータをエクスポートする
	if config.IsLocal() {
		if err := dev.Run(r); err != nil {
			logger.Error(ctx, "failed to export metadata", map[string]any{
				"error": err,
			})
		}
	}

	// ---- HTTP サーバーの起動とシャットダウン処理 ----
	// HTTP Server の作成
	addr := fmt.Sprintf(":%s", config.ApiPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 120 * time.Second, // 大きな画像データを返す場合に備えて延長
		IdleTimeout:  60 * time.Second,
	}

	// サーバーをゴルーチンで起動
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info(ctx, "Server listening", map[string]any{
			"addr": server.Addr,
		})
		serverErrors <- server.ListenAndServe()
	}()

	// シグナル待機
	select {
	case err := <-serverErrors:
		// サーバー起動エラー
		logger.Error(ctx, "Server error", map[string]any{
			"error": err,
		})
		return
	case <-ctx.Done():
		// シャットダウンシグナル受信
		logger.Info(ctx, "Shutting down server", map[string]any{
			"signal": "SIGINT or SIGTERM",
		})

		// グレースフルシャットダウンの実行
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error(ctx, "Failed to gracefully shutdown server", map[string]any{
				"error": err,
			})
			if err := server.Close(); err != nil {
				logger.Error(ctx, "Failed to close server", map[string]any{
					"error": err,
				})
			}
			return
		}

		logger.Info(ctx, "Server shutdown complete", map[string]any{})
	}
}
