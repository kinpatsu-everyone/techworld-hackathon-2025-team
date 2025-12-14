package gcs

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// Client はGCSクライアントです
type Client struct {
	client     *storage.Client
	bucketName string
	baseURL    string
}

// NewClient は新しいGCSクライアントを作成します
// credentialsJSONが空の場合、Application Default Credentials (ADC) を使用します
// ADCは以下の順序で認証情報を探します:
//  1. 環境変数 GOOGLE_APPLICATION_CREDENTIALS で指定されたJSONファイル
//  2. gcloud CLIの認証情報 (gcloud auth application-default login)
//  3. GCE/GKEのサービスアカウント（クラウド環境の場合）
func NewClient(ctx context.Context, bucketName, baseURL string, credentialsJSON []byte) (*Client, error) {
	var opts []option.ClientOption
	if len(credentialsJSON) > 0 {
		opts = append(opts, option.WithCredentialsJSON(credentialsJSON))
	}
	// credentialsJSONが空の場合、optsも空になり、storage.NewClientは自動的にADCを使用します

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &Client{
		client:     client,
		bucketName: bucketName,
		baseURL:    baseURL,
	}, nil
}

// UploadImage は画像データをGCSにアップロードし、URLを返します
// objectPath: GCS内のオブジェクトパス（例: "monsters/{uuid}/original.jpg"）
// imageData: アップロードする画像データ
// mimeType: 画像のMIMEタイプ（例: "image/jpeg", "image/png"）
// makePublic: trueの場合、オブジェクトを公開読み取り可能にする（falseの場合は署名付きURLを使用）
// 戻り値: 公開URLまたは署名付きURL
//
// 署名付きURLについて:
// - ブラウザやフロントエンドから直接表示可能（<img src="署名付きURL">で表示できる）
// - URLを知っている人なら誰でもアクセス可能（URLが漏洩すると誰でも見られる）
// - 有効期限がある（デフォルト: 1年間）
// - バケットを公開しなくても動作する（セキュリティ上推奨）
// - サービスアカウントの秘密鍵が必要（GCSCredentialsJSONに含まれる）
func (c *Client) UploadImage(ctx context.Context, objectPath string, imageData []byte, mimeType string, makePublic bool) (string, error) {
	if c.bucketName == "" {
		return "", fmt.Errorf("bucket name is required")
	}

	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(objectPath)

	// オブジェクトの書き込み
	writer := obj.NewWriter(ctx)
	writer.ContentType = mimeType
	writer.CacheControl = "public, max-age=31536000" // 1年間キャッシュ

	// 画像データを書き込む
	if _, err := writer.Write(imageData); err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to write image data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 公開アクセスを設定する場合
	if makePublic {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return "", fmt.Errorf("failed to set public ACL: %w", err)
		}

		// 公開URLを生成
		var publicURL string
		if c.baseURL != "" {
			// カスタムベースURLが指定されている場合
			publicURL = fmt.Sprintf("%s/%s", c.baseURL, objectPath)
		} else {
			// デフォルトのGCS公開URL
			publicURL = fmt.Sprintf("https://images.kinpatsu.fanlav.net/%s", objectPath)
		}
		return publicURL, nil
	}

	// 公開アクセスを許可しない場合、署名付きURLを生成（1年間有効）
	signedURL, err := c.GetSignedURL(ctx, objectPath, 365*24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return signedURL, nil
}

// UploadImageFromReader はio.Readerから画像データを読み込んでGCSにアップロードします
// objectPath: GCS内のオブジェクトパス
// reader: 画像データを読み込むReader
// mimeType: 画像のMIMEタイプ
// makePublic: trueの場合、オブジェクトを公開読み取り可能にする（falseの場合は署名付きURLを使用）
// 戻り値: 公開URLまたは署名付きURL
func (c *Client) UploadImageFromReader(ctx context.Context, objectPath string, reader io.Reader, mimeType string, makePublic bool) (string, error) {
	if c.bucketName == "" {
		return "", fmt.Errorf("bucket name is required")
	}

	bucket := c.client.Bucket(c.bucketName)
	obj := bucket.Object(objectPath)

	// オブジェクトの書き込み
	writer := obj.NewWriter(ctx)
	writer.ContentType = mimeType
	writer.CacheControl = "public, max-age=31536000" // 1年間キャッシュ

	// 画像データをコピー
	if _, err := io.Copy(writer, reader); err != nil {
		writer.Close()
		return "", fmt.Errorf("failed to copy image data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// 公開アクセスを設定する場合
	if makePublic {
		if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
			return "", fmt.Errorf("failed to set public ACL: %w", err)
		}

		// 公開URLを生成
		var publicURL string
		if c.baseURL != "" {
			// カスタムベースURLが指定されている場合
			publicURL = fmt.Sprintf("%s/%s", c.baseURL, objectPath)
		} else {
			// デフォルトのGCS公開URL
			publicURL = fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.bucketName, objectPath)
		}
		return publicURL, nil
	}

	// 公開アクセスを許可しない場合、署名付きURLを生成（1年間有効）
	signedURL, err := c.GetSignedURL(ctx, objectPath, 365*24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return signedURL, nil
}

// GenerateObjectPath はモンスターIDと画像タイプからGCSオブジェクトパスを生成します
// monsterID: モンスターID（UUID）
// imageType: 画像タイプ（"original" または "generated"）
// extension: ファイル拡張子（例: "jpg", "png"）
// 戻り値: GCSオブジェクトパス
func GenerateObjectPath(monsterID, imageType, extension string) string {
	// 拡張子にドットがない場合は追加
	if extension != "" && extension[0] != '.' {
		extension = "." + extension
	}
	return filepath.Join("monsters", monsterID, fmt.Sprintf("%s%s", imageType, extension))
}

// GetExtensionFromMimeType はMIMEタイプからファイル拡張子を取得します
func GetExtensionFromMimeType(mimeType string) string {
	switch mimeType {
	case "image/jpeg", "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	case "image/webp":
		return "webp"
	default:
		return "jpg" // デフォルト
	}
}

// Close はGCSクライアントを閉じます
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// GetSignedURL は署名付きURLを生成します（一時的なアクセス用）
// objectPath: GCS内のオブジェクトパス
// expiration: URLの有効期限
// 戻り値: 署名付きURL
//
// 署名付きURLの特徴:
// - ブラウザやフロントエンドから直接表示可能（<img src="...">で表示できる）
// - URLを知っている人なら誰でもアクセス可能
// - 有効期限がある（期限切れ後はアクセス不可）
// - バケットを公開しなくても動作する
// - サービスアカウントの秘密鍵で署名される（GCSCredentialsJSONに含まれる必要がある）
func (c *Client) GetSignedURL(ctx context.Context, objectPath string, expiration time.Duration) (string, error) {
	if c.bucketName == "" {
		return "", fmt.Errorf("bucket name is required")
	}

	bucket := c.client.Bucket(c.bucketName)

	opts := &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(expiration),
	}

	url, err := bucket.SignedURL(objectPath, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	return url, nil
}
