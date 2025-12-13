package gemini

import (
	"context"
	"fmt"

	"google.golang.org/genai"
)

// Client はGemini APIクライアントです
type Client struct {
	client *genai.Client
	model  string
}

// NewClient は新しいGemini APIクライアントを作成します
func NewClient(apiKey, model string) (*Client, error) {
	if model == "" {
		model = "gemini-3-pro-image-preview"
	}

	ctx := context.Background()
	config := &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	}

	client, err := genai.NewClient(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Client{
		client: client,
		model:  model,
	}, nil
}

// GenerateContent は画像生成リクエストを送信します
func (c *Client) GenerateContent(ctx context.Context, prompt string) (*genai.GenerateContentResponse, error) {
	contents := genai.Text(prompt)

	result, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		contents,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	return result, nil
}

// AnalyzeImage は画像データとテキストプロンプトを送信して、分析結果をテキストで返します
func (c *Client) AnalyzeImage(ctx context.Context, prompt string, imageData []byte, mimeType string) (*genai.GenerateContentResponse, error) {
	// 画像データが空の場合はエラー
	if len(imageData) == 0 {
		return nil, fmt.Errorf("image data is required")
	}

	// MIMEタイプが指定されていない場合はデフォルトを設定
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	// テキストと画像データを含むPartsを作成
	parts := []*genai.Part{
		{Text: prompt},
	}

	// 画像データを追加
	parts = append(parts, &genai.Part{
		InlineData: &genai.Blob{
			Data:     imageData,
			MIMEType: mimeType,
		},
	})

	// Contentを作成
	contents := []*genai.Content{
		{
			Parts: parts,
			Role:  genai.RoleUser,
		},
	}

	result, err := c.client.Models.GenerateContent(
		ctx,
		c.model,
		contents,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	return result, nil
}
