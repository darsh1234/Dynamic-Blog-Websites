# Go_Gin_Blog_Platform

A full-stack blog platform with a Go (Gin) backend, React frontend, and PostgreSQL.

## Stack
- Backend: Go, Gin, PostgreSQL, JWT, structured logging
- Frontend: React + Vite
- Containers: Docker / Docker Compose
- Deployment docs: AWS ECR + ECS (backend), Elastic Beanstalk (frontend), RDS PostgreSQL

## Features
- Monorepo structure (`backend`, `frontend`, `docs`)
- JWT auth (`register`, `login`, `refresh`, `logout`)
- Password reset request/confirm
- Role-based authorization (`reader`, `author`, `admin`)
- Posts CRUD with pagination and ownership checks
- Admin user listing and role update endpoints
- Structured API error responses
- React UI for auth, posts, and admin user-role management
- Docker Compose setup for local Postgres + backend + frontend

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

## Local Setup

### Run with Docker Compose
```bash
cd Go_Gin_Blog_Platform
docker compose up --build
```

Apply schema migration:
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

### Run services directly (optional)
Backend:
```bash
cd Go_Gin_Blog_Platform/backend
psql "postgres://postgres:postgres@localhost:5432/blog_platform?sslmode=disable" < migrations/000001_init.up.sql
go run ./cmd/api
```

Frontend:
```bash
cd Go_Gin_Blog_Platform/frontend
npm install
npm run dev
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

## Additional Docs
- Design and setup: `Go_Gin_Blog_Platform/docs/design-setup.md`
- AWS deployment mapping: `Go_Gin_Blog_Platform/docs/aws-deployment-mapping.md`
- Deployment artifacts: `Go_Gin_Blog_Platform/docs/deployment/`
