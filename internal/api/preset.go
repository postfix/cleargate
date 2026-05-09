package api

import (
	"encoding/json"
	"net/http"

	"github.com/postfix/cleargate/internal/models"
	"github.com/postfix/cleargate/internal/repository"
)

type PresetHandler struct {
	repo *repository.PresetRepository
}

func NewPresetHandler(repo *repository.PresetRepository) *PresetHandler {
	return &PresetHandler{
		repo: repo,
	}
}

func (h *PresetHandler) HandleSavePreset(w http.ResponseWriter, r *http.Request) {
	var preset models.Preset
	if err := json.NewDecoder(r.Body).Decode(&preset); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Save(&preset); err != nil {
		http.Error(w, "failed to save preset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "preset_saved"})
}

func (h *PresetHandler) HandleListPresets(w http.ResponseWriter, r *http.Request) {
	toolID := r.URL.Query().Get("tool_id")
	
	presets, err := h.repo.ListByTool(toolID)
	if err != nil {
		http.Error(w, "failed to list presets", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(presets)
}

func (h *PresetHandler) HandleDeletePreset(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing preset id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "failed to delete preset", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}
