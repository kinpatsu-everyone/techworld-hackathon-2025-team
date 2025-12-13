ATLAS_ENV=local
OUTPUT=../client/lib/client.ts
BASE_URL=https://backend-api-713089770976.asia-northeast1.run.app

.PHONY: help build run test lint clean docker-up docker-down docker-logs docker-build-prod deploy deploy-tag tf-init tf-plan tf-apply tf-destroy sqlc generate

help: ## Display this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	@echo "Building application..."
	@go build -o main ./cmd/main.go

run: ## Run the application locally
	@echo "Running application..."
	@go run ./cmd/main.go

test: ## Run tests
	@echo "Running tests..."
	@go test ./... -v

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	@gofmt -w .
	@go mod tidy

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -f main coverage.out coverage.html
	@rm -rf tmp/

docker-up: ## Start docker compose services
	@echo "Starting docker compose services..."
	@docker-compose --env-file .env.local up

docker-down: ## Stop docker compose services
	@echo "Stopping docker compose services..."
	@docker-compose down

docker-logs: ## Show docker compose logs
	@docker-compose logs -f

docker-restart: docker-down docker-up ## Restart docker compose services

sqlc: ## Generate Go code from SQL queries
	@echo "Generating sqlc code..."
	@command -v sqlc >/dev/null 2>&1 || { echo "sqlc not found, installing..."; go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; }
	@sqlc generate
	@echo "sqlc code generation completed!"

generate: ## Generate TypeScript API client from router definitions
	@echo "Generating TypeScript API client..."
	@go run ./cmd/generate/main.go
	@echo "TypeScript API client generation completed!"

generate-client: ## Generate TypeScript client with custom options (usage: make generate-client OUTPUT=./client.ts BASE_URL=https://api.example.com)
	@echo "Generating TypeScript API client..."
	@go run ./cmd/generate/main.go -output $(or $(OUTPUT),.api/client.ts) -base-url $(or $(BASE_URL),http://localhost:8080)
	@echo "TypeScript API client generation completed!"

.PHONY: migrate
migrate:
	@echo "スキーマの差分からマイグレーションファイルを生成します..."
	@command -v atlas >/dev/null 2>&1 || { echo "atlas not found, installing..."; go install ariga.io/atlas/cmd/atlas@latest; }
	@if atlas migrate diff --env $(ATLAS_ENV); then \
		echo "マイグレーションファイルの生成が完了しました！"; \
	else \
		echo "マイグレーションファイルの生成に失敗しました。"; \
	fi

.PHONY: apply
apply:
	@echo "マイグレーションファイルを適用します..."
	@command -v atlas >/dev/null 2>&1 || { echo "atlas not found, installing..."; go install ariga.io/atlas/cmd/atlas@latest; }
	@if atlas migrate apply --env $(ATLAS_ENV); then \
		echo "マイグレーションファイルの適用が完了しました！"; \
	else \
		echo "マイグレーションファイルの適用に失敗しました。"; \
	fi

# =============================================================================
# Production Build & Deploy (GCP Artifact Registry + Cloud Run)
# =============================================================================

docker-build-prod: ## Build production Docker image locally
	@echo "Building production image..."
	@docker build -f Dockerfile.production --platform linux/amd64 -t backend-api:latest .

deploy: ## Build, push to Artifact Registry, and deploy to Cloud Run (latest tag)
	@./scripts/deploy.sh latest

deploy-tag: ## Deploy with specific tag (usage: make deploy-tag TAG=v1.0.0)
	@./scripts/deploy.sh $(TAG)

# =============================================================================
# Terraform (GCP Infrastructure)
# =============================================================================

TF_DIR := terraform/environments/production

tf-init: ## Initialize Terraform
	@echo "Initializing Terraform..."
	@terraform -chdir=$(TF_DIR) init

tf-plan: ## Plan Terraform changes
	@echo "Planning Terraform changes..."
	@terraform -chdir=$(TF_DIR) plan

tf-apply: ## Apply Terraform changes
	@echo "Applying Terraform changes..."
	@terraform -chdir=$(TF_DIR) apply

tf-destroy: ## Destroy Terraform resources
	@echo "Destroying Terraform resources..."
	@terraform -chdir=$(TF_DIR) destroy
