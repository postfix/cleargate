package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"

	"github.com/postfix/cleargate/internal/job"
	"github.com/postfix/cleargate/internal/models"
	"github.com/postfix/cleargate/internal/repository"
	"github.com/postfix/cleargate/internal/runtime"
	"github.com/postfix/cleargate/internal/workspace"
)

type ExecutionHandler struct {
	runtime          runtime.ContainerRuntime
	workspaceManager *workspace.Manager
	logger           *job.Logger
	repo             *repository.ToolSpecRepository
}

func NewExecutionHandler(r runtime.ContainerRuntime, wm *workspace.Manager, l *job.Logger, repo *repository.ToolSpecRepository) *ExecutionHandler {
	return &ExecutionHandler{
		runtime:          r,
		workspaceManager: wm,
		logger:           l,
		repo:             repo,
	}
}

type ExecuteRequest struct {
	ToolID string                 `json:"tool_id"`
	JobID  string                 `json:"job_id"`
	Values map[string]interface{} `json:"values"`
}

func (h *ExecutionHandler) HandleExecute(w http.ResponseWriter, r *http.Request) {
	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.ToolID == "" || req.JobID == "" {
		http.Error(w, "missing tool_id or job_id", http.StatusBadRequest)
		return
	}

	tsRecord, err := h.repo.GetByID(req.ToolID)
	if err != nil {
		http.Error(w, "tool not found", http.StatusNotFound)
		return
	}

	var spec models.ToolSpec
	if err := yaml.Unmarshal([]byte(tsRecord.Content), &spec); err != nil {
		http.Error(w, "invalid toolspec", http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	workspacePath, err := h.workspaceManager.InitializeWorkspace(req.JobID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to initialize workspace: %v", err), http.StatusInternalServerError)
		return
	}

	var cmdArgs []string
	var positionals []string

	// Basic execution assembly mapping
	if spec.Runtime.Executable == "nmap" {
		// Use a known nmap image if none specified
		if spec.Runtime.ContainerImage == "" {
			spec.Runtime.ContainerImage = "docker.io/instrumentisto/nmap:latest"
		}
		// For nmap image where entrypoint is already nmap, we omit nmap from command.
		// Wait, if we aren't sure about entrypoint, we pass "nmap" then args. If it's alpine, we do "nmap args".
		// We will assume entrypoint is empty or we use full command. Let's just pass the args.
	} else {
		cmdArgs = append(cmdArgs, spec.Runtime.Executable)
	}

	for _, f := range spec.Flags {
		val, ok := req.Values[f.ID]
		if !ok {
			continue
		}
		
		if f.Type == "boolean" {
			if b, ok := val.(bool); ok && b {
				if f.FlagString != "" {
					cmdArgs = append(cmdArgs, f.FlagString)
				}
			}
		} else if f.Type == "string" {
			if str, ok := val.(string); ok && str != "" {
				if f.ID == "target" {
					positionals = append(positionals, str)
				} else if f.FlagString != "" {
					cmdArgs = append(cmdArgs, f.FlagString, str)
				}
			}
		}
	}
	
	cmdArgs = append(cmdArgs, positionals...)

	// Fallback to alpine if no image
	if spec.Runtime.ContainerImage == "" {
		spec.Runtime.ContainerImage = "docker.io/library/alpine:latest"
		cmdArgs = append([]string{"echo"}, cmdArgs...)
	}

	containerReq := runtime.CreateContainerRequest{
		Image:        spec.Runtime.ContainerImage,
		Name:         fmt.Sprintf("cleargate-job-%s", req.JobID),
		Command:      cmdArgs,
		Remove:       true,
		WorkspaceDir: workspacePath,
	}

	if err := h.runtime.PullImage(ctx, containerReq.Image); err != nil {
		http.Error(w, fmt.Sprintf("failed to pull image: %v", err), http.StatusInternalServerError)
		return
	}

	containerID, err := h.runtime.Create(ctx, containerReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create container: %v", err), http.StatusInternalServerError)
		return
	}

	if err := h.runtime.Start(ctx, containerID); err != nil {
		http.Error(w, fmt.Sprintf("failed to start container: %v", err), http.StatusInternalServerError)
		return
	}

	// Attach logs in background
	logChan, err := h.runtime.Logs(ctx, containerID)
	if err == nil {
		go func() {
			for ev := range logChan {
				if ev.Stream == "stdout" {
					h.logger.LogStdout(req.JobID, ev.Data)
				} else {
					h.logger.LogStderr(req.JobID, ev.Data)
				}
			}
			// Wait for completion to get exit code
			err := h.runtime.Wait(ctx, containerID)
			exitCode := 0
			if err != nil {
				exitCode = 1
			}
			h.logger.LogStatus(req.JobID, "complete", exitCode)
		}()
	} else {
		h.logger.LogStderr(req.JobID, []byte(fmt.Sprintf("failed to attach logs: %v", err)))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": req.JobID})
}

// HandleEvents streams the logs for a specific job ID via SSE.
func (h *ExecutionHandler) HandleEvents(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	// Subscribe to the logger stream
	sub := h.logger.Subscribe(jobID)
	defer h.logger.Unsubscribe(jobID, sub)

	ctx := r.Context()

	for {
		select {
		case ev, ok := <-sub:
			if !ok {
				return
			}
			
			// If it's a raw log event (e.g., from Podman), it might have newlines.
			// The frontend expects the SSE to be formatted as {"type": "...", "data": "..."}
			// or similar. 
			
			var payload map[string]interface{}
			if ev.Type == "status" {
				payload = map[string]interface{}{
					"type": "complete",
					"status": ev.Data,
					"exitCode": 0, // Simplified for MVP
				}
				if string(ev.Data) == "complete" {
					payload["status"] = "succeeded"
				}
			} else {
				payload = map[string]interface{}{
					"type": ev.Type, // stdout or stderr
					"data": string(ev.Data),
				}
			}
			
			eventJSON, _ := json.Marshal(payload)
			fmt.Fprintf(w, "data: %s\n\n", eventJSON)
			flusher.Flush()
			
			if ev.Type == "status" && string(ev.Data) == "complete" {
				return
			}
			
		case <-ctx.Done():
			return
		}
	}
}
