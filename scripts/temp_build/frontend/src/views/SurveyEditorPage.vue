<template>
  <div class="survey-editor-page">
    <div class="container">
      <div class="page-header">
        <h1>{{ isEditMode ? 'Edit Survey' : 'Create Survey' }}</h1>
      </div>
      
      <div v-if="loading" class="loading">Loading survey data...</div>
      
      <div v-else-if="error" class="error-text">{{ error }}</div>
      
      <div v-else class="survey-editor">
        <form @submit.prevent="saveSurvey" class="survey-form">
          <!-- Survey Basic Info -->
          <div class="form-section">
            <h2>Survey Information</h2>
            
            <div class="form-group">
              <label for="survey-title">Title</label>
              <input 
                type="text" 
                id="survey-title" 
                v-model="survey.title" 
                required 
                placeholder="Enter survey title"
              />
            </div>
            
            <div class="form-group">
              <label for="survey-description">Description (optional)</label>
              <textarea 
                id="survey-description" 
                v-model="survey.description" 
                placeholder="Enter survey description"
                rows="3"
              ></textarea>
            </div>
          </div>
          
          <!-- Questions Section -->
          <div class="form-section">
            <div class="section-header">
              <h2>Questions</h2>
              <button 
                type="button" 
                @click="addQuestion" 
                class="btn-secondary"
              >
                Add Question
              </button>
            </div>
            
            <!-- No Questions Message -->
            <div v-if="survey.questions.length === 0" class="empty-questions">
              <p>No questions added yet. Click "Add Question" to start building your survey.</p>
            </div>
            
            <!-- Questions List -->
            <div v-else class="questions-list">
              <div 
                v-for="(question, index) in survey.questions" 
                :key="question.question_id" 
                class="question-item"
              >
                <div class="question-header">
                  <div class="question-number">Question {{ index + 1 }}</div>
                  <div class="question-actions">
                    <button 
                      type="button" 
                      @click="moveQuestion(index, -1)" 
                      :disabled="index === 0"
                      class="btn-icon"
                      title="Move up"
                    >
                      ↑
                    </button>
                    <button 
                      type="button" 
                      @click="moveQuestion(index, 1)" 
                      :disabled="index === survey.questions.length - 1"
                      class="btn-icon"
                      title="Move down"
                    >
                      ↓
                    </button>
                    <button 
                      type="button" 
                      @click="removeQuestion(index)" 
                      class="btn-icon btn-delete-icon"
                      title="Delete question"
                    >
                      ×
                    </button>
                  </div>
                </div>
                
                <div class="question-body">
                  <!-- Question Text -->
                  <div class="form-group">
                    <label :for="`question-text-${index}`">Question Text</label>
                    <input 
                      type="text" 
                      :id="`question-text-${index}`" 
                      v-model="question.text" 
                      required 
                      placeholder="Enter question text"
                    />
                  </div>
                  
                  <!-- Question Type -->
                  <div class="form-group">
                    <label :for="`question-type-${index}`">Question Type</label>
                    <select 
                      :id="`question-type-${index}`" 
                      v-model="question.type" 
                      @change="handleQuestionTypeChange(question)"
                    >
                      <option value="single-choice">Single Choice</option>
                      <option value="multiple-choice">Multiple Choice</option>
                      <option value="open-text">Open Text</option>
                      <option value="scale">Scale</option>
                      <option value="matrix-single">Matrix (Single Choice)</option>
                      <option value="matrix-multiple">Matrix (Multiple Choice)</option>
                    </select>
                  </div>
                  
                  <!-- Required -->
                  <div class="form-group checkbox-group">
                    <input 
                      type="checkbox" 
                      :id="`question-required-${index}`" 
                      v-model="question.is_required" 
                    />
                    <label :for="`question-required-${index}`">Required</label>
                  </div>
                  
                  <!-- Answer Options for choice questions -->
                  <div 
                    v-if="['single-choice', 'multiple-choice'].includes(question.type)" 
                    class="answer-options"
                  >
                    <h4>Answer Options</h4>
                    
                    <div 
                      v-for="(option, optionIndex) in question.options || []" 
                      :key="optionIndex" 
                      class="option-item"
                    >
                      <input 
                        type="text" 
                        v-model="option.text" 
                        :placeholder="`Option ${optionIndex + 1}`"
                        required
                      />
                      <button 
                        type="button" 
                        @click="removeOption(question, optionIndex)" 
                        class="btn-icon btn-delete-icon"
                        title="Delete option"
                      >
                        ×
                      </button>
                    </div>
                    
                    <button 
                      type="button" 
                      @click="addOption(question)" 
                      class="btn-small"
                    >
                      Add Option
                    </button>
                  </div>
                  
                  <!-- Scale settings -->
                  <div v-if="question.type === 'scale'" class="scale-settings">
                    <h4>Scale Settings</h4>
                    
                    <div class="scale-inputs">
                      <div class="form-group">
                        <label :for="`scale-min-${index}`">Min Value</label>
                        <input 
                          type="number" 
                          :id="`scale-min-${index}`" 
                          v-model.number="question.scale_settings.min" 
                          required 
                          min="0"
                          max="10"
                        />
                      </div>
                      
                      <div class="form-group">
                        <label :for="`scale-max-${index}`">Max Value</label>
                        <input 
                          type="number" 
                          :id="`scale-max-${index}`" 
                          v-model.number="question.scale_settings.max" 
                          required 
                          min="1"
                          max="10"
                        />
                      </div>
                    </div>
                    
                    <div class="scale-labels">
                      <div class="form-group">
                        <label :for="`scale-min-label-${index}`">Min Label</label>
                        <input 
                          type="text" 
                          :id="`scale-min-label-${index}`" 
                          v-model="question.scale_settings.minLabel" 
                          placeholder="e.g., Not at all likely"
                        />
                      </div>
                      
                      <div class="form-group">
                        <label :for="`scale-max-label-${index}`">Max Label</label>
                        <input 
                          type="text" 
                          :id="`scale-max-label-${index}`" 
                          v-model="question.scale_settings.maxLabel" 
                          placeholder="e.g., Extremely likely"
                        />
                      </div>
                    </div>
                  </div>
                  
                  <!-- Matrix settings -->
                  <div v-if="['matrix-single', 'matrix-multiple'].includes(question.type)" class="matrix-settings">
                    <h4>Matrix Settings</h4>
                    
                    <div class="matrix-rows">
                      <h5>Rows</h5>
                      <div 
                        v-for="(row, rowIndex) in question.matrix_rows || []" 
                        :key="`row-${rowIndex}`" 
                        class="matrix-item"
                      >
                        <input 
                          type="text" 
                          v-model="question.matrix_rows[rowIndex]" 
                          :placeholder="`Row ${rowIndex + 1}`"
                          required
                        />
                        <button 
                          type="button" 
                          @click="removeMatrixRow(question, rowIndex)" 
                          class="btn-icon btn-delete-icon"
                          title="Delete row"
                        >
                          ×
                        </button>
                      </div>
                      
                      <button 
                        type="button" 
                        @click="addMatrixRow(question)" 
                        class="btn-small"
                      >
                        Add Row
                      </button>
                    </div>
                    
                    <div class="matrix-columns">
                      <h5>Columns</h5>
                      <div 
                        v-for="(column, columnIndex) in question.matrix_columns || []" 
                        :key="`column-${columnIndex}`" 
                        class="matrix-item"
                      >
                        <input 
                          type="text" 
                          v-model="question.matrix_columns[columnIndex]" 
                          :placeholder="`Column ${columnIndex + 1}`"
                          required
                        />
                        <button 
                          type="button" 
                          @click="removeMatrixColumn(question, columnIndex)" 
                          class="btn-icon btn-delete-icon"
                          title="Delete column"
                        >
                          ×
                        </button>
                      </div>
                      
                      <button 
                        type="button" 
                        @click="addMatrixColumn(question)" 
                        class="btn-small"
                      >
                        Add Column
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Form Actions -->
          <div class="form-actions">
            <button type="submit" class="btn" :disabled="saving">
              {{ saving ? 'Saving...' : 'Save Survey' }}
            </button>
            <button type="button" @click="cancel" class="btn-secondary">
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import surveyService from '@/services/survey'
import { 
  Survey, 
  Question, 
  CreateSurveyRequest, 
  QuestionType,
  Option,
  ScaleSettings
} from '@/types/survey'

