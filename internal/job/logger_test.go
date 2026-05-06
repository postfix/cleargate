package job

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/postfix/cleargate/internal/runtime"
	"github.com/postfix/cleargate/internal/workspace"
)

func TestCaptureLogs(t *testing.T) {
	tempBase := t.TempDir()
	wm := workspace.NewManager(workspace.Config{BasePath: tempBase})
	jobID := "job_logger_123"
	wm.InitializeWorkspace(jobID)

	logger := NewLogger(wm)
	logChan := make(chan runtime.LogEvent)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Run capture in background
	errChan := make(chan error)
	go func() {
		errChan <- logger.CaptureLogs(ctx, jobID, logChan)
	}()

	// Send some logs
	logChan <- runtime.LogEvent{Stream: "stdout", Data: []byte("hello")}
	logChan <- runtime.LogEvent{Stream: "stderr", Data: []byte("error")}
	close(logChan)

	// Wait for completion
	if err := <-errChan; err != nil {
		t.Fatalf("CaptureLogs failed: %v", err)
	}

	// Verify files
	stdoutContent, _ := os.ReadFile(filepath.Join(tempBase, jobID, "logs", "stdout.log"))
	if string(stdoutContent) != "hello\n" {
		t.Errorf("stdout content mismatch: got %q", stdoutContent)
	}

	stderrContent, _ := os.ReadFile(filepath.Join(tempBase, jobID, "logs", "stderr.log"))
	if string(stderrContent) != "error\n" {
		t.Errorf("stderr content mismatch: got %q", stderrContent)
	}
}
