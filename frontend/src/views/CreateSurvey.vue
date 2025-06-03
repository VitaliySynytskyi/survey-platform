<template>
  <div class="create-survey">
    <v-row>
      <v-col cols="12">
        <v-card class="mb-4 pa-4">
          <div class="d-flex align-center">
            <v-btn icon class="mr-2" to="/dashboard">
              <v-icon>mdi-arrow-left</v-icon>
            </v-btn>
            <h1 class="text-h4">Create New Survey</h1>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <v-row>
      <v-col cols="12">
        <v-card class="pa-6">
          <v-alert
            v-if="error"
            type="error"
            class="mb-4"
            closable
            @click:close="error = ''"
          >
            {{ error }}
          </v-alert>

          <v-form @submit.prevent="saveSurvey" ref="form" v-model="valid">
            <v-row>
              <v-col cols="12" md="8">
                <h2 class="text-h5 mb-4">Survey Details</h2>
                
                <v-text-field
                  v-model="survey.title"
                  label="Survey Title"
                  :rules="titleRules"
                  required
                  variant="outlined"
                  density="comfortable"
                  class="mb-4"
                ></v-text-field>
                
                <v-textarea
                  v-model="survey.description"
                  label="Description"
                  :rules="descriptionRules"
                  required
                  variant="outlined"
                  density="comfortable"
                  class="mb-4"
                  rows="3"
                ></v-textarea>

                <v-switch
                  v-model="survey.is_active"
                  label="Make survey active immediately"
                  color="success"
                  inset
                  class="mb-4"
                ></v-switch>
              </v-col>

              <v-col cols="12" md="4">
                <v-card class="pa-4 mb-4" variant="outlined">
                  <v-card-title class="text-subtitle-1">Survey Settings</v-card-title>
                  <v-card-text>
                    <v-switch
                      v-model="survey.settings.require_login"
                      label="Require login to respond"
                      density="compact"
                      color="primary"
                    ></v-switch>
                    
                    <v-switch
                      v-model="survey.settings.allow_anonymous"
                      label="Allow anonymous responses"
                      density="compact"
                      color="primary"
                      :disabled="survey.settings.require_login"
                    ></v-switch>
                    
                    <v-switch
                      v-model="survey.settings.one_response_per_user"
                      label="One response per user"
                      density="compact"
                      color="primary"
                      :disabled="survey.settings.allow_anonymous && !survey.settings.require_login"
                    ></v-switch>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>

            <v-divider class="my-6"></v-divider>

            <h2 class="text-h5 mb-4">Questions</h2>
            
            <v-alert
              v-if="survey.questions.length === 0"
              type="info"
              class="mb-4"
            >
              No questions added yet. Add your first question below.
            </v-alert>

            <div v-for="(question, index) in survey.questions" :key="index" class="mb-6">
              <v-card class="question-card pa-4" variant="outlined">
                <div class="d-flex justify-space-between align-center mb-2">
                  <span class="text-subtitle-1 font-weight-bold">Question {{ index + 1 }}</span>
                  <v-btn
                    icon
                    color="error"
                    variant="text"
                    @click="removeQuestion(index)"
                    size="small"
                  >
                    <v-icon>mdi-delete</v-icon>
                  </v-btn>
                </div>

                <v-text-field
                  v-model="question.text"
                  label="Question Text"
                  :rules="questionTextRules"
                  required
                  variant="outlined"
                  density="comfortable"
                  class="mb-2"
                ></v-text-field>

                <v-select
                  v-model="question.type"
                  label="Question Type"
                  :items="questionTypes"
                  item-title="text"
                  item-value="value"
                  variant="outlined"
                  density="comfortable"
                  class="mb-2"
                  @update:model-value="onQuestionTypeChange(question)"
                ></v-select>

                <v-switch
                  v-model="question.required"
                  label="Required question"
                  color="primary"
                  density="compact"
                  class="mb-2"
                ></v-switch>

                <!-- Options for multiple choice, checkbox, or dropdown questions -->
                <template v-if="['multiple_choice', 'checkbox', 'dropdown'].includes(question.type)">
                  <div class="mt-2 mb-1 text-subtitle-2">Options</div>
                  
                  <div v-for="(option, optIndex) in question.options" :key="optIndex" class="d-flex align-center mb-2">
                    <v-text-field
                      v-model="question.options[optIndex]"
                      label="Option text"
                      variant="outlined"
                      density="compact"
                      hide-details
                      class="mr-2"
                    ></v-text-field>
                    
                    <v-btn
                      icon
                      color="error"
                      variant="text"
                      @click="removeOption(question, optIndex)"
                      size="small"
                      density="compact"
                    >
                      <v-icon>mdi-close</v-icon>
                    </v-btn>
                  </div>
                  
                  <v-btn
                    prepend-icon="mdi-plus"
                    variant="text"
                    size="small"
                    @click="addOption(question)"
                    class="mt-1"
                    color="primary"
                  >
                    Add Option
                  </v-btn>
                </template>
              </v-card>
            </div>

            <v-btn
              prepend-icon="mdi-plus"
              variant="outlined"
              color="primary"
              @click="addQuestion"
              class="mb-6"
            >
              Add Question
            </v-btn>

            <v-divider class="my-6"></v-divider>

            <div class="d-flex justify-end">
              <v-btn
                to="/dashboard"
                variant="outlined"
                color="grey"
                class="mr-2"
              >
                Cancel
              </v-btn>
              
              <v-btn
                type="submit"
                color="primary"
                :loading="loading"
                :disabled="!valid || loading || survey.questions.length === 0"
              >
                Create Survey
              </v-btn>
            </div>
          </v-form>
        </v-card>
      </v-col>
    </v-row>

    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar.show" :color="snackbar.color" :timeout="3000">
      {{ snackbar.text }}
      <template v-slot:actions>
        <v-btn variant="text" icon="mdi-close" @click="snackbar.show = false"></v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script>
