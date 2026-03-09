package apperr

import (
	"fmt"
	"sort"
	"strings"
)

func RenderKV(err AppError) string {
	parts := []string{
		"ERROR",
		fmt.Sprintf("code=%s", err.Code),
		fmt.Sprintf("retryable=%t", err.Retryable),
		fmt.Sprintf("message=%q", err.Message),
		fmt.Sprintf("hint=%q", err.Hint),
		fmt.Sprintf("cause=%q", err.Cause),
	}

	keys := make([]string, 0, len(err.Meta))
	for k := range err.Meta {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s=%q", k, err.Meta[k]))
	}

	return strings.Join(parts, " ")
}
