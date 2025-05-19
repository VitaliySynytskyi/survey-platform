<template>
  <div class="edit-survey">
    <v-row>
      <v-col cols="12">
        <v-card class="mb-4 pa-4">
          <div class="d-flex align-center">
            <v-btn icon class="mr-2" to="/dashboard">
              <v-icon>mdi-arrow-left</v-icon>
            </v-btn>
            <h1 class="text-h4">Edit Survey: {{ survey.title || 'Loading...' }}</h1>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-if="loading && !initialLoadComplete">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <div class="mt-2">Loading survey details...</div>
      </v-col>
    </v-row>

    <v-row v-else-if="error">
      <v-col cols="12">
        <v-alert type="error" closable @click:close="error = ''">
          {{ error }}
        </v-alert>
      </v-col>
    </v-row>

    <v-row v-else>
      <v-col cols="12">
        <v-card class="pa-6">
          <v-form @submit.prevent="updateSurveyHandler" ref="form" v-model="valid">
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
                  label="Survey is active"
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
            <v-alert v-if="survey.questions.length === 0" type="info" class="mb-4">
              No questions in this survey. Add questions below.
            </v-alert>

            <div v-for="(question, index) in survey.questions" :key="question.id || index" class="mb-6">
              <v-card class="question-card pa-4" variant="outlined">
                <div class="d-flex justify-space-between align-center mb-2">
                  <span class="text-subtitle-1 font-weight-bold">Question {{ index + 1 }}</span>
                  <v-btn icon color="error" variant="text" @click="removeQuestion(index)" size="small">
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
                ></v-select>
                <v-switch
                  v-model="question.required"
                  label="Required question"
                  color="primary"
                  density="compact"
                  class="mb-2"
                ></v-switch>
                <template v-if="['multiple_choice', 'checkbox', 'dropdown'].includes(question.type)">
                  <div class="mt-2 mb-1 text-subtitle-2">Options</div>
                  <div v-for="(optionObject, optIndex) in question.options" :key="optIndex" class="d-flex align-center mb-2">
                    <v-text-field
                      v-model="optionObject.text"
                      label="Option text"
                      variant="outlined"
                      density="compact"
                      hide-details
                      class="mr-2"
                    ></v-text-field>
                    <v-btn icon color="error" variant="text" @click="removeOption(question, optIndex)" size="small" density="compact">
                      <v-icon>mdi-close</v-icon>
                    </v-btn>
                  </div>
                  <v-btn prepend-icon="mdi-plus" variant="text" size="small" @click="addOption(question)" class="mt-1" color="primary">
                    Add Option
                  </v-btn>
                </template>
              </v-card>
            </div>

            <v-btn prepend-icon="mdi-plus" variant="outlined" color="primary" @click="addQuestion" class="mb-6">
              Add Question
            </v-btn>

            <v-divider class="my-6"></v-divider>

            <div class="d-flex justify-end">
              <v-btn to="/dashboard" variant="outlined" color="grey" class="mr-2">Cancel</v-btn>
              <v-btn type="submit" color="primary" :loading="saving" :disabled="!valid || saving || survey.questions.length === 0">
                Save Changes
              </v-btn>
            </div>
          </v-form>
        </v-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar.show" :color="snackbar.color" :timeout="3000">
      {{ snackbar.text }}
      <template v-slot:actions>
        <v-btn variant="text" icon="mdi-close" @click="snackbar.show = false"></v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script>
