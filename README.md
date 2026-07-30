# Context Bundle Builder

「AIに渡す情報を設計するツール」— comet-taskAI ロードマップ Product E。

ファイル・URL・Obsidianノート・前タスクの実行結果など複数ソースから、AIへ渡す「Bundle」を組み立てる。
実行前に「AIが実際に何を受け取るか」をプレビューできる。Bundleは保存・複製してテンプレートとして再利用できる。

詳細は [docs/spec.md](docs/spec.md) を参照。

## 現在のステータス: Phase 3（Wails + Vue3 UI）完了

- [x] Phase 0: プロジェクト立ち上げ
- [x] Phase 1: データモデル・Bundle CRUD API（APIキー認証・ブートストラップ認証）
- [x] Phase 2: ソースコネクタ（file/url/obsidian_note/task_output）・解決エンジン・プレビューAPI
- [x] Phase 3: Wails + Vue3 UI（Bundle一覧・編集・プレビューパネル・設定）
- [ ] Phase 4: 仕上げ・署名・配布・LP

## 使い方（デスクトップアプリ）

1. [Releases](../../releases) から自分のOS用のビルドをダウンロードして起動する
2. Bundle一覧から新規作成 → ファイル/URL/Obsidianノート/タスク出力のソースを追加
3. 「解決する」を押すと、AIが実際に受け取るテキストがプレビューに表示される
4. コピーしてAIに貼り付ける、または `.md` として書き出す
5. 外部ツール（AI Scheduler等）から使う場合は、Settings画面でAPIキーを発行してAPIエンドポイントに接続する

アプリはウインドウを閉じている間もAPIを起動したまま待ち受け続ける。完全に終了するにはSettings画面の「Quit」を使う。

## 使い方（開発・ヘッドレスサーバー）

```bash
go mod tidy   # 依存解決
make run      # :8422 でAPIサーバー起動（SQLite: context-bundle-builder.db、cmd/cbbserve）
make ui       # frontend/ の vite dev サーバー起動
make smoke    # bootstrap鍵発行 → Bundle作成 → 解決 → プレビュー の一連を確認する自己完結テスト
```

デスクトップアプリとしてビルドするには `wails build`（`wails.json` 参照）。

### APIキー認証

`AgentKey`が0件のときのみ `POST /api/v1/keys` を未認証で許可する（最初の1件を発行するため）。
1件発行された時点で以降は `Authorization: Bearer <key>` が必須になる。

### Bundleの作成・解決

```bash
curl -X POST localhost:8422/api/v1/bundles \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "PRレビュー用",
    "refs": [
      {"kind": "file", "ref": "/path/to/README.md"},
      {"kind": "url", "ref": "https://example.com/spec"}
    ]
  }'

curl -X POST localhost:8422/api/v1/bundles/{id}/resolve \
  -H "Authorization: Bearer $API_KEY"
```

## API

| メソッド | パス | 用途 |
|---|---|---|
| POST | `/api/v1/keys` | APIキー発行（ブートストラップ時のみ未認証） |
| GET | `/api/v1/keys` | 発行済みキー一覧 |
| DELETE | `/api/v1/keys/{id}` | キー失効 |
| POST | `/api/v1/bundles` | Bundle作成 |
| GET | `/api/v1/bundles` | 一覧 |
| GET | `/api/v1/bundles/{id}` | 単体取得 |
| PATCH | `/api/v1/bundles/{id}` | 更新（name/description/refs） |
| DELETE | `/api/v1/bundles/{id}` | 削除 |
| POST | `/api/v1/bundles/{id}/duplicate` | 複製 |
| POST | `/api/v1/bundles/{id}/resolve` | 全refsを解決してResolvedContextItem[]を返す |
| GET | `/api/v1/bundles/{id}/preview` | resolveと同内容（UIのプレビューパネル向けの名前） |
| GET | `/api/v1/settings` | ソースコネクタ設定（Obsidian vaultルート・task_outputアダプターURL） |
| PATCH | `/api/v1/settings` | 設定更新 |

## ディレクトリ構成

```
internal/db/        GORMモデル（Bundle/AgentKey/AppSettings/ResolveCacheEntry）・SQLite初期化
internal/api/        REST API（keys/bundles/resolve/settings）+ 認証ミドルウェア
internal/resolve/    ContextRef → ResolvedItem の解決エンジン（file/url/obsidian_note/task_output）
cmd/cbbserve/         ヘッドレスHTTPサーバー（Wailsを介さない開発・デバッグ用）
cmd/smoketest/       bootstrap→bundle作成→解決→プレビューの通しスモークテスト
frontend/            Wails + Vue3 + Vite + Pinia UI
main.go / app.go     Wailsデスクトップアプリのエントリポイント
docs/                設計ドキュメント
```
