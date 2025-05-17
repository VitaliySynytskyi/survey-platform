<template>
  <div class="survey-results-page">
    <div class="container py-4">
      <div v-if="loading" class="loading-container text-center py-5">
        <div class="spinner">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <p class="mt-3 text-muted">Loading survey results...</p>
      </div>

      <div v-else-if="error" class="error-container alert alert-danger my-4" role="alert">
        <h4 class="alert-heading mb-2">Error</h4>
        <p class="mb-0">{{ error }}</p>
        <div class="mt-3">
          <button @click="fetchResults" class="btn btn-outline-danger">Try Again</button>
          <button @click="goBack" class="btn btn-outline-secondary ms-2">Back to Surveys</button>
        </div>
      </div>

      <div v-else-if="surveyResults" class="survey-results-container">
        <div class="survey-header bg-white rounded shadow-sm p-4 mb-4">
          <div class="d-flex justify-content-between align-items-start mb-4">
            <div>
              <h1 class="mb-2">{{ surveyResults.title }}</h1>
              <p class="description text-muted">{{ surveyResults.description || 'No description provided' }}</p>
            </div>
            <div class="d-flex">
              <button @click="exportToCsv" class="btn btn-outline-primary me-2">
                <i class="bi bi-download me-1"></i> Export
              </button>
              <router-link :to="`/my-surveys`" class="btn btn-outline-secondary">
                <i class="bi bi-arrow-left me-1"></i> Back
              </router-link>
            </div>
          </div>

          <div class="survey-stats d-flex flex-wrap">
            <div class="stat-card">
              <div class="stat-value">{{ surveyResults.total_responses }}</div>
              <div class="stat-label">Total Responses</div>
            </div>
            <div v-if="surveyResults.completion_rate" class="stat-card">
              <div class="stat-value">{{ Math.round(surveyResults.completion_rate) }}%</div>
              <div class="stat-label">Completion Rate</div>
            </div>
            <div v-if="surveyResults.average_time_spent" class="stat-card">
              <div class="stat-value">{{ surveyResults.average_time_spent }}</div>
              <div class="stat-label">Avg. Time Spent</div>
            </div>
          </div>
        </div>

        <div v-if="surveyResults.total_responses === 0" class="no-responses bg-white rounded shadow-sm p-5 text-center">
          <div class="mb-4">
            <i class="bi bi-inbox text-muted" style="font-size: 4rem;"></i>
          </div>
          <h3 class="mb-3">No responses yet</h3>
          <p class="text-muted mb-4">No responses have been collected for this survey yet.</p>
          <div class="d-flex justify-content-center">
            <router-link :to="`/surveys/${surveyId}`" class="btn btn-outline-primary me-3">
              View Survey
            </router-link>
            <button @click="shareSurvey" class="btn btn-primary">
              <i class="bi bi-share me-1"></i> Share Survey
            </button>
          </div>
        </div>

        <div v-else class="question-results">
          <div 
            v-for="question in surveyResults.question_results" 
            :key="question.question_id"
            class="question-result bg-white rounded shadow-sm p-4 mb-4"
          >
            <h2 class="mb-2">{{ question.title }}</h2>
            <p class="total-answers mb-4">{{ question.total_answers }} answers</p>

            <!-- Single Choice / Multiple Choice Questions -->
            <div 
              v-if="['single-choice', 'multiple-choice'].includes(question.type)" 
              class="chart-container"
            >
              <div class="bar-chart">
                <div 
                  v-for="(option, optionId) in question.analytics.option_counts" 
                  :key="optionId"
                  class="chart-bar-container mb-3"
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
              <div class="scale-stats mb-4">
                <div class="row row-cols-2 row-cols-md-4 g-3">
                  <div class="col">
                    <div class="stat-box">
                      <span class="stat-title">Average</span>
                      <span class="stat-value">{{ question.analytics.average?.toFixed(2) }}</span>
                    </div>
                  </div>
                  <div class="col">
                    <div class="stat-box">
                      <span class="stat-title">Median</span>
                      <span class="stat-value">{{ question.analytics.median?.toFixed(1) }}</span>
                    </div>
                  </div>
                  <div class="col">
                    <div class="stat-box">
                      <span class="stat-title">Min</span>
                      <span class="stat-value">{{ question.analytics.min }}</span>
                    </div>
                  </div>
                  <div class="col">
                    <div class="stat-box">
                      <span class="stat-title">Max</span>
                      <span class="stat-value">{{ question.analytics.max }}</span>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Scale distribution visualization -->
              <div class="scale-distribution mb-3">
                <h5 class="mb-3">Distribution</h5>
                <div class="distribution-chart d-flex align-items-end mb-2">
                  <div 
                    v-for="(count, value) in getScaleDistribution(question.response_data)" 
                    :key="value" 
                    class="distribution-bar"
                    :style="{ 
                      height: `${getScaleBarHeight(count, question.response_data)}%`,
                      backgroundColor: getScaleBarColor(value)
                    }"
                    :title="`${value}: ${count} responses`"
                  ></div>
                </div>
                <div class="distribution-labels d-flex justify-content-between">
                  <div 
                    v-for="value in getScaleRange(question.analytics.min, question.analytics.max)" 
                    :key="value"
                    class="scale-label"
                  >{{ value }}</div>
                </div>
              </div>
            </div>

            <!-- Open Text Questions -->
            <div v-else-if="question.type === 'open-text'" class="open-text-results">
              <div class="text-responses mb-4">
                <h5 class="mb-3">Responses <span class="text-muted">({{ question.response_data.length }})</span></h5>
                <div class="card">
                  <ul class="list-group list-group-flush">
                    <li 
                      v-for="(response, index) in question.response_data.slice(0, showAllText ? undefined : 5)" 
                      :key="index"
                      class="list-group-item text-response"
                    >
                      {{ response }}
                    </li>
                  </ul>
                </div>
                <div class="mt-3 text-center">
                  <button 
                    v-if="!showAllText && question.response_data.length > 5" 
                    @click="showAllText = true" 
                    class="btn btn-sm btn-link"
                  >
                    <i class="bi bi-chevron-down me-1"></i> Show all {{ question.response_data.length }} responses
                  </button>
                  <button 
                    v-else-if="showAllText && question.response_data.length > 5"
                    @click="showAllText = false" 
                    class="btn btn-sm btn-link"
                  >
                    <i class="bi bi-chevron-up me-1"></i> Show less
                  </button>
                </div>
              </div>

              <!-- Word Frequency -->
              <div v-if="question.analytics.word_frequency" class="word-frequency">
                <h5 class="mb-3">Common Words</h5>
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

        <div class="actions d-flex justify-content-between mt-4">
          <router-link :to="`/surveys/${surveyId}/responses`" class="btn btn-primary">
            <i class="bi bi-people me-1"></i> View Individual Responses
          </router-link>
          <router-link :to="`/my-surveys`" class="btn btn-outline-secondary">
            <i class="bi bi-arrow-left me-1"></i> Back to My Surveys
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import surveyService from '@/services/survey'

