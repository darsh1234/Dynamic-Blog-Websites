# Elastic Beanstalk Frontend Deploy (Sample)

This runbook deploys the frontend to Elastic Beanstalk as a single Docker container.

## 1) Build and Push Frontend Image
```bash
cd Go_Gin_Blog_Platform/frontend
aws ecr get-login-password --region <REGION> | docker login --username AWS --password-stdin <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com

docker build -f Dockerfile.prod \
  --build-arg VITE_API_BASE_URL=https://<BACKEND_DOMAIN>/api/v1 \
  --build-arg VITE_APP_NAME="Blog Platform" \
  --build-arg VITE_BRAND_THEME=blog_a \
  -t go-gin-blog-frontend:<TAG> .
docker tag go-gin-blog-frontend:<TAG> <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/go-gin-blog-frontend:<TAG>
docker push <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/go-gin-blog-frontend:<TAG>
```

## 2) Create Deployment Bundle
- Copy `elastic-beanstalk-dockerrun.sample.json` to `Dockerrun.aws.json`.
- Replace sample values for account, region, and image tag.
- Zip and deploy to EB environment.

## 3) Frontend Config Strategy
- `VITE_*` variables are compiled into static assets at image build time.
- For each environment/brand, build and push a dedicated image tag with the correct `--build-arg` values.
- Keep `Dockerrun.aws.json` image tag aligned with the brand/environment build.

## 4) Health and Smoke Checks
- Frontend URL returns application shell
- `/healthz` path returns `ok`
- Login and posts flows succeed against backend API

## 5) Rollback
- Redeploy previous application version from EB versions list.
