# backend-template

## 概要

GoとMySQLを使用したバックエンドアプリケーションの開発テンプレートです。

## 技術スタック

- **Go** v1.25.5
- **PostgreSQL** 16
- **Redis** 7
- **Docker** & Docker Compose
- **Air** (ホットリロード)
- **golangci-lint** (静的解析)

## 特徴

- ✅ 標準の `net/http` パッケージを使用したHTTPサーバー
- ✅ グレースフルシャットダウン対応（Ctrl+C）
- ✅ `.env.local` ファイルからの設定読み込み
- ✅ `slog` パッケージによる構造化ログ
- ✅ Docker Compose でのローカル開発環境
- ✅ Air によるホットリロード対応
- ✅ PostgreSQL と Redis コンテナ

## ディレクトリ構成

```
.
├── config/              # 設定管理パッケージ
│   └── config.go
├── main.go              # メインアプリケーション
├── .env.local.example   # 環境変数のサンプル
├── .golangci.yml        # golangci-lint設定
├── .air.toml            # Air設定（ホットリロード）
├── Dockerfile           # 本番用Dockerfile
├── Dockerfile.dev       # 開発用Dockerfile
├── docker-compose.yaml  # Docker Compose設定
└── README.md
```

## セットアップ

### 前提条件

- Go 1.25.5 以上
- Docker & Docker Compose
- golangci-lint（オプション）

### 1. 環境変数の設定

`.env.local.example` をコピーして `.env.local` を作成します：

```bash
cp .env.local.example .env.local
```

必要に応じて `.env.local` の値を編集してください。

### 2. ローカル開発（Docker Compose使用）

Docker Compose を使用してアプリケーション、PostgreSQL、Redis を起動します：

```bash
docker-compose up
```

これにより以下が起動します：
- Go アプリケーション（ポート 8080、ホットリロード有効）
- PostgreSQL（ポート 5432）
- Redis（ポート 6379）

コードを編集すると自動的に再ビルド・再起動されます。

### 3. ローカル開発（Go直接実行）

依存関係をインストール：

```bash
go mod download
```

アプリケーションを実行：

```bash
go run main.go
```

## エンドポイント

### `GET /`
ルートエンドポイント。

**レスポンス：**
```json
{"message":"Welcome to backend-template"}
```

### `GET /health`
ヘルスチェックエンドポイント。

**レスポンス：**
```json
{"status":"ok"}
```

## ビルド

### 本番用バイナリのビルド

```bash
go build -o main .
```

### Dockerイメージのビルド

```bash
docker build -t backend-template:latest .
```

## テスト

```bash
go test ./...
```

## 静的解析

golangci-lint を使用してコードの静的解析を実行：

```bash
golangci-lint run
```

## 設定

環境変数は `.env.local` ファイルまたはシステムの環境変数から読み込まれます。

| 変数名 | デフォルト値 | 説明 |
|--------|-------------|------|
| `SERVER_HOST` | `0.0.0.0` | サーバーのホスト |
| `SERVER_PORT` | `8080` | サーバーのポート |
| `DB_HOST` | `localhost` | データベースのホスト |
| `DB_PORT` | `5432` | データベースのポート |
| `DB_USER` | `postgres` | データベースのユーザー名 |
| `DB_PASSWORD` | `postgres` | データベースのパスワード |
| `DB_NAME` | `backend_db` | データベース名 |
| `DB_SSLMODE` | `disable` | SSL接続モード |
| `REDIS_HOST` | `localhost` | Redisのホスト |
| `REDIS_PORT` | `6379` | Redisのポート |
| `REDIS_PASSWORD` | `` | Redisのパスワード |
| `REDIS_DB` | `0` | Redisのデータベース番号 |
| `LOG_LEVEL` | `info` | ログレベル（debug/info/warn/error） |

## グレースフルシャットダウン

アプリケーションは `SIGINT`（Ctrl+C）または `SIGTERM` シグナルを受け取ると、グレースフルシャットダウンを実行します。現在処理中のリクエストが完了するまで最大30秒待機してからシャットダウンします。

## ライセンス

MIT

