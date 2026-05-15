package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/postfix/cleargate/internal/models"
	"github.com/postfix/cleargate/internal/repository"
	"github.com/postfix/cleargate/internal/workspace"
)

func TestUnknownFlagRejection(t *testing.T) {
	repo, _ := repository.NewToolSpecRepository("?access_mode=READ_WRITE")
	defer repo.Close()

	// Seed a valid tool
	_ = repo.SaveDraft(&models.ToolSpec{
		Metadata: models.Metadata{Name: "testtool", Version: "1.0"},
		Flags: []models.Flag{
			{ID: "allowed-flag", Type: "string"},
		},
		Inputs: []models.Input{
			{ID: "allowed-input", Type: "file"},
		},
	})
	repo.Approve("testtool")

	wm := workspace.NewManager(workspace.Config{BasePath: t.TempDir()})
	handler := NewExecutionHandler(nil, wm, nil, repo, nil, nil)

	reqPayload := map[string]interface{}{
		"tool_id": "testtool",
		"job_id":  "job_123",
		"values": map[string]interface{}{
			"allowed-flag":   "value1",
			"unknown-flag":   "value2",
			"unknown-input":  "value3",
		},
	}
	body, _ := json.Marshal(reqPayload)
	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.HandleExecute(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 Bad Request, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	errors, ok := resp["errors"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected errors object in response")
	}

	if _, ok := errors["unknown-flag"]; !ok {
		t.Errorf("Expected validation error for unknown-flag")
	}
	if _, ok := errors["unknown-input"]; !ok {
		t.Errorf("Expected validation error for unknown-input")
	}
	if _, ok := errors["allowed-flag"]; ok {
		t.Errorf("Did not expect validation error for allowed-flag")
	}
}
