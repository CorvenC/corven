package docker

import (
	"strings"
	"testing"
)

func TestReadBuildStream(t *testing.T) {
	stream := `{"stream":"Step 1/2 : FROM busybox\n"}
{"stream":" ---> abc\n"}
{"aux":{"ID":"sha256:1234"}}
{"stream":"Successfully built 1234\n"}
`
	var out strings.Builder
	id, err := readBuildStream(strings.NewReader(stream), &out)
	if err != nil {
		t.Fatalf("readBuildStream: %v", err)
	}
	if id != "sha256:1234" {
		t.Errorf("id = %q, want %q", id, "sha256:1234")
	}
	if !strings.Contains(out.String(), "Step 1/2 : FROM busybox") || !strings.Contains(out.String(), "Successfully built") {
		t.Errorf("log = %q, missing build lines", out.String())
	}
}

func TestReadBuildStreamError(t *testing.T) {
	stream := `{"stream":"Step 1/2 : RUN false\n"}
{"errorDetail":{"code":1,"message":"The command '/bin/sh -c false' returned a non-zero code: 1"},"error":"x"}
`
	_, err := readBuildStream(strings.NewReader(stream), &strings.Builder{})
	if err == nil || !strings.Contains(err.Error(), "non-zero code") {
		t.Fatalf("err = %v, want the daemon's error message", err)
	}
}

func TestReadBuildStreamNoID(t *testing.T) {
	_, err := readBuildStream(strings.NewReader(`{"stream":"done\n"}`), &strings.Builder{})
	if err == nil {
		t.Fatal("expected an error when no image ID is reported")
	}
}

func TestLabelsAddsManaged(t *testing.T) {
	in := map[string]string{"app": "shop"}
	got := labels(in)
	if got[ManagedLabel] != "true" || got["app"] != "shop" {
		t.Errorf("labels = %v", got)
	}
	if _, changed := in[ManagedLabel]; changed {
		t.Error("labels modified its input")
	}
}
