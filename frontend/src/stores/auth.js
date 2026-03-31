import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authApi } from '@/api/auth'
import { storage } from '@/utils/storage'

export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref(storage.get('access_token'))
  const user = ref(null)
  const loading = ref(false)
  const error = ref(null)
  
  const isAuthenticated = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  
  async function login(email, password) {
    loading.value = true
    error.value = null
    
    try {
      const response = await authApi.login(email, password)
      accessToken.value = response.access_token
      storage.set('access_token', response.access_token)
      
      await fetchProfile()
      
      return response
    } catch (err) {
      error.value = err.message || 'Ошибка входа'
      throw err
    } finally {
      loading.value = false
    }
  }
  
  async function register(userData) {
    loading.value = true
    error.value = null
    
    try {
      const response = await authApi.register(userData)
      return response
    } catch (err) {
      error.value = err.message || 'Ошибка регистрации'
      throw err
    } finally {
      loading.value = false
    }
  }
  
  async function fetchProfile() {
    try {
      const response = await authApi.getProfile()
      user.value = response
      return response
    } catch (err) {
      // Игнорируем ошибки профиля
      return null
    }
  }
  
  async function refreshToken() {
    try {
      const response = await authApi.refresh()
      accessToken.value = response.access_token
      storage.set('access_token', response.access_token)
      return response
    } catch (err) {
      logout()
      throw err
    }
  }
  
  function logout() {
    accessToken.value = null
    user.value = null
    storage.remove('access_token')
  }
  
  return {
    accessToken,
    user,
    loading,
    error,
    isAuthenticated,
    isAdmin,
    login,
    register,
    fetchProfile,
    refreshToken,
    logout
  }
})
