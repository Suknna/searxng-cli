package apperr

import (
	"context"
	"errors"
	"testing"

	"searxng-cli/internal/read"
)

func TestClassifyVariants(t *testing.T) {
	tests := []struct {
		name string
		err  error
		code Code
	}{
		{name: "dns", err: errors.New("lookup x: no such host"), code: CodeNetworkDNS},
		{name: "timeout", err: context.DeadlineExceeded, code: CodeNetworkTimeout},
		{name: "reset", err: errors.New("EOF"), code: CodeNetworkReset},
		{name: "tls", err: errors.New("tls: handshake failure"), code: CodeTLSHandshake},
		{name: "decode", err: &DecodeError{Err: errors.New("invalid character")}, code: CodeResponseDecode},
		{name: "config", err: &ConfigError{Err: errors.New("context missing")}, code: CodeConfigInvalid},
		{name: "template", err: &TemplateError{Err: errors.New("bad template")}, code: CodeTemplateBad},
		{name: "http", err: &HTTPStatusError{StatusCode: 503, Status: "503 Service Unavailable"}, code: CodeHTTPNon2xx},
		{name: "auth decode", err: &AuthDecodeError{Err: errors.New("illegal base64")}, code: CodeAuthDecodeInvalid},
		{name: "auth config", err: &AuthConfigError{Err: errors.New("missing auth")}, code: CodeAuthConfigInvalid},
		{name: "robots", err: &read.RobotsDisallowedError{URL: "https://example.org"}, code: CodeRobotsDisallowed},
		{name: "content type", err: &read.ContentTypeUnsupportedError{ContentType: "application/json"}, code: CodeContentTypeUnsupported},
		{name: "extract empty", err: &read.ExtractEmptyError{}, code: CodeExtractEmpty},
		{name: "fetch failed", err: &read.FetchFailedError{Err: errors.New("dial error")}, code: CodeFetchFailed},
		{name: "unknown", err: errors.New("boom"), code: CodeInternal},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := FromError(tc.err)
			if got.Code != tc.code {
				t.Fatalf("code=%s", got.Code)
			}
		})
	}
}
