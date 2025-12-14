package handler

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kinpatsu-everyone/backend-template/config"
	"github.com/kinpatsu-everyone/backend-template/enum"
	"github.com/kinpatsu-everyone/backend-template/internal/gcs"
	"github.com/kinpatsu-everyone/backend-template/internal/gemini"
	"github.com/kinpatsu-everyone/backend-template/internal/mysql"
	"github.com/kinpatsu-everyone/backend-template/pkg/outologger"
)

// CreateMonsterRequest はMonster登録リクエストです
type CreateMonsterRequest struct {
	Nickname  string                `multipart:"nickname"`  // ニックネーム
	Latitude  float64               `multipart:"latitude"`  // 緯度(-90.0 ~ 90.0)
	Longitude float64               `multipart:"longitude"` // 経度(-180.0 ~ 180.0)
	Image     *multipart.FileHeader `multipart:"image"`     // 画像ファイル
}

// Validate はリクエストのバリデーションを行います
func (r CreateMonsterRequest) Validate() error {
	if r.Nickname == "" {
		return fmt.Errorf("nickname is required")
	}
	if r.Image == nil {
		return fmt.Errorf("image is required")
	}
	// TODO: 緯度・経度の範囲チェック（-90.0 ~ 90.0, -180.0 ~ 180.0）
	return nil
}

// CreateMonsterResponse はMonster登録レスポンスです
type CreateMonsterResponse struct {
	MonsterID         string `json:"monsterid"`           // モンスターID(UUID)
	TrashType         string `json:"trash_type"`          // ごみ種判別結果
	GeneratedImageURL string `json:"generated_image_url"` // 生成されたモンスター画像のGCS URL
	OriginalImageURL  string `json:"original_image_url"`  // 元のごみ箱画像のGCS URL
}

