package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/postfix/cleargate/internal/workspace"
)

// UploadHandler handles streaming multipart uploads directly to disk.
type UploadHandler struct {
	workspaceManager *workspace.Manager
}

func NewUploadHandler(wm *workspace.Manager) *UploadHandler {
	return &UploadHandler{workspaceManager: wm}
}

// HandleUpload streams a file to the job's input directory.
// Expected query param or path var: job_id
func (h *UploadHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	// For MVP, we'll get job_id from query params.
	// In a real framework, this would be a path variable e.g., /api/jobs/{job_id}/upload
	jobID := r.URL.Query().Get("job_id")
	if jobID == "" {
		http.Error(w, "missing job_id parameter", http.StatusBadRequest)
		return
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "invalid multipart request", http.StatusBadRequest)
		return
	}

	inputDir := h.workspaceManager.GetPath(jobID, "input")
	// Ensure the workspace exists
	if _, err := os.Stat(inputDir); os.IsNotExist(err) {
		http.Error(w, "job workspace not found", http.StatusNotFound)
		return
	}

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, "error reading multipart form", http.StatusInternalServerError)
			return
		}

		// Only handle file uploads
		if part.FileName() == "" {
			continue
		}

		// Sanitize filename to prevent directory traversal
		safeFileName := filepath.Base(part.FileName())
		destPath := filepath.Join(inputDir, safeFileName)

		// Check if it's strictly within the input directory
		if filepath.Dir(destPath) != inputDir {
			http.Error(w, "invalid file path", http.StatusBadRequest)
			return
		}

		dst, err := os.Create(destPath)
		if err != nil {
			http.Error(w, "failed to create file on disk", http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(dst, part)
		dst.Close()

		if err != nil {
			http.Error(w, "failed to write file to disk", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Upload successful for job %s\n", jobID)
}
