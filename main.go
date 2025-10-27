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
		Name:        "read_pdf",
		Description: "Reads the full text content from a PDF file and returns it as a string",
	}, handlers.ReadHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_in_pdf",
		Description: "Searches for specific text in a PDF file and returns all matches with their page numbers and surrounding context",
	}, handlers.SearchHandler)

	log.Println("Starting server...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}

