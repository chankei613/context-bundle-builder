# Context Bundle Builder — 仕様書

> 作成: 2026-07-29
> ステータス: 設計フェーズ

---

## 1. 製品概要

**「AIに渡す情報を設計するツール」** — ファイル・URL・Obsidianノート・前タスクの実行結果など
複数ソースから、AIへ渡す「Context Bundle」を組み立て、実行前に「AIが実際に何を受け取るか」を
プレビューできるデスクトップアプリ。

### 解決する問題

- 毎回同じファイルを手でコピー&ペーストしてAIに渡している
- 何をAIに渡したか（渡していないか）を振り返れない
- 良い指示の構成（どのファイル＋どのURL＋どのメモを渡すと成果物の質が上がるか）を
  再利用・チームで共有する手段がない
- AI Scheduler（Product B）でスケジュール実行する際、毎回のContext組み立てが手作業になる

### ソリューション

`ContextRef`（file/url/obsidian_note/task_output）を集めた「Bundle」を作成・保存し、
`resolve` した結果を1つのテキストとしてプレビュー・コピーできるようにする。
Bundleは名前を付けて保存でき、複製・テンプレート化して他のプロジェクトにも使い回せる。

---

## 2. コアコンセプト

### データモデルはcomet-taskAIの既存スキーマをそのまま使う

`cometinc/comet-taskAI/schema/types.ts` に既に定義されている型を正規ソースとする。
新規に型を発明しない。

```typescript
// 既存（comet-taskAI/schema/types.ts より）
interface ContextRef {
  kind: "file" | "url" | "obsidian_note" | "task_output"
  ref: string // ファイルパス / URL / ノートID / task#id.output_key
}

interface ResolvedContextItem {
  ref: ContextRef
  content: string       // 事前解決済みの本文
  resolved_at: string
}
```

本製品が新規に追加するのは「Bundle」というContextRefの集合と、それを解決・プレビューする機構のみ。

```typescript
interface Bundle {
  id: string
  name: string
  description: string
  refs: ContextRef[]       // 追加順を保持（並び替え可能）
  created_at: string
  updated_at: string
}

interface ResolvedBundle {
  bundle_id: string
  items: ResolvedContextItem[]   // refs と同じ順序で解決
  errors: { ref: ContextRef; message: string }[]  // 解決に失敗したものは落とさず記録する
  preview_text: string           // items[].content を連結した最終テキスト
  char_count: number
  estimated_tokens: number       // 簡易換算（4文字≒1トークン）。外部API呼び出し不要
  resolved_at: string
}
```

**失敗を握りつぶさない方針**: 1つのrefが解決できなくても（ファイルが消えていた、URLが404等）
Bundle全体の解決は止めない。`errors[]` に記録し、プレビューにはその旨を明示する。

---

## 3. 機能一覧

### Phase 1 (Bundle CRUD)

| 機能 | 説明 |
|------|------|
| Bundle作成/一覧/取得/更新/削除 | 名前 + ContextRef[] の基本CRUD |
| 複製 | 既存Bundleを複製して別名で保存（テンプレート化の起点） |
| APIキー管理 | ブートストラップ発行（0件時のみ未認証）・失効 |

### Phase 2 (解決エンジン)

| 機能 | 説明 |
|------|------|
| file コネクタ | ローカルファイルパスを読む |
| url コネクタ | HTTP GET、タイムアウト・サイズ上限、簡易HTML→テキスト抽出 |
| obsidian_note コネクタ | 設定したvaultルート配下の.mdファイルをパスで読む |
| task_output コネクタ | `task#id.output_key` 形式。comet-taskAI/execution-ledgerが未接続でも壊れない（空+警告） |
| Bundle全体の解決 | 全refsを解決してResolvedContextItem[]を返す。1件失敗しても続行 |
| プレビュー生成 | 連結テキスト・文字数・概算トークン数 |
| 解決結果キャッシュ | 同一refをTTL内なら再取得しない |

### Phase 3 (UI)

| 機能 | 説明 |
|------|------|
| Bundle一覧ビュー | 作成・複製・削除 |
| Bundle編集ビュー | ContextRef追加（ファイル選択/URL入力/ノート選択/task_output入力）、ドラッグで並び替え |
| プレビューパネル | 解決結果をソースごとに折りたたみ表示 + 全体連結テキスト + 文字数/概算トークン数 |
| コピー/エクスポート | プレビューテキストをクリップボードコピー、`.md`書き出し |
| APIキー管理UI | 発行・一覧・失効 |

### Phase 4 (拡張候補・本リリースの対象外)

| 機能 | 説明 |
|------|------|
| AI Scheduler連携 | スケジュール実行前にBundleを自動解決して注入する（Product B側の対応が必要） |
| Bundleのエクスポート/インポート | チーム間共有（JSON） |

---

## 4. UX フロー

```
起動
 └── Bundle一覧ビュー（既定画面）
      ├── 新規作成 → Bundle編集ビュー
      │     ├── ContextRef追加（file/url/obsidian_note/task_output）
      │     ├── 並び替え
      │     └── プレビューパネルで即座に解決結果を確認
      ├── 複製 → 別名で保存
      └── コピー → プレビューテキストをクリップボードへ

APIキー管理
 └── SettingsView → キー発行（ブートストラップ後は認証必須）・失効
```

---

## 5. データストア

SQLite（ローカル、`~/.context-bundle-builder/bundles.db`）

```sql
bundles (
  id, name, description, refs JSON, created_at, updated_at
)
agent_keys (id, name, api_key_hash, created_at, revoked_at)
resolve_cache (
  ref_kind, ref_value, content, resolved_at, expires_at
)
```

`refs` はJSONカラム（GORM `serializer:json`、harness-manager/execution-ledgerと同じパターン）。
`resolve_cache` はrefのkind+valueで一意化し、TTL経過後は再取得する。
