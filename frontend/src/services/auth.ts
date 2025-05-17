import apiClient from './axios'
import { User } from '@/types/user'

interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

interface TokenResponse {
  access_token: string
  refresh_token: string
  expires_in: number
}

const authService = {
  /**
   * Register a new user
   */
  async register(email: string, password: string): Promise<void> {
    console.log("Sending registration request with:", { email, password });
    try {
      const response = await apiClient.post('/auth/register', { email, password });
      console.log("Registration response:", response.data);
    } catch (error: any) {
      console.error("Registration error:", error.response?.data || error.message);
      throw error;
    }
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
  async logout(refreshTokenValue: string | null): Promise<void> {
    const payload = refreshTokenValue ? { refresh_token: refreshTokenValue } : {};
    await apiClient.post('/auth/logout', payload)
  },

  /**
   * Get new access token using refresh token
   */
  async refreshToken(refreshToken: string): Promise<TokenResponse> {
    const response = await apiClient.post<TokenResponse>('/auth/refresh', { refresh_token: refreshToken })
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