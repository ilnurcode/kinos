import axios from 'axios'

const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json'
  },
  withCredentials: true
})

// Request interceptor - добавить JWT токен
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('access_token')
    
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor - handle 401 и refresh token
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config
    
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      try {
        // Попытка refresh токена
        const response = await axios.post(
          `${import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api'}/users/refresh`,
          {},
          { withCredentials: true }
        )
        
        const newToken = response.data.access_token
        localStorage.setItem('access_token', newToken)
        
        originalRequest.headers.Authorization = `Bearer ${newToken}`
        return apiClient(originalRequest)
      } catch (refreshError) {
        // Refresh не удался - logout
        localStorage.removeItem('access_token')
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }
    
    return Promise.reject(error)
  }
)

export default apiClient
