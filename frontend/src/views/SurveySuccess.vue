<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12 rounded-xl">
          <v-card-text class="pa-6 text-center">
            <v-icon size="80" color="success" class="mb-4">mdi-check-circle-outline</v-icon>
            <h2 class="text-h4 font-weight-bold mb-3">Thank You!</h2>
            <p class="text-subtitle-1 mb-6 text-medium-emphasis">Your survey response has been successfully submitted.</p>
            
            <div class="d-flex flex-column flex-sm-row justify-center gap-3">
              <v-btn 
                color="primary" 
                to="/dashboard" 
                size="large" 
                variant="flat"
                rounded="pill"
                prepend-icon="mdi-view-dashboard"
                class="px-6 mb-2 mb-sm-0"
              >
                View All Surveys
              </v-btn>
              <v-btn 
                v-if="surveyId && canViewAnalytics" 
                :to="`/surveys/${surveyId}/analytics`" 
                color="secondary"
                size="large"
                variant="tonal"
                rounded="pill"
                prepend-icon="mdi-chart-bar"
                class="px-6"
              >
                View Analytics
              </v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router';
import { computed, onMounted, ref } from 'vue';
import { useAuthStore } from '../store/auth';
import { surveyApi } from '../services/api';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const surveyId = computed(() => route.params.id);
const isLoading = ref(false);
const isOwnSurvey = ref(false);

// Check if user is authenticated
const isAuthenticated = computed(() => authStore.isAuthenticated);

// Check if user owns this survey (or is admin)
const checkSurveyOwnership = async () => {
  if (!isAuthenticated.value || !surveyId.value) return false;
  
  isLoading.value = true;
  try {
    const response = await surveyApi.getSurvey(surveyId.value);
    const survey = response.data;
    
    // Check if current user is the creator or is admin
    isOwnSurvey.value = survey?.creator_id === authStore.user?.id || 
                        (authStore.user?.roles && authStore.user.roles.includes('admin'));
  } catch (error) {
    console.error('Error checking survey ownership:', error);
    isOwnSurvey.value = false;
  } finally {
    isLoading.value = false;
  }
};

// Show analytics button only for authenticated users who own the survey or are admins
const canViewAnalytics = computed(() => isAuthenticated.value && isOwnSurvey.value);

// Redirect to dashboard if coming from survey completion and user is authenticated
onMounted(async () => {
  if (isAuthenticated.value) {
    await checkSurveyOwnership();
  } else {
    // We could store a flag in localStorage to indicate they've just completed a survey
    localStorage.setItem('completedSurvey', 'true');
  }
});

</script>

<style scoped>
.fill-height {
  min-height: 80vh;
}
</style> 