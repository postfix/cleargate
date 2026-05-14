package runtime

import (
	"context"
)

// ContainerID is a string representing the unique identifier of a container.
type ContainerID string

// CreateContainerRequest contains parameters for creating a new container.
type CreateContainerRequest struct {
	Image   string
	Name    string
	Command []string
	Remove  bool
	
	WorkspaceDir   string
	CapAdd         []string
	TimeoutSeconds int
}

// LogEvent represents a single line or chunk of log output.
type LogEvent struct {
	Stream string // "stdout" or "stderr"
	Data   []byte
}

// ContainerRuntime defines the interface for interacting with container sandboxes.
type ContainerRuntime interface {
	PullImage(ctx context.Context, image string) error
	Create(ctx context.Context, req CreateContainerRequest) (ContainerID, error)
	Start(ctx context.Context, id ContainerID) error
	Wait(ctx context.Context, id ContainerID) (int, error)
	Stop(ctx context.Context, id ContainerID) error
	GracefulStop(ctx context.Context, id ContainerID) error
	Logs(ctx context.Context, id ContainerID) (<-chan LogEvent, error)
	// Add other methods (Inspect, Remove, CopyTo, CopyFrom) as needed
}

// Ensure the interface is satisfied by the podman client (checked in tests)
