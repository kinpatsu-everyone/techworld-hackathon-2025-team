# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- Initial backend template implementation
- Go 1.25 module setup with `github.com/kinpatsu-everyone/backend-template`
- Production-level HTTP server using standard `net/http` package
- Graceful shutdown support (Ctrl+C) with 30-second timeout
- Configuration management package with `.env.local` file support
- Environment variable utility functions: `GetEnv`, `GetEnvAsInt`, `GetEnvAsBool`, `MustGetEnv`
- Structured logging using `slog` package with JSON output
- HTTP endpoints: `/` (root) and `/health` (health check)
- Docker Compose setup with Go app, PostgreSQL 16, and Redis 7
- Air hot reload configuration for development
- golangci-lint configuration with 25+ linters
- Comprehensive test suite for config and main packages
- Makefile with common development tasks
- Docker and Docker Compose configurations
- Development Dockerfile with Air for hot reload
- Production Dockerfile with multi-stage build
- Documentation in Japanese (README.md, QUICKSTART.md)
- .dockerignore for optimized Docker builds
- .gitignore with Go-specific patterns

### Configuration
- Server configuration (host, port)
- Database configuration (PostgreSQL)
- Redis configuration
- Logging configuration (log level)

### Developer Experience
- Hot reload in development mode
- Docker Compose for local development
- Comprehensive documentation
- Quick start guide
- Example environment variables
