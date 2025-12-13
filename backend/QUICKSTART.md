# クイックスタートガイド

## 🚀 最速スタート（Docker使用）

```bash
# 1. 環境変数ファイルをコピー
cp .env.local.example .env.local

# 2. Docker Composeで起動
docker-compose up

# 3. 動作確認
curl http://localhost:8080
curl http://localhost:8080/health
```

## 💻 ローカル開発（Go直接実行）

```bash
# 依存関係のダウンロード
go mod download

# アプリケーションの実行
go run main.go

# または、ビルドして実行
make build
./main
```

## 📝 よく使うコマンド

```bash
# テスト実行
make test

# カバレッジ付きテスト
make test-coverage

# フォーマット
make fmt

# 静的解析
make vet

# Linter実行（golangci-lintが必要）
make lint

# クリーンアップ
make clean
```

## 🐳 Docker コマンド

```bash
# 起動
make docker-up

# 停止
make docker-down

# ログ確認
make docker-logs

# 再起動
make docker-restart
```

## 📦 サービス

Docker Composeで起動すると、以下のサービスが利用可能になります：

- **アプリケーション**: http://localhost:8080
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## 🔧 開発のヒント

### ホットリロード

Docker Composeで起動すると、コードの変更が自動的に検出され、アプリケーションが再起動します。

### デバッグログ

`.env.local`で`LOG_LEVEL=debug`に設定すると、詳細なログが出力されます。

### データベース接続

PostgreSQLへの接続情報：

```
Host: localhost (Docker外) / db (Docker内)
Port: 5432
User: postgres
Password: postgres
Database: backend_db
```

### Redis接続

Redisへの接続情報：

```
Host: localhost (Docker外) / redis (Docker内)
Port: 6379
```

## 🧪 テスト

```bash
# 全テスト実行
go test ./... -v

# 特定パッケージのテスト
go test ./config -v

# カバレッジレポート生成
make test-coverage
# ブラウザで coverage.html を開く
```

## 🛠️ トラブルシューティング

### ポートが使用中の場合

`.env.local`の`SERVER_PORT`を変更してください。

### Docker コンテナが起動しない

```bash
# コンテナとボリュームを削除して再起動
docker-compose down -v
docker-compose up --build
```

### 依存関係の問題

```bash
# go.modを整理
go mod tidy

# キャッシュをクリア
go clean -modcache
```
