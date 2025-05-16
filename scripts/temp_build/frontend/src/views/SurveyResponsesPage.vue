<template>
  <div class="survey-responses-page">
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading survey responses...</p>
    </div>

    <div v-else-if="error" class="error-container">
      <h2>Error</h2>
      <p>{{ error }}</p>
      <button @click="fetchResponses" class="btn btn-primary">Try Again</button>
    </div>

    <div v-else-if="responsesData" class="survey-responses-container">
      <div class="responses-header">
        <h1>Individual Responses</h1>
        <p v-if="survey">{{ survey.title }}</p>
      </div>

      <div class="filters-container">
        <div class="date-filters">
          <div class="date-filter">
            <label for="start-date">Start Date</label>
            <input 
              type="date" 
              id="start-date" 
              v-model="filters.startDate"
              @change="applyFilters"
            >
          </div>
          <div class="date-filter">
            <label for="end-date">End Date</label>
            <input 
              type="date" 
              id="end-date" 
              v-model="filters.endDate"
              @change="applyFilters"
            >
          </div>
          <button 
            v-if="hasFilters" 
            @click="clearFilters" 
            class="btn btn-text"
          >
            Clear Filters
          </button>
        </div>
      </div>

      <div v-if="responsesData.total_count === 0" class="no-responses">
        <p>No responses match your criteria.</p>
      </div>

      <div v-else class="responses-list">
        <div 
          v-for="response in responsesData.responses" 
          :key="response.id"
          class="response-card"
        >
          <div class="response-meta">
            <div class="response-id">
              <span v-if="response.respondent_id">User: {{ response.respondent_id }}</span>
              <span v-else-if="response.anonymous_id">Anonymous: {{ response.anonymous_id }}</span>
              <span v-else>Anonymous</span>
            </div>
            <div class="response-date">
              {{ formatDate(response.submitted_at) }}
            </div>
          </div>

          <div class="response-answers">
            <div 
              v-for="answer in response.answers" 
              :key="answer.question_id"
              class="response-answer"
            >
              <div class="question-title">
                {{ getQuestionTitle(answer.question_id) }}
              </div>
              <div class="answer-value">
                <!-- Single Choice -->
                <div v-if="getQuestionType(answer.question_id) === 'single-choice'">
                  {{ getOptionText(answer.question_id, answer.value as string) }}
                </div>
                
                <!-- Multiple Choice -->
                <div v-else-if="getQuestionType(answer.question_id) === 'multiple-choice'">
                  <ul class="option-list">
                    <li v-for="option in answer.value as string[]" :key="option">
                      {{ getOptionText(answer.question_id, option) }}
                    </li>
                  </ul>
                </div>
                
                <!-- Scale -->
                <div v-else-if="getQuestionType(answer.question_id) === 'scale'">
                  <div class="scale-value">{{ answer.value }}</div>
                </div>
                
                <!-- Open Text -->
                <div v-else-if="getQuestionType(answer.question_id) === 'open-text'">
                  <div class="text-value">{{ answer.value }}</div>
                </div>
                
                <!-- Matrix -->
                <div v-else-if="getQuestionType(answer.question_id) === 'matrix'">
                  <div class="matrix-value">
                    <pre>{{ JSON.stringify(answer.value, null, 2) }}</pre>
                  </div>
                </div>
                
                <!-- Default -->
                <div v-else>
                  <pre>{{ JSON.stringify(answer.value, null, 2) }}</pre>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="responsesData && responsesData.total_pages > 1" class="pagination">
        <button 
          :disabled="currentPage === 1" 
          @click="prevPage" 
          class="btn btn-outline"
        >
          Previous
        </button>
        
        <span class="pagination-info">
          Page {{ currentPage }} of {{ responsesData.total_pages }}
        </span>
        
        <button 
          :disabled="currentPage === responsesData.total_pages" 
          @click="nextPage" 
          class="btn btn-outline"
        >
          Next
        </button>
      </div>

      <div class="actions">
        <router-link :to="`/surveys/${surveyId}/results`" class="btn btn-secondary">
          Back to Survey Results
        </router-link>
        <router-link :to="`/my-surveys`" class="btn btn-outline">
          Back to My Surveys
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import surveyService from '@/services/survey'

const route = useRoute()
const router = useRouter()
const surveyId = computed(() => route.params.id as string)

// State
const survey = ref<any>(null)
const responsesData = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const currentPage = ref(1)

// Filters
const filters = reactive({
  startDate: '',
  endDate: '',
})

const hasFilters = computed(() => {
  return filters.startDate !== '' || filters.endDate !== ''
})

// Fetch survey details
const fetchSurvey = async () => {
  try {
    survey.value = await surveyService.getSurvey(surveyId.value)
  } catch (err: any) {
    console.error('Failed to fetch survey:', err)
    // We can continue even if this fails
  }
}

