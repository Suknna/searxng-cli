package render

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"searxng-cli/internal/search"
)

const (
	titleMax    = 80
	urlMax      = 120
	contentMax  = 160
	templateMax = 200
)

func MarkdownTable(results []search.Result, tmplText string, limit int) (string, error) {
	if limit <= 0 {
		limit = len(results)
	}
	if limit > len(results) {
		limit = len(results)
	}

	tmpl, err := template.New("row").Parse(tmplText)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	b.WriteString("| # | title | url | content | template |\n")
	b.WriteString("| --- | --- | --- | --- | --- |\n")

	for i := 0; i < limit; i++ {
		r := results[i]
		rendered, err := renderTemplate(tmpl, r)
		if err != nil {
			return "", err
		}

		fmt.Fprintf(&b, "| %d | %s | %s | %s | %s |\n",
			i+1,
			normalizeCell(r.Title, titleMax),
			normalizeCell(r.URL, urlMax),
			normalizeCell(r.Content, contentMax),
			normalizeCell(rendered, templateMax),
		)
	}

	return b.String(), nil
}

func renderTemplate(t *template.Template, r search.Result) (string, error) {
	var out bytes.Buffer
	if err := t.Execute(&out, r); err != nil {
		return "", err
	}
	return out.String(), nil
}

func normalizeCell(s string, max int) string {
	s = strings.NewReplacer("\n", " ", "\r", " ", "\t", " ").Replace(s)
	s = strings.Join(strings.Fields(s), " ")
	s = strings.ReplaceAll(s, "|", "\\|")

	if max > 0 {
		r := []rune(s)
		if len(r) > max {
			s = string(r[:max]) + "..."
		}
	}
	return s
}
