<template>
  <div class="survey-results-page">
    <div v-if="loading" class="loading-container">
      <div class="spinner"></div>
      <p>Loading survey results...</p>
    </div>

    <div v-else-if="error" class="error-container">
      <h2>Error</h2>
      <p>{{ error }}</p>
      <button @click="fetchResults" class="btn btn-primary">Try Again</button>
    </div>

    <div v-else-if="surveyResults" class="survey-results-container">
      <div class="survey-header">
        <h1>{{ surveyResults.title }}</h1>
        <p class="description">{{ surveyResults.description }}</p>
        <div class="survey-stats">
          <div class="stat-item">
            <span class="stat-value">{{ surveyResults.total_responses }}</span>
            <span class="stat-label">Total Responses</span>
          </div>
          <div v-if="surveyResults.completion_rate" class="stat-item">
            <span class="stat-value">{{ Math.round(surveyResults.completion_rate) }}%</span>
            <span class="stat-label">Completion Rate</span>
          </div>
          <div v-if="surveyResults.average_time_spent" class="stat-item">
            <span class="stat-value">{{ surveyResults.average_time_spent }}</span>
            <span class="stat-label">Avg. Time</span>
          </div>
        </div>
      </div>

      <div v-if="surveyResults.total_responses === 0" class="no-responses">
        <p>No responses have been collected for this survey yet.</p>
      </div>

      <div v-else class="question-results">
        <div 
          v-for="question in surveyResults.question_results" 
          :key="question.question_id"
          class="question-result"
        >
          <h2>{{ question.title }}</h2>
          <p class="total-answers">{{ question.total_answers }} answers</p>

          <!-- Single Choice / Multiple Choice Questions -->
          <div 
            v-if="['single-choice', 'multiple-choice'].includes(question.type)" 
            class="chart-container"
          >
            <div class="bar-chart">
              <div 
                v-for="(option, optionId) in question.analytics.option_counts" 
                :key="optionId"
                class="chart-bar-container"
              >
                <div class="chart-label">{{ option.text }}</div>
                <div class="chart-bar-wrapper">
                  <div 
                    class="chart-bar" 
                    :style="{ width: `${option.percent}%` }"
                    :class="{ 'low-percent': option.percent < 5 }"
                  ></div>
                  <span class="chart-value">{{ option.count }} ({{ Math.round(option.percent) }}%)</span>
                </div>
              </div>
            </div>
          </div>

          <!-- Scale Questions -->
          <div v-else-if="question.type === 'scale'" class="scale-results">
            <div class="scale-stats">
              <div class="stat-box">
                <span class="stat-title">Average</span>
                <span class="stat-value">{{ question.analytics.average?.toFixed(2) }}</span>
              </div>
              <div class="stat-box">
                <span class="stat-title">Median</span>
                <span class="stat-value">{{ question.analytics.median?.toFixed(2) }}</span>
              </div>
              <div class="stat-box">
                <span class="stat-title">Min</span>
                <span class="stat-value">{{ question.analytics.min }}</span>
              </div>
              <div class="stat-box">
                <span class="stat-title">Max</span>
                <span class="stat-value">{{ question.analytics.max }}</span>
              </div>
            </div>
            <!-- Here would be a distribution chart component -->
          </div>

          <!-- Open Text Questions -->
          <div v-else-if="question.type === 'open-text'" class="open-text-results">
            <div class="text-responses">
              <h3>Responses</h3>
              <ul class="response-list">
                <li 
                  v-for="(response, index) in question.response_data.slice(0, showAllText ? undefined : 5)" 
                  :key="index"
                  class="text-response"
                >
                  {{ response }}
                </li>
              </ul>
              <button 
                v-if="!showAllText && question.response_data.length > 5" 
                @click="showAllText = true" 
                class="btn btn-text"
              >
                Show all {{ question.response_data.length }} responses
              </button>
              <button 
                v-else-if="showAllText"
                @click="showAllText = false" 
                class="btn btn-text"
              >
                Show less
              </button>
            </div>

            <!-- Word Frequency -->
            <div v-if="question.analytics.word_frequency" class="word-frequency">
              <h3>Common Words</h3>
              <div class="word-cloud">
                <span 
                  v-for="(count, word) in topWords(question.analytics.word_frequency)" 
                  :key="word"
                  class="word-item"
                  :style="{ fontSize: `${Math.min(count * 0.5 + 1, 2.5)}em` }"
                >
                  {{ word }}
                </span>
              </div>
            </div>
          </div>

          <!-- Matrix Questions -->
          <div v-else-if="question.type === 'matrix'" class="matrix-results">
            <!-- Matrix visualization would go here -->
            <p>Matrix visualization not yet implemented</p>
          </div>

          <!-- Default for other question types -->
          <div v-else class="default-results">
            <pre>{{ JSON.stringify(question.response_data, null, 2) }}</pre>
          </div>
        </div>
      </div>

      <div class="actions">
        <router-link :to="`/surveys/${surveyId}/responses`" class="btn btn-secondary">
          View Individual Responses
        </router-link>
        <router-link :to="`/my-surveys`" class="btn btn-outline">
          Back to My Surveys
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import surveyService from '@/services/survey'

