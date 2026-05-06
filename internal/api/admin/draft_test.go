package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/postfix/cleargate/internal/llm"
	"github.com/postfix/cleargate/internal/repository"
)

func TestAdminHandlers(t *testing.T) {
	repo, err := repository.NewToolSpecRepository("?access_mode=READ_WRITE")
	if err != nil {
		t.Fatalf("Failed to init repo: %v", err)
	}
	defer repo.Close()

	assistant := llm.NewMockAssistant(`
apiVersion: cleargate.dev/v1
metadata:
  name: testtool
  version: "1.0.0"
`)

	handler := NewAdminHandler(assistant, repo)

	// 1. Create Draft
	body := []byte(`{"help_text": "test"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/admin/tools/draft", bytes.NewReader(body))
	rr := httptest.NewRecorder()
	handler.HandleCreateDraft(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected 201 Created, got %v", rr.Code)
	}

	// 2. List Drafts
	req = httptest.NewRequest(http.MethodGet, "/api/admin/tools/drafts", nil)
	rr = httptest.NewRecorder()
	handler.HandleListDrafts(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", rr.Code)
	}

	var drafts []repository.ToolSpecRecord
	json.NewDecoder(rr.Body).Decode(&drafts)

	if len(drafts) != 1 {
		t.Fatalf("Expected 1 draft, got %d", len(drafts))
	}
	if drafts[0].ID != "testtool-1.0.0" {
		t.Errorf("Expected draft ID 'testtool-1.0.0', got %s", drafts[0].ID)
	}

	// 3. Approve
	req = httptest.NewRequest(http.MethodPost, "/api/admin/tools/testtool-1.0.0/approve", nil)
	rr = httptest.NewRecorder()
	handler.HandleApproveDraft(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected 200 OK, got %v", rr.Code)
	}

	// 4. List Drafts again (should be 0)
	req = httptest.NewRequest(http.MethodGet, "/api/admin/tools/drafts", nil)
	rr = httptest.NewRecorder()
	handler.HandleListDrafts(rr, req)

	json.NewDecoder(rr.Body).Decode(&drafts)
	if len(drafts) != 0 {
		t.Errorf("Expected 0 drafts after approval, got %d", len(drafts))
	}
}
