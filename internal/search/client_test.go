package search

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"searxng-cli/internal/auth"
)

func TestFetch_UsesOnlyQAndFormatJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			t.Fatalf("path = %s", r.URL.Path)
		}
		vals := r.URL.Query()
		if vals.Get("q") != "searxng" {
			t.Fatalf("q = %q", vals.Get("q"))
		}
		if vals.Get("format") != "json" {
			t.Fatalf("format = %q", vals.Get("format"))
		}
		if len(vals) != 2 {
			t.Fatalf("unexpected query params: %v", vals)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[{"title":"t","url":"u","content":"c"}]}`))
	}))
	defer ts.Close()

	got, err := Fetch(context.Background(), ts.URL, "searxng", 2*time.Second, auth.Options{})
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
	if len(got) != 1 || got[0].Title != "t" {
		t.Fatalf("results = %#v", got)
	}
}

func TestFetch_AppliesAPIKeyAuth(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("X-Test-Key"); got != "secret" {
			t.Fatalf("header=%q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"results":[]}`))
	}))
	defer ts.Close()

	_, err := Fetch(context.Background(), ts.URL, "searxng", time.Second, auth.Options{
		Mode:         "api_key",
		APIKeyHeader: "X-Test-Key",
		APIKey:       base64.StdEncoding.EncodeToString([]byte("secret")),
	})
	if err != nil {
		t.Fatalf("Fetch error: %v", err)
	}
}
