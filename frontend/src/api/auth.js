import apiClient from './client'

export const authApi = {
  async login(email, password) {
    const response = await apiClient.post('/users/login', { email, password })
    return response.data
  },
  
  async register(userData) {
    const response = await apiClient.post('/users/register', userData)
    return response.data
  },
  
  async getProfile() {
    const response = await apiClient.get('/profile')
    return response.data
  },
  
  async updateProfile(profileData) {
    const response = await apiClient.put('/profile', profileData)
    return response.data
  },
  
  async refresh() {
    const response = await apiClient.post('/users/refresh')
    return response.data
  },
  
  async revoke() {
    await apiClient.post('/users/revoke')
  }
}
