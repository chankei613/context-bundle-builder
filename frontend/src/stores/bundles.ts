import { defineStore } from 'pinia'
import {
  ListBundles,
  GetBundle,
  CreateBundle,
  UpdateBundle,
  DeleteBundle,
  DuplicateBundle,
  ResolveBundle,
} from '../../wailsjs/go/main/App'
import { db, api } from '../../wailsjs/go/models'

export const useBundlesStore = defineStore('bundles', {
  state: () => ({
    bundles: [] as db.Bundle[],
    current: null as db.Bundle | null,
    resolved: null as api.ResolveResult | null,
    loading: false,
    resolving: false,
    error: null as string | null,
  }),
  actions: {
    async list() {
      this.loading = true
      this.error = null
      try {
        this.bundles = (await ListBundles()) ?? []
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async load(id: string) {
      this.loading = true
      this.error = null
      this.resolved = null
      try {
        this.current = await GetBundle(id)
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async create(name: string): Promise<db.Bundle | null> {
      this.error = null
      try {
        const b = await CreateBundle(name, '', [])
        await this.list()
        return b
      } catch (e) {
        this.error = String(e)
        return null
      }
    },
    async save(id: string, patch: { name?: string; description?: string; refs?: db.ContextRef[] }) {
      this.error = null
      try {
        this.current = await UpdateBundle(
          id,
          patch.name ?? null,
          patch.description ?? null,
          patch.refs ?? null,
        )
      } catch (e) {
        this.error = String(e)
      }
    },
    async remove(id: string) {
      this.error = null
      try {
        await DeleteBundle(id)
        await this.list()
      } catch (e) {
        this.error = String(e)
      }
    },
    async duplicate(id: string) {
      this.error = null
      try {
        await DuplicateBundle(id)
        await this.list()
      } catch (e) {
        this.error = String(e)
      }
    },
    async resolve(id: string) {
      this.resolving = true
      this.error = null
      try {
        this.resolved = await ResolveBundle(id)
      } catch (e) {
        this.error = String(e)
      } finally {
        this.resolving = false
      }
    },
    clearCurrent() {
      this.current = null
      this.resolved = null
    },
  },
})
