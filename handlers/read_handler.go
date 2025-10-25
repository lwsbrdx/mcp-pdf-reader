package handlers

import (
	"context"
	"mcp-pdf-reader/internal"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ReadInput struct {
	Path string `json:"path" jsonschema:"Chemin vers le pdf cibler pour en récupérer le contenu"`
}

type ReadOutput struct {
	Content string `json:"content" jsonschema:"Contenu du pdf sous forme de string en markdown"`
}

func ReadHandler(ctx context.Context, request *mcp.CallToolRequest, input ReadInput) (
	*mcp.CallToolResult,
	ReadOutput,
	error,
) {
	content := internal.Read(input.Path)

	return &mcp.CallToolResult{}, ReadOutput{Content: content}, nil
}
