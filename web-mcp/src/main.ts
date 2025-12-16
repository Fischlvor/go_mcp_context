import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import { pinia } from '@/stores'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'

// 导入基础样式（shadcn-vue + Tailwind）
import '@/assets/base.css'

const app = createApp(App)

app.use(pinia).use(router).use(ElementPlus)

app.mount('#app')
