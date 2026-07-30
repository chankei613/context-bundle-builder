// Package resolve は ContextRef（file/url/obsidian_note/task_output）を実際に取得し、
// ResolvedItem に変換する。1件の失敗でBundle全体を止めない方針を徹底する（docs/spec.md参照）。
package resolve

import (
	"context"
	"time"

	"github.com/chankei613/context-bundle-builder/internal/db"
	"gorm.io/gorm"
)

type ResolvedItem struct {
	Ref        db.ContextRef `json:"ref"`
	Content    string        `json:"content"`
	ResolvedAt time.Time     `json:"resolved_at"`
}

type ResolveError struct {
	Ref     db.ContextRef `json:"ref"`
	Message string        `json:"message"`
}

// Options は解決エンジンの実行時設定。全てオプショナルで、未設定のコネクタは
// ハードエラーではなく ResolveError として記録され、他のrefsの解決は継続する。
type Options struct {
	ObsidianVaultRoot string        // 空なら obsidian_note は未設定エラー
	TaskOutputBaseURL string        // 空なら task_output は未接続扱い
	HTTPTimeout       time.Duration // 0ならデフォルト10秒
	MaxBytes          int64         // 0ならデフォルト2MB
	Cache             *gorm.DB      // nilならキャッシュ無効
	CacheTTL          time.Duration // 0ならデフォルト15分
}

func (o Options) withDefaults() Options {
	if o.HTTPTimeout == 0 {
		o.HTTPTimeout = 10 * time.Second
	}
	if o.MaxBytes == 0 {
		o.MaxBytes = 2 << 20 // 2MB
	}
	if o.CacheTTL == 0 {
		o.CacheTTL = 15 * time.Minute
	}
	return o
}

// Result はBundle全体の解決結果。
type Result struct {
	Items           []ResolvedItem
	Errors          []ResolveError
	PreviewText     string
	CharCount       int
	EstimatedTokens int
}

// Resolve は refs を順番に解決する。1件失敗しても他は続行する。
func Resolve(ctx context.Context, refs []db.ContextRef, opts Options) Result {
	opts = opts.withDefaults()

	res := Result{
		Items:  make([]ResolvedItem, 0, len(refs)),
		Errors: make([]ResolveError, 0),
	}

	var preview string
	for _, ref := range refs {
		content, err := resolveOne(ctx, ref, opts)
		if err != nil {
			res.Errors = append(res.Errors, ResolveError{Ref: ref, Message: err.Error()})
			continue
		}
		item := ResolvedItem{Ref: ref, Content: content, ResolvedAt: time.Now()}
		res.Items = append(res.Items, item)
		if preview != "" {
			preview += "\n\n---\n\n"
		}
		preview += formatSourceHeader(ref) + "\n" + content
	}

	res.PreviewText = preview
	res.CharCount = len([]rune(preview))
	res.EstimatedTokens = estimateTokens(preview)
	return res
}

func resolveOne(ctx context.Context, ref db.ContextRef, opts Options) (string, error) {
	if content, hit := cacheLookup(opts.Cache, ref.Kind, ref.Ref); hit {
		return content, nil
	}

	content, err := dispatch(ctx, ref, opts)
	if err != nil {
		return "", err
	}

	cacheStore(opts.Cache, ref.Kind, ref.Ref, content, opts.CacheTTL)
	return content, nil
}

func dispatch(ctx context.Context, ref db.ContextRef, opts Options) (string, error) {
	switch ref.Kind {
	case db.RefKindFile:
		return resolveFile(ref.Ref, opts)
	case db.RefKindURL:
		return resolveURL(ctx, ref.Ref, opts)
	case db.RefKindObsidianNote:
		return resolveObsidianNote(ref.Ref, opts)
	case db.RefKindTaskOutput:
		return resolveTaskOutput(ctx, ref.Ref, opts)
	default:
		return "", errUnknownKind
	}
}

func formatSourceHeader(ref db.ContextRef) string {
	return "## [" + string(ref.Kind) + "] " + ref.Ref
}

// estimateTokens は4文字≒1トークンの簡易換算（外部API呼び出し不要）。
func estimateTokens(s string) int {
	n := len([]rune(s))
	if n == 0 {
		return 0
	}
	return (n + 3) / 4
}
