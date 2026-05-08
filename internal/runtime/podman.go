package runtime

import (
	"context"
	"fmt"
	"os"
	gruntime "runtime"
	"go.podman.io/podman/v6/pkg/bindings"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/specgen"
)

type PodmanRuntime struct {
	connCtx context.Context
}

func NewPodmanRuntime() (*PodmanRuntime, error) {
	uri := os.Getenv("CONTAINER_HOST")
	if uri == "" {
		if gruntime.GOOS == "darwin" {
			uri = fmt.Sprintf("unix:///Users/%s/.local/share/containers/podman/machine/podman.sock", os.Getenv("USER"))
		} else {
			uri = fmt.Sprintf("unix://%s/podman/podman.sock", os.Getenv("XDG_RUNTIME_DIR"))
		}
	}

	ctx, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to podman at %s: %w", uri, err)
	}

	return &PodmanRuntime{connCtx: ctx}, nil
}

func (p *PodmanRuntime) PullImage(ctx context.Context, image string) error {
	_, err := images.Pull(p.connCtx, image, nil)
	return err
}

func (p *PodmanRuntime) Create(ctx context.Context, req CreateContainerRequest) (ContainerID, error) {
	s := specgen.NewSpecGenerator(req.Image, false)
	s.Name = req.Name
	s.Command = req.Command

	// Hardened security profile by default
	s.CapDrop = []string{"ALL"}
	
	// Need pointers for some boolean options
	trueVal := true
	s.ReadOnlyFilesystem = &trueVal
	s.NoNewPrivileges = &trueVal

	// Optionally you can mount the workspace here if passed in req.
	// For now, if the workspace mount needs to be configured, the req needs a Mounts field.
	// We'll add it later if the API handler provides it.

	rm := req.Remove
	s.Remove = &rm

	createResponse, err := containers.CreateWithSpec(p.connCtx, s, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	return ContainerID(createResponse.ID), nil
}

func (p *PodmanRuntime) Start(ctx context.Context, id ContainerID) error {
	return containers.Start(p.connCtx, string(id), nil)
}

func (p *PodmanRuntime) Wait(ctx context.Context, id ContainerID) error {
	// Wait logic. It returns an exit code.
	_, err := containers.Wait(p.connCtx, string(id), nil)
	return err
}

func (p *PodmanRuntime) Logs(ctx context.Context, id ContainerID) (<-chan LogEvent, error) {
	ch := make(chan LogEvent)

	// In real implementation, this requires streaming. For the MVP, we might need a custom pipe or using containers.Logs
	// with stdout and stderr channels. Let's do a simple placeholder that just closes for now, as logs via bindings
	// typically requires passing io.Writers.
	
	go func() {
		defer close(ch)
		// We'll wire up real logs if needed, but the current LogEvent structure isn't directly matching Podman's io.Writers.
		// For MVP, we'll just not block.
	}()

	return ch, nil
}
