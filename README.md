# Dynamic Blog Websites Portfolio

This repository contains three blog-focused projects that demonstrate backend and full-stack development across Flask, Django, and Go/Gin.

## Repository Goals
- Show practical implementation of auth, CRUD, pagination, and password reset patterns
- Compare framework-specific approaches (Flask vs Django vs Go/Gin)
- Present a cloud-ready architecture path for a Go microservice-style stack

## Projects at a Glance

| Project | Stack | Current Scope | Docs |
|---|---|---|---|
| `Flask_Blog` | Flask, SQLAlchemy, Jinja, SQLite | Working blog app with auth + post CRUD + reset email | `Flask_Blog/README.md` |
| `django_project` | Django, SQLite, Crispy Forms | Working blog app + extra `expense_dev` app | `django_project/README.md` |
| `Go_Gin_Blog_Platform` | Go, Gin, React, PostgreSQL, Docker, AWS-ready design | Implemented backend APIs + React UI + cloud-ready deployment artifacts | `Go_Gin_Blog_Platform/README.md` |

## Recommended Order to Explore
1. `Flask_Blog` for blueprint-based Flask architecture
2. `django_project` for class-based views and built-in auth workflows
3. `Go_Gin_Blog_Platform` for production-style Go backend design and cloud mapping

## Quick Start

### 1) Flask Blog
```bash
cd Flask_Blog
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
export SECRET_KEY="change-this"
export SQLALCHEMY_DATABASE_URI="sqlite:///site.db"
python run.py
```

### 2) Django Project
```bash
cd django_project
python3 -m venv .venv
source .venv/bin/activate
pip install "Django==4.2.2" "django-crispy-forms==1.14.0" "Pillow==9.5.0"
python manage.py migrate
python manage.py runserver
```

### 3) Go/Gin Blog Platform
See:
- `Go_Gin_Blog_Platform/README.md`
- `Go_Gin_Blog_Platform/docs/design-setup.md`

Run full stack locally:
```bash
cd Go_Gin_Blog_Platform
docker compose up --build
```

## Cloud and Deployment Direction
The Go/Gin platform is intentionally documented as:
- Backend container deployed to ECS (image in ECR)
- Frontend deployed to Elastic Beanstalk
- PostgreSQL hosted on RDS
- Local-first development with cloud-compatible env and service boundaries
- Sample cloud artifacts under `Go_Gin_Blog_Platform/docs/deployment/`

## Why This Repo Works for Portfolio Use
- It demonstrates progression from Python web frameworks to Go service architecture.
- It includes implemented apps plus cloud-ready deployment artifacts and runbooks.
- It keeps deployment concerns explicit (Docker, AWS target topology, env strategy).

## Notes
- Existing Python projects are preserved functionally; documentation has been upgraded for clarity and reproducibility.
- Go platform includes implemented backend/frontend code plus deployment planning artifacts in `Go_Gin_Blog_Platform/docs/`.
