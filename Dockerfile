# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build static binary stripped of debug information
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/api/main.go

# Minimal runtime stage
FROM alpine:latest

WORKDIR /app

# Copy compiled binary and static web assets
COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

EXPOSE 8080

CMD ["./server"]
