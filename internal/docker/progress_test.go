package docker

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	controlapi "github.com/moby/buildkit/api/services/control"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func status(t *testing.T, st *controlapi.StatusResponse) []byte {
	t.Helper()
	b, err := st.MarshalVT()
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestProgress(t *testing.T) {
	start := time.Date(2026, 12, 1, 10, 0, 0, 0, time.UTC)
	var out strings.Builder
	p := newProgress(&out)

	// Step a starts; step b is cached.
	err := p.handle(status(t, &controlapi.StatusResponse{Vertexes: []*controlapi.Vertex{
		{Digest: "a", Name: "[1/2] FROM busybox", Started: timestamppb.New(start)},
		{Digest: "b", Name: "[2/2] COPY index.html /www/", Started: timestamppb.New(start), Cached: true},
	}}))
	if err != nil {
		t.Fatal(err)
	}
	// Step a logs a line and finishes 1.5 s later. Repeated vertexes must not print twice.
	err = p.handle(status(t, &controlapi.StatusResponse{
		Vertexes: []*controlapi.Vertex{
			{Digest: "a", Name: "[1/2] FROM busybox", Started: timestamppb.New(start), Completed: timestamppb.New(start.Add(1500 * time.Millisecond))},
			{Digest: "b", Name: "[2/2] COPY index.html /www/", Started: timestamppb.New(start), Cached: true},
		},
		Logs: []*controlapi.VertexLog{{Vertex: "a", Msg: []byte("pulling layer\n")}},
	}))
	if err != nil {
		t.Fatal(err)
	}

	want := "#1 [1/2] FROM busybox\n" +
		"#2 [2/2] COPY index.html /www/\n" +
		"#2 CACHED\n" +
		"#1 DONE 1.5s\n" +
		"#1 pulling layer\n"
	if out.String() != want {
		t.Errorf("output:\n%s\nwant:\n%s", out.String(), want)
	}
}

func TestProgressError(t *testing.T) {
	var out strings.Builder
	p := newProgress(&out)
	err := p.handle(status(t, &controlapi.StatusResponse{Vertexes: []*controlapi.Vertex{
		{Digest: "x", Name: "[2/2] RUN false", Started: timestamppb.Now(), Error: "exit code: 1"},
	}}))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "#1 ERROR: exit code: 1") {
		t.Errorf("output = %q, want the step error", out.String())
	}
}

func TestReadBuildStreamBuildKit(t *testing.T) {
	trace := base64.StdEncoding.EncodeToString(status(t, &controlapi.StatusResponse{Vertexes: []*controlapi.Vertex{
		{Digest: "a", Name: "[1/1] FROM scratch", Started: timestamppb.Now(), Cached: true},
	}}))
	stream := fmt.Sprintf(`{"id":"moby.buildkit.trace","aux":%q}
{"id":"moby.image.id","aux":{"ID":"sha256:feed"}}
`, trace)

	var out strings.Builder
	id, err := readBuildStream(strings.NewReader(stream), &out)
	if err != nil {
		t.Fatalf("readBuildStream: %v", err)
	}
	if id != "sha256:feed" {
		t.Errorf("id = %q, want sha256:feed", id)
	}
	if !strings.Contains(out.String(), "#1 [1/1] FROM scratch") || !strings.Contains(out.String(), "#1 CACHED") {
		t.Errorf("log = %q, missing the BuildKit steps", out.String())
	}
}
