package job

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/postfix/cleargate/internal/runtime"
	"github.com/postfix/cleargate/internal/workspace"
)

type LogMessage struct {
	Type string
	Data []byte
}

// Logger handles piping container log events to files and streaming.
type Logger struct {
	workspaceManager *workspace.Manager
	subs             map[string][]chan LogMessage
	mu               sync.Mutex
}

func NewLogger(wm *workspace.Manager) *Logger {
	return &Logger{
		workspaceManager: wm,
		subs:             make(map[string][]chan LogMessage),
	}
}

func (l *Logger) Subscribe(jobID string) <-chan LogMessage {
	l.mu.Lock()
	defer l.mu.Unlock()
	ch := make(chan LogMessage, 100)
	l.subs[jobID] = append(l.subs[jobID], ch)
	return ch
}

func (l *Logger) Unsubscribe(jobID string, ch <-chan LogMessage) {
	l.mu.Lock()
	defer l.mu.Unlock()
	subs := l.subs[jobID]
	for i, sub := range subs {
		if sub == ch {
			l.subs[jobID] = append(subs[:i], subs[i+1:]...)
			break
		}
	}
}

func (l *Logger) broadcast(jobID string, msg LogMessage) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, sub := range l.subs[jobID] {
		select {
		case sub <- msg:
		default: // non-blocking
		}
	}
}

func (l *Logger) LogStdout(jobID string, data []byte) {
	l.broadcast(jobID, LogMessage{Type: "stdout", Data: data})
}

func (l *Logger) LogStderr(jobID string, data []byte) {
	l.broadcast(jobID, LogMessage{Type: "stderr", Data: data})
}

func (l *Logger) LogStatus(jobID string, status string, exitCode int) {
	l.broadcast(jobID, LogMessage{Type: "status", Data: []byte(status)})
}

// CaptureLogs reads from the LogEvent channel and writes to files (legacy/fallback).
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
				return nil
			}
			if event.Stream == "stdout" {
				stdoutFile.Write(event.Data)
				stdoutFile.Write([]byte("\n"))
				l.LogStdout(jobID, event.Data)
			} else if event.Stream == "stderr" {
				stderrFile.Write(event.Data)
				stderrFile.Write([]byte("\n"))
				l.LogStderr(jobID, event.Data)
			}
		}
	}
}
