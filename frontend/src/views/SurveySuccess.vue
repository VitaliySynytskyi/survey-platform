<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-12">
          <v-toolbar color="success" dark flat>
            <v-toolbar-title class="text-center">Submission Successful</v-toolbar-title>
          </v-toolbar>
          <v-card-text class="pa-5">
            <div class="text-center">
              <v-icon size="64" color="success" class="mb-4">mdi-check-circle-outline</v-icon>
              <h2 class="text-h5 mb-3">Thank You!</h2>
              <p class="mb-5">Your survey response has been successfully submitted.</p>
              <v-btn color="primary" to="/" class="mr-2">Go to Home</v-btn>
              <v-btn v-if="surveyId" :to="`/surveys/${surveyId}/analytics`" color="secondary">View Analytics</v-btn>
            </div>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { useRoute } from 'vue-router';
import { computed } from 'vue';

const route = useRoute();
const surveyId = computed(() => route.params.id);

</script>

<style scoped>
.fill-height {
  min-height: 80vh;
}
</style> 