<template>
  <v-container fluid class="fill-height pa-0">
    <v-row class="fill-height ma-0">
      <!-- Left side - illustration -->
      <v-col cols="12" md="6" class="d-none d-md-flex bg-primary pa-0 login-bg">
        <div class="d-flex flex-column justify-center align-center w-100 h-100 px-12 text-white">
          <div class="text-center">
            <h1 class="text-h3 font-weight-bold mb-6">Welcome Back</h1>
            <p class="text-body-1 mb-8">Sign in to continue your journey with Survey Platform.</p>
            <div class="login-illustration">
              <svg viewBox="0 0 500 400" xmlns="http://www.w3.org/2000/svg" class="illustration">
                <!-- Background shapes -->
                <circle cx="250" cy="200" r="150" fill="rgba(255,255,255,0.1)" />
                <circle cx="250" cy="200" r="100" fill="rgba(255,255,255,0.15)" />
                
                <!-- Document/Survey -->
                <rect x="150" y="100" width="200" height="250" rx="10" fill="white" />
                
                <!-- Document Header -->
                <rect x="170" y="120" width="160" height="20" rx="5" fill="#e6e6ff" />
                
                <!-- Document Content -->
                <rect x="170" y="160" width="160" height="15" rx="3" fill="#e6e6ff" />
                <rect x="170" y="185" width="160" height="15" rx="3" fill="#e6e6ff" />
                
                <!-- Radio Buttons -->
                <circle cx="180" cy="220" r="8" fill="#4361ee" stroke="white" stroke-width="2" />
                <rect x="195" y="215" width="120" height="10" rx="3" fill="#e6e6ff" />
                
                <circle cx="180" cy="250" r="8" fill="#e6e6ff" stroke="#4361ee" stroke-width="2" />
                <rect x="195" y="245" width="120" height="10" rx="3" fill="#e6e6ff" />
                
                <circle cx="180" cy="280" r="8" fill="#e6e6ff" stroke="#4361ee" stroke-width="2" />
                <rect x="195" y="275" width="120" height="10" rx="3" fill="#e6e6ff" />
                
                <!-- Submit Button -->
                <rect x="190" y="310" width="120" height="25" rx="12.5" fill="#4361ee" />
                
                <!-- Person/User -->
                <circle cx="385" cy="180" r="30" fill="#7209b7" />
                <path d="M385 140 v-20" stroke="white" stroke-width="3" />
                <path d="M360 200 a 25 25 0 0 1 50 0" fill="#7209b7" />
                
                <!-- Decorative Elements -->
                <circle cx="100" cy="100" r="20" fill="rgba(255,255,255,0.2)" />
                <circle cx="400" cy="320" r="15" fill="rgba(255,255,255,0.2)" />
                <circle cx="130" cy="350" r="25" fill="rgba(255,255,255,0.1)" />
              </svg>
            </div>
          </div>
        </div>
      </v-col>

      <!-- Right side - login form -->
      <v-col cols="12" md="6" class="d-flex align-center justify-center pa-8 pa-sm-12">
        <v-card class="login-card mx-auto" flat max-width="450">
          <v-card-title class="text-h4 font-weight-bold text-center pb-2">
            Log In
          </v-card-title>
          <v-card-subtitle class="text-center pb-6">
            Enter your credentials to access your account
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
          
          <v-form @submit.prevent="login" ref="form" v-model="valid" class="mt-4">
            <v-text-field
              v-model="username"
              label="Username"
              :rules="usernameRules"
              prepend-inner-icon="mdi-account"
              variant="outlined"
              density="comfortable"
              class="mb-4"
              required
              autofocus
            ></v-text-field>
            
            <v-text-field
              v-model="password"
              label="Password"
              :rules="passwordRules"
              :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
              prepend-inner-icon="mdi-lock"
              @click:append-inner="showPassword = !showPassword"
              :type="showPassword ? 'text' : 'password'"
              variant="outlined"
              density="comfortable"
              class="mb-2"
              required
            ></v-text-field>
            
            <div class="d-flex justify-space-between align-center mb-8">
              <v-checkbox
                v-model="rememberMe"
                label="Remember me"
                hide-details
                density="compact"
              ></v-checkbox>
              <a href="#" class="text-decoration-none text-primary">Forgot password?</a>
            </div>
            
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
              <span class="text-body-1 font-weight-bold">Log In</span>
            </v-btn>
            
            <div class="mt-6 text-center">
              <span class="text-body-2">Don't have an account?</span> 
              <router-link to="/register" class="text-decoration-none font-weight-bold ml-1">
                Sign up now
              </router-link>
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
    const showPassword = ref(false);
    const rememberMe = ref(false);
    
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
          password: password.value,
          rememberMe: rememberMe.value
        });
        
        // Redirect to dashboard or the requested page
        const redirectPath = route.query.redirect || '/dashboard';
        router.push(redirectPath);
      } catch (err) {
        error.value = err.error || 'Authentication failed. Please check your credentials.';
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
      showPassword,
      rememberMe,
      usernameRules,
      passwordRules,
      login
    };
  }
};
</script>

<style scoped>
.login-bg {
  background: linear-gradient(135deg, var(--v-primary-base) 0%, var(--v-secondary-base) 100%);
  position: relative;
}

.login-bg::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-image: url("data:image/svg+xml,%3Csvg width='20' height='20' viewBox='0 0 20 20' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='%23ffffff' fill-opacity='0.05' fill-rule='evenodd'%3E%3Ccircle cx='3' cy='3' r='3'/%3E%3Ccircle cx='13' cy='13' r='3'/%3E%3C/g%3E%3C/svg%3E");
}

.login-card {
  border-radius: 16px;
  width: 100%;
}

.v-btn {
  letter-spacing: 0.5px;
  text-transform: none;
  border-radius: 8px;
}

.theme--light.v-text-field .v-input__slot:hover {
  border-color: var(--v-primary-base) !important;
}

.login-illustration {
  width: 100%;
  max-width: 400px;
  margin: 0 auto;
}

.illustration {
  width: 100%;
  height: auto;
  filter: drop-shadow(0 10px 20px rgba(0,0,0,0.15));
  animation: float 6s ease-in-out infinite;
}

@keyframes float {
  0% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-20px);
  }
  100% {
    transform: translateY(0px);
  }
}
</style> 