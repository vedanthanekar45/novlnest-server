FROM postgres:latest

COPY db/sql/schema/*.sql /docker-entrypoint-initdb.d/