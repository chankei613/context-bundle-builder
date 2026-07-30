package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// Server は全ロジックの実体。HTTPハンドラーとWailsネイティブバインディングの
// 両方がこの同じ Server のメソッドを呼ぶことで、UIとAPIの挙動がズレないようにする。
type Server struct {
	DB *gorm.DB

	// ObsidianVaultRoot / TaskOutputBaseURL は resolve.Options に渡す実行時設定。
	// UI(Settings)から変更可能にする想定。空のままでも他のコネクタは動く。
	ObsidianVaultRoot string
	TaskOutputBaseURL string
}

func New(conn *gorm.DB) *Server {
	return &Server{DB: conn}
}

func (s *Server) Router() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1/keys", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB, "/api/v1/keys"))
		r.Post("/", s.httpIssueKey)
		r.Get("/", s.httpListKeys)
		r.Delete("/{id}", s.httpRevokeKey)
	})

	r.Route("/api/v1/bundles", func(r chi.Router) {
		r.Use(APIKeyAuth(s.DB))
		r.Post("/", s.httpCreateBundle)
		r.Get("/", s.httpListBundles)
		r.Get("/{id}", s.httpGetBundle)
		r.Patch("/{id}", s.httpUpdateBundle)
		r.Delete("/{id}", s.httpDeleteBundle)
		r.Post("/{id}/duplicate", s.httpDuplicateBundle)
		r.Post("/{id}/resolve", s.httpResolveBundle)
		r.Get("/{id}/preview", s.httpPreviewBundle)
	})

	return r
}

// NewRouter はcmd/cbbserve（単体HTTPサーバー）向けの簡易コンストラクタ。
func NewRouter(conn *gorm.DB) http.Handler {
	return New(conn).Router()
}
