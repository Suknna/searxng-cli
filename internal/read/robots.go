package read

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/temoto/robotstxt"
)

func CheckRobots(ctx context.Context, targetURL, userAgent string, timeout time.Duration) error {
	u, err := url.Parse(targetURL)
	if err != nil {
		return err
	}
	robotsURL := *u
	robotsURL.Path = path.Join("/", "robots.txt")
	robotsURL.RawQuery = ""

	client := &http.Client{Timeout: timeout}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, robotsURL.String(), nil)
	if err != nil {
		return err
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	rb, err := robotstxt.FromResponse(resp)
	if err != nil {
		return fmt.Errorf("parse robots: %w", err)
	}
	agent := userAgent
	if agent == "" {
		agent = "*"
	}
	group := rb.FindGroup(agent)
	if !group.Test(u.Path) {
		return &RobotsDisallowedError{URL: targetURL}
	}
	return nil
}
