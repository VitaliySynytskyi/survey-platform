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
          <v-card class="rounded-xl image-shadow pa-0" elevation="10">
            <div class="survey-platform-preview">
              <div class="preview-overlay"></div>
              <div class="preview-content">
                <!-- Dashboard Preview with Real Visual Elements -->
                <svg viewBox="0 0 1000 550" xmlns="http://www.w3.org/2000/svg" class="platform-illustration">
                  <!-- Background -->
                  <linearGradient id="bg-gradient" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" stop-color="#4361ee" />
                    <stop offset="100%" stop-color="#7209b7" />
                  </linearGradient>
                  <rect x="0" y="0" width="1000" height="550" fill="url(#bg-gradient)" />
                  
                  <!-- Decorative patterns -->
                  <circle cx="150" cy="100" r="50" fill="rgba(255,255,255,0.1)" />
                  <circle cx="850" cy="450" r="70" fill="rgba(255,255,255,0.1)" />
                  <circle cx="950" cy="150" r="30" fill="rgba(255,255,255,0.1)" />
                  <circle cx="200" cy="500" r="40" fill="rgba(255,255,255,0.1)" />
                  
                  <!-- Main dashboard panel -->
                  <rect x="80" y="70" width="550" height="400" rx="12" fill="white" filter="drop-shadow(0 10px 15px rgba(0,0,0,0.2))" />
                  
                  <!-- Dashboard header -->
                  <rect x="100" y="90" width="150" height="30" rx="6" fill="#4361ee" />
                  <rect x="400" y="90" width="200" height="20" rx="5" fill="#f1f3f9" />
                  
                  <!-- Dashboard stats -->
                  <rect x="100" y="140" width="120" height="80" rx="8" fill="#f1f3f9" />
                  <rect x="240" y="140" width="120" height="80" rx="8" fill="#f1f3f9" />
                  <rect x="380" y="140" width="120" height="80" rx="8" fill="#f1f3f9" />
                  <circle cx="130" cy="165" r="15" fill="#4361ee" />
                  <circle cx="270" cy="165" r="15" fill="#7209b7" />
                  <circle cx="410" cy="165" r="15" fill="#06d6a0" />
                  <rect x="115" y="190" width="90" height="8" rx="4" fill="#adb5bd" />
                  <rect x="115" y="205" width="60" height="8" rx="4" fill="#adb5bd" />
                  <rect x="255" y="190" width="90" height="8" rx="4" fill="#adb5bd" />
                  <rect x="255" y="205" width="60" height="8" rx="4" fill="#adb5bd" />
                  <rect x="395" y="190" width="90" height="8" rx="4" fill="#adb5bd" />
                  <rect x="395" y="205" width="60" height="8" rx="4" fill="#adb5bd" />
                  
                  <!-- Dashboard surveys list -->
                  <rect x="100" y="240" width="500" height="60" rx="8" fill="#f1f3f9" />
                  <rect x="100" y="310" width="500" height="60" rx="8" fill="#f1f3f9" />
                  <rect x="100" y="380" width="500" height="60" rx="8" fill="#f1f3f9" />
                  
                  <rect x="120" y="255" width="150" height="15" rx="3" fill="#adb5bd" />
                  <rect x="120" y="280" width="100" height="10" rx="2" fill="#adb5bd" />
                  <rect x="500" y="260" width="80" height="25" rx="12.5" fill="#4361ee" />
                  
                  <rect x="120" y="325" width="150" height="15" rx="3" fill="#adb5bd" />
                  <rect x="120" y="350" width="100" height="10" rx="2" fill="#adb5bd" />
                  <rect x="500" y="330" width="80" height="25" rx="12.5" fill="#7209b7" />
                  
                  <rect x="120" y="395" width="150" height="15" rx="3" fill="#adb5bd" />
                  <rect x="120" y="420" width="100" height="10" rx="2" fill="#adb5bd" />
                  <rect x="500" y="400" width="80" height="25" rx="12.5" fill="#06d6a0" />
                  
                  <!-- Survey form -->
                  <rect x="660" y="70" width="260" height="280" rx="12" fill="white" filter="drop-shadow(0 10px 15px rgba(0,0,0,0.2))" />
                  <rect x="680" y="90" width="220" height="25" rx="5" fill="#4361ee" />
                  <rect x="680" y="130" width="220" height="15" rx="3" fill="#f1f3f9" />
                  
                  <!-- Survey Form Radio Options -->
                  <circle cx="690" cy="170" r="8" fill="#4361ee" />
                  <rect x="710" y="165" width="170" height="10" rx="2" fill="#f1f3f9" />
                  
                  <circle cx="690" cy="200" r="8" fill="white" stroke="#4361ee" stroke-width="2" />
                  <rect x="710" y="195" width="170" height="10" rx="2" fill="#f1f3f9" />
                  
                  <circle cx="690" cy="230" r="8" fill="white" stroke="#4361ee" stroke-width="2" />
                  <rect x="710" y="225" width="170" height="10" rx="2" fill="#f1f3f9" />
                  
                  <circle cx="690" cy="260" r="8" fill="white" stroke="#4361ee" stroke-width="2" />
                  <rect x="710" y="255" width="170" height="10" rx="2" fill="#f1f3f9" />
                  
                  <rect x="730" y="290" width="120" height="40" rx="20" fill="#4361ee" />
                  
                  <!-- Analytics Panel -->
                  <rect x="660" y="370" width="260" height="180" rx="12" fill="white" filter="drop-shadow(0 10px 15px rgba(0,0,0,0.2))" />
                  <rect x="680" y="390" width="220" height="20" rx="4" fill="#f1f3f9" />
                  
                  <!-- Bar Chart -->
                  <rect x="690" y="430" width="30" height="100" rx="4" fill="#4361ee" opacity="0.8" />
                  <rect x="730" y="450" width="30" height="80" rx="4" fill="#7209b7" opacity="0.8" />
                  <rect x="770" y="420" width="30" height="110" rx="4" fill="#06d6a0" opacity="0.8" />
                  <rect x="810" y="470" width="30" height="60" rx="4" fill="#ff9f1c" opacity="0.8" />
                  
                  <!-- Animation elements (these will be animated with CSS) -->
                  <circle cx="500" cy="40" r="10" fill="white" opacity="0.3" class="float-element" />
                  <circle cx="600" cy="200" r="15" fill="white" opacity="0.2" class="float-element" />
                  <circle cx="700" cy="500" r="20" fill="white" opacity="0.1" class="float-element" />
                </svg>
              </div>
            </div>
          </v-card>
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
  box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1) !important;
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
  background-color: var(--v-surface-variant-base);
}

