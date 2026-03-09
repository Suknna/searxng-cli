package read

import (
	"context"
	"io"
	"net/http"
	"strings"
)

func FetchHTML(ctx context.Context, target string, opts Options) (FetchedPage, error) {
	client := &http.Client{Timeout: opts.Timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return FetchedPage{}, &FetchFailedError{Err: err}
	}
	if opts.UserAgent != "" {
		req.Header.Set("User-Agent", opts.UserAgent)
	}

	resp, err := client.Do(req)
	if err != nil {
		return FetchedPage{}, &FetchFailedError{Err: err}
	}
	defer resp.Body.Close()

	r := io.Reader(resp.Body)
	if opts.MaxBytes > 0 {
		r = io.LimitReader(resp.Body, opts.MaxBytes)
	}
	body, err := io.ReadAll(r)
	if err != nil {
		return FetchedPage{}, &FetchFailedError{Err: err}
	}

	ct := strings.ToLower(resp.Header.Get("Content-Type"))
	if !strings.Contains(ct, "text/html") {
		return FetchedPage{}, &ContentTypeUnsupportedError{ContentType: ct}
	}

	return FetchedPage{
		URL:         resp.Request.URL.String(),
		StatusCode:  resp.StatusCode,
		ContentType: ct,
		HTML:        string(body),
	}, nil
}
