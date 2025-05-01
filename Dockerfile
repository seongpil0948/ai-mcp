# Path: Dockerfile
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./
COPY modules/*/go.mod modules/*/go.sum ./modules/

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/mcpserver ./cmd/mcpserver

# Create final lightweight container
FROM alpine:3.19

# Add ca certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the built executable from the builder stage
COPY --from=builder /app/mcpserver .

# Copy configurations
COPY configs/ ./configs/

# Set environment variables
ENV GIN_MODE=release

# Expose HTTP port (if running in HTTP mode)
EXPOSE 8080

# Command to run
ENTRYPOINT ["/app/mcpserver"]