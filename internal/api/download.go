package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/postfix/cleargate/internal/workspace"
)

// DownloadHandler handles artifact downloads securely.
type DownloadHandler struct {
	workspaceManager *workspace.Manager
}

func NewDownloadHandler(wm *workspace.Manager) *DownloadHandler {
	return &DownloadHandler{workspaceManager: wm}
}

// HandleDownload serves an artifact file from the job's output directory.
// Expected query params or path vars: job_id, filename
func (h *DownloadHandler) HandleDownload(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("job_id")
	filename := r.URL.Query().Get("filename")

	if jobID == "" || filename == "" {
		http.Error(w, "missing job_id or filename", http.StatusBadRequest)
		return
	}

	outputDir := h.workspaceManager.GetPath(jobID, "output")
	
	// Prevent directory traversal: Clean the path and check if it has directory separators
	// We only allow downloading direct children of the output directory for this MVP.
	cleanName := filepath.Clean(filename)
	if strings.Contains(cleanName, string(filepath.Separator)) || cleanName == ".." || cleanName == "." {
		http.Error(w, "invalid filename: directory traversal is not allowed", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(outputDir, cleanName)

	// Check if file exists
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) || info.IsDir() {
		http.Error(w, "artifact not found", http.StatusNotFound)
		return
	}

	http.ServeFile(w, r, filePath)
}
