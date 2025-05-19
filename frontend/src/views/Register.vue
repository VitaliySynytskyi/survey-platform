<template>
  <v-container fluid class="fill-height pa-0">
    <v-row class="fill-height ma-0">
      <!-- Left side - illustration -->
      <v-col cols="12" md="6" class="d-none d-md-flex bg-secondary pa-0 register-bg">
        <div class="d-flex flex-column justify-center align-center w-100 h-100 px-12 text-white">
          <div class="text-center">
            <h1 class="text-h3 font-weight-bold mb-6">Join Our Community</h1>
            <p class="text-body-1 mb-8">Create an account to start building and sharing your own surveys.</p>
            <div class="register-illustration">
              <svg viewBox="0 0 500 400" xmlns="http://www.w3.org/2000/svg" class="illustration">
                <!-- Background shapes -->
                <circle cx="250" cy="200" r="150" fill="rgba(255,255,255,0.1)" />
                <circle cx="250" cy="200" r="100" fill="rgba(255,255,255,0.15)" />
                
                <!-- Person 1 -->
                <circle cx="160" cy="180" r="30" fill="#ffffff" />
                <path d="M160 140 v-15" stroke="white" stroke-width="3" />
                <path d="M135 200 a 25 25 0 0 1 50 0" fill="#ffffff" />
                
                <!-- Person 2 -->
                <circle cx="250" cy="180" r="35" fill="#ffffff" />
                <path d="M250 140 v-15" stroke="white" stroke-width="3" />
                <path d="M220 205 a 30 30 0 0 1 60 0" fill="#ffffff" />
                
                <!-- Person 3 -->
                <circle cx="340" cy="180" r="30" fill="#ffffff" />
                <path d="M340 140 v-15" stroke="white" stroke-width="3" />
                <path d="M315 200 a 25 25 0 0 1 50 0" fill="#ffffff" />
                
                <!-- Survey form - person 1 -->
                <rect x="120" y="230" width="80" height="100" rx="5" fill="white" />
                <rect x="130" y="240" width="60" height="8" rx="2" fill="#e6e6ff" />
                <rect x="130" y="255" width="60" height="8" rx="2" fill="#e6e6ff" />
                <circle cx="135" cy="275" r="5" fill="#7209b7" />
                <rect x="145" y="270" width="40" height="8" rx="2" fill="#e6e6ff" />
                
                <!-- Survey form - person 2 -->
                <rect x="210" y="230" width="80" height="100" rx="5" fill="white" />
                <rect x="220" y="240" width="60" height="8" rx="2" fill="#e6e6ff" />
                <rect x="220" y="255" width="60" height="8" rx="2" fill="#e6e6ff" />
                <circle cx="225" cy="275" r="5" fill="#7209b7" />
                <rect x="235" y="270" width="40" height="8" rx="2" fill="#e6e6ff" />
                
                <!-- Survey form - person 3 -->
                <rect x="300" y="230" width="80" height="100" rx="5" fill="white" />
                <rect x="310" y="240" width="60" height="8" rx="2" fill="#e6e6ff" />
                <rect x="310" y="255" width="60" height="8" rx="2" fill="#e6e6ff" />
                <circle cx="315" cy="275" r="5" fill="#7209b7" />
                <rect x="325" y="270" width="40" height="8" rx="2" fill="#e6e6ff" />
                
                <!-- Connecting lines -->
                <path d="M160 220 L160 230" stroke="white" stroke-width="2" />
                <path d="M250 220 L250 230" stroke="white" stroke-width="2" />
                <path d="M340 220 L340 230" stroke="white" stroke-width="2" />
                
                <!-- Decorative Elements -->
                <circle cx="100" cy="100" r="20" fill="rgba(255,255,255,0.2)" />
                <circle cx="400" cy="320" r="15" fill="rgba(255,255,255,0.2)" />
                <circle cx="130" cy="350" r="25" fill="rgba(255,255,255,0.1)" />
              </svg>
            </div>
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

.register-illustration {
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