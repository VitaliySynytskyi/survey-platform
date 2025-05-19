<template>
    <v-container>
      <v-row>
        <v-col cols="12">
          <h1 class="text-h4 mb-4">Survey Analytics: {{ surveyTitle }}</h1>
          <p v-if="analyticsData && analyticsData.survey_id" class="text-subtitle-1 grey--text text--darken-1">Survey ID: {{ analyticsData.survey_id }}</p>
        </v-col>
      </v-row>
  
      <v-row v-if="loading">
        <v-col cols="12" class="text-center py-10">
          <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
          <p class="mt-4 text-subtitle-1">Loading analytics data...</p>
        </v-col>
      </v-row>
  
      <v-row v-if="error">
        <v-col cols="12">
          <v-alert type="error" prominent border="left" elevation="2">
            <div class="text-h6">Error Loading Analytics</div>
            <p>{{ error }}</p>
          </v-alert>
        </v-col>
      </v-row>
  
      <v-row v-if="!loading && !error && analyticsData">
        <v-col cols="12" md="4">
          <v-card class="mb-5 elevation-2">
            <v-card-item>
                <v-card-title class="text-h6">Overall Summary</v-card-title>
            </v-card-item>
            <v-divider></v-divider>
            <v-card-text class="pa-5">
              <div class="d-flex align-center mb-3">
                <v-icon color="primary" class="mr-3">mdi-account-multiple</v-icon>
                <span class="text-body-1"><strong>Total Responses:</strong> {{ analyticsData.total_responses }}</span>
              </div>
               <div class="d-flex align-center">
                <v-icon color="primary" class="mr-3">mdi-list-status</v-icon>
                <span class="text-body-1"><strong>Total Questions:</strong> {{ analyticsData.question_analytics ? analyticsData.question_analytics.length : 0 }}</span>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <v-col cols="12" md="8" v-if="analyticsData.total_responses > 0">
             <v-alert type="info" variant="tonal" density="compact" class="mb-5">
                Percentages for question options are calculated based on the number of users who answered that specific question.
            </v-alert>
        </v-col>
  
        <v-col cols="12" v-for="question in analyticsData.question_analytics" :key="question.question_id">
          <v-card class="mb-5 elevation-2">
            <v-card-item>
                 <v-card-title class="text-h6">{{ question.question_text }}</v-card-title>
                 <v-card-subtitle>Question Type: <v-chip size="small" color="teal" label class="ml-1">{{ question.question_type }}</v-chip></v-card-subtitle>
            </v-card-item>
            <v-divider></v-divider>
            <v-card-text class="pa-5">
              <div v-if="isChartable(question.question_type) && question.options_summary && question.options_summary.length > 0">
                <Bar :data="getChartData(question)" :options="chartOptions" style="max-height: 300px;" />
              </div>
              <div v-else-if="isChartable(question.question_type) && (!question.options_summary || question.options_summary.length === 0)">
                 <p class="text-center grey--text py-5">No responses recorded for this question yet to display a chart.</p>
              </div>

              <div v-if="(question.question_type === 'text' || question.question_type === 'paragraph' || question.question_type === 'short_answer' || question.question_type === 'date') && question.text_responses">
                <h4 class="text-subtitle-1 font-weight-medium mb-2">Text Responses ({{ question.text_responses.length }}):</h4>
                <v-list density="compact" v-if="question.text_responses.length > 0">
                  <v-list-item v-for="(response, index) in question.text_responses.slice(0, displayedTextResponsesCount)" :key="index" class="pa-1 mb-1 border rounded">
                     <v-list-item-title class="text-body-2">{{ response.response }}</v-list-item-title>
                  </v-list-item>
                </v-list>
                 <v-btn 
                    v-if="question.text_responses.length > displayedTextResponsesCount"
                    @click="displayedTextResponsesCount += 10"
                    small variant="outlined" color="primary" class="mt-2">
                    Show More
                </v-btn>
                <p v-else-if="question.text_responses.length === 0" class="text-center grey--text py-5">No text responses recorded for this question yet.</p>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

         <v-col cols="12" v-if="showRawData">
            <v-divider class="my-6"></v-divider>
            <div class="d-flex justify-space-between align-center mb-2">
                <h3 class="text-h5">Raw Analytics Data</h3>
                 <v-btn icon @click="showRawData = false"><v-icon>mdi-close</v-icon></v-btn>
            </div>
            <pre class="pa-3 elevation-1 rounded">{{ JSON.stringify(analyticsData, null, 2) }}</pre>
          </v-col>
      </v-row>
      <v-row v-if="!loading && !error && analyticsData && analyticsData.total_responses === 0">
        <v-col cols="12">
            <v-alert type="info" prominent border="left" elevation="2" icon="mdi-information-outline">
                <div class="text-h6">No Responses Yet</div>
                <p>This survey currently has no responses. Analytics will be shown here once responses are submitted.</p>
            </v-alert>
        </v-col>
      </v-row>
    </v-container>
  </template>
  
  <script setup>
  import { ref, onMounted, computed } from 'vue';
  import { useRoute } from 'vue-router';
  import { surveyApi } from '@/services/api'; 
  import { useAuthStore } from '@/store/auth';
  import { Bar } from 'vue-chartjs';
  import {
    Chart as ChartJS,
    Title,
    Tooltip,
    Legend,
    BarElement,
    CategoryScale,
    LinearScale
  } from 'chart.js';
  
  ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale);
  
  const route = useRoute();
  const authStore = useAuthStore();
  
  const surveyId = ref(route.params.id);
  const surveyTitle = ref('Loading survey title...');
  const analyticsData = ref(null);
  const loading = ref(true);
  const error = ref(null);
  const showRawData = ref(false); // To toggle raw data view, for debugging
  const displayedTextResponsesCount = ref(5); // For text responses pagination

  const chartColors = [
    'rgba(54, 162, 235, 0.8)', // Blue
    'rgba(255, 99, 132, 0.8)', // Red
    'rgba(75, 192, 192, 0.8)', // Green
    'rgba(255, 206, 86, 0.8)', // Yellow
    'rgba(153, 102, 255, 0.8)', // Purple
    'rgba(255, 159, 64, 0.8)',  // Orange
    'rgba(199, 199, 199, 0.8)', // Grey
    'rgba(83, 102, 255, 0.8)',  // Indigo
    'rgba(255, 99, 255, 0.8)',  // Pink
    'rgba(102, 255, 83, 0.8)'   // Light Green
  ];

  const chartOptions = computed(() => ({
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false, // Usually not needed for single dataset bar charts per question
      },
      tooltip: {
        callbacks: {
          label: function(context) {
            let label = context.dataset.label || '';
            if (label) {
              label += ': ';
            }
            if (context.parsed.y !== null) {
              label += context.parsed.y + ' responses';
            }
            const percentage = context.dataset.percentages ? context.dataset.percentages[context.dataIndex] : null;
            if (percentage !== null && percentage !== undefined) {
                label += ` (${percentage.toFixed(1)}%)`;
            }
            return label;
          }
        }
      }
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          precision: 0 // Ensure y-axis shows whole numbers for counts
        }
      }
    }
  }));

  const isChartable = (questionType) => {
      return ['single_choice', 'multiple_choice', 'dropdown', 'linear_scale', 'checkbox'].includes(questionType);
  }
  
  const getChartData = (question) => {
    if (!question || !question.options_summary) {
      return { labels: [], datasets: [] };
    }
    const labels = question.options_summary.map(opt => opt.option_text);
    const data = question.options_summary.map(opt => opt.count);
    const percentages = question.options_summary.map(opt => opt.percentage);
  
    return {
      labels,
      datasets: [
        {
          label: 'Responses',
          backgroundColor: chartColors.slice(0, labels.length), // Use a slice of colors
          borderColor: chartColors.map(color => color.replace('0.8', '1')), // Darker border
          borderWidth: 1,
          data: data,
          percentages: percentages // Store percentages here for tooltip
        }
      ]
    };
  };
  
  onMounted(async () => {
    if (!surveyId.value) {
      error.value = 'Survey ID is missing from the route.';
      loading.value = false;
      return;
    }
  
    // No need to check authStore.token here, as the route has meta: { requiresAuth: true }
    // which should be handled by router guards.
  
    try {
      loading.value = true;
      error.value = null;
      
      const surveyDetailsPromise = surveyApi.getSurvey(surveyId.value);
      const analyticsDataPromise = surveyApi.getSurveyAnalytics(surveyId.value);

      const [surveyDetailsResponse, analyticsResponse] = await Promise.allSettled([
          surveyDetailsPromise,
          analyticsDataPromise
      ]);

      if (surveyDetailsResponse.status === 'fulfilled') {
          surveyTitle.value = surveyDetailsResponse.value.data.title;
      } else {
          console.error("Failed to fetch survey title:", surveyDetailsResponse.reason);
          surveyTitle.value = 'Survey Details Unavailable';
          // Optionally set a partial error if only title fails but analytics succeed
      }

      if (analyticsResponse.status === 'fulfilled') {
          analyticsData.value = analyticsResponse.value.data;
          console.log('Full analytics data received:', JSON.parse(JSON.stringify(analyticsData.value))); // Log the data
      } else {
          throw analyticsResponse.reason; // This will be caught by the main catch block
      }
  
    } catch (err) {
      console.error("Error fetching survey analytics:", err);
      if (err.response) {
        error.value = `Error ${err.response.status}: ${err.response.data.error || err.response.data.message || err.message}`;
      } else {
        error.value = `Failed to fetch survey analytics: ${err.message}`;
      }
    } finally {
      loading.value = false;
    }
  });
  
  </script>
  
  <style scoped>
  pre {
    background-color: #f0f0f0;
    padding: 15px;
    border-radius: 5px;
    white-space: pre-wrap;      
    white-space: -moz-pre-wrap; 
    white-space: -pre-wrap;     
    white-space: -o-pre-wrap;   
    word-wrap: break-word;      
  }
  .v-card-item {
      padding: 16px;
  }
  .border.rounded {
      border: 1px solid #e0e0e0 !important;
      border-radius: 4px !important;
  }
  </style>