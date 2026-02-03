# Release Checklist (Cloud-Ready)

## Pre-Release
- [ ] Backend tests pass (`go test ./...`)
- [ ] Frontend build passes (`npm run build`)
- [ ] Migration scripts reviewed
- [ ] Environment variables confirmed for target environment
- [ ] JWT secrets available in SSM/Secrets Manager

## Deploy Backend (ECS)
- [ ] Build and push backend image to ECR
- [ ] Register new ECS task definition revision
- [ ] Update ECS service to new revision
- [ ] Confirm ALB target health and `/api/v1/healthz`

## Deploy Frontend (Elastic Beanstalk)
- [ ] Build and push frontend image (`Dockerfile.prod`) to ECR
- [ ] Update `Dockerrun.aws.json` image tag
- [ ] Deploy new EB application version
- [ ] Validate app loads and API calls succeed

## Post-Deploy Validation
- [ ] Register/login/refresh/logout smoke test
- [ ] Create/update/delete post smoke test
- [ ] Admin role update smoke test
- [ ] Password reset request/confirm smoke test
- [ ] CloudWatch logs checked for critical errors

## Rollback Readiness
- [ ] Previous ECS task revision known
- [ ] Previous EB application version available
- [ ] Rollback owner assigned
