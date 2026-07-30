# 技術選定

**決定日:** 2026-07-29
**ステータス:** 確定

---

## 決定

| レイヤー | 採用 | 理由 |
|---|---|---|
| Desktop基盤 | Wails v2 | A/B/K/C/Dの5製品すべてで実績あり。統合時の摩擦が最小 |
| Backend | Go 1.22+ | 同上 |
| Frontend | Vue 3 + Vite + Pinia | 同上（Nuxtは不採用。SSR不要） |
| Styling | UnoCSS | 同上 |
| DB | SQLite + GORM | 同上。個人〜小規模チームの規模で十分 |
| CI | GitHub Actions | Go 1.23 / macos-14 を**最初から**採用（harness-manager/execution-ledgerで踏んだ落とし穴を回避） |
| 配布 | `.app` / `.exe` シングルバイナリ + コード署名・公証 | 同じApple Developer証明書を`build-release.sh`で流用（環境変数で認証情報を渡す。直書きしない） |

## 却下した案

- **専用ベクトルDB/RAG基盤**: 本製品は「AIに渡す情報を人間が設計する」ツールであり、
  自動検索・埋め込み類似度検索は範囲外。ローカルファースト・シンプルさ優先の製品方針に反する
- **Nuxt SSR**: 他製品と同じ理由で不要

## 他製品からの流用ポイント・既知の落とし穴の事前回避

execution-ledgerのPhase 4で複数の落とし穴を踏んで学んだため、Context Bundle Builderでは
**Phase 0の時点から以下を先取りして対応する**（後回しにしない）：

1. **APIキー認証**: harness-manager / agent-config-manager / execution-ledger と同じ実装パターン
   （Bearer token → SHA-256ハッシュ照合 + ブートストラップ認証）をそのまま踏襲する
2. **`frontend/wailsjs` は最初からgitignoreしない**（コミット対象にする）。
   CIのfrontend-unitジョブがクリーンチェックアウトでtypecheckできず全滅する問題を、
   execution-ledgerではPhase 4で事後的に発見・修正した。今回はPhase 3のUI初期化時点で
   `.gitignore` に `frontend/wailsjs/` を書かない
3. **`.golangci.yml` はv2形式（`version: "2"` + `formatters`セクション）で最初から書く**。
   旧形式のまま放置すると `golangci-lint run` がconfig読み込みエラーで即死する
4. **CIのgo-testジョブは `./internal/... ./cmd/...` にスコープを限定する**。
   ルートパッケージに `//go:embed all:frontend/dist` があると、frontend/distが存在しない
   クリーンチェックアウトでは `go build ./...` / `go vet ./...` が失敗する
5. **release.yml**: harness-managerのv0.1.1で修正済みの構成（Go 1.23、notarytool用zip提出、
   windows拡張子なし対策、macos-14）をそのままコピーする
6. **`build-release.sh` / `wails.json` に個人のメールアドレス・実名を直書きしない**。
   `APPLE_ID` / `TEAM_ID` / `DEVELOPER_NAME` / `APP_PASSWORD` は全て環境変数で渡す。
   `wails.json` の `author` は `name: "comet"` のみとし、`email` フィールドは入れない
