# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go-based backend template using PostgreSQL and Redis. Uses standard library `net/http` for HTTP server with graceful shutdown support. Configuration loaded from `.env.local` files. Structured logging via `slog` with JSON output.

## Development Commands

### Building and Running
```bash
make build           # Build production binary
make run             # Run locally with Go
go run main.go       # Alternative direct run
```

### Testing
```bash
make test            # Run all tests
make test-coverage   # Generate HTML coverage report
go test ./config -v  # Test specific package
```

### Code Quality
```bash
make lint            # Run golangci-lint (extensive checks enabled)
make fmt             # Format code and tidy modules
make vet             # Run go vet
```

### Docker Development
```bash
make docker-up       # Start all services (app, PostgreSQL, Redis)
make docker-down     # Stop services
make docker-logs     # Follow logs
make docker-restart  # Restart all services
docker-compose up    # Start with hot reload enabled
```

### Cleanup
```bash
make clean           # Remove build artifacts and coverage files
```

## Architecture

### Application Structure

- **main.go**: Entry point with HTTP server setup, route registration, and graceful shutdown handling
  - `setupLogger()`: Configures slog handler based on LOG_LEVEL
  - `setupRoutes()`: Registers HTTP handlers
  - Graceful shutdown: 30-second timeout for in-flight requests

- **config/config.go**: Configuration management with env file loading
  - `Load()`: Reads from `.env.local` then falls back to system env vars
  - Helper functions: `GetEnv()`, `GetEnvAsInt()`, `GetEnvAsBool()`, `MustGetEnv()`
  - Config structs: ServerConfig, DatabaseConfig, RedisConfig, LogConfig

### Configuration Management

Environment variables are loaded with this precedence:
1. System environment variables (highest priority)
2. `.env.local` file values (only if system env not set)
3. Default values in code (lowest priority)

The `loadEnvFile()` function skips empty lines and comments, strips quotes, and only sets vars not already in the environment.

### Docker Environment

Three services defined in `docker-compose.yaml`:
- **app**: Uses `Dockerfile.dev` with Air for hot reload, mounts code volume
- **db**: PostgreSQL 16 on port 5432, healthcheck with `pg_isready`
- **redis**: Redis 7 on port 6379, healthcheck with `redis-cli ping`

Hot reload: Air watches `.go`, `.tpl`, `.tmpl`, `.html` files and rebuilds to `tmp/main` on changes. Excludes test files.

### Network Configuration

Important hostname differences:
- **From host machine**: Use `localhost` for DB_HOST and REDIS_HOST
- **Inside Docker**: Use `db` and `redis` service names (docker-compose.yaml sets these)

### Linting Configuration

golangci-lint (.golangci.yml) has 25+ linters enabled including:
- Security: gosec
- Style: revive, gocritic
- Performance: prealloc
- Correctness: errcheck, staticcheck, govet
- Import organization: `github.com/kinpatsu-everyone/backend-template` as local prefix

## Key Environment Variables

Required for database connection:
- `DB_HOST`: `localhost` (host) or `db` (Docker)
- `DB_PORT`: 5432
- `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`

Required for Redis:
- `REDIS_HOST`: `localhost` (host) or `redis` (Docker)
- `REDIS_PORT`: 6379
- `REDIS_PASSWORD`, `REDIS_DB`

Server config:
- `SERVER_HOST`: Default `0.0.0.0`
- `SERVER_PORT`: Default 8080

Logging:
- `LOG_LEVEL`: `debug`, `info`, `warn`, `error` (default: `info`)

## Testing Practices

Tests are located alongside source files (`*_test.go`):
- `config/config_test.go`: Tests env var loading and parsing
- `main_test.go`: Tests HTTP handlers

Use table-driven tests. See existing test files for patterns.

## Common Development Workflows

### Initial Setup
```bash
cp .env.local.example .env.local
# Edit .env.local as needed
docker-compose up
```

### Adding New Handlers
1. Define handler function in main.go (or create new package)
2. Register in `setupRoutes()` function in main.go
3. Set `Content-Type: application/json` header
4. Use `slog` for logging with structured fields
5. Add tests in `main_test.go`

### Adding Configuration
1. Add field to appropriate Config struct in config/config.go
2. Add default value in `Load()` function using `GetEnv()` or `GetEnvAsInt()`
3. Document in `.env.local.example`
4. Add test in `config/config_test.go`

### HTTP Server Patterns
- Use `http.ServeMux` for routing (standard library)
- Set timeouts: ReadTimeout, WriteTimeout (15s), IdleTimeout (60s)
- Return JSON responses with appropriate status codes
- Log requests with structured fields: "method", "path", etc.

## Module and Import Management

Module path: `github.com/kinpatsu-everyone/backend-template`

When adding dependencies:
```bash
go get <package>
make fmt  # Runs go mod tidy
```
