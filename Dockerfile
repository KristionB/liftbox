# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build server
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./server/main.go

# Build client
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/client ./client/cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/server /app/server
COPY --from=builder /app/client /app/client

# Expose server port
EXPOSE 8080

# Default to running server
CMD ["./server"]

