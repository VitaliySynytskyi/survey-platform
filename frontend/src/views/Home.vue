<template>
  <div class="home">
    <!-- Hero Section -->
    <v-row justify="center" align="center" class="hero-section py-12">
      <v-col cols="12" md="10" lg="8" class="text-center">
        <h1 class="text-h2 font-weight-bold mb-4 animate__animated animate__fadeInDown">
          Create Powerful Surveys with Ease
        </h1>
        <p class="text-body-1 mb-8 animate__animated animate__fadeIn animate__delay-1s">
          Design beautiful surveys, collect valuable insights, and make data-driven decisions 
          with our intuitive platform.
        </p>
        
        <div class="d-flex flex-wrap justify-center gap-3 animate__animated animate__fadeInUp animate__delay-1s">
          <v-btn
            to="/register"
            color="primary"
            size="large"
            rounded="pill"
            class="px-6 py-3"
            variant="flat"
            v-if="!isAuthenticated"
          >
            <v-icon start>mdi-rocket-launch</v-icon>
            Get Started for Free
          </v-btn>
          <v-btn
            to="/login"
            color="secondary"
            size="large"
            rounded="pill"
            class="px-6 py-3 ml-4"
            variant="outlined"
            v-if="!isAuthenticated"
          >
            <v-icon start>mdi-login</v-icon>
            Log In
          </v-btn>
          <v-btn
            to="/dashboard"
            color="primary"
            size="large"
            rounded="pill"
            class="px-6 py-3"
            variant="flat"
            v-if="isAuthenticated"
          >
            <v-icon start>mdi-view-dashboard</v-icon>
            Go to Dashboard
          </v-btn>
        </div>
        
        <div class="mt-16 preview-image animate__animated animate__zoomIn animate__delay-2s">
          <v-img
            src="https://placehold.co/1200x600/4361ee/ffffff?text=Survey+Platform"
            alt="Survey Platform Preview"
            class="rounded-xl image-shadow"
            cover
          ></v-img>
        </div>
      </v-col>
    </v-row>
    
    <!-- Features Section -->
    <v-container class="py-16">
      <v-row class="mb-12 text-center">
        <v-col cols="12">
          <span class="text-overline text-primary font-weight-bold">FEATURES</span>
          <h2 class="text-h3 font-weight-bold mt-2">Everything You Need</h2>
          <p class="text-body-1 mx-auto mt-3" style="max-width: 600px">
            Our platform provides all the tools you need to create, distribute, and analyze surveys.
          </p>
        </v-col>
      </v-row>
      
      <v-row>
        <v-col cols="12" md="4">
          <v-card elevation="0" variant="flat" class="feature-card h-100 pa-5">
            <div class="d-flex justify-center mb-6">
              <v-avatar size="80" color="primary" class="elevation-3">
                <v-icon size="40" color="white">mdi-pencil-ruler</v-icon>
              </v-avatar>
            </div>
            <v-card-title class="text-center text-h5 font-weight-bold">Intuitive Survey Builder</v-card-title>
            <v-card-text class="text-center">
              Create professional surveys in minutes with our drag-and-drop editor. Choose from multiple question types and customize your design.
            </v-card-text>
          </v-card>
        </v-col>
        
        <v-col cols="12" md="4">
          <v-card elevation="0" variant="flat" class="feature-card h-100 pa-5">
            <div class="d-flex justify-center mb-6">
              <v-avatar size="80" color="secondary" class="elevation-3">
                <v-icon size="40" color="white">mdi-share-variant</v-icon>
              </v-avatar>
            </div>
            <v-card-title class="text-center text-h5 font-weight-bold">Seamless Distribution</v-card-title>
            <v-card-text class="text-center">
              Share your surveys instantly via email, social media, or embedded on your website. Track respondents in real-time.
            </v-card-text>
          </v-card>
        </v-col>
        
        <v-col cols="12" md="4">
          <v-card elevation="0" variant="flat" class="feature-card h-100 pa-5">
            <div class="d-flex justify-center mb-6">
              <v-avatar size="80" color="accent" class="elevation-3">
                <v-icon size="40" color="white">mdi-chart-bar</v-icon>
              </v-avatar>
            </div>
            <v-card-title class="text-center text-h5 font-weight-bold">Powerful Analytics</v-card-title>
            <v-card-text class="text-center">
              Visualize survey results with beautiful charts and graphs. Export data in multiple formats for further analysis.
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
    
    <!-- Testimonials Section -->
    <v-container fluid class="testimonials-section py-16">
      <v-row justify="center">
        <v-col cols="12" md="10" lg="8">
          <v-row class="mb-10 text-center">
            <v-col cols="12">
              <span class="text-overline text-primary font-weight-bold">TESTIMONIALS</span>
              <h2 class="text-h3 font-weight-bold mt-2">What Our Users Say</h2>
            </v-col>
          </v-row>
          
          <v-row>
            <v-col cols="12" md="6" class="mb-6">
              <v-card class="testimonial-card pa-6 h-100" elevation="3">
                <v-card-text class="text-body-1 font-italic">
                  "This platform has transformed how we collect customer feedback. The analytics tools make it so easy to gain meaningful insights from our surveys."
                </v-card-text>
                <v-divider class="my-4"></v-divider>
                <div class="d-flex align-center">
                  <v-avatar color="primary" size="48" class="mr-3">
                    <v-icon color="white">mdi-account</v-icon>
                  </v-avatar>
                  <div>
                    <div class="font-weight-bold">Sarah Johnson</div>
                    <div class="text-caption">Marketing Manager, TechCorp</div>
                  </div>
                </div>
              </v-card>
            </v-col>
            
            <v-col cols="12" md="6" class="mb-6">
              <v-card class="testimonial-card pa-6 h-100" elevation="3">
                <v-card-text class="text-body-1 font-italic">
                  "We use Survey Platform for all our market research. The ease of creating surveys and analyzing results has significantly improved our workflow."
                </v-card-text>
                <v-divider class="my-4"></v-divider>
                <div class="d-flex align-center">
                  <v-avatar color="secondary" size="48" class="mr-3">
                    <v-icon color="white">mdi-account</v-icon>
                  </v-avatar>
                  <div>
                    <div class="font-weight-bold">Michael Chen</div>
                    <div class="text-caption">Research Director, GlobalData</div>
                  </div>
                </div>
              </v-card>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-container>
    
    <!-- CTA Section -->
    <v-container fluid class="cta-section py-16 text-center">
      <v-row justify="center">
        <v-col cols="12" md="8" lg="6">
          <h2 class="text-h3 font-weight-bold mb-6">Ready to Get Started?</h2>
          <p class="text-body-1 mb-8">
            Join thousands of organizations who trust our platform for their survey needs.
          </p>
          <v-btn
            to="/register"
            color="primary"
            size="large"
            rounded="pill"
            class="px-8 py-3"
            variant="flat"
            v-if="!isAuthenticated"
          >
            <v-icon start>mdi-rocket-launch</v-icon>
            Create Your First Survey
          </v-btn>
          <v-btn
            to="/dashboard"
            color="primary"
            size="large"
            rounded="pill"
            class="px-8 py-3"
            variant="flat"
            v-if="isAuthenticated"
          >
            <v-icon start>mdi-view-dashboard</v-icon>
            Go to Dashboard
          </v-btn>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script>
import { computed } from 'vue';
import { useAuthStore } from '../store/auth';
import 'animate.css';

export default {
  name: 'Home',
  setup() {
    const authStore = useAuthStore();
    const isAuthenticated = computed(() => authStore.isAuthenticated);
    
    return {
      isAuthenticated
    };
  }
};
</script>

<style scoped>
.hero-section {
  min-height: 80vh;
  display: flex;
  align-items: center;
}

.preview-image {
  max-width: 1000px;
  margin: 0 auto;
}

.image-shadow {
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.feature-card {
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  border-radius: 16px;
}

.feature-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.1);
}

.testimonials-section {
  background-color: rgba(67, 97, 238, 0.05);
}

.testimonial-card {
  border-radius: 16px;
  transition: transform 0.3s ease;
}

.testimonial-card:hover {
  transform: translateY(-5px);
}

.cta-section {
  background: linear-gradient(135deg, var(--v-primary-base) 0%, var(--v-secondary-base) 100%);
  color: white;
  border-radius: 0;
}
</style> 