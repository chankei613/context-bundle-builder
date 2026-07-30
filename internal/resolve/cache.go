package resolve

import (
	"time"

	"github.com/chankei613/context-bundle-builder/internal/db"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// cacheLookup は期限内のキャッシュがあれば内容とヒットしたかを返す。
// キャッシュ自体の読み取り失敗はキャッシュミス扱いにして解決を継続する。
func cacheLookup(conn *gorm.DB, kind db.RefKind, ref string) (string, bool) {
	if conn == nil {
		return "", false
	}
	var entry db.ResolveCacheEntry
	err := conn.Where("ref_kind = ? AND ref_value = ? AND expires_at > ?", kind, ref, time.Now()).
		First(&entry).Error
	if err != nil {
		return "", false
	}
	if entry.Err != "" {
		return "", false
	}
	return entry.Content, true
}

// cacheStore はキャッシュへ書き込む（同一ref+kindは上書き）。書き込み失敗は無視する
// （キャッシュは最適化であり、失敗しても解決結果自体には影響させない）。
func cacheStore(conn *gorm.DB, kind db.RefKind, ref, content string, ttl time.Duration) {
	if conn == nil {
		return
	}
	now := time.Now()
	var existing db.ResolveCacheEntry
	err := conn.Where("ref_kind = ? AND ref_value = ?", kind, ref).First(&existing).Error
	if err == nil {
		existing.Content = content
		existing.Err = ""
		existing.ResolvedAt = now
		existing.ExpiresAt = now.Add(ttl)
		conn.Save(&existing)
		return
	}
	entry := db.ResolveCacheEntry{
		ID:         uuid.NewString(),
		RefKind:    kind,
		RefValue:   ref,
		Content:    content,
		ResolvedAt: now,
		ExpiresAt:  now.Add(ttl),
	}
	conn.Create(&entry)
}
