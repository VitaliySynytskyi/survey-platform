<template>
  <div class="take-survey">
    <v-row v-if="loading">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <div class="mt-4">Loading survey...</div>
      </v-col>
    </v-row>

    <v-row v-else-if="error">
      <v-col cols="12">
        <v-alert type="error" closable>
          {{ error }}
        </v-alert>
      </v-col>
    </v-row>

    <template v-else-if="!submitted">
      <v-row>
        <v-col cols="12">
          <v-card class="mb-4 pa-6">
            <v-card-title class="text-h4 mb-2">{{ survey.title }}</v-card-title>
            <v-card-subtitle v-if="survey.created_by" class="mb-4">
              Created by {{ survey.created_by }}
            </v-card-subtitle>
            <v-card-text>
              <p class="text-body-1">{{ survey.description }}</p>
              <v-alert v-if="survey.is_active === false" type="warning" density="compact" class="mt-4">
                This survey is currently inactive.
              </v-alert>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12" md="8" lg="7" class="mx-auto">
          <v-form @submit.prevent="submitSurvey" ref="form" v-model="valid" :disabled="survey.is_active === false">
            <v-card class="mb-4 pa-6" v-for="(question, index) in survey.questions" :key="index">
              <div class="d-flex align-center mb-2">
                <span class="text-h6">{{ index + 1 }}. {{ question.text }}</span>
                <v-chip v-if="question.required" color="primary" size="small" class="ml-2">Required</v-chip>
              </div>

              <!-- Short Answer -->
              <v-text-field
                v-if="question.type === 'short_answer'"
                v-model="responses[question.id]"
                label="Your answer"
                variant="outlined"
                density="comfortable"
                :rules="question.required ? [v => !!v || 'This field is required'] : []"
              ></v-text-field>

              <!-- Paragraph -->
              <v-textarea
                v-else-if="question.type === 'paragraph'"
                v-model="responses[question.id]"
                label="Your answer"
                variant="outlined"
                rows="3"
                :rules="question.required ? [v => !!v || 'This field is required'] : []"
              ></v-textarea>

              <!-- Multiple Choice -->
              <v-radio-group
                v-else-if="question.type === 'multiple_choice'"
                v-model="responses[question.id]"
                :rules="question.required ? [v => !!v || 'This field is required'] : []"
              >
                <v-radio
                  v-for="(option, optIndex) in question.options"
                  :key="optIndex"
                  :label="option.text"
                  :value="option.text"
                ></v-radio>
              </v-radio-group>

              <!-- Checkbox -->
              <template v-else-if="question.type === 'checkbox'">
                <v-checkbox
                  v-for="(option, optIndex) in question.options"
                  :key="optIndex"
                  :model-value="isOptionSelected(question.id, option.text)"
                  @update:model-value="updateCheckboxValue(question.id, option.text, $event)"
                  :label="option.text"
                  density="comfortable"
                  hide-details
                  class="mb-1"
                ></v-checkbox>
                <div 
                  v-if="question.required && !getCheckboxArray(question.id).length" 
                  class="text-error text-caption mt-2"
                >
                  Please select at least one option
                </div>
              </template>

              <!-- Dropdown -->
              <v-select
                v-else-if="question.type === 'dropdown'"
                v-model="responses[question.id]"
                :items="question.options"
                item-title="text"
                item-value="text"
                label="Select an option"
                variant="outlined"
                density="comfortable"
                :rules="question.required ? [v => !!v || 'This field is required'] : []"
              ></v-select>

              <!-- Linear Scale -->
              <div v-else-if="question.type === 'linear_scale'" class="mt-4">
                <v-slider
                  v-model="responses[question.id]"
                  min="1"
                  max="5"
                  :step="1"
                  :ticks="true"
                  :tick-labels="['1', '2', '3', '4', '5']"
                  show-ticks="always"
                  thumb-label="always"
                  :rules="question.required ? [v => !!v || 'This field is required'] : []"
                ></v-slider>
                <div class="d-flex justify-space-between text-caption mt-n4">
                  <span>Not at all</span>
                  <span>Very much</span>
                </div>
              </div>

              <!-- Date -->
              <v-text-field
                v-else-if="question.type === 'date'"
                v-model="responses[question.id]"
                label="Your answer"
                type="date"
                variant="outlined"
                density="comfortable"
                :rules="question.required ? [v => !!v || 'This field is required'] : []"
              ></v-text-field>
            </v-card>

            <v-card class="pa-6 mb-4">
              <div class="d-flex justify-end">
                <v-btn
                  type="submit"
                  color="primary"
                  size="large"
                  :loading="submitting"
                  :disabled="!valid || submitting || survey.is_active === false"
                >
                  Submit Response
                </v-btn>
              </div>
            </v-card>
          </v-form>
        </v-col>
      </v-row>
    </template>

    <v-row v-else>
      <v-col cols="12" md="8" lg="6" class="mx-auto">
        <v-card class="pa-6 text-center">
          <v-card-title class="text-h4 mb-4">Thank You!</v-card-title>
          <v-card-text>
            <v-icon color="success" size="64" class="mb-4">mdi-check-circle</v-icon>
            <p class="text-h6 mb-6">Your response has been recorded.</p>
            <div class="d-flex justify-center">
              <v-btn color="primary" to="/">Return to Home</v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import { surveyApi } from '../services/api';

