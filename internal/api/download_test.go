package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/postfix/cleargate/internal/workspace"
)

func TestHandleDownload(t *testing.T) {
	tempBase := t.TempDir()
	wm := workspace.NewManager(workspace.Config{BasePath: tempBase})
	jobID := "job_456"
	wm.InitializeWorkspace(jobID)

	outputDir := wm.GetPath(jobID, "output")
	validFile := filepath.Join(outputDir, "result.json")
	os.WriteFile(validFile, []byte(`{"status":"ok"}`), 0644)

	handler := NewDownloadHandler(wm)

	tests := []struct {
		name         string
		jobID        string
		filename     string
		expectedCode int
	}{
		{"Valid file", jobID, "result.json", http.StatusOK},
		{"Missing job_id", "", "result.json", http.StatusBadRequest},
		{"Missing filename", jobID, "", http.StatusBadRequest},
		{"Path traversal absolute", jobID, "/etc/passwd", http.StatusBadRequest},
		{"Path traversal relative", jobID, "../metadata/job.json", http.StatusBadRequest},
		{"File not found", jobID, "missing.txt", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/download?job_id="+tt.jobID+"&filename="+tt.filename, nil)
			rr := httptest.NewRecorder()

			handler.HandleDownload(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("HandleDownload() status code = %v, expected %v", rr.Code, tt.expectedCode)
			}
		})
	}
}
