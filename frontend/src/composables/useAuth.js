import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

export function useAuth() {
  const authStore = useAuthStore()
  
  const isAuthenticated = computed(() => authStore.isAuthenticated)
  const isAdmin = computed(() => authStore.isAdmin)
  const user = computed(() => authStore.user)
  const loading = computed(() => authStore.loading)
  const error = computed(() => authStore.error)
  
  async function login(email, password) {
    return authStore.login(email, password)
  }
  
  async function register(userData) {
    return authStore.register(userData)
  }
  
  function logout() {
    authStore.logout()
  }
  
  async function loadProfile() {
    return authStore.fetchProfile()
  }
  
  return {
    isAuthenticated,
    isAdmin,
    user,
    loading,
    error,
    login,
    register,
    logout,
    loadProfile
  }
}
