# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o blog-api ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary dari build stage
COPY --from=builder /app/blog-api .

# Expose port
EXPOSE 8080

# Run aplikasi
CMD ["./blog-api"]