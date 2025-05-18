<template>
  <v-app>
    <v-app-bar color="primary" dark app>
      <v-toolbar-title>
        <router-link to="/" class="text-decoration-none text-white">Survey Platform</router-link>
      </v-toolbar-title>

      <v-spacer></v-spacer>

      <template v-if="!isAuthenticated">
        <v-btn to="/login" variant="text" class="mx-1">Login</v-btn>
        <v-btn to="/register" variant="text" class="mx-1">Register</v-btn>
      </template>
      <template v-else>
        <v-btn to="/dashboard" variant="text" class="mx-1">Dashboard</v-btn>
        <v-menu>
          <template v-slot:activator="{ props }">
            <v-btn icon v-bind="props">
              <v-avatar color="surface-variant" size="32">
                <v-icon>mdi-account</v-icon>
              </v-avatar>
            </v-btn>
          </template>
          <v-list>
            <v-list-item @click="logout">
              <v-list-item-title>Logout</v-list-item-title>
            </v-list-item>
          </v-list>
        </v-menu>
      </template>
    </v-app-bar>

    <v-main>
      <v-container fluid>
        <router-view />
      </v-container>
    </v-main>

    <v-footer app class="bg-primary text-center d-flex justify-center" dark>
      <div>© {{ new Date().getFullYear() }} Survey Platform - All rights reserved</div>
    </v-footer>
  </v-app>
</template>

<script>
import { useAuthStore } from './store/auth';
import { computed } from 'vue';

export default {
  name: 'App',
  setup() {
    const authStore = useAuthStore();
    
    const isAuthenticated = computed(() => authStore.isAuthenticated);
    
    const logout = () => {
      authStore.logout();
    };
    
    return {
      isAuthenticated,
      logout
    };
  }
};
</script>

<style>
:root {
  --primary-color: #3f51b5;
  --secondary-color: #ff4081;
  --background-color: #f5f5f5;
  --text-color: #333333;
}

body {
  font-family: 'Roboto', sans-serif;
  background-color: var(--background-color);
  color: var(--text-color);
}

.page-container {
  padding: 24px;
}

.text-decoration-none {
  text-decoration: none;
}
</style> 