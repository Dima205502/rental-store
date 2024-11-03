FROM postgres:latest

ENV POSTGRES_USER="KinDeR"
ENV POSTGRES_PASSWORD="lksdADops"
ENV POSTGRES_DB="Shop"

COPY init.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

WORKDIR /var/lib/postgresql/data

CMD ["postgres"]