const route = useRoute()
const router = useRouter()
const surveyId = computed(() => route.params.id as string)

// State
const surveyResults = ref<any>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const showAllText = ref(false)

// Fetch survey results
const fetchResults = async () => {
  loading.value = true
  error.value = null
  
  try {
    surveyResults.value = await surveyService.getSurveyResults(surveyId.value)
  } catch (err: any) {
    console.error('Failed to fetch survey results:', err)
    error.value = err.response?.data?.error || 'Failed to load survey results'
    
    // If the error is 403 or 401, redirect to login
    if (err.response?.status === 401 || err.response?.status === 403) {
      router.push('/login')
    }
  } finally {
    loading.value = false
  }
}

// Helper to get top words (up to 15)
const topWords = (wordFrequency: Record<string, number>) => {
  return Object.entries(wordFrequency)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 15)
    .reduce((obj, [word, count]) => {
      obj[word] = count
      return obj
    }, {} as Record<string, number>)
}

// Fetch results on mount
onMounted(fetchResults)
</script>

<style scoped>
.survey-results-page {
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

.survey-header {
  margin-bottom: 30px;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.survey-stats {
  display: flex;
  gap: 20px;
  margin-top: 15px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 10px 15px;
  background-color: #f9f9f9;
  border-radius: 5px;
}

.stat-value {
  font-size: 1.4em;
  font-weight: bold;
  color: #2c3e50;
}

.stat-label {
  font-size: 0.8em;
  color: #666;
}

.no-responses {
  text-align: center;
  padding: 40px;
  background-color: #f9f9f9;
  border-radius: 5px;
}

.question-result {
  margin-bottom: 40px;
  padding: 20px;
  background-color: white;
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.05);
}

.chart-container {
  margin-top: 15px;
}

.bar-chart {
  width: 100%;
}

.chart-bar-container {
  margin-bottom: 15px;
  display: flex;
  align-items: center;
}

.chart-label {
  width: 150px;
  margin-right: 10px;
  text-align: right;
  font-size: 0.9em;
  color: #333;
}

.chart-bar-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
}

.chart-bar {
  height: 25px;
  background-color: #3498db;
  border-radius: 3px;
  min-width: 2%;
  transition: width 0.5s ease-in-out;
}

.chart-value {
  margin-left: 10px;
  font-size: 0.9em;
  min-width: 80px;
}

.low-percent {
  min-width: 2%;
}

.scale-stats {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
  margin-top: 15px;
}

.stat-box {
  padding: 15px;
  background-color: #f9f9f9;
  border-radius: 5px;
  min-width: 100px;
  text-align: center;
}

.text-responses {
  margin-top: 15px;
}

.response-list {
  list-style: none;
  padding: 0;
}

.text-response {
  padding: 10px;
  margin-bottom: 10px;
  background-color: #f9f9f9;
  border-radius: 3px;
  border-left: 3px solid #3498db;
}

.word-cloud {
  margin-top: 15px;
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  padding: 20px;
  background-color: #f9f9f9;
  border-radius: 5px;
}

.word-item {
  color: #2c3e50;
  display: inline-block;
  padding: 5px;
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

.total-answers {
  color: #666;
  font-size: 0.9em;
  margin-top: -10px;
}
</style> 