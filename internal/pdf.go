package internal

import (
	"bytes"
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
