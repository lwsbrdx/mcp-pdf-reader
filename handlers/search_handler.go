package handlers

import (
	"context"
	"mcp-pdf-reader/internal"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchInput struct {
	Path          string `json:"path" jsonschema:"Path to the PDF file to search in"`
	Query         string `json:"query" jsonschema:"Text to search for in the PDF"`
	Page          *int   `json:"page,omitempty" jsonschema:"Optional: specific page number to search in"`
	CaseSensitive bool   `json:"case_sensitive,omitempty" jsonschema:"Optional: whether to perform case-sensitive search (default: false)"`
	ContextLength int    `json:"context_length,omitempty" jsonschema:"Optional: number of characters of context around each match (default: 50)"`
}

type SearchOutput struct {
	Matches    []internal.Match `json:"matches"`
	TotalCount int              `json:"total_count"`
}

func SearchHandler(ctx context.Context, request *mcp.CallToolRequest, input SearchInput) (
	*mcp.CallToolResult,
	SearchOutput,
	error,
) {
	// Set default context length if not provided
	contextLen := input.ContextLength
	if contextLen == 0 {
		contextLen = 50
	}

	matches := internal.Search(
		input.Path,
		input.Query,
		input.Page,
		input.CaseSensitive,
		contextLen,
	)

	// Ensure matches is never nil (return empty array instead)
	if matches == nil {
		matches = []internal.Match{}
	}

	return &mcp.CallToolResult{}, SearchOutput{
		Matches:    matches,
		TotalCount: len(matches),
	}, nil
}
