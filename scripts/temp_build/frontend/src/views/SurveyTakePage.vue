<template>
  <div class="survey-take-page">
    <div class="container py-5">
      <div v-if="loading" class="text-center py-5">
        <div class="spinner-border" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="mt-3">Loading survey...</p>
      </div>

      <div v-else-if="error" class="alert alert-danger" role="alert">
        <h4 class="alert-heading">Error loading survey!</h4>
        <p>{{ error }}</p>
      </div>

      <div v-else-if="submitted" class="submitted-container text-center py-5">
        <div class="submitted-icon mb-4">
          <i class="bi bi-check-circle-fill text-success" style="font-size: 4rem;"></i>
        </div>
        <h2 class="mb-3">Thank you for your response!</h2>
        <p class="mb-4">Your answers have been submitted successfully.</p>
        <button @click="goToHome" class="btn btn-primary">Return to Home</button>
      </div>

      <div v-else-if="survey" class="survey-container">
        <div class="survey-header mb-4">
          <h1 class="survey-title">{{ survey.title }}</h1>
          <p v-if="survey.description" class="survey-description">{{ survey.description }}</p>
        </div>

        <form @submit.prevent="submitSurvey" class="survey-form">
          <div v-for="question in sortedQuestions" :key="question.id" class="question-container mb-4 p-4 border rounded">
            <div class="question-header mb-3">
              <h3 class="question-title">
                {{ question.title }}
                <span v-if="question.required" class="text-danger">*</span>
              </h3>
              <p v-if="question.description" class="question-description text-muted">{{ question.description }}</p>
            </div>

            <!-- Single Choice Question -->
            <div v-if="question.type === 'single-choice'" class="question-options">
              <div v-for="option in sortOptions(question.options)" :key="option.id" class="form-check mb-2">
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
              <div v-for="option in sortOptions(question.options)" :key="option.id" class="form-check mb-2">
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
              <div v-if="question.required && validationErrors[question.id]" class="text-danger">
                Please select at least one option
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
              <div class="scale-container d-flex justify-content-between mb-2">
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
                  <label :for="`scale-${question.id}-${value}`" class="btn btn-outline-primary">{{ value }}</label>
                </span>
              </div>
              <div class="scale-labels d-flex justify-content-between">
                <small>Low</small>
                <small>High</small>
              </div>
            </div>
          </div>

          <div class="survey-actions mt-5 text-center">
            <button type="submit" class="btn btn-primary btn-lg" :disabled="submitting">
              <span v-if="submitting" class="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
              {{ submitting ? 'Submitting...' : 'Submit Response' }}
            </button>
          </div>
        </form>
      </div>

      <div v-else class="alert alert-warning text-center" role="alert">
        Survey not found or has been removed.
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
    
    // Computed properties
    const sortedQuestions = computed(() => {
      if (!survey.value || !survey.value.questions) return []
      return [...survey.value.questions].sort((a, b) => a.order - b.order)
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
      sortOptions,
      validateRequired,
      submitSurvey,
      goToHome
    }
  }
}
</script>

<style scoped>
.survey-container {
  max-width: 800px;
  margin: 0 auto;
}

.question-container {
  background-color: #f9f9f9;
  transition: transform 0.2s, box-shadow 0.2s;
}

.question-container:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.scale-container {
  margin-top: 15px;
  padding: 10px 0;
}

.scale-point {
  text-align: center;
  flex: 1;
}

.submitted-container {
  max-width: 600px;
  margin: 0 auto;
}
</style> 