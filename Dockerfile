FROM postgres:latest

ENV POSTGRES_USER="KinDeR"
ENV POSTGRES_PASSWORD="Dimaaaa"
ENV POSTGRES_DB="draft"

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

WORKDIR /var/lib/postgresql/data

CMD ["postgres"]
