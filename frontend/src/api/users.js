import apiClient from './client'

export const usersApi = {
  async getProfile() {
    const response = await apiClient.get('/profile')
    return response.data
  },
  
  async updateProfile(profileData) {
    const response = await apiClient.put('/profile', profileData)
    return response.data
  },
  
  async getUsers(params = {}) {
    const response = await apiClient.get('/admin/users', { params })
    return response.data
  },
  
  async updateRole(userId, role) {
    const response = await apiClient.put('/admin/users/role', { 
      user_id: userId, 
      role 
    })
    return response.data
  }
}
