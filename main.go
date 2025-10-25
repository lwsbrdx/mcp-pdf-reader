package main

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"mcp-pdf-reader/handlers"
)

func main() {
	server := mcp.NewServer(&mcp.Implementation{Name: "pdf-reader", Version: "v1.0.0"}, nil)

	mcp.AddTool(server, &mcp.Tool{
		Name: "read_pdf",
		Description: "Reads the content of a PDF file and gives the result as a string",
	}, handlers.ReadHandler);

	log.Println("Starting server...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

