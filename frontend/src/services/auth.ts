import apiClient from './axios'
import { User } from '@/types/user'

interface LoginResponse {
  user: User
  accessToken: string
  refreshToken: string
}

interface TokenResponse {
  accessToken: string
  refreshToken: string
}

const authService = {
  /**
   * Register a new user
   */
  async register(email: string, password: string): Promise<void> {
    await apiClient.post('/auth/register', { email, password })
  },

  /**
   * Login user and get tokens
   */
  async login(email: string, password: string): Promise<LoginResponse> {
    const response = await apiClient.post<LoginResponse>('/auth/login', { email, password })
    return response.data
  },

  /**
   * Logout user
   */
  async logout(token: string): Promise<void> {
    await apiClient.post('/auth/logout', {}, {
      headers: { Authorization: `Bearer ${token}` }
    })
  },

  /**
   * Get new access token using refresh token
   */
  async refreshToken(refreshToken: string): Promise<TokenResponse> {
    const response = await apiClient.post<TokenResponse>('/auth/refresh', { refreshToken })
    return response.data
  },

  /**
   * Get current user profile
   */
  async getUserProfile(token: string): Promise<User> {
    const response = await apiClient.get<User>('/auth/me', {
      headers: { Authorization: `Bearer ${token}` }
    })
    return response.data
  }
}

export default authService 