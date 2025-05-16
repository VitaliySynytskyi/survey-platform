<template>
  <div class="login-page">
    <div class="container">
      <div class="auth-form-container">
        <h1 class="text-center">Login</h1>
        
        <form @submit.prevent="handleLogin" class="auth-form">
          <div class="form-group">
            <label for="email">Email</label>
            <input 
              type="email" 
              id="email" 
              v-model="email" 
              required 
              :disabled="loading"
              placeholder="Enter your email"
            />
          </div>
          
          <div class="form-group">
            <label for="password">Password</label>
            <input 
              type="password" 
              id="password" 
              v-model="password" 
              required 
              :disabled="loading"
              placeholder="Enter your password"
            />
          </div>
          
          <div v-if="error" class="error-text mb-2">{{ error }}</div>
          
          <button type="submit" class="btn" :disabled="loading">
            {{ loading ? 'Logging in...' : 'Login' }}
          </button>
          
          <div class="auth-links mt-3">
            <p>Don't have an account? <router-link to="/register">Register</router-link></p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const error = ref('')

const handleLogin = async () => {
  error.value = ''
  loading.value = true
  
  try {
    await authStore.login(email.value, password.value)
    router.push('/')
  } catch (e: any) {
    error.value = e.message || 'Login failed. Please check your credentials.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  padding: 2rem 0;
}

.auth-form-container {
  max-width: 450px;
  margin: 0 auto;
  padding: 2rem;
  background-color: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
}

.auth-form-container h1 {
  margin-bottom: 1.5rem;
  color: var(--primary-color);
}

.auth-form .btn {
  width: 100%;
  padding: 0.75rem;
  font-size: 1.1rem;
  margin-top: 1rem;
}

.auth-links {
  text-align: center;
}
</style> 