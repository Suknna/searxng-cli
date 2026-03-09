package search

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"searxng-cli/internal/apperr"
	"searxng-cli/internal/auth"
)

func TestFetchNon2xxReturnsHTTPStatusError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad gateway", http.StatusBadGateway)
	}))
	defer ts.Close()

	_, err := Fetch(context.Background(), ts.URL, "q", 2*time.Second, auth.Options{})
	if err == nil {
		t.Fatal("expected error")
	}

	var statusErr *apperr.HTTPStatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("unexpected error type: %T", err)
	}
	if statusErr.StatusCode != 502 {
		t.Fatalf("status=%d", statusErr.StatusCode)
	}
}
