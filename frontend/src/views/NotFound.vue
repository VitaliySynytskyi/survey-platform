<template>
  <div class="not-found">
    <v-container fluid class="fill-height">
      <v-row justify="center" align="center">
        <v-col cols="12" sm="10" md="8" lg="6" class="text-center">
          <div class="animate__animated animate__fadeIn">
            <div class="display-1 font-weight-bold error-code mb-4">404</div>
            <h1 class="text-h3 font-weight-bold mb-3">Page Not Found</h1>
            <p class="text-body-1 mb-8">The page you are looking for might have been removed, had its name changed, or is temporarily unavailable.</p>
            
            <v-img
              src="https://placehold.co/600x400/4cc9f0/121212?text=Lost+in+Space"
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
                to="/"
                class="px-6"
              >
                Go Home
              </v-btn>
              <v-btn 
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
  name: 'NotFound',
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
.not-found {
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