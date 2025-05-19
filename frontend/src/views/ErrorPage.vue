<template>
  <div class="error-page">
    <v-container fluid class="fill-height">
      <v-row justify="center" align="center">
        <v-col cols="12" sm="10" md="8" lg="6" class="text-center">
          <div class="animate__animated animate__fadeIn">
            <div class="display-1 font-weight-bold error-code mb-4">{{ code }}</div>
            <h1 class="text-h3 font-weight-bold mb-3">{{ title }}</h1>
            <p class="text-body-1 mb-8">{{ message }}</p>
            
            <v-img
              :src="imageUrl"
              max-width="400"
              class="mx-auto mb-10 animate__animated animate__pulse animate__infinite animate__slow rounded-xl"
            ></v-img>
            
            <div class="d-flex flex-wrap justify-center gap-4">
              <v-btn 
                color="primary" 
                size="large" 
                rounded="pill"
                variant="flat"
                prepend-icon="mdi-home"
                :to="homePath"
                class="px-6"
              >
                {{ homeText }}
              </v-btn>
              <v-btn 
                v-if="showBackButton"
                color="secondary" 
                size="large" 
                rounded="pill"
                variant="outlined"
                prepend-icon="mdi-arrow-left"
                @click="goBack"
                class="px-6"
              >
                Go Back
              </v-btn>
              <v-btn 
                v-if="actionPath"
                :color="actionColor" 
                size="large" 
                rounded="pill"
                variant="flat"
                :prepend-icon="actionIcon"
                :to="actionPath"
                class="px-6"
              >
                {{ actionText }}
              </v-btn>
            </div>
          </div>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import { useRouter } from 'vue-router';
import 'animate.css';

export default {
  name: 'ErrorPage',
  props: {
    code: {
      type: String,
      default: '404'
    },
    title: {
      type: String,
      default: 'Page Not Found'
    },
    message: {
      type: String,
      default: 'The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.'
    },
    imageUrl: {
      type: String,
      default: 'https://placehold.co/600x400/4cc9f0/121212?text=Error'
    },
    homePath: {
      type: String,
      default: '/'
    },
    homeText: {
      type: String,
      default: 'Go Home'
    },
    showBackButton: {
      type: Boolean,
      default: true
    },
    actionPath: {
      type: String,
      default: ''
    },
    actionText: {
      type: String,
      default: 'Try Again'
    },
    actionIcon: {
      type: String,
      default: 'mdi-refresh'
    },
    actionColor: {
      type: String,
      default: 'info'
    }
  },
  setup() {
    const router = useRouter();
    
    const goBack = () => {
      router.go(-1);
    };
    
    return {
      goBack
    };
  }
}
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  background-color: var(--v-background);
}

.error-code {
  font-size: 120px;
  line-height: 1;
  background: linear-gradient(90deg, var(--v-primary-base), var(--v-secondary-base));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 900;
  letter-spacing: -2px;
  margin-bottom: 20px;
}
</style> 