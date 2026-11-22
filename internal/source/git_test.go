package source

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// newRepo creates a local repository with two commits on "main" and returns its file:// URL
// and the SHA of the first commit. Skips the test when git is not installed.
func newRepo(t *testing.T) (repoURL, firstSHA string) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		base := []string{"-c", "user.name=Test", "-c", "user.email=test@example.com", "-c", "commit.gpgsign=false"}
		cmd := exec.Command("git", append(base, args...)...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}

	run("init", "--quiet", "--initial-branch=main")
	run("config", "uploadpack.allowReachableSHA1InWant", "true")
	write(t, filepath.Join(dir, "version.txt"), "one")
	run("add", ".")
	run("commit", "--quiet", "-m", "first")
	firstSHA = run("rev-parse", "HEAD")
	write(t, filepath.Join(dir, "version.txt"), "two")
	run("commit", "--quiet", "-am", "second")

	return "file://" + filepath.ToSlash(dir), firstSHA
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestFetchBranch(t *testing.T) {
	repoURL, _ := newRepo(t)
	dir := filepath.Join(t.TempDir(), "work")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sha, err := Fetch(ctx, FetchOptions{URL: repoURL, Ref: "main", Dir: dir, AllowFileURLs: true})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if len(sha) != 40 {
		t.Errorf("sha = %q, want a full 40-char SHA", sha)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "version.txt"))
	if string(got) != "two" {
		t.Errorf("version.txt = %q, want the latest commit", got)
	}
}

func TestFetchCommit(t *testing.T) {
	repoURL, first := newRepo(t)
	dir := filepath.Join(t.TempDir(), "work")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sha, err := Fetch(ctx, FetchOptions{URL: repoURL, Ref: first, Dir: dir, AllowFileURLs: true})
	if err != nil {
		t.Fatalf("Fetch: %v", err)
	}
	if sha != first {
		t.Errorf("sha = %q, want %q", sha, first)
	}
	got, _ := os.ReadFile(filepath.Join(dir, "version.txt"))
	if string(got) != "one" {
		t.Errorf("version.txt = %q, want the first commit", got)
	}
}

func TestFetchUnknownRef(t *testing.T) {
	repoURL, _ := newRepo(t)
	_, err := Fetch(context.Background(), FetchOptions{
		URL: repoURL, Ref: "does-not-exist", Dir: filepath.Join(t.TempDir(), "w"), AllowFileURLs: true,
	})
	if err == nil || !strings.Contains(err.Error(), "git fetch") {
		t.Fatalf("err = %v, want a git fetch error", err)
	}
}

func TestValidate(t *testing.T) {
	notEmpty := t.TempDir()
	write(t, filepath.Join(notEmpty, "x"), "x")

	tests := []struct {
		name string
		opts FetchOptions
	}{
		{"option injection in ref", FetchOptions{URL: "https://gitlab.com/a/b.git", Ref: "--upload-pack=evil", Dir: t.TempDir()}},
		{"revision range", FetchOptions{URL: "https://gitlab.com/a/b.git", Ref: "main..dev", Dir: t.TempDir()}},
		{"empty ref", FetchOptions{URL: "https://gitlab.com/a/b.git", Ref: "", Dir: t.TempDir()}},
		{"ssh URL", FetchOptions{URL: "ssh://git@gitlab.com/a/b.git", Ref: "main", Dir: t.TempDir()}},
		{"file URL without opt-in", FetchOptions{URL: "file:///tmp/repo", Ref: "main", Dir: t.TempDir()}},
		{"ext transport", FetchOptions{URL: "ext::sh -c touch% /tmp/pwned", Ref: "main", Dir: t.TempDir()}},
		{"missing dir", FetchOptions{URL: "https://gitlab.com/a/b.git", Ref: "main"}},
		{"dir not empty", FetchOptions{URL: "https://gitlab.com/a/b.git", Ref: "main", Dir: notEmpty}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validate(tt.opts); err == nil {
				t.Error("expected a validation error")
			}
		})
	}

	ok := FetchOptions{URL: "https://gitlab.com/corvenc/corven.git", Ref: "feat/x-1.2", Dir: t.TempDir()}
	if err := validate(ok); err != nil {
		t.Errorf("valid options rejected: %v", err)
	}
}

func TestRedact(t *testing.T) {
	got := redact("https://user:secret@gitlab.com/a/b.git")
	if strings.Contains(got, "secret") {
		t.Errorf("redact leaked the password: %q", got)
	}
}
