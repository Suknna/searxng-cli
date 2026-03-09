package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadEffectivePriority(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yml")

	content := `apiVersion: searxng-cli/v1
kind: Config
current-context: default
contexts:
  default:
    server: "https://example.local"
    timeout: "20s"
    limit: 7
    template: "Title={{.Title}}"`

	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	override := Overrides{
		ConfigPath: cfgPath,
		Server:     "https://flag.local",
		Timeout:    3 * time.Second,
		Limit:      intPtr(4),
		Template:   strPtr("Title={{.URL}}"),
	}

	eff, _, err := LoadEffective(override)
	if err != nil {
		t.Fatalf("LoadEffective error: %v", err)
	}

	if eff.Server != "https://flag.local" {
		t.Fatalf("server = %q", eff.Server)
	}
	if eff.Timeout != 3*time.Second {
		t.Fatalf("timeout = %s", eff.Timeout)
	}
	if eff.Limit != 4 {
		t.Fatalf("limit = %d", eff.Limit)
	}
	if eff.Template != "Title={{.URL}}" {
		t.Fatalf("template = %q", eff.Template)
	}
}

func TestDefaultConfigServer(t *testing.T) {
	cfg := DefaultConfig()
	if got := cfg.Contexts["default"].Server; got != "https://searx.example.org/" {
		t.Fatalf("default server = %q", got)
	}
}

func TestDefaultConfigAuthModeNone(t *testing.T) {
	cfg := DefaultConfig()
	if got := cfg.Contexts["default"].Auth.Mode; got != "none" {
		t.Fatalf("auth mode = %q", got)
	}
}

func TestLoadEffectiveAuthPriority(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yml")

	content := `apiVersion: searxng-cli/v1
kind: Config
current-context: default
contexts:
  default:
    server: "https://example.local"
    timeout: "20s"
    limit: 7
    template: "Title={{.Title}}"
    auth:
      mode: "api_key"
      api_key_header: "X-From-Config"
      api_key: "Y29uZmlnLWtleQ=="`

	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv("SEARXNG_AUTH_MODE", "basic")
	t.Setenv("SEARXNG_AUTH_HEADER", "X-From-Env")
	t.Setenv("SEARXNG_AUTH_API_KEY", "ZW52LWtleQ==")

	override := Overrides{
		ConfigPath: cfgPath,
		AuthMode:   strPtr("api_key"),
		AuthHeader: strPtr("X-From-Flag"),
		AuthAPIKey: strPtr("ZmxhZy1rZXk="),
	}

	eff, _, err := LoadEffective(override)
	if err != nil {
		t.Fatalf("LoadEffective error: %v", err)
	}

	if eff.AuthMode != "api_key" {
		t.Fatalf("auth mode = %q", eff.AuthMode)
	}
	if eff.AuthHeader != "X-From-Flag" {
		t.Fatalf("auth header = %q", eff.AuthHeader)
	}
	if eff.AuthAPIKey != "ZmxhZy1rZXk=" {
		t.Fatalf("auth api key = %q", eff.AuthAPIKey)
	}
}

func intPtr(v int) *int       { return &v }
func strPtr(v string) *string { return &v }
