<template>
  <div class="dashboard">
    <v-row>
      <v-col cols="12">
        <v-card class="mb-4 pa-4">
          <div class="d-flex align-center justify-space-between">
            <h1 class="text-h4">My Surveys</h1>
            <v-btn color="primary" prepend-icon="mdi-plus" to="/surveys/create">
              Create New Survey
            </v-btn>
          </div>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-if="loading">
      <v-col cols="12" class="text-center">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
      </v-col>
    </v-row>

    <v-row v-else-if="error">
      <v-col cols="12">
        <v-alert type="error" closable>
          {{ error }}
        </v-alert>
      </v-col>
    </v-row>

    <v-row v-else-if="surveys.length === 0">
      <v-col cols="12">
        <v-card class="pa-6 text-center">
          <v-card-title class="text-h5 mb-2">No surveys yet</v-card-title>
          <v-card-text>
            <p>You haven't created any surveys yet. Start by creating your first survey.</p>
          </v-card-text>
          <v-card-actions class="justify-center">
            <v-btn color="primary" prepend-icon="mdi-plus" to="/surveys/create">
              Create Your First Survey
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <v-row v-else>
      <v-col v-for="survey in surveys" :key="survey.id" cols="12" md="6" lg="4">
        <v-card class="h-100 survey-card">
          <v-card-title class="text-h5">{{ survey.title }}</v-card-title>
          <v-card-subtitle>
            Created: {{ formatDate(survey.created_at) }}
          </v-card-subtitle>

          <v-card-text>
            <p class="mb-2">{{ survey.description }}</p>
            <v-chip class="mr-2" :color="survey.is_active ? 'success' : 'grey'">
              {{ survey.is_active ? 'Active' : 'Inactive' }}
            </v-chip>
            <v-chip color="info">
              {{ survey.questions_count }} questions
            </v-chip>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions>
            <v-btn color="primary" variant="text" :to="`/surveys/${survey.id}`">
              <v-icon left>mdi-eye</v-icon> View
            </v-btn>
            <v-btn color="secondary" variant="text" :to="`/surveys/${survey.id}/edit`">
              <v-icon left>mdi-pencil</v-icon> Edit
            </v-btn>
            <v-btn color="info" variant="text" :to="`/surveys/${survey.id}/analytics`">
              <v-icon left>mdi-chart-bar</v-icon> Analytics
            </v-btn>
            <v-spacer></v-spacer>
            <v-menu>
              <template v-slot:activator="{ props }">
                <v-btn icon v-bind="props">
                  <v-icon>mdi-dots-vertical</v-icon>
                </v-btn>
              </template>
              <v-list>
                <v-list-item :to="`/surveys/${survey.id}/responses`">
                  <v-list-item-title>
                    <v-icon left>mdi-chart-bar</v-icon>
                    View Responses
                  </v-list-item-title>
                </v-list-item>
                <v-list-item @click="copyShareLink(survey.id)">
                  <v-list-item-title>
                    <v-icon left>mdi-share-variant</v-icon>
                    Share Survey
                  </v-list-item-title>
                </v-list-item>
                <v-list-item @click="toggleStatus(survey)">
                  <v-list-item-title>
                    <v-icon left>{{ survey.is_active ? 'mdi-close-circle' : 'mdi-check-circle' }}</v-icon>
                    {{ survey.is_active ? 'Deactivate' : 'Activate' }}
                  </v-list-item-title>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item @click="confirmDelete(survey)" class="text-error">
                  <v-list-item-title class="text-error">
                    <v-icon left color="error">mdi-delete</v-icon>
                    Delete
                  </v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- Delete confirmation dialog -->
    <v-dialog v-model="deleteDialog" max-width="500px">
      <v-card>
        <v-card-title class="text-h5">Confirm Delete</v-card-title>
        <v-card-text>
          Are you sure you want to delete the survey "{{ selectedSurvey?.title }}"? This action cannot be undone.
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" variant="text" @click="deleteSurvey" :loading="deleteLoading">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

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
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../store/auth';
import axios from '../utils/axiosConfig';
import { surveyApi } from '../services/api';
import { useRouter } from 'vue-router';

