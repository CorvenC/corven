// Package source fetches the code that Corven deploys.
//
// For now it shallow-fetches one commit of a Git repository with the git command line,
// which the platform image ships. Private repositories and provider APIs come later.
package source

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// FetchOptions describes what to fetch and where.
type FetchOptions struct {
	URL string // repository URL; https only unless AllowFileURLs is set
	Ref string // branch, tag or full commit SHA
	Dir string // empty or missing folder that receives the working tree

	// AllowFileURLs accepts file:// URLs. Tests use it with local repositories;
	// the platform never sets it.
	AllowFileURLs bool
}

// refPattern accepts branch names, tags and SHAs, and rejects anything that git could
// read as an option or a revision expression.
var refPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]{0,254}$`)

// Fetch downloads exactly one commit (no history) into opts.Dir and returns its full SHA.
func Fetch(ctx context.Context, opts FetchOptions) (string, error) {
	if err := validate(opts); err != nil {
		return "", err
	}
	if err := os.MkdirAll(opts.Dir, 0o750); err != nil {
		return "", fmt.Errorf("fetch: create %s: %w", opts.Dir, err)
	}

	steps := [][]string{
		{"init", "--quiet"},
		{"fetch", "--quiet", "--depth=1", "--no-tags", "--", opts.URL, opts.Ref},
		{"checkout", "--quiet", "--detach", "FETCH_HEAD"},
	}
	for _, args := range steps {
		if _, err := git(ctx, opts.Dir, args...); err != nil {
			return "", fmt.Errorf("fetch %s at %s: %w", redact(opts.URL), opts.Ref, err)
		}
	}

	sha, err := git(ctx, opts.Dir, "rev-parse", "HEAD")
	if err != nil {
		return "", fmt.Errorf("fetch: read commit: %w", err)
	}
	return sha, nil
}

func validate(opts FetchOptions) error {
	if opts.Dir == "" {
		return errors.New("fetch: dir is required")
	}
	if !refPattern.MatchString(opts.Ref) || strings.Contains(opts.Ref, "..") {
		return fmt.Errorf("fetch: invalid ref %q", opts.Ref)
	}
	u, err := url.Parse(opts.URL)
	if err != nil {
		return fmt.Errorf("fetch: invalid URL: %w", err)
	}
	switch {
	case u.Scheme == "https" && u.Host != "":
	case u.Scheme == "file" && opts.AllowFileURLs:
	default:
		return fmt.Errorf("fetch: unsupported URL %q (only https)", redact(opts.URL))
	}
	if entries, err := os.ReadDir(opts.Dir); err == nil && len(entries) > 0 {
		return fmt.Errorf("fetch: %s is not empty", opts.Dir)
	}
	return nil
}

// git runs one git command in dir and returns its trimmed stdout.
// Prompts are disabled so a private repo fails fast instead of waiting for a password.
func git(ctx context.Context, dir string, args ...string) (string, error) {
	//nolint:gosec // G204: the binary is fixed; URL and ref are checked in validate() and passed after "--".
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0", "GIT_ASKPASS=true", "LC_ALL=C")
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return "", fmt.Errorf("git %s: %s", args[0], msg)
	}
	return strings.TrimSpace(stdout.String()), nil
}

// redact hides credentials embedded in a URL before it reaches logs or errors.
func redact(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.User == nil {
		return raw
	}
	return u.Redacted()
}