// Simple ID generator
const generateId = () => {
  return Math.random().toString(36).substring(2, 15) + 
         Math.random().toString(36).substring(2, 15)
}

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// State
const loading = ref(false)
const saving = ref(false)
const error = ref('')
const survey = ref<CreateSurveyRequest>({
  title: '',
  description: '',
  questions: []
})

// Computed
const isEditMode = computed(() => !!route.params.id)
const surveyId = computed(() => route.params.id as string)

// Initialize
onMounted(async () => {
  if (isEditMode.value) {
    await loadSurvey()
  }
})

// Load existing survey
const loadSurvey = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const existingSurvey = await surveyService.getSurvey(surveyId.value)
    survey.value = {
      title: existingSurvey.title,
      description: existingSurvey.description || '',
      questions: existingSurvey.questions
    }
  } catch (e: any) {
    error.value = e.message || 'Failed to load survey'
  } finally {
    loading.value = false
  }
}

// Save survey
const saveSurvey = async () => {
  if (!authStore.user) {
    error.value = 'You must be logged in to save a survey'
    return
  }
  
  saving.value = true
  error.value = ''
  
  try {
    if (isEditMode.value) {
      await surveyService.updateSurvey(surveyId.value, survey.value)
    } else {
      await surveyService.createSurvey(survey.value)
    }
    
    router.push('/my-surveys')
  } catch (e: any) {
    error.value = e.message || 'Failed to save survey'
    saving.value = false
  }
}

