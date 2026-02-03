# AWS Deployment Mapping (Cloud-Ready)

This document defines how localhost architecture maps to AWS deployment.

## 1) Services
- Backend: Go/Gin API container on ECS (Fargate)
- Frontend: React build served from Elastic Beanstalk Docker environment
- Database: PostgreSQL on RDS
- Registry: ECR for backend/frontend images
- Logs: CloudWatch
- Secrets: AWS SSM Parameter Store or Secrets Manager

## 2) Deployment Artifacts in This Repo
- ECS task definition sample: `docs/deployment/ecs-task-definition.sample.json`
- ECS deploy runbook: `docs/deployment/ecs-deploy-runbook.md`
- EB deploy runbook: `docs/deployment/elastic-beanstalk-frontend.md`
- EB Dockerrun sample: `docs/deployment/elastic-beanstalk-dockerrun.sample.json`
- Environment matrix: `docs/deployment/env-matrix.md`
- Release checklist: `docs/deployment/release-checklist.md`

## 3) Environment Strategy
One codebase supports two app variants via env values:
- `APP_VARIANT=blog_a`
- `APP_VARIANT=blog_b`

Recommended environments:
- `dev`
- `staging-blog-a`
- `staging-blog-b`
- `prod-blog-a`
- `prod-blog-b`

## 4) Backend (ECS) Variable Groups
Runtime:
- `APP_ENV`, `PORT`, `DATABASE_URL`, `APP_VARIANT`, `FRONTEND_BASE_URL`

Security/Auth:
- `JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`
- `JWT_ACCESS_TTL_MINUTES`, `JWT_REFRESH_TTL_HOURS`, `PASSWORD_RESET_TTL_MINUTES`

Email:
- `EMAIL_PROVIDER=ses`, `EMAIL_FROM`, `AWS_REGION`, `AWS_SES_FROM_ARN`
- SES sender uses an adapter boundary; wire AWS SDK implementation before production use

Networking/CORS:
- `CORS_ALLOWED_ORIGINS`

## 5) Frontend (Elastic Beanstalk) Variables
- `VITE_API_BASE_URL`
- `VITE_APP_NAME`
- `VITE_BRAND_THEME`

## 6) Networking Baseline
- RDS in private subnets
- ECS tasks in private subnets
- ALB in public subnets routing to ECS tasks
- EB environment reachable via HTTPS
- Security groups only allow required traffic paths

## 7) Health and Operations
- Backend health endpoint: `/api/v1/healthz`
- ALB target group should use `/api/v1/healthz`
- CloudWatch log groups for backend and frontend
- Rolling/blue-green deployment strategy per environment

## 8) Security Baseline
- No hardcoded secrets in source control
- IAM roles least privilege (task execution + app role)
- TLS enforced at public entrypoints
- JWT secrets and DB credentials rotated per environment
