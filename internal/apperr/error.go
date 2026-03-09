package apperr

import "fmt"

type Code string

const (
	CodeNetworkDNS             Code = "NETWORK_DNS"
	CodeNetworkTimeout         Code = "NETWORK_TIMEOUT"
	CodeNetworkReset           Code = "NETWORK_RESET"
	CodeTLSHandshake           Code = "TLS_HANDSHAKE"
	CodeHTTPNon2xx             Code = "HTTP_NON_2XX"
	CodeResponseDecode         Code = "RESPONSE_DECODE"
	CodeConfigInvalid          Code = "CONFIG_INVALID"
	CodeTemplateBad            Code = "TEMPLATE_INVALID"
	CodeAuthDecodeInvalid      Code = "AUTH_DECODE_INVALID"
	CodeAuthConfigInvalid      Code = "AUTH_CONFIG_INVALID"
	CodeRobotsDisallowed       Code = "ROBOTS_DISALLOWED"
	CodeContentTypeUnsupported Code = "CONTENT_TYPE_UNSUPPORTED"
	CodeExtractEmpty           Code = "EXTRACT_EMPTY"
	CodeFetchFailed            Code = "FETCH_FAILED"
	CodeInternal               Code = "INTERNAL_UNKNOWN"
)

type AppError struct {
	Code      Code
	Retryable bool
	Message   string
	Hint      string
	Cause     string
	Meta      map[string]string
}

func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Cause)
}

type withMeta struct {
	err  error
	meta map[string]string
}

func (e *withMeta) Error() string { return e.err.Error() }
func (e *withMeta) Unwrap() error { return e.err }

func Annotate(err error, meta map[string]string) error {
	if err == nil || len(meta) == 0 {
		return err
	}
	return &withMeta{err: err, meta: meta}
}

type HTTPStatusError struct {
	StatusCode int
	Status     string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("non-2xx response: %s", e.Status)
}

type DecodeError struct{ Err error }

func (e *DecodeError) Error() string { return fmt.Sprintf("decode response: %v", e.Err) }
func (e *DecodeError) Unwrap() error { return e.Err }

type ConfigError struct{ Err error }

func (e *ConfigError) Error() string { return fmt.Sprintf("config error: %v", e.Err) }
func (e *ConfigError) Unwrap() error { return e.Err }

type TemplateError struct{ Err error }

func (e *TemplateError) Error() string { return fmt.Sprintf("template error: %v", e.Err) }
func (e *TemplateError) Unwrap() error { return e.Err }

type AuthDecodeError struct{ Err error }

func (e *AuthDecodeError) Error() string { return fmt.Sprintf("auth decode invalid: %v", e.Err) }
func (e *AuthDecodeError) Unwrap() error { return e.Err }

type AuthConfigError struct{ Err error }

func (e *AuthConfigError) Error() string { return fmt.Sprintf("auth config invalid: %v", e.Err) }
func (e *AuthConfigError) Unwrap() error { return e.Err }
