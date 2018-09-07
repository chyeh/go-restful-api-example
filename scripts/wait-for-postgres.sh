#!/bin/sh
set -e

until PGPASSWORD=$4 psql -h $1 -U $2 -d $3 sslmode=disable -c '\q'; do
  >&2 echo "Postgres is unavailable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $5