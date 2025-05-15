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
      
      user.value = response.user
      accessToken.value = response.accessToken
      refreshToken.value = response.refreshToken

      // Save to localStorage
      localStorage.setItem('user', JSON.stringify(response.user))
      localStorage.setItem('accessToken', response.accessToken)
      localStorage.setItem('refreshToken', response.refreshToken)

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

    try {
      if (accessToken.value) {
        await authService.logout(accessToken.value)
      }
    } catch (e) {
      console.error('Logout error:', e)
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
      
      accessToken.value = response.accessToken
      refreshToken.value = response.refreshToken

      // Update localStorage
      localStorage.setItem('accessToken', response.accessToken)
      localStorage.setItem('refreshToken', response.refreshToken)

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