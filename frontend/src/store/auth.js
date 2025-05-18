import { defineStore } from 'pinia'
import axios from '../utils/axiosConfig'
import router from '../router'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: localStorage.getItem('token') || null,
    refreshToken: localStorage.getItem('refreshToken') || null,
  }),
  
  getters: {
    isAuthenticated: (state) => !!state.token,
    getUser: (state) => state.user,
    getToken: (state) => state.token,
  },
  
  actions: {
    async register(userData) {
      try {
        const response = await axios.post('/api/v1/auth/register', userData)
        
        this.setAuthData(response.data)
        return response.data
      } catch (error) {
        throw error.response ? error.response.data : error
      }
    },
    
    async login(credentials) {
      try {
        const response = await axios.post('/api/v1/auth/login', credentials)
        
        this.setAuthData(response.data)
        return response.data
      } catch (error) {
        throw error.response ? error.response.data : error
      }
    },
    
    async fetchUser() {
      try {
        if (!this.token) return null
        
        const response = await axios.get('/api/v1/users/me')
        this.user = response.data
        return response.data
      } catch (error) {
        // If unauthorized, logout the user
        if (error.response && error.response.status === 401) {
          this.logout()
        }
        throw error.response ? error.response.data : error
      }
    },
    
    async refreshAuthToken() {
      try {
        if (!this.refreshToken) throw new Error('No refresh token available')
        
        const response = await axios.post('/api/v1/auth/refresh', {
          refresh_token: this.refreshToken
        })
        
        this.setAuthData(response.data)
        return response.data
      } catch (error) {
        this.logout()
        throw error.response ? error.response.data : error
      }
    },
    
    setAuthData(data) {
      this.token = data.token
      this.refreshToken = data.refresh_token
      this.user = data.user
      
      localStorage.setItem('token', data.token)
      localStorage.setItem('refreshToken', data.refresh_token)
      
      // Set axios default headers
      axios.defaults.headers.common['Authorization'] = `Bearer ${data.token}`
    },
    
    logout() {
      this.user = null
      this.token = null
      this.refreshToken = null
      
      localStorage.removeItem('token')
      localStorage.removeItem('refreshToken')
      
      // Remove axios default headers
      delete axios.defaults.headers.common['Authorization']
      
      // Redirect to login page
      router.push('/login')
    }
  }
})

// Set up axios interceptors for auth token refreshing
axios.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config
    
    // If the error is due to an expired token and we haven't already tried to refresh it
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      try {
        const authStore = useAuthStore()
        await authStore.refreshAuthToken()
        
        // Retry the original request with the new token
        originalRequest.headers['Authorization'] = `Bearer ${authStore.getToken}`
        return axios(originalRequest)
      } catch (refreshError) {
        // If refresh failed, logout and redirect to login
        return Promise.reject(refreshError)
      }
    }
    
    return Promise.reject(error)
  }
) 