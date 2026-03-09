package read

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	readability "github.com/mackee/go-readability"
)

func ExtractContent(html, baseURL, format string) (string, error) {
	article, err := readability.Extract(html, readability.DefaultOptions())
	if err != nil {
		return "", err
	}

	var markdown string
	if article.Root != nil {
		markdown = readability.ToMarkdown(article.Root)
	}
	if strings.TrimSpace(markdown) == "" {
		markdown, err = fallbackMarkdown(html)
		if err != nil {
			return "", err
		}
	}

	cleanMarkdown := sanitizeMarkdown(markdown)
	cleanText := sanitizeText(cleanMarkdown)
	if cleanText == "" {
		return "", &ExtractEmptyError{}
	}

	if strings.EqualFold(format, "text") {
		return cleanText, nil
	}
	return cleanMarkdown, nil
}

func fallbackMarkdown(html string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}
	body := strings.TrimSpace(doc.Find("body").Text())
	if body == "" {
		return "", fmt.Errorf("empty body")
	}
	return body, nil
}
