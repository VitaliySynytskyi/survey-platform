<template>
  <div class="home-page">
    <div class="hero">
      <h1>Survey Platform</h1>
      <p>Create and share surveys easily</p>
      <div class="hero-actions">
        <router-link v-if="isAuthenticated" to="/surveys/new" class="btn">Create Survey</router-link>
        <router-link v-else to="/register" class="btn">Get Started</router-link>
      </div>
    </div>
    
    <div class="features">
      <div class="container">
        <h2 class="text-center mb-4">Features</h2>
        <div class="features-grid">
          <div class="feature">
            <h3>Create Surveys</h3>
            <p>Design custom surveys with various question types</p>
          </div>
          <div class="feature">
            <h3>Collect Responses</h3>
            <p>Share your surveys and collect responses easily</p>
          </div>
          <div class="feature">
            <h3>Analyze Results</h3>
            <p>View and export survey results</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
</script>

<style scoped>
.home-page {
  text-align: center;
}

.hero {
  padding: 3rem 1rem;
  background-color: var(--light-gray);
  margin-bottom: 2rem;
}

.hero h1 {
  font-size: 2.5rem;
  margin-bottom: 1rem;
  color: var(--primary-color);
}

.hero p {
  font-size: 1.2rem;
  margin-bottom: 2rem;
  color: var(--dark-gray);
}

.hero-actions {
  margin-top: 2rem;
}

.hero-actions .btn {
  font-size: 1.1rem;
  padding: 0.75rem 1.5rem;
}

.features {
  padding: 3rem 0;
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-top: 2rem;
}

.feature {
  padding: 1.5rem;
  background-color: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
}

.feature h3 {
  color: var(--primary-color);
  margin-bottom: 0.5rem;
}
</style> 