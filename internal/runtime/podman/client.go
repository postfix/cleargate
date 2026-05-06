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
	ch := make(chan runtime.LogEvent)
	
	// Implementation note for MVP: 
	// The podman bindings containers.Logs method requires channels for stdout and stderr
	// and runs asynchronously.
	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	go func() {
		defer close(ch)
		
		opts := &containers.LogOptions{
			Follow:   func() *bool { b := true; return &b }(),
			Stdout:   func() *bool { b := true; return &b }(),
			Stderr:   func() *bool { b := true; return &b }(),
		}
		
		// Run containers.Logs in a goroutine because it blocks while following
		go func() {
			_ = containers.Logs(c.connCtx, string(id), opts, stdoutChan, stderrChan)
			close(stdoutChan)
			close(stderrChan)
		}()

		// Multiplex the channels
		for {
			select {
			case out, ok := <-stdoutChan:
				if ok {
					ch <- runtime.LogEvent{Stream: "stdout", Data: []byte(out)}
				} else {
					stdoutChan = nil
				}
			case errOut, ok := <-stderrChan:
				if ok {
					ch <- runtime.LogEvent{Stream: "stderr", Data: []byte(errOut)}
				} else {
					stderrChan = nil
				}
			}
			
			if stdoutChan == nil && stderrChan == nil {
				break
			}
		}
	}()
	
	return ch, nil
}

// Compile-time check to ensure Client implements ContainerRuntime
var _ runtime.ContainerRuntime = (*Client)(nil)
