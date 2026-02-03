# django_project

A Django workspace containing two apps:
- `blog`: full blog workflow (auth, posts CRUD, profile, pagination, password reset)
- `expense_dev`: simple secondary app with its own route/template

## Features

### Blog app (`blog` + `users`)
- User registration and login/logout
- Profile update (username, email, profile image)
- Password reset flow (email-based)
- Post create/read/update/delete
- Per-user post listings
- Pagination on homepage and user pages

### Expense app (`expense_dev`)
- Standalone app entry route and template
- Useful as an example of multi-app Django project structure

## Tech Stack
- Python 3
- Django 4.2
- SQLite (default local DB)
- Crispy Forms
- Pillow for profile image processing

## Project Layout
```text
django_project/
  manage.py
  django_project/
    settings.py
    urls.py
  blog/
  users/
  expense_dev/
  media/
  db.sqlite3
```

## Prerequisites
- Python 3.10+
- `pip`

## Local Setup

1) Create and activate a virtual environment.

```bash
cd django_project
python3 -m venv .venv
source .venv/bin/activate
```

2) Install dependencies.

Option A (exact repo requirements):
```bash
pip install -r requirements.txt
```

Option B (minimal runtime set):
```bash
pip install "Django==4.2.2" "django-crispy-forms==1.14.0" "Pillow==9.5.0"
```

3) Set optional email variables for password reset.

```bash
export EMAIL_USER="your-email@example.com"
export EMAIL_PASS="your-email-app-password"
```

4) Apply migrations and run server.

```bash
python manage.py migrate
python manage.py runserver
```

App runs at `http://127.0.0.1:8000` by default.

## Important Routes
- Blog home: `/`
- Blog about: `/about/`
- Post detail: `/post/<id>/`
- Create/update/delete posts: `/post/new/`, `/post/<id>/update/`, `/post/<id>/delete/`
- User posts: `/user/<username>`
- Register/login/logout: `/register/`, `/login/`, `/logout/`
- Profile: `/profile/`
- Password reset: `/password-reset/` flow routes
- Expense app: `/expense_dev`

## Notes on Auth and Reset
- Uses Django auth views for login/logout/password reset.
- Email backend is SMTP-based in settings and reads credentials from env vars.

## Troubleshooting
- If media files do not load, confirm `DEBUG=True` for local and that `MEDIA_ROOT` exists.
- If password reset email is not sent, verify SMTP creds and sender policy.
- If migrations fail, remove stale local DB only if you intend to reset local data.

## Deployment Notes
For production readiness:
- Set `DEBUG=False` and secure `ALLOWED_HOSTS`
- Move secrets into environment/secret manager
- Use PostgreSQL
- Serve static/media via CDN or object storage
- Run with Gunicorn/ASGI behind reverse proxy
