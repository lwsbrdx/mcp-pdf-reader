FROM golang:1.25-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o pdf-server main.go

# Final stage - minimal image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 mcp && \
    adduser -D -u 1000 -G mcp mcp

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/pdf-server .

# Change ownership
RUN chown -R mcp:mcp /app

# Switch to non-root user
USER mcp

# Expose no ports (stdio communication)
# The server communicates via stdin/stdout

# Set entrypoint
ENTRYPOINT ["/app/pdf-server"]

