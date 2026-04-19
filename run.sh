#!/usr/bin/env bash
set -euo pipefail

DB_NAME=dashboard
DB_USER=dashboard
DB_PASSWORD=dashboard
DB_PORT=5432
DB_HOST=localhost
CONTAINER_NAME=dashboard-postgres

# Provision the postgres container
if podman ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
  if podman ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo "postgres container already running"
  else
    echo "starting stopped postgres container..."
    podman start "$CONTAINER_NAME"
  fi
else
  echo "creating postgres container..."
  podman run -d \
    --name "$CONTAINER_NAME" \
    -e POSTGRES_DB="$DB_NAME" \
    -e POSTGRES_USER="$DB_USER" \
    -e POSTGRES_PASSWORD="$DB_PASSWORD" \
    -p "${DB_PORT}:5432" \
    postgres:17
fi

# Wait for postgres to be ready
echo "waiting for postgres to be ready..."
for i in $(seq 1 30); do
  if podman exec "$CONTAINER_NAME" pg_isready -U "$DB_USER" -q 2>/dev/null; then
    echo "postgres ready"
    break
  fi
  if [ "$i" -eq 30 ]; then
    echo "postgres did not become ready in time" >&2
    exit 1
  fi
  sleep 1
done

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Database URL: $DATABASE_URL"

go run ./cmd/web/ \
  --address 127.0.0.1 \
  --port 8080 \
  --template-dir ./templates \
  --static-dir ./static \
  --reload \
  --database-url "$DATABASE_URL" \
  --ollama-http-address "http://localhost:11434" \
  --log-level 0
