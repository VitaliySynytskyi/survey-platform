<template>
  <v-app :theme="theme">
    <!-- Modern App Bar -->
    <v-app-bar flat density="comfortable" elevation="1">
      <v-app-bar-nav-icon @click="drawer = !drawer" class="d-md-none"></v-app-bar-nav-icon>
      
      <v-app-bar-title class="font-weight-bold">
        <router-link to="/" class="text-decoration-none">
          <span class="text-primary">Survey</span><span class="text-secondary">Platform</span>
        </router-link>
      </v-app-bar-title>

      <v-spacer></v-spacer>

      <!-- Desktop Navigation -->
      <div class="d-none d-md-flex">
        <template v-if="!isAuthenticated">
          <v-btn to="/login" variant="text" class="mx-1">
            <v-icon start>mdi-login</v-icon> Login
          </v-btn>
          <v-btn to="/register" color="primary" variant="tonal" class="mx-1">
            <v-icon start>mdi-account-plus</v-icon> Register
          </v-btn>
        </template>
        <template v-else>
          <v-btn to="/dashboard" variant="text" class="mx-1">
            <v-icon start>mdi-view-dashboard</v-icon> Dashboard
          </v-btn>
          <v-btn to="/surveys/create" color="primary" variant="tonal" class="mx-1">
            <v-icon start>mdi-plus</v-icon> New Survey
          </v-btn>
        </template>
      </div>

      <!-- Theme Toggle -->
      <v-btn icon @click="toggleTheme" class="ml-2">
        <v-icon>{{ theme === 'dark' ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
      </v-btn>

      <!-- User Menu -->
      <template v-if="isAuthenticated">
        <v-menu transition="slide-y-transition">
          <template v-slot:activator="{ props }">
            <v-btn icon class="ml-2" v-bind="props">
              <v-avatar color="primary" size="36">
                <v-icon color="white">mdi-account</v-icon>
              </v-avatar>
            </v-btn>
          </template>
          <v-list density="compact" width="200" elevation="3" rounded="lg">
            <v-list-item to="/profile" prepend-icon="mdi-account-circle" title="Profile"></v-list-item>
            <v-list-item prepend-icon="mdi-cog" title="Settings"></v-list-item>
            <v-divider></v-divider>
            <v-list-item @click="logout" prepend-icon="mdi-logout" title="Logout" color="error"></v-list-item>
          </v-list>
        </v-menu>
      </template>
    </v-app-bar>

    <!-- Navigation Drawer (mobile) -->
    <v-navigation-drawer v-model="drawer" temporary>
      <v-list>
        <v-list-item
          title="Home"
          to="/"
          prepend-icon="mdi-home"
        ></v-list-item>
        
        <v-divider class="my-2"></v-divider>
        
        <template v-if="!isAuthenticated">
          <v-list-item
            title="Login"
            to="/login"
            prepend-icon="mdi-login"
          ></v-list-item>
          <v-list-item
            title="Register"
            to="/register"
            prepend-icon="mdi-account-plus"
          ></v-list-item>
        </template>
        
        <template v-else>
          <v-list-item
            title="Dashboard"
            to="/dashboard"
            prepend-icon="mdi-view-dashboard"
          ></v-list-item>
          <v-list-item
            title="Create Survey"
            to="/surveys/create"
            prepend-icon="mdi-plus-circle"
          ></v-list-item>
          <v-list-item
            title="Profile"
            to="/profile"
            prepend-icon="mdi-account-circle"
          ></v-list-item>
          <v-divider class="my-2"></v-divider>
          <v-list-item
            title="Logout"
            @click="logout"
            prepend-icon="mdi-logout"
            color="error"
          ></v-list-item>
        </template>
      </v-list>
    </v-navigation-drawer>

    <!-- Main Content -->
    <v-main>
      <v-container fluid class="pa-4 pa-sm-6">
        <transition name="fade" mode="out-in">
          <router-view />
        </transition>
      </v-container>
    </v-main>

    <!-- Call to Action Banner -->
    <v-sheet v-if="!isAuthenticated" class="py-8 px-4 primary text-center text-white" elevation="2">
      <v-container>
        <h2 class="text-h4 mb-3 text-white">Ready to Get Started?</h2>
        <p class="text-subtitle-1 mb-6 text-white">Join thousands of organizations who trust our platform for their survey needs.</p>
        <v-btn 
          size="large" 
          variant="flat" 
          color="white" 
          class="text-primary font-weight-bold px-6"
          to="/register"
          rounded="pill"
        >
          Create Free Account
        </v-btn>
      </v-container>
    </v-sheet>

    <!-- Simplified Footer -->
    <v-footer app height="auto" class="py-4" absolute>
      <v-container>
        <div class="d-flex flex-column flex-md-row align-center justify-space-between w-100">
          <div class="mb-3 mb-md-0 text-center text-md-left">
            <span class="text-primary font-weight-bold">Survey</span><span class="text-secondary font-weight-bold">Platform</span>
            <span class="text-body-2 ml-2">© {{ new Date().getFullYear() }}</span>
          </div>
          <div class="d-flex">
            <v-btn icon variant="text" href="https://github.com" target="_blank" aria-label="GitHub" size="small">
              <v-icon>mdi-github</v-icon>
            </v-btn>
            <v-btn icon variant="text" href="https://twitter.com" target="_blank" aria-label="Twitter" size="small">
              <v-icon>mdi-twitter</v-icon>
            </v-btn>
            <v-btn icon variant="text" href="https://linkedin.com" target="_blank" aria-label="LinkedIn" size="small">
              <v-icon>mdi-linkedin</v-icon>
            </v-btn>
          </div>
        </div>
      </v-container>
    </v-footer>
  </v-app>
</template>

<script>
import { useAuthStore } from './store/auth';
import { computed, ref } from 'vue';

export default {
  name: 'App',
  setup() {
    const authStore = useAuthStore();
    const drawer = ref(false);
    const theme = ref(localStorage.getItem('theme') || 'light');
    
    const isAuthenticated = computed(() => authStore.isAuthenticated);
    
    const logout = () => {
      drawer.value = false;
      authStore.logout();
    };
    
    const toggleTheme = () => {
      theme.value = theme.value === 'light' ? 'dark' : 'light';
      localStorage.setItem('theme', theme.value);
    };
    
    return {
      isAuthenticated,
      logout,
      drawer,
      theme,
      toggleTheme
    };
  }
};
</script>

<style>
:root {
  --transition-speed: 0.3s;
}

/* Global styles */
.v-application {
  font-family: 'Inter', 'Roboto', sans-serif;
}

/* Page transition */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--transition-speed) ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Button transitions */
.v-btn {
  transition: all var(--transition-speed) ease;
}

/* Card hover effect */
.v-card {
  transition: transform var(--transition-speed), box-shadow var(--transition-speed);
}

.v-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 14px 28px rgba(0,0,0,0.15), 0 10px 10px rgba(0,0,0,0.12) !important;
}

/* Custom container for sections */
.section-container {
  padding: 24px;
  border-radius: 16px;
  margin-bottom: 24px;
}

/* Footer links */
.footer-link {
  color: inherit;
  text-decoration: none;
  opacity: 0.8;
  transition: opacity 0.2s ease;
}

.footer-link:hover {
  opacity: 1;
  color: var(--v-primary-base);
}

/* Custom scrollbar */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
}

::-webkit-scrollbar-thumb {
  background: var(--v-primary-base);
  border-radius: 4px;
}

::-webkit-scrollbar-thumb:hover {
  background: var(--v-secondary-base);
}
</style> 