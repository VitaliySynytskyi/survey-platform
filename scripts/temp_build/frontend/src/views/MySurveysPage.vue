<template>
  <div class="my-surveys-page">
    <div class="container">
      <div class="page-header">
        <h1>My Surveys</h1>
        <router-link to="/surveys/new" class="btn">Create New Survey</router-link>
      </div>
      
      <div v-if="loading" class="loading">Loading surveys...</div>
      
      <div v-else-if="error" class="error-text">{{ error }}</div>
      
      <div v-else-if="surveys.length === 0" class="empty-state">
        <p>You haven't created any surveys yet.</p>
        <router-link to="/surveys/new" class="btn">Create Your First Survey</router-link>
      </div>
      
      <div v-else class="surveys-list">
        <div v-for="survey in surveys" :key="survey.id" class="survey-card">
          <div class="survey-card-header">
            <h3>{{ survey.title }}</h3>
            <div class="survey-meta">
              <span class="questions-count">{{ survey.questions.length }} questions</span>
              <span class="date">Created: {{ formatDate(survey.created_at) }}</span>
            </div>
          </div>
          
          <p v-if="survey.description" class="survey-description">
            {{ survey.description }}
          </p>
          
          <div class="survey-actions">
            <router-link :to="`/surveys/${survey.id}/edit`" class="btn-edit">
              <span class="btn-icon">✏️</span>
              Edit
            </router-link>
            <button @click="confirmDelete(survey)" class="btn-delete">
              <span class="btn-icon">🗑️</span>
              Delete
            </button>
          </div>
        </div>
        
        <!-- Pagination controls -->
        <div v-if="totalPages > 1" class="pagination">
          <button 
            :disabled="currentPage === 1" 
            @click="changePage(currentPage - 1)" 
            class="pagination-btn"
          >
            Previous
          </button>
          
          <span class="pagination-info">
            Page {{ currentPage }} of {{ totalPages }}
          </span>
          
          <button 
            :disabled="currentPage === totalPages" 
            @click="changePage(currentPage + 1)" 
            class="pagination-btn"
          >
            Next
          </button>
        </div>
      </div>
    </div>
    
    <!-- Delete confirmation modal (simplified) -->
    <div v-if="showDeleteModal" class="modal-overlay">
      <div class="modal-content">
        <h3>Confirm Delete</h3>
        <p>Are you sure you want to delete survey "{{ surveyToDelete?.title }}"?</p>
        <p>This action cannot be undone.</p>
        <div class="modal-actions">
          <button @click="deleteConfirmed" class="btn-delete">Delete</button>
          <button @click="cancelDelete" class="btn-cancel">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import surveyService from '@/services/survey'
import { Survey } from '@/types/survey'

const router = useRouter()
const authStore = useAuthStore()

const surveys = ref<Survey[]>([])
const loading = ref(true)
const error = ref('')
const currentPage = ref(1)
const totalSurveys = ref(0)
const perPage = ref(10)

// For delete modal
const showDeleteModal = ref(false)
const surveyToDelete = ref<Survey | null>(null)

// Computed properties for pagination
const totalPages = computed(() => Math.ceil(totalSurveys.value / perPage.value))

onMounted(async () => {
  await fetchSurveys()
})

const fetchSurveys = async () => {
  if (!authStore.user) {
    return
  }
  
  loading.value = true
  error.value = ''
  
  try {
    const response = await surveyService.getUserSurveys(
      authStore.user.id, 
      currentPage.value, 
      perPage.value
    )
    
    surveys.value = response.surveys
    totalSurveys.value = response.totalCount
  } catch (e: any) {
    error.value = e.message || 'Failed to fetch surveys'
  } finally {
    loading.value = false
  }
}

const changePage = async (page: number) => {
  currentPage.value = page
  await fetchSurveys()
  // Scroll to top when changing page
  window.scrollTo(0, 0)
}

const confirmDelete = (survey: Survey) => {
  surveyToDelete.value = survey
  showDeleteModal.value = true
}

const cancelDelete = () => {
  showDeleteModal.value = false
  surveyToDelete.value = null
}

const deleteConfirmed = async () => {
  if (!surveyToDelete.value) return
  
  loading.value = true
  
  try {
    await surveyService.deleteSurvey(surveyToDelete.value.id)
    // Remove from the local list
    surveys.value = surveys.value.filter(s => s.id !== surveyToDelete.value?.id)
    showDeleteModal.value = false
    surveyToDelete.value = null
    
    // Reload if we emptied the page
    if (surveys.value.length === 0 && currentPage.value > 1) {
      currentPage.value -= 1
      await fetchSurveys()
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to delete survey'
  } finally {
    loading.value = false
  }
}

// Format date helper function
const formatDate = (dateString: string) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString()
}
</script>

<style scoped>
.my-surveys-page {
  padding: 2rem 0;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  margin: 0;
  color: var(--primary-color);
}

.loading {
  text-align: center;
  padding: 2rem;
  color: var(--gray);
}

.empty-state {
  text-align: center;
  padding: 3rem 0;
}

.empty-state p {
  margin-bottom: 1.5rem;
  font-size: 1.2rem;
  color: var(--dark-gray);
}

.surveys-list {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

.survey-card {
  background-color: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: 1.5rem;
}

.survey-card-header {
  margin-bottom: 1rem;
}

.survey-card-header h3 {
  color: var(--primary-color);
  margin: 0 0 0.5rem 0;
}

.survey-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.9rem;
  color: var(--dark-gray);
}

.survey-description {
  margin-bottom: 1.5rem;
  color: var(--text-color);
}

.survey-actions {
  display: flex;
  gap: 1rem;
}

.btn-edit,
.btn-delete {
  padding: 0.5rem 1rem;
  border-radius: var(--border-radius);
  font-size: 0.9rem;
  display: inline-flex;
  align-items: center;
  cursor: pointer;
}

.btn-edit {
  background-color: var(--secondary-color);
  color: var(--white);
  text-decoration: none;
}

.btn-delete {
  background-color: #f44336;
  color: var(--white);
  border: none;
}

.btn-icon {
  margin-right: 0.5rem;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 2rem;
  gap: 1rem;
}

.pagination-btn {
  padding: 0.5rem 1rem;
  border: 1px solid var(--gray);
  background-color: var(--white);
  border-radius: var(--border-radius);
  cursor: pointer;
}

.pagination-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.pagination-info {
  color: var(--dark-gray);
}

/* Modal styles */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  max-width: 500px;
  width: 90%;
  padding: 2rem;
}

.modal-content h3 {
  margin-top: 0;
  color: var(--text-color);
  margin-bottom: 1rem;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 1.5rem;
}

.btn-cancel {
  padding: 0.5rem 1rem;
  border: 1px solid var(--gray);
  background-color: var(--white);
  border-radius: var(--border-radius);
  cursor: pointer;
}

@media (min-width: 768px) {
  .surveys-list {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style> 