/* Survey Platform Preview Styles */
.survey-platform-preview {
  position: relative;
  width: 100%;
  height: 550px;
  overflow: hidden;
  border-radius: 16px;
}

.preview-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 1;
}

.preview-content {
  position: relative;
  z-index: 2;
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
}

.platform-illustration {
  width: 100%;
  height: 100%;
  filter: drop-shadow(0 10px 20px rgba(0,0,0,0.15));
}

/* Animation for floating elements */
.float-element {
  animation: float 8s ease-in-out infinite;
}

.float-element:nth-child(2) {
  animation-delay: 1s;
  animation-duration: 7s;
}

.float-element:nth-child(3) {
  animation-delay: 2s;
  animation-duration: 9s;
}

@keyframes float {
  0% {
    transform: translateY(0) translateX(0);
  }
  25% {
    transform: translateY(-15px) translateX(10px);
  }
  50% {
    transform: translateY(0) translateX(20px);
  }
  75% {
    transform: translateY(15px) translateX(10px);
  }
  100% {
    transform: translateY(0) translateX(0);
  }
}

/* Dashboard Preview */
.preview-dashboard {
  width: 65%;
  height: 100%;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 12px;
  padding: 20px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.preview-logo {
  width: 120px;
  height: 30px;
  background: var(--v-primary-base);
  border-radius: 6px;
}

.preview-nav {
  width: 200px;
  height: 20px;
  background: rgba(0, 0, 0, 0.1);
  border-radius: 10px;
}

.preview-stats {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.preview-stat-card {
  width: 30%;
  height: 80px;
  background: var(--v-surface-variant-lighten5);
  border-radius: 8px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.05);
}

.preview-surveys {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.preview-survey-item {
  height: 60px;
  background: var(--v-surface-base);
  border-radius: 8px;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.05);
}

/* Survey Form Preview */
.preview-survey-form {
  width: 30%;
  height: 60%;
  background: white;
  border-radius: 12px;
  padding: 15px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  display: flex;
  flex-direction: column;
  gap: 15px;
  margin-left: 20px;
}

.preview-survey-title {
  height: 30px;
  background: var(--v-primary-lighten5);
  border-radius: 6px;
  margin-bottom: 10px;
}

.preview-question {
  height: 60px;
  background: var(--v-surface-variant-lighten5);
  border-radius: 6px;
}

.preview-options {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.preview-option {
  height: 40px;
  background: rgba(0, 0, 0, 0.05);
  border-radius: 6px;
}

.preview-button {
  height: 40px;
  width: 60%;
  margin: 0 auto;
  background: var(--v-primary-base);
  border-radius: 20px;
}

/* Analytics Preview */
.preview-analytics {
  width: 30%;
  height: 35%;
  background: white;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.1);
  padding: 15px;
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-left: 20px;
  margin-top: 20px;
  align-self: flex-end;
}

.preview-chart-header {
  height: 20px;
  background: var(--v-info-lighten5);
  border-radius: 4px;
  margin-bottom: 5px;
}

.preview-chart {
  flex: 1;
  background: linear-gradient(135deg, var(--v-primary-lighten5), var(--v-secondary-lighten5));
  border-radius: 6px;
}

.preview-data-summary {
  height: 30px;
  background: var(--v-surface-variant-lighten5);
  border-radius: 4px;
}

@media (max-width: 960px) {
  .preview-dashboard {
    width: 100%;
    height: 60%;
  }
  
  .preview-survey-form, .preview-analytics {
    width: 48%;
    margin-left: 0;
    margin-top: 20px;
    height: 35%;
  }
}

@media (max-width: 600px) {
  .survey-platform-preview {
    height: 700px;
  }
  
  .preview-content {
    flex-direction: column;
  }
  
  .preview-dashboard, .preview-survey-form, .preview-analytics {
    width: 100%;
    margin-left: 0;
    height: auto;
    margin-bottom: 20px;
  }
  
  .preview-stats {
    flex-direction: column;
    gap: 10px;
  }
  
  .preview-stat-card {
    width: 100%;
  }
}
</style> 