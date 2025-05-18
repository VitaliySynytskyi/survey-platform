<template>
  <div class="survey-responses">
    <v-row>
      <v-col cols="12">
        <v-card class="mb-4 pa-4">
          <div class="d-flex align-center">
            <v-btn icon class="mr-2" to="/dashboard">
              <v-icon>mdi-arrow-left</v-icon>
            </v-btn>
            <h1 class="text-h4 mr-4">Survey Responses</h1>
            <v-chip
              v-if="surveyDetails.title"
              color="primary"
              size="large"
            >
              {{ surveyDetails.title }}
            </v-chip>
            <v-spacer></v-spacer>
            <v-btn
              color="primary"
              variant="outlined"
              :to="`/surveys/${props.id}/edit`"
              prepend-icon="mdi-pencil"
              class="ml-2"
            >
              Edit Survey
            </v-btn>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-if="loading">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <div class="mt-2">Loading survey data...</div>
      </v-col>
    </v-row>

    <v-row v-else-if="error">
      <v-col cols="12">
        <v-alert type="error" closable @click:close="error = ''">
          {{ error }}
        </v-alert>
      </v-col>
    </v-row>

    <template v-else>
      <v-row>
        <v-col cols="12">
          <v-card class="mb-4">
            <v-card-title class="text-h5 mb-2">Overview</v-card-title>
            <v-card-text>
              <v-row>
                <v-col cols="12" md="3">
                  <v-card class="pa-4 text-center">
                    <div class="text-h4 mb-1">{{ responses.length }}</div>
                    <div class="text-subtitle-1">Total Responses</div>
                  </v-card>
                </v-col>
                <v-col cols="12" md="3">
                  <v-card class="pa-4 text-center">
                    <div class="text-h4 mb-1">{{ surveyDetails.responseRate !== undefined ? surveyDetails.responseRate + '%' : 'N/A' }}</div>
                    <div class="text-subtitle-1">Completion Rate</div>
                  </v-card>
                </v-col>
                <v-col cols="12" md="3">
                  <v-card class="pa-4 text-center">
                    <div class="text-h4 mb-1">{{ surveyDetails.avgCompletionTime || 'N/A' }}</div>
                    <div class="text-subtitle-1">Avg. Time to Complete</div>
                  </v-card>
                </v-col>
                <v-col cols="12" md="3">
                  <v-card class="pa-4 text-center">
                    <div class="text-h4 mb-1">{{ lastResponseDateDisplay }}</div>
                    <div class="text-subtitle-1">Last Response</div>
                  </v-card>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <v-row>
        <v-col cols="12">
          <v-card>
            <v-tabs v-model="activeTab">
              <v-tab value="summary">Summary</v-tab>
              <v-tab value="individual">Individual Responses</v-tab>
              <v-tab value="export">Export Data</v-tab>
            </v-tabs>
            
            <v-card-text>
              <v-window v-model="activeTab">
                <!-- Summary Tab -->
                <v-window-item value="summary">
                  <div v-if="!surveyDetails.questions || surveyDetails.questions.length === 0" class="text-center pa-4">
                    No questions in this survey or survey data not loaded.
                  </div>
                  
                  <div v-else class="question-summaries">
                    <div v-for="(question, index) in surveyDetails.questions" :key="question.id" class="mb-8">
                      <h3 class="text-h6 mb-2">{{ index + 1 }}. {{ question.text }}</h3>
                      
                      <template v-if="['multiple_choice', 'checkbox', 'dropdown'].includes(question.type)">
                        <v-card class="pa-4" variant="outlined">
                          <div v-for="(option, optIndex) in question.options" :key="optIndex" class="mb-2">
                            <div class="d-flex align-center mb-1">
                              <div class="text-body-1 mr-2" style="width: 200px;">{{ option.text }}</div>
                              <div class="mr-2" style="flex: 1;">
                                <v-progress-linear
                                  :model-value="getOptionPercentage(question.id, option.text, question.type)"
                                  height="20"
                                  color="primary"
                                  rounded
                                >
                                  <template v-slot:default="{ value }">
                                    <strong>{{ Math.ceil(value) }}%</strong>
                                  </template>
                                </v-progress-linear>
                              </div>
                              <div class="text-body-2" style="width: 50px;">
                                {{ getOptionCount(question.id, option.text, question.type) }}
                              </div>
                            </div>
                          </div>
                        </v-card>
                      </template>
                      <template v-else-if="question.type === 'short_answer'">
                        <v-card class="pa-4" variant="outlined">
                          <v-list density="compact">
                            <v-list-subheader>Responses Summary</v-list-subheader>
                            <div v-if="getAggregatedShortAnswers(question.id).length === 0" class="text-grey pa-2">
                              No responses for this question yet.
                            </div>
                            <v-list-item
                              v-for="(item, aggIdx) in getAggregatedShortAnswers(question.id)"
                              :key="aggIdx"
                            >
                              <v-list-item-title>"{{ item.value }}"</v-list-item-title>
                              <template v-slot:append>
                                <v-chip color="primary" size="small">{{ item.count }}</v-chip>
                              </template>
                            </v-list-item>
                          </v-list>
                        </v-card>
                      </template>
                      <template v-else-if="question.type === 'paragraph'">
                        <v-card class="pa-4" variant="outlined">
                           <v-list density="compact">
                            <v-list-subheader>Responses (first {{ maxTextResponsesToShow }} shown)</v-list-subheader>
                            <div v-if="getTextResponsesForQuestion(question.id).length === 0" class="text-grey pa-2">
                              No text responses for this question yet.
                            </div>
                            <v-list-item 
                              v-for="(answer, ansIndex) in getTextResponsesForQuestion(question.id).slice(0, maxTextResponsesToShow)" 
                              :key="ansIndex"
                              :title="answer.value"
                            ></v-list-item>
                            <v-list-item v-if="getTextResponsesForQuestion(question.id).length > maxTextResponsesToShow">
                              ...and {{ getTextResponsesForQuestion(question.id).length - maxTextResponsesToShow }} more.
                            </v-list-item>
                          </v-list>
                        </v-card>
                      </template>
                      <!-- Linear Scale -->
                      <template v-else-if="question.type === 'linear_scale'">
                        <v-card class="pa-4" variant="outlined">
                           <v-list-subheader class="mb-2">Responses Distribution</v-list-subheader>
                           <div v-if="getLinearScaleDistribution(question.id).length === 0" class="text-grey pa-2">
                              No responses for this question yet.
                           </div>
                           <div v-for="(item, scaleIdx) in getLinearScaleDistribution(question.id)" :key="scaleIdx" class="d-flex align-center mb-1">
                              <div class="text-body-1 mr-2" style="width: 60px;">Scale {{ item.scaleValue }}:</div>
                              <div class="mr-2" style="flex: 1;">
                                <v-progress-linear
                                  :model-value="item.percentage"
                                  height="20"
                                  color="secondary"
                                  rounded
                                >
                                  <template v-slot:default="{ value }">
                                    <strong>{{ Math.ceil(value) }}%</strong>
                                  </template>
                                </v-progress-linear>
                              </div>
                              <div class="text-body-2" style="width: 80px;">
                                {{ item.count }} response(s)
                              </div>
                           </div>
                        </v-card>
                      </template>
                      <!-- Date -->
                      <template v-else-if="question.type === 'date'">
                         <v-card class="pa-4" variant="outlined">
                           <v-list density="compact">
                            <v-list-subheader>Date Responses (Top {{maxTextResponsesToShow}} shown with counts)</v-list-subheader>
                             <div v-if="getAggregatedDateAnswers(question.id).length === 0" class="text-grey pa-2">
                                No date responses for this question yet.
                             </div>
                            <v-list-item 
                              v-for="(item, dateIdx) in getAggregatedDateAnswers(question.id).slice(0, maxTextResponsesToShow)" 
                              :key="dateIdx"
                            >
                              <v-list-item-title>{{ formatDate(item.value, { year: 'numeric', month: 'long', day: 'numeric' }) }}</v-list-item-title>
                              <template v-slot:append>
                                <v-chip color="primary" size="small">{{ item.count }}</v-chip>
                              </template>
                            </v-list-item>
                             <v-list-item v-if="getAggregatedDateAnswers(question.id).length > maxTextResponsesToShow">
                               ...and {{ getAggregatedDateAnswers(question.id).length - maxTextResponsesToShow }} more date entries.
                             </v-list-item>
                           </v-list>
                        </v-card>
                      </template>
                    </div>
                  </div>
                </v-window-item>
                
                <!-- Individual Responses Tab -->
                <v-window-item value="individual">
                  <div v-if="responses.length === 0" class="text-center pa-4">
                    No responses received yet.
                  </div>
                  
                  <div v-else>
                    <v-data-table
                      :headers="individualResponseHeaders"
                      :items="responses"
                      class="elevation-1"
                      item-value="id"
                    >
                      <template v-slot:item.submittedAt="{ item }">
                        {{ formatDate(item.submittedAt) }}
                      </template>
                       <template v-slot:item.userId="{ item }">
                        {{ item.userId || 'Anonymous' }}
                      </template>
                      <template v-slot:item.id="{ item }">
                        {{ item.id || item._id }}
                      </template>
                      <template v-slot:item.actions="{ item }">
                        <v-btn
                          icon
                          variant="text"
                          size="small"
                          @click="viewResponseDetails(item)"
                        >
                          <v-icon>mdi-eye</v-icon>
                        </v-btn>
                      </template>
                    </v-data-table>
                  </div>
                </v-window-item>
                
                <!-- Export Data Tab -->
                <v-window-item value="export">
                  <v-card flat>
                    <v-card-title>Export Survey Responses</v-card-title>
                    <v-card-text>
                      <p class="mb-4">
                        Download all responses for this survey in CSV format.
                      </p>
                      <v-btn 
                        color="primary"
                        @click="exportResponsesCSV"
                        :loading="exporting"
                        prepend-icon="mdi-download-box-outline"
                      >
                        Export to CSV
                      </v-btn>
                      <v-alert v-if="exportError" type="error" class="mt-4" closable @click:close="exportError = null">
                        {{ exportError }}
                      </v-alert>
                    </v-card-text>
                  </v-card>
                </v-window-item>
              </v-window>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>
    
    <!-- Individual Response Dialog -->
    <v-dialog v-model="dialog" max-width="700px">
      <v-card v-if="selectedResponse">
        <v-card-title>
          <span class="text-h5">Response Details (ID: {{ selectedResponse.id || selectedResponse._id }})</span>
          <v-spacer></v-spacer>
          <span class="text-subtitle-2">Submitted: {{ formatDate(selectedResponse.submittedAt) }}</span>
        </v-card-title>
        <v-card-text v-if="surveyDetails.questions">
          <v-list dense>
            <v-list-item v-for="question in surveyDetails.questions" :key="question.id">
              <v-list-item-title class="font-weight-bold">{{ question.text }}</v-list-item-title>
              <v-list-item-subtitle>
                {{ getAnswerForQuestion(selectedResponse, question.id) }}
              </v-list-item-subtitle>
            </v-list-item>
          </v-list>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" variant="text" @click="dialog = false">Close</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- Snackbar for notifications -->
    <v-snackbar v-model="snackbar.show" :color="snackbar.color" :timeout="3000">
      {{ snackbar.text }}
      <template v-slot:actions>
        <v-btn icon="mdi-close" variant="text" @click="snackbar.show = false"></v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../store/auth';
