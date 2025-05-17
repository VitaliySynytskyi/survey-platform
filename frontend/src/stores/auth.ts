import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import authService from '@/services/auth'
import { User } from '@/types/user'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(null)
  const refreshToken = ref<string | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  // Initialize state from localStorage
  const initFromStorage = () => {
    const storedUser = localStorage.getItem('user')
    const storedAccessToken = localStorage.getItem('accessToken')
    const storedRefreshToken = localStorage.getItem('refreshToken')

    if (storedUser) {
      try {
        user.value = JSON.parse(storedUser)
      } catch (e) {
        console.error('Failed to parse stored user data:', e)
        localStorage.removeItem('user')
      }
    }

    if (storedAccessToken) accessToken.value = storedAccessToken
    if (storedRefreshToken) refreshToken.value = storedRefreshToken
  }

  // Call init
  initFromStorage()

  // Getters
  const isAuthenticated = computed(() => !!accessToken.value)

  // Actions
  const register = async (email: string, password: string) => {
    loading.value = true
    error.value = null

    try {
      await authService.register(email, password)
      return true
    } catch (e: any) {
      error.value = e.message || 'Registration failed'
      throw e
    } finally {
      loading.value = false
    }
  }

  const login = async (email: string, password: string) => {
    loading.value = true
    error.value = null

    try {
      const response = await authService.login(email, password)
      
      // Store tokens from response
      accessToken.value = response.access_token
      refreshToken.value = response.refresh_token

      // Save tokens to localStorage
      localStorage.setItem('accessToken', response.access_token)
      localStorage.setItem('refreshToken', response.refresh_token)

      // Fetch user profile with the access token
      try {
        await fetchUser()
      } catch (userError) {
        console.error('Failed to fetch user profile:', userError)
        // Continue with login even if profile fetch fails
      }

      return true
    } catch (e: any) {
      error.value = e.message || 'Login failed'
      throw e
    } finally {
      loading.value = false
    }
  }

  const logout = async () => {
    loading.value = true
    error.value = null
    const currentRefreshToken = refreshToken.value; // Get refreshToken before clearing

    try {
      // Pass the refresh token to the authService.logout method
      // The authService.logout method now expects the refresh token
      await authService.logout(currentRefreshToken) 
    } catch (e) {
      console.error('Logout error:', e)
      // Even if logout API call fails, proceed to clear local state and storage
    } finally {
      // Reset state
      user.value = null
      accessToken.value = null
      refreshToken.value = null

      // Clear localStorage
      localStorage.removeItem('user')
      localStorage.removeItem('accessToken')
      localStorage.removeItem('refreshToken')

      loading.value = false
    }
  }

  const refreshSession = async () => {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    loading.value = true
    error.value = null

    try {
      const response = await authService.refreshToken(refreshToken.value)
      
      accessToken.value = response.access_token
      refreshToken.value = response.refresh_token

      // Update localStorage
      localStorage.setItem('accessToken', response.access_token)
      localStorage.setItem('refreshToken', response.refresh_token)

      return true
    } catch (e: any) {
      error.value = e.message || 'Session refresh failed'
      // If refresh fails, logout the user
      await logout()
      throw e
    } finally {
      loading.value = false
    }
  }

  const fetchUser = async () => {
    if (!accessToken.value) {
      throw new Error('No access token available')
    }

    loading.value = true
    error.value = null

    try {
      const userData = await authService.getUserProfile(accessToken.value)
      user.value = userData
      localStorage.setItem('user', JSON.stringify(userData))
      return userData
    } catch (e: any) {
      error.value = e.message || 'Failed to fetch user profile'
      throw e
    } finally {
      loading.value = false
    }
  }

  return {
    user,
    accessToken,
    refreshToken,
    loading,
    error,
    isAuthenticated,
    register,
    login,
    logout,
    refreshSession,
    fetchUser
  }
}) 