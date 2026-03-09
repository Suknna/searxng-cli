package read

import "time"

type Options struct {
	Format        string
	Timeout       time.Duration
	RespectRobots bool
	MaxBytes      int64
	Retry         int
	UserAgent     string
}

type FetchedPage struct {
	URL         string
	StatusCode  int
	ContentType string
	HTML        string
}
