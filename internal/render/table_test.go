package render

import (
	"strings"
	"testing"

	"searxng-cli/internal/search"
)

func TestRenderTable_ColumnsAndNormalize(t *testing.T) {
	rows, err := MarkdownTable([]search.Result{{
		Title:   "a|b\n",
		URL:     "https://x.example/path",
		Content: "hello\tworld",
	}}, "Title={{.Title}} URL={{.URL}} Content={{.Content}}", 10)
	if err != nil {
		t.Fatalf("MarkdownTable error: %v", err)
	}

	if !strings.Contains(rows, "| # | title | url | content | template |") {
		t.Fatalf("missing header: %s", rows)
	}
	if !strings.Contains(rows, "a\\|b") {
		t.Fatalf("pipe not escaped: %s", rows)
	}
	if strings.Contains(rows, "\n ") {
		t.Fatalf("newline not normalized: %s", rows)
	}
}
