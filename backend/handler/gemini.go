package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/kinpatsu-everyone/backend-template/config"
	"github.com/kinpatsu-everyone/backend-template/pkg/gemini"
	"github.com/kinpatsu-everyone/backend-template/pkg/outologger"
	"google.golang.org/genai"
)

// GenerateImageRequest は画像生成リクエストです
type GenerateImageRequest struct {
	Prompt string `json:"prompt"`
	Model  string `json:"model,omitempty"`
}

// Validate はリクエストのバリデーションを行います
func (r GenerateImageRequest) Validate() error {
	if r.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	return nil
}

// GenerateImageResponse は画像生成レスポンスです
type GenerateImageResponse struct {
	Candidates []CandidateResponse `json:"candidates"`
}

// CandidateResponse は生成された候補のレスポンスです
type CandidateResponse struct {
	Content      ContentResponse `json:"content"`
	FinishReason string          `json:"finish_reason,omitempty"`
}

// ContentResponse はコンテンツのレスポンスです
type ContentResponse struct {
	Role  string         `json:"role"`
	Parts []PartResponse `json:"parts"`
}

// PartResponse はパートのレスポンスです
type PartResponse struct {
	Text      string     `json:"text,omitempty"`
	ImageData *ImageData `json:"image_data,omitempty"`
	FileData  *FileData  `json:"file_data,omitempty"`
}

// ImageData は画像データのレスポンスです
type ImageData struct {
	MimeType string `json:"mime_type"`
	Data     string `json:"data"` // base64エンコードされた画像データ
}

// FileData はファイルデータのレスポンスです
type FileData struct {
	MimeType string `json:"mime_type"`
	FileURI  string `json:"file_uri"`
}

// GenerateImage は画像生成テスト用ハンドラーです
func GenerateImage(ctx context.Context, req *GenerateImageRequest) (*GenerateImageResponse, error) {
	model := req.Model
	if model == "" {
		model = "gemini-3-pro-image-preview"
	}
	logger := outologger.GetLogger()
	logger.Info(ctx, "creating gemini client", map[string]any{
		"model": model,
	})
	client, err := gemini.NewClient(config.GeminiAPIKey, model)
	if err != nil {
		logger := outologger.GetLogger()
		logger.Info(ctx, "failed to create gemini client", map[string]any{
			"error": err,
		})
		// panic(err)
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	resp, err := client.GenerateContent(ctx, req.Prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// genai.GenerateContentResponseをハンドラーのレスポンス形式に変換
	candidates := make([]CandidateResponse, 0, len(resp.Candidates))
	for _, cand := range resp.Candidates {
		parts := make([]PartResponse, 0)
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				partResp := PartResponse{}

				// テキストデータ
				if part.Text != "" {
					partResp.Text = part.Text
				}

				// インライン画像データ（base64エンコード）
				if part.InlineData != nil {
					// バイトデータをbase64エンコード
					base64Data := base64.StdEncoding.EncodeToString(part.InlineData.Data)
					partResp.ImageData = &ImageData{
						MimeType: part.InlineData.MIMEType,
						Data:     base64Data,
					}
				}

				// ファイルデータ（URI）
				if part.FileData != nil {
					partResp.FileData = &FileData{
						MimeType: part.FileData.MIMEType,
						FileURI:  part.FileData.FileURI,
					}
				}

				// 何かしらのデータがある場合のみ追加
				if partResp.Text != "" || partResp.ImageData != nil || partResp.FileData != nil {
					parts = append(parts, partResp)
				}
			}
		}
		candidates = append(candidates, CandidateResponse{
			Content: ContentResponse{
				Role:  getRole(cand.Content),
				Parts: parts,
			},
			FinishReason: string(cand.FinishReason),
		})
	}

	return &GenerateImageResponse{
		Candidates: candidates,
	}, nil
}

