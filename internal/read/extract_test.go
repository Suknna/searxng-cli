package read

import (
	"strings"
	"testing"
)

func TestExtractContentMarkdownIsSanitized(t *testing.T) {
	html := `<html><body><article><h1>Title</h1><p>Hello <span>world</span></p></article></body></html>`
	out, err := ExtractContent(html, "https://example.org", "markdown")
	if err != nil {
		t.Fatalf("ExtractContent error: %v", err)
	}
	if strings.Contains(out, "<span>") || strings.Contains(out, "<article>") {
		t.Fatalf("markdown still has html: %q", out)
	}
}

func TestExtractContentTextHasNoHTML(t *testing.T) {
	html := `<html><body><article><h1>Title</h1><p>Hello <span>world</span></p></article></body></html>`
	out, err := ExtractContent(html, "https://example.org", "text")
	if err != nil {
		t.Fatalf("ExtractContent error: %v", err)
	}
	if strings.Contains(out, "<") || strings.Contains(out, ">") {
		t.Fatalf("text still has html: %q", out)
	}
}
