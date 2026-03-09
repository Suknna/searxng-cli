package read

import (
	"regexp"
	"strings"
)

var tagPattern = regexp.MustCompile(`<[^>]+>`)

func sanitizeMarkdown(md string) string {
	cleaned := tagPattern.ReplaceAllString(md, "")
	return normalizeWhitespace(cleaned)
}

func sanitizeText(text string) string {
	cleaned := tagPattern.ReplaceAllString(text, "")
	return normalizeWhitespace(cleaned)
}

func normalizeWhitespace(s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
