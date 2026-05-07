package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/postfix/cleargate/internal/models"
)

type PresetHandler struct {
	mu      sync.RWMutex
	presets []models.Preset
}

func NewPresetHandler() *PresetHandler {
	return &PresetHandler{
		presets: make([]models.Preset, 0),
	}
}

func (h *PresetHandler) HandleSavePreset(w http.ResponseWriter, r *http.Request) {
	var preset models.Preset
	if err := json.NewDecoder(r.Body).Decode(&preset); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.mu.Lock()
	h.presets = append(h.presets, preset)
	h.mu.Unlock()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "preset_saved"})
}

func (h *PresetHandler) HandleListPresets(w http.ResponseWriter, r *http.Request) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.presets)
}
