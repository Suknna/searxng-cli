package auth

import (
	"encoding/base64"
	"net/http"
	"testing"
)

func TestApplyAPIKeyDecoded(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.org", nil)
	err := Apply(req, Options{Mode: "api_key", APIKeyHeader: "X-Custom-Key", APIKey: base64.StdEncoding.EncodeToString([]byte("secret"))})
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	if got := req.Header.Get("X-Custom-Key"); got != "secret" {
		t.Fatalf("header = %q", got)
	}
}

func TestApplyBasicDecoded(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.org", nil)
	err := Apply(req, Options{Mode: "basic", Username: base64.StdEncoding.EncodeToString([]byte("alice")), Password: base64.StdEncoding.EncodeToString([]byte("p@ss"))})
	if err != nil {
		t.Fatalf("Apply error: %v", err)
	}
	u, p, ok := req.BasicAuth()
	if !ok || u != "alice" || p != "p@ss" {
		t.Fatalf("basic auth = %v %q %q", ok, u, p)
	}
}

func TestApplyInvalidBase64(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "https://example.org", nil)
	err := Apply(req, Options{Mode: "api_key", APIKeyHeader: "X-Api-Key", APIKey: "%%%"})
	if err == nil {
		t.Fatal("expected error")
	}
}
