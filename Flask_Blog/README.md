# Flask_Blog

A dynamic blog application built with Flask, SQLAlchemy, and server-rendered templates.

## Features
- User registration, login, logout
- Password hashing with Flask-Bcrypt
- Password reset flow via email token
- Profile update with image upload and resize
- Post create/read/update/delete
- Per-user post feeds
- Home page pagination
- Error pages for `403`, `404`, `500`

## Tech Stack
- Python 3
- Flask, Flask-Login, Flask-WTF, Flask-Mail, Flask-SQLAlchemy
- SQLite (default local option)
- Jinja2 templates + Bootstrap styling

## Project Layout
```text
Flask_Blog/
  run.py
  requirements.txt
  flaskblog/
    __init__.py
    config.py
    models.py
    users/
    posts/
    main/
    templates/
    static/
  instance/
```

## Prerequisites
- Python 3.10+
- `pip`

## Local Setup

1) Create and activate a virtual environment.

```bash
cd Flask_Blog
python3 -m venv .venv
source .venv/bin/activate
```

2) Install dependencies.

```bash
pip install -r requirements.txt
```

3) Set required environment variables.

```bash
export SECRET_KEY="change-this-for-local"
export SQLALCHEMY_DATABASE_URI="sqlite:///site.db"
export EMAIL_USER="your-email@example.com"
export EMAIL_PASS="your-email-app-password"
```

4) Run the app.

```bash
python run.py
```

App runs at `http://127.0.0.1:5000` by default.

## Important Routes
- `/register`, `/login`, `/logout`
- `/account`
- `/reset_password`, `/reset_password/<token>`
- `/post/new`
- `/post/<post_id>`
- `/post/<post_id>/update`
- `/post/<post_id>/delete`
- `/user/<username>`

## Auth and Password Reset
- Passwords are hashed before persistence.
- Reset tokens are generated using timed token signing.
- Reset emails are sent using SMTP credentials from env vars.

## Troubleshooting
- If reset emails fail, confirm `EMAIL_USER` / `EMAIL_PASS` and SMTP access.
- If DB errors appear on startup, verify `SQLALCHEMY_DATABASE_URI` format.
- If profile image upload fails, confirm Pillow was installed correctly.

## Deployment Notes
For production hardening:
- Move secrets to a secret manager
- Disable debug mode
- Use PostgreSQL/MySQL instead of SQLite
- Run behind Gunicorn + reverse proxy
- Configure secure cookie/session settings
