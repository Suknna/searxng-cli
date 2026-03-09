package cmd

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadCommandDefaultMarkdown(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			_, _ = w.Write([]byte("User-agent: *\nAllow: /\n"))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body><article><h1>Title</h1><p>Hello <span>world</span></p></article></body></html>`))
	}))
	defer ts.Close()

	out, err := executeForTest("read", ts.URL)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if strings.Contains(out, "<span>") {
		t.Fatalf("output still has html tags: %s", out)
	}
}

func TestReadCommandTextFormat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			_, _ = w.Write([]byte("User-agent: *\nAllow: /\n"))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		_, _ = w.Write([]byte(`<html><body><article><h1>Title</h1><p>Hello <span>world</span></p></article></body></html>`))
	}))
	defer ts.Close()

	out, err := executeForTest("read", ts.URL, "--format", "text")
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if strings.Contains(out, "<") || strings.Contains(out, ">") {
		t.Fatalf("text output still has html tags: %s", out)
	}
}