export default {
  name: 'SurveyResultsPage',
  
  setup() {
    const route = useRoute()
    const router = useRouter()
    const surveyId = computed(() => route.params.id)

    // State
    const surveyResults = ref(null)
    const loading = ref(true)
    const error = ref(null)
    const showAllText = ref(false)

    // Fetch survey results
    const fetchResults = async () => {
      loading.value = true
      error.value = null
      
      try {
        surveyResults.value = await surveyService.getSurveyResults(surveyId.value)
      } catch (err) {
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
    const topWords = (wordFrequency) => {
      return Object.entries(wordFrequency)
        .sort((a, b) => b[1] - a[1])
        .slice(0, 15)
        .reduce((obj, [word, count]) => {
          obj[word] = count
          return obj
        }, {})
    }

    // Utility function to go back
    const goBack = () => {
      router.push('/my-surveys')
    }

    // Share survey function
    const shareSurvey = () => {
      const surveyLink = `${window.location.origin}/surveys/${surveyId.value}/take`
      
      if (navigator.clipboard) {
        navigator.clipboard.writeText(surveyLink)
          .then(() => {
            alert('Survey link copied to clipboard!')
          })
          .catch(() => {
            // Fallback if clipboard API fails
            prompt('Copy the survey link:', surveyLink)
          })
      } else {
        // Fallback for browsers without clipboard API
        prompt('Copy the survey link:', surveyLink)
      }
    }

    // Export to CSV function
    const exportToCsv = () => {
      if (!surveyResults.value) return
      
      try {
        // Generate CSV content
        let csvContent = "data:text/csv;charset=utf-8,"
        
        // Add headers
        csvContent += "Question,Type,Total Answers\n"
        
        // Add data rows
        surveyResults.value.question_results.forEach((question) => {
          csvContent += `"${question.title}","${question.type}",${question.total_answers}\n`
          
          // For choice questions, add option breakdown
          if (['single-choice', 'multiple-choice'].includes(question.type) && question.analytics.option_counts) {
            csvContent += "Option,Count,Percentage\n"
            
            for (const optionId in question.analytics.option_counts) {
              const option = question.analytics.option_counts[optionId]
              csvContent += `"${option.text}",${option.count},${option.percent.toFixed(2)}%\n`
            }
            
            csvContent += "\n" // Add empty line between questions
          }
          
          // For scale questions
          if (question.type === 'scale') {
            csvContent += "Statistic,Value\n"
            csvContent += `Average,${question.analytics.average?.toFixed(2)}\n`
            csvContent += `Median,${question.analytics.median?.toFixed(2)}\n`
            csvContent += `Min,${question.analytics.min}\n`
            csvContent += `Max,${question.analytics.max}\n\n`
          }
        })
        
        // Create download link
        const encodedUri = encodeURI(csvContent)
        const link = document.createElement("a")
        link.setAttribute("href", encodedUri)
        link.setAttribute("download", `survey_results_${surveyId.value}.csv`)
        document.body.appendChild(link)
        
        // Trigger download
        link.click()
        
        // Clean up
        document.body.removeChild(link)
      } catch (error) {
        console.error('Error exporting to CSV:', error)
        alert('Failed to export results to CSV')
      }
    }

    // Scale visualization helper functions
    const getScaleDistribution = (responseData) => {
      // Count occurrences of each value
      const distribution = {}
      
      responseData.forEach(response => {
        const value = Number(response)
        distribution[value] = (distribution[value] || 0) + 1
      })
      
      return distribution
    }

    const getScaleBarHeight = (count, responseData) => {
      // Calculate percentage height based on max count
      const maxCount = Math.max(...Object.values(getScaleDistribution(responseData)))
      return (count / maxCount) * 100
    }

    const getScaleBarColor = (value) => {
      // Color gradient based on value
      const colors = [
        '#FF4B4B', // 1 - red
        '#FFAB49', // 2 - orange
        '#FFD700', // 3 - yellow
        '#70C1B3', // 4 - teal
        '#30A0E0'  // 5 - blue
      ]
      
      return colors[Number(value) - 1] || '#30A0E0'
    }

    const getScaleRange = (min, max) => {
      // Generate array of numbers from min to max
      const range = []
      for (let i = min; i <= max; i++) {
        range.push(i)
      }
      return range
    }

    // Fetch results on mount
    onMounted(fetchResults)

    return {
      surveyResults,
      loading,
      error,
      showAllText,
      surveyId,
      fetchResults,
      topWords,
      goBack,
      shareSurvey,
      exportToCsv,
      getScaleDistribution,
      getScaleBarHeight,
      getScaleBarColor,
      getScaleRange
    }
  }
}
</script>

<style scoped>
.survey-results-page {
  padding: 20px 0;
}

.loading-container, 
.error-container {
  max-width: 600px;
  margin: 2rem auto;
  padding: 2rem;
}

.spinner {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.survey-stats {
  gap: 15px;
}

.stat-card {
  background-color: #f8f9fa;
  border-radius: 8px;
  padding: 15px 25px;
  text-align: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  min-width: 130px;
  transition: transform 0.3s, box-shadow 0.3s;
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--primary-color);
  margin-bottom: 5px;
}

.stat-label {
  font-size: 0.9rem;
  color: #6c757d;
}

.stat-box {
  padding: 15px;
  border-radius: 6px;
  text-align: center;
  background-color: #f8f9fa;
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.stat-title {
  font-size: 0.85rem;
  color: #6c757d;
  margin-bottom: 5px;
}

/* Chart styles */
.chart-bar-container {
  display: flex;
  align-items: center;
  margin-bottom: 15px;
}

.chart-label {
  width: 200px;
  margin-right: 15px;
  text-align: right;
  font-size: 0.95rem;
  color: #495057;
  font-weight: 500;
}

.chart-bar-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
}

.chart-bar {
  height: 30px;
  background-color: var(--primary-color);
  border-radius: 4px;
  min-width: 2%;
  transition: width 0.6s cubic-bezier(0.22, 0.61, 0.36, 1);
}

.chart-value {
  margin-left: 15px;
  font-size: 0.9rem;
  min-width: 90px;
  color: #495057;
}

.low-percent {
  min-width: 15px;
}

/* Scale distribution */
.scale-distribution {
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 8px;
}

.distribution-chart {
  height: 200px;
  gap: 10px;
}

.distribution-bar {
  flex: 1;
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
  transition: height 0.5s;
}

.distribution-labels {
  margin-top: 10px;
}

.scale-label {
  flex: 1;
  text-align: center;
  font-size: 0.9rem;
  color: #6c757d;
}

/* Text responses */
.text-response {
  padding: 15px;
  background-color: #fff;
  border-left: none !important;
  border-right: none !important;
}

.text-response:first-child {
  border-top: none !important;
}

/* Word cloud */
.word-cloud {
  padding: 25px;
  background-color: #f8f9fa;
  border-radius: 8px;
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  align-items: center;
  gap: 15px;
}

.word-item {
  display: inline-block;
  padding: 5px 10px;
  color: var(--primary-color);
  font-weight: 500;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .chart-label {
    width: 120px;
    font-size: 0.85rem;
  }
  
  .chart-bar {
    height: 25px;
  }
  
  .chart-value {
    font-size: 0.8rem;
    margin-left: 10px;
  }
  
  .distribution-chart {
    height: 150px;
  }
}
</style> 