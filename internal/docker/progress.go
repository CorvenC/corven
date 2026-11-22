package docker

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	controlapi "github.com/moby/buildkit/api/services/control"
)

// progress turns BuildKit status updates into a plain-text log similar to
// `docker build --progress=plain`:
//
//	#2 [1/2] FROM docker.io/library/busybox:1.37
//	#2 DONE 0.8s
//	#3 [2/2] COPY index.html /www/index.html
//	#3 CACHED
type progress struct {
	out   io.Writer
	steps map[string]*step // by vertex digest
	next  int
}

type step struct {
	num      int
	started  bool
	finished bool
}

func newProgress(out io.Writer) *progress {
	return &progress{out: out, steps: map[string]*step{}, next: 1}
}

// handle decodes one "moby.buildkit.trace" payload and writes the new lines.
func (p *progress) handle(payload []byte) error {
	var st controlapi.StatusResponse
	if err := st.UnmarshalVT(payload); err != nil {
		return fmt.Errorf("decode buildkit status: %w", err)
	}

	var b strings.Builder
	for _, v := range st.Vertexes {
		s := p.step(v.Digest)
		if v.Started != nil && !s.started {
			s.started = true
			fmt.Fprintf(&b, "#%d %s\n", s.num, v.Name)
		}
		if s.finished {
			continue
		}
		switch {
		case v.Cached:
			s.finished = true
			fmt.Fprintf(&b, "#%d CACHED\n", s.num)
		case v.Error != "":
			s.finished = true
			fmt.Fprintf(&b, "#%d ERROR: %s\n", s.num, v.Error)
		case v.Completed != nil && v.Started != nil:
			s.finished = true
			took := v.Completed.AsTime().Sub(v.Started.AsTime()).Seconds()
			fmt.Fprintf(&b, "#%d DONE %.1fs\n", s.num, took)
		}
	}
	for _, l := range st.Logs {
		s := p.step(l.Vertex)
		for _, line := range bytes.Split(bytes.TrimRight(l.Msg, "\n"), []byte("\n")) {
			fmt.Fprintf(&b, "#%d %s\n", s.num, line)
		}
	}

	if b.Len() == 0 {
		return nil
	}
	_, err := io.WriteString(p.out, b.String())
	return err
}

func (p *progress) step(digest string) *step {
	s, ok := p.steps[digest]
	if !ok {
		s = &step{num: p.next}
		p.next++
		p.steps[digest] = s
	}
	return s
}
