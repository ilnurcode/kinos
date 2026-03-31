import { defineStore } from 'pinia'
import { ref } from 'vue'
import { usersApi } from '@/api/users'

export const useUsersStore = defineStore('users', () => {
  const users = ref([])
  const loading = ref(false)
  const error = ref(null)
  const totalCount = ref(0)
  
  async function fetchUsers(params = {}) {
    loading.value = true
    error.value = null
    
    try {
      const response = await usersApi.getUsers(params)
      users.value = response.users || []
      totalCount.value = response.total || users.value.length
    } catch (err) {
      error.value = err.message || 'Ошибка загрузки пользователей'
    } finally {
      loading.value = false
    }
  }
  
  async function updateUserRole(userId, role) {
    try {
      await usersApi.updateRole(userId, role)
      await fetchUsers()
    } catch (err) {
      error.value = err.message || 'Ошибка обновления роли'
      throw err
    }
  }
  
  return {
    users,
    loading,
    error,
    totalCount,
    fetchUsers,
    updateUserRole
  }
})
