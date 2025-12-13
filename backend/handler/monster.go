package handler

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/kinpatsu-everyone/backend-template/internal/mysql"
)

// CreateMonsterRequest はMonster登録リクエストです
type CreateMonsterRequest struct {
	Nickname  string                `multipart:"nickname"`  // ニックネーム
	Latitude  string                `multipart:"latitude"`  // 緯度(-90.0 ~ 90.0)
	Longitude string                `multipart:"longitude"` // 経度(-180.0 ~ 180.0)
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
	MonsterID string `json:"monsterid"` // モンスターID(UUID)
}

// CreateMonster はMonster登録ハンドラーです
// 処理内容:
// 1. Monsterの永続化（ニックネーム、緯度、経度を保存）
// 2. AnalyzeAndGenerateImageMultipartの画像分析・画像生成処理を実行
// 3. 生成された画像のURLをMonsterのImageurlに保存
// 4. MonsterのPKをレスポンスで返す
func CreateMonster(ctx context.Context, req *CreateMonsterRequest) (*CreateMonsterResponse, error) {
	// TODO: 実装予定
	// 1. Monsterの永続化処理
	// 2. AnalyzeAndGenerateImageMultipartの呼び出し
	// 3. 生成画像のURL保存
	// 4. MonsterのPKを返す

	// 現在はmockデータを返す
	return &CreateMonsterResponse{
		MonsterID: "00000000-0000-0000-0000-000000000001", // mock UUID
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
	ID            string `json:"id"`             // モンスターID(UUID)
	Nickname      string `json:"nickname"`       // ニックネーム
	Latitude      string `json:"latitude"`       // 緯度(-90.0 ~ 90.0, nullable)
	Longitude     string `json:"longitude"`      // 経度(-180.0 ~ 180.0, nullable)
	TrashCategory string `json:"trash_category"` // ゴミ種別("指定なし", "燃えるゴミ", "不燃ごみ", "缶", "瓶", "ペットボトル")
	ImageURL      string `json:"image_url"`      // 画像のURL (https://images.kinpatsu.fanlav.net/monsters/{uuid}/model.png)
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
				Latitude:      "35.6812",
				Longitude:     "139.7671",
				TrashCategory: mysql.TrashCategoryToString(1), // 燃えるゴミ
				ImageURL:      "https://images.kinpatsu.fanlav.net/monsters/00000000-0000-0000-0000-000000000001/model.png",
			},
			{
				ID:            "00000000-0000-0000-0000-000000000002",
				Nickname:      "缶モンスター",
				Latitude:      "35.6823",
				Longitude:     "139.7682",
				TrashCategory: mysql.TrashCategoryToString(3), // 缶
				ImageURL:      "https://images.kinpatsu.fanlav.net/monsters/00000000-0000-0000-0000-000000000002/model.png",
			},
			{
				ID:            "00000000-0000-0000-0000-000000000003",
				Nickname:      "ペットボトルモンスター",
				Latitude:      "",
				Longitude:     "",
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
	Monster MonsterItem `json:"monster"` // Monster情報
}

// GetMonster はMonster一件取得ハンドラーです
// 処理内容:
// 1. パスパラメータからMonsterのPKを取得
// 2. データベースからMonsterを取得
// 3. Monsterの画像URLを生成 (https://images.kinpatsu.fanlav.net/monsters/{uuid}/model.png)
// 4. Monsterの分類種別を取得
// 5. レスポンスとして返す
func GetMonster(ctx context.Context, req *GetMonsterRequest) (*GetMonsterResponse, error) {
	// TODO: 実装予定
	// 1. パスパラメータからMonsterのPKを取得（現在はリクエストボディから取得）
	// 2. データベースからMonsterを取得
	// 3. Monsterの画像URLを生成 (https://images.kinpatsu.fanlav.net/monsters/{uuid}/model.png)
	// 4. Monsterの分類種別を取得（Monstertrashcategoryテーブルから）
	// 5. レスポンスとして返す

	// 現在はmockデータを返す
	return &GetMonsterResponse{
		Monster: MonsterItem{
			ID:            req.ID,
			Nickname:      "燃えるゴミモンスター",
			Latitude:      "35.6812",
			Longitude:     "139.7671",
			TrashCategory: mysql.TrashCategoryToString(1), // 燃えるゴミ
			ImageURL:      fmt.Sprintf("https://images.kinpatsu.fanlav.net/monsters/%s/model.png", req.ID),
		},
	}, nil
}
