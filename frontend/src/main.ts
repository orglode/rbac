import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './app/router'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/message-box/style/css'
import './styles/index.css'

const app = createApp(App)
app.use(createPinia())
app.use(router)
app.mount('#app')