import axios from '../utils/axiosConfig';
import { saveAs } from 'file-saver';

const props = defineProps({
  id: { type: String, required: true }
});

const authStore = useAuthStore();

const loading = ref(true);
const error = ref('');
const surveyDetails = ref({ title: '', questions: [], settings: {}, responseRate: undefined, avgCompletionTime: undefined, lastResponseDate: undefined });
const responses = ref([]);
const activeTab = ref('summary');

const snackbar = ref({ show: false, text: '', color: 'success' });

// For individual response dialog
const dialog = ref(false);
const selectedResponse = ref(null);
const maxTextResponsesToShow = 5; // For summary view

// For export tab
const exportFormat = ref('csv');
const exporting = ref(false);
const exportOptions = ref({
  includeTimestamp: true,
  includeUserId: true,
});
const exportError = ref(null);

const individualResponseHeaders = [
  { title: 'Submission ID', key: 'id' }, // Changed from _id to id based on MongoDB model
  { title: 'User ID', key: 'userId' },
  { title: 'Submitted At', key: 'submittedAt' },
  { title: 'Actions', key: 'actions', sortable: false },
];

const fetchSurveyData = async () => {
  loading.value = true;
  error.value = '';
  try {
    // Fetch survey details (title, questions)
    const surveyDetailResponse = await axios.get(`/api/v1/surveys/${props.id}`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    });
    surveyDetails.value = surveyDetailResponse.data;
    if (!surveyDetails.value.questions) surveyDetails.value.questions = [];
    if (!surveyDetails.value.settings) surveyDetails.value.settings = {};

    // Fetch responses for the survey
    const responsesResponse = await axios.get(`/api/v1/surveys/${props.id}/responses`, {
      headers: { Authorization: `Bearer ${authStore.token}` }
    });
    responses.value = responsesResponse.data;

  } catch (err) {
    console.error('Error fetching survey data:', err);
    error.value = err.response?.data?.error || 'Failed to load survey data. Please try again.';
    if (err.response?.status === 404 && (!surveyDetails.value || !surveyDetails.value.title)) {
        error.value = "Survey details not found. It may have been deleted or the ID is incorrect.";
    } else if (err.response?.status === 404) {
        error.value = "Failed to load responses. The survey details were found, but responses could not be retrieved or there are no responses yet.";
        responses.value = []; // Ensure responses are empty if not found
    }
  } finally {
    loading.value = false;
  }
};

