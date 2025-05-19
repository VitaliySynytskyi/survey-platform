<template>
  <v-card class="survey-card h-100 rounded-xl" elevation="2">
    <div class="status-indicator" :class="survey.is_active ? 'active' : 'inactive'"></div>
    
    <!-- Owner Badge -->
    <v-badge
      v-if="!isOwnSurvey"
      content="Created by others"
      color="info"
      offset-x="15"
      offset-y="15"
      location="top end"
    ></v-badge>
    
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
      <v-btn 
        variant="text" 
        color="secondary" 
        :to="`/surveys/${survey.id}/edit`" 
        size="small"
        :disabled="!canEdit"
        v-tooltip="!canEdit ? 'You can only edit your own surveys' : ''"
      >
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
          <v-list-item 
            @click="$emit('toggle-status', survey)"
            :disabled="!canEdit"
            v-if="canEdit"
          >
            <template v-slot:prepend>
              <v-icon>{{ survey.is_active ? 'mdi-close-circle' : 'mdi-check-circle' }}</v-icon>
            </template>
            <v-list-item-title>{{ survey.is_active ? 'Deactivate' : 'Activate' }}</v-list-item-title>
          </v-list-item>
          <v-divider v-if="canDelete"></v-divider>
          <v-list-item 
            @click="$emit('delete', survey)" 
            class="text-error"
            :disabled="!canDelete"
            v-if="canDelete"
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
</template>

<script>
export default {
  name: 'SurveyCard',
  props: {
    survey: {
      type: Object,
      required: true
    },
    isOwnSurvey: {
      type: Boolean,
      default: true
    },
    canEdit: {
      type: Boolean,
      default: true
    },
    canDelete: {
      type: Boolean,
      default: true
    }
  },
  emits: ['delete', 'toggle-status', 'share'],
  methods: {
    formatDate(dateString) {
      if (!dateString) return 'N/A';
      const options = { year: 'numeric', month: 'short', day: 'numeric' };
      try {
        return new Date(dateString).toLocaleDateString(undefined, options);
      } catch (e) {
        return dateString;
      }
    },
    copyShareLink(surveyId) {
      const link = `${window.location.origin}/surveys/${surveyId}`;
      navigator.clipboard.writeText(link).then(() => {
        this.$emit('share', {
          success: true,
          message: 'Survey share link copied to clipboard!'
        });
      }).catch(err => {
        console.error('Failed to copy share link:', err);
        this.$emit('share', {
          success: false,
          message: 'Failed to copy link. Please try again.'
        });
      });
    }
  }
};
</script>

<style scoped>
.survey-card {
  position: relative;
  overflow: hidden;
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
</style> 