// Add a new question
const addQuestion = () => {
  const newQuestion: Question = {
    question_id: generateId(),
    text: '',
    type: QuestionType.SingleChoice,
    is_required: false,
    options: [
      { value: generateId(), text: '' },
      { value: generateId(), text: '' }
    ]
  }
  
  survey.value.questions.push(newQuestion)
}

// Remove a question
const removeQuestion = (index: number) => {
  survey.value.questions.splice(index, 1)
}

// Move a question up or down
const moveQuestion = (index: number, direction: number) => {
  const newIndex = index + direction
  
  if (newIndex < 0 || newIndex >= survey.value.questions.length) {
    return
  }
  
  const questions = [...survey.value.questions]
  const temp = questions[index]
  questions[index] = questions[newIndex]
  questions[newIndex] = temp
  
  survey.value.questions = questions
}

// Handle question type change
const handleQuestionTypeChange = (question: Question) => {
  // Reset options and settings based on question type
  switch (question.type) {
    case QuestionType.SingleChoice:
    case QuestionType.MultipleChoice:
      question.options = question.options?.length ? question.options : [
        { value: generateId(), text: '' },
        { value: generateId(), text: '' }
      ]
      // Clear other settings
      delete question.scale_settings
      delete question.matrix_rows
      delete question.matrix_columns
      break
      
    case QuestionType.OpenText:
      // Clear all options and settings
      delete question.options
      delete question.scale_settings
      delete question.matrix_rows
      delete question.matrix_columns
      break
      
    case QuestionType.Scale:
      // Set default scale settings
      question.scale_settings = {
        min: 0,
        max: 10,
        minLabel: '',
        maxLabel: ''
      }
      // Clear other options
      delete question.options
      delete question.matrix_rows
      delete question.matrix_columns
      break
      
    case QuestionType.MatrixSingle:
    case QuestionType.MatrixMultiple:
      // Set default matrix settings
      question.matrix_rows = question.matrix_rows?.length ? question.matrix_rows : ['', '']
      question.matrix_columns = question.matrix_columns?.length ? question.matrix_columns : ['', '']
      // Clear other options
      delete question.options
      delete question.scale_settings
      break
  }
}

