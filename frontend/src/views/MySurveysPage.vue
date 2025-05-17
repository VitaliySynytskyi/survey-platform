<template>
  <div class="my-surveys-page">
    <div class="container py-4">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1 class="mb-0">My Surveys</h1>
        <router-link to="/surveys/create" class="btn btn-primary">
          <i class="bi bi-plus-lg me-2"></i>Create New Survey
        </router-link>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="d-flex flex-column align-items-center justify-content-center py-5">
        <div class="spinner-border text-primary mb-3" role="status">
          <span class="visually-hidden">Loading...</span>
        </div>
        <p class="text-muted">Loading your surveys...</p>
      </div>

      <!-- Error State -->
      <div v-else-if="error" class="alert alert-danger shadow-sm" role="alert">
        <h4 class="alert-heading mb-2">
          <i class="bi bi-exclamation-triangle-fill me-2"></i>Error
        </h4>
        <p>{{ error }}</p>
        <hr>
        <div class="d-flex">
          <button @click="fetchSurveys" class="btn btn-outline-danger">
            <i class="bi bi-arrow-repeat me-2"></i>Try Again
          </button>
          <router-link to="/login" v-if="!authStore.isAuthenticated" class="btn btn-primary ms-2">
            Log In
          </router-link>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="surveys.length === 0 && !error" class="empty-state text-center py-5 bg-light rounded shadow-sm">
        <div class="empty-icon mb-4">
          <i class="bi bi-clipboard-plus" style="font-size: 4rem; color: #6c757d;"></i>
        </div>
        <h3>You haven't created any surveys yet</h3>
        <p class="text-muted mb-4">Get started by creating your first survey</p>
        <router-link to="/surveys/create" class="btn btn-primary btn-lg">
          <i class="bi bi-plus-lg me-2"></i>Create Your First Survey
        </router-link>
      </div>

      <!-- Survey List -->
      <div v-else class="row">
        <div v-for="survey in surveys" :key="survey.id" class="col-md-6 col-lg-4 mb-4">
          <div class="card survey-card h-100 shadow-sm">
            <div class="card-header bg-white border-bottom">
              <div class="d-flex justify-content-between align-items-center">
                <span class="badge rounded-pill" :class="getStatusClass(survey)">{{ getStatusText(survey) }}</span>
                <div class="dropdown">
                  <button class="btn btn-sm btn-light" type="button" data-bs-toggle="dropdown" aria-expanded="false">
                    <i class="bi bi-three-dots-vertical"></i>
                  </button>
                  <ul class="dropdown-menu dropdown-menu-end">
                    <li>
                      <router-link :to="`/surveys/${survey.id}/edit`" class="dropdown-item">
                        <i class="bi bi-pencil me-2"></i>Edit
                      </router-link>
                    </li>
                    <li>
                      <router-link :to="`/surveys/${survey.id}/results`" class="dropdown-item">
                        <i class="bi bi-bar-chart me-2"></i>Results
                      </router-link>
                    </li>
                    <li><hr class="dropdown-divider"></li>
                    <li>
                      <a href="#" @click.prevent="confirmDeleteSurvey(survey.id)" class="dropdown-item text-danger">
                        <i class="bi bi-trash me-2"></i>Delete
                      </a>
                    </li>
                  </ul>
                </div>
              </div>
            </div>
            <div class="card-body d-flex flex-column">
              <h5 class="card-title mb-2">{{ survey.title }}</h5>
              <p class="card-text text-muted mb-3 flex-grow-1">
                {{ survey.description || 'No description provided' }}
              </p>
              <div class="survey-meta small text-muted mb-3">
                <div class="meta-item mb-1">
                  <i class="bi bi-calendar me-2"></i>Created: {{ formatDate(survey.created_at) }}
                </div>
                <div class="meta-item mb-1" v-if="survey.updated_at">
                  <i class="bi bi-pencil me-2"></i>Updated: {{ formatDate(survey.updated_at) }}
                </div>
                <div class="meta-item" v-if="survey.questions?.length">
                  <i class="bi bi-list-check me-2"></i>Questions: {{ survey.questions.length }}
                </div>
              </div>
            </div>
            <div class="card-footer bg-white border-top d-flex justify-content-between">
              <router-link :to="`/surveys/${survey.id}/take`" class="btn btn-primary">
                <i class="bi bi-pencil-square me-1"></i>Take
              </router-link>
              <button class="btn btn-outline-secondary ms-2" @click="copyShareLink(survey.id)">
                <i class="bi bi-share me-1"></i>Share
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, Ref } from 'vue';
import { useRouter } from 'vue-router';
import surveyService from '@/services/survey';
import { useAuthStore } from '@/stores/auth';
import type { Survey } from '@/types/survey';

