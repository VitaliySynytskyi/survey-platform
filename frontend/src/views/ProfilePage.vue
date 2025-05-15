<template>
  <div class="profile-page">
    <div class="container">
      <h1>Profile</h1>
      <div v-if="loading" class="loading">Loading profile data...</div>
      
      <div v-else-if="user" class="profile-container">
        <div class="profile-info">
          <h2>User Information</h2>
          <div class="info-item">
            <strong>Email:</strong> {{ user.email }}
          </div>
          <div class="info-item">
            <strong>Role:</strong> {{ user.role }}
          </div>
          <div class="info-item">
            <strong>Member since:</strong> {{ formatDate(user.createdAt) }}
          </div>
        </div>
        
        <div class="profile-stats">
          <h2>Account Statistics</h2>
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-value">{{ surveyCount }}</div>
              <div class="stat-label">Surveys Created</div>
            </div>
            <!-- Add more stats here in the future -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const loading = ref(true)
const error = ref('')
const surveyCount = ref(0)

const user = computed(() => authStore.user)

onMounted(async () => {
  try {
    await authStore.fetchUser()
    // Here you would also fetch survey count or other user stats
    // For example: const stats = await userStatsService.getUserStats()
    // surveyCount.value = stats.surveyCount
    surveyCount.value = 0 // Placeholder
  } catch (e: any) {
    error.value = e.message || 'Failed to load profile'
  } finally {
    loading.value = false
  }
})

// Format date helper function
const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString()
}
</script>

<style scoped>
.profile-page {
  padding: 2rem 0;
}

.profile-page h1 {
  margin-bottom: 1.5rem;
  color: var(--primary-color);
}

.loading {
  text-align: center;
  padding: 2rem;
  color: var(--gray);
}

.profile-container {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

@media (min-width: 768px) {
  .profile-container {
    grid-template-columns: 1fr 1fr;
  }
}

.profile-info, .profile-stats {
  background-color: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: 1.5rem;
}

.profile-info h2, .profile-stats h2 {
  margin-bottom: 1rem;
  color: var(--primary-color);
  font-size: 1.3rem;
}

.info-item {
  margin-bottom: 0.5rem;
  padding: 0.5rem 0;
  border-bottom: 1px solid var(--light-gray);
}

.info-item:last-child {
  border-bottom: none;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(120px, 1fr));
  gap: 1rem;
}

.stat-item {
  text-align: center;
  background-color: var(--light-gray);
  padding: 1rem;
  border-radius: var(--border-radius);
}

.stat-value {
  font-size: 1.8rem;
  font-weight: bold;
  color: var(--primary-color);
}

.stat-label {
  color: var(--dark-gray);
  font-size: 0.9rem;
  margin-top: 0.5rem;
}
</style> 