// Package db はContext Bundle BuilderのGORMモデルとSQLite初期化を提供する。
// ContextRef / ResolvedContextItem は comet-taskAI の schema/types.ts を正規ソースとし、
// 新規に型を発明しない。docs/spec.md 参照。
package db

import "time"

type RefKind string

const (
	RefKindFile         RefKind = "file"
	RefKindURL          RefKind = "url"
	RefKindObsidianNote RefKind = "obsidian_note"
	RefKindTaskOutput   RefKind = "task_output"
)

// ContextRef は comet-taskAI/schema/types.ts の ContextRef と1:1対応する。
type ContextRef struct {
	Kind RefKind `json:"kind"`
	Ref  string  `json:"ref"`
}

// Bundle は名前付きの ContextRef の集合（テンプレート）。
type Bundle struct {
	ID          string       `gorm:"primaryKey" json:"id"`
	Name        string       `gorm:"index" json:"name"`
	Description string       `json:"description"`
	Refs        []ContextRef `gorm:"serializer:json" json:"refs"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// AgentKey — Bundle CRUD/解決APIを叩くためのAPIキー。ハッシュのみ保存する。
type AgentKey struct {
	ID         string     `gorm:"primaryKey" json:"id"`
	Name       string     `json:"name"`
	APIKeyHash string     `json:"-"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

// ResolveCacheEntry は同一refの再取得を避けるためのTTL付きキャッシュ。
type ResolveCacheEntry struct {
	ID         string    `gorm:"primaryKey" json:"-"`
	RefKind    RefKind   `gorm:"index:idx_ref,unique" json:"-"`
	RefValue   string    `gorm:"index:idx_ref,unique" json:"-"`
	Content    string    `json:"-"`
	Err        string    `json:"-"`
	ResolvedAt time.Time `json:"-"`
	ExpiresAt  time.Time `gorm:"index" json:"-"`
}