export default {
  name: 'TakeSurvey',
  props: {
    id: {
      type: String,
      required: true
    }
  },
  setup(props) {
    const route = useRoute();
    const router = useRouter();
    const authStore = useAuthStore();
    const form = ref(null);
    const valid = ref(true);
    const loading = ref(true);
    const submitting = ref(false);
    const error = ref('');
    const submitted = ref(false);
    
    const survey = ref({
      id: '',
      title: '',
      description: '',
      created_by: '',
      is_active: true, // Assume active by default until fetched
      questions: []
    });
    
    const responses = reactive({});
    const checkboxResponses = reactive({});
    
    // Mock survey data
    const mockSurvey = {
      id: '1',
      title: 'Customer Satisfaction Survey',
      description: 'We value your feedback! Please take a few minutes to tell us about your experience with our product and services.',
      created_by: 'Admin User',
      questions: [
        {
          id: 'q1',
          text: 'How satisfied are you with our product?',
          type: 'linear_scale',
          required: true,
          options: []
        },
        {
          id: 'q2',
          text: 'What features do you like most?',
          type: 'checkbox',
          required: true,
          options: ['User Interface', 'Performance', 'Reliability', 'Customer Support', 'Price']
        },
        {
          id: 'q3',
          text: 'How did you hear about us?',
          type: 'multiple_choice',
          required: true,
          options: ['Search Engine', 'Social Media', 'Friend/Colleague', 'Advertisement', 'Other']
        },
        {
          id: 'q4',
          text: 'When did you start using our product?',
          type: 'date',
          required: false,
          options: []
        },
        {
          id: 'q5',
          text: 'What could we improve?',
          type: 'paragraph',
          required: false,
          options: []
        }
      ]
    };
    
    const fetchSurveyDetails = async () => {
      loading.value = true;
      error.value = '';
      try {
        const response = await surveyApi.getSurvey(props.id);
        survey.value = response.data;
        if (survey.value && survey.value.is_active === false) {
          error.value = 'This survey is currently inactive and cannot accept responses.';
          // Do not proceed to initialize responses if survey is not active
          loading.value = false; // Stop loading indicator
          return; // Stop further execution
        }
        // Initialize responses object based on fetched questions
        if (survey.value && survey.value.questions) {
          survey.value.questions.forEach(q => {
            if (q.type === 'checkbox') {
              checkboxResponses[q.id] = []; // Initialize as array for checkboxes
              responses[q.id] = []; // Also keep a parallel structure in responses if needed, or solely use checkboxResponses
            } else {
              responses[q.id] = null; // Default for others
            }
          });
        }
      } catch (err) {
        console.error('Error fetching survey details:', err);
        if (err.response?.status === 404) {
          error.value = 'Survey not found. It may have been deleted or the link is incorrect.';
        } else {
          error.value = 'Failed to load survey. Please try again later.';
        }
      } finally {
        if (error.value === 'This survey is currently not active and cannot accept responses.') {
          // If error is due to inactive survey, loading is already set to false
        } else {
          loading.value = false;
        }
      }
    };
    
    const getCheckboxArray = (questionId) => {
      // console.log(`[TakeSurvey] getCheckboxArray for qID: ${questionId}`);
      if (!checkboxResponses[questionId]) {
        // console.log(`[TakeSurvey] Initializing checkboxResponses for qID: ${questionId}`);
        checkboxResponses[questionId] = [];
      }
      return checkboxResponses[questionId];
    };
    
    const isOptionSelected = (questionId, optionText) => {
      const selected = getCheckboxArray(questionId).includes(optionText);
      // console.log(`[TakeSurvey] isOptionSelected for qID: ${questionId}, option: ${optionText}, selected: ${selected}`);
      return selected;
    };
    
    const updateCheckboxValue = (questionId, optionText, value) => {
      // console.log(`[TakeSurvey DEBUG] updateCheckboxValue CALLED - qID: ${questionId}, optionText: ${optionText}, value (should be boolean): ${value}`);
      const currentOptions = getCheckboxArray(questionId);
      if (value) { // if checked
        if (!currentOptions.includes(optionText)) {
          checkboxResponses[questionId].push(optionText);
        }
      } else { // if unchecked
        checkboxResponses[questionId] = currentOptions.filter(item => item !== optionText);
      }
      // Use JSON.stringify and parse to avoid console.log showing a live reference that might change
      // console.log(`[TakeSurvey DEBUG] checkboxResponses for qID ${questionId} AFTER update:`, JSON.parse(JSON.stringify(checkboxResponses[questionId])));
    };
    
    const submitSurvey = async () => {
      if (!valid.value) return;
      
      submitting.value = true;
      error.value = ''; // Clear previous errors
      
      // Combine regular responses with checkbox responses
      const allResponses = {
        ...responses,
        ...Object.keys(checkboxResponses).reduce((acc, key) => {
          acc[key] = checkboxResponses[key];
          return acc;
        }, {})
      };
      
      const submission = {
        surveyId: parseInt(props.id, 10), // Ensure surveyId is an integer
        answers: Object.keys(allResponses).map(questionIdString => ({
          questionId: parseInt(questionIdString, 10), // Ensure questionId is an integer
          value: allResponses[questionIdString]
        }))
      };
      
      try {
        // In production, uncomment this
        await surveyApi.submitResponse(submission, {
          headers: authStore.isAuthenticated ? { Authorization: `Bearer ${authStore.token}` } : {}
        });
        
        // For development, simulate successful submission - REMOVE THIS
        // setTimeout(() => {
        submitted.value = true;
        submitting.value = false;
        // }, 1500);
      } catch (err) {
        console.error('Error submitting response:', err);
        error.value = 'Failed to submit your response. Please try again.';
        submitting.value = false;
      }
    };
    
    onMounted(() => {
      fetchSurveyDetails();
    });
    
    return {
      survey,
      form,
      valid,
      loading,
      submitting,
      error,
      responses,
      submitted,
      getCheckboxArray,
      isOptionSelected,
      updateCheckboxValue,
      submitSurvey
    };
  }
};
</script>

<style scoped>
.take-survey {
  margin-bottom: 60px;
}
</style> 