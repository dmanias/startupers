#!/bin/sh
set -e

if [ -n "$POSTGRES_URL" ]; then
  echo "Connect with URL."
  until psql "$POSTGRES_URL" -c ";" > /dev/null; do
    sleep 2
  done
else
  echo "Connect with args."
  until PGPASSWORD="$POSTGRES_PASSWORD" psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c ";" > /dev/null; do
    sleep 2
  done
fi
echo "Database is up"
