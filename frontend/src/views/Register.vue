<template>
  <v-container>
    <v-row justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="pa-6 elevation-10">
          <v-card-title class="text-h4 text-center pb-4">
            Create an Account
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
          
          <v-form @submit.prevent="register" ref="form" v-model="valid">
            <v-text-field
              v-model="username"
              label="Username"
              :rules="usernameRules"
              required
              prepend-inner-icon="mdi-account"
              variant="outlined"
              density="comfortable"
              class="mb-2"
            ></v-text-field>
            
            <v-text-field
              v-model="email"
              label="Email"
              :rules="emailRules"
              type="email"
              required
              prepend-inner-icon="mdi-email"
              variant="outlined"
              density="comfortable"
              class="mb-2"
            ></v-text-field>
            
            <v-text-field
              v-model="password"
              label="Password"
              :rules="passwordRules"
              type="password"
              required
              prepend-inner-icon="mdi-lock"
              variant="outlined"
              density="comfortable"
              class="mb-2"
            ></v-text-field>
            
            <v-text-field
              v-model="confirmPassword"
              label="Confirm Password"
              :rules="[...passwordRules, passwordMatchRule]"
              type="password"
              required
              prepend-inner-icon="mdi-lock-check"
              variant="outlined"
              density="comfortable"
              class="mb-4"
            ></v-text-field>
            
            <v-btn
              type="submit"
              color="primary"
              block
              :loading="loading"
              :disabled="!valid || loading"
              class="mt-4"
              size="large"
            >
              Register
            </v-btn>
            
            <div class="mt-6 text-center">
              Already have an account? 
              <router-link to="/login" class="text-primary">Login</router-link>
            </div>
          </v-form>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../store/auth';

export default {
  name: 'Register',
  setup() {
    const authStore = useAuthStore();
    const router = useRouter();
    
    const username = ref('');
    const email = ref('');
    const password = ref('');
    const confirmPassword = ref('');
    const error = ref('');
    const loading = ref(false);
    const valid = ref(false);
    const form = ref(null);
    
    const usernameRules = [
      v => !!v || 'Username is required',
      v => v.length >= 3 || 'Username must be at least 3 characters',
    ];
    
    const emailRules = [
      v => !!v || 'Email is required',
      v => /.+@.+\..+/.test(v) || 'Email must be valid',
    ];
    
    const passwordRules = [
      v => !!v || 'Password is required',
      v => v.length >= 6 || 'Password must be at least 6 characters',
    ];
    
    const passwordMatchRule = value => {
      return value === password.value || 'Passwords do not match';
    };
    
    const register = async () => {
      if (!valid.value) return;
      
      loading.value = true;
      error.value = '';
      
      try {
        await authStore.register({
          username: username.value,
          email: email.value,
          password: password.value
        });
        
        // Redirect to login page after successful registration
        router.push('/login');
      } catch (err) {
        error.value = err.error || 'Failed to register. Please try again later.';
      } finally {
        loading.value = false;
      }
    };
    
    return {
      username,
      email,
      password,
      confirmPassword,
      error,
      loading,
      valid,
      form,
      usernameRules,
      emailRules,
      passwordRules,
      passwordMatchRule,
      register
    };
  }
};
</script>

<style scoped>
.v-card {
  border-radius: 12px;
}
</style> 