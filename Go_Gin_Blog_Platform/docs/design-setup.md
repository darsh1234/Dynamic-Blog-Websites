# Go_Gin_Blog_Platform - Design and Setup Document

## 1) Project Objective
Build a localhost-first blog platform that demonstrates production-oriented backend engineering in Go while staying cloud-ready for AWS deployment.

Resume alignment:
- Go + Gin REST API
- JWT auth (access + refresh)
- Role-based authorization (`admin`, `author`, `reader`)
- Password reset flow
- Posts CRUD + pagination
- PostgreSQL + SQL migrations
- Dockerized local run
- Deployment model documented for ECS (backend), Elastic Beanstalk (frontend), and RDS (database)

## 2) Scope
In scope:
- Full stack monorepo structure (`backend`, `frontend`)
- Backend API contract, service layering, and data model
- Local development setup and runbook
- Cloud-ready deployment mapping for AWS
- Testing and acceptance criteria

Out of scope (for now):
- Terraform/IaC implementation
- Live AWS deployment execution from this repo
- Advanced moderation workflows, comments, media uploads

## 3) High-Level Architecture

Localhost runtime:
- `frontend` (React + Vite) calls `backend` REST API
- `backend` (Go + Gin) uses PostgreSQL
- Optional local Docker network (frontend + backend + postgres)

Cloud-ready target:
- Backend container -> AWS ECS service
- Frontend static app/container -> AWS Elastic Beanstalk
- PostgreSQL -> AWS RDS (PostgreSQL)
- Email reset provider -> local stub in dev, SES adapter for cloud integration

## 4) Monorepo Layout

```text
Go_Gin_Blog_Platform/
  backend/
    cmd/api/
    internal/
      auth/
      config/
      db/
      email/
      logging/
      models/
      repository/
      service/
      transport/http/
    migrations/
    Dockerfile
    .env.example
  frontend/
    src/
    Dockerfile
    Dockerfile.prod
    .env.example
  docs/
    design-setup.md
    aws-deployment-mapping.md
    deployment/
  docker-compose.yml
  README.md
```

## 5) Backend Layering Responsibilities
- `cmd/api`: application entrypoint, startup lifecycle
- `internal/config`: env loading + strict validation
- `internal/logging`: structured logger (JSON)
- `internal/db`: DB connection and migration integration
- `internal/models`: GORM models
- `internal/repository`: database access patterns
- `internal/service`: business rules, orchestration
- `internal/auth`: JWT mint/verify, password hashing, role checks
- `internal/email`: `EmailSender` interface + local stub + SES adapter
- `internal/transport/http`: Gin router, handlers, middleware

## 6) API Contract (`/api/v1`)

### Health
- `GET /healthz`
  - 200 response:
  ```json
  {"status":"ok"}
  ```

### Auth
- `POST /auth/register`
- `POST /auth/login`
- `POST /auth/refresh`
- `POST /auth/logout`
- `POST /auth/password-reset/request`
- `POST /auth/password-reset/confirm`

Token model:
- Access token: short-lived JWT used in `Authorization: Bearer <token>`
- Refresh token: rotating token; stored server-side as hash

### Posts
- `GET /posts?page=&limit=`
- `GET /posts/:id`
- `POST /posts` (author/admin)
- `PATCH /posts/:id` (author owner/admin)
- `DELETE /posts/:id` (author owner/admin)

### Admin
- `GET /admin/users` (admin)
- `PATCH /admin/users/:id/role` (admin)

### Error Envelope
All controlled errors follow:
```json
{"error":{"code":"<string>","message":"<human-readable>","details":{}}}
```

## 7) Authorization Matrix
- `reader`
  - Can: view posts, manage own auth session
  - Cannot: create/update/delete posts, manage users
- `author`
  - Can: create posts, edit/delete own posts, view posts, manage own auth session
  - Cannot: edit/delete others' posts, manage users
- `admin`
  - Can: all post actions, list users, change roles

## 8) Data Model

`users`
- `id` (uuid, pk)
- `email` (unique)
- `password_hash`
- `role` (`admin|author|reader`)
- `created_at`, `updated_at`