// Fetch survey responses
const fetchResponses = async () => {
  loading.value = true
  error.value = null
  
  try {
    // First, get the survey itself to have question info
    if (!survey.value) {
      await fetchSurvey()
    }
    
    // Get responses with current filters
    responsesData.value = await surveyService.getIndividualResponses(
      surveyId.value,
      currentPage.value,
      10, // limit
      filters.startDate || undefined,
      filters.endDate || undefined
    )
  } catch (err: any) {
    console.error('Failed to fetch survey responses:', err)
    error.value = err.response?.data?.error || 'Failed to load survey responses'
    
    // If the error is 403 or 401, redirect to login
    if (err.response?.status === 401 || err.response?.status === 403) {
      router.push('/login')
    }
  } finally {
    loading.value = false
  }
}

// Helper to get question info
const getQuestionTitle = (questionId: string) => {
  if (!survey.value || !survey.value.questions) return 'Unknown Question'
  
  const question = survey.value.questions.find((q: any) => q.id === questionId)
  return question ? question.title : 'Unknown Question'
}

const getQuestionType = (questionId: string) => {
  if (!survey.value || !survey.value.questions) return ''
  
  const question = survey.value.questions.find((q: any) => q.id === questionId)
  return question ? question.type : ''
}

const getOptionText = (questionId: string, optionId: string) => {
  if (!survey.value || !survey.value.questions) return optionId
  
  const question = survey.value.questions.find((q: any) => q.id === questionId)
  if (!question || !question.options) return optionId
  
  const option = question.options.find((o: any) => o.id === optionId)
  return option ? option.text : optionId
}

// Format date
const formatDate = (dateString: string) => {
  return new Date(dateString).toLocaleString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Pagination
const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    fetchResponses()
  }
}

const nextPage = () => {
  if (responsesData.value && currentPage.value < responsesData.value.total_pages) {
    currentPage.value++
    fetchResponses()
  }
}

// Filter functions
const applyFilters = () => {
  currentPage.value = 1 // Reset to first page
  fetchResponses()
}

const clearFilters = () => {
  filters.startDate = ''
  filters.endDate = ''
  currentPage.value = 1
  fetchResponses()
}

// Fetch data on mount
onMounted(fetchResponses)
</script>

<style scoped>
.survey-responses-page {
  padding: 20px;
  max-width: 1000px;
  margin: 0 auto;
}

.loading-container, .error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 300px;
}

.spinner {
  border: 4px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  border-top: 4px solid #3498db;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.responses-header {
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.filters-container {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 5px;
}

.date-filters {
  display: flex;
  gap: 15px;
  align-items: center;
}

.date-filter {
  display: flex;
  flex-direction: column;
}

.date-filter label {
  font-size: 0.8em;
  margin-bottom: 5px;
  color: #666;
}

.date-filter input {
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

.no-responses {
  text-align: center;
  padding: 40px;
  background-color: #f9f9f9;
  border-radius: 5px;
}

.response-card {
  margin-bottom: 25px;
  padding: 20px;
  background-color: white;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.response-meta {
  display: flex;
  justify-content: space-between;
  padding-bottom: 10px;
  margin-bottom: 15px;
  border-bottom: 1px solid #eee;
  font-size: 0.9em;
}

.response-id {
  font-weight: bold;
  color: #2c3e50;
}

.response-date {
  color: #666;
}

.response-answer {
  margin-bottom: 15px;
  padding-bottom: 15px;
  border-bottom: 1px dashed #eee;
}

.response-answer:last-child {
  border-bottom: none;
}

.question-title {
  font-weight: bold;
  margin-bottom: 5px;
}

.answer-value {
  padding: 10px;
  background-color: #f9f9f9;
  border-radius: 4px;
}

.option-list {
  margin: 0;
  padding-left: 20px;
}

.text-value {
  white-space: pre-wrap;
}

.scale-value {
  font-size: 1.2em;
  font-weight: bold;
}

.pagination {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 20px;
  margin-bottom: 20px;
}

.pagination-info {
  margin: 0 15px;
}

.btn {
  display: inline-block;
  padding: 8px 16px;
  margin-right: 10px;
  border-radius: 4px;
  font-size: 14px;
  cursor: pointer;
  text-decoration: none;
  text-align: center;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background-color: #3498db;
  color: white;
  border: none;
}

.btn-secondary {
  background-color: #2c3e50;
  color: white;
  border: none;
}

.btn-outline {
  background-color: transparent;
  color: #2c3e50;
  border: 1px solid #2c3e50;
}

.btn-text {
  background-color: transparent;
  color: #3498db;
  border: none;
  padding: 5px 0;
}

.actions {
  margin-top: 30px;
  display: flex;
  justify-content: flex-end;
}
</style> 