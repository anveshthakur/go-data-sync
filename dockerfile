# use the latest postgres image
FROM postgres:latest

# set environment variables
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB concurrency


# Expose the postgres port
EXPOSE 5432