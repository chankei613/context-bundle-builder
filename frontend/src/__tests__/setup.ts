import { vi } from 'vitest'

// Node 25+ の実験的 localStorage がhappy-domのwindow.localStorageと衝突し
// `--localstorage-file` 未指定だと getItem 等がthrowする既知の非互換。
// メモリ上の簡易実装に差し替えて回避する（execution-ledgerで判明した対処法）。
class MemoryStorage implements Storage {
  private store = new Map<string, string>()
  get length() { return this.store.size }
  clear() { this.store.clear() }
  getItem(key: string) { return this.store.has(key) ? this.store.get(key)! : null }
  key(index: number) { return Array.from(this.store.keys())[index] ?? null }
  removeItem(key: string) { this.store.delete(key) }
  setItem(key: string, value: string) { this.store.set(key, String(value)) }
}
const memoryStorage = new MemoryStorage()
Object.defineProperty(globalThis, 'localStorage', { value: memoryStorage, configurable: true })
Object.defineProperty(window, 'localStorage', { value: memoryStorage, configurable: true })

// Wails ランタイムのモック — テスト環境では no-op
vi.mock('../../wailsjs/runtime/runtime', () => ({
  EventsOn: vi.fn(),
  EventsOff: vi.fn(),
  EventsEmit: vi.fn(),
  BrowserOpenURL: vi.fn(),
}))

vi.mock('../../wailsjs/go/main/App', () => ({
  GetAppVersion: vi.fn().mockResolvedValue('0.1.0'),
  GetAPIURL: vi.fn().mockResolvedValue('http://127.0.0.1:8422'),
  Quit: vi.fn().mockResolvedValue(undefined),
  ListBundles: vi.fn().mockResolvedValue([]),
  GetBundle: vi.fn().mockResolvedValue({ id: '1', name: '', description: '', refs: [] }),
  CreateBundle: vi.fn().mockResolvedValue({ id: '1', name: 'test', description: '', refs: [] }),
  UpdateBundle: vi.fn().mockResolvedValue({ id: '1', name: 'test', description: '', refs: [] }),
  DeleteBundle: vi.fn().mockResolvedValue(undefined),
  DuplicateBundle: vi.fn().mockResolvedValue({ id: '2', name: 'test (copy)', description: '', refs: [] }),
  ResolveBundle: vi.fn().mockResolvedValue({ bundle_id: '1', items: [], errors: [], preview_text: '', char_count: 0, estimated_tokens: 0 }),
  GetSettings: vi.fn().mockResolvedValue({ obsidian_vault_root: '', task_output_base_url: '' }),
  UpdateSettings: vi.fn().mockResolvedValue({ obsidian_vault_root: '', task_output_base_url: '' }),
  ListKeys: vi.fn().mockResolvedValue([]),
  IssueKey: vi.fn().mockResolvedValue({ id: '1', name: 'test', api_key: 'test-key' }),
  RevokeKey: vi.fn().mockResolvedValue(undefined),
}))
