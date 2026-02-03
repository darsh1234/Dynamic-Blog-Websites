# ECS Backend Deploy Runbook (Sample)

This runbook is for deploying `backend` to ECS Fargate using ECR images.

## 1) Prerequisites
- ECR repository exists: `go-gin-blog-backend`
- ECS cluster exists
- ECS service exists and is connected to an ALB
- RDS PostgreSQL endpoint is reachable from ECS subnets/security groups
- SSM/Secrets Manager contains JWT secrets

## 2) Build and Push Image
```bash
cd Go_Gin_Blog_Platform/backend
aws ecr get-login-password --region <REGION> | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com

docker build -t go-gin-blog-backend:<TAG> .
docker tag go-gin-blog-backend:<TAG> <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/go-gin-blog-backend:<TAG>
docker push <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/go-gin-blog-backend:<TAG>
```

## 3) Register New Task Definition Revision
- Update `ecs-task-definition.sample.json` sample values.
- Register:
```bash
cd Go_Gin_Blog_Platform
aws ecs register-task-definition --cli-input-json file://docs/deployment/ecs-task-definition.sample.json
```

## 4) Update ECS Service
```bash
aws ecs update-service \
  --cluster <CLUSTER_NAME> \
  --service <SERVICE_NAME> \
  --task-definition go-gin-blog-backend:<REVISION> \
  --force-new-deployment
```

## 5) Verify Deployment
- ECS task reaches `RUNNING`
- ALB target health is healthy
- Backend health responds:
```bash
curl https://<BACKEND_DOMAIN>/api/v1/healthz
```
- Smoke test auth and posts routes

## 6) Rollback
```bash
aws ecs update-service \
  --cluster <CLUSTER_NAME> \
  --service <SERVICE_NAME> \
  --task-definition go-gin-blog-backend:<PREVIOUS_REVISION> \
  --force-new-deployment
```
