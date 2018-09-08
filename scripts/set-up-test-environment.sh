#!/bin/bash

export COMPOSE_PROJECT_NAME=$1
docker-compose up -d

until PGPASSWORD=hellofresh psql -h localhost -U hellofresh sslmode=disable -c '\q' &>/dev/null; do
  >&2 echo "postgres is not ready, try again"
  sleep 1
done

./init-db-schema.sh postgres://hellofresh:hellofresh@localhost:5432/hellofresh >/dev/null
./init-db-user-data.sh postgres://hellofresh:hellofresh@localhost:5432/hellofresh >/dev/null

until curl localhost &>/dev/null; do
  >&2 echo "go application is not ready, try again"
  sleep 1
done