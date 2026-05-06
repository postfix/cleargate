//go:build remote
// +build remote

package podman

import (
	"context"
	"fmt"

	"go.podman.io/podman/v6/libpod/define"
	"go.podman.io/podman/v6/pkg/bindings"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/specgen"

	"github.com/postfix/cleargate/internal/runtime"
)

type Client struct {
	connCtx context.Context
}

// NewClient creates a new Podman client connected to the given socket URI.
func NewClient(ctx context.Context, socketURI string) (*Client, error) {
	conn, err := bindings.NewConnection(ctx, socketURI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to podman: %w", err)
	}
	return &Client{connCtx: conn}, nil
}

func (c *Client) PullImage(ctx context.Context, image string) error {
	_, err := images.Pull(c.connCtx, image, &images.PullOptions{})
	return err
}

func (c *Client) Create(ctx context.Context, req runtime.CreateContainerRequest) (runtime.ContainerID, error) {
	spec := specgen.NewSpecGenerator(req.Image, false)
	spec.Name = req.Name
	spec.Command = req.Command
	if req.Remove {
		rm := true
		spec.Remove = &rm
	}

	createResp, err := containers.CreateWithSpec(c.connCtx, spec, &containers.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return runtime.ContainerID(createResp.ID), nil
}

func (c *Client) Start(ctx context.Context, id runtime.ContainerID) error {
	return containers.Start(c.connCtx, string(id), &containers.StartOptions{})
}

func (c *Client) Wait(ctx context.Context, id runtime.ContainerID) error {
	_, err := containers.Wait(c.connCtx, string(id), &containers.WaitOptions{
		Condition: []define.ContainerStatus{
			define.ContainerStateExited,
			define.ContainerStateStopped,
		},
	})
	return err
}

func (c *Client) Logs(ctx context.Context, id runtime.ContainerID) (<-chan runtime.LogEvent, error) {
	// Minimal stub for Logs, as full implementation requires wiring stdout/stderr channels.
	// For MVP we just return an empty channel to satisfy the interface.
	ch := make(chan runtime.LogEvent)
	close(ch)
	return ch, nil
}

// Compile-time check to ensure Client implements ContainerRuntime
var _ runtime.ContainerRuntime = (*Client)(nil)
