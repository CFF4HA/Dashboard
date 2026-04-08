#!/bin/bash

# Function to shut down background processes
cleanup() {
  echo "Shutting down services..."
  kill "$WEB" "$API" "$PUBCHEM" 2>/dev/null
  exit
}

# Execute cleanup function on Ctrl+C or script exit
trap cleanup SIGINT SIGTERM EXIT

# Starting the postgres server
podman run -d --rm --replace --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres postgres:latest
POSTGRES_URL="postgres://postgres:postgres@localhost:5432/postgres"
SQLALCHEMY_URL="postgresql://postgres:postgres@localhost:5432/postgres"

sleep 5 # Wait for the database to initialize

# start the python server
(cd internal/pubchem && .venv/bin/python3 server.py --port 8082 --db $SQLALCHEMY_URL) &
PUBCHEM=$!

# Start the web server
(go run cmd/web/main.go --backend "http://localhost:8081" --reload) &
WEB=$!
echo "Running the website (pid=$WEB)!"

# Start the API server
(go run cmd/backend/main.go --db "$POSTGRES_URL" --port 8081 --pubchem "http://localhost:8082/") &
API=$!
echo "Running the api (pid=$API)!"

# Start Caddy in the foreground
echo "Starting Caddy proxy..."
caddy run --config Caddyfile &
CADDY_PID=$!

# Wait for all processes to finish
wait "$WEB" "$API" "$CADDY_PID" "$PUBCHEM"
