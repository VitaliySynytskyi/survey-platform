<template>
  <div class="dashboard">
    <!-- Dashboard Header with Stats -->
    <v-card class="mb-6 pa-6 rounded-xl" elevation="2">
      <v-row>
        <v-col cols="12" md="8">
          <h1 class="text-h4 font-weight-bold mb-2">My Surveys Dashboard</h1>
          <p class="text-body-1 text-medium-emphasis">
            Manage your surveys, view responses, and analyze results.
          </p>
        </v-col>
        <v-col cols="12" md="4" class="d-flex align-center justify-end">
          <v-btn
            color="primary"
            prepend-icon="mdi-plus"
            to="/surveys/create"
            size="large"
            rounded="pill"
            variant="flat"
            class="px-6"
          >
            Create New Survey
          </v-btn>
        </v-col>
      </v-row>

      <!-- Statistics Cards -->
      <v-row class="mt-6">
        <v-col cols="12" sm="6" md="3">
          <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="primary">
            <div class="d-flex align-center mb-2">
              <v-avatar color="primary" size="36" class="mr-3">
                <v-icon color="white">mdi-poll</v-icon>
              </v-avatar>
              <span class="text-body-2 font-weight-medium">Total Surveys</span>
            </div>
            <div class="text-h3 font-weight-bold mt-2">{{ surveys.length }}</div>
          </v-card>
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="success">
            <div class="d-flex align-center mb-2">
              <v-avatar color="success" size="36" class="mr-3">
                <v-icon color="white">mdi-check-circle</v-icon>
              </v-avatar>
              <span class="text-body-2 font-weight-medium">Active Surveys</span>
            </div>
            <div class="text-h3 font-weight-bold mt-2">{{ activeSurveys }}</div>
          </v-card>
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="info">
            <div class="d-flex align-center mb-2">
              <v-avatar color="info" size="36" class="mr-3">
                <v-icon color="white">mdi-account-multiple</v-icon>
              </v-avatar>
              <span class="text-body-2 font-weight-medium">Total Responses</span>
            </div>
            <div class="text-h3 font-weight-bold mt-2">{{ totalResponses }}</div>
          </v-card>
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="secondary">
            <div class="d-flex align-center mb-2">
              <v-avatar color="secondary" size="36" class="mr-3">
                <v-icon color="white">mdi-calendar-today</v-icon>
              </v-avatar>
              <span class="text-body-2 font-weight-medium">Latest Survey</span>
            </div>
            <div class="text-subtitle-1 font-weight-bold mt-2 text-truncate">
              {{ latestSurveyTitle || 'N/A' }}
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-card>

    <!-- Filters and Search -->
    <v-card class="mb-6 pa-4 rounded-xl" elevation="2">
      <v-row>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            v-model="searchQuery"
            label="Search surveys"
            prepend-inner-icon="mdi-magnify"
            variant="outlined"
            density="comfortable"
            hide-details
            placeholder="Search by title or description"
          ></v-text-field>
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="statusFilter"
            label="Status"
            prepend-inner-icon="mdi-filter-variant"
            variant="outlined"
            density="comfortable"
            hide-details
            :items="[
              { title: 'All', value: 'all' },
              { title: 'Active', value: 'active' },
              { title: 'Inactive', value: 'inactive' }
            ]"
          ></v-select>
        </v-col>
        <v-col cols="12" sm="6" md="3">
          <v-select
            v-model="sortBy"
            label="Sort by"
            prepend-inner-icon="mdi-sort"
            variant="outlined"
            density="comfortable"
            hide-details
            :items="[
              { title: 'Newest first', value: 'newest' },
              { title: 'Oldest first', value: 'oldest' },
              { title: 'Title (A-Z)', value: 'title_asc' },
              { title: 'Title (Z-A)', value: 'title_desc' }
            ]"
          ></v-select>
        </v-col>
        <v-col cols="12" sm="6" md="2" class="d-flex align-center">
          <v-btn
            color="primary"
            variant="text"
            @click="resetFilters"
            prepend-icon="mdi-refresh"
            block
          >
            Reset
          </v-btn>
        </v-col>
      </v-row>
    </v-card>

    <!-- Loading State -->
    <v-row v-if="loading">
      <v-col cols="12" class="text-center py-12">
        <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        <div class="text-h6 mt-4">Loading your surveys...</div>
      </v-col>
    </v-row>

    <!-- Error State -->
    <v-row v-else-if="error">
      <v-col cols="12">
        <v-alert
          type="error"
          variant="tonal"
          border="start"
          closable
          class="rounded-lg"
        >
          <div class="d-flex align-center">
            <v-icon size="24" class="mr-3">mdi-alert-circle</v-icon>
            <div>
              <div class="text-subtitle-1 font-weight-bold">Error loading surveys</div>
              <div>{{ error }}</div>
            </div>
          </div>
          <template v-slot:append>
            <v-btn 
              color="error" 
              variant="text" 
              @click="fetchSurveys"
              prepend-icon="mdi-refresh"
            >
              Retry
            </v-btn>
          </template>
        </v-alert>
      </v-col>
    </v-row>

    <!-- Empty State -->
    <v-row v-else-if="filteredSurveys.length === 0 && !searchQuery && statusFilter === 'all'">
      <v-col cols="12">
        <v-card class="pa-12 text-center rounded-xl" elevation="0" variant="outlined">
          <v-avatar color="primary" size="72" class="mb-6">
            <v-icon size="36" color="white">mdi-poll</v-icon>
          </v-avatar>
          <v-card-title class="text-h4 mb-2">No surveys yet</v-card-title>
          <v-card-text>
            <p class="text-body-1 mb-6">You haven't created any surveys yet. Start by creating your first survey.</p>
          </v-card-text>
          <v-card-actions class="justify-center">
            <v-btn 
              color="primary" 
              prepend-icon="mdi-plus" 
              to="/surveys/create"
              rounded="pill"
              size="large"
              variant="flat"
              class="px-6"
            >
              Create Your First Survey
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- No Results for Search -->
    <v-row v-else-if="filteredSurveys.length === 0">
      <v-col cols="12">
        <v-card class="pa-8 text-center rounded-xl" elevation="0" variant="outlined">
          <v-avatar color="info" size="56" class="mb-4">
            <v-icon size="32" color="white">mdi-file-search</v-icon>
          </v-avatar>
          <v-card-title class="text-h5 mb-2">No matching surveys</v-card-title>
          <v-card-text>
            <p class="text-body-1 mb-4">No surveys match your current filters.</p>
          </v-card-text>
          <v-card-actions class="justify-center">
            <v-btn 
              color="primary" 
              variant="text" 
              @click="resetFilters"
              prepend-icon="mdi-refresh"
            >
              Reset Filters
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- Survey Grid -->
    <v-row v-else>
      <v-col v-for="survey in filteredSurveys" :key="survey.id" cols="12" sm="6" lg="4" class="survey-card-col">
        <v-card class="h-100 survey-card rounded-xl" elevation="2">
          <div class="status-indicator" :class="survey.is_active ? 'active' : 'inactive'"></div>
          <v-card-item>
            <template v-slot:prepend>
              <v-avatar color="primary" size="40" class="mr-3">
                <v-icon color="white">mdi-poll</v-icon>
              </v-avatar>
            </template>
            <v-card-title class="text-h5 mb-1 text-truncate">{{ survey.title }}</v-card-title>
            <v-card-subtitle>
              <div class="d-flex align-center">
                <v-icon size="16" class="mr-1">mdi-calendar</v-icon>
                <span>{{ formatDate(survey.created_at) }}</span>
              </div>
            </v-card-subtitle>
          </v-card-item>

          <v-card-text>
            <p class="mb-4 survey-description">{{ survey.description }}</p>
            <div class="d-flex flex-wrap gap-2 mb-2">
              <v-chip
                size="small"
                :color="survey.is_active ? 'success' : 'grey'"
                variant="tonal"
                class="text-caption"
              >
                {{ survey.is_active ? 'Active' : 'Inactive' }}
              </v-chip>
              <v-chip
                size="small"
                color="info"
                variant="tonal"
                class="text-caption"
              >
                {{ survey.questions_count }} questions
              </v-chip>
              <v-chip
                size="small"
                color="primary"
                variant="tonal"
                class="text-caption"
              >
                {{ survey.responses_count || 0 }} responses
              </v-chip>
            </div>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions class="pa-4">
            <v-btn variant="text" color="primary" :to="`/surveys/${survey.id}`" size="small">
              <v-icon>mdi-eye</v-icon>
            </v-btn>
            <v-btn variant="text" color="secondary" :to="`/surveys/${survey.id}/edit`" size="small">
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn variant="text" color="info" :to="`/surveys/${survey.id}/analytics`" size="small">
              <v-icon>mdi-chart-bar</v-icon>
            </v-btn>
            <v-btn variant="text" @click="copyShareLink(survey.id)" size="small">
              <v-icon>mdi-share-variant</v-icon>
            </v-btn>
            <v-spacer></v-spacer>
            <v-menu location="bottom end">
              <template v-slot:activator="{ props }">
                <v-btn icon v-bind="props" size="small">
                  <v-icon>mdi-dots-vertical</v-icon>
                </v-btn>
              </template>
              <v-list density="compact" min-width="200">
                <v-list-item :to="`/surveys/${survey.id}/responses`">
                  <template v-slot:prepend>
                    <v-icon>mdi-format-list-bulleted</v-icon>
                  </template>
                  <v-list-item-title>View Responses</v-list-item-title>
                </v-list-item>
                <v-list-item @click="toggleStatus(survey)">
                  <template v-slot:prepend>
                    <v-icon>{{ survey.is_active ? 'mdi-close-circle' : 'mdi-check-circle' }}</v-icon>
                  </template>
                  <v-list-item-title>{{ survey.is_active ? 'Deactivate' : 'Activate' }}</v-list-item-title>
                </v-list-item>
                <v-divider></v-divider>
                <v-list-item @click="confirmDelete(survey)" class="text-error">
                  <template v-slot:prepend>
                    <v-icon color="error">mdi-delete</v-icon>
                  </template>
                  <v-list-item-title class="text-error">Delete</v-list-item-title>
                </v-list-item>
              </v-list>
            </v-menu>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>

    <!-- Pagination -->
    <div class="d-flex justify-center mt-6" v-if="filteredSurveys.length > 0">
      <v-pagination
        v-model="page"
        :length="totalPages"
        :total-visible="5"
        rounded="circle"
      ></v-pagination>
    </div>

    <!-- Delete confirmation dialog -->
    <v-dialog v-model="deleteDialog" max-width="500px" class="rounded-xl">
      <v-card class="rounded-xl pa-2">
        <v-card-title class="text-h5 px-4 pt-4">Confirm Delete</v-card-title>
        <v-card-text class="px-4 pb-0">
          Are you sure you want to delete the survey "{{ selectedSurvey?.title }}"? This action cannot be undone.
        </v-card-text>
        <v-card-actions class="pa-4">
          <v-spacer></v-spacer>
          <v-btn color="grey" variant="text" @click="deleteDialog = false">Cancel</v-btn>
          <v-btn color="error" variant="flat" @click="deleteSurvey" :loading="deleteLoading">Delete</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <!-- Snackbar for notifications -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      :timeout="3000"
      location="top"
      rounded="pill"
    >
      <div class="d-flex align-center">
        <v-icon 
          class="mr-2" 
          :icon="snackbar.color === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle'"
        ></v-icon>
        {{ snackbar.text }}
      </div>
      <template v-slot:actions>
        <v-btn variant="text" icon="mdi-close" @click="snackbar.show = false"></v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<script>