// AnalyzeImageRequest は画像分析リクエストです
type AnalyzeImageRequest struct {
	ImageData string `json:"image_data"`          // base64エンコードされた画像データ
	MimeType  string `json:"mime_type,omitempty"` // 画像のMIMEタイプ（例: "image/jpeg", "image/png"）
	Model     string `json:"model,omitempty"`
}

// Validate はリクエストのバリデーションを行います
func (r AnalyzeImageRequest) Validate() error {
	if r.ImageData == "" {
		return fmt.Errorf("image_data is required")
	}
	return nil
}

// AnalyzeImageResponse は画像分析レスポンスです（テキストのみ）
type AnalyzeImageResponse struct {
	Text string `json:"text"`
}

// AnalyzeImage は画像分析テスト用ハンドラーです
// ゴミ箱の写真から分別種類を判定します
func AnalyzeImage(ctx context.Context, req *AnalyzeImageRequest) (*AnalyzeImageResponse, error) {
	model := req.Model
	if model == "" {
		// 画像分析には通常のGeminiモデルを使用
		model = "gemini-2.5-flash"
	}
	logger := outologger.GetLogger()
	logger.Info(ctx, "creating gemini client for image analysis", map[string]any{
		"model": model,
	})
	client, err := gemini.NewClient(config.GeminiAPIKey, model)
	if err != nil {
		logger.Error(ctx, "failed to create gemini client", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	// base64エンコードされた画像データをデコード
	imageBytes, err := base64.StdEncoding.DecodeString(req.ImageData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image data: %w", err)
	}

	// MIMEタイプが指定されていない場合はデフォルトを設定
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	// 固定プロンプト: ゴミ箱の写真から分別種類を判定（中身も含む）
	prompt := AnalyzeTrashBinPrompt

	resp, err := client.AnalyzeImage(ctx, prompt, imageBytes, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	// レスポンスからテキストを抽出
	text := ""
	if len(resp.Candidates) > 0 && resp.Candidates[0].Content != nil {
		for _, part := range resp.Candidates[0].Content.Parts {
			if part.Text != "" {
				text += part.Text
			}
		}
	}

	return &AnalyzeImageResponse{
		Text: text,
	}, nil
}

// AnalyzeAndGenerateImageRequest は画像分析と画像生成を統合したリクエストです
type AnalyzeAndGenerateImageRequest struct {
	ImageData string `json:"image_data"`          // base64エンコードされた画像データ
	MimeType  string `json:"mime_type,omitempty"` // 画像のMIMEタイプ（例: "image/jpeg", "image/png"）
	Model     string `json:"model,omitempty"`     // 画像生成用のモデル
}

// Validate はリクエストのバリデーションを行います
func (r AnalyzeAndGenerateImageRequest) Validate() error {
	if r.ImageData == "" {
		return fmt.Errorf("image_data is required")
	}
	return nil
}

// AnalyzeAndGenerateImageResponse は画像分析と画像生成を統合したレスポンスです
type AnalyzeAndGenerateImageResponse struct {
	ImageData string `json:"image_data"` // base64エンコードされた画像データ
	MimeType  string `json:"mime_type"`  // 画像のMIMEタイプ
}

// AnalyzeAndGenerateImageMultipartRequest はmultipart/form-dataで画像分析と画像生成を統合したリクエストです
type AnalyzeAndGenerateImageMultipartRequest struct {
	Image *multipart.FileHeader `multipart:"image"`                        // 画像ファイル
	Model string                `json:"model,omitempty" multipart:"model"` // 画像生成用のモデル
}

// Validate はリクエストのバリデーションを行います
func (r AnalyzeAndGenerateImageMultipartRequest) Validate() error {
	if r.Image == nil {
		return fmt.Errorf("image is required")
	}
	return nil
}

// AnalyzeAndGenerateImageMultipart はmultipart/form-dataで画像分析と画像生成を統合したハンドラーです
// ゴミ箱の写真を分析し、分別種をテーマにしたモンスターキャラクターを生成します
func AnalyzeAndGenerateImageMultipart(ctx context.Context, req *AnalyzeAndGenerateImageMultipartRequest) (*AnalyzeAndGenerateImageResponse, error) {
	logger := outologger.GetLogger()

	// Step 1: 画像ファイルを開く
	file, err := req.Image.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	// Step 2: 画像データを読み込む
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Step 3: MIMEタイプを取得（Content-Typeヘッダーから、またはファイル名から推測）
	mimeType := req.Image.Header.Get("Content-Type")
	if mimeType == "" {
		// ファイル名から推測
		filename := req.Image.Filename
		if strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
			mimeType = "image/jpeg"
		} else if strings.HasSuffix(filename, ".png") {
			mimeType = "image/png"
		} else if strings.HasSuffix(filename, ".gif") {
			mimeType = "image/gif"
		} else if strings.HasSuffix(filename, ".webp") {
			mimeType = "image/webp"
		} else {
			mimeType = "image/jpeg" // デフォルト
		}
	}

	// Step 4: 画像分析用のクライアントを作成
	analysisModel := "gemini-2.5-flash"
	analysisClient, err := gemini.NewClient(config.GeminiAPIKey, analysisModel)
	if err != nil {
		logger.Error(ctx, "failed to create analysis client", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to create analysis client: %w", err)
	}

	// Step 5: 画像分析を実行
	logger.Info(ctx, "analyzing trash bin image", map[string]any{
		"mime_type": mimeType,
		"filename":  req.Image.Filename,
		"size":      len(imageBytes),
	})
	analysisResp, err := analysisClient.AnalyzeImage(ctx, AnalyzeTrashBinPrompt, imageBytes, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	// 分析結果のテキストを抽出
	analysisText := ""
	if len(analysisResp.Candidates) > 0 && analysisResp.Candidates[0].Content != nil {
		for _, part := range analysisResp.Candidates[0].Content.Parts {
			if part.Text != "" {
				analysisText += part.Text
			}
		}
	}

	// Step 6: 分析結果から分別種を抽出
	trashType := "unknown"
	var analysisResult trashAnalysisResult

	// JSONをパース（テキスト内にJSONが含まれている可能性があるため、抽出を試みる）
	// まず、JSON部分を抽出
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

	logger.Info(ctx, "trash type determined", map[string]any{
		"trash_type": trashType,
	})

	// Step 7: 画像生成用のクライアントを作成
	generateModel := req.Model
	generateModel = "gemini-3-pro-image-preview"

	generateClient, err := gemini.NewClient(config.GeminiAPIKey, generateModel)
	if err != nil {
		logger.Error(ctx, "failed to create generate client", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to create generate client: %w", err)
	}

	// Step 8: 分別種をテーマにした画像生成プロンプトを作成
	generatePrompt := fmt.Sprintf(GenerateTrashMonsterPromptTemplate, trashType)

	logger.Info(ctx, "generating monster image", map[string]any{
		"model":      generateModel,
		"trash_type": trashType,
		"prompt":     generatePrompt,
	})

	// Step 9: 画像生成を実行
	generateResp, err := generateClient.GenerateContent(ctx, generatePrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// Step 10: 生成画像からバイナリデータを抽出
	var generatedImageData []byte
	var imageMimeType string = "image/png" // デフォルト

	for _, cand := range generateResp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				// インライン画像データを取得
				if part.InlineData != nil {
					generatedImageData = part.InlineData.Data
					imageMimeType = part.InlineData.MIMEType
					break
				}
			}
			if len(generatedImageData) > 0 {
				break
			}
		}
	}

	if len(generatedImageData) == 0 {
		return nil, fmt.Errorf("generated image data not found")
	}

	// 画像データをbase64エンコード
	base64ImageData := base64.StdEncoding.EncodeToString(generatedImageData)

	return &AnalyzeAndGenerateImageResponse{
		ImageData: base64ImageData,
		MimeType:  imageMimeType,
	}, nil
}

// trashAnalysisResult は画像分析結果のJSON構造です
type trashAnalysisResult struct {
	TrashType   string `json:"trash_type"`
	Color       string `json:"color"`
	Markings    string `json:"markings"`
	Contents    string `json:"contents"`
	Reasoning   string `json:"reasoning"`
	Description string `json:"description"`
	Confidence  string `json:"confidence"`
}

// AnalyzeAndGenerateImage は画像分析と画像生成を統合したハンドラーです
// ゴミ箱の写真を分析し、分別種をテーマにしたモンスターキャラクターを生成します
func AnalyzeAndGenerateImage(ctx context.Context, req *AnalyzeAndGenerateImageRequest) (*AnalyzeAndGenerateImageResponse, error) {
	logger := outologger.GetLogger()

	// Step 1: 画像分析用のクライアントを作成
	analysisModel := "gemini-2.5-flash"
	analysisClient, err := gemini.NewClient(config.GeminiAPIKey, analysisModel)
	if err != nil {
		logger.Error(ctx, "failed to create analysis client", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to create analysis client: %w", err)
	}

	// base64エンコードされた画像データをデコード
	imageBytes, err := base64.StdEncoding.DecodeString(req.ImageData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64 image data: %w", err)
	}

	// MIMEタイプが指定されていない場合はデフォルトを設定
	mimeType := req.MimeType
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	// Step 2: 画像分析を実行
	logger.Info(ctx, "analyzing trash bin image", map[string]any{
		"mime_type": mimeType,
	})
	analysisResp, err := analysisClient.AnalyzeImage(ctx, AnalyzeTrashBinPrompt, imageBytes, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	// 分析結果のテキストを抽出
	analysisText := ""
	if len(analysisResp.Candidates) > 0 && analysisResp.Candidates[0].Content != nil {
		for _, part := range analysisResp.Candidates[0].Content.Parts {
			if part.Text != "" {
				analysisText += part.Text
			}
		}
	}

	// Step 3: 分析結果から分別種を抽出
	trashType := "unknown"
	var analysisResult trashAnalysisResult

	// JSONをパース（テキスト内にJSONが含まれている可能性があるため、抽出を試みる）
	// まず、JSON部分を抽出
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

	logger.Info(ctx, "trash type determined", map[string]any{
		"trash_type": trashType,
	})

	// Step 4: 画像生成用のクライアントを作成
	generateModel := req.Model
	if generateModel == "" {
		generateModel = "gemini-3-pro-image-preview"
	}
	generateClient, err := gemini.NewClient(config.GeminiAPIKey, generateModel)
	if err != nil {
		logger.Error(ctx, "failed to create generate client", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to create generate client: %w", err)
	}

	// Step 5: 分別種をテーマにした画像生成プロンプトを作成
	generatePrompt := fmt.Sprintf(GenerateTrashMonsterPromptTemplate, trashType, trashType)

	logger.Info(ctx, "generating monster image", map[string]any{
		"model":      generateModel,
		"trash_type": trashType,
	})

	// Step 6: 画像生成を実行
	generateResp, err := generateClient.GenerateContent(ctx, generatePrompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	// Step 7: 生成画像からバイナリデータを抽出
	var imageData []byte
	var imageMimeType string = "image/png" // デフォルト

	for _, cand := range generateResp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				// インライン画像データを取得
				if part.InlineData != nil {
					imageData = part.InlineData.Data
					imageMimeType = part.InlineData.MIMEType
					break
				}
			}
			if len(imageData) > 0 {
				break
			}
		}
	}

	if len(imageData) == 0 {
		return nil, fmt.Errorf("generated image data not found")
	}

	// 画像データをbase64エンコード
	base64ImageData := base64.StdEncoding.EncodeToString(imageData)

	return &AnalyzeAndGenerateImageResponse{
		ImageData: base64ImageData,
		MimeType:  imageMimeType,
	}, nil
}

// getRole はContentからroleを取得します
func getRole(content *genai.Content) string {
	if content == nil {
		return ""
	}
	return string(content.Role)
}