const formatDate = (dateString, customOptions) => {
  if (!dateString) return 'N/A';
  const defaultOptions = { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' };
  const optionsToUse = customOptions || defaultOptions;
  try {
    return new Date(dateString).toLocaleDateString(undefined, optionsToUse);
  } catch (e) {
    return dateString; // Return original if parsing fails
  }
};

const lastResponseDateDisplay = computed(() => {
  if (!responses.value || responses.value.length === 0) {
    return 'N/A';
  }
  // Assuming responses are not necessarily sorted
  let latestDate = new Date(0);
  responses.value.forEach(response => {
    if (response.submittedAt) {
      const submittedDate = new Date(response.submittedAt);
      if (submittedDate > latestDate) {
        latestDate = submittedDate;
      }
    }
  });
  return latestDate.getTime() === new Date(0).getTime() ? 'N/A' : formatDate(latestDate.toISOString());
});

// --- Methods for Summary Tab ---
const getOptionCount = (questionId, optionText, questionType) => {
  let count = 0;
  responses.value.forEach(response => {
    const answer = response.answers.find(a => a.questionId === questionId);
    if (answer) {
      if (questionType === 'checkbox' && Array.isArray(answer.value)) {
        if (answer.value.includes(optionText)) {
          count++;
        }
      } else if (answer.value === optionText) {
        count++;
      }
    }
  });
  return count;
};

const getOptionPercentage = (questionId, optionText, questionType) => {
  const count = getOptionCount(questionId, optionText, questionType);
  
  if (questionType === 'checkbox') {
    // For checkboxes, percentage is out of total responses, as one response can have multiple checks.
    return responses.value.length > 0 ? (count / responses.value.length) * 100 : 0;
  } else {
    // For single-select (multiple_choice, dropdown), percentage is out of those who answered THIS question.
    const totalResponsesForQuestion = responses.value.filter(r => 
      r.answers.some(a => a.questionId === questionId && a.value !== null && a.value !== undefined && (Array.isArray(a.value) ? a.value.length > 0 : a.value !== ''))
    ).length;
    return totalResponsesForQuestion > 0 ? (count / totalResponsesForQuestion) * 100 : 0;
  }
};

const getTextResponsesForQuestion = (questionId) => {
  return responses.value
    .map(response => response.answers.find(a => a.questionId === questionId))
    .filter(answer => answer && answer.value && typeof answer.value === 'string' && answer.value.trim() !== '') 
    .map(answer => answer); // Returns the answer object {questionId, value}
};

const getAggregatedShortAnswers = (questionId) => {
  const answersMap = new Map();
  responses.value.forEach(response => {
    const answer = response.answers.find(a => a.questionId === questionId);
    if (answer && typeof answer.value === 'string' && answer.value.trim() !== '') {
      answersMap.set(answer.value, (answersMap.get(answer.value) || 0) + 1);
    }
  });
  return Array.from(answersMap, ([value, count]) => ({ value, count })).sort((a,b) => b.count - a.count);
};

const getLinearScaleDistribution = (questionId) => {
  const scaleCounts = new Map();
  let totalResponsesForQuestion = 0;
  responses.value.forEach(response => {
    const answer = response.answers.find(a => a.questionId === questionId);
    if (answer && answer.value !== null && answer.value !== undefined && answer.value !== '') {
      totalResponsesForQuestion++;
      const scaleValue = answer.value.toString(); // Ensure it's a string for map key
      scaleCounts.set(scaleValue, (scaleCounts.get(scaleValue) || 0) + 1);
    }
  });

  const distribution = [];
  // Assuming scale is 1-5 from TakeSurvey.vue, but ideally this should come from question.settings
  const assumedScalePoints = surveyDetails.value?.questions?.find(q => q.id === questionId)?.linear_scale_max || 5; // example default
  for (let i = 1; i <= assumedScalePoints; i++) {
    const valStr = i.toString();
    const count = scaleCounts.get(valStr) || 0;
    distribution.push({
      scaleValue: valStr,
      count: count,
      percentage: totalResponsesForQuestion > 0 ? (count / totalResponsesForQuestion) * 100 : 0
    });
  }
  return distribution;
};

const getAggregatedDateAnswers = (questionId) => {
  const answersMap = new Map();
  responses.value.forEach(response => {
    const answer = response.answers.find(a => a.questionId === questionId);
    if (answer && typeof answer.value === 'string' && answer.value.trim() !== '') {
        try {
            // Normalize date to YYYY-MM-DD for consistent grouping
            const dateValue = new Date(answer.value).toISOString().split('T')[0];
            answersMap.set(dateValue, (answersMap.get(dateValue) || 0) + 1);
        } catch (e) {
            // Could be an invalid date string, perhaps log or count as 'invalid'
            answersMap.set(answer.value, (answersMap.get(answer.value) || 0) + 1); 
        }
    }
  });
  return Array.from(answersMap, ([value, count]) => ({ value, count })).sort((a,b) => b.count - a.count);
};
// --- End Methods for Summary Tab ---


// --- Methods for Individual Response Dialog ---
const viewResponseDetails = (responseItem) => {
  selectedResponse.value = responseItem;
  dialog.value = true;
};

const getAnswerForQuestion = (response, questionId) => {
  if (!response || !response.answers) return 'No data';
  const answer = response.answers.find(a => a.questionId === questionId);
  return answer ? answer.value : 'No answer provided';
};
// --- End Methods for Individual Response Dialog ---

const showSnackbar = (text, color = 'success') => {
  snackbar.value.text = text;
  snackbar.value.color = color;
  snackbar.value.show = true;
};

const exportResponsesCSV = async () => {
  exporting.value = true;
  exportError.value = null;
  try {
    const response = await axios.get(`/api/v1/surveys/${props.id}/responses/export`, {
      responseType: 'blob',
      headers: { Authorization: `Bearer ${authStore.token}` }
    });

    const contentDisposition = response.headers['content-disposition'];
    let filename = `survey_${props.id}_responses.csv`;
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/);
      if (filenameMatch != null && filenameMatch[1]) {
        filename = filenameMatch[1].replace(/['"]/g, '');
      }
    }

    saveAs(new Blob(["\uFEFF", response.data], { type: 'text/csv;charset=utf-8;' }), filename);

  } catch (err) {
    console.error("Error exporting CSV:", err);
    if (err.response && err.response.status === 204) {
      exportError.value = "No responses found to export.";
    } else if (err.response && err.response.data) {
        try {
            // If the error response is a blob, try to read it as text for a JSON error message
            const errorText = await (err.response.data).text();
            const errorJson = JSON.parse(errorText);
            exportError.value = `Failed to export CSV: ${errorJson.error || 'Server error'}`;
        } catch (parseError) {
            exportError.value = "Failed to export CSV. An unknown error occurred while parsing the error response.";
        }
    } else {
        exportError.value = "Failed to export CSV. An unknown error occurred.";
    }
  } finally {
    exporting.value = false;
  }
};

onMounted(() => {
  fetchSurveyData();
});

</script>

<style scoped>
.survey-responses .v-card {
  margin-bottom: 16px;
}
.question-summaries .v-card {
  margin-bottom: 24px;
}
</style> 