import { ref, reactive, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';
import axios from '../utils/axiosConfig';

export default {
  name: 'EditSurvey',
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
    const valid = ref(false);
    const loading = ref(true);
    const initialLoadComplete = ref(false);
    const saving = ref(false);
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
      id: props.id,
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
    const descriptionRules = [v => !!v || 'Description is required'];
    const questionTextRules = [v => !!v || 'Question text is required'];

    const fetchSurveyDetails = async () => {
      loading.value = true;
      error.value = '';
      try {
        const response = await axios.get(`/api/v1/surveys/${props.id}`, {
          headers: { Authorization: `Bearer ${authStore.token}` }
        });
        const fetchedSurvey = response.data;
        
        // Ensure settings object exists
        if (!fetchedSurvey.settings) {
          fetchedSurvey.settings = { require_login: false, allow_anonymous: true, one_response_per_user: false };
        }
        // Ensure questions array exists and process options
        if (!fetchedSurvey.questions) {
           fetchedSurvey.questions = [];
        } else {
          fetchedSurvey.questions = fetchedSurvey.questions.map(q => {
            if (['multiple_choice', 'checkbox', 'dropdown'].includes(q.type)) {
              if (q.options && q.options.length > 0) {
                q.options = q.options.map(opt => {
                  if (typeof opt === 'string') {
                    return { text: opt };
                  } else if (typeof opt === 'object' && opt !== null && typeof opt.text === 'string') {
                    return { text: opt.text }; // Ensure it's in the {text: ...} format we want
                  }
                  return { text: '' }; // Fallback for malformed options
                });
              } else {
                q.options = []; // Ensure it's an empty array if no options
              }
            }
            return q;
          });
        }
        Object.assign(survey, fetchedSurvey);

      } catch (err) {
        console.error('Error fetching survey details:', err);
        error.value = err.response?.data?.error || 'Failed to load survey details. Please try again.';
        if (err.response?.status === 404) {
            error.value = "Survey not found. It may have been deleted or the ID is incorrect.";
        }
      } finally {
        loading.value = false;
        initialLoadComplete.value = true;
      }
    };

    const addQuestion = () => {
      const newQuestion = {
        // id: null, // Let backend assign ID for new questions if necessary
        text: '',
        type: 'short_answer',
        required: false,
        options: [] // Initialize with empty options array
      };
      // For types that use options, one could conditionally add an initial empty option here if desired
      // e.g., if (['multiple_choice', 'checkbox', 'dropdown'].includes(newQuestion.type)) {
      //   newQuestion.options.push('');
      // }
      survey.questions.push(newQuestion);
    };

    const removeQuestion = (index) => {
      survey.questions.splice(index, 1);
    };

    const addOption = (question) => {
      if (!question.options) question.options = [];
      question.options.push({ text: '' });
    };

    const removeOption = (question, index) => {
      if (question.options) { // Ensure options array exists
        question.options.splice(index, 1);
      }
    };

    const updateSurveyHandler = async () => {
      // Validate form first
      const { valid: formIsValid } = await form.value.validate();
      if (!formIsValid || survey.questions.length === 0) {
        if (survey.questions.length === 0) {
          showSnackbar('Survey must have at least one question.', 'error');
        }
        return;
      }
      
      saving.value = true;
      error.value = '';
      try {
        const payload = {
            ...survey, // Spread existing survey properties like title, description, settings
            id: survey.id, // Ensure id is part of the payload
            questions: survey.questions.map(q => {
              let processedOptions = [];
              if (['multiple_choice', 'checkbox', 'dropdown'].includes(q.type)) {
                // Ensure options are strings and not empty before sending to backend
                processedOptions = q.options ? q.options.map(opt => opt.text.trim()).filter(optText => optText !== '') : [];
              }
              return {
                id: q.id, // Include existing question ID if it exists (for updates)
                text: q.text,
                type: q.type,
                required: q.required,
                options: processedOptions 
              };
            })
        };
        await axios.put(`/api/v1/surveys/${props.id}`, payload, {
          headers: { Authorization: `Bearer ${authStore.token}` }
        });
        showSnackbar('Survey updated successfully!', 'success');
        await fetchSurveyDetails(); // Re-fetch data after successful save
        router.push('/dashboard');
      } catch (err) {
        console.error('Error updating survey:', err);
        error.value = err.response?.data?.error || 'Failed to update survey. Please try again.';
      } finally {
        saving.value = false;
      }
    };

    const showSnackbar = (text, color = 'success') => {
      snackbar.value.text = text;
      snackbar.value.color = color;
      snackbar.value.show = true;
    };

    onMounted(() => {
      fetchSurveyDetails();
    });

    return {
      form,
      valid,
      loading,
      initialLoadComplete,
      saving,
      error,
      survey,
      questionTypes,
      titleRules,
      descriptionRules,
      questionTextRules,
      addQuestion,
      removeQuestion,
      addOption,
      removeOption,
      updateSurveyHandler,
      snackbar,
      showSnackbar
    };
  }
};
</script>

<style scoped>
.question-card {
  border-left: 4px solid var(--primary-color);
}
</style> 