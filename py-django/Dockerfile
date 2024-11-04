FROM python:3.12-slim

WORKDIR /py-django

COPY requirements.txt .

RUN apt-get update && \
    apt-get install curl -y && \
    pip install --no-cache-dir -r requirements.txt

COPY . .

EXPOSE 8000

ENV DJANGO_SETTINGS_MODULE=videos.settings \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

CMD ["sh", "-c",  \
    "python manage.py makemigrations && \
    python manage.py migrate && \
    python manage.py collectstatic --noinput && \
    python manage.py init_admin && \
    gunicorn videos.wsgi:application --bind 0.0.0.0:8000"]
