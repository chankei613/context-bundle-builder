<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSettingsStore } from '@/stores/settings'
import { useKeysStore } from '@/stores/keys'
import { useI18n } from '@/i18n'
import { Quit } from '../../wailsjs/go/main/App'

const { t } = useI18n()
const settings = useSettingsStore()
const keys = useKeysStore()
const newKeyName = ref('')
const copied = ref(false)
const saved = ref(false)

onMounted(() => {
  settings.load()
  keys.load()
})

async function save() {
  await settings.save()
  saved.value = true
  setTimeout(() => (saved.value = false), 1500)
}

async function issue() {
  if (!newKeyName.value.trim()) return
  await keys.issue(newKeyName.value.trim())
  newKeyName.value = ''
}

function copyKey() {
  if (!keys.lastIssuedKey) return
  navigator.clipboard.writeText(keys.lastIssuedKey.apiKey)
  copied.value = true
  setTimeout(() => (copied.value = false), 1500)
}

function copyURL() {
  navigator.clipboard.writeText(keys.apiURL)
}

function confirmQuit() {
  if (confirm(t('settings.quit.confirm'))) Quit()
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 max-w-xl overflow-y-auto h-full">
    <h2 class="text-sm font-semibold">{{ t('settings.title') }}</h2>

    <div v-if="settings.error || keys.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ settings.error || keys.error }}
    </div>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('settings.api.title') }}</h3>
      <p class="text-xs text-muted-foreground">{{ t('settings.api.desc') }}</p>
      <div class="flex items-center gap-2">
        <code class="flex-1 text-xs bg-gray-50 border border-border rounded px-2 py-1.5 truncate">{{ keys.apiURL }}</code>
        <button @click="copyURL" class="text-xs px-2 py-1.5 border border-border rounded hover:bg-gray-50">{{ t('settings.keys.copy') }}</button>
      </div>
    </section>

    <section class="space-y-3">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('settings.sources.title') }}</h3>

      <div class="space-y-1">
        <label class="text-xs">{{ t('settings.sources.obsidianRoot') }}</label>
        <p class="text-xs text-muted-foreground">{{ t('settings.sources.obsidianRoot.desc') }}</p>
        <input v-model="settings.obsidianVaultRoot" class="w-full text-sm border border-border rounded px-2 py-1.5 font-mono" placeholder="/Users/you/ObsidianVault" />
      </div>

      <div class="space-y-1">
        <label class="text-xs">{{ t('settings.sources.taskOutputURL') }}</label>
        <p class="text-xs text-muted-foreground">{{ t('settings.sources.taskOutputURL.desc') }}</p>
        <input v-model="settings.taskOutputBaseURL" class="w-full text-sm border border-border rounded px-2 py-1.5 font-mono" placeholder="http://localhost:9000" />
      </div>

      <button @click="save" :disabled="settings.saving" class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40">
        {{ saved ? t('settings.saved') : t('settings.save') }}
      </button>
    </section>

    <section class="space-y-2">
      <h3 class="text-xs font-semibold text-muted-foreground">{{ t('settings.keys.title') }}</h3>

      <div v-if="keys.lastIssuedKey" class="border rounded p-3 space-y-1" style="border-color: #fab219; background: #fffaf0">
        <div class="text-xs font-medium">{{ t('settings.keys.issued') }}</div>
        <div class="flex items-center gap-2">
          <code class="flex-1 text-xs bg-white border border-border rounded px-2 py-1 truncate">{{ keys.lastIssuedKey.apiKey }}</code>
          <button @click="copyKey" class="text-xs px-2 py-1 border border-border rounded hover:bg-gray-50">
            {{ copied ? '✓' : t('settings.keys.copy') }}
          </button>
        </div>
      </div>

      <div class="flex gap-2">
        <input
          v-model="newKeyName"
          @keyup.enter="issue"
          :placeholder="t('settings.keys.name')"
          class="flex-1 text-sm border border-border rounded px-2 py-1.5"
        />
        <button @click="issue" class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white">{{ t('settings.keys.issue') }}</button>
      </div>

      <div v-if="keys.loading" class="text-xs text-muted-foreground">{{ t('loading') }}</div>
      <div v-else-if="keys.keys.length === 0" class="text-xs text-muted-foreground">{{ t('settings.keys.empty') }}</div>
      <ul v-else class="space-y-1">
        <li
          v-for="k in keys.keys"
          :key="k.id"
          class="flex items-center justify-between text-sm border border-border rounded px-3 py-2"
        >
          <div>
            <div class="font-medium">{{ k.name }}</div>
            <div class="text-xs text-muted-foreground">{{ fmt(k.created_at) }}</div>
          </div>
          <span v-if="k.revoked_at" class="text-xs text-muted-foreground">{{ t('settings.keys.revoked') }}</span>
          <button v-else @click="keys.revoke(k.id)" class="text-xs text-red-600 hover:underline">{{ t('settings.keys.revoke') }}</button>
        </li>
      </ul>
    </section>

    <section class="pt-4 border-t border-border flex items-center justify-between">
      <span class="text-xs text-muted-foreground">{{ t('settings.version') }}: {{ keys.appVersion }}</span>
      <button @click="confirmQuit" class="text-xs px-3 py-1.5 border border-border rounded hover:bg-gray-50">{{ t('settings.quit') }}</button>
    </section>
  </div>
</template>
