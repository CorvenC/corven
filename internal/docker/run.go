package docker

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

// Default resource limits for a service (spec: app-deployment, "Resource limits").
const (
	DefaultMemoryBytes = 512 << 20     // 512 MB
	DefaultNanoCPUs    = 1_000_000_000 // 1 CPU
)

// RunOptions describes a container to start.
type RunOptions struct {
	Image       string            // image ID or reference
	Name        string            // container name; empty lets Docker choose
	Env         []string          // KEY=value pairs
	Port        uint16            // container port to publish on 127.0.0.1; 0 publishes nothing
	Labels      map[string]string // extra labels; the managed label is always added
	MemoryBytes int64             // 0 means DefaultMemoryBytes
	NanoCPUs    int64             // 0 means DefaultNanoCPUs
}

// Container is a started container.
type Container struct {
	ID       string
	HostPort string // port on 127.0.0.1 mapped to RunOptions.Port, empty if none was published
}

// Run creates and starts a container. A published port is bound to 127.0.0.1 on a free host port,
// so nothing is exposed to the network until the proxy routes to it.
func (c *Client) Run(ctx context.Context, opts RunOptions) (Container, error) {
	if opts.Image == "" {
		return Container{}, errors.New("run: image is required")
	}

	cfg := &container.Config{
		Image:  opts.Image,
		Env:    opts.Env,
		Labels: labels(opts.Labels),
	}
	host := &container.HostConfig{
		SecurityOpt: []string{"no-new-privileges:true"},
		Resources: container.Resources{
			Memory:   orDefault(opts.MemoryBytes, DefaultMemoryBytes),
			NanoCPUs: orDefault(opts.NanoCPUs, DefaultNanoCPUs),
		},
	}

	var port network.Port
	if opts.Port != 0 {
		var ok bool
		port, ok = network.PortFrom(opts.Port, network.TCP)
		if !ok {
			return Container{}, fmt.Errorf("run: invalid port %d", opts.Port)
		}
		cfg.ExposedPorts = network.PortSet{port: {}}
		host.PortBindings = network.PortMap{
			port: {{HostIP: netip.MustParseAddr("127.0.0.1"), HostPort: ""}},
		}
	}

	created, err := c.api.ContainerCreate(ctx, client.ContainerCreateOptions{
		Name:       opts.Name,
		Config:     cfg,
		HostConfig: host,
	})
	if err != nil {
		return Container{}, fmt.Errorf("create container from %s: %w", opts.Image, err)
	}

	if _, err := c.api.ContainerStart(ctx, created.ID, client.ContainerStartOptions{}); err != nil {
		_ = c.Remove(context.WithoutCancel(ctx), created.ID)
		return Container{}, fmt.Errorf("start container %s: %w", created.ID, err)
	}

	out := Container{ID: created.ID}
	if opts.Port != 0 {
		inspect, err := c.api.ContainerInspect(ctx, created.ID, client.ContainerInspectOptions{})
		if err != nil {
			return out, fmt.Errorf("inspect container %s: %w", created.ID, err)
		}
		bindings := inspect.Container.NetworkSettings.Ports[port]
		if len(bindings) == 0 {
			return out, fmt.Errorf("container %s: port %d was not published", created.ID, opts.Port)
		}
		out.HostPort = bindings[0].HostPort
	}
	return out, nil
}

// Remove stops and deletes a container, including its anonymous volumes.
func (c *Client) Remove(ctx context.Context, id string) error {
	if _, err := c.api.ContainerRemove(ctx, id, client.ContainerRemoveOptions{Force: true, RemoveVolumes: true}); err != nil {
		return fmt.Errorf("remove container %s: %w", id, err)
	}
	return nil
}

func orDefault(v, def int64) int64 {
	if v == 0 {
		return def
	}
	return v
}