import { ref, onMounted, computed, watch } from 'vue';
import { useAuthStore } from '../store/auth';
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
    
    // Filter and pagination state
    const searchQuery = ref('');
    const statusFilter = ref('all');
    const sortBy = ref('newest');
    const page = ref(1);
    const pageSize = ref(6); // Number of surveys per page

    const fetchSurveys = async () => {
      loading.value = true;
      error.value = '';
      try {
        const response = await surveyApi.getSurveys();
        surveys.value = response.data.map(survey => ({
          ...survey,
          responses_count: Math.floor(Math.random() * 100) // Placeholder until we have real data
        }));
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
      const link = `${window.location.origin}/surveys/${surveyId}`;
      navigator.clipboard.writeText(link).then(() => {
        showSnackbar('Survey share link copied to clipboard!', 'success');
      }).catch(err => {
        console.error('Failed to copy share link:', err);
        showSnackbar('Failed to copy link. Please try again.', 'error');
      });
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
        showSnackbar(`Survey "${selectedSurvey.value.title}" has been deleted.`, 'success');
        deleteDialog.value = false;
      } catch (err) {
        console.error('Error deleting survey:', err);
        showSnackbar('Failed to delete survey. Please try again.', 'error');
      } finally {
        deleteLoading.value = false;
        selectedSurvey.value = null;
      }
    };

    const toggleStatus = async (survey) => {
      try {
        const updatedSurvey = {
          ...survey,
          is_active: !survey.is_active
        };
        
        await surveyApi.updateSurvey(survey.id, updatedSurvey);
        
        // Update local state
        const index = surveys.value.findIndex(s => s.id === survey.id);
        if (index !== -1) {
          surveys.value[index].is_active = !survey.is_active;
        }
        
        showSnackbar(
          `Survey "${survey.title}" has been ${survey.is_active ? 'deactivated' : 'activated'}.`,
          'success'
        );
      } catch (err) {
        console.error('Error updating survey status:', err);
        showSnackbar('Failed to update survey status. Please try again.', 'error');
      }
    };

    const showSnackbar = (text, color = 'success') => {
      snackbar.value = {
        show: true,
        text,
        color
      };
    };

    // Filter and sort
    const filteredSurveys = computed(() => {
      let result = [...surveys.value];
      
      // Apply status filter
      if (statusFilter.value !== 'all') {
        const isActive = statusFilter.value === 'active';
        result = result.filter(survey => survey.is_active === isActive);
      }
      
      // Apply search filter
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase();
        result = result.filter(survey => 
          survey.title.toLowerCase().includes(query) || 
          (survey.description && survey.description.toLowerCase().includes(query))
        );
      }
      
      // Apply sorting
      result.sort((a, b) => {
        switch (sortBy.value) {
          case 'newest':
            return new Date(b.created_at) - new Date(a.created_at);
          case 'oldest':
            return new Date(a.created_at) - new Date(b.created_at);
          case 'title_asc':
            return a.title.localeCompare(b.title);
          case 'title_desc':
            return b.title.localeCompare(a.title);
          default:
            return 0;
        }
      });
      
      return result;
    });

    // Apply pagination
    const paginatedSurveys = computed(() => {
      const start = (page.value - 1) * pageSize.value;
      const end = start + pageSize.value;
      return filteredSurveys.value.slice(start, end);
    });

    const totalPages = computed(() => {
      return Math.ceil(filteredSurveys.value.length / pageSize.value);
    });

    const resetFilters = () => {
      searchQuery.value = '';
      statusFilter.value = 'all';
      sortBy.value = 'newest';
      page.value = 1;
    };

    // Statistics getters
    const activeSurveys = computed(() => {
      return surveys.value.filter(survey => survey.is_active).length;
    });

    const totalResponses = computed(() => {
      return surveys.value.reduce((total, survey) => total + (survey.responses_count || 0), 0);
    });

    const latestSurveyTitle = computed(() => {
      if (surveys.value.length === 0) return null;
      
      const sorted = [...surveys.value].sort(
        (a, b) => new Date(b.created_at) - new Date(a.created_at)
      );
      
      return sorted[0]?.title;
    });

    // Reset page when filters change
    watch([searchQuery, statusFilter, sortBy], () => {
      page.value = 1;
    });

    onMounted(() => {
      fetchSurveys();
    });

    return {
      surveys,
      loading,
      error,
      deleteDialog,
      selectedSurvey,
      deleteLoading,
      snackbar,
      searchQuery,
      statusFilter,
      sortBy,
      page,
      filteredSurveys: paginatedSurveys,
      totalPages,
      activeSurveys,
      totalResponses,
      latestSurveyTitle,
      formatDate,
      copyShareLink,
      confirmDelete,
      deleteSurvey,
      toggleStatus,
      resetFilters,
      fetchSurveys
    };
  }
};
</script>

<style scoped>
.dashboard {
  animation: fadeIn 0.5s ease-out;
}

.survey-card-col {
  transition: transform 0.3s ease-out;
}

.survey-card {
  position: relative;
  overflow: hidden;
  height: 100%;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.survey-card:hover {
  transform: translateY(-8px);
}

.survey-description {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.status-indicator {
  position: absolute;
  top: 0;
  right: 0;
  width: 100%;
  height: 4px;
}

.status-indicator.active {
  background-color: var(--v-success-base);
}

.status-indicator.inactive {
  background-color: var(--v-grey-base);
}

.stat-card {
  border-radius: 12px;
  transition: transform 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style> 