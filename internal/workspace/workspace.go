package workspace

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config holds workspace configuration.
type Config struct {
	BasePath string
}

// DefaultConfig provides a default configuration for workspaces.
func DefaultConfig() Config {
	home, _ := os.UserHomeDir()
	return Config{
		BasePath: filepath.Join(home, ".cleargate", "jobs"),
	}
}

// Manager handles the creation and path resolution for job workspaces.
type Manager struct {
	config Config
}

// NewManager creates a new workspace manager.
func NewManager(cfg Config) *Manager {
	return &Manager{config: cfg}
}

// InitializeWorkspace creates the deterministic directory layout for a given job.
// It returns the absolute path to the root of the job's workspace.
func (m *Manager) InitializeWorkspace(jobID string) (string, error) {
	jobRoot := filepath.Join(m.config.BasePath, jobID)

	dirs := []string{
		filepath.Join(jobRoot, "input"),
		filepath.Join(jobRoot, "output"),
		filepath.Join(jobRoot, "logs"),
		filepath.Join(jobRoot, "metadata"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return jobRoot, nil
}

// GetPath returns the path to a specific subdirectory within a job's workspace.
func (m *Manager) GetPath(jobID, subDir string) string {
	return filepath.Join(m.config.BasePath, jobID, subDir)
}
