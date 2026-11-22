package docker

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

// newTestClient connects to the local Docker daemon or skips the test when there is none
// (for example in CI jobs without Docker). Run with -short to always skip.
func newTestClient(t *testing.T) *Client {
	t.Helper()
	if testing.Short() {
		t.Skip("integration test: needs Docker")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err := New(ctx)
	if err != nil {
		t.Skipf("integration test: Docker not available: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func TestBuildAndRun(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	tag := "corven-test/hello:" + strings.ToLower(strings.ReplaceAll(t.Name(), "/", "-"))
	var log strings.Builder
	id, err := c.Build(ctx, BuildOptions{ContextDir: "testdata/hello", Tag: tag, Output: &log})
	if err != nil {
		t.Fatalf("Build: %v\n%s", err, log.String())
	}
	t.Cleanup(func() { _ = c.RemoveImage(context.Background(), id) })
	if !strings.HasPrefix(id, "sha256:") {
		t.Errorf("image id = %q, want sha256:...", id)
	}

	ctr, err := c.Run(ctx, RunOptions{Image: id, Port: 8080, Labels: map[string]string{"test": t.Name()}})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	t.Cleanup(func() { _ = c.Remove(context.Background(), ctr.ID) })
	if ctr.HostPort == "" {
		t.Fatal("no host port published")
	}

	body := waitForHTTP(t, "http://127.0.0.1:"+ctr.HostPort+"/", 20*time.Second)
	if !strings.Contains(body, "hello from corven") {
		t.Errorf("body = %q, want the test page", body)
	}
}

func TestBuildFailureReportsDaemonError(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err := c.Build(ctx, BuildOptions{ContextDir: "testdata/hello", Dockerfile: "missing.Dockerfile", Tag: "corven-test/broken:latest"})
	if err == nil {
		t.Fatal("expected an error for a missing Dockerfile")
	}
}

// waitForHTTP polls url until it answers 200 or the timeout passes, and returns the body.
func waitForHTTP(t *testing.T, url string, timeout time.Duration) string {
	t.Helper()
	deadline := time.Now().Add(timeout)
	httpClient := &http.Client{Timeout: 2 * time.Second}
	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := httpClient.Get(url)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			_ = resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return string(body)
			}
			lastErr = &statusError{resp.StatusCode}
		} else {
			lastErr = err
		}
		time.Sleep(300 * time.Millisecond)
	}
	t.Fatalf("GET %s did not answer 200 within %s: %v", url, timeout, lastErr)
	return ""
}

type statusError struct{ code int }

func (e *statusError) Error() string { return http.StatusText(e.code) }
