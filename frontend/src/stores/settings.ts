import { defineStore } from 'pinia'
import { GetSettings, UpdateSettings } from '../../wailsjs/go/main/App'

export const useSettingsStore = defineStore('settings', {
  state: () => ({
    obsidianVaultRoot: '',
    taskOutputBaseURL: '',
    loading: false,
    saving: false,
    error: null as string | null,
  }),
  actions: {
    async load() {
      this.loading = true
      this.error = null
      try {
        const s = await GetSettings()
        this.obsidianVaultRoot = s.obsidian_vault_root
        this.taskOutputBaseURL = s.task_output_base_url
      } catch (e) {
        this.error = String(e)
      } finally {
        this.loading = false
      }
    },
    async save() {
      this.saving = true
      this.error = null
      try {
        const s = await UpdateSettings(this.obsidianVaultRoot, this.taskOutputBaseURL)
        this.obsidianVaultRoot = s.obsidian_vault_root
        this.taskOutputBaseURL = s.task_output_base_url
      } catch (e) {
        this.error = String(e)
      } finally {
        this.saving = false
      }
    },
  },
})
