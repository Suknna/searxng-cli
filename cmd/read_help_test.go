package cmd

import (
	"strings"
	"testing"
)

func TestReadHelp(t *testing.T) {
	out, err := executeForTest("read", "--help")
	if err != nil {
		t.Fatalf("read help error: %v", err)
	}
	if !strings.Contains(out, "--format") {
		t.Fatalf("missing format flag: %s", out)
	}
	if !strings.Contains(out, "base64") {
		t.Fatalf("missing base64 guidance: %s", out)
	}
}
