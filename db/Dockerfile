FROM postgres:12-alpine

LABEL author ="Thibault"
LABEL description = "database initialisation"
LABEL version = "1.0"

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=123456
ENV POSTGRES_DB=todolist_db

COPY init.sql /docker-entrypoint-initdb.d/