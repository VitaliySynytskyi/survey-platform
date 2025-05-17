<template>
  <div class="survey-take-page">
    <div class="container py-4">
      <div v-if="loading" class="loading-container text-center py-5">
        <div class="spinner">
          <div class="spinner-border text-primary" role="status">
            <span class="visually-hidden">Loading...</span>
          </div>
        </div>
        <p class="mt-3 text-muted">Loading survey...</p>
      </div>

      <div v-else-if="error" class="error-container alert alert-danger my-4 shadow-sm" role="alert">
        <h4 class="alert-heading mb-2">Error loading survey!</h4>
        <p class="mb-0">{{ error }}</p>
        <div class="mt-3">
          <button @click="fetchSurvey" class="btn btn-outline-danger">Try Again</button>
          <button @click="goToHome" class="btn btn-outline-secondary ms-2">Return Home</button>
        </div>
      </div>

      <div v-else-if="submitted" class="submitted-container text-center py-5 my-5 bg-white rounded shadow">
        <div class="submitted-icon mb-4">
          <i class="bi bi-check-circle-fill text-success" style="font-size: 5rem;"></i>
        </div>
        <h2 class="mb-3">Thank you for your response!</h2>
        <p class="mb-4 text-secondary lead">Your answers have been submitted successfully.</p>
        <div class="d-flex justify-content-center gap-3">
          <button @click="goToHome" class="btn btn-primary px-4 py-2">Return to Home</button>
          <button @click="takeSurveyAgain" class="btn btn-outline-primary px-4 py-2">Take Survey Again</button>
        </div>
      </div>

      <div v-else-if="survey" class="survey-container bg-white rounded shadow p-4 p-md-5">
        <div class="survey-header mb-5 pb-3 border-bottom">
          <h1 class="survey-title mb-3">{{ survey.title }}</h1>
          <p v-if="survey.description" class="survey-description lead">{{ survey.description }}</p>
        </div>

        <!-- Progress indicator -->
        <div class="progress-container mb-4" v-if="survey.questions && survey.questions.length > 0">
          <div class="d-flex justify-content-between mb-2">
            <span class="text-muted">Question {{ currentQuestionNumber }} of {{ survey.questions.length }}</span>
            <span class="text-muted">{{ Math.round((currentQuestionNumber / survey.questions.length) * 100) }}% completed</span>
          </div>
          <div class="progress" style="height: 8px;">
            <div 
              class="progress-bar" 
              :style="{ width: `${(currentQuestionNumber / survey.questions.length) * 100}%` }"
            ></div>
          </div>
        </div>

        <form @submit.prevent="submitSurvey" class="survey-form">
          <div v-for="(question, index) in sortedQuestions" :key="question.id" 
              class="question-container mb-4 p-4 border rounded"
              :class="{'active-question': index === currentQuestionIndex}">
            <div class="question-header mb-3">
              <h3 class="question-title d-flex align-items-center">
                {{ question.title }}
                <span v-if="question.required" class="text-danger ms-2">*</span>
              </h3>
              <p v-if="question.description" class="question-description text-muted mt-2">{{ question.description }}</p>
            </div>

            <!-- Single Choice Question -->
            <div v-if="question.type === 'single-choice'" class="question-options">
              <div v-for="option in sortOptions(question.options)" :key="option.id" class="form-check mb-3">
                <input
                  :id="`option-${option.id}`"
                  type="radio"
                  class="form-check-input"
                  :name="`question-${question.id}`"
                  :value="option.id"
                  v-model="answers[question.id]"
                  :required="question.required"
                >
                <label :for="`option-${option.id}`" class="form-check-label">{{ option.text }}</label>
              </div>
            </div>

            <!-- Multiple Choice Question -->
            <div v-else-if="question.type === 'multiple-choice'" class="question-options">
              <div v-for="option in sortOptions(question.options)" :key="option.id" class="form-check mb-3">
                <input
                  :id="`option-${option.id}`"
                  type="checkbox"
                  class="form-check-input"
                  :value="option.id"
                  v-model="multipleChoiceAnswers[question.id]"
                  @change="validateRequired(question)"
                >
                <label :for="`option-${option.id}`" class="form-check-label">{{ option.text }}</label>
              </div>
              <div v-if="question.required && validationErrors[question.id]" class="text-danger mt-2 small">
                <i class="bi bi-exclamation-circle me-1"></i>Please select at least one option
              </div>
            </div>

            <!-- Open Text Question -->
            <div v-else-if="question.type === 'open-text'" class="question-text">
              <textarea
                :id="`question-${question.id}`"
                class="form-control"
                rows="3"
                v-model="answers[question.id]"
                :required="question.required"
                placeholder="Your answer..."
              ></textarea>
            </div>

            <!-- Scale Question -->
            <div v-else-if="question.type === 'scale'" class="question-scale">
              <div class="scale-container d-flex justify-content-between align-items-center mb-3">
                <span class="scale-label">Low</span>
                <div class="scale-buttons">
                  <span v-for="value in 5" :key="value" class="scale-point">
                    <input
                      :id="`scale-${question.id}-${value}`"
                      type="radio"
                      class="btn-check"
                      :name="`question-${question.id}`"
                      :value="value.toString()"
                      v-model="answers[question.id]"
                      :required="question.required"
                    >
                    <label :for="`scale-${question.id}-${value}`" class="btn btn-outline-primary scale-btn">{{ value }}</label>
                  </span>
                </div>
                <span class="scale-label">High</span>
              </div>
            </div>

            <!-- Navigation buttons for each question -->
            <div class="question-navigation mt-4 d-flex justify-content-between">
              <button v-if="index > 0" 
                type="button" 
                class="btn btn-outline-secondary" 
                @click="currentQuestionIndex = index - 1">
                <i class="bi bi-chevron-left me-2"></i>Previous
              </button>
              <div class="flex-grow-1"></div>
              <button v-if="index < sortedQuestions.length - 1" 
                type="button" 
                class="btn btn-outline-primary" 
                @click="currentQuestionIndex = index + 1">
                Next<i class="bi bi-chevron-right ms-2"></i>
              </button>
            </div>
          </div>

          <div class="survey-actions mt-5 d-flex justify-content-center">
            <button type="submit" class="btn btn-primary btn-lg px-5 py-3" :disabled="submitting">
              <span v-if="submitting" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
              {{ submitting ? 'Submitting...' : 'Submit Response' }}
            </button>
          </div>
        </form>
      </div>

      <div v-else class="alert alert-warning text-center my-5 shadow-sm" role="alert">
        <i class="bi bi-exclamation-triangle-fill me-2"></i>
        <strong>Survey not found or has been removed.</strong>
        <div class="mt-3">
          <button @click="goToHome" class="btn btn-warning">Return to Home</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import surveyService from '@/services/survey'

