FROM postgres:latest

# Устанавливаем переменные окружения для пользователя Postgres
ENV POSTGRES_USER="KinDeR"
ENV POSTGRES_PASSWORD="Dimaaaa"
ENV POSTGRES_DB="draft"

# Копируем скрипты инициализации базы данных (опционально)
# COPY init.sql /docker-entrypoint-initdb.d/

# Открываем порт для доступа к базе данных
EXPOSE 5432

# Устанавливаем рабочую директорию
WORKDIR /var/lib/postgresql/data

# Стартовая команда (используется стандартный entrypoint Postgres)
CMD ["postgres"]
