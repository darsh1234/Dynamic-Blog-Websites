# Frontend (React + Vite)

React frontend baseline for the Go/Gin blog platform.

## What It Includes
- Auth screens: login, register, forgot/reset password
- Token-aware API client with refresh token retry flow
- Posts UI: list with pagination, create, edit, delete
- Admin UI: list users and update role (`admin` only)
- Route protection:
  - authenticated routes for `/posts`
  - role-based route for `/admin/users`

## Tech Stack
- React 18
- React Router 6
- Vite 5

## Local Setup

```bash
cd Go_Gin_Blog_Platform/frontend
cp .env.example .env
npm install
npm run dev
```

App default URL: `http://localhost:5173`

## Docker Setup (From Monorepo Root)

```bash
cd Go_Gin_Blog_Platform
docker compose up --build
```

This starts:
- frontend on `http://localhost:5173`
- backend on `http://localhost:8080`
- postgres on `localhost:5432`

For Elastic Beanstalk production-style deploy, use:
- `Dockerfile.prod`
- `nginx.conf`
- `../docs/deployment/elastic-beanstalk-frontend.md`

Production image build example:
```bash
docker build -f Dockerfile.prod \
  --build-arg VITE_API_BASE_URL=https://<BACKEND_DOMAIN>/api/v1 \
  --build-arg VITE_APP_NAME="Blog Platform" \
  --build-arg VITE_BRAND_THEME=blog_a \
  -t go-gin-blog-frontend:prod .
```

## Environment Variables
`.env.example`:

```env
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_APP_NAME=Blog Platform
VITE_BRAND_THEME=blog_a
```

## API Expectations
Backend should expose:
- `/api/v1/auth/*` for register/login/refresh/logout/reset
- `/api/v1/posts` CRUD + pagination
- `/api/v1/admin/users` list + role update

## Notes
- Access and refresh tokens are stored in `localStorage` for local development convenience.
- On `401`, the client attempts `/auth/refresh` once and retries the original request.
- UI is intentionally lightweight but production-shaped so you can extend quickly.
