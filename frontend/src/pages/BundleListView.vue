<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBundlesStore } from '@/stores/bundles'
import { useI18n } from '@/i18n'

const { t } = useI18n()
const store = useBundlesStore()
const router = useRouter()

const newName = ref('')
const creating = ref(false)

onMounted(() => store.list())

async function createBundle() {
  if (!newName.value.trim()) return
  creating.value = true
  const b = await store.create(newName.value.trim())
  creating.value = false
  newName.value = ''
  if (b) router.push(`/bundles/${b.id}`)
}

async function remove(id: string) {
  if (!confirm(t('bundles.card.delete.confirm'))) return
  await store.remove(id)
}

function fmt(v: any): string {
  const d = new Date(v)
  return isNaN(d.getTime()) ? String(v) : d.toLocaleString()
}
</script>

<template>
  <div class="p-6 space-y-6 overflow-y-auto h-full">
    <div class="flex items-center justify-between">
      <h2 class="text-sm font-semibold">{{ t('bundles.title') }}</h2>
    </div>

    <div v-if="store.error" class="text-sm border rounded px-3 py-2 border-red-300 text-red-600">
      {{ t('error.prefix') }}{{ store.error }}
      <button @click="store.list" class="ml-2 underline">{{ t('error.retry') }}</button>
    </div>

    <div class="flex gap-2 max-w-md">
      <input
        v-model="newName"
        @keyup.enter="createBundle"
        :placeholder="t('bundles.new.name')"
        class="flex-1 text-sm border border-border rounded px-2 py-1.5"
      />
      <button
        @click="createBundle"
        :disabled="creating || !newName.trim()"
        class="text-sm px-3 py-1.5 rounded bg-gray-900 text-white disabled:opacity-40"
      >{{ t('bundles.new.create') }}</button>
    </div>

    <div v-if="store.loading" class="text-sm text-muted-foreground">{{ t('loading') }}</div>
    <div v-else-if="store.bundles.length === 0" class="text-sm text-muted-foreground">{{ t('bundles.empty') }}</div>

    <div v-else class="grid gap-3" style="grid-template-columns: repeat(auto-fill, minmax(260px, 1fr))">
      <div
        v-for="b in store.bundles"
        :key="b.id"
        class="border border-border rounded-lg p-4 flex flex-col gap-2"
      >
        <div>
          <div class="font-medium text-sm truncate">{{ b.name }}</div>
          <div v-if="b.description" class="text-xs text-muted-foreground truncate">{{ b.description }}</div>
        </div>
        <div class="text-xs text-muted-foreground">
          {{ t('bundles.card.refs', { n: b.refs?.length ?? 0 }) }} · {{ fmt(b.updated_at) }}
        </div>
        <div class="flex gap-2 mt-auto pt-2">
          <button @click="router.push(`/bundles/${b.id}`)" class="text-xs px-2 py-1 rounded border border-border hover:bg-gray-50">
            {{ t('bundles.card.open') }}
          </button>
          <button @click="store.duplicate(b.id)" class="text-xs px-2 py-1 rounded border border-border hover:bg-gray-50">
            {{ t('bundles.card.duplicate') }}
          </button>
          <button @click="remove(b.id)" class="text-xs px-2 py-1 rounded border border-border text-red-600 hover:bg-red-50 ml-auto">
            {{ t('bundles.card.delete') }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