export default {
  name: 'SurveyTakePage',
  
  setup() {
    const route = useRoute()
    const router = useRouter()
    const surveyId = route.params.id
    
    // State
    const survey = ref(null)
    const loading = ref(true)
    const error = ref(null)
    const submitting = ref(false)
    const submitted = ref(false)
    const answers = reactive({})
    const multipleChoiceAnswers = reactive({})
    const validationErrors = reactive({})
    const currentQuestionIndex = ref(0)
    
    // Computed properties
    const sortedQuestions = computed(() => {
      if (!survey.value || !survey.value.questions) return []
      return [...survey.value.questions].sort((a, b) => a.order - b.order)
    })
    
    const currentQuestionNumber = computed(() => {
      return currentQuestionIndex.value + 1
    })
    
    // Methods
    const fetchSurvey = async () => {
      loading.value = true
      error.value = null
      
      try {
        const response = await surveyService.getPublicSurvey(surveyId)
        survey.value = response
        
        // Initialize multiple choice answers
        survey.value.questions.forEach(question => {
          if (question.type === 'multiple-choice') {
            multipleChoiceAnswers[question.id] = []
          }
        })
      } catch (err) {
        console.error('Error fetching survey:', err)
        error.value = 'Unable to load survey. Please try again later.'
      } finally {
        loading.value = false
      }
    }
    
    const sortOptions = (options) => {
      if (!options) return []
      return [...options].sort((a, b) => a.order - b.order)
    }
    
    const validateRequired = (question) => {
      if (question.required && question.type === 'multiple-choice') {
        validationErrors[question.id] = multipleChoiceAnswers[question.id].length === 0
      }
    }
    
    const validateAllRequired = () => {
      let isValid = true
      
      survey.value.questions.forEach(question => {
        if (question.required) {
          if (question.type === 'multiple-choice') {
            const isAnswered = multipleChoiceAnswers[question.id] && multipleChoiceAnswers[question.id].length > 0
            validationErrors[question.id] = !isAnswered
            if (!isAnswered) isValid = false
          }
        }
      })
      
      return isValid
    }
    
    const prepareResponses = () => {
      const response = {
        survey_id: surveyId,
        anonymous_id: generateAnonymousId(), // Generate a random ID for anonymous users
        answers: []
      }
      
      // Process single-choice, open-text, and scale answers
      for (const questionId in answers) {
        response.answers.push({
          question_id: questionId,
          value: answers[questionId]
        })
      }
      
      // Process multiple-choice answers
      for (const questionId in multipleChoiceAnswers) {
        if (multipleChoiceAnswers[questionId].length > 0) {
          response.answers.push({
            question_id: questionId,
            value: multipleChoiceAnswers[questionId]
          })
        }
      }
      
      return response
    }
    
    const generateAnonymousId = () => {
      return 'anon_' + Math.random().toString(36).substring(2, 15)
    }
    
    const submitSurvey = async () => {
      if (!validateAllRequired()) {
        alert('Please answer all required questions.')
        return
      }
      
      submitting.value = true
      
      try {
        const responseData = prepareResponses()
        await surveyService.submitSurveyResponse(surveyId, responseData)
        submitted.value = true
      } catch (err) {
        console.error('Error submitting response:', err)
        alert('Failed to submit your response. Please try again.')
      } finally {
        submitting.value = false
      }
    }
    
    const goToHome = () => {
      router.push('/')
    }
    
    const takeSurveyAgain = () => {
      // Reset everything
      submitted.value = false
      currentQuestionIndex.value = 0
      
      // Clear all answers
      for (const key in answers) delete answers[key]
      
      // Reset multiple choice answers
      for (const key in multipleChoiceAnswers) {
        multipleChoiceAnswers[key] = []
      }
      
      // Clear validation errors
      for (const key in validationErrors) delete validationErrors[key]
    }
    
    // Lifecycle hooks
    onMounted(() => {
      fetchSurvey()
    })
    
    return {
      survey,
      loading,
      error,
      submitting,
      submitted,
      answers,
      multipleChoiceAnswers,
      validationErrors,
      sortedQuestions,
      currentQuestionNumber,
      currentQuestionIndex,
      sortOptions,
      validateRequired,
      submitSurvey,
      goToHome,
      takeSurveyAgain,
      fetchSurvey
    }
  }
}
</script>

