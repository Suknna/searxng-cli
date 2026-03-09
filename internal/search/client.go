package search

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"searxng-cli/internal/apperr"
	"searxng-cli/internal/auth"
)

type Result struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

type response struct {
	Results []Result `json:"results"`
}

func Fetch(ctx context.Context, server, query string, timeout time.Duration, authOpts auth.Options) ([]Result, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("query is empty")
	}

	base, err := url.Parse(strings.TrimRight(server, "/"))
	if err != nil {
		return nil, err
	}
	base.Path = "/search"

	params := url.Values{}
	params.Set("q", query)
	params.Set("format", "json")
	base.RawQuery = params.Encode()

	hc := &http.Client{Timeout: timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	if err != nil {
		return nil, err
	}
	if err := auth.Apply(req, authOpts); err != nil {
		return nil, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, &apperr.HTTPStatusError{StatusCode: resp.StatusCode, Status: resp.Status}
	}

	var payload response
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, &apperr.DecodeError{Err: err}
	}
	return payload.Results, nil
}
