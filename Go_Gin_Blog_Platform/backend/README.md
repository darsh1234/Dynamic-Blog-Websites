# Backend (Go + Gin API)

This is the backend service for `Go_Gin_Blog_Platform`.

## What Is Implemented
- JWT auth flow: register, login, refresh, logout
- Password reset request/confirm flow
- Role model: `admin`, `author`, `reader`
- Posts CRUD with pagination and ownership checks
- Admin endpoints for listing users and updating roles
- Structured JSON logging and health endpoint
- PostgreSQL schema migrations (SQL files)

## Structure
- `cmd/api`: app bootstrap and dependency wiring
- `internal/auth`: JWT and hashing utilities
- `internal/config`: env parsing and validation
- `internal/db`: PostgreSQL + GORM connection
- `internal/repository`: data access layer
- `internal/service`: business logic layer
- `internal/transport/http`: Gin handlers and middleware
- `migrations`: SQL migrations

## Local Commands
```bash
cd Go_Gin_Blog_Platform/backend
GOTOOLCHAIN=local CGO_ENABLED=0 go test ./...
psql "postgres://postgres:postgres@localhost:5432/blog_platform?sslmode=disable" < migrations/000001_init.up.sql
go run ./cmd/api
```

Full architecture and deployment docs:
- `../docs/design-setup.md`
- `../docs/aws-deployment-mapping.md`
