import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createRouter, createWebHashHistory } from 'vue-router'
import 'virtual:uno.css'
import '@unocss/reset/tailwind.css'
import './assets/globals.css'
import App from './App.vue'
import BundleListView from './pages/BundleListView.vue'
import BundleEditView from './pages/BundleEditView.vue'
import HelpView from './pages/HelpView.vue'
import SettingsView from './pages/SettingsView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    { path: '/', redirect: '/bundles' },
    { path: '/bundles', component: BundleListView },
    { path: '/bundles/:id', component: BundleEditView },
    { path: '/help', component: HelpView },
    { path: '/settings', component: SettingsView },
  ],
})

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
