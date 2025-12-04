import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { pinia } from '@/stores'

// 导入基础样式（shadcn-vue + Tailwind）
import '@/assets/base.css'

const app = createApp(App)

app.use(pinia).use(router)

app.mount('#app')
