FROM postgres:14.2

ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=password

COPY ./scripts/01-init-db.sh /docker-entrypoint-initdb.d/
COPY ./scripts/02-data.sql /docker-entrypoint-initdb.d/

EXPOSE 5432

CMD ["postgres"]