// Add option to a question
const addOption = (question: Question) => {
  if (!question.options) {
    question.options = []
  }
  
  question.options.push({ value: generateId(), text: '' })
}

// Remove option from a question
const removeOption = (question: Question, index: number) => {
  if (question.options && question.options.length > 1) {
    question.options.splice(index, 1)
  }
}

// Add matrix row
const addMatrixRow = (question: Question) => {
  if (!question.matrix_rows) {
    question.matrix_rows = []
  }
  
  question.matrix_rows.push('')
}

// Remove matrix row
const removeMatrixRow = (question: Question, index: number) => {
  if (question.matrix_rows && question.matrix_rows.length > 1) {
    question.matrix_rows.splice(index, 1)
  }
}

// Add matrix column
const addMatrixColumn = (question: Question) => {
  if (!question.matrix_columns) {
    question.matrix_columns = []
  }
  
  question.matrix_columns.push('')
}

// Remove matrix column
const removeMatrixColumn = (question: Question, index: number) => {
  if (question.matrix_columns && question.matrix_columns.length > 1) {
    question.matrix_columns.splice(index, 1)
  }
}

// Cancel editing
const cancel = () => {
  router.push('/my-surveys')
}
</script>

<style scoped>
.survey-editor-page {
  padding: 2rem 0;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  color: var(--primary-color);
}

.loading, .error-text {
  text-align: center;
  padding: 2rem;
}

.survey-editor {
  background-color: var(--white);
  border-radius: var(--border-radius);
  box-shadow: var(--box-shadow);
  padding: 2rem;
}

.form-section {
  margin-bottom: 2rem;
  padding-bottom: 2rem;
  border-bottom: 1px solid var(--light-gray);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.empty-questions {
  text-align: center;
  padding: 2rem;
  background-color: var(--light-gray);
  border-radius: var(--border-radius);
}

.questions-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.question-item {
  background-color: var(--light-gray);
  border-radius: var(--border-radius);
  padding: 1.5rem;
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.question-number {
  font-weight: bold;
  font-size: 1.1rem;
}

.question-actions {
  display: flex;
  gap: 0.5rem;
}

.question-body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.checkbox-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.checkbox-group input {
  width: auto;
}

.answer-options, .scale-settings, .matrix-settings {
  background-color: var(--white);
  padding: 1rem;
  border-radius: var(--border-radius);
  margin-top: 0.5rem;
}

.answer-options h4, .scale-settings h4, .matrix-settings h4 {
  margin-bottom: 1rem;
  font-size: 1rem;
}

.matrix-settings h5 {
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
}

.matrix-rows, .matrix-columns {
  margin-bottom: 1.5rem;
}

.option-item, .matrix-item {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.5rem;
}

.scale-inputs, .scale-labels {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  margin-bottom: 1rem;
}

.btn-secondary {
  background-color: var(--white);
  color: var(--primary-color);
  border: 1px solid var(--primary-color);
}

.btn-secondary:hover {
  background-color: var(--light-gray);
}

.btn-small {
  font-size: 0.9rem;
  padding: 0.3rem 0.6rem;
}

.btn-icon {
  width: 30px;
  height: 30px;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--white);
  border: 1px solid var(--gray);
  border-radius: var(--border-radius);
  font-size: 1.2rem;
  color: var(--dark-gray);
}

.btn-icon:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-delete-icon {
  color: var(--error-color);
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 1.5rem;
}
</style> 