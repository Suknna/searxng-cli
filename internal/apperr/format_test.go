package apperr

import (
	"strings"
	"testing"
)

func TestRenderKV(t *testing.T) {
	err := AppError{
		Code:      CodeNetworkReset,
		Retryable: true,
		Message:   "connection closed by peer",
		Hint:      "check proxy path",
		Cause:     "EOF",
		Meta:      map[string]string{"server": "https://example.org"},
	}
	line := RenderKV(err)
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