// CreateMonster はMonster登録ハンドラーです
// 処理内容:
// 1. Monsterの永続化（ニックネーム、緯度、経度を保存）
// 2. AnalyzeAndGenerateImageMultipartの画像分析・画像生成処理を実行
// 3. 生成された画像のURLをMonsterのImageurlに保存
// 4. MonsterのPKをレスポンスで返す
func CreateMonster(ctx context.Context, req *CreateMonsterRequest) (*CreateMonsterResponse, error) {
	logger := outologger.GetLogger()
	queries := mysql.GetQueries()

	// 1. Monsterの永続化処理
	monsterID := uuid.New().String()
	monsterTrashCategoryID := uuid.New().String()

	// 画像ファイルを開く
	file, err := req.Image.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer file.Close()

	// 画像データを読み込む
	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// MIMEタイプを取得
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

	// 緯度・経度をsql.NullStringに変換
	var latitude, longitude sql.NullString
	if req.Latitude != 0 {
		latitude = sql.NullString{String: fmt.Sprintf("%f", req.Latitude), Valid: true}
	}
	if req.Longitude != 0 {
		longitude = sql.NullString{String: fmt.Sprintf("%f", req.Longitude), Valid: true}
	}

	// 初期状態で空のURLでMonsterを作成（後で更新）
	_, err = queries.CreateMonster(ctx, mysql.CreateMonsterParams{
		Monsterid:                monsterID,
		Nickname:                 req.Nickname,
		Originaltrashbinimageurl: "", // GCSアップロードは後で実装
		Generatedmonsterimageurl: "", // GCSアップロードは後で実装
		Latitude:                 latitude,
		Longitude:                longitude,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create monster: %w", err)
	}

	// 2. 画像からゴミ種別を判定
	// gemini/client.go:149 の AnalyzeTrashBinImage メソッドを利用
	analysisModel := "gemini-2.5-flash"
	analysisClient, err := gemini.NewClient(config.GeminiAPIKey, analysisModel)
	if err != nil {
		logger.Error(ctx, "failed to create analysis client", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to create analysis client: %w", err)
	}

	logger.Info(ctx, "analyzing trash bin image", map[string]any{
		"mime_type": mimeType,
		"filename":  req.Image.Filename,
		"size":      len(imageBytes),
	})

	// AnalyzeTrashBinImage を呼び出し（gemini/client.go:149）
	trashType, _, _, err := analysisClient.AnalyzeTrashBinImage(ctx, imageBytes, mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze image: %w", err)
	}

	logger.Info(ctx, "trash type determined", map[string]any{
		"trash_type": trashType,
	})

	// ゴミ種別の文字列をuint8に変換
	trashCategory := enum.StringToTrashCategoryEnum(trashType)

	// 3. ゴミ種別をMonstertrashcategoryに保存
	_, err = queries.CreateMonsterTrashCategory(ctx, mysql.CreateMonsterTrashCategoryParams{
		Monstertrashcategoryid: monsterTrashCategoryID,
		Monsterid:              monsterID,
		Trashcategory:          uint8(trashCategory),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create monster trash category: %w", err)
	}

	// 4. monsterの画像を生成
	// gemini/client.go:188 の GenerateMonsterImage メソッドを利用
	generateModel := "gemini-3-pro-image-preview"
	generateClient, err := gemini.NewClient(config.GeminiAPIKey, generateModel)
	if err != nil {
		logger.Error(ctx, "failed to create generate client", map[string]any{
			"error": err,
		})
		return nil, fmt.Errorf("failed to create generate client: %w", err)
	}

	logger.Info(ctx, "generating monster image", map[string]any{
		"model":      generateModel,
		"trash_type": trashType,
	})

	// GenerateMonsterImage を呼び出し（gemini/client.go:188）
	generatedImageData, generatedMimeType, err := generateClient.GenerateMonsterImage(ctx, trashType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate image: %w", err)
	}

	logger.Info(ctx, "monster image generated", map[string]any{
		"size":      len(generatedImageData),
		"mime_type": generatedMimeType,
	})

	// 5. 生成画像のURL保存（GCSにアップロード・公開URL方式）
	var generatedImagePath string
	var originalImagePath string

	if config.GCSBucketName != "" {
		// GCSクライアントを作成
		// 認証方法:
		// - GCS_CREDENTIALS_JSONが設定されている場合: サービスアカウントキーを使用
		// - GCS_CREDENTIALS_JSONが未設定の場合: Application Default Credentials (ADC) を使用
		//   ADCは以下の順序で認証情報を探します:
		//   1. 環境変数 GOOGLE_APPLICATION_CREDENTIALS で指定されたJSONファイル
		//   2. gcloud CLIの認証情報 (gcloud auth application-default login)
		//   3. GCE/GKEのサービスアカウント（クラウド環境の場合）
		var credentialsJSON []byte
		if config.GCSCredentialsJSON != "" {
			credentialsJSON = []byte(config.GCSCredentialsJSON)
		}

		gcsClient, err := gcs.NewClient(ctx, config.GCSBucketName, config.GCSBaseURL, credentialsJSON)
		if err != nil {
			logger.Error(ctx, "failed to create GCS client", map[string]any{
				"error": err,
			})
			// GCSクライアントの作成に失敗しても処理は続行（空文字列のまま）
		} else {
			defer gcsClient.Close()

			// 生成されたモンスター画像をアップロード（パスのみ保存）
			generatedExtension := gcs.GetExtensionFromMimeType(generatedMimeType)
			generatedObjectPath := gcs.GenerateGeneratedImagePath(monsterID, generatedExtension)
			generatedImagePath, err = gcsClient.UploadImageWithPath(ctx, generatedObjectPath, generatedImageData, generatedMimeType)
			if err != nil {
				logger.Error(ctx, "failed to upload generated image to GCS", map[string]any{
					"error":       err,
					"object_path": generatedObjectPath,
				})
				generatedImagePath = ""
			} else {
				logger.Info(ctx, "generated image uploaded to GCS", map[string]any{
					"path":        generatedImagePath,
					"object_path": generatedObjectPath,
				})
			}

			// 元のゴミ箱画像をアップロード（パスのみ保存）
			originalExtension := gcs.GetExtensionFromMimeType(mimeType)
			originalObjectPath := gcs.GenerateOriginalImagePath(monsterID, originalExtension)
			originalImagePath, err = gcsClient.UploadImageWithPath(ctx, originalObjectPath, imageBytes, mimeType)
			if err != nil {
				logger.Error(ctx, "failed to upload original image to GCS", map[string]any{
					"error":       err,
					"object_path": originalObjectPath,
				})
				originalImagePath = ""
			} else {
				logger.Info(ctx, "original image uploaded to GCS", map[string]any{
					"path":        originalImagePath,
					"object_path": originalObjectPath,
				})
			}
		}
	} else {
		logger.Info(ctx, "GCS bucket name not configured, skipping image upload", nil)
	}

	// Monsterの画像パスを更新
	_, err = queries.UpdateMonster(ctx, mysql.UpdateMonsterParams{
		Nickname:                 req.Nickname,
		Originaltrashbinimageurl: originalImagePath,
		Generatedmonsterimageurl: generatedImagePath,
		Latitude:                 latitude,
		Longitude:                longitude,
		Monsterid:                monsterID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update monster with image paths: %w", err)
	}

	// 6. MonsterのPKとパスを返す（クライアントには署名付きURLは返さない）
	return &CreateMonsterResponse{
		MonsterID:         monsterID,
		TrashType:         trashType,
		GeneratedImageURL: generatedImagePath,
		OriginalImageURL:  originalImagePath,
	}, nil
}

// GetMonstersRequest はMonster一覧取得リクエストです
type GetMonstersRequest struct{}

// Validate はリクエストのバリデーションを行います
func (r GetMonstersRequest) Validate() error {
	return nil
}

// MonsterItem はMonster一覧の各アイテムです
type MonsterItem struct {
	ID            string  `json:"id"`             // モンスターID(UUID)
	Nickname      string  `json:"nickname"`       // ニックネーム
	Latitude      float64 `json:"latitude"`       // 緯度(-90.0 ~ 90.0, nullable)
	Longitude     float64 `json:"longitude"`      // 経度(-180.0 ~ 180.0, nullable)
	TrashCategory string  `json:"trash_category"` // ゴミ種別("指定なし", "燃えるゴミ", "不燃ごみ", "缶", "瓶", "ペットボトル")
	ImageURL      string  `json:"image_url"`      // 画像のURL (https://images.kinpatsu.fanlav.net/monsters/{uuid}/model.png)
}

// GetMonstersResponse はMonster一覧取得レスポンスです
type GetMonstersResponse struct {
	Monsters []MonsterItem `json:"monsters"` // Monsterの配列
}

// GetMonsters はMonster一覧取得ハンドラーです
// 処理内容:
// 1. データベースからMonster一覧を取得
// 2. 各Monsterの画像URLを生成 (https://images.kinpatsu.fanlav.net/monsters/{uuid}/model.png)
// 3. 各Monsterの分類種別を取得
// 4. レスポンスとして配列を返す
func GetMonsters(ctx context.Context, _ *GetMonstersRequest) (*GetMonstersResponse, error) {
	// TODO: 実装予定
	// 1. データベースからMonster一覧を取得
	// 2. 各Monsterの画像URLを生成 (https://images.kinpatsu.fanlav.net/monsters/{uuid}/model.png)
	// 3. 各Monsterの分類種別を取得（Monstertrashcategoryテーブルから）
	// 4. レスポンスとして配列を返す

	// 現在はmockデータを返す
	return &GetMonstersResponse{
		Monsters: []MonsterItem{
			{
				ID:            "00000000-0000-0000-0000-000000000001",
				Nickname:      "燃えるゴミモンスター",
				Latitude:      35.6812,
				Longitude:     139.7671,
				TrashCategory: mysql.TrashCategoryToString(1), // 燃えるゴミ
				ImageURL:      "https://images.kinpatsu.fanlav.net/monsters/00000000-0000-0000-0000-000000000001/model.png",
			},
			{
				ID:            "00000000-0000-0000-0000-000000000002",
				Nickname:      "缶モンスター",
				Latitude:      35.6823,
				Longitude:     139.7682,
				TrashCategory: mysql.TrashCategoryToString(3), // 缶
				ImageURL:      "https://images.kinpatsu.fanlav.net/monsters/00000000-0000-0000-0000-000000000002/model.png",
			},
			{
				ID:            "00000000-0000-0000-0000-000000000003",
				Nickname:      "ペットボトルモンスター",
				Latitude:      0,
				Longitude:     0,
				TrashCategory: mysql.TrashCategoryToString(5), // ペットボトル
				ImageURL:      "https://images.kinpatsu.fanlav.net/monsters/00000000-0000-0000-0000-000000000003/model.png",
			},
		},
	}, nil
}

// GetMonsterRequest はMonster一件取得リクエストです
type GetMonsterRequest struct {
	ID string `json:"id"` // モンスターID(UUID)
}

// Validate はリクエストのバリデーションを行います
func (r GetMonsterRequest) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("id is required")
	}
	return nil
}

// GetMonsterResponse はMonster一件取得レスポンスです
type GetMonsterResponse struct {
	Monster           MonsterItem `json:"monster"`             // Monster情報
	OriginalImageURL  string      `json:"original_image_url"`  // 元のゴミ箱画像の署名付きURL
	GeneratedImageURL string      `json:"generated_image_url"` // 生成されたモンスター画像の署名付きURL
}

// GetMonster はMonster一件取得ハンドラーです
// 処理内容:
// 1. パスパラメータからMonsterのPKを取得
// 2. データベースからMonsterを取得
// 3. 保存されたGCSパスから署名付きURLを生成
// 4. Monsterの分類種別を取得
// 5. レスポンスとして返す
func GetMonster(ctx context.Context, req *GetMonsterRequest) (*GetMonsterResponse, error) {
	logger := outologger.GetLogger()
	queries := mysql.GetQueries()

	// 1. データベースからMonsterを取得
	monster, err := queries.GetMonster(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get monster: %w", err)
	}

	// 2. ゴミ種別を取得（MonsterTrashCategoryテーブルから）
	trashCategories, err := queries.ListMonsterTrashCategories(ctx, req.ID)
	if err != nil {
		logger.Error(ctx, "failed to get trash categories", map[string]any{
			"error":      err,
			"monster_id": req.ID,
		})
		// エラーがあっても続行（デフォルト値を使用）
	}

	trashCategory := "指定なし"
	if len(trashCategories) > 0 {
		trashCategory = mysql.TrashCategoryToString(trashCategories[0].Trashcategory)
	}

	// 3. 保存されているパスから署名付きURLを生成
	var generatedImageURL string
	var originalImageURL string

	if config.GCSBucketName != "" && (monster.Generatedmonsterimageurl != "" || monster.Originaltrashbinimageurl != "") {
		var credentialsJSON []byte
		if config.GCSCredentialsJSON != "" {
			credentialsJSON = []byte(config.GCSCredentialsJSON)
		}

		gcsClient, err := gcs.NewClient(ctx, config.GCSBucketName, config.GCSBaseURL, credentialsJSON)
		if err != nil {
			logger.Error(ctx, "failed to create GCS client", map[string]any{
				"error": err,
			})
		} else {
			defer gcsClient.Close()

			// 生成画像の署名付きURLを取得
			if monster.Generatedmonsterimageurl != "" {
				generatedImageURL, err = gcsClient.GetSignedURL(ctx, monster.Generatedmonsterimageurl, 24*time.Hour)
				if err != nil {
					logger.Error(ctx, "failed to generate signed URL for generated image", map[string]any{
						"error": err,
						"path":  monster.Generatedmonsterimageurl,
					})
					generatedImageURL = ""
				} else {
					logger.Info(ctx, "generated signed URL for generated image", map[string]any{
						"path": monster.Generatedmonsterimageurl,
					})
				}
			}

			// 元画像の署名付きURLを取得
			if monster.Originaltrashbinimageurl != "" {
				originalImageURL, err = gcsClient.GetSignedURL(ctx, monster.Originaltrashbinimageurl, 24*time.Hour)
				if err != nil {
					logger.Error(ctx, "failed to generate signed URL for original image", map[string]any{
						"error": err,
						"path":  monster.Originaltrashbinimageurl,
					})
					originalImageURL = ""
				} else {
					logger.Info(ctx, "generated signed URL for original image", map[string]any{
						"path": monster.Originaltrashbinimageurl,
					})
				}
			}
		}
	}

	// 4. レスポンスを返す
	var latitude, longitude float64
	if monster.Latitude.Valid {
		latitude, _ = strconv.ParseFloat(monster.Latitude.String, 64)
	}
	if monster.Longitude.Valid {
		longitude, _ = strconv.ParseFloat(monster.Longitude.String, 64)
	}

	return &GetMonsterResponse{
		Monster: MonsterItem{
			ID:            monster.Monsterid,
			Nickname:      monster.Nickname,
			Latitude:      latitude,
			Longitude:     longitude,
			TrashCategory: trashCategory,
			ImageURL:      generatedImageURL, // 生成画像の署名付きURL
		},
		OriginalImageURL:  originalImageURL,  // 元画像の署名付きURL
		GeneratedImageURL: generatedImageURL, // 生成画像の署名付きURL
	}, nil
}
