<template>
  <div class="app">
    <header class="header">
      <div class="container header-container">
        <div class="logo">
          <router-link to="/">Survey Platform</router-link>
        </div>
        <nav class="main-nav">
          <router-link to="/" class="nav-link">Home</router-link>
          <template v-if="isAuthenticated">
            <router-link to="/my-surveys" class="nav-link">My Surveys</router-link>
            <router-link to="/profile" class="nav-link">Profile</router-link>
            <button @click="logout" class="nav-link logout-btn">Logout</button>
          </template>
          <template v-else>
            <router-link to="/login" class="nav-link">Login</router-link>
            <router-link to="/register" class="nav-link">Register</router-link>
          </template>
        </nav>
      </div>
    </header>
    <main class="main">
      <div class="container">
        <router-view />
      </div>
    </main>
    <footer class="footer">
      <div class="container">
        <p>&copy; {{ new Date().getFullYear() }} Survey Platform. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const isAuthenticated = computed(() => authStore.isAuthenticated)

const logout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('Logout error:', error)
  }
}
</script>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  background-color: var(--white);
  box-shadow: var(--box-shadow);
  padding: 1rem 0;
}

.header-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logo {
  font-size: 1.5rem;
  font-weight: bold;
}

.logo a {
  color: var(--primary-color);
}

.logo a:hover {
  text-decoration: none;
}

.main-nav {
  display: flex;
  align-items: center;
}

.nav-link {
  margin-left: 1.5rem;
  color: var(--text-color);
  font-weight: 500;
}

.nav-link:hover {
  color: var(--primary-color);
  text-decoration: none;
}

.nav-link.router-link-active {
  color: var(--primary-color);
}

.logout-btn {
  background: none;
  border: none;
  color: var(--text-color);
  font-weight: 500;
  font-size: 1rem;
  cursor: pointer;
  padding: 0;
  margin-left: 1.5rem;
}

.logout-btn:hover {
  color: var(--primary-color);
  background: none;
}

.main {
  flex: 1;
  padding: 2rem 0;
}

.footer {
  background-color: var(--white);
  padding: 1.5rem 0;
  margin-top: auto;
  text-align: center;
  color: var(--dark-gray);
}
</style> 