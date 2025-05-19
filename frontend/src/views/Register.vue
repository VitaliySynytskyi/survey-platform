<template>
  <v-container fluid class="fill-height pa-0">
    <v-row class="fill-height ma-0">
      <!-- Left side - illustration -->
      <v-col cols="12" md="6" class="d-none d-md-flex bg-secondary pa-0 register-bg">
        <div class="d-flex flex-column justify-center align-center w-100 h-100 px-12 text-white">
          <div class="text-center">
            <h1 class="text-h3 font-weight-bold mb-6">Join Our Community</h1>
            <p class="text-body-1 mb-8">Create an account to start building and sharing your own surveys.</p>
            <v-img
              src="https://placehold.co/500x400/ffffff/7209b7?text=Create+Account"
              alt="Registration Illustration"
              class="mx-auto mt-8"
              max-width="400"
            ></v-img>
          </div>
        </div>
      </v-col>

      <!-- Right side - registration form -->
      <v-col cols="12" md="6" class="d-flex align-center justify-center pa-8 pa-sm-12">
        <v-card class="register-card mx-auto" flat max-width="450">
          <v-card-title class="text-h4 font-weight-bold text-center pb-2">
            Create Account
          </v-card-title>
          <v-card-subtitle class="text-center pb-6">
            Join today and start creating surveys
          </v-card-subtitle>
          
          <v-alert
            v-if="error"
            type="error"
            class="mb-6"
            closable
            variant="tonal"
            @click:close="error = ''"
          >
            {{ error }}
          </v-alert>
          
          <v-form @submit.prevent="register" ref="form" v-model="valid" class="mt-4">
            <v-text-field
              v-model="username"
              label="Username"
              :rules="usernameRules"
              required
              prepend-inner-icon="mdi-account"
              variant="outlined"
              density="comfortable"
              class="mb-4"
              hint="Choose a unique username (min. 3 characters)"
              persistent-hint
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
              class="mb-4"
            ></v-text-field>
            
            <v-text-field
              v-model="password"
              label="Password"
              :rules="passwordRules"
              :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
              prepend-inner-icon="mdi-lock"
              @click:append-inner="showPassword = !showPassword"
              :type="showPassword ? 'text' : 'password'"
              required
              variant="outlined"
              density="comfortable"
              class="mb-4"
              hint="Password must be at least 6 characters long"
              persistent-hint
            ></v-text-field>
            
            <v-text-field
              v-model="confirmPassword"
              label="Confirm Password"
              :rules="[...passwordRules, passwordMatchRule]"
              :append-inner-icon="showConfirmPassword ? 'mdi-eye-off' : 'mdi-eye'"
              prepend-inner-icon="mdi-lock-check"
              @click:append-inner="showConfirmPassword = !showConfirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              required
              variant="outlined"
              density="comfortable"
              class="mb-4"
            ></v-text-field>
            
            <v-checkbox
              v-model="agreeTerms"
              :rules="[v => !!v || 'You must agree to continue']"
              label="I agree to the Terms of Service and Privacy Policy"
              hide-details
              class="mb-6"
            ></v-checkbox>
            
            <v-btn
              type="submit"
              color="primary"
              block
              :loading="loading"
              :disabled="!valid || loading"
              class="py-6 mb-4"
              height="48"
              variant="flat"
            >
              <span class="text-body-1 font-weight-bold">Create Account</span>
            </v-btn>
            
            <div class="mt-6 text-center">
              <span class="text-body-2">Already have an account?</span>
              <router-link to="/login" class="text-decoration-none font-weight-bold ml-1">
                Log in
              </router-link>
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
    const showPassword = ref(false);
    const showConfirmPassword = ref(false);
    const agreeTerms = ref(false);
    
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
        router.push({
          path: '/login',
          query: { registered: 'success' }
        });
      } catch (err) {
        error.value = err.error || 'Registration failed. Please try again or use a different username/email.';
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
      showPassword,
      showConfirmPassword,
      agreeTerms,
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
.register-bg {
  background: linear-gradient(135deg, var(--v-secondary-base) 0%, var(--v-accent-base) 100%);
  position: relative;
}

.register-bg::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url("data:image/svg+xml,%3Csvg width='20' height='20' viewBox='0 0 20 20' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='%23ffffff' fill-opacity='0.05' fill-rule='evenodd'%3E%3Ccircle cx='3' cy='3' r='3'/%3E%3Ccircle cx='13' cy='13' r='3'/%3E%3C/g%3E%3C/svg%3E");
}

.register-card {
  border-radius: 16px;
  width: 100%;
}

.v-btn {
  letter-spacing: 0.5px;
  text-transform: none;
  border-radius: 8px;
}

.v-text-field--outlined .v-input__slot:hover {
  border-color: var(--v-primary-base) !important;
}
</style> 