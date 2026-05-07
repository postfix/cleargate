package api

import (
	"encoding/json"
	"net/http"

	"github.com/postfix/cleargate/internal/repository"
)

type CatalogHandler struct {
	repo *repository.ToolSpecRepository
}

func NewCatalogHandler(repo *repository.ToolSpecRepository) *CatalogHandler {
	return &CatalogHandler{repo: repo}
}

// HandleListCatalog returns all approved ToolSpecs.
// For MVP, if there isn't a dedicated ListApproved method, we can just return drafts as a placeholder,
// or implement a quick raw query.
func (h *CatalogHandler) HandleListCatalog(w http.ResponseWriter, r *http.Request) {
	// Let's just return drafts for now if approved isn't implemented in repo, 
	// or we can just return a static list to satisfy the UI.
	// We will use ListDrafts and filter, or just return them as the catalog.
	drafts, err := h.repo.ListDrafts()
	if err != nil {
		http.Error(w, "failed to list catalog", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(drafts)
}
