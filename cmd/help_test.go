package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpMentionsScopeForRootAndSearch(t *testing.T) {
	needle := "fuzzy search results"

	rootOut, err := executeForTest("--help")
	if err != nil {
		t.Fatalf("root help error: %v", err)
	}
	if !strings.Contains(rootOut, needle) {
		t.Fatalf("root help missing notice: %s", rootOut)
	}
	if strings.Contains(rootOut, "agent-browser") || strings.Contains(rootOut, "playwright") {
		t.Fatalf("root help should not suggest browser tools: %s", rootOut)
	}
	if !strings.Contains(rootOut, "search first, then read") {
		t.Fatalf("root help missing workflow guidance: %s", rootOut)
	}
	if !strings.Contains(rootOut, "If you already have a URL, call read directly") {
		t.Fatalf("root help missing direct read guidance: %s", rootOut)
	}

	searchOut, err := executeForTest("search", "--help")
	if err != nil {
		t.Fatalf("search help error: %v", err)
	}
	if !strings.Contains(searchOut, needle) {
		t.Fatalf("search help missing notice: %s", searchOut)
	}
	if !strings.Contains(searchOut, "key=value") {
		t.Fatalf("search help missing error format: %s", searchOut)
	}
	if !strings.Contains(searchOut, "# | title | url | content | template") {
		t.Fatalf("search help missing output structure: %s", searchOut)
	}
	if strings.Contains(rootOut, "--verbose") {
		t.Fatalf("root help should not include verbose flag: %s", rootOut)
	}
	for _, flag := range []string{"--auth-mode", "--auth-header", "--auth-api-key", "--auth-username", "--auth-password"} {
		if !strings.Contains(searchOut, flag) {
			t.Fatalf("search help missing %s: %s", flag, searchOut)
		}
	}
	if !strings.Contains(searchOut, "base64") {
		t.Fatalf("search help missing base64 auth guidance: %s", searchOut)
	}
}

func executeForTest(args ...string) (string, error) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	return buf.String(), err
}
