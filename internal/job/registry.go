package job

import (
	"sync"
	"time"
)

// JobStatus represents the current state of a job.
type JobStatus string

const (
	JobStatusRunning   JobStatus = "running"
	JobStatusSucceeded JobStatus = "succeeded"
	JobStatusFailed    JobStatus = "failed"
)

// JobRecord tracks a single execution.
type JobRecord struct {
	JobID     string    `json:"job_id"`
	ToolID    string    `json:"tool_id"`
	Status    JobStatus `json:"status"`
	ExitCode  int       `json:"exit_code"`
	StartedAt time.Time `json:"started_at"`
}

// Registry is an in-memory store of active and recent jobs.
type Registry struct {
	mu   sync.RWMutex
	jobs map[string]*JobRecord // keyed by job ID
}

// NewRegistry creates a new job registry.
func NewRegistry() *Registry {
	return &Registry{
		jobs: make(map[string]*JobRecord),
	}
}

// Register adds a new running job to the registry.
func (r *Registry) Register(jobID, toolID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.jobs[jobID] = &JobRecord{
		JobID:     jobID,
		ToolID:    toolID,
		Status:    JobStatusRunning,
		StartedAt: time.Now(),
	}
}

// Complete marks a job as finished.
func (r *Registry) Complete(jobID string, exitCode int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if j, ok := r.jobs[jobID]; ok {
		if exitCode == 0 {
			j.Status = JobStatusSucceeded
		} else {
			j.Status = JobStatusFailed
		}
		j.ExitCode = exitCode
	}
}

// GetByTool returns recent jobs for a given tool ID.
func (r *Registry) GetByTool(toolID string) []JobRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []JobRecord
	for _, j := range r.jobs {
		if j.ToolID == toolID {
			result = append(result, *j)
		}
	}
	return result
}

// Get returns a single job by ID.
func (r *Registry) Get(jobID string) *JobRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if j, ok := r.jobs[jobID]; ok {
		cpy := *j
		return &cpy
	}
	return nil
}
