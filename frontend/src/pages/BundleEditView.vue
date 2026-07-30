<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useBundlesStore } from '@/stores/bundles'
import { useI18n } from '@/i18n'
import { db } from '../../wailsjs/go/models'

const { t } = useI18n()
const store = useBundlesStore()
const route = useRoute()
const router = useRouter()

const id = computed(() => String(route.params.id))

const nameDraft = ref('')
const descDraft = ref('')
const refs = ref<db.ContextRef[]>([])
let saveTimer: ReturnType<typeof setTimeout> | null = null

const newKind = ref<'file' | 'url' | 'obsidian_note' | 'task_output'>('file')
const newRef = ref('')

const copied = ref(false)

async function loadBundle() {
  store.clearCurrent()
  await store.load(id.value)
  if (store.current) {
    nameDraft.value = store.current.name
    descDraft.value = store.current.description
    refs.value = [...(store.current.refs ?? [])]
  }
}

onMounted(loadBundle)
watch(id, loadBundle)

function scheduleSave() {
  if (saveTimer) clearTimeout(saveTimer)
  saveTimer = setTimeout(() => {
    store.save(id.value, { name: nameDraft.value, description: descDraft.value, refs: refs.value })
  }, 400)
}

function addRef() {
  if (!newRef.value.trim()) return
  refs.value.push({ kind: newKind.value, ref: newRef.value.trim() } as db.ContextRef)
  newRef.value = ''
  scheduleSave()
}

function removeRef(i: number) {
  refs.value.splice(i, 1)
  scheduleSave()
}

function moveRef(i: number, dir: -1 | 1) {
  const j = i + dir
  if (j < 0 || j >= refs.value.length) return
  const tmp = refs.value[i]
  refs.value[i] = refs.value[j]
  refs.value[j] = tmp
  scheduleSave()
}

async function resolve() {
  await store.save(id.value, { name: nameDraft.value, description: descDraft.value, refs: refs.value })
  await store.resolve(id.value)
}

async function copyPreview() {
  if (!store.resolved) return
  await navigator.clipboard.writeText(store.resolved.preview_text)
  copied.value = true
  setTimeout(() => (copied.value = false), 1500)
}

function exportMarkdown() {
  if (!store.resolved) return
  const blob = new Blob([store.resolved.preview_text], { type: 'text/markdown' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `${nameDraft.value || 'bundle'}.md`
  a.click()
  URL.revokeObjectURL(url)
}

const placeholderKey = computed(() => `edit.refs.placeholder.${newKind.value}`)
</script>

<template>
  <div class="flex h-full">
    <div class="flex-1 overflow-y-auto p-6 space-y-5 max-w-xl border-r border-border">
      <button @click="router.push('/bundles')" class="text-xs text-muted-foreground hover:text-foreground underline">
        &larr; {{ t('edit.back') }}
      </button>

      <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
        {{ t('error.prefix') }}{{ store.error }}
        <button @click="loadBundle" class="ml-2 underline">{{ t('error.retry') }}</button>
      </div>

      <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

      <template v-else-if="store.current">
        <div class="space-y-1">
          <label class="text-xs text-muted-foreground">{{ t('edit.name') }}</label>
          <input v-model="nameDraft" @input="scheduleSave" class="w-full text-sm border border-border rounded px-2 py-1.5 font-medium" />
        </div>
        <div class="space-y-1">
          <label class="text-xs text-muted-foreground">{{ t('edit.description') }}</label>
          <textarea v-model="descDraft" @input="scheduleSave" rows="2" class="w-full text-sm border border-border rounded px-2 py-1.5" />
        </div>

        <div class="space-y-2">
          <h3 class="text-xs font-semibold text-muted-foreground">{{ t('edit.refs.title') }}</h3>

          <div v-if="refs.length === 0" class="text-xs text-muted-foreground">{{ t('edit.refs.empty') }}</div>
          <ul v-else class="space-y-1">
            <li
              v-for="(r, i) in refs"
              :key="i"
              class="flex items-center gap-2 text-sm border border-border rounded px-2 py-1.5"
            >
              <span class="text-xs shrink-0 px-1.5 py-0.5 rounded bg-gray-100 text-gray-600">{{ t('edit.refs.kind.' + r.kind) }}</span>
              <span class="flex-1 truncate font-mono text-xs">{{ r.ref }}</span>
              <button @click="moveRef(i, -1)" :disabled="i === 0" class="text-muted-foreground disabled:opacity-20" :title="t('edit.refs.moveUp')">&uarr;</button>
              <button @click="moveRef(i, 1)" :disabled="i === refs.length - 1" class="text-muted-foreground disabled:opacity-20" :title="t('edit.refs.moveDown')">&darr;</button>
              <button @click="removeRef(i)" class="text-red-500 hover:text-red-700" :title="t('edit.refs.remove')">&times;</button>
            </li>
          </ul>

          <div class="flex gap-2 pt-1">
            <select v-model="newKind" class="text-sm border border-border rounded px-2 py-1.5">
              <option value="file">{{ t('edit.refs.kind.file') }}</option>
              <option value="url">{{ t('edit.refs.kind.url') }}</option>
              <option value="obsidian_note">{{ t('edit.refs.kind.obsidian_note') }}</option>
              <option value="task_output">{{ t('edit.refs.kind.task_output') }}</option>
            </select>
            <input
              v-model="newRef"
              @keyup.enter="addRef"
              :placeholder="t(placeholderKey)"
              class="flex-1 text-sm border border-border rounded px-2 py-1.5 font-mono"
            />
            <button @click="addRef" class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white shrink-0">{{ t('edit.refs.add') }}</button>
          </div>
        </div>
      </template>
    </div>

    <div class="flex-1 overflow-y-auto p-6 space-y-3">
      <div class="flex items-center justify-between">
        <h3 class="text-sm font-semibold">{{ t('preview.title') }}</h3>
        <button
          @click="resolve"
          :disabled="store.resolving"
          class="text-xs px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
        >{{ t('preview.resolve') }}</button>
      </div>

      <div v-if="store.resolving" class="text-sm text-muted-foreground">{{ t('loading') }}</div>

      <template v-else-if="store.resolved">
        <div class="flex items-center justify-between text-xs text-muted-foreground">
          <span>{{ t('preview.stats', { chars: store.resolved.char_count, tokens: store.resolved.estimated_tokens }) }}</span>
          <div class="flex gap-2">
            <button @click="copyPreview" class="px-2 py-1 rounded border border-border hover:bg-gray-50">
              {{ copied ? t('preview.copied') : t('preview.copy') }}
            </button>
            <button @click="exportMarkdown" class="px-2 py-1 rounded border border-border hover:bg-gray-50">
              {{ t('preview.export') }}
            </button>
          </div>
        </div>

        <div v-if="store.resolved.errors?.length" class="text-xs border rounded px-3 py-2 border-amber-300 text-amber-700 bg-amber-50">
          {{ t('preview.errors.title', { n: store.resolved.errors.length }) }}
          <ul class="mt-1 list-disc list-inside">
            <li v-for="(e, i) in store.resolved.errors" :key="i" class="font-mono truncate">{{ e.ref.ref }}: {{ e.message }}</li>
          </ul>
        </div>

        <pre class="text-xs font-mono whitespace-pre-wrap border border-border rounded p-3 bg-gray-50 max-h-[70vh] overflow-y-auto">{{ store.resolved.preview_text || '—' }}</pre>
      </template>

      <div v-else class="text-sm text-muted-foreground">{{ t('preview.empty') }}</div>
    </div>
  </div>
</template>