`posts`
- `id` (uuid, pk)
- `author_id` (fk -> users.id)
- `title`
- `content`
- `status` (`draft|published`)
- `created_at`, `updated_at`

`refresh_tokens`
- `id` (uuid, pk)
- `user_id` (fk -> users.id)
- `token_hash`
- `expires_at`
- `revoked_at` (nullable)
- `created_at`

`password_reset_tokens`
- `id` (uuid, pk)
- `user_id` (fk -> users.id)
- `token_hash`
- `expires_at`
- `used_at` (nullable)
- `created_at`

Indexes:
- `users(email)` unique
- `posts(author_id, created_at desc)`
- `refresh_tokens(user_id, revoked_at)`
- `password_reset_tokens(user_id, used_at)`

## 9) Configuration

Backend required env vars:
- `APP_ENV` (`local|staging|production`)
- `PORT` (default `8080`)
- `DATABASE_URL` (PostgreSQL DSN)
- `JWT_ACCESS_SECRET`
- `JWT_REFRESH_SECRET`
- `JWT_ACCESS_TTL_MINUTES`
- `JWT_REFRESH_TTL_HOURS`
- `PASSWORD_RESET_TTL_MINUTES`
- `EMAIL_PROVIDER` (`stub|ses`)
- `EMAIL_FROM`
- `AWS_REGION` (only for SES/cloud)
- `AWS_SES_FROM_ARN` (only for SES/cloud)
- `CORS_ALLOWED_ORIGINS`
- `APP_VARIANT` (supports one codebase deployed as two brands/apps)
- `FRONTEND_BASE_URL` (base URL used in password reset links)
- `REQUEST_TIMEOUT_SECONDS`

Frontend required env vars:
- `VITE_API_BASE_URL`
- `VITE_APP_NAME`
- `VITE_BRAND_THEME`

## 10) Local Setup Workflow

Prerequisites:
- Go 1.21+
- Node.js 20+
- Docker + Docker Compose
- PostgreSQL client tools (optional)

Local workflow:
1. Start PostgreSQL (docker compose)
2. Run backend migrations
3. Start backend API
4. Start frontend Vite app
5. Validate health and auth endpoints

Local password reset behavior:
- Default uses `EmailSender` stub implementation
- Stub logs reset link/token to console
- SES adapter is included and can be wired to AWS SDK for live cloud sends

## 11) AWS Deployment Mapping (Cloud-Ready)

Backend (ECS):
- Container image in ECR
- ECS task definition with env vars/secrets
- Security group allows inbound from ALB only
- Logs to CloudWatch

Frontend (Elastic Beanstalk):
- Build artifacts or container deploy
- Env points to backend base URL
- Separate EB environments for app variants if needed

Database (RDS PostgreSQL):
- Private subnet deployment
- Connection from ECS security group
- Credentials in Secrets Manager/SSM

Email (SES):
- `EmailSender` has an SES adapter boundary; complete AWS SDK wiring before enabling live sends
- Verified sender identity

## 12) Testing Strategy

Backend unit tests:
- token creation/validation
- password hashing
- role authorization checks
- service validation and error mapping

Backend integration tests:
- register/login/refresh/logout flows
- password reset request/confirm
- posts CRUD + pagination
- ownership and role enforcement

Failure tests:
- DB unavailable returns controlled 5xx
- invalid/expired token paths
- revoked refresh token paths

Documentation acceptance tests:
- all READMEs have setup, env vars, run commands, troubleshooting

## 13) Current Coverage
- Backend bootstraps config, logging, database connectivity, and health checks
- Auth supports register/login/refresh/logout with role-aware access control
- Password reset request/confirm endpoints are implemented with local email stub behavior
- Posts API supports CRUD, pagination, and author ownership checks
- Admin API supports user listing and role updates
- Frontend includes auth, posts management, and admin role-management screens
- Docker setup runs PostgreSQL, backend, and frontend locally
- Deployment docs map the stack to ECS, Elastic Beanstalk, and RDS

## 14) Done Criteria
- Project runs locally end-to-end
- API contract implemented and tested
- Role and reset flows verified
- Docs are complete enough for recruiter review
- Architecture clearly maps to ECS/EB/RDS deployment model
