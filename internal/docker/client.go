// Package docker is Corven's thin layer over the Docker Engine API.
//
// It builds images from a local folder, runs them as containers and removes them.
// Everything it creates carries the [ManagedLabel] so Corven can find and clean up its own objects.
package docker

import (
	"context"
	"fmt"

	"github.com/moby/moby/client"
)

// ManagedLabel marks every image and container created by Corven.
const ManagedLabel = "dev.corven.managed"

// Client talks to a Docker daemon.
type Client struct {
	api *client.Client
}

// New connects to the Docker daemon configured by the environment (DOCKER_HOST and friends,
// or the default local socket) and checks that it answers.
func New(ctx context.Context) (*Client, error) {
	api, err := client.New(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("create docker client: %w", err)
	}
	if _, err := api.Ping(ctx, client.PingOptions{}); err != nil {
		_ = api.Close()
		return nil, fmt.Errorf("reach docker daemon: %w", err)
	}
	return &Client{api: api}, nil
}

// Close releases the connection to the daemon.
func (c *Client) Close() error {
	return c.api.Close()
}

// labels returns the given labels plus the managed label, without modifying the input.
func labels(extra map[string]string) map[string]string {
	out := make(map[string]string, len(extra)+1)
	for k, v := range extra {
		out[k] = v
	}
	out[ManagedLabel] = "true"
	return out
}
