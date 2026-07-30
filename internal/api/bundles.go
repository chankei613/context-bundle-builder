package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/chankei613/context-bundle-builder/internal/db"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

var errBundleNotFound = &apiError{"bundle not found"}

type CreateBundleInput struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Refs        []db.ContextRef `json:"refs"`
}

// CreateBundle は新しいBundleを作成する（HTTP・ネイティブバインディング共用）。
func (s *Server) CreateBundle(in CreateBundleInput) (db.Bundle, error) {
	if in.Name == "" {
		return db.Bundle{}, errNameRequired
	}
	now := time.Now()
	b := db.Bundle{
		ID:          uuid.NewString(),
		Name:        in.Name,
		Description: in.Description,
		Refs:        in.Refs,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if b.Refs == nil {
		b.Refs = []db.ContextRef{}
	}
	if err := s.DB.Create(&b).Error; err != nil {
		return db.Bundle{}, err
	}
	return b, nil
}

func (s *Server) ListBundles() ([]db.Bundle, error) {
	var bundles []db.Bundle
	err := s.DB.Order("updated_at desc").Find(&bundles).Error
	return bundles, err
}

func (s *Server) GetBundle(id string) (db.Bundle, error) {
	var b db.Bundle
	if err := s.DB.First(&b, "id = ?", id).Error; err != nil {
		return db.Bundle{}, errBundleNotFound
	}
	return b, nil
}

type UpdateBundleInput struct {
	Name        *string          `json:"name"`
	Description *string          `json:"description"`
	Refs        *[]db.ContextRef `json:"refs"`
}

func (s *Server) UpdateBundle(id string, in UpdateBundleInput) (db.Bundle, error) {
	b, err := s.GetBundle(id)
	if err != nil {
		return db.Bundle{}, err
	}
	if in.Name != nil {
		b.Name = *in.Name
	}
	if in.Description != nil {
		b.Description = *in.Description
	}
	if in.Refs != nil {
		b.Refs = *in.Refs
	}
	b.UpdatedAt = time.Now()
	if err := s.DB.Save(&b).Error; err != nil {
		return db.Bundle{}, err
	}
	return b, nil
}

func (s *Server) DeleteBundle(id string) error {
	res := s.DB.Delete(&db.Bundle{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errBundleNotFound
	}
	return nil
}

// DuplicateBundle は既存Bundleを複製して別名で保存する（テンプレート化の起点）。
func (s *Server) DuplicateBundle(id string) (db.Bundle, error) {
	orig, err := s.GetBundle(id)
	if err != nil {
		return db.Bundle{}, err
	}
	now := time.Now()
	dup := db.Bundle{
		ID:          uuid.NewString(),
		Name:        orig.Name + " (copy)",
		Description: orig.Description,
		Refs:        orig.Refs,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.DB.Create(&dup).Error; err != nil {
		return db.Bundle{}, err
	}
	return dup, nil
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpCreateBundle(w http.ResponseWriter, r *http.Request) {
	var body CreateBundleInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	b, err := s.CreateBundle(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}

func (s *Server) httpListBundles(w http.ResponseWriter, r *http.Request) {
	bundles, err := s.ListBundles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, bundles)
}

func (s *Server) httpGetBundle(w http.ResponseWriter, r *http.Request) {
	b, err := s.GetBundle(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) httpUpdateBundle(w http.ResponseWriter, r *http.Request) {
	var body UpdateBundleInput
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	b, err := s.UpdateBundle(chi.URLParam(r, "id"), body)
	if err != nil {
		if err == errBundleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, b)
}

func (s *Server) httpDeleteBundle(w http.ResponseWriter, r *http.Request) {
	if err := s.DeleteBundle(chi.URLParam(r, "id")); err != nil {
		if err == errBundleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) httpDuplicateBundle(w http.ResponseWriter, r *http.Request) {
	b, err := s.DuplicateBundle(chi.URLParam(r, "id"))
	if err != nil {
		if err == errBundleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, b)
}
