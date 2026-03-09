package apperr

import (
	"context"
	"errors"
	"net"
	"strings"

	"searxng-cli/internal/read"
)

func FromError(err error) AppError {
	if err == nil {
		return AppError{Code: CodeInternal, Message: "unknown error", Hint: "check logs"}
	}

	meta := map[string]string{}
	var wrapped *withMeta
	if errors.As(err, &wrapped) {
		for k, v := range wrapped.meta {
			meta[k] = v
		}
	}

	if app, ok := err.(AppError); ok {
		if app.Meta == nil {
			app.Meta = meta
		}
		return fillDefaults(app, err.Error())
	}

	var httpErr *HTTPStatusError
	if errors.As(err, &httpErr) {
		retryable := httpErr.StatusCode >= 500
		meta["status"] = httpErr.Status
		return fillDefaults(AppError{Code: CodeHTTPNon2xx, Retryable: retryable, Meta: meta}, err.Error())
	}

	var cfgErr *ConfigError
	if errors.As(err, &cfgErr) {
		return fillDefaults(AppError{Code: CodeConfigInvalid, Meta: meta}, err.Error())
	}

	var tmplErr *TemplateError
	if errors.As(err, &tmplErr) {
		return fillDefaults(AppError{Code: CodeTemplateBad, Meta: meta}, err.Error())
	}

	var authDecodeErr *AuthDecodeError
	if errors.As(err, &authDecodeErr) {
		return fillDefaults(AppError{Code: CodeAuthDecodeInvalid, Meta: meta}, err.Error())
	}

	var authConfigErr *AuthConfigError
	if errors.As(err, &authConfigErr) {
		return fillDefaults(AppError{Code: CodeAuthConfigInvalid, Meta: meta}, err.Error())
	}

	var robotsErr *read.RobotsDisallowedError
	if errors.As(err, &robotsErr) {
		return fillDefaults(AppError{Code: CodeRobotsDisallowed, Meta: meta}, err.Error())
	}

	var ctErr *read.ContentTypeUnsupportedError
	if errors.As(err, &ctErr) {
		return fillDefaults(AppError{Code: CodeContentTypeUnsupported, Meta: meta}, err.Error())
	}

	var emptyErr *read.ExtractEmptyError
	if errors.As(err, &emptyErr) {
		return fillDefaults(AppError{Code: CodeExtractEmpty, Meta: meta}, err.Error())
	}

	var fetchErr *read.FetchFailedError
	if errors.As(err, &fetchErr) {
		return fillDefaults(AppError{Code: CodeFetchFailed, Retryable: true, Meta: meta}, err.Error())
	}

	var decErr *DecodeError
	if errors.As(err, &decErr) {
		return fillDefaults(AppError{Code: CodeResponseDecode, Meta: meta}, err.Error())
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fillDefaults(AppError{Code: CodeNetworkTimeout, Retryable: true, Meta: meta}, err.Error())
	}

	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return fillDefaults(AppError{Code: CodeNetworkTimeout, Retryable: true, Meta: meta}, err.Error())
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return fillDefaults(AppError{Code: CodeNetworkDNS, Meta: meta}, err.Error())
	}

	s := strings.ToLower(err.Error())
	if strings.Contains(s, "x509") || strings.Contains(s, "tls") || strings.Contains(s, "handshake") {
		return fillDefaults(AppError{Code: CodeTLSHandshake, Meta: meta}, err.Error())
	}
	if strings.Contains(s, "no such host") {
		return fillDefaults(AppError{Code: CodeNetworkDNS, Meta: meta}, err.Error())
	}
	if strings.Contains(s, "connection reset") || strings.Contains(s, "broken pipe") || s == "eof" || strings.Contains(s, ": eof") {
		return fillDefaults(AppError{Code: CodeNetworkReset, Retryable: true, Meta: meta}, err.Error())
	}

	return fillDefaults(AppError{Code: CodeInternal, Meta: meta}, err.Error())
}

func fillDefaults(a AppError, cause string) AppError {
	a.Cause = cause
	if a.Meta == nil {
		a.Meta = map[string]string{}
	}
	if a.Message == "" {
		a.Message = defaultMessage(a.Code)
	}
	if a.Hint == "" {
		a.Hint = defaultHint(a.Code)
	}
	return a
}

func defaultMessage(code Code) string {
	switch code {
	case CodeNetworkDNS:
		return "dns resolution failed"
	case CodeNetworkTimeout:
		return "network request timed out"
	case CodeNetworkReset:
		return "connection closed by peer"
	case CodeTLSHandshake:
		return "tls handshake failed"
	case CodeHTTPNon2xx:
		return "upstream returned non-2xx"
	case CodeResponseDecode:
		return "failed to parse upstream response"
	case CodeConfigInvalid:
		return "invalid configuration"
	case CodeTemplateBad:
		return "invalid template"
	case CodeAuthDecodeInvalid:
		return "invalid base64 authentication value"
	case CodeAuthConfigInvalid:
		return "invalid authentication configuration"
	case CodeRobotsDisallowed:
		return "robots policy disallows this URL"
	case CodeContentTypeUnsupported:
		return "content type is not html"
	case CodeExtractEmpty:
		return "extracted content is empty"
	case CodeFetchFailed:
		return "failed to fetch page content"
	default:
		return "unexpected internal error"
	}
}

func defaultHint(code Code) string {
	switch code {
	case CodeNetworkDNS:
		return "check host name and DNS settings"
	case CodeNetworkTimeout:
		return "check network path or increase timeout"
	case CodeNetworkReset:
		return "check proxy, VPN, or network interception"
	case CodeTLSHandshake:
		return "check certificate chain or TLS interception"
	case CodeHTTPNon2xx:
		return "check upstream health and query"
	case CodeResponseDecode:
		return "check upstream response format"
	case CodeConfigInvalid:
		return "check config file and context"
	case CodeTemplateBad:
		return "check template syntax"
	case CodeAuthDecodeInvalid:
		return "provide valid base64-encoded authentication values"
	case CodeAuthConfigInvalid:
		return "check auth mode and required auth fields"
	case CodeRobotsDisallowed:
		return "disable robots check only when policy allows"
	case CodeContentTypeUnsupported:
		return "target URL must return text/html"
	case CodeExtractEmpty:
		return "try another page or fallback extractor"
	case CodeFetchFailed:
		return "check URL reachability and retry settings"
	default:
		return "inspect cause and retry"
	}
}
