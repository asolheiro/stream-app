# Use a imagem base do Python
FROM python:3.12-slim

WORKDIR /py-django

EXPOSE 8000

ENV DJANGO_SETTINGS_MODULE=videos.settings \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

COPY supervisord.conf /etc/supervisord.conf
COPY requirements.txt .

RUN apt update \
    && apt install curl supervisor gdal-bin -y --no-install-recommends \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir -r requirements.txt
    

COPY . /py-django

RUN python manage.py makemigrations && \
    python manage.py collectstatic --noinput


CMD python manage.py migrate \
    && python manage.py init_admin \
    && supervisord -c /etc/supervisord.conf



# FROM python:3.12-slim

# WORKDIR /py-django

# COPY requirements.txt .

# RUN apt-get update && \
#     apt-get install curl -y && \
#     pip install --no-cache-dir -r requirements.txt

# RUN apt-get update && \
#     apt-get install -y supervisor

# COPY . .

# EXPOSE 8000

# ENV DJANGO_SETTINGS_MODULE=videos.settings \
#     PYTHONUNBUFFERED=1 \
#     PYTHONDONTWRITEBYTECODE=1

# RUN python manage.py makemigrations && \
#     python manage.py collectstatic --noinput
    
# CMD ["python manage.py migrate && \
#     python manage.py init_admin"]

# ENTRYPOINT gunicorn videos.wsgi:application --bind 0.0.0.0:8000
