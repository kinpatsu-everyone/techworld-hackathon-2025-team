package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"github.com/kinpatsu-everyone/backend-template/pkg/outologger"
	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter"
)

const (
	UserTag outorouter.Tag = "User"
)

type CreateUserRequest struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

func (r CreateUserRequest) Validate() error {
	return nil
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

func CreateUser(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	return &CreateUserResponse{UserID: "testUserID"}, nil
}

// RegisterAPI は API エンドポイントを登録し、HTTP ハンドラーを返します。
func RegisterAPI(ctx context.Context) http.Handler {
	dev := outorouter.NewDevConfig(
		outorouter.WithMetadataFilePath(".api/metadata.json"),
		// TypeScriptクライアントコードを生成する
		// Expo/React Native向けの型安全なAPIクライアントを自動生成
		outorouter.WithTypeScriptClient(".api/client.ts", "http://localhost:8080"),
	)

	logger := outologger.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
	outologger.SetLogger(logger)
	r := outorouter.New(
		outorouter.WithLogger(outologger.GetLogger()),
	)

	outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[CreateUserRequest, CreateUserResponse]{
		Domain:      "user",
		Version:     1,
		MethodName:  "CreateUser",
		Summary:     "",
		Description: "",
		Tags:        outorouter.RegisterTags(UserTag),
		Handler:     CreateUser,
	})

	// Development モードの場合、メタデータをエクスポートする
	if err := dev.Run(r); err != nil {
		logger.Error(ctx, "failed to export metadata", map[string]any{
			"error": err,
		})
		return nil
	}

	return r.Handler()
}

func main() {
	ctx := context.Background()

	handler := RegisterAPI(ctx)
	logger := outologger.GetLogger()
	if logger == nil {
		logger = outologger.NewSlogLogger(slog.New(slog.NewJSONHandler(os.Stdout, nil)))
		outologger.SetLogger(logger)
	}

	if handler == nil {
		logger.Error(ctx, "failed to initialize router", nil)
		os.Exit(1)
	}

	addr := ":8080"
	logger.Info(ctx, "starting example HTTP server", map[string]any{
		"addr": addr,
	})

	if err := http.ListenAndServe(addr, handler); err != nil {
		logger.Error(ctx, "http server stopped", map[string]any{
			"error": err,
		})
	}
}
