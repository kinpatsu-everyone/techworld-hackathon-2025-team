# Gemini API テスト手順

## 環境変数の設定

`.env.local` または環境変数に以下を設定してください：

```bash
GEMINI_API_KEY=your-api-key-here
```

## サーバーの起動

```bash
go run cmd/main.go
```

または

```bash
docker-compose up
```

## APIリクエスト例

### 1. 基本的な画像生成リクエスト

```bash
curl -X POST http://localhost:8080/v1/gemini/GenerateImage \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Generate a 4K photorealistic image of a yellow banana floating in space with Earth in the background, add text overlay: \"Nano Banana Pro\""
  }'
```

### 2. モデルを指定したリクエスト

```bash
curl -X POST http://localhost:8080/v1/gemini/GenerateImage \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A beautiful sunset over the ocean",
    "model": "gemini-3-pro-image-preview"
  }'
```

### 3. レスポンス例

```json
{
  "candidates": [
    {
      "content": {
        "role": "model",
        "parts": [
          {
            "text": "...",
            "image_data": {
              "mime_type": "image/png",
              "data": "iVBORw0KGgoAAAANSUhEUgAA..."
            }
          }
        ]
      },
      "finish_reason": "STOP"
    }
  ]
}
```

## レスポンスの説明

- `candidates`: 生成された候補のリスト
- `content.role`: コンテンツの役割（"user" または "model"）
- `content.parts`: コンテンツのパーツ（テキスト、画像データなど）
- `parts.text`: テキストコンテンツ（存在する場合）
- `parts.image_data`: 画像データ（base64エンコード、存在する場合）
  - `mime_type`: 画像のMIMEタイプ（例: "image/png"）
  - `data`: base64エンコードされた画像データ
- `parts.file_data`: ファイルデータ（URI、存在する場合）
  - `mime_type`: ファイルのMIMEタイプ
  - `file_uri`: ファイルのURI
- `finish_reason`: 生成が終了した理由（"STOP", "MAX_TOKENS" など）

## 画像データの使用方法

レスポンスに含まれるbase64エンコードされた画像データは、以下のように使用できます：

```javascript
// JavaScript例
const imageData = response.candidates[0].content.parts[0].image_data.data;
const image = document.createElement('img');
image.src = `data:${response.candidates[0].content.parts[0].image_data.mime_type};base64,${imageData}`;
document.body.appendChild(image);
```

```python
# Python例
import base64
from PIL import Image
from io import BytesIO

image_data = response['candidates'][0]['content']['parts'][0]['image_data']['data']
image_bytes = base64.b64decode(image_data)
image = Image.open(BytesIO(image_bytes))
image.show()
```

## 画像分析API

### 1. 基本的な画像分析リクエスト

画像データをbase64エンコードして送信し、分析結果をテキストで受け取ります。

```bash
# 画像ファイルをbase64エンコード（macOS/Linux）
IMAGE_BASE64=$(base64 -i path/to/image.jpg)

curl -X POST http://localhost:8080/gemini/v1/AnalyzeImage \
  -H "Content-Type: application/json" \
  -d "{
    \"prompt\": \"この画像に写っているものを詳しく説明してください\",
    \"image_data\": \"$IMAGE_BASE64\",
    \"mime_type\": \"image/jpeg\"
  }"
```

### 1-1. Postmanで送信する場合

**リクエスト設定:**
- Method: `POST`
- URL: `http://localhost:8080/gemini/v1/AnalyzeImage`
- Headers: `Content-Type: application/json`
- Body: `raw` を選択し、`JSON` を選択

**Body (JSON):**
```json
{
  "prompt": "この画像に写っているものを詳しく説明してください",
  "image_data": "ここにbase64エンコードされた画像データを貼り付けてください",
  "mime_type": "image/jpeg"
}
```

**画像をbase64エンコードする方法:**

1. **オンラインツールを使用:**
   - https://base64.guru/converter/encode/image
   - 画像をアップロードしてbase64文字列を取得

2. **コマンドラインでエンコード:**
   ```bash
   # macOS/Linux
   base64 -i image.jpg
   
   # または
   cat image.jpg | base64
   ```

3. **出力結果をそのまま `image_data` フィールドに貼り付け**

**完全な例（小さなサンプル画像の場合）:**
```json
{
  "prompt": "この画像に写っているものを詳しく説明してください",
  "image_data": "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
  "mime_type": "image/png"
}
```

**注意:** base64エンコードされた画像データは非常に長くなるため、実際の画像を使用する場合は、上記の方法でエンコードした文字列を使用してください。

### 2. Pythonで画像をbase64エンコードして送信

```python
import base64
import requests
import json

# 画像ファイルを読み込んでbase64エンコード
with open("path/to/image.jpg", "rb") as image_file:
    image_data = base64.b64encode(image_file.read()).decode('utf-8')

# APIリクエスト
response = requests.post(
    "http://localhost:8080/gemini/v1/AnalyzeImage",
    headers={"Content-Type": "application/json"},
    json={
        "prompt": "この画像に写っているものを詳しく説明してください",
        "image_data": image_data,
        "mime_type": "image/jpeg"
    }
)

print(response.json())
```

### 3. レスポンス例

```json
{
  "text": "この画像には、美しい夕日が海の向こうに沈んでいく様子が写っています。空はオレンジ色とピンク色のグラデーションで彩られており、雲が柔らかく浮かんでいます。海面は穏やかで、夕日の光を反射して輝いています。全体的に非常に平和で美しい風景です。"
}
```

### 4. モデルを指定したリクエスト

```bash
curl -X POST http://localhost:8080/gemini/v1/AnalyzeImage \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "この画像の色使いと構図について分析してください",
    "image_data": "iVBORw0KGgoAAAANSUhEUgAA...",
    "mime_type": "image/png",
    "model": "gemini-2.5-flash"
  }'
```

### 5. 画像ファイルをbase64エンコードする方法

**macOS/Linux:**
```bash
base64 -i image.jpg
```

**Windows (PowerShell):**
```powershell
[Convert]::ToBase64String([IO.File]::ReadAllBytes("image.jpg"))
```

**Node.js:**
```javascript
const fs = require('fs');
const imageBuffer = fs.readFileSync('image.jpg');
const base64Image = imageBuffer.toString('base64');
console.log(base64Image);
```

**Python:**
```python
import base64

with open("image.jpg", "rb") as image_file:
    base64_image = base64.b64encode(image_file.read()).decode('utf-8')
    print(base64_image)
```

## エンドポイント一覧

- `POST /gemini/v1/GenerateImage` - テキストプロンプトから画像を生成
- `POST /gemini/v1/AnalyzeImage` - 画像データを分析してテキストで返す


