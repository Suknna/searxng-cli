package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestHelpMentionsScopeForRootAndSearch(t *testing.T) {
	needle := "模糊搜索结果"

	rootOut, err := executeForTest("--help")
	if err != nil {
		t.Fatalf("root help error: %v", err)
	}
	if !strings.Contains(rootOut, needle) {
		t.Fatalf("root help missing notice: %s", rootOut)
	}

	searchOut, err := executeForTest("search", "--help")
	if err != nil {
		t.Fatalf("search help error: %v", err)
	}
	if !strings.Contains(searchOut, needle) {
		t.Fatalf("search help missing notice: %s", searchOut)
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
