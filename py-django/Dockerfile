# Use uma imagem Python oficial
FROM python:3.12-slim

# Define o diretório de trabalho
WORKDIR /py-django

# Copia o arquivo de requisitos e instala as dependências
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copia o restante do código da aplicação para o contêiner
COPY . .

# Expõe a porta padrão do Django (8000)
EXPOSE 8000

# Define variáveis de ambiente para Django
ENV DJANGO_SETTINGS_MODULE=videos.settings \
    PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1

# Executa as migrações e inicia o servidor usando o Gunicorn
CMD ["sh", "-c", "python manage.py migrate && gunicorn videos.wsgi:application --bind 0.0.0.0:8000"]
