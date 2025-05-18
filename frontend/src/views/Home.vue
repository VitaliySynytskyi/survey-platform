<template>
  <div class="home">
    <v-row justify="center" align="center">
      <v-col cols="12" md="8" class="text-center">
        <h1 class="text-h2 mb-6">Welcome to Survey Platform</h1>
        <p class="text-body-1 mb-12">
          Create and share surveys with ease. Collect valuable feedback from your target audience.
        </p>
        
        <v-row justify="center" class="mt-6">
          <v-col cols="12" md="6" class="d-flex justify-center">
            <v-btn
              to="/login"
              color="primary"
              x-large
              class="mx-2"
              v-if="!isAuthenticated"
            >
              Log In
            </v-btn>
            <v-btn
              to="/register"
              color="secondary"
              x-large
              class="mx-2"
              v-if="!isAuthenticated"
            >
              Register
            </v-btn>
            <v-btn
              to="/dashboard"
              color="primary"
              x-large
              v-if="isAuthenticated"
            >
              Go to Dashboard
            </v-btn>
          </v-col>
        </v-row>
      </v-col>
    </v-row>
    
    <v-row class="mt-16">
      <v-col cols="12" md="4">
        <v-card class="pa-4" height="100%">
          <v-card-title class="text-center">Create Surveys</v-card-title>
          <v-card-text class="text-center">
            Easily create surveys with multiple question types. Customize the survey to your needs.
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="4">
        <v-card class="pa-4" height="100%">
          <v-card-title class="text-center">Share with Respondents</v-card-title>
          <v-card-text class="text-center">
            Share your survey with a simple link. Collect responses from anyone, anywhere.
          </v-card-text>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="4">
        <v-card class="pa-4" height="100%">
          <v-card-title class="text-center">Analyze Results</v-card-title>
          <v-card-text class="text-center">
            View and analyze the collected responses. Get insights from the gathered data.
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useAuthStore } from '../store/auth';

export default {
  name: 'Home',
  setup() {
    const authStore = useAuthStore();
    const isAuthenticated = computed(() => authStore.isAuthenticated);
    
    return {
      isAuthenticated
    };
  }
};
</script> 