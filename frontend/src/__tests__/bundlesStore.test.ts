import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useBundlesStore } from '@/stores/bundles'
import { ListBundles, ResolveBundle } from '../../wailsjs/go/main/App'

describe('bundles store error handling', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    vi.mocked(ListBundles).mockReset()
    vi.mocked(ResolveBundle).mockReset()
  })

  it('captures a failed list() as store.error and clears loading', async () => {
    vi.mocked(ListBundles).mockRejectedValueOnce(new Error('network down'))
    const store = useBundlesStore()

    await store.list()

    expect(store.loading).toBe(false)
    expect(store.error).toContain('network down')
  })

  it('clears the previous error on a successful retry', async () => {
    vi.mocked(ListBundles).mockRejectedValueOnce(new Error('network down'))
    const store = useBundlesStore()
    await store.list()
    expect(store.error).not.toBeNull()

    vi.mocked(ListBundles).mockResolvedValueOnce([])
    await store.list()

    expect(store.error).toBeNull()
  })

  it('captures a failed resolve() without throwing', async () => {
    vi.mocked(ResolveBundle).mockRejectedValueOnce(new Error('resolve failed'))
    const store = useBundlesStore()

    await store.resolve('bundle-1')

    expect(store.resolving).toBe(false)
    expect(store.error).toContain('resolve failed')
  })
})
