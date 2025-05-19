<template>
  <div class="profile">
    <v-row>
      <v-col cols="12" md="4">
        <!-- Profile Card -->
        <v-card class="mb-6 pa-4 rounded-xl" elevation="2">
          <div class="d-flex flex-column align-center mb-6">
            <v-avatar color="primary" size="120" class="mb-4">
              <v-icon size="64" color="white">mdi-account</v-icon>
            </v-avatar>
            <h2 class="text-h4 font-weight-bold">{{ user?.name || user?.username }}</h2>
            <p class="text-body-1 text-medium-emphasis">{{ user?.email }}</p>
            <v-chip color="primary" size="small" class="mt-2">{{ user?.role === 'admin' ? 'Administrator' : 'User' }}</v-chip>
          </div>
          
          <v-list class="transparent">
            <v-list-item>
              <template v-slot:prepend>
                <v-icon color="primary">mdi-account-outline</v-icon>
              </template>
              <v-list-item-title>Username</v-list-item-title>
              <v-list-item-subtitle>{{ user?.username }}</v-list-item-subtitle>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-icon color="primary">mdi-email-outline</v-icon>
              </template>
              <v-list-item-title>Email</v-list-item-title>
              <v-list-item-subtitle>{{ user?.email }}</v-list-item-subtitle>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-icon color="primary">mdi-badge-account-outline</v-icon>
              </template>
              <v-list-item-title>Role</v-list-item-title>
              <v-list-item-subtitle>{{ user?.role === 'admin' ? 'Administrator' : 'User' }}</v-list-item-subtitle>
            </v-list-item>
            
            <v-list-item>
              <template v-slot:prepend>
                <v-icon color="primary">mdi-clock-outline</v-icon>
              </template>
              <v-list-item-title>Member Since</v-list-item-title>
              <v-list-item-subtitle>{{ formatDate(user?.created_at) }}</v-list-item-subtitle>
            </v-list-item>
          </v-list>
          
          <v-card-actions class="mt-2">
            <v-btn 
              variant="flat" 
              color="primary" 
              block 
              @click="editProfileDialog = true"
              prepend-icon="mdi-account-edit"
            >
              Edit Profile
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
      
      <v-col cols="12" md="8">
        <!-- Stats Section -->
        <v-card class="mb-6 pa-6 rounded-xl" elevation="2">
          <h3 class="text-h5 font-weight-bold mb-4">Your Activity</h3>
          
          <v-row>
            <v-col cols="12" sm="6" lg="3">
              <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="primary">
                <div class="d-flex align-center mb-2">
                  <v-avatar color="primary" size="36" class="mr-3">
                    <v-icon color="white">mdi-poll</v-icon>
                  </v-avatar>
                  <span class="text-body-2 font-weight-medium">My Surveys</span>
                </div>
                <div class="text-h3 font-weight-bold mt-2">{{ stats.surveysCount }}</div>
              </v-card>
            </v-col>
            
            <v-col cols="12" sm="6" lg="3">
              <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="success">
                <div class="d-flex align-center mb-2">
                  <v-avatar color="success" size="36" class="mr-3">
                    <v-icon color="white">mdi-check-circle</v-icon>
                  </v-avatar>
                  <span class="text-body-2 font-weight-medium">Active Surveys</span>
                </div>
                <div class="text-h3 font-weight-bold mt-2">{{ stats.activeSurveysCount }}</div>
              </v-card>
            </v-col>
            
            <v-col cols="12" sm="6" lg="3">
              <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="info">
                <div class="d-flex align-center mb-2">
                  <v-avatar color="info" size="36" class="mr-3">
                    <v-icon color="white">mdi-account-multiple</v-icon>
                  </v-avatar>
                  <span class="text-body-2 font-weight-medium">Total Responses</span>
                </div>
                <div class="text-h3 font-weight-bold mt-2">{{ stats.responsesCount }}</div>
              </v-card>
            </v-col>
            
            <v-col cols="12" sm="6" lg="3">
              <v-card class="stat-card pa-4 h-100" elevation="0" variant="tonal" color="warning">
                <div class="d-flex align-center mb-2">
                  <v-avatar color="warning" size="36" class="mr-3">
                    <v-icon color="white">mdi-chart-areaspline</v-icon>
                  </v-avatar>
                  <span class="text-body-2 font-weight-medium">Avg. Response Rate</span>
                </div>
                <div class="text-h3 font-weight-bold mt-2">{{ stats.responseRate }}%</div>
              </v-card>
            </v-col>
          </v-row>
        </v-card>
        
        <!-- Recent Activity -->
        <v-card class="mb-6 pa-6 rounded-xl" elevation="2">
          <div class="d-flex align-center justify-space-between mb-4">
            <h3 class="text-h5 font-weight-bold">Recent Activity</h3>
            <v-btn 
              variant="text" 
              color="primary" 
              to="/dashboard"
              prepend-icon="mdi-view-dashboard"
            >
              View All
            </v-btn>
          </div>
          
          <v-timeline density="compact" align="start">
            <v-timeline-item
              v-for="(activity, index) in recentActivity"
              :key="index"
              :dot-color="activity.color"
              size="small"
            >
              <div class="d-flex align-center mb-1">
                <span class="text-subtitle-2 font-weight-medium">{{ activity.title }}</span>
                <v-spacer></v-spacer>
                <span class="text-caption text-medium-emphasis">{{ formatTimeAgo(activity.timestamp) }}</span>
              </div>
              <p class="text-body-2 mb-0">{{ activity.description }}</p>
            </v-timeline-item>
            
            <v-timeline-item v-if="recentActivity.length === 0">
              <p class="text-body-2 font-italic">No recent activity found.</p>
            </v-timeline-item>
          </v-timeline>
        </v-card>
        
        <!-- Account Security -->
        <v-card class="rounded-xl" elevation="2">
          <v-card-title class="text-h5 font-weight-bold px-6 pt-6">Account Security</v-card-title>
          
          <v-card-text class="pa-6">
            <v-list>
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon color="primary" size="24">mdi-lock-outline</v-icon>
                </template>
                <v-list-item-title>Password</v-list-item-title>
                <v-list-item-subtitle>Last changed: {{ passwordLastChanged }}</v-list-item-subtitle>
                <template v-slot:append>
                  <v-btn 
                    size="small" 
                    variant="text" 
                    color="primary"
                    @click="changePasswordDialog = true"
                  >
                    Change
                  </v-btn>
                </template>
              </v-list-item>
              
              <v-divider></v-divider>
              
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon color="primary" size="24">mdi-bell-outline</v-icon>
                </template>
                <v-list-item-title>Notifications</v-list-item-title>
                <v-list-item-subtitle>Manage your notification preferences</v-list-item-subtitle>
                <template v-slot:append>
                  <v-switch
                    v-model="notificationsEnabled"
                    color="primary"
                    hide-details
                    density="compact"
                  ></v-switch>
                </template>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
    
    <!-- Edit Profile Dialog -->
    <v-dialog v-model="editProfileDialog" max-width="600px" persistent>
      <v-card class="pa-4 rounded-xl">
        <v-card-title class="text-h5 font-weight-bold pb-2">Edit Profile</v-card-title>
        
        <v-card-text>
          <v-form ref="editProfileForm" v-model="editProfileFormValid">
            <v-text-field
              v-model="editedProfile.name"
              label="Display Name"
              prepend-inner-icon="mdi-account"
              variant="outlined"
              class="mb-3"
            ></v-text-field>
            
            <v-text-field
              v-model="editedProfile.username"
              label="Username"
              prepend-inner-icon="mdi-account-outline"
              variant="outlined"
              class="mb-3"
              :rules="[v => !!v || 'Username is required']"
            ></v-text-field>
            
            <v-text-field
              v-model="editedProfile.email"
              label="Email"
              prepend-inner-icon="mdi-email-outline"
              variant="outlined"
              class="mb-3"
              :rules="[
                v => !!v || 'Email is required',
                v => /.+@.+\..+/.test(v) || 'Email must be valid'
              ]"
            ></v-text-field>
          </v-form>
        </v-card-text>
        
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn 
            variant="text" 
            color="grey"
            @click="editProfileDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn 
            variant="flat" 
            color="primary"
            :disabled="!editProfileFormValid || updateLoading"
            :loading="updateLoading"
            @click="updateProfile"
          >
            Save Changes
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
    
    <!-- Change Password Dialog -->
    <v-dialog v-model="changePasswordDialog" max-width="600px" persistent>
      <v-card class="pa-4 rounded-xl">
        <v-card-title class="text-h5 font-weight-bold pb-2">Change Password</v-card-title>
        
        <v-card-text>
          <v-form ref="passwordForm" v-model="passwordFormValid">
            <v-text-field
              v-model="passwordData.currentPassword"
              label="Current Password"
              prepend-inner-icon="mdi-lock-outline"
              variant="outlined"
              type="password"
              class="mb-3"
              :rules="[v => !!v || 'Current password is required']"
            ></v-text-field>
            
            <v-text-field
              v-model="passwordData.newPassword"
              label="New Password"
              prepend-inner-icon="mdi-lock-outline"
              variant="outlined"
              type="password"
              class="mb-3"
              :rules="[
                v => !!v || 'New password is required',
                v => v && v.length >= 8 || 'Password must be at least 8 characters',
                v => v !== passwordData.currentPassword || 'New password must be different'
              ]"
            ></v-text-field>
            
            <v-text-field
              v-model="passwordData.confirmPassword"
              label="Confirm New Password"
              prepend-inner-icon="mdi-lock-outline"
              variant="outlined"
              type="password"
              class="mb-3"
              :rules="[
                v => !!v || 'Please confirm your password',
                v => v === passwordData.newPassword || 'Passwords do not match'
              ]"
            ></v-text-field>
          </v-form>
        </v-card-text>
        
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn 
            variant="text" 
            color="grey"
            @click="changePasswordDialog = false"
          >
            Cancel
          </v-btn>
          <v-btn 
            variant="flat" 
            color="primary"
            :disabled="!passwordFormValid || passwordLoading"
            :loading="passwordLoading"
            @click="updatePassword"
          >
            Update Password
          </v-btn>
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
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '../store/auth';
import { authApi, surveyApi } from '../services/api';

