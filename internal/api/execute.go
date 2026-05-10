package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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
	registry         *job.Registry
	auditRepo        *repository.AuditRepository
}

func NewExecutionHandler(r runtime.ContainerRuntime, wm *workspace.Manager, l *job.Logger, repo *repository.ToolSpecRepository, reg *job.Registry, auditRepo *repository.AuditRepository) *ExecutionHandler {
	return &ExecutionHandler{
		runtime:          r,
		workspaceManager: wm,
		logger:           l,
		repo:             repo,
		registry:         reg,
		auditRepo:        auditRepo,
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

	// Validate unknown flags
	allowedKeys := make(map[string]bool)
	for _, f := range spec.Flags {
		allowedKeys[f.ID] = true
	}
	for _, in := range spec.Inputs {
		allowedKeys[in.ID] = true
	}

	for key := range req.Values {
		if !allowedKeys[key] {
			http.Error(w, fmt.Sprintf("Unknown flag/input: %s", key), http.StatusBadRequest)
			return
		}
	}

	var cmdArgs []string

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

	isPositional := make(map[string]int)
	for _, p := range spec.Positionals {
		isPositional[p.Source] = p.Order
	}
	posArgs := make(map[int]string)

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
				if order, ok := isPositional[f.ID]; ok {
					posArgs[order] = str
				} else if f.FlagString != "" {
					cmdArgs = append(cmdArgs, f.FlagString, str)
				}
			}
		}
	}

	for _, in := range spec.Inputs {
		val, ok := req.Values[in.ID]
		if !ok {
			continue
		}
		if str, ok := val.(string); ok && str != "" {
			path := filepath.Join("/workspace/input", str)
			if order, ok := isPositional[in.ID]; ok {
				posArgs[order] = path
			} else if in.FlagString != "" {
				cmdArgs = append(cmdArgs, in.FlagString, path)
			}
		}
	}

	for _, out := range spec.Outputs {
		path := filepath.Join("/workspace/output", out.Path)
		if order, ok := isPositional[out.ID]; ok {
			posArgs[order] = path
		} else if out.FlagString != "" {
			cmdArgs = append(cmdArgs, out.FlagString, path)
		}
	}

	var orders []int
	for o := range posArgs {
		orders = append(orders, o)
	}
	sort.Ints(orders)
	for _, o := range orders {
		cmdArgs = append(cmdArgs, posArgs[o])
	}

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

	if spec.Runtime.Executable == "nmap" {
		containerReq.CapAdd = []string{"CAP_NET_RAW"}
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

	// Register job in the registry
	h.registry.Register(req.JobID, req.ToolID)

	timeout := spec.Runtime.TimeoutSeconds
	if timeout <= 0 {
		timeout = 300 // default 5 minutes
	}
	execCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)

	// Attach logs in background
	logChan, err := h.runtime.Logs(execCtx, containerID)
	if err == nil {
		go func() {
			defer cancel()
			var stdoutBytes, stderrBytes int64
			for ev := range logChan {
				if ev.Stream == "stdout" {
					stdoutBytes += int64(len(ev.Data))
					h.logger.LogStdout(req.JobID, ev.Data)
				} else {
					stderrBytes += int64(len(ev.Data))
					h.logger.LogStderr(req.JobID, ev.Data)
				}
			}
			// Wait for completion to get exit code
			exitCode, err := h.runtime.Wait(execCtx, containerID)
			if err != nil {
				exitCode = -1
				h.logger.LogStderr(req.JobID, []byte(fmt.Sprintf("wait error (possibly timeout): %v", err)))
			}
			h.registry.Complete(req.JobID, exitCode)
			h.logger.LogStatus(req.JobID, "complete", exitCode)
			
			h.auditRepo.Log(&models.AuditLog{
				JobID:     req.JobID,
				ToolID:    req.ToolID,
				ExitCode:  exitCode,
				CreatedAt: time.Now(),
			})

			// Write metadata.json to workspace
			h.writeJobMetadata(req.JobID, exitCode, stdoutBytes, stderrBytes)
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
				parts := strings.SplitN(string(ev.Data), ":", 2)
				status := "failed"
				var exitCode int = -1
				if len(parts) == 2 {
					if parts[0] == "complete" {
						fmt.Sscanf(parts[1], "%d", &exitCode)
						if exitCode == 0 {
							status = "succeeded"
						}
					}
				}
				
				payload = map[string]interface{}{
					"type": "complete",
					"status": status,
					"exitCode": exitCode,
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

// HandleCancelJob stops a running job by killing its container.
func (h *ExecutionHandler) HandleCancelJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	containerName := fmt.Sprintf("cleargate-job-%s", jobID)
	err := h.runtime.Stop(context.Background(), runtime.ContainerID(containerName))
	if err != nil {
		h.logger.LogStderr(jobID, []byte("Job cancelled by user.\n"))
		http.Error(w, fmt.Sprintf("failed to stop job: %v", err), http.StatusInternalServerError)
		return
	}

	h.logger.LogStatus(jobID, "complete", 137) // 137 is SIGKILL
	h.registry.Complete(jobID, 137)
	w.WriteHeader(http.StatusOK)
}

// HandleListJobs returns active/recent jobs, optionally filtered by tool_id.
func (h *ExecutionHandler) HandleListJobs(w http.ResponseWriter, r *http.Request) {
	toolID := r.URL.Query().Get("tool_id")
	var jobs []job.JobRecord
	if toolID != "" {
		jobs = h.registry.GetByTool(toolID)
	}
	if jobs == nil {
		jobs = []job.JobRecord{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// writeJobMetadata scans the output directory and writes metadata.json to the workspace root.
func (h *ExecutionHandler) writeJobMetadata(jobID string, exitCode int, stdoutBytes, stderrBytes int64) {
	outputDir := h.workspaceManager.GetPath(jobID, "output")
	var outputFiles []string

	entries, err := os.ReadDir(outputDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				outputFiles = append(outputFiles, e.Name())
			}
		}
	}
	if outputFiles == nil {
		outputFiles = []string{}
	}

	meta := models.JobMetadata{
		ExitCode:    exitCode,
		StdoutBytes: stdoutBytes,
		StderrBytes: stderrBytes,
		OutputFiles: outputFiles,
	}

	metaBytes, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return
	}

	metaPath := filepath.Join(h.workspaceManager.GetPath(jobID, ""), "metadata.json")
	os.WriteFile(metaPath, metaBytes, 0644)
}

// HandleJobMetadata serves the metadata.json for a completed job.
func (h *ExecutionHandler) HandleJobMetadata(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	metaPath := filepath.Join(h.workspaceManager.GetPath(jobID, ""), "metadata.json")
	data, err := os.ReadFile(metaPath)
	if err != nil {
		http.Error(w, "metadata not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
