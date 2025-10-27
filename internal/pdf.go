package internal

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

func Read(filepath string) string {
	pdf.DebugOn = true

	file, reader, err := pdf.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var buf bytes.Buffer
	b, err := reader.GetPlainText()
	if err != nil {
		panic(err)
	}
	buf.ReadFrom(b)
	content := strings.TrimSpace(buf.String())

	return content
}

type Match struct {
	Page       int    `json:"page"`
	Context    string `json:"context"`
	MatchStart int    `json:"match_start"`
	MatchEnd   int    `json:"match_end"`
}

func Search(path, query string, page *int, caseSensitive bool, contextLen int) []Match {
	file, reader, err := pdf.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	totalPages := reader.NumPage()
	matches := make([]Match, 0) // Initialize with empty slice instead of nil

	startPage, endPage := 1, totalPages
	if page != nil {
		startPage, endPage = *page, *page
	}

	for pageNum := startPage; pageNum <= endPage; pageNum++ {
		pageText := extractPageText(reader, pageNum)
		// extractSentences ??
		// iterate over sentences to put previous sentences in the context for the IA

		searchText := pageText
		searchQuery := query
		if !caseSensitive {
			searchText = strings.ToLower(pageText)
			searchQuery = strings.ToLower(query)
		}

		idx := 0
		for {
			pos := strings.Index(searchText[idx:], searchQuery)
			if pos == -1 {
				break
			}

			actualPos := idx + pos

			start := max(0, actualPos-contextLen)
			end := min(len(pageText), actualPos+len(query)+contextLen)
			context := pageText[start:end]

			matches = append(matches, Match{
				Page:       pageNum,
				Context:    context,
				MatchStart: actualPos - start,
				MatchEnd:   actualPos - start + len(query),
			})

			idx = actualPos + len(query)
		}
	}

	return matches
}

func extractPageText(reader *pdf.Reader, pageNum int) string {
	page := reader.Page(pageNum)
	if page.V.IsNull() {
		panic(fmt.Sprintf("Page not found, file has %d pages", reader.NumPage()))
	}

	content, err := page.GetPlainText(nil)
	if err != nil {
		panic(err)
	}

	return content
}