interface ExtendedSurvey extends Survey {
  is_published?: boolean;
  is_draft?: boolean;
}

const surveys: Ref<ExtendedSurvey[]> = ref([]);
const loading: Ref<boolean> = ref(true);
const error: Ref<string | null> = ref(null);
const authStore = useAuthStore();
const router = useRouter();

const fetchSurveys = async () => {
  if (!authStore.user?.id) {
    error.value = 'Please log in to view your surveys';
    loading.value = false;
    return;
  }
  try {
    loading.value = true;
    const response = await surveyService.getUserSurveys(authStore.user.id, 1, 100);
    surveys.value = response.surveys;
  } catch (err: any) {
    console.error('Failed to fetch surveys:', err);
    error.value = err.response?.data?.error || err.message || 'Failed to load surveys';
    if (err.response?.status === 401 || err.response?.status === 403) {
        router.push('/login');
    }
  } finally {
    loading.value = false;
  }
};

const formatDate = (dateString?: string): string => {
  if (!dateString) return 'N/A';
  return new Date(dateString).toLocaleDateString();
};

const confirmDeleteSurvey = async (surveyId: string) => {
  if (window.confirm('Are you sure you want to delete this survey? This action cannot be undone.')) {
    try {
      await surveyService.deleteSurvey(surveyId);
      surveys.value = surveys.value.filter((s: Survey) => s.id !== surveyId);
    } catch (err: any) {
      console.error('Failed to delete survey:', err);
      alert(err.response?.data?.error || 'Failed to delete survey');
    }
  }
};

const getStatusClass = (survey: Survey): string => {
  if (survey.is_published) {
    return 'bg-success';
  } else if (survey.is_draft) {
    return 'bg-secondary';
  } else {
    return 'bg-primary';
  }
};

const getStatusText = (survey: Survey): string => {
  if (survey.is_published) {
    return 'Published';
  } else if (survey.is_draft) {
    return 'Draft';
  } else {
    return 'Active';
  }
};

const copyShareLink = (surveyId: string) => {
  const shareUrl = `${window.location.origin}/surveys/${surveyId}/take`;
  navigator.clipboard.writeText(shareUrl)
    .then(() => {
      alert('Share link copied to clipboard!');
    })
    .catch(() => {
      alert('Failed to copy link. The share URL is: ' + shareUrl);
    });
};

onMounted(() => {
  if (authStore.isAuthenticated) {
    fetchSurveys();
  } else {
    error.value = "Please log in to see your surveys";
    loading.value = false;
  }
});
</script>

<style scoped>
.my-surveys-page {
  background-color: #f8f9fa;
  min-height: calc(100vh - 70px);
}

.survey-card {
  transition: transform 0.2s, box-shadow 0.2s;
  border-radius: 8px;
  overflow: hidden;
  height: 100%;
}

.survey-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1) !important;
}

.empty-state {
  max-width: 600px;
  margin: 0 auto;
  padding: 3rem;
}

.dropdown-item {
  padding: 0.5rem 1rem;
}

.dropdown-item:hover {
  background-color: #f8f9fa;
}

.dropdown-item.text-danger:hover {
  background-color: #f8d7da;
}

.card-title {
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.card-text {
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
}

.survey-meta {
  font-size: 0.85rem;
}

.meta-item {
  display: flex;
  align-items: center;
}

@media (max-width: 768px) {
  .my-surveys-page {
    padding: 1rem 0;
  }
}
</style> 