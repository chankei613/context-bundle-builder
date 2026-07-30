package api

import (
	"encoding/json"
	"net/http"

	"github.com/chankei613/context-bundle-builder/internal/db"
)

// GetSettings はAppSettingsを返す（未作成なら空値で作る）。
func (s *Server) GetSettings() (db.AppSettings, error) {
	var settings db.AppSettings
	err := s.DB.FirstOrCreate(&settings, db.AppSettings{ID: db.AppSettingsID}).Error
	if err != nil {
		return db.AppSettings{}, err
	}
	s.ObsidianVaultRoot = settings.ObsidianVaultRoot
	s.TaskOutputBaseURL = settings.TaskOutputBaseURL
	return settings, nil
}

type UpdateSettingsInput struct {
	ObsidianVaultRoot *string `json:"obsidian_vault_root"`
	TaskOutputBaseURL *string `json:"task_output_base_url"`
}

// UpdateSettings は設定を更新し、Serverの実行時設定にも即座に反映する
// （再起動しなくてもresolveの挙動に反映されるようにするため）。
func (s *Server) UpdateSettings(in UpdateSettingsInput) (db.AppSettings, error) {
	settings, err := s.GetSettings()
	if err != nil {
		return db.AppSettings{}, err
	}
	if in.ObsidianVaultRoot != nil {
		settings.ObsidianVaultRoot = *in.ObsidianVaultRoot
	}
	if in.TaskOutputBaseURL != nil {
		settings.TaskOutputBaseURL = *in.TaskOutputBaseURL
	}
	if err := s.DB.Save(&settings).Error; err != nil {
		return db.AppSettings{}, err
	}
	s.ObsidianVaultRoot = settings.ObsidianVaultRoot
	s.TaskOutputBaseURL = settings.TaskOutputBaseURL
	return settings, nil
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.GetSettings()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (s *Server) httpUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var body UpdateSettingsInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	settings, err := s.UpdateSettings(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}
