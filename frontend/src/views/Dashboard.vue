<template>
  <div class="dashboard">
    <!-- Dashboard Header with Stats -->
    <v-card class="mb-6 pa-6 rounded-xl" elevation="2">
      <v-row>
        <v-col cols="12" md="8">
          <h1 class="text-h4 font-weight-bold mb-2">Surveys Dashboard</h1>
          <p class="text-body-1 text-medium-emphasis">
            Browse surveys, manage your own, view responses, and analyze results.
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

      <!-- Survey View Tabs -->
      <v-tabs
        v-model="viewTab"
        color="primary"
        align-tabs="center"
        class="mt-4"
        density="comfortable"
      >
        <v-tab value="my-surveys" prepend-icon="mdi-account">My Surveys</v-tab>
        <v-tab value="all-surveys" prepend-icon="mdi-earth">All Surveys</v-tab>
      </v-tabs>

      <!-- Statistics Cards -->
      <v-row class="mt-6">
        <v-col cols="12" sm="6" md="3">
          <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="primary">
            <div class="d-flex align-center mb-2">
              <v-avatar color="primary" size="36" class="mr-3">
                <v-icon color="white">mdi-poll</v-icon>
              </v-avatar>
              <span class="text-body-2 font-weight-medium">Total {{ viewTab === 'my-surveys' ? 'My' : 'All' }} Surveys</span>
            </div>
            <div class="text-h3 font-weight-bold mt-2">{{ totalSurveysCount }}</div>
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
              <div class="text-body-2 font-weight-medium">Total Responses</div>
            </div>
            <div class="text-h3 font-weight-bold mt-2">
              <div v-if="loadingResponseCounts">
                <v-progress-circular indeterminate size="24" color="white" class="mr-2"></v-progress-circular>
                <span class="text-body-1">Loading...</span>
              </div>
              <span v-else>{{ calculateTotalResponses() }}</span>
            </div>
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
        <div class="text-h6 mt-4">Loading surveys...</div>
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

    <!-- Empty State - My Surveys -->
    <v-row v-else-if="viewTab === 'my-surveys' && filteredSurveys.length === 0 && !searchQuery && statusFilter === 'all'">
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

    <!-- Empty State - All Surveys -->
    <v-row v-else-if="viewTab === 'all-surveys' && filteredSurveys.length === 0 && !searchQuery && statusFilter === 'all'">
      <v-col cols="12">
        <v-card class="pa-12 text-center rounded-xl" elevation="0" variant="outlined">
          <v-avatar color="primary" size="72" class="mb-6">
            <v-icon size="36" color="white">mdi-poll</v-icon>
          </v-avatar>
          <v-card-title class="text-h4 mb-2">No surveys available</v-card-title>
          <v-card-text>
            <p class="text-body-1 mb-6">There are no surveys available in the system yet.</p>
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
              Create First Survey
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
          
          <!-- Owner Badge -->
          <div v-if="!isOwnSurvey(survey)" class="ownership-badge">
            <v-chip
              color="info"
              size="small"
              label
              class="font-weight-medium"
            >
              Created by others
            </v-chip>
          </div>
          
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
                <v-chip 
                  v-if="survey.created_by_username" 
                  class="ml-2" 
                  size="x-small" 
                  color="grey"
                  label
                >
                  {{ survey.created_by_username }}
                </v-chip>
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
                class="text-caption d-flex align-center"
              >
                <v-icon size="16" start>mdi-message-reply-text</v-icon>
                <v-progress-circular v-if="loadingResponseCounts" indeterminate size="12" width="2" color="primary" class="mr-1"></v-progress-circular>
                <span class="font-weight-medium ml-1">{{ survey.responses_count || 0 }} responses</span>
              </v-chip>
            </div>
          </v-card-text>

          <v-divider></v-divider>

          <v-card-actions class="pa-4">
            <v-btn variant="text" color="primary" :to="`/surveys/${survey.id}`" size="small">
              <v-icon>mdi-eye</v-icon>
            </v-btn>
            <v-btn 
              variant="text" 
              color="secondary" 
              :to="`/surveys/${survey.id}/edit`" 
              size="small"
              :disabled="!canEdit(survey)"
              v-tooltip="!canEdit(survey) ? 'You can only edit your own surveys' : null"
            >
              <v-icon>mdi-pencil</v-icon>
            </v-btn>
            <v-btn 
              variant="text" 
              color="info" 
              :to="`/surveys/${survey.id}/analytics`" 
              size="small"
              :disabled="!canViewAnalytics(survey)"
              v-tooltip="!canViewAnalytics(survey) ? 'You can only view analytics for your own surveys' : null"
            >
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
                <v-list-item 
                  :to="`/surveys/${survey.id}/responses`" 
                  :disabled="!canViewResponses(survey)"
                  v-tooltip="!canViewResponses(survey) ? 'You can only view responses for your own surveys' : null"
                >
                  <template v-slot:prepend>
                    <v-icon>mdi-format-list-bulleted</v-icon>
                  </template>
                  <v-list-item-title>View Responses</v-list-item-title>
                </v-list-item>
                <v-list-item 
                  @click="toggleStatus(survey)" 
                  :disabled="!canEdit(survey)"
                  v-if="canEdit(survey)"
                >
                  <template v-slot:prepend>
                    <v-icon>{{ survey.is_active ? 'mdi-close-circle' : 'mdi-check-circle' }}</v-icon>
                  </template>
                  <v-list-item-title>{{ survey.is_active ? 'Deactivate' : 'Activate' }}</v-list-item-title>
                </v-list-item>
                <v-divider v-if="canDelete(survey)"></v-divider>
                <v-list-item 
                  @click="confirmDelete(survey)" 
                  class="text-error"
                  :disabled="!canDelete(survey)"
                  v-if="canDelete(survey)"
                >
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
    const loadingResponseCounts = ref(false);
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
    
    // View mode and filter state
    const viewTab = ref('my-surveys');
    const searchQuery = ref('');
    const statusFilter = ref('all');
    const sortBy = ref('newest');
    const page = ref(1);
    const pageSize = ref(6); // Number of surveys per page
    const totalSurveysCount = ref(0); // For server-side pagination

    // Check if the current user has admin role
    const isAdmin = computed(() => {
      return authStore.user?.role === 'admin';
    });

    // Check if a survey is owned by the current user
    const isOwnSurvey = (survey) => {
      console.log('isOwnSurvey check:', 
        'Survey Creator ID:', survey.creator_id, 'Type:', typeof survey.creator_id, 
        'Auth User ID:', authStore.user?.id, 'Type:', typeof authStore.user?.id, 
        'Comparison Result:', survey.creator_id === authStore.user?.id
      );
      return survey.creator_id === authStore.user?.id;
    };

    // Determine if user can edit a survey (owner or admin)
    const canEdit = (survey) => {
      return isOwnSurvey(survey) || isAdmin.value;
    };

    // Determine if user can delete a survey (owner or admin)
    const canDelete = (survey) => {
      return isOwnSurvey(survey) || isAdmin.value;
    };

    // Determine if user can view analytics (owner or admin)
    const canViewAnalytics = (survey) => {
      return isOwnSurvey(survey) || isAdmin.value;
    };

    // Determine if user can view responses (owner or admin)
    const canViewResponses = (survey) => {
      return isOwnSurvey(survey) || isAdmin.value;
    };

    const fetchSurveys = async () => {
      loading.value = true;
      error.value = '';
      try {
        let response;
        const apiParams = {
          page: page.value,
          limit: pageSize.value,
          include_responses: true // Explicitly request response counts
        };
        // Future improvement: Pass filter/sort params to backend if supported
        // if (searchQuery.value) apiParams.search = searchQuery.value;
        // if (statusFilter.value !== 'all') apiParams.status = statusFilter.value;
        // if (sortBy.value) apiParams.sort_by = mapSortKey(sortBy.value);
        
        if (viewTab.value === 'my-surveys') {
          response = await surveyApi.getUserSurveys(apiParams);
        } else {
          response = await surveyApi.getAllSurveys(apiParams);
        }
        
        console.log('API Response Data:', response.data); // For debugging response structure

        if (response.data && Array.isArray(response.data.data) && typeof response.data.total === 'number') {
          surveys.value = response.data.data.map(survey => ({
            ...survey,
            responses_count: typeof survey.responses_count === 'number' ? survey.responses_count : 0
          }));
          totalSurveysCount.value = response.data.total;
          
          // Fetch individual response counts if not included in the survey data
          if (!surveys.value.some(survey => survey.responses_count > 0) && surveys.value.length > 0) {
            await fetchResponseCounts();
          }
        } else if (response.data && Array.isArray(response.data.items) && typeof response.data.total === 'number') {
          // Alternative common structure { items: [], total: X }
          surveys.value = response.data.items.map(survey => ({
            ...survey,
            responses_count: typeof survey.responses_count === 'number' ? survey.responses_count : 0
          }));
          totalSurveysCount.value = response.data.total;
          
          // Fetch individual response counts if not included in the survey data
          if (!surveys.value.some(survey => survey.responses_count > 0) && surveys.value.length > 0) {
            await fetchResponseCounts();
          }
        } else {
          console.error('Unexpected API response structure for surveys:', response.data);
          surveys.value = [];
          totalSurveysCount.value = 0;
          error.value = 'Failed to load surveys due to unexpected data format from server.';
        }
        
        // Debug response counts
        console.log('Survey data with response counts:', surveys.value);

      } catch (err) {
        console.error('Error fetching surveys:', err);
        const errorMessage = err.response?.data?.error || 'Failed to load surveys. Please try again later.';
        if (err.response?.status === 401) {
            error.value = 'Your session may have expired. Please try logging out and logging back in.';
        } else {
            error.value = errorMessage;
        }
        surveys.value = []; // Clear surveys on error
        totalSurveysCount.value = 0; // Reset count on error
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
      if (!canDelete(survey)) {
        showSnackbar('You do not have permission to delete this survey.', 'error');
        return;
      }
      
      selectedSurvey.value = survey;
      deleteDialog.value = true;
    };

    const deleteSurvey = async () => {
      if (!selectedSurvey.value) return;
      
      if (!canDelete(selectedSurvey.value)) {
        showSnackbar('You do not have permission to delete this survey.', 'error');
        deleteDialog.value = false;
        return;
      }
      
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
      if (!canEdit(survey)) {
        showSnackbar('You do not have permission to edit this survey.', 'error');
        return;
      }
      
      try {
        // Use the dedicated status update endpoint instead of full survey update
        await surveyApi.updateSurveyStatus(survey.id, !survey.is_active);
        
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

    // Filter and sort (now operates on the current page of data from `surveys.value`)
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
    // const paginatedSurveys = computed(() => { // This is no longer needed with server-side pagination
    //   const start = (page.value - 1) * pageSize.value;
    //   const end = start + pageSize.value;
    //   return filteredSurveys.value.slice(start, end);
    // });

    const totalPages = computed(() => {
      if (pageSize.value === 0) return 0;
      return Math.ceil(totalSurveysCount.value / pageSize.value);
    });

    const resetFilters = () => {
      searchQuery.value = '';
      statusFilter.value = 'all';
      sortBy.value = 'newest';
      if (page.value !== 1) {
        page.value = 1;
      } else {
        fetchSurveys();
      }
    };

    // Statistics getters
    const activeSurveys = computed(() => {
      return surveys.value.filter(survey => survey.is_active).length;
    });

    const calculateTotalResponses = () => {
      return surveys.value.reduce((total, survey) => total + (survey.responses_count || 0), 0);
    };
    
    // Fetch response counts for each survey individually
    const fetchResponseCounts = async () => {
      if (surveys.value.length === 0) return;
      
      loadingResponseCounts.value = true;
      try {
        // Create an array of promises to fetch response counts in parallel
        const countPromises = surveys.value.map(async (survey) => {
          try {
            const response = await surveyApi.getSurveyResponses(survey.id, { count_only: true });
            if (response.data && typeof response.data.count === 'number') {
              return {
                surveyId: survey.id,
                count: response.data.count
              };
            }
            return null;
          } catch (err) {
            console.error(`Failed to fetch response count for survey ${survey.id}:`, err);
            return null;
          }
        });
        
        // Wait for all requests to complete
        const results = await Promise.all(countPromises);
        
        // Update survey response counts
        results.forEach(result => {
          if (result) {
            const surveyIndex = surveys.value.findIndex(s => s.id === result.surveyId);
            if (surveyIndex !== -1) {
              surveys.value[surveyIndex].responses_count = result.count;
            }
          }
        });
        
        console.log('Updated survey response counts:', surveys.value.map(s => ({ 
          id: s.id, 
          title: s.title, 
          responses: s.responses_count 
        })));
      } catch (err) {
        console.error('Error fetching response counts:', err);
      } finally {
        loadingResponseCounts.value = false;
      }
    };

    const latestSurveyTitle = computed(() => {
      if (surveys.value.length === 0) return null;
      
      const sorted = [...surveys.value].sort(
        (a, b) => new Date(b.created_at) - new Date(a.created_at)
      );
      
      return sorted[0]?.title;
    });

    // Reset page when filters change
    // watch([searchQuery, statusFilter, sortBy], () => { // OLD logic
    //   page.value = 1;
    //   // fetchSurveys(); // Fetch surveys after filter changes and page is reset to 1
    // });

    // // Refetch surveys when view mode changes // OLD logic
    // watch(viewTab, () => {
    //   page.value = 1; // Reset to page 1 when tab changes
    //   fetchSurveys();
    //   resetFilters(); 
    // });
    
    // // Refetch surveys when page changes due to pagination component // OLD logic
    // watch(page, (newPage, oldPage) => {
    //   if (newPage !== oldPage) {
    //     fetchSurveys();
    //   }
    // }, { immediate: false }); 

    // ----- NEW REVISED WATCHERS AND RELATED LOGIC -----

    // Handles pagination clicks or programmatic page changes
    watch(page, (newPage, oldPage) => {
      if (newPage !== oldPage) {
        fetchSurveys().then(() => {
          // Refresh response counts when page changes
          if (surveys.value.length > 0) {
            fetchResponseCounts();
          }
        });
      }
    }, { immediate: false }); // immediate: false is important if onMounted also fetches

    // Watch for changes in individual filters or sort
    watch([searchQuery, statusFilter, sortBy], () => {
      // When a filter changes, we always want to go to page 1 and refetch.
      if (page.value !== 1) {
        page.value = 1; // Let page watcher handle the fetch
      } else {
        fetchSurveys().then(() => {
          // Refresh response counts when filters change
          if (surveys.value.length > 0) {
            fetchResponseCounts();
          }
        });
      }
    });

    // Watch for tab changes
    watch(viewTab, () => {
      // When tab changes, reset filters and fetch.
      // This means setting page to 1 as well.
      searchQuery.value = ''; // Manually reset filters tied to the tab change
      statusFilter.value = 'all';
      sortBy.value = 'newest';
      if (page.value !== 1) {
        page.value = 1; // Let page watcher handle the fetch
      } else {
        fetchSurveys().then(() => {
          // Refetch response counts when the view tab changes
          if (surveys.value.length > 0) {
            fetchResponseCounts();
          }
        });
      }
    });
    // ----- END NEW REVISED WATCHERS ----- 

    // Call fetchResponseCounts after initial surveys are loaded
    onMounted(async () => {
      if (!authStore.user) {
        try {
          await authStore.fetchUser();
        } catch (error) {
          console.error("Failed to fetch user details", error);
        }
      }
      
      await fetchSurveys();
      
      // Explicitly fetch response counts after initial surveys load
      if (surveys.value.length > 0) {
        fetchResponseCounts();
      }
    });

    return {
      surveys, // surveys ref now holds current page data
      loading,
      loadingResponseCounts,
      error,
      deleteDialog,
      selectedSurvey,
      deleteLoading,
      snackbar,
      viewTab,
      searchQuery,
      statusFilter,
      sortBy,
      page, // Bound to v-pagination
      // filteredSurveys: paginatedSurveys, // OLD: This was client-side pagination
      filteredSurveys: filteredSurveys, // NEW: This is the computed prop for filtered/sorted current page data
      totalPages, // Driven by totalSurveysCount from server
      activeSurveys, // Note: This and other stats are now for the current page
      calculateTotalResponses,
      latestSurveyTitle,
      totalSurveysCount, // Expose for the template (total surveys card)
      isAdmin,
      isOwnSurvey,
      canEdit,
      canDelete,
      canViewAnalytics,
      canViewResponses,
      formatDate,
      copyShareLink,
      confirmDelete,
      deleteSurvey,
      toggleStatus,
      resetFilters,
      fetchSurveys,
      fetchResponseCounts,
      showSnackbar
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

.survey-badge {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 1;
}

.ownership-badge {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
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