FROM postgres:latest

ENV POSTGRES_USER="KinDeR"
ENV POSTGRES_PASSWORD="Ajklsdha"
ENV POSTGRES_DB="Auth"

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

WORKDIR /var/lib/postgresql/data

CMD ["postgres"]
