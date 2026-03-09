package read

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckRobotsDisallow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/robots.txt" {
			_, _ = w.Write([]byte("User-agent: *\nDisallow: /blocked\n"))
			return
		}
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	err := CheckRobots(context.Background(), ts.URL+"/blocked", "*", time.Second)
	if err == nil {
		t.Fatal("expected disallow error")
	}
}
