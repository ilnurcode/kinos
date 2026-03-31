import 'bootstrap/dist/css/bootstrap.min.css'
import { createPinia } from 'pinia'
import { createApp } from 'vue'
import App from './App.vue'
import './assets/css/main.scss'
import router from './router'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

// После загрузки приложения проверяем есть ли токен и загружаем профиль
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore()
const token = localStorage.getItem('access_token')
if (token) {
	authStore.fetchProfile().catch(() => {
		console.log('Failed to fetch profile')
	})
}

app.mount('#app')
