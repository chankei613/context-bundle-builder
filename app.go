package main

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/chankei613/context-bundle-builder/internal/api"
	"github.com/chankei613/context-bundle-builder/internal/db"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// apiAddr はBundle CRUD/解決APIの待ち受けアドレス。将来的にAI Scheduler等の外部プロセスが
// 解決済みBundleを取得できるよう、ウインドウの表示/非表示に関わらずこのHTTPサーバーは動き続ける
// （UI自体はこのHTTPを経由せず、下記のネイティブバインディング経由で同じ *api.Server を直接呼ぶ）。
const apiAddr = "127.0.0.1:8422"

// App はWailsのバインディング。実処理は internal/api.Server が持っており、
// ここはWails固有の初期化・エラー通知と、UI向けのネイティブバインディングだけを担当する。
// 同じ Server を cmd/cbbserve のHTTP APIも使っているので、UIとAPIで挙動がズレない。
type App struct {
	ctx    context.Context
	server *api.Server
	srv    *http.Server
	ready  bool
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	dataDir := appDataDir()
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		runtime.LogErrorf(ctx, "data dir error: %s", err)
		return
	}

	conn, err := db.Init(filepath.Join(dataDir, "context-bundle-builder.db"))
	if err != nil {
		runtime.LogErrorf(ctx, "db init error: %s", err)
		return
	}
	a.server = api.New(conn)

	a.srv = &http.Server{Addr: apiAddr, Handler: a.server.Router()}
	go func() {
		if err := a.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			runtime.LogErrorf(ctx, "api server error: %s", err)
		}
	}()

	a.ready = true
	runtime.LogInfof(ctx, "Context Bundle Builder ready (api: http://%s, data: %s)", apiAddr, dataDir)
}

func (a *App) shutdown(ctx context.Context) {
	if a.srv != nil {
		_ = a.srv.Close()
	}
	if a.server != nil {
		if sqlDB, err := a.server.DB.DB(); err == nil {
			_ = sqlDB.Close()
		}
	}
}

var errNotReady = &notReadyError{}

type notReadyError struct{}

func (e *notReadyError) Error() string { return "app not ready — check startup logs" }

// ─── フロントエンドへ公開するメソッド ──────────────────────────────────────────

// GetAppVersion はアプリのバージョン文字列を返す。
func (a *App) GetAppVersion() string {
	return AppVersion
}

// GetAPIURL は外部プロセスがBundle CRUD/解決APIを叩く先のベースURLを返す（Settings画面に表示する）。
func (a *App) GetAPIURL() string {
	return "http://" + apiAddr
}

func (a *App) ListBundles() ([]db.Bundle, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListBundles()
}

func (a *App) GetBundle(id string) (db.Bundle, error) {
	if !a.ready {
		return db.Bundle{}, errNotReady
	}
	return a.server.GetBundle(id)
}

func (a *App) CreateBundle(name string, description string, refs []db.ContextRef) (db.Bundle, error) {
	if !a.ready {
		return db.Bundle{}, errNotReady
	}
	return a.server.CreateBundle(api.CreateBundleInput{Name: name, Description: description, Refs: refs})
}

func (a *App) UpdateBundle(id string, name *string, description *string, refs *[]db.ContextRef) (db.Bundle, error) {
	if !a.ready {
		return db.Bundle{}, errNotReady
	}
	return a.server.UpdateBundle(id, api.UpdateBundleInput{Name: name, Description: description, Refs: refs})
}

func (a *App) DeleteBundle(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.DeleteBundle(id)
}

func (a *App) DuplicateBundle(id string) (db.Bundle, error) {
	if !a.ready {
		return db.Bundle{}, errNotReady
	}
	return a.server.DuplicateBundle(id)
}

func (a *App) ResolveBundle(id string) (api.ResolveResult, error) {
	if !a.ready {
		return api.ResolveResult{}, errNotReady
	}
	return a.server.ResolveBundle(id)
}

func (a *App) GetSettings() (db.AppSettings, error) {
	if !a.ready {
		return db.AppSettings{}, errNotReady
	}
	return a.server.GetSettings()
}

func (a *App) UpdateSettings(obsidianVaultRoot *string, taskOutputBaseURL *string) (db.AppSettings, error) {
	if !a.ready {
		return db.AppSettings{}, errNotReady
	}
	return a.server.UpdateSettings(api.UpdateSettingsInput{
		ObsidianVaultRoot: obsidianVaultRoot,
		TaskOutputBaseURL: taskOutputBaseURL,
	})
}

func (a *App) ListKeys() ([]db.AgentKey, error) {
	if !a.ready {
		return nil, errNotReady
	}
	return a.server.ListKeys()
}

func (a *App) IssueKey(name string) (api.IssueKeyResult, error) {
	if !a.ready {
		return api.IssueKeyResult{}, errNotReady
	}
	return a.server.IssueKey(name)
}

func (a *App) RevokeKey(id string) error {
	if !a.ready {
		return errNotReady
	}
	return a.server.RevokeKey(id)
}

// Quit はアプリを完全終了する（Settings 画面から呼ぶ）。
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}

func appDataDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".context-bundle-builder")
}