<style scoped>
.survey-container {
  max-width: 800px;
  margin: 0 auto;
}

.loading-container, 
.error-container {
  max-width: 600px;
  margin: 2rem auto;
  padding: 2rem;
}

.question-container {
  background-color: #f9f9f9;
  transition: transform 0.2s, box-shadow 0.2s;
  border-radius: 8px;
  border-color: #eaeaea !important;
  position: relative;
  display: none;
}

.question-container.active-question {
  display: block;
}

.question-container:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 15px rgba(0, 0, 0, 0.08);
}

.scale-container {
  margin-top: 15px;
  padding: 10px 0;
}

.scale-buttons {
  display: flex;
  gap: 10px;
}

.scale-btn {
  border-radius: 50%;
  width: 50px;
  height: 50px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-weight: bold;
}

.scale-label {
  color: #6c757d;
  font-size: 0.9rem;
  width: 60px;
}

.submitted-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 3rem;
}

.form-check {
  padding: 10px;
  border-radius: 5px;
  margin-bottom: 10px;
  transition: background-color 0.2s;
}

.form-check:hover {
  background-color: #f0f0f0;
}

.form-check-input:checked + .form-check-label {
  font-weight: 500;
  color: #0d6efd;
}

.form-check-label {
  cursor: pointer;
  padding-left: 0.5rem;
}

.spinner {
  display: flex;
  justify-content: center;
  margin-bottom: 20px;
}

.progress {
  height: 8px;
  background-color: #e9ecef;
  border-radius: 1rem;
  overflow: hidden;
}

.progress-bar {
  background-color: #0d6efd;
  transition: width 0.6s ease;
}

.question-title {
  color: #333;
  margin-bottom: 5px;
}

.question-description {
  font-size: 0.9rem;
}

.question-navigation {
  margin-top: 20px;
  padding-top: 15px;
  border-top: 1px solid #eee;
}

@media (max-width: 768px) {
  .scale-btn {
    width: 40px;
    height: 40px;
    font-size: 0.9rem;
  }
  
  .question-title {
    font-size: 1.25rem;
  }
  
  .scale-buttons {
    gap: 5px;
  }
  
  .scale-label {
    width: 40px;
    font-size: 0.8rem;
  }
}
</style> 