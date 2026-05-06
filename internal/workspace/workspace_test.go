package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInitializeWorkspace(t *testing.T) {
	tempBase := t.TempDir()
	manager := NewManager(Config{BasePath: tempBase})

	jobID := "job_test_123"
	root, err := manager.InitializeWorkspace(jobID)
	if err != nil {
		t.Fatalf("InitializeWorkspace failed: %v", err)
	}

	expectedDirs := []string{"input", "output", "logs", "metadata"}
	for _, dir := range expectedDirs {
		path := filepath.Join(root, dir)
		info, err := os.Stat(path)
		if os.IsNotExist(err) {
			t.Errorf("Expected directory %s does not exist", path)
		} else if !info.IsDir() {
			t.Errorf("Path %s is not a directory", path)
		}
	}
}
