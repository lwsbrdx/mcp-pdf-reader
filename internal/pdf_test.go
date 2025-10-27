package internal_test

import (
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
		{
			name:     "Read dummy file",
			filepath: "../samples/pdf-sample-0.pdf",
			want:     "Dummy PDF file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := internal.Read(tt.filepath)

			if got != tt.want {
				t.Errorf("Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		path          string
		query         string
		page          *int
		caseSensitive bool
		contextLen    int
		want          []internal.Match
	}{
		{
			name:          "Search in lorem file",
			path:          "../samples/pdf-sample-lorem.pdf",
			query:         "Fusce nec tellus sed",
			caseSensitive: true,
			contextLen:    50,
			want: []internal.Match{
				{
					1,
					"\nimperdiet. Duis sagittis ipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris \nmassa. Vestibulum lac",
					50,
					70,
				},
				{
					2,
					"mperdiet. Duis sagittis ipsum. \nPraesent mauris. \nFusce nec tellus sed augue semper porta. Mauris massa. Vestibulum laci",
					50,
					70,
				},
				{
					2,
					"imperdiet. Duis sagittis \nipsum. Praesent mauris. Fusce nec tellus sed augue semper porta. Mauris massa. Vestibulum \nlac",
					50,
					70,
				},
				{
					3,
					"mperdiet. Duis sagittis ipsum. \nPraesent mauris. \nFusce nec tellus sed augue semper porta. Mauris massa. Vestibulum laci",
					50,
					70,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatches := internal.Search(tt.path, tt.query, tt.page, tt.caseSensitive, tt.contextLen)
			var found bool

			for _, match := range gotMatches {
				// fmt.Println(match)
				for _, want := range tt.want {
					// fmt.Println(want)
					if match.Page == want.Page &&
						match.MatchStart == want.MatchStart &&
						match.MatchEnd == want.MatchEnd &&
						match.Context == want.Context {
						found = true
					}
				}

				if found == false {
					t.Errorf("Search() = %v, want %v", match, tt.want)
				}
			}
		})
	}
}
