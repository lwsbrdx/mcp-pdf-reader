package internal

import (
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/klippa-app/go-pdfium/references"
	"github.com/klippa-app/go-pdfium/requests"
	"github.com/klippa-app/go-pdfium/webassembly"
	pdfium "github.com/klippa-app/go-pdfium"
)

var (
	pool     pdfium.Pool
	poolOnce sync.Once
)

// initPool initializes the PDFium WebAssembly pool once
func initPool() {
	poolOnce.Do(func() {
		var err error
		pool, err = webassembly.Init(webassembly.Config{
			MinIdle:  1,
			MaxIdle:  3,
			MaxTotal: 5,
		})
		if err != nil {
			log.Fatal("Failed to initialize PDFium pool:", err)
		}
	})
}

func Read(filepath string) string {
	initPool()

	instance, err := pool.GetInstance(10000)
	if err != nil {
		panic(err)
	}
	defer instance.Close()

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		FilePath: &filepath,
	})
	if err != nil {
		panic(err)
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{Document: doc.Document})
	if err != nil {
		panic(err)
	}

	var content strings.Builder
	for i := 0; i < pageCount.PageCount; i++ {
		pageText := extractPageText(instance, doc.Document, i)
		content.WriteString(pageText)
		if i < pageCount.PageCount-1 {
			content.WriteString("\n\n")
		}
	}

	return content.String()
}

// cleanText normalizes whitespace (go-pdfium produces clean text, so minimal processing needed)
func cleanText(text string) string {
	// Just normalize any excessive line breaks
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	return strings.TrimSpace(text)
}

type Match struct {
	Page       int    `json:"page"`
	Context    string `json:"context"`
	MatchStart int    `json:"match_start"`
	MatchEnd   int    `json:"match_end"`
}

func Search(path, query string, page *int, caseSensitive bool, contextLen int) []Match {
	initPool()

	instance, err := pool.GetInstance(10000)
	if err != nil {
		panic(err)
	}
	defer instance.Close()

	doc, err := instance.OpenDocument(&requests.OpenDocument{
		FilePath: &path,
	})
	if err != nil {
		panic(err)
	}
	defer instance.FPDF_CloseDocument(&requests.FPDF_CloseDocument{Document: doc.Document})

	pageCount, err := instance.FPDF_GetPageCount(&requests.FPDF_GetPageCount{Document: doc.Document})
	if err != nil {
		panic(err)
	}

	totalPages := pageCount.PageCount
	matches := make([]Match, 0)

	startPage, endPage := 1, totalPages
	if page != nil {
		startPage, endPage = *page, *page
	}

	for pageNum := startPage; pageNum <= endPage; pageNum++ {
		pageText := extractPageText(instance, doc.Document, pageNum-1) // go-pdfium uses 0-based indexing

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

func extractPageText(instance pdfium.Pdfium, document references.FPDF_DOCUMENT, pageIndex int) string {
	page, err := instance.FPDF_LoadPage(&requests.FPDF_LoadPage{
		Document: document,
		Index:    pageIndex,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to load page %d: %v", pageIndex, err))
	}
	defer instance.FPDF_ClosePage(&requests.FPDF_ClosePage{Page: page.Page})

	textPage, err := instance.FPDFText_LoadPage(&requests.FPDFText_LoadPage{
		Page: requests.Page{ByReference: &page.Page},
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to load text page %d: %v", pageIndex, err))
	}
	defer instance.FPDFText_ClosePage(&requests.FPDFText_ClosePage{TextPage: textPage.TextPage})

	charCount, err := instance.FPDFText_CountChars(&requests.FPDFText_CountChars{TextPage: textPage.TextPage})
	if err != nil {
		panic(fmt.Sprintf("Failed to count chars on page %d: %v", pageIndex, err))
	}

	pageText, err := instance.FPDFText_GetText(&requests.FPDFText_GetText{
		TextPage:   textPage.TextPage,
		StartIndex: 0,
		Count:      charCount.Count,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to get text from page %d: %v", pageIndex, err))
	}

	return cleanText(pageText.Text)
}
