package internal_test

import (
	"fmt"
	"mcp-pdf-reader/internal"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		filepath string
		want     string
	}{
		// TODO: Add test cases.
		{
			name: "Read dummy file",
			filepath: "../samples/pdf-sample-0.pdf",
			want: "Dummy PDF file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := internal.Read(tt.filepath)
			
			fmt.Println(got)
			fmt.Println(tt.want)

			if got != tt.want {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

