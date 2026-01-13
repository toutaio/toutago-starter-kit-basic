# Development Dockerfile
FROM golang:1.23-alpine AS base

# Install development tools
RUN apk add --no-cache git make

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Use Air for hot reload
CMD ["air", "-c", ".air.toml"]
