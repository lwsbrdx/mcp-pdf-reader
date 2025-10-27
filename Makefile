.PHONY: build docker-build docker-run docker-clean all

# Build local binary
build:
	@go build -o pdf-server main.go

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t mcp-pdf-reader:latest .
	@echo "✅ Image built successfully"
	@docker images mcp-pdf-reader:latest

# Run Docker container (interactive test)
docker-run:
	@docker run --rm -i mcp-pdf-reader:latest

# Clean Docker images
docker-clean:
	@docker rmi mcp-pdf-reader:latest

# Build everything
all: build docker-build