export default {
  name: 'Dashboard',
  setup() {
    const authStore = useAuthStore();
    const surveys = ref([]);
    const loading = ref(false);
    const error = ref('');
    const deleteDialog = ref(false);
    const selectedSurvey = ref(null);
    const deleteLoading = ref(false);
    const snackbar = ref({
      show: false,
      text: '',
      color: 'success'
    });
    const router = useRouter();

    const fetchSurveys = async () => {
      loading.value = true;
      error.value = '';
      try {
        const response = await surveyApi.getSurveys();
        surveys.value = response.data;
      } catch (err) {
        console.error('Error fetching surveys:', err);
        const errorMessage = err.response?.data?.error || 'Failed to load surveys. Please try again later.';
        if (err.response?.status === 401) {
            error.value = 'Your session may have expired. Please try logging out and logging back in.';
        } else {
            error.value = errorMessage;
        }
      } finally {
        loading.value = false;
      }
    };

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A';
      const options = { year: 'numeric', month: 'short', day: 'numeric' };
      try {
        return new Date(dateString).toLocaleDateString(undefined, options);
      } catch (e) {
        return dateString;
      }
    };

    const copyShareLink = (surveyId) => {
      const link = `${window.location.origin}/surveys/take/${surveyId}`;
      navigator.clipboard.writeText(link).then(() => {
        showSnackbar('Survey share link copied to clipboard!', 'success');
      }).catch(err => {
        console.error('Failed to copy share link:', err);
        showSnackbar('Failed to copy link. Please try again.', 'error');
      });
    };

    const toggleStatus = async (survey) => {
      const originalStatus = survey.is_active;
      const newStatus = !survey.is_active;
      survey.is_active = newStatus;
      
      try {
        await surveyApi.updateSurveyStatus(survey.id, newStatus);
        showSnackbar(
          `Survey "${survey.title}" ${newStatus ? 'activated' : 'deactivated'} successfully.`,
          'success'
        );
      } catch (err) {
        survey.is_active = originalStatus;
        console.error('Error updating survey status:', err);
        const errMessage = err.response?.data?.error || 'Failed to update survey status.';
        showSnackbar(errMessage, 'error');
      }
    };

    const confirmDelete = (survey) => {
      selectedSurvey.value = survey;
      deleteDialog.value = true;
    };

    const deleteSurvey = async () => {
      if (!selectedSurvey.value) return;
      deleteLoading.value = true;
      try {
        await surveyApi.deleteSurvey(selectedSurvey.value.id);
        surveys.value = surveys.value.filter(s => s.id !== selectedSurvey.value.id);
        showSnackbar(`Survey "${selectedSurvey.value.title}" deleted successfully.`, 'success');
        deleteDialog.value = false;
      } catch (err) {
        console.error('Error deleting survey:', err);
        const errMessage = err.response?.data?.error || 'Failed to delete survey.';
        showSnackbar(errMessage, 'error');
      } finally {
        deleteLoading.value = false;
        selectedSurvey.value = null;
      }
    };

    const showSnackbar = (text, color = 'success') => {
      snackbar.value = {
        show: true,
        text,
        color
      };
    };

    const navigateToAnalytics = (surveyId) => {
      router.push({ name: 'SurveyAnalytics', params: { id: surveyId } });
    };

    onMounted(async () => {
      if (authStore.token) {
        loading.value = true;
        error.value = '';
        try {
          await authStore.fetchUser();
          if (authStore.isAuthenticated) {
            await fetchSurveys();
          } else {
            error.value = 'Your session is invalid. Please log in again.';
          }
        } catch (err) {
          console.error('Dashboard setup error:', err);
          if (!error.value) {
             error.value = 'Failed to initialize dashboard. Please try logging in again.';
          }
        } finally {
          loading.value = false;
        }
      } else {
        error.value = 'You are not authenticated. Please log in.';
        loading.value = false; 
      }
    });

    return {
      surveys,
      loading,
      error,
      formatDate,
      copyShareLink,
      toggleStatus,
      confirmDelete,
      deleteDialog,
      selectedSurvey,
      deleteSurvey,
      deleteLoading,
      snackbar,
      showSnackbar,
      authStore,
      navigateToAnalytics
    };
  }
};
</script>

<style scoped>
.survey-card {
  transition: transform 0.3s, box-shadow 0.3s;
}

.survey-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.15) !important;
}
</style> 