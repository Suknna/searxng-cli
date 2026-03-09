package read

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchHTMLOK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<html><body>ok</body></html>"))
	}))
	defer ts.Close()

	res, err := FetchHTML(context.Background(), ts.URL, Options{Timeout: time.Second})
	if err != nil {
		t.Fatalf("FetchHTML error: %v", err)
	}
	if !strings.Contains(res.HTML, "ok") {
		t.Fatalf("unexpected html: %q", res.HTML)
	}
}

func TestFetchHTMLUnsupportedContentType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"x":1}`))
	}))
	defer ts.Close()

	_, err := FetchHTML(context.Background(), ts.URL, Options{Timeout: time.Second})
	if err == nil {
		t.Fatal("expected error")
	}
}
