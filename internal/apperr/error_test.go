package apperr

import "testing"

func TestDefaultsForCode(t *testing.T) {
	app := FromError(assertErr("no such host"))
	if app.Code != CodeNetworkDNS {
		t.Fatalf("code = %s", app.Code)
	}
	if app.Message == "" {
		t.Fatal("message should not be empty")
	}
	if app.Hint == "" {
		t.Fatal("hint should not be empty")
	}
}

type assertErr string

func (e assertErr) Error() string { return string(e) }