export default {
  name: 'ProfileView',
  
  setup() {
    const authStore = useAuthStore();
    const user = ref(null);
    const stats = ref({
      surveysCount: 0,
      activeSurveysCount: 0,
      responsesCount: 0,
      responseRate: 0
    });
    
    // Recent activity
    const recentActivity = ref([]);
    
    // Edit profile dialog
    const editProfileDialog = ref(false);
    const editProfileForm = ref(null);
    const editProfileFormValid = ref(false);
    const editedProfile = ref({
      name: '',
      username: '',
      email: ''
    });
    const updateLoading = ref(false);
    
    // Password change dialog
    const changePasswordDialog = ref(false);
    const passwordForm = ref(null);
    const passwordFormValid = ref(false);
    const passwordData = ref({
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    });
    const passwordLoading = ref(false);
    const passwordLastChanged = ref('Never');
    
    // Notification settings
    const notificationsEnabled = ref(true);
    
    // Snackbar
    const snackbar = ref({
      show: false,
      text: '',
      color: 'success'
    });
    
    // Fetch user data
    const fetchUserData = async () => {
      try {
        // Get user profile
        const userData = await authApi.getCurrentUser();
        user.value = userData.data;
        
        // Populate edit form with current data
        editedProfile.value = {
          name: user.value.name || '',
          username: user.value.username || '',
          email: user.value.email || ''
        };
        
        // Set password last changed date
        passwordLastChanged.value = user.value.password_updated_at 
          ? formatDate(user.value.password_updated_at) 
          : 'Never';
        
        // Fetch user statistics
        await fetchUserStats();
        
        // Fetch recent activity
        await fetchRecentActivity();
      } catch (error) {
        console.error('Error fetching user data:', error);
        showSnackbar('Failed to load profile data', 'error');
      }
    };
    
    // Fetch user statistics
    const fetchUserStats = async () => {
      stats.value = { surveysCount: 0, activeSurveysCount: 0, responsesCount: 0, responseRate: 0 }; // Reset stats
      try {
        const surveysResponse = await surveyApi.getUserSurveys({ limit: 1000 }); // Fetch all surveys
        
        if (surveysResponse.data && surveysResponse.data.data) {
          const userSurveys = surveysResponse.data.data;
          stats.value.surveysCount = userSurveys.length;
          
          const activeUserSurveys = userSurveys.filter(s => s.is_active);
          stats.value.activeSurveysCount = activeUserSurveys.length;
          
          let totalResponsesAllSurveys = 0;
          let activeSurveysWithResponses = 0;

          if (userSurveys.length > 0) {
            const surveyResponseDataPromises = userSurveys.map(async (survey) => {
              try {
                const response = await surveyApi.getSurveyResponses(survey.id, { count_only: true });
                let count = 0;
                if (response.data) {
                  if (typeof response.data.count === 'number') {
                    count = response.data.count;
                  } else if (Array.isArray(response.data)) {
                    count = response.data.length;
                  }
                }
                return { surveyId: survey.id, is_active: survey.is_active, count: count };
              } catch (err) {
                console.error(`Failed to fetch response count for survey ${survey.id} in profile:`, err);
                return { surveyId: survey.id, is_active: survey.is_active, count: 0 };
              }
            });
            
            const surveyResponseData = await Promise.all(surveyResponseDataPromises);

            surveyResponseData.forEach(data => {
              totalResponsesAllSurveys += data.count;
              if (data.is_active && data.count > 0) {
                activeSurveysWithResponses++;
              }
            });
          }
          
          stats.value.responsesCount = totalResponsesAllSurveys; // Total responses from all surveys
          
          // Calculate Avg. Response Rate as: 
          // (Number of active surveys with at least one response / Total number of active surveys) * 100
          if (stats.value.activeSurveysCount > 0) {
            stats.value.responseRate = Math.round((activeSurveysWithResponses / stats.value.activeSurveysCount) * 100);
          } else {
            stats.value.responseRate = 0; // Avoid division by zero if no active surveys
          }
        }
      } catch (error) {
        console.error('Error fetching user stats for profile:', error);
      }
    };
    
    // Fetch recent activity
    const fetchRecentActivity = async () => {
      try {
        // Get user surveys for activity timeline
        const params = { page: 1, limit: 5 }; // Just get the 5 most recent
        const surveysResponse = await surveyApi.getUserSurveys(params);
        
        if (surveysResponse.data && surveysResponse.data.data) {
          const surveys = surveysResponse.data.data;
          
          // Create activity items for recent surveys
          recentActivity.value = surveys.map(survey => ({
            title: `Survey "${survey.title}"`,
            description: survey.is_active 
              ? `Created and published` 
              : `Created but not yet published`,
            timestamp: survey.created_at,
            color: survey.is_active ? 'success' : 'grey'
          }));
        }
      } catch (error) {
        console.error('Error fetching recent activity:', error);
      }
    };
    
    // Update user profile
    const updateProfile = async () => {
      if (!editProfileFormValid.value) return;
      
      updateLoading.value = true;
      try {
        await authApi.updateCurrentUser(editedProfile.value);
        
        // Update local user data
        user.value = {
          ...user.value,
          ...editedProfile.value
        };
        
        // Update auth store
        await authStore.fetchUser();
        
        showSnackbar('Profile updated successfully', 'success');
        editProfileDialog.value = false;
      } catch (error) {
        console.error('Error updating profile:', error);
        showSnackbar('Failed to update profile', 'error');
      } finally {
        updateLoading.value = false;
      }
    };
    
    // Update password
    const updatePassword = async () => {
      if (!passwordFormValid.value) return;
      
      passwordLoading.value = true;
      try {
        // Make API call to update password
        await authApi.updateCurrentUser({
          current_password: passwordData.value.currentPassword,
          new_password: passwordData.value.newPassword
        });
        
        // Reset form
        passwordData.value = {
          currentPassword: '',
          newPassword: '',
          confirmPassword: ''
        };
        
        showSnackbar('Password updated successfully', 'success');
        changePasswordDialog.value = false;
        
        // Update last changed date
        passwordLastChanged.value = formatDate(new Date());
      } catch (error) {
        console.error('Error updating password:', error);
        showSnackbar(
          error.response?.data?.error || 'Failed to update password', 
          'error'
        );
      } finally {
        passwordLoading.value = false;
      }
    };
    
    // Format date
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A';
      const options = { year: 'numeric', month: 'short', day: 'numeric' };
      try {
        return new Date(dateString).toLocaleDateString(undefined, options);
      } catch (e) {
        return dateString;
      }
    };
    
    // Format time ago
    const formatTimeAgo = (dateString) => {
      if (!dateString) return 'N/A';
      
      try {
        const date = new Date(dateString);
        const now = new Date();
        const seconds = Math.floor((now - date) / 1000);
        
        let interval = Math.floor(seconds / 31536000);
        if (interval >= 1) {
          return interval === 1 ? '1 year ago' : `${interval} years ago`;
        }
        
        interval = Math.floor(seconds / 2592000);
        if (interval >= 1) {
          return interval === 1 ? '1 month ago' : `${interval} months ago`;
        }
        
        interval = Math.floor(seconds / 86400);
        if (interval >= 1) {
          return interval === 1 ? '1 day ago' : `${interval} days ago`;
        }
        
        interval = Math.floor(seconds / 3600);
        if (interval >= 1) {
          return interval === 1 ? '1 hour ago' : `${interval} hours ago`;
        }
        
        interval = Math.floor(seconds / 60);
        if (interval >= 1) {
          return interval === 1 ? '1 minute ago' : `${interval} minutes ago`;
        }
        
        return seconds < 10 ? 'just now' : `${Math.floor(seconds)} seconds ago`;
      } catch (e) {
        return dateString;
      }
    };
    
    // Show snackbar
    const showSnackbar = (text, color = 'success') => {
      snackbar.value = {
        show: true,
        text,
        color
      };
    };
    
    onMounted(() => {
      fetchUserData();
    });
    
    return {
      user,
      stats,
      recentActivity,
      editProfileDialog,
      editProfileForm,
      editProfileFormValid,
      editedProfile,
      updateLoading,
      changePasswordDialog,
      passwordForm,
      passwordFormValid,
      passwordData,
      passwordLoading,
      passwordLastChanged,
      notificationsEnabled,
      snackbar,
      formatDate,
      formatTimeAgo,
      updateProfile,
      updatePassword,
      showSnackbar
    };
  }
};
</script>

<style scoped>
.profile {
  animation: fadeIn 0.5s ease-out;
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