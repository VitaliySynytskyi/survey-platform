import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './utils/axiosConfig' // Import axios configuration

// Import global styles
import './assets/styles/global.css'

// Vuetify
import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'
import { createVuetify } from 'vuetify'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'
import 'animate.css'

// Custom theme with modern colors
const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light',
    themes: {
      light: {
        dark: false,
        colors: {
          primary: '#4361ee',      // Modern blue
          secondary: '#3a0ca3',    // Deep purple
          accent: '#4cc9f0',       // Light blue
          error: '#ef476f',        // Soft red
          warning: '#ffd166',      // Yellow
          info: '#118ab2',         // Teal
          success: '#06d6a0',      // Mint green
          background: '#f8f9fa',   // Light background
          surface: '#ffffff'       // White surface
        }
      },
      dark: {
        dark: true,
        colors: {
          primary: '#4cc9f0',      // Light blue
          secondary: '#7209b7',    // Purple
          accent: '#4361ee',       // Blue
          error: '#ef476f',        // Soft red
          warning: '#ffd166',      // Yellow
          info: '#118ab2',         // Teal
          success: '#06d6a0',      // Mint green
          background: '#121212',   // Dark background
          surface: '#1e1e1e'       // Dark surface
        }
      }
    }
  },
  icons: {
    defaultSet: 'mdi'
  }
})

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(vuetify)

app.mount('#app') 