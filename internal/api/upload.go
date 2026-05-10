package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/postfix/cleargate/internal/job"
	"github.com/postfix/cleargate/internal/models"
	"github.com/postfix/cleargate/internal/repository"
	"github.com/postfix/cleargate/internal/workspace"
)

// UploadHandler handles streaming multipart uploads directly to disk.
type UploadHandler struct {
	workspaceManager *workspace.Manager
	repo             *repository.ToolSpecRepository
	registry         *job.Registry
}

func NewUploadHandler(wm *workspace.Manager, repo *repository.ToolSpecRepository, reg *job.Registry) *UploadHandler {
	return &UploadHandler{workspaceManager: wm, repo: repo, registry: reg}
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

	// Look up the job to find its tool_id
	jobRecord := h.registry.Get(jobID)
	if jobRecord == nil {
		http.Error(w, "unknown job_id", http.StatusNotFound)
		return
	}

	// Load the ToolSpec to get input constraints
	tsRecord, err := h.repo.GetByID(jobRecord.ToolID)
	if err != nil {
		http.Error(w, "tool not found for job", http.StatusNotFound)
		return
	}

	var spec models.ToolSpec
	if err := yaml.Unmarshal([]byte(tsRecord.Content), &spec); err != nil {
		http.Error(w, "invalid toolspec", http.StatusInternalServerError)
		return
	}

	// Build a lookup map for input constraints by ID
	inputMap := make(map[string]*models.Input)
	for i := range spec.Inputs {
		inputMap[spec.Inputs[i].ID] = &spec.Inputs[i]
	}

	reader, err := r.MultipartReader()
	if err != nil {
		http.Error(w, "invalid multipart request", http.StatusBadRequest)
		return
	}

	// Ensure the workspace exists
	if _, err := h.workspaceManager.InitializeWorkspace(jobID); err != nil {
		http.Error(w, "failed to initialize job workspace", http.StatusInternalServerError)
		return
	}
	inputDir := h.workspaceManager.GetPath(jobID, "input")

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

		// Look up the input definition by form field name
		inputDef, ok := inputMap[part.FormName()]
		if !ok {
			http.Error(w, fmt.Sprintf("undeclared input: %s", part.FormName()), http.StatusBadRequest)
			return
		}

		// Enforce allowed extensions
		if len(inputDef.AllowedExtensions) > 0 {
			ext := strings.ToLower(filepath.Ext(part.FileName()))
			allowed := false
			for _, a := range inputDef.AllowedExtensions {
				if strings.ToLower(a) == ext {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, fmt.Sprintf("file extension %q not allowed for input %s", ext, inputDef.ID), http.StatusBadRequest)
				return
			}
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

		// Enforce size limit via io.LimitReader
		var written int64
		if inputDef.MaxSizeMB > 0 {
			limitBytes := int64(inputDef.MaxSizeMB) * 1024 * 1024
			written, err = io.Copy(dst, io.LimitReader(part, limitBytes+1))
			dst.Close()
			if err != nil {
				os.Remove(destPath)
				http.Error(w, "failed to write file to disk", http.StatusInternalServerError)
				return
			}
			if written > limitBytes {
				os.Remove(destPath)
				http.Error(w, fmt.Sprintf("file exceeds max size of %dMB for input %s", inputDef.MaxSizeMB, inputDef.ID), http.StatusBadRequest)
				return
			}
		} else {
			written, err = io.Copy(dst, part)
			dst.Close()
			if err != nil {
				http.Error(w, "failed to write file to disk", http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Upload successful for job %s\n", jobID)
}
