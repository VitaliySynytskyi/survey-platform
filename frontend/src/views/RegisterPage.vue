<template>
  <div class="register-page">
    <div class="container">
      <div class="auth-form-container">
        <h1 class="text-center">Register</h1>
        
        <form @submit.prevent="handleRegister" class="auth-form">
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
          
          <div class="form-group">
            <label for="confirmPassword">Confirm Password</label>
            <input 
              type="password" 
              id="confirmPassword" 
              v-model="confirmPassword" 
              required 
              :disabled="loading"
              placeholder="Confirm your password"
            />
            <div v-if="!passwordsMatch && confirmPassword" class="error-text">
              Passwords do not match
            </div>
          </div>
          
          <div v-if="error" class="error-text mb-2">{{ error }}</div>
          
          <button 
            type="submit" 
            class="btn" 
            :disabled="loading || !passwordsMatch || !email || !password || !confirmPassword"
          >
            {{ loading ? 'Registering...' : 'Register' }}
          </button>
          
          <div class="auth-links mt-3">
            <p>Already have an account? <router-link to="/login">Login</router-link></p>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')

const passwordsMatch = computed(() => {
  if (!confirmPassword.value) return true
  return password.value === confirmPassword.value
})

const handleRegister = async () => {
  if (!passwordsMatch.value) {
    error.value = 'Passwords do not match'
    return
  }
  
  error.value = ''
  loading.value = true
  
  try {
    await authStore.register(email.value, password.value)
    // After successful registration, redirect to login
    router.push('/login')
  } catch (e: any) {
    error.value = e.message || 'Registration failed. Please try again.'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
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