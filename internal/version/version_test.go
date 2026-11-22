package version

import (
	"runtime"
	"strings"
	"testing"
)

func TestString(t *testing.T) {
	oldVersion, oldCommit, oldDate := Version, Commit, Date
	t.Cleanup(func() { Version, Commit, Date = oldVersion, oldCommit, oldDate })

	Version, Commit, Date = "v1.2.3", "abc1234", "2026-12-31"

	got := String("corven")
	for _, want := range []string{"corven", "v1.2.3", "abc1234", "2026-12-31", runtime.GOOS + "/" + runtime.GOARCH} {
		if !strings.Contains(got, want) {
			t.Errorf("String() = %q, missing %q", got, want)
		}
	}
}

func TestDefaults(t *testing.T) {
	if Version != "dev" {
		t.Skip("version injected at build time")
	}
	if got := String("corven"); !strings.HasPrefix(got, "corven dev ") {
		t.Errorf("String() = %q, want prefix %q", got, "corven dev ")
	}
}
