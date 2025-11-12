FROM migrate/migrate:v4.15.2

COPY ./api/data/migrations /migrations
