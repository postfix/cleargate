package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/postfix/cleargate/internal/models"
	"github.com/postfix/cleargate/internal/repository"
	"github.com/postfix/cleargate/internal/workspace"
)

func TestHandleUpload(t *testing.T) {
	tempBase := t.TempDir()
	wm := workspace.NewManager(workspace.Config{BasePath: tempBase})
	jobID := "job_123"
	wm.InitializeWorkspace(jobID)

	repo, err := repository.NewToolSpecRepository("?access_mode=READ_WRITE")
	if err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	defer repo.Close()

	// Seed a test ToolSpec so GetByID doesn't fail
	err = repo.SaveDraft(&models.ToolSpec{
		Metadata: models.Metadata{Name: "testtool", Version: "1.0"},
		Inputs: []models.Input{
			{ID: "file"},
		},
	})
	if err != nil {
		t.Fatalf("Failed to seed toolspec: %v", err)
	}
	repo.Approve("testtool")

	// Create a multipart body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("hello world"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload?job_id="+jobID+"&tool_id=testtool", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := NewUploadHandler(wm, repo, nil)
	handler.HandleUpload(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Verify file was written
	expectedFile := filepath.Join(tempBase, jobID, "input", "test.txt")
	content, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read expected file: %v", err)
	}

	if string(content) != "hello world" {
		t.Errorf("File content mismatch. got %q want %q", content, "hello world")
	}
}
