import { ref } from 'vue'

export type Locale = 'en' | 'ja'

// localStorage に保存して再起動後も維持する
const saved = localStorage.getItem('locale') as Locale | null
const locale = ref<Locale>(saved === 'en' || saved === 'ja' ? saved : 'ja')

const messages: Record<Locale, Record<string, string>> = {
  en: {
    'app.subtitle': 'Context Bundle Builder',
    'lang.toggle': 'JA',
    'nav.bundles': 'Bundles',
    'nav.settings': 'Settings',

    'error.prefix': 'Error: ',
    'error.retry': 'Retry',
    'loading': 'Loading…',

    'bundles.title': 'Bundles',
    'bundles.empty': 'No bundles yet. Create one to start collecting context.',
    'bundles.new': 'New bundle',
    'bundles.new.name': 'Bundle name',
    'bundles.new.create': 'Create',
    'bundles.card.refs': '{n} refs',
    'bundles.card.open': 'Open',
    'bundles.card.duplicate': 'Duplicate',
    'bundles.card.delete': 'Delete',
    'bundles.card.delete.confirm': 'Delete this bundle? This cannot be undone.',

    'edit.back': 'Back to bundles',
    'edit.name': 'Name',
    'edit.description': 'Description',
    'edit.refs.title': 'Sources',
    'edit.refs.empty': 'No sources yet. Add a file, URL, Obsidian note, or task output below.',
    'edit.refs.add': 'Add source',
    'edit.refs.kind.file': 'File',
    'edit.refs.kind.url': 'URL',
    'edit.refs.kind.obsidian_note': 'Obsidian note',
    'edit.refs.kind.task_output': 'Task output',
    'edit.refs.placeholder.file': '/path/to/file.md',
    'edit.refs.placeholder.url': 'https://example.com/doc',
    'edit.refs.placeholder.obsidian_note': 'projects/foo/notes.md (relative to vault root)',
    'edit.refs.placeholder.task_output': 'task#123.output_key',
    'edit.refs.remove': 'Remove',
    'edit.refs.moveUp': 'Move up',
    'edit.refs.moveDown': 'Move down',

    'preview.title': 'Preview — what the AI receives',
    'preview.resolve': 'Resolve',
    'preview.empty': 'Add sources and press Resolve to see what will be sent.',
    'preview.stats': '{chars} chars · ~{tokens} tokens',
    'preview.copy': 'Copy to clipboard',
    'preview.copied': 'Copied',
    'preview.export': 'Export as .md',
    'preview.errors.title': '{n} source(s) failed to resolve',

    'settings.title': 'Settings',
    'settings.api.title': 'API endpoint',
    'settings.api.desc': 'External tools (e.g. AI Scheduler) can read resolved bundles from this address.',
    'settings.sources.title': 'Source connectors',
    'settings.sources.obsidianRoot': 'Obsidian vault root',
    'settings.sources.obsidianRoot.desc': 'Absolute path to the vault. Notes are resolved relative to this.',
    'settings.sources.taskOutputURL': 'Task output adapter URL',
    'settings.sources.taskOutputURL.desc': 'Base URL of a small adapter exposing GET {base}/resolve?ref=... Leave empty if unused.',
    'settings.save': 'Save',
    'settings.saved': 'Saved',
    'settings.keys.title': 'API keys',
    'settings.keys.name': 'Key name',
    'settings.keys.issue': 'Issue key',
    'settings.keys.issued': 'Key issued — copy it now, it will not be shown again',
    'settings.keys.copy': 'Copy',
    'settings.keys.revoke': 'Revoke',
    'settings.keys.revoked': 'Revoked',
    'settings.keys.empty': 'No keys issued yet. Issue one to allow external tools to use the API.',
    'settings.version': 'Version',
    'settings.quit': 'Quit',
    'settings.quit.confirm': 'Quit the app? External access via the API will stop until you reopen it.',
  },
  ja: {
    'app.subtitle': 'Context Bundle Builder',
    'lang.toggle': 'EN',
    'nav.bundles': 'Bundle',
    'nav.settings': '設定',

    'error.prefix': 'エラー: ',
    'error.retry': '再試行',
    'loading': '読み込み中…',

    'bundles.title': 'Bundle一覧',
    'bundles.empty': 'まだBundleがありません。作成してソースを集め始めましょう。',
    'bundles.new': '新規作成',
    'bundles.new.name': 'Bundle名',
    'bundles.new.create': '作成',
    'bundles.card.refs': '{n}件のソース',
    'bundles.card.open': '開く',
    'bundles.card.duplicate': '複製',
    'bundles.card.delete': '削除',
    'bundles.card.delete.confirm': 'このBundleを削除しますか？元に戻せません。',

    'edit.back': 'Bundle一覧へ戻る',
    'edit.name': '名前',
    'edit.description': '説明',
    'edit.refs.title': 'ソース',
    'edit.refs.empty': 'まだソースがありません。下からファイル・URL・Obsidianノート・タスク出力を追加してください。',
    'edit.refs.add': 'ソースを追加',
    'edit.refs.kind.file': 'ファイル',
    'edit.refs.kind.url': 'URL',
    'edit.refs.kind.obsidian_note': 'Obsidianノート',
    'edit.refs.kind.task_output': 'タスク出力',
    'edit.refs.placeholder.file': '/path/to/file.md',
    'edit.refs.placeholder.url': 'https://example.com/doc',
    'edit.refs.placeholder.obsidian_note': 'projects/foo/notes.md（vaultルートからの相対パス）',
    'edit.refs.placeholder.task_output': 'task#123.output_key',
    'edit.refs.remove': '削除',
    'edit.refs.moveUp': '上へ',
    'edit.refs.moveDown': '下へ',

    'preview.title': 'プレビュー — AIが実際に受け取る内容',
    'preview.resolve': '解決する',
    'preview.empty': 'ソースを追加して「解決する」を押すと、送信される内容が表示されます。',
    'preview.stats': '{chars}文字 · 約{tokens}トークン',
    'preview.copy': 'クリップボードにコピー',
    'preview.copied': 'コピーしました',
    'preview.export': '.mdとして書き出し',
    'preview.errors.title': '{n}件のソースが解決に失敗しました',

    'settings.title': '設定',
    'settings.api.title': 'APIエンドポイント',
    'settings.api.desc': '外部ツール（AI Scheduler等）はこのアドレスから解決済みBundleを取得できます。',
    'settings.sources.title': 'ソースコネクタ',
    'settings.sources.obsidianRoot': 'Obsidian vaultルート',
    'settings.sources.obsidianRoot.desc': 'vaultの絶対パス。ノートはここからの相対パスで解決されます。',
    'settings.sources.taskOutputURL': 'タスク出力アダプターURL',
    'settings.sources.taskOutputURL.desc': 'GET {base}/resolve?ref=... を提供する小さなアダプターのベースURL。使わない場合は空のままで構いません。',
    'settings.save': '保存',
    'settings.saved': '保存しました',
    'settings.keys.title': 'APIキー',
    'settings.keys.name': 'キー名',
    'settings.keys.issue': 'キーを発行',
    'settings.keys.issued': 'キーを発行しました — この場では二度と表示されないので今すぐコピーしてください',
    'settings.keys.copy': 'コピー',
    'settings.keys.revoke': '失効',
    'settings.keys.revoked': '失効済み',
    'settings.keys.empty': 'まだキーがありません。外部ツールにAPIを使わせる場合は発行してください。',
    'settings.version': 'バージョン',
    'settings.quit': '終了',
    'settings.quit.confirm': 'アプリを終了しますか？再度開くまでAPI経由の外部アクセスは停止します。',
  },
}

export function useI18n() {
  function t(key: string, params?: Record<string, string | number>): string {
    let msg = messages[locale.value][key] ?? messages.en[key] ?? key
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        msg = msg.replace(`{${k}}`, String(v))
      }
    }
    return msg
  }

  function toggleLocale() {
    locale.value = locale.value === 'en' ? 'ja' : 'en'
    localStorage.setItem('locale', locale.value)
  }

  return { t, locale, toggleLocale }
}
