package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/postfix/cleargate/internal/models"
	"github.com/postfix/cleargate/internal/repository"
)

type MockAssistantInvalid struct{}

func (m *MockAssistantInvalid) GenerateDraft(ctx context.Context, prompt string) (*models.ToolSpec, error) {
	// Returns a ToolSpec missing the required Runtime.TimeoutSeconds and Metadata.Name
	return &models.ToolSpec{
		APIVersion: "v1",
		Kind:       "Tool",
		Metadata: models.Metadata{
			Version: "1.0",
		},
		Runtime: models.Runtime{
			Executable: "nmap",
		},
	}, nil
}

func TestDraftValidation(t *testing.T) {
	repo, _ := repository.NewToolSpecRepository("?access_mode=READ_WRITE")
	defer repo.Close()

	assistant := &MockAssistantInvalid{}
	handler := NewAdminHandler(assistant, repo, nil)

	reqBody := `{"prompt": "Make an nmap tool"}`
	req := httptest.NewRequest(http.MethodPost, "/admin/draft", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	handler.HandleCreateDraft(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected 400 Bad Request, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	errors, ok := resp["errors"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected errors object in response")
	}

	if _, hasMetaName := errors["ToolSpec.Metadata.Name"]; !hasMetaName {
		t.Errorf("Expected validation error for ToolSpec.Metadata.Name, got: %v", errors)
	}
	if _, hasTimeout := errors["ToolSpec.Runtime.TimeoutSeconds"]; !hasTimeout {
		t.Errorf("Expected validation error for ToolSpec.Runtime.TimeoutSeconds, got: %v", errors)
	}
}
