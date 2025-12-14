package gemini

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// AnalyzeTrashBinPrompt はゴミ箱の写真から分別種類を判定するためのプロンプトです
const AnalyzeTrashBinPrompt = `この画像に写っているゴミ箱を分析してください。
ゴミ箱の色、マーク、中身、その他の特徴を観察し、このゴミ箱がどの分別種類（燃えるゴミ、不燃ごみ、缶、瓶、ペットボトルなど）に該当するかを判定してください。

以下のJSON形式で回答してください：
{
  "trash_type": "分別種類（例: 燃えるゴミ、不燃ごみ、缶、瓶、ペットボトルなど）",
  "color": "ゴミ箱の主な色",
  "markings": "ゴミ箱に書かれているマークや文字",
  "contents": "ゴミ箱の中身の説明（見える範囲で）",
  "reasoning": "判定理由の説明",
  "description": "ゴミ箱の全体的な説明",
  "confidence": "判定の確信度（high, medium, low）"
}

判定が難しい場合は、"trash_type"を"unknown"としてください。`

// GenerateTrashMonsterPromptTemplate は分別種をテーマにしたモンスターキャラクター生成用のプロンプトテンプレートです
// プレースホルダー: %s = trashType (2箇所)
const GenerateTrashMonsterPromptTemplate = `この写真のゴミ箱について、景観や背景はそのままにリアルなモンスターに生まれ変わらせてください。
以下の条件を守ってください：
・日本文化をモチーフにした動物や妖怪・幻獣など
・分別種「%s」をテーマにしたデザイン
`

// TrashAnalysisResult は画像分析結果のJSON構造です
type TrashAnalysisResult struct {
	TrashType   string `json:"trash_type"`
	Color       string `json:"color"`
	Markings    string `json:"markings"`
	Contents    string `json:"contents"`
	Reasoning   string `json:"reasoning"`
	Description string `json:"description"`
	Confidence  string `json:"confidence"`
}

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

// AnalyzeTrashBinImage はゴミ箱の画像を分析して分別種類を判定します
// 戻り値: 分別種類（trash_type）、分析結果のテキスト、解析されたJSON構造
func (c *Client) AnalyzeTrashBinImage(ctx context.Context, imageData []byte, mimeType string) (trashType string, analysisText string, result TrashAnalysisResult, err error) {
	// 画像分析を実行
	analysisResp, err := c.AnalyzeImage(ctx, AnalyzeTrashBinPrompt, imageData, mimeType)
	if err != nil {
		return "", "", TrashAnalysisResult{}, fmt.Errorf("failed to analyze image: %w", err)
	}

	// 分析結果のテキストを抽出
	analysisText = ""
	if len(analysisResp.Candidates) > 0 && analysisResp.Candidates[0].Content != nil {
		for _, part := range analysisResp.Candidates[0].Content.Parts {
			if part.Text != "" {
				analysisText += part.Text
			}
		}
	}

	// 分析結果から分別種を抽出
	trashType = "unknown"
	var analysisResult TrashAnalysisResult

	// JSONをパース（テキスト内にJSONが含まれている可能性があるため、抽出を試みる）
	jsonStart := strings.Index(analysisText, "{")
	jsonEnd := strings.LastIndex(analysisText, "}")
	if jsonStart >= 0 && jsonEnd > jsonStart {
		jsonStr := analysisText[jsonStart : jsonEnd+1]
		if err := json.Unmarshal([]byte(jsonStr), &analysisResult); err == nil {
			trashType = analysisResult.TrashType
			if trashType == "" {
				trashType = "unknown"
			}
		}
	}

	return trashType, analysisText, analysisResult, nil
}

// GenerateMonsterImage は分別種をテーマにしたモンスターキャラクターの画像を生成します
// 戻り値: 生成された画像データ（バイナリ）、MIMEタイプ
func (c *Client) GenerateMonsterImage(ctx context.Context, trashType string) (imageData []byte, mimeType string, err error) {
	// 分別種をテーマにした画像生成プロンプトを作成
	generatePrompt := fmt.Sprintf(GenerateTrashMonsterPromptTemplate, trashType)

	// 画像生成を実行
	generateResp, err := c.GenerateContent(ctx, generatePrompt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate image: %w", err)
	}

	// 生成画像からバイナリデータを抽出
	imageData = nil
	mimeType = "image/png" // デフォルト

	for _, cand := range generateResp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				// インライン画像データを取得
				if part.InlineData != nil {
					imageData = part.InlineData.Data
					mimeType = part.InlineData.MIMEType
					break
				}
			}
			if len(imageData) > 0 {
				break
			}
		}
	}

	if len(imageData) == 0 {
		return nil, "", fmt.Errorf("generated image data not found")
	}

	return imageData, mimeType, nil
}

// AnalyzeAndGenerateMonsterImage はゴミ箱の画像を分析し、分別種をテーマにしたモンスターキャラクターの画像を生成します
// generateClient: 画像生成用のクライアント（分析用とは別のモデルを使用する場合）
// 戻り値: 分別種類、分析結果、生成された画像データ（バイナリ）、MIMEタイプ
func (c *Client) AnalyzeAndGenerateMonsterImage(ctx context.Context, imageData []byte, mimeType string, generateClient *Client) (string, TrashAnalysisResult, []byte, string, error) {
	// Step 1: 画像分析（現在のクライアントを使用）
	trashType, _, analysisResult, err := c.AnalyzeTrashBinImage(ctx, imageData, mimeType)
	if err != nil {
		return "", TrashAnalysisResult{}, nil, "", fmt.Errorf("failed to analyze image: %w", err)
	}

	// Step 2: 画像生成（指定されたクライアントを使用、nilの場合は現在のクライアントを使用）
	genClient := generateClient
	if genClient == nil {
		genClient = c
	}

	generatedImageData, generatedMimeType, err := genClient.GenerateMonsterImage(ctx, trashType)
	if err != nil {
		return trashType, analysisResult, nil, "", fmt.Errorf("failed to generate image: %w", err)
	}

	return trashType, analysisResult, generatedImageData, generatedMimeType, nil
}

// DecodeBase64Image はbase64エンコードされた画像データをデコードします
func DecodeBase64Image(base64Data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(base64Data)
}

// EncodeImageToBase64 は画像データをbase64エンコードします
func EncodeImageToBase64(imageData []byte) string {
	return base64.StdEncoding.EncodeToString(imageData)
}
