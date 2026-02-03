# Dynamic Blog Websites

This repository contains three blog projects:
- `Flask_Blog`
- `django_project`
- `Go_Gin_Blog_Platform`

## Projects

| Project | Stack | Docs |
|---|---|---|
| `Flask_Blog` | Flask, SQLAlchemy, Jinja, SQLite | `Flask_Blog/README.md` |
| `django_project` | Django, SQLite, Crispy Forms | `django_project/README.md` |
| `Go_Gin_Blog_Platform` | Go, Gin, React, PostgreSQL, Docker, AWS deployment docs | `Go_Gin_Blog_Platform/README.md` |

## Quick Start

### Flask
```bash
cd Flask_Blog
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
python run.py
```

### Django
```bash
cd django_project
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
python manage.py migrate
python manage.py runserver
```

### Go + Gin + React
```bash
cd Go_Gin_Blog_Platform
docker compose up --build
```
