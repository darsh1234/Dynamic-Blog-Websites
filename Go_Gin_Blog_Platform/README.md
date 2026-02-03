# Go_Gin_Blog_Platform

A localhost-first, AWS-cloud-ready blog platform built to demonstrate production-oriented Go backend work.

## Stack
- Backend: Go, Gin, PostgreSQL, JWT, structured logging
- Frontend: React + Vite
- Containers: Docker / Docker Compose
- Cloud target: AWS ECR + ECS (backend), Elastic Beanstalk (frontend), RDS PostgreSQL

## Implemented in This Folder
- Monorepo structure (`backend`, `frontend`, `docs`)
- Backend Go API with:
  - `/api/v1/healthz` with database ping check
  - Implemented auth endpoints (`register`, `login`, `refresh`, `logout`, password reset request/confirm)
  - Role middleware for protected routes
  - Implemented posts CRUD + pagination with ownership checks
  - Implemented admin users list + role update endpoints
  - Standard error envelope for API failures
- Frontend baseline with:
  - auth pages (login/register/forgot/reset)
  - posts dashboard (list/create/edit/delete with pagination)
  - admin user-role management screen
  - token refresh retry behavior on `401`
- Config + logging + PostgreSQL connection bootstrap
- Initial PostgreSQL SQL migration files
- Dockerfiles (backend + frontend) + docker-compose (postgres + backend + frontend)
- Detailed design and AWS mapping documents

## Monorepo Layout
```text
Go_Gin_Blog_Platform/
  backend/
    cmd/api/main.go
    internal/config/
    internal/logging/
    internal/transport/http/
    migrations/
    .env.example
    Dockerfile
  frontend/
    .env.example
    Dockerfile
    Dockerfile.prod
  docs/
    design-setup.md
    aws-deployment-mapping.md
    deployment/
  docker-compose.yml
```

## Quick Start (Local Full Stack)

### Option A: Run with Docker Compose
```bash
cd Go_Gin_Blog_Platform
docker compose up --build
```

Apply initial schema migration once PostgreSQL is up:
```bash
docker compose exec -T postgres psql -U postgres -d blog_platform < backend/migrations/000001_init.up.sql
```

Local URLs:
- Frontend: `http://localhost:5173`
- Backend health: `http://localhost:8080/api/v1/healthz`

Backend health check:
```bash
curl http://localhost:8080/api/v1/healthz
```

### Option B: Run backend directly
```bash
cd Go_Gin_Blog_Platform/backend
psql "postgres://postgres:postgres@localhost:5432/blog_platform?sslmode=disable" < migrations/000001_init.up.sql
go run ./cmd/api
```

## API Surface (`/api/v1`)
- Auth:
  - `POST /auth/register`
  - `POST /auth/login`
  - `POST /auth/refresh`
  - `POST /auth/logout`
  - `POST /auth/password-reset/request`
  - `POST /auth/password-reset/confirm`
- Posts:
  - `GET /posts`
  - `GET /posts/:id`
  - `POST /posts`
  - `PATCH /posts/:id`
  - `DELETE /posts/:id`
- Admin:
  - `GET /admin/users`
  - `PATCH /admin/users/:id/role`

Current implementation status:
- Auth endpoints are implemented and return token pairs.
- Posts CRUD + pagination endpoints are implemented.
- Admin user management endpoints are implemented.
- React frontend baseline is implemented and wired to backend API contracts.

## Validation Commands
Backend:
```bash
cd Go_Gin_Blog_Platform/backend
GOTOOLCHAIN=local CGO_ENABLED=0 go test ./...
```

Frontend:
```bash
cd Go_Gin_Blog_Platform/frontend
npm run build
```

## Role Model
- `reader`: read-only
- `author`: manage own posts
- `admin`: full post access + user role management

## Environment Model
Backend env examples are in:
- `Go_Gin_Blog_Platform/backend/.env.example`

Frontend env examples are in:
- `Go_Gin_Blog_Platform/frontend/.env.example`

## Core Documentation
- Design and setup: `Go_Gin_Blog_Platform/docs/design-setup.md`
- AWS deployment mapping: `Go_Gin_Blog_Platform/docs/aws-deployment-mapping.md`
- Deployment artifacts: `Go_Gin_Blog_Platform/docs/deployment/`

## Project Highlights
This project demonstrates:
- Gin REST API design
- JWT + roles + password reset flow
- PostgreSQL migration strategy
- Dockerized local workflow
- AWS-targeted deployment model (ECS/EB/RDS)
