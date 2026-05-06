package job

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/postfix/cleargate/internal/runtime"
	"github.com/postfix/cleargate/internal/workspace"
)

// Logger handles piping container log events to files in the job workspace.
type Logger struct {
	workspaceManager *workspace.Manager
}

func NewLogger(wm *workspace.Manager) *Logger {
	return &Logger{workspaceManager: wm}
}

// CaptureLogs reads from the LogEvent channel and writes to stdout.log and stderr.log.
func (l *Logger) CaptureLogs(ctx context.Context, jobID string, logChan <-chan runtime.LogEvent) error {
	logsDir := l.workspaceManager.GetPath(jobID, "logs")

	stdoutFile, err := os.OpenFile(filepath.Join(logsDir, "stdout.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open stdout.log: %w", err)
	}
	defer stdoutFile.Close()

	stderrFile, err := os.OpenFile(filepath.Join(logsDir, "stderr.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open stderr.log: %w", err)
	}
	defer stderrFile.Close()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-logChan:
			if !ok {
				return nil // Channel closed, logging finished
			}
			if event.Stream == "stdout" {
				stdoutFile.Write(event.Data)
				stdoutFile.Write([]byte("\n")) // Ensure newline
			} else if event.Stream == "stderr" {
				stderrFile.Write(event.Data)
				stderrFile.Write([]byte("\n"))
			}
		}
	}
}
