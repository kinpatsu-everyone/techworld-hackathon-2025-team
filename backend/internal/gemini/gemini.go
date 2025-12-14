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
const GenerateTrashMonsterPromptTemplate = `[アップロードされたゴミ箱の画像]をベースに、その形や色を活かしたモンスターのイラストを生成してください。
**1. モンスター基本情報**

* **ゴミの分類**:

    %s

**2. 属性と特徴**

* **テーマとなる動物**: [以下の項目に当てはまる動物を選択してください]
動物園や水族館にいる動物をランダムに選択。

| ゴミの分類 | 属性 | モンスターの具体的な特徴 |
| :--- | :--- | :--- |
| **燃えるゴミ** | 炎属性 | primary color: 赤系の色、secondary color: ゴミ箱の色 |
| **プラ・その他** | 闇属性 | primary color: 紫系の色、secondary color: ゴミ箱の色 |
| **ペットボトル** | 水属性 | primary color: 青系の色、secondary color: ゴミ箱の色 |
| **びん・缶** | 光属性 | primary color: 黄色系の色、secondary color: ゴミ箱の色 |
| **新聞紙** | 木属性 | primary color: 緑系の色、secondary color: ゴミ箱の色 |

色は[以下のいずれかを選択して入力してください]の特徴を持つ質感にする
- ちょっとグラデーションあり
- ちょっと光沢あり
- マット系の色


**3. ポーズと背景**

* **ポーズ**: アップロードされた**ゴミ箱の形状を生かした**形で、キャラがポーズをとるのも良さそうです。お任せします。例として、ゴミ箱の中に丸まっていたり、周りに巻き付いたりする構図を想定しています。
* **背景**: 背景は、入力写真そのものにして、キャラを追記する形にしてください。

**4. スタイル指定（イラストレーション）**

* **スタイル**: 可愛すぎず、リアルで怖すぎない、コミックぽい感じ。
* **枠線**: **非常に太く、明確な枠線（太いアウトライン）**を使用し、キャラクターを際立たせること。
* **着色**: 明るくクリアな、シンプルな**セル画塗り**（Cell Shading）。複雑なテクスチャやグラデーションは避ける。
* **構図**: 正方形のイラスト。(元画像が長方形ならゴミ箱が写るようにクロップする)

<!-- **5. 出力形式**

* 生成されたイラストの下部に、以下のテキストを明確に表示すること。

    * "モンスター名: [ここに記入したモンスター名]"
    * "ゴミの分類: [ここに記入したゴミの分類]" -->

---
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
