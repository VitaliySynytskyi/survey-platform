<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="pa-6">
          <v-card-title class="text-h4 text-center pb-6">
            Login
          </v-card-title>
          
          <v-alert
            v-if="error"
            type="error"
            class="mb-4"
            closable
            @click:close="error = ''"
          >
            {{ error }}
          </v-alert>
          
          <v-form @submit.prevent="login" ref="form" v-model="valid">
            <v-text-field
              v-model="username"
              label="Username"
              :rules="usernameRules"
              required
              autofocus
            ></v-text-field>
            
            <v-text-field
              v-model="password"
              label="Password"
              :rules="passwordRules"
              type="password"
              required
            ></v-text-field>
            
            <v-btn
              type="submit"
              color="primary"
              block
              :loading="loading"
              :disabled="!valid || loading"
              class="mt-4"
            >
              Login
            </v-btn>
            
            <div class="mt-4 text-center">
              Don't have an account? 
              <router-link to="/register">Register</router-link>
            </div>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useAuthStore } from '../store/auth';

export default {
  name: 'Login',
  setup() {
    const authStore = useAuthStore();
    const router = useRouter();
    const route = useRoute();
    
    const username = ref('');
    const password = ref('');
    const error = ref('');
    const loading = ref(false);
    const valid = ref(false);
    const form = ref(null);
    
    const usernameRules = [
      v => !!v || 'Username is required',
    ];
    
    const passwordRules = [
      v => !!v || 'Password is required',
    ];
    
    const login = async () => {
      if (!valid.value) return;
      
      loading.value = true;
      error.value = '';
      
      try {
        await authStore.login({
          username: username.value,
          password: password.value
        });
        
        // Redirect to dashboard or the requested page
        const redirectPath = route.query.redirect || '/dashboard';
        router.push(redirectPath);
      } catch (err) {
        error.value = err.error || 'Failed to login. Please check your credentials.';
      } finally {
        loading.value = false;
      }
    };
    
    return {
      username,
      password,
      error,
      loading,
      valid,
      form,
      usernameRules,
      passwordRules,
      login
    };
  }
};
</script> 