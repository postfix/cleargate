package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/postfix/cleargate/internal/job"
	"github.com/postfix/cleargate/internal/models"
	"github.com/postfix/cleargate/internal/repository"
	"github.com/postfix/cleargate/internal/runtime"
	"github.com/postfix/cleargate/internal/workspace"
)

type MockRuntime struct {
	runtime.ContainerRuntime
	gracefulStopCalled bool
}

func (m *MockRuntime) Create(ctx context.Context, req runtime.CreateContainerRequest) (runtime.ContainerID, error) {
	return "test-container-id", nil
}

func (m *MockRuntime) Start(ctx context.Context, id runtime.ContainerID) error {
	return nil
}

func (m *MockRuntime) Logs(ctx context.Context, id runtime.ContainerID) (<-chan runtime.LogEvent, error) {
	ch := make(chan runtime.LogEvent)
	// Don't close so the test hangs and allows timeout
	return ch, nil
}

func (m *MockRuntime) Wait(ctx context.Context, id runtime.ContainerID) (int, error) {
	<-ctx.Done()
	return 137, nil
}

func (m *MockRuntime) GracefulStop(ctx context.Context, id runtime.ContainerID) error {
	m.gracefulStopCalled = true
	return nil
}

func (m *MockRuntime) PullImage(ctx context.Context, image string) error {
	return nil
}

func TestGracefulTimeout(t *testing.T) {
	repo, _ := repository.NewToolSpecRepository("?access_mode=READ_WRITE")
	defer repo.Close()

	_ = repo.SaveDraft(&models.ToolSpec{
		Metadata: models.Metadata{Name: "testtool", Version: "1.0"},
		Runtime: models.Runtime{
			Executable:     "sleep",
			TimeoutSeconds: 1, // Set timeout to 1 second
		},
	})
	repo.Approve("testtool")

	wm := workspace.NewManager(workspace.Config{BasePath: t.TempDir()})
	reg := job.NewRegistry()
	logger := job.NewLogger(wm)
	mockRt := &MockRuntime{}
	handler := NewExecutionHandler(mockRt, wm, logger, repo, reg, nil)

	reqPayload := map[string]interface{}{
		"tool_id": "testtool",
		"job_id":  "job_123",
		"values":  map[string]interface{}{},
	}
	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.HandleExecute(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rr.Code)
	}

	// Wait for the timeout to trigger (1 second + small buffer)
	time.Sleep(1200 * time.Millisecond)

	if !mockRt.gracefulStopCalled {
		t.Errorf("Expected GracefulStop to be called after timeout")
	}
}
