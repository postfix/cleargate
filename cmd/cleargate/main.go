package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/postfix/cleargate/internal/api"
	"github.com/postfix/cleargate/internal/api/admin"
	"github.com/postfix/cleargate/internal/job"
	"github.com/postfix/cleargate/internal/llm"
	"github.com/postfix/cleargate/internal/repository"
	"github.com/postfix/cleargate/internal/runtime"
	"github.com/postfix/cleargate/internal/workspace"
)

func main() {
	toolsDir := flag.String("tools-dir", "./tools", "Directory containing ToolSpec YAMLs")
	flag.Parse()

	log.Println("Starting ClearGate Server...")

	// 1. Load Configuration & Initialize Dependencies
	dbPath := "cleargate.db" // In MVP, could use an env var or a config file
	workspaceCfg := workspace.DefaultConfig()
	workspaceManager := workspace.NewManager(workspaceCfg)

	// ToolSpec DB
	repo, err := repository.NewToolSpecRepository(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer repo.Close()

	if err := repo.SyncFromDirectory(*toolsDir); err != nil {
		log.Printf("Warning: failed to sync tools directory: %v", err)
	}

	// Preset and Audit DBs (reusing DuckDB connection from repo for simplicity, wait, Repo encapsulates *sql.DB but doesn't expose it. We should just open another or expose it.
	// Actually, wait, ToolSpecRepository has db private. Let's see if we can use the same path.)
	// Wait, duckdb only supports single writer per file unless read_only is true for others.
	// So we should expose the *sql.DB from ToolSpecRepository or open a new one?
	// It's better to use a single connection. We can just expose `repo.DB()`. Let's add that to toolspec_repo.go later, or open a separate db file, e.g. "cleargate-presets.db". For MVP, let's open "cleargate.db" again? No, DuckDB will error.
	// I'll update toolspec_repo to expose DB() via a getter or just make db public.
	// Actually, I can just create separate DB files for MVP to avoid touching toolspec_repo, or I can update toolspec_repo.go.
	// Let's assume I will update toolspec_repo.go to expose DB().
	presetRepo, err := repository.NewPresetRepository(repo.DB())
	if err != nil {
		log.Fatalf("Failed to initialize preset repo: %v", err)
	}
	auditRepo, err := repository.NewAuditRepository(repo.DB())
	if err != nil {
		log.Fatalf("Failed to initialize audit repo: %v", err)
	}

	// Podman Runtime
	var runtimeClient runtime.ContainerRuntime
	podmanClient, err := runtime.NewPodmanRuntime()
	if err != nil {
		log.Printf("Warning: failed to connect to Podman (%v). Falling back to DummyRuntime.", err)
		runtimeClient = &DummyRuntime{}
	} else {
		log.Println("Successfully connected to Podman socket.")
		runtimeClient = podmanClient
	}

	// Logger & Registry
	jobLogger := job.NewLogger(workspaceManager)
	jobRegistry := job.NewRegistry()

	// LLM Assistant (Using Mock for MVP)
	var assistant llm.TemplateAssistant = llm.NewMockAssistant(`
apiVersion: cleargate.dev/v1alpha1
kind: ToolSpec
metadata:
  name: demo-tool
  version: "1.0.0"
  description: "A generated demo tool"
  owner: "admin"
runtime:
  executable: "echo"
  argv0: "echo"
`)

	// 2. Instantiate Handlers
	uploadHandler := api.NewUploadHandler(workspaceManager)
	downloadHandler := api.NewDownloadHandler(workspaceManager)
	executeHandler := api.NewExecutionHandler(runtimeClient, workspaceManager, jobLogger, repo, jobRegistry, auditRepo)
	catalogHandler := api.NewCatalogHandler(repo)
	presetHandler := api.NewPresetHandler(presetRepo)
	
	// Admin handler requires assistant. Passing nil is risky, but for MVP we assume it's set or it returns 500.
	// Actually we should create a dummy assistant if it's nil, but the interface checks will catch it.
	adminHandler := admin.NewAdminHandler(assistant, repo, auditRepo)

	// 3. Register Routes
	mux := http.NewServeMux()
	
	mux.HandleFunc("POST /api/upload", uploadHandler.HandleUpload)
	mux.HandleFunc("GET /api/download", downloadHandler.HandleDownload)
	mux.HandleFunc("POST /api/execute", executeHandler.HandleExecute)
	mux.HandleFunc("GET /api/jobs/{id}/events", executeHandler.HandleEvents)
	mux.HandleFunc("GET /api/jobs", executeHandler.HandleListJobs)
	mux.HandleFunc("GET /api/catalog", catalogHandler.HandleListCatalog)
	
	mux.HandleFunc("POST /api/presets", presetHandler.HandleSavePreset)
	mux.HandleFunc("GET /api/presets", presetHandler.HandleListPresets)
	mux.HandleFunc("DELETE /api/presets", presetHandler.HandleDeletePreset)

	mux.HandleFunc("POST /api/admin/drafts", adminHandler.HandleCreateDraft)
	mux.HandleFunc("GET /api/admin/drafts", adminHandler.HandleListDrafts)
	mux.HandleFunc("POST /api/admin/tools/{id}/approve", adminHandler.HandleApproveDraft)
	mux.HandleFunc("GET /api/admin/audit", adminHandler.HandleListAuditLogs)

	// 4. Serve Static SPA with React Router fallback
	spaDir := "./web/dist"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Try to serve the exact file first
		path := filepath.Join(spaDir, r.URL.Path)
		info, err := os.Stat(path)
		if err == nil && !info.IsDir() {
			http.ServeFile(w, r, path)
			return
		}
		// Fallback: serve index.html so React Router handles the route
		http.ServeFile(w, r, filepath.Join(spaDir, "index.html"))
	})

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}

// DummyRuntime implements runtime.ContainerRuntime for MVP testing
type DummyRuntime struct{}

func (d *DummyRuntime) PullImage(ctx context.Context, image string) error { return nil }
func (d *DummyRuntime) Create(ctx context.Context, req runtime.CreateContainerRequest) (runtime.ContainerID, error) {
	return runtime.ContainerID("dummy-id"), nil
}
func (d *DummyRuntime) Start(ctx context.Context, id runtime.ContainerID) error { return nil }
func (d *DummyRuntime) Wait(ctx context.Context, id runtime.ContainerID) (int, error) { return 0, nil }
func (d *DummyRuntime) Logs(ctx context.Context, id runtime.ContainerID) (<-chan runtime.LogEvent, error) {
	ch := make(chan runtime.LogEvent)
	go func() {
		defer close(ch)
		ch <- runtime.LogEvent{Stream: "stdout", Data: []byte("Dummy log execution\n")}
	}()
	return ch, nil
}
