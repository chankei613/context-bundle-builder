# 開発計画

**予測期間:** 2〜3週間相当（ロードマップ見積もり）。過去実績（A〜D）に基づき短縮を狙う。

| Phase | 内容 |
|---|---|
| Phase 0 | プロジェクト立ち上げ（Go初期化・docs・GitHub repo） |
| Phase 1 | データモデル・Bundle CRUD API（APIキー認証・ブートストラップ認証） |
| Phase 2 | ソースコネクタ（file/url/obsidian_note/task_output）・解決エンジン・プレビューAPI |
| Phase 3 | Wails + Vue3 UI（Bundle一覧・編集・プレビューパネル・コピー導線） |
| Phase 4 | 仕上げ・署名・配布・LP |

## 優先順位の根拠

本製品の価値は「解決エンジンが正しく動く」ことに尽きる（UIはその上に乗る薄いレイヤー）。
execution-ledgerと同じ判断で、Phase 1-2（CRUD + 解決エンジン）をUIより先に固め、
`curl`ベースの手動テストで4種類のコネクタが正しく動くことを検証してからUIに進む。

Phase 4はexecution-ledgerで判明した既知の落とし穴（wailsjs未コミット・golangci-lint v2形式・
CI embedスコープ・個人情報直書き）を tech-stack.md に明記済みなので、**Phase 0〜3の時点から
先取りして対応し、Phase 4での手戻りを無くす**のが今回の狙い。
