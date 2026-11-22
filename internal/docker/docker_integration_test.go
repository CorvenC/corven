package docker

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitlab.com/corvenc/corven/internal/source"
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

func TestFetchAndBuildWithBuildKit(t *testing.T) {
	c := newTestClient(t)
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	// A local repository holding the test app, fetched like a remote one.
	repo := t.TempDir()
	for _, name := range []string{"Dockerfile", "index.html"} {
		data, err := os.ReadFile(filepath.Join("testdata", "hello", name))
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(repo, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	gitRun(t, repo, "init", "--quiet", "--initial-branch=main")
	gitRun(t, repo, "add", ".")
	gitRun(t, repo, "commit", "--quiet", "-m", "app")

	work := filepath.Join(t.TempDir(), "src")
	sha, err := source.Fetch(ctx, source.FetchOptions{URL: "file://" + filepath.ToSlash(repo), Ref: "main", Dir: work, AllowFileURLs: true})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}

	var log strings.Builder
	id, err := c.Build(ctx, BuildOptions{ContextDir: work, Tag: "corven-test/fetched:" + sha[:12], Output: &log})
	if err != nil {
		t.Fatalf("Build: %v\n%s", err, log.String())
	}
	t.Cleanup(func() { _ = c.RemoveImage(context.Background(), id) })

	// BuildKit progress lines look like "#1 [internal] load build definition from Dockerfile".
	if !strings.Contains(log.String(), "#1 ") || !strings.Contains(log.String(), "COPY index.html") {
		t.Errorf("build log does not look like BuildKit progress:\n%s", log.String())
	}
}

func TestBuildKitFailingStep(t *testing.T) {
	c := newTestClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM busybox:1.37\nRUN echo boom && false\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var log strings.Builder
	_, err := c.Build(ctx, BuildOptions{ContextDir: dir, Tag: "corven-test/failing:latest", Output: &log})
	if err == nil {
		t.Fatalf("expected the build to fail\n%s", log.String())
	}
	if !strings.Contains(log.String(), "boom") {
		t.Errorf("log should include the failing step output:\n%s", log.String())
	}
}

func gitRun(t *testing.T, dir string, args ...string) {
	t.Helper()
	base := []string{"-c", "user.name=Test", "-c", "user.email=test@example.com", "-c", "commit.gpgsign=false"}
	cmd := exec.Command("git", append(base, args...)...)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
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
