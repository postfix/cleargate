package admin

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/postfix/cleargate/internal/llm"
	"github.com/postfix/cleargate/internal/repository"
)

type AdminHandler struct {
	assistant llm.TemplateAssistant
	repo      *repository.ToolSpecRepository
}

func NewAdminHandler(assistant llm.TemplateAssistant, repo *repository.ToolSpecRepository) *AdminHandler {
	return &AdminHandler{assistant: assistant, repo: repo}
}

type DraftRequest struct {
	HelpText string `json:"help_text"`
}

func (h *AdminHandler) HandleCreateDraft(w http.ResponseWriter, r *http.Request) {
	var req DraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	spec, err := h.assistant.GenerateDraft(r.Context(), req.HelpText)
	if err != nil {
		http.Error(w, "failed to generate draft", http.StatusInternalServerError)
		return
	}

	if err := h.repo.SaveDraft(spec); err != nil {
		http.Error(w, "failed to save draft", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "draft_created"})
}

func (h *AdminHandler) HandleListDrafts(w http.ResponseWriter, r *http.Request) {
	drafts, err := h.repo.ListDrafts()
	if err != nil {
		http.Error(w, "failed to list drafts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drafts)
}

func (h *AdminHandler) HandleApproveDraft(w http.ResponseWriter, r *http.Request) {
	// Simple path parsing since we aren't using an external router in MVP
	// Expecting: /api/admin/tools/{id}/approve
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	id := parts[4]

	if err := h.repo.Approve(id); err != nil {
		http.Error(w, "failed to approve draft or not found", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "approved"})
}
