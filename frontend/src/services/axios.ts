import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

// Create Axios instance
const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1',
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
apiClient.interceptors.request.use(
  config => {
    const authStore = useAuthStore()
    
    // If token exists and Authorization header is not set
    if (authStore.accessToken && config.headers && !config.headers.Authorization) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }
    
    return config
  },
  error => {
    return Promise.reject(error)
  }
)

// Response interceptor
apiClient.interceptors.response.use(
  response => response,
  async error => {
    // Check if error was 401 Unauthorized
    if (error.response && error.response.status === 401) {
      const authStore = useAuthStore()
      const originalRequest = error.config
      
      // If refresh token is available and we haven't tried to refresh yet
      if (authStore.refreshToken && !originalRequest._retry) {
        originalRequest._retry = true
        
        try {
          // Try to refresh the token
          await authStore.refreshSession()
          
          // Update Authorization header with the new token
          originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`
          
          // Retry the original request
          return apiClient(originalRequest)
        } catch (refreshError) {
          // If refresh failed, redirect to login
          await authStore.logout()
          window.location.href = '/login'
          return Promise.reject(refreshError)
        }
      }
    }
    
    return Promise.reject(error)
  }
)

export default apiClient 