import { ref, reactive, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import axios from '../utils/axiosConfig';

export default {
  name: 'CreateSurvey',
  setup() {
    const router = useRouter();
    const authStore = useAuthStore();
    const form = ref(null);
    const valid = ref(false);
    const loading = ref(false);
    const error = ref('');
    
    const questionTypes = [
      { text: 'Short Answer', value: 'short_answer' },
      { text: 'Paragraph', value: 'paragraph' },
      { text: 'Multiple Choice', value: 'multiple_choice' },
      { text: 'Checkbox', value: 'checkbox' },
      { text: 'Dropdown', value: 'dropdown' },
      { text: 'Linear Scale', value: 'linear_scale' },
      { text: 'Date', value: 'date' }
    ];
    
    const survey = reactive({
      title: '',
      description: '',
      is_active: true,
      settings: {
        require_login: false,
        allow_anonymous: true,
        one_response_per_user: false
      },
      questions: []
    });
    
    const snackbar = ref({
      show: false,
      text: '',
      color: 'success'
    });
    
    const titleRules = [
      v => !!v || 'Title is required',
      v => v.length <= 100 || 'Title must be less than 100 characters'
    ];
    
    const descriptionRules = [
      v => !!v || 'Description is required'
    ];
    
    const questionTextRules = [
      v => !!v || 'Question text is required'
    ];
    
    const addQuestion = () => {
      survey.questions.push({
        text: '',
        type: 'short_answer',
        required: false,
        options: []
      });
    };
    
    const removeQuestion = (index) => {
      survey.questions.splice(index, 1);
    };
    
    const addOption = (question) => {
      if (!question.options) question.options = [];
      question.options.push('');
      
      // Auto-focus on the new option input after Vue updates the DOM
      nextTick(() => {
        const questionIndex = survey.questions.indexOf(question);
        const optionInputs = document.querySelectorAll(`.question-card:nth-child(${questionIndex + 1}) input[label*="Option"]`);
        if (optionInputs.length > 0) {
          optionInputs[optionInputs.length - 1].focus();
        }
      });
    };
    
    const removeOption = (question, index) => {
      if (question.options && question.options.length > 1) {
        question.options.splice(index, 1);
      }
    };
    
    const onQuestionTypeChange = (question) => {
      // Initialize options for question types that need them
      if (['multiple_choice', 'checkbox', 'dropdown'].includes(question.type)) {
        if (!question.options || question.options.length === 0) {
          question.options = ['', ''];
        }
      } else {
        // Clear options for question types that don't need them
        question.options = [];
      }
    };
    
    const saveSurvey = async () => {
      if (!valid.value || survey.questions.length === 0) {
        if (survey.questions.length === 0) {
          showSnackbar('Survey must have at least one question.', 'error');
        } else {
          showSnackbar('Please fix all validation errors before saving.', 'error');
        }
        return;
      }
      
      // Additional validation for questions with options
      for (let i = 0; i < survey.questions.length; i++) {
        const question = survey.questions[i];
        if (['multiple_choice', 'checkbox', 'dropdown'].includes(question.type)) {
          const validOptions = question.options.filter(opt => opt && opt.trim() !== '');
          if (validOptions.length < 2) {
            showSnackbar(`Question ${i + 1} must have at least 2 options.`, 'error');
            return;
          }
        }
      }
      
      loading.value = true;
      error.value = '';
      
      try {
        // Process questions to format options correctly for backend
        const processedSurvey = {
          ...survey,
          questions: survey.questions.map((question, index) => {
            let processedOptions = [];
            if (['multiple_choice', 'checkbox', 'dropdown'].includes(question.type)) {
              // Convert array of strings to array of non-empty strings
              processedOptions = question.options.filter(opt => opt && opt.trim() !== '');
            }
            return {
              text: question.text,
              type: question.type,
              required: question.required,
              order_num: index + 1,
              options: processedOptions
            };
          })
        };
        
        const response = await axios.post('/api/v1/surveys', processedSurvey, {
          headers: { Authorization: `Bearer ${authStore.token}` }
        });
        
        showSnackbar('Survey created successfully!', 'success');
        // Small delay to show success message before redirect
        setTimeout(() => {
          router.push('/dashboard');
        }, 1000);
      } catch (err) {
        console.error('Error creating survey:', err);
        const errorMsg = err.response?.data?.error || 'Failed to create survey. Please try again.';
        error.value = errorMsg;
        showSnackbar(errorMsg, 'error');
      } finally {
        loading.value = false;
      }
    };
    
    const showSnackbar = (text, color = 'success') => {
      snackbar.value = {
        show: true,
        text,
        color
      };
    };
    
    return {
      form,
      valid,
      loading,
      error,
      survey,
      questionTypes,
      snackbar,
      titleRules,
      descriptionRules,
      questionTextRules,
      addQuestion,
      removeQuestion,
      addOption,
      removeOption,
      onQuestionTypeChange,
      saveSurvey
    };
  }
};
</script>

<style scoped>
.question-card {
  border-left: 4px solid var(--primary-color);
}
</style> 