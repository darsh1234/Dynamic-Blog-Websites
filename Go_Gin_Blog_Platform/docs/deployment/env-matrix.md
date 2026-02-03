# Environment Matrix (Sample)

## Backend Variables by Environment

| Variable | Local | Staging | Production |
|---|---|---|---|
| `APP_ENV` | `local` | `staging` | `production` |
| `DATABASE_URL` | local postgres | RDS staging endpoint | RDS prod endpoint |
| `JWT_ACCESS_SECRET` | local secret | SSM/Secrets Manager | SSM/Secrets Manager |
| `JWT_REFRESH_SECRET` | local secret | SSM/Secrets Manager | SSM/Secrets Manager |
| `EMAIL_PROVIDER` | `stub` | `ses` | `ses` |
| `EMAIL_FROM` | `no-reply@localhost` | staging sender | production sender |
| `AWS_REGION` | optional | required | required |
| `APP_VARIANT` | `blog_a` | `blog_a` or `blog_b` | `blog_a` or `blog_b` |
| `FRONTEND_BASE_URL` | `http://localhost:5173` | staging frontend URL | production frontend URL |
| `CORS_ALLOWED_ORIGINS` | localhost only | staging frontend URLs | production frontend URLs |

## Frontend Variables by Environment

| Variable | Local | Staging | Production |
|---|---|---|---|
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | staging backend URL | production backend URL |
| `VITE_APP_NAME` | `Blog Platform` | branded app name | branded app name |
| `VITE_BRAND_THEME` | `blog_a` | `blog_a` or `blog_b` | `blog_a` or `blog_b` |

## Variant Mapping

| Variant | Example Domain | Frontend Name |
|---|---|---|
| `blog_a` | `blog-a.example.com` | `Blog Platform A` |
| `blog_b` | `blog-b.example.com` | `Blog Platform B` |
