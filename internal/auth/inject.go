package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"searxng-cli/internal/apperr"
)

type Options struct {
	Mode         string
	APIKeyHeader string
	APIKey       string
	Username     string
	Password     string
}

func Apply(req *http.Request, opts Options) error {
	mode := strings.ToLower(strings.TrimSpace(opts.Mode))
	switch mode {
	case "", "none":
		return nil
	case "api_key":
		if strings.TrimSpace(opts.APIKeyHeader) == "" {
			return &apperr.AuthConfigError{Err: fmt.Errorf("api key header is required")}
		}
		decoded, err := decodeBase64(opts.APIKey)
		if err != nil {
			return err
		}
		req.Header.Set(opts.APIKeyHeader, decoded)
		return nil
	case "basic":
		user, err := decodeBase64(opts.Username)
		if err != nil {
			return err
		}
		pass, err := decodeBase64(opts.Password)
		if err != nil {
			return err
		}
		req.SetBasicAuth(user, pass)
		return nil
	default:
		return &apperr.AuthConfigError{Err: fmt.Errorf("unsupported mode %q", opts.Mode)}
	}
}

func decodeBase64(v string) (string, error) {
	if strings.TrimSpace(v) == "" {
		return "", &apperr.AuthConfigError{Err: fmt.Errorf("missing base64 value")}
	}
	b, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return "", &apperr.AuthDecodeError{Err: err}
	}
	return string(b), nil
}
