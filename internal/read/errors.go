package read

import "fmt"

type RobotsDisallowedError struct{ URL string }

func (e *RobotsDisallowedError) Error() string {
	return fmt.Sprintf("robots disallowed: %s", e.URL)
}

type ContentTypeUnsupportedError struct{ ContentType string }

func (e *ContentTypeUnsupportedError) Error() string {
	return fmt.Sprintf("content type unsupported: %s", e.ContentType)
}

type ExtractEmptyError struct{}

func (e *ExtractEmptyError) Error() string { return "extract result is empty" }

type FetchFailedError struct{ Err error }

func (e *FetchFailedError) Error() string { return fmt.Sprintf("fetch failed: %v", e.Err) }
func (e *FetchFailedError) Unwrap() error { return e.Err }
