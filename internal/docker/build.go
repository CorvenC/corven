package docker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/moby/moby/api/types/build"
	"github.com/moby/moby/api/types/jsonstream"
	"github.com/moby/moby/client"
)

// Message IDs the daemon uses for BuildKit output.
const (
	auxTrace   = "moby.buildkit.trace" // aux is a base64 BuildKit StatusResponse
	auxImageID = "moby.image.id"       // aux is {"ID": "sha256:..."}
)

// BuildOptions describes an image build from a local folder.
type BuildOptions struct {
	ContextDir string            // folder sent to the daemon as the build context
	Dockerfile string            // path inside ContextDir; empty means "Dockerfile"
	Tag        string            // image reference to create, for example "corven/shop:abc1234"
	Labels     map[string]string // extra image labels; the managed label is always added
	Output     io.Writer         // receives the build log; nil discards it
}

// Build builds an image with BuildKit and returns its ID. The build log streams to opts.Output
// as it happens. Layers are cached by the daemon, so rebuilding unchanged steps is fast.
func (c *Client) Build(ctx context.Context, opts BuildOptions) (string, error) {
	if opts.ContextDir == "" || opts.Tag == "" {
		return "", errors.New("build: context dir and tag are required")
	}
	dockerfile := opts.Dockerfile
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	out := opts.Output
	if out == nil {
		out = io.Discard
	}

	buildContext, err := tarDir(opts.ContextDir)
	if err != nil {
		return "", err
	}

	res, err := c.api.ImageBuild(ctx, buildContext, client.ImageBuildOptions{
		Tags:       []string{opts.Tag},
		Dockerfile: dockerfile,
		Labels:     labels(opts.Labels),
		Remove:     true,
		Version:    build.BuilderBuildKit,
	})
	if err != nil {
		return "", fmt.Errorf("build %s: %w", opts.Tag, err)
	}
	defer res.Body.Close()

	id, err := readBuildStream(res.Body, out)
	if err != nil {
		return "", fmt.Errorf("build %s: %w", opts.Tag, err)
	}
	return id, nil
}

// readBuildStream copies the daemon's JSON build messages to out as plain text,
// returns the image ID announced at the end, and turns an error message into a Go error.
// It understands both BuildKit traces and classic builder lines.
func readBuildStream(body io.Reader, out io.Writer) (string, error) {
	dec := json.NewDecoder(body)
	prog := newProgress(out)
	var id string
	for {
		var msg jsonstream.Message
		if err := dec.Decode(&msg); errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return "", fmt.Errorf("read build output: %w", err)
		}

		if msg.Error != nil {
			return "", errors.New(msg.Error.Message)
		}
		if msg.Stream != "" {
			if _, err := io.WriteString(out, msg.Stream); err != nil {
				return "", fmt.Errorf("write build output: %w", err)
			}
		}
		if msg.Aux != nil && msg.ID == auxTrace {
			var payload []byte // JSON base64 string, decoded by encoding/json
			if err := json.Unmarshal(*msg.Aux, &payload); err != nil {
				return "", fmt.Errorf("read build trace: %w", err)
			}
			if err := prog.handle(payload); err != nil {
				return "", err
			}
			continue
		}
		if msg.Aux != nil && (msg.ID == auxImageID || msg.ID == "") {
			var aux struct {
				ID string `json:"ID"`
			}
			if json.Unmarshal(*msg.Aux, &aux) == nil && aux.ID != "" {
				id = aux.ID
			}
		}
	}
	if id == "" {
		return "", errors.New("daemon did not report an image ID")
	}
	return id, nil
}

// RemoveImage deletes an image by ID or reference.
func (c *Client) RemoveImage(ctx context.Context, ref string) error {
	if _, err := c.api.ImageRemove(ctx, ref, client.ImageRemoveOptions{Force: true, PruneChildren: true}); err != nil {
		return fmt.Errorf("remove image %s: %w", ref, err)
	}
	return nil
}
