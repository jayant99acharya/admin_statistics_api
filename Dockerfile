# Development Dockerfile with Go toolchain
FROM golang:1.21-alpine

# Install git and other tools
RUN apk add --no-cache git ca-certificates tzdata curl

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8090

# Command to run (can be overridden)
CMD ["go", "run", "main.go"]