package read

import (
	"context"
	"time"
)

func ReadURL(ctx context.Context, target string, opts Options) (string, error) {
	if opts.Timeout <= 0 {
		opts.Timeout = 10 * time.Second
	}
	if opts.Format == "" {
		opts.Format = "markdown"
	}
	if opts.UserAgent == "" {
		opts.UserAgent = "searxng-cli/1.0"
	}

	if opts.RespectRobots {
		if err := CheckRobots(ctx, target, opts.UserAgent, opts.Timeout); err != nil {
			return "", err
		}
	}

	attempts := opts.Retry + 1
	var lastErr error
	for i := 0; i < attempts; i++ {
		page, err := FetchHTML(ctx, target, opts)
		if err != nil {
			lastErr = err
			continue
		}
		out, err := ExtractContent(page.HTML, page.URL, opts.Format)
		if err != nil {
			return "", err
		}
		return out, nil
	}
	return "", lastErr
}
