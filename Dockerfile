FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/dashboard ./cmd/web/
RUN cp /app/dashboard /app/server
RUN cp /app/server /bin/server

EXPOSE 8080

CMD /app/dashboard \
  --address 0.0.0.0 \
  --port 8080 \
  --template-dir /app/templates \
  --static-dir /app/static \
  --log-level 0 \
  --database-url "$DATABASE_URL" \
  --ollama-http-address "$OLLAMA_HTTP_ADDRESS"
