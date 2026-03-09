package cmd

import (
	"errors"
	"strings"
	"testing"

	"searxng-cli/internal/apperr"
)

func TestRenderCLIError(t *testing.T) {
	err := apperr.Annotate(errors.New("EOF"), map[string]string{"server": "https://example.org"})
	line := renderCLIError(err)

	if !strings.Contains(line, "code=NETWORK_RESET") {
		t.Fatalf("line=%s", line)
	}
	if !strings.Contains(line, "retryable=true") {
		t.Fatalf("line=%s", line)
	}
	if !strings.Contains(line, "server=\"https://example.org\"") {
		t.Fatalf("line=%s", line)
	}
}
