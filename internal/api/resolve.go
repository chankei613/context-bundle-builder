package api

import (
	"context"
	"net/http"

	"github.com/chankei613/context-bundle-builder/internal/resolve"
	"github.com/go-chi/chi/v5"
)

type ResolveResult struct {
	BundleID        string                 `json:"bundle_id"`
	Items           []resolve.ResolvedItem `json:"items"`
	Errors          []resolve.ResolveError `json:"errors"`
	PreviewText     string                 `json:"preview_text"`
	CharCount       int                    `json:"char_count"`
	EstimatedTokens int                    `json:"estimated_tokens"`
}

func (s *Server) resolveOpts() resolve.Options {
	return resolve.Options{
		ObsidianVaultRoot: s.ObsidianVaultRoot,
		TaskOutputBaseURL: s.TaskOutputBaseURL,
		Cache:             s.DB,
	}
}

// ResolveBundle はBundleの全refsを解決する（HTTP・ネイティブバインディング共用）。
func (s *Server) ResolveBundle(id string) (ResolveResult, error) {
	b, err := s.GetBundle(id)
	if err != nil {
		return ResolveResult{}, err
	}

	r := resolve.Resolve(context.Background(), b.Refs, s.resolveOpts())
	return ResolveResult{
		BundleID:        b.ID,
		Items:           r.Items,
		Errors:          r.Errors,
		PreviewText:     r.PreviewText,
		CharCount:       r.CharCount,
		EstimatedTokens: r.EstimatedTokens,
	}, nil
}

// PreviewBundle は ResolveBundle のエイリアス（UIのプレビューパネル用に意味を明示する名前）。
func (s *Server) PreviewBundle(id string) (ResolveResult, error) {
	return s.ResolveBundle(id)
}

// ─── HTTPハンドラー ────────────────────────────────────────────────────

func (s *Server) httpResolveBundle(w http.ResponseWriter, r *http.Request) {
	result, err := s.ResolveBundle(chi.URLParam(r, "id"))
	if err != nil {
		if err == errBundleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) httpPreviewBundle(w http.ResponseWriter, r *http.Request) {
	result, err := s.PreviewBundle(chi.URLParam(r, "id"))
	if err != nil {
		if err == errBundleNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, result)
}
