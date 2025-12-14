package router

import (
	"net/http"

	"github.com/kinpatsu-everyone/backend-template/handler"
	"github.com/kinpatsu-everyone/backend-template/pkg/outorouter"
)

func Build(r *outorouter.Router) (http.Handler, error) {
	outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[handler.HealthzRequest, handler.HealthzResponse]{
		Domain:      "healthz",
		Version:     1,
		MethodName:  "Healthz",
		Summary:     "Health Check Endpoint",
		Description: "Returns a simple health check response.",
		Tags:        outorouter.RegisterTags("Health"),
		Handler:     handler.Healthz,
	})

	// 画像生成用単体テストエンドポイント
	outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[handler.GenerateImageRequest, handler.GenerateImageResponse]{
		Domain:      "gemini",
		Version:     1,
		MethodName:  "GenerateImage",
		Summary:     "Generate Image using Gemini 3 Pro Image",
		Description: "Generates an image from a text prompt using Google Gemini 3 Pro Image API.",
		Tags:        outorouter.RegisterTags("AI", "Image"),
		Handler:     handler.GenerateImage,
	})

	// 画像分析用単体テストエンドポイント
	outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[handler.AnalyzeImageRequest, handler.AnalyzeImageResponse]{
		Domain:      "gemini",
		Version:     1,
		MethodName:  "AnalyzeImage",
		Summary:     "Analyze Image using Gemini",
		Description: "Analyzes an image with a text prompt and returns text analysis results using Google Gemini API.",
		Tags:        outorouter.RegisterTags("AI", "Image", "Analysis"),
		Handler:     handler.AnalyzeImage,
	})

	//　画像分析と画像生成を統合したテストエンドポイント (TODO: ゴミ箱データ登録処理と統合する)
	// // JSON版（base64エンコードされた画像データを受け取る）
	// outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[handler.AnalyzeAndGenerateImageRequest, handler.AnalyzeAndGenerateImageResponse]{
	// 	Domain:      "gemini",
	// 	Version:     1,
	// 	MethodName:  "AnalyzeAndGenerateImage",
	// 	Summary:     "Analyze Trash Bin and Generate Monster Character",
	// 	Description: "Analyzes a trash bin image to determine trash type, then generates a monster character themed on that trash type. The character is animal-motifed and based on the trash bin image. Returns base64 encoded image data.",
	// 	Tags:        outorouter.RegisterTags("AI", "Image", "Analysis", "Generation"),
	// 	Handler:     handler.AnalyzeAndGenerateImage,
	// })

	// Multipart版（multipart/form-dataで画像ファイルを受け取る）
	//  (TODO: ゴミ箱データ登録処理と統合する)
	outorouter.RegisterMultipartEndpoint(r, outorouter.MultipartEndpoint[handler.AnalyzeAndGenerateImageMultipartRequest, handler.AnalyzeAndGenerateImageResponse]{
		Domain:      "gemini",
		Version:     1,
		MethodName:  "AnalyzeAndGenerateImage",
		Summary:     "Analyze Trash Bin and Generate Monster Character (Multipart)",
		Description: "Analyzes a trash bin image (sent as multipart/form-data) to determine trash type, then generates a monster character themed on that trash type. The character is animal-motifed and based on the trash bin image. Returns base64 encoded image data.",
		Tags:        outorouter.RegisterTags("AI", "Image", "Analysis", "Generation"),
		Handler:     handler.AnalyzeAndGenerateImageMultipart,
		MaxMemory:   32 * 1024 * 1024, // 32MB
	})

	// Monster登録エンドポイント
	outorouter.RegisterMultipartEndpoint(r, outorouter.MultipartEndpoint[handler.CreateMonsterRequest, handler.CreateMonsterResponse]{
		Domain:      "monster",
		Version:     1,
		MethodName:  "CreateMonster",
		Summary:     "Create Monster",
		Description: "Creates a new monster by analyzing a trash bin image, generating a monster character, and persisting the monster data. Returns the monster ID.",
		Tags:        outorouter.RegisterTags("Monster", "AI", "Image"),
		Handler:     handler.CreateMonster,
		MaxMemory:   32 * 1024 * 1024, // 32MB
	})

	// Monster一覧取得エンドポイント（生成画像）
	outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[handler.GetMonstersRequest, handler.GetMonstersResponse]{
		Domain:      "monster",
		Version:     1,
		MethodName:  "GetMonsters",
		Summary:     "Get Monsters",
		Description: "Returns a list of all monsters with their ID, nickname, latitude, longitude, trash category, and generated monster image URL.",
		Tags:        outorouter.RegisterTags("Monster"),
		Handler:     handler.GetMonsters,
	})

	// ゴミ箱一覧取得エンドポイント（元画像）
	outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[handler.GetTrashsRequest, handler.GetTrashsResponse]{
		Domain:      "trash",
		Version:     1,
		MethodName:  "GetTrashs",
		Summary:     "Get Trashs",
		Description: "Returns a list of all trash bins with their ID, nickname, latitude, longitude, trash category, and original trash bin image URL.",
		Tags:        outorouter.RegisterTags("Trash"),
		Handler:     handler.GetTrashs,
	})

	// Monster一件取得エンドポイント
	outorouter.RegisterUnaryJSONEndpoint(r, outorouter.UnaryJSONEndpoint[handler.GetMonsterRequest, handler.GetMonsterResponse]{
		Domain:      "monster",
		Version:     1,
		MethodName:  "GetMonster",
		Summary:     "Get Monster",
		Description: "Returns a single monster by ID with its nickname, latitude, longitude, trash category, and image URL.",
		Tags:        outorouter.RegisterTags("Monster"),
		Handler:     handler.GetMonster,
	})

	return r.Handler(), nil
}
