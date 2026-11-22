package docker

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func TestTarDir(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"Dockerfile":        "FROM scratch\n",
		"src/main.txt":      "hello",
		".git/config":       "must not be sent",
		"src/.git/HEAD":     "nested .git is skipped too",
		"src/deep/file.txt": "deep",
	}
	for name, content := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
	}

	r, err := tarDir(dir)
	if err != nil {
		t.Fatalf("tarDir: %v", err)
	}

	got := map[string]string{}
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("read tar: %v", err)
		}
		body, _ := io.ReadAll(tr)
		got[hdr.Name] = string(body)
	}

	want := []string{"Dockerfile", "src/", "src/deep/", "src/deep/file.txt", "src/main.txt"}
	var names []string
	for name := range got {
		names = append(names, name)
	}
	sort.Strings(names)
	if strings.Join(names, ",") != strings.Join(want, ",") {
		t.Fatalf("entries = %v, want %v", names, want)
	}
	if got["src/main.txt"] != "hello" {
		t.Errorf("src/main.txt = %q, want %q", got["src/main.txt"], "hello")
	}
}

func TestTarDirMissingFolder(t *testing.T) {
	if _, err := tarDir(filepath.Join(t.TempDir(), "nope")); err == nil {
		t.Fatal("expected an error for a missing folder")
	}
}
