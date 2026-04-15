#!/bin/bash

CF_CLIENT_ID="5130be68df401afe4e6571fe6a87a377.access"
CF_CLIENT_SECRET="4246438b20ce00fecaea7f0f4471989f412d6e8b1d77d63b8d7f05d4cd6acaa1"

# Function to shut down background processes
cleanup() {
  echo "Shutting down services..."
  kill "$WEB" "$PUBCHEM" 2>/dev/null
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
(go run cmd/web/main.go --backend "http://localhost:8081" --backend http://localhost:8082 --reload --llm https://ollama.godiegogo.me --db $POSTGRES_URL --cf_client_id $CF_CLIENT_ID --cf_client_secret $CF_CLIENT_SECRET) &
#(go run cmd/web/main.go --backend "http://localhost:8081" --llm http://localhost:11434 --backend http://localhost:8082 --reload --db $POSTGRES_URL) &
WEB=$!
echo "Running the website (pid=$WEB)!"

# Start Caddy in the foreground
echo "Starting Caddy proxy..."
caddy run --config Caddyfile &
CADDY_PID=$!

# Wait for all processes to finish
wait "$WEB" "$CADDY_PID" "$PUBCHEM"
