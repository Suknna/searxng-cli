package read

import (
	"strings"
	"testing"
)

func TestSanitizeMarkdownRemovesHTMLTags(t *testing.T) {
	input := "# Title\n<div>noise</div>\nParagraph <span>x</span>"
	got := sanitizeMarkdown(input)
	if strings.Contains(got, "<div>") || strings.Contains(got, "<span>") {
		t.Fatalf("contains html tags: %q", got)
	}
}
