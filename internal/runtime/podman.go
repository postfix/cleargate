package runtime

import (
	"context"
	"fmt"
	"os"
	gruntime "runtime"
	"time"

	"github.com/opencontainers/runtime-spec/specs-go"
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
	exists, err := images.Exists(p.connCtx, image, nil)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = images.Pull(p.connCtx, image, nil)
	return err
}

func (p *PodmanRuntime) Create(ctx context.Context, req CreateContainerRequest) (ContainerID, error) {
	s := specgen.NewSpecGenerator(req.Image, false)
	s.Name = req.Name
	s.Command = req.Command

	// Hardened security profile by default
	s.CapDrop = []string{"ALL"}
	if len(req.CapAdd) > 0 {
		s.CapAdd = req.CapAdd
	}

	// Need pointers for some boolean options
	trueVal := true
	s.ReadOnlyFilesystem = &trueVal
	s.NoNewPrivileges = &trueVal

	if req.WorkspaceDir != "" {
		s.WorkDir = "/workspace"
		s.Mounts = append(s.Mounts, specs.Mount{
			Destination: "/workspace/input",
			Type:        "bind",
			Source:      req.WorkspaceDir + "/input",
			Options:     []string{"ro", "bind"},
		}, specs.Mount{
			Destination: "/workspace/output",
			Type:        "bind",
			Source:      req.WorkspaceDir + "/output",
			Options:     []string{"rw", "bind"},
		})
	}

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

func (p *PodmanRuntime) Wait(ctx context.Context, id ContainerID) (int, error) {
	// Wait logic. It returns an exit code.
	exitCode, err := containers.Wait(p.connCtx, string(id), nil)
	return int(exitCode), err
}

func (p *PodmanRuntime) Logs(ctx context.Context, id ContainerID) (<-chan LogEvent, error) {
	ch := make(chan LogEvent)

	stdoutChan := make(chan string)
	stderrChan := make(chan string)

	opts := new(containers.LogOptions).WithFollow(true).WithStdout(true).WithStderr(true)

	// Fetch logs in background and close the wrapper channel when done
	go func() {
		defer close(ch)

		// Run containers.Logs in its own goroutine to feed stdoutChan/stderrChan
		go func() {
			err := containers.Logs(p.connCtx, string(id), opts, stdoutChan, stderrChan)
			if err != nil {
				fmt.Printf("Error getting logs for %s: %v\n", id, err)
			}
			close(stdoutChan)
			close(stderrChan)
		}()

		for {
			select {
			case out, ok := <-stdoutChan:
				if ok {
					ch <- LogEvent{Stream: "stdout", Data: []byte(out)}
				} else {
					stdoutChan = nil
				}
			case errOut, ok := <-stderrChan:
				if ok {
					ch <- LogEvent{Stream: "stderr", Data: []byte(errOut)}
				} else {
					stderrChan = nil
				}
			case <-ctx.Done():
				return
			}

			if stdoutChan == nil && stderrChan == nil {
				break
			}
		}
	}()

	return ch, nil
}

func (p *PodmanRuntime) Stop(ctx context.Context, id ContainerID) error {
	return containers.Kill(p.connCtx, string(id), nil)
}

func (p *PodmanRuntime) GracefulStop(ctx context.Context, id ContainerID) error {
	sigTerm := "SIGTERM"
	err := containers.Kill(p.connCtx, string(id), new(containers.KillOptions).WithSignal(sigTerm))
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	sigKill := "SIGKILL"
	_ = containers.Kill(p.connCtx, string(id), new(containers.KillOptions).WithSignal(sigKill))
	return nil
}
