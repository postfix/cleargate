package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/postfix/cleargate/internal/workspace"
)

func TestHandleUpload(t *testing.T) {
	tempBase := t.TempDir()
	wm := workspace.NewManager(workspace.Config{BasePath: tempBase})
	jobID := "job_123"
	wm.InitializeWorkspace(jobID)

	// Create a multipart body
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("hello world"))
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload?job_id="+jobID, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	handler := NewUploadHandler(wm)
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
