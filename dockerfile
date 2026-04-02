# Build stage
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app "./cmd"

# Run stage
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/app .

COPY --from=builder /app/internal/database/migrations/migrations ./internal/database/migrations/migrations

EXPOSE $SERVER_PORT

CMD ["./app"]