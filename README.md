# MCP PDF Reader

> ⚠️ **Work in Progress**: This project is under active development. Features and APIs may change.

A Model Context Protocol (MCP) server in Go for reading and extracting text from PDF files. This server enables AI assistants like Claude to read and analyze PDF documents.

## Features

- ✅ Read full text content from PDF files
- 🔄 Read text content from specific pages or page ranges *(planned)*
- 🔄 Read PDF metadata (author, title, creation date, etc.) *(planned)*
- 🔄 Get the total page count of a PDF *(planned)*
- 🔄 Process multiple PDF sources (local paths or URLs) in a single request *(planned)*

## Installation

### Prerequisites

- Go 1.21 or higher
- Git

### Build from source

```bash
# Clone the repository
git clone https://github.com/lwsbrdx/mcp-pdf.git
cd mcp-pdf

# Build the server
make build
# or
go build -o pdf-server main.go
```

### Using Docker

```bash
# Build the Docker image
docker build -t mcp-pdf-reader:latest .

# Run with Docker (interactive mode for testing)
echo '{"jsonrpc":"2.0","id":1,"method":"initialize",...}' | docker run --rm -i mcp-pdf-reader:latest
```

The Docker image is optimized using multi-stage builds and weighs only ~20MB.

## Usage

### With Claude Code

Add the server to your Claude Code configuration file (`~/.config/claude-code/mcp_config.json`):

```json
{
  "mcpServers": {
    "pdf-reader": {
      "command": "/path/to/mcp-pdf/pdf-server"
    }
  }
}
```

Restart Claude Code, and the PDF reading tool will be available.

### Testing the server

You can test the server using the official MCP client:

```bash
# Install the Go SDK examples
git clone https://github.com/modelcontextprotocol/go-sdk.git

# Use the listfeatures client
cd go-sdk/examples/client/listfeatures
go run main.go /path/to/pdf-server
```

## Available Tools

### `read_pdf`

Reads the full text content from a PDF file.

**Parameters:**
- `path` (string): Path to the PDF file

**Returns:**
- `content` (string): Extracted text content from the PDF

**Example usage:**
```json
{
  "name": "read_pdf",
  "arguments": {
    "path": "/path/to/document.pdf"
  }
}
```

## Development

### Running tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests for a specific package
go test ./internal -v
```

### Project structure

```
mcp-pdf-reader/
├── main.go              # MCP server entry point
├── handlers/            # MCP tool handlers
│   └── read_handler.go  # Handler for read_pdf tool
├── internal/            # Internal packages
│   ├── pdf.go          # PDF reading logic
│   └── pdf_test.go     # Tests
├── samples/            # Sample PDF files for testing
└── TODO.md            # Project guide and documentation
```

## Roadmap

### Phase 1: Core Features (Current)
- [x] Basic PDF text extraction
- [x] MCP server implementation
- [x] Single file reading support
- [ ] Error handling improvements

### Phase 2: Advanced Reading
- [ ] Read specific pages or page ranges
- [ ] Extract PDF metadata (author, title, dates, etc.)
- [ ] Get total page count
- [ ] Support for password-protected PDFs

### Phase 3: Multi-source Support
- [ ] Process multiple PDF files in one request
- [ ] Support for PDF URLs (download and read)
- [ ] Caching mechanism for remote PDFs

### Phase 4: Performance & Deployment
- [x] Docker image for easy deployment
- [ ] Performance benchmarks
- [ ] Memory optimization for large PDFs
- [ ] Streaming support for very large files

### Phase 5: Additional Features
- [ ] Search text within PDFs
- [ ] Extract images from PDFs
- [ ] PDF structure analysis (TOC, bookmarks)
- [ ] OCR support for scanned PDFs

## Docker Support

### Building the image

```bash
docker build -t mcp-pdf-reader:latest .
```

### Configuration with Claude Code

To use the Docker image with Claude Code, update your MCP configuration:

```json
{
  "mcpServers": {
    "pdf-reader": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "mcp-pdf-reader:latest"]
    }
  }
}
```

### Image details

- **Base image**: Alpine Linux (minimal, secure)
- **Size**: ~20MB (multi-stage build)
- **User**: Runs as non-root user `mcp` (UID 1000)
- **Transport**: stdio (standard input/output)

## Benchmarks *(Planned)*

Performance benchmarks will be added once the core features are stable.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [modelcontextprotocol/go-sdk](https://github.com/modelcontextprotocol/go-sdk)
- PDF parsing using [ledongthuc/pdf](https://github.com/ledongthuc/pdf) and [pdfcpu](https://github.com/pdfcpu/pdfcpu)

## Related Projects

- [Model Context Protocol](https://modelcontextprotocol.io) - Official MCP documentation
- [Claude Code](https://claude.com/claude-code) - AI-powered CLI tool
