<template>
  <div class="app-container">
    <header class="main-header">
      <nav class="navbar navbar-expand-lg navbar-dark">
        <div class="container">
          <router-link class="navbar-brand" to="/">
            <i class="bi bi-clipboard-data me-2"></i>Survey Platform
          </router-link>
          
          <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNav">
            <span class="navbar-toggler-icon"></span>
          </button>
          
          <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav me-auto">
              <li class="nav-item">
                <router-link class="nav-link" to="/">Home</router-link>
              </li>
              <li class="nav-item" v-if="isAuthenticated">
                <router-link class="nav-link" to="/my-surveys">My Surveys</router-link>
              </li>
            </ul>
            
            <ul class="navbar-nav">
              <template v-if="isAuthenticated">
                <li class="nav-item dropdown">
                  <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                    <i class="bi bi-person-circle me-1"></i>{{ user ? user.username : 'User' }}
                  </a>
                  <ul class="dropdown-menu dropdown-menu-end">
                    <li>
                      <router-link class="dropdown-item" to="/profile">
                        <i class="bi bi-person me-2"></i>Profile
                      </router-link>
                    </li>
                    <li><hr class="dropdown-divider"></li>
                    <li>
                      <a class="dropdown-item" href="#" @click.prevent="logout">
                        <i class="bi bi-box-arrow-right me-2"></i>Logout
                      </a>
                    </li>
                  </ul>
                </li>
              </template>
              <template v-else>
                <li class="nav-item">
                  <router-link class="nav-link" to="/login">Log In</router-link>
                </li>
                <li class="nav-item">
                  <router-link class="nav-link btn btn-primary btn-sm text-white ms-2 px-3" to="/register">Sign Up</router-link>
                </li>
              </template>
            </ul>
          </div>
        </div>
      </nav>
    </header>

    <main class="main-content">
      <router-view />
    </main>

    <footer class="main-footer">
      <div class="container">
        <div class="row">
          <div class="col-md-6">
            <p>&copy; {{ new Date().getFullYear() }} Survey Platform. All rights reserved.</p>
          </div>
          <div class="col-md-6 text-end">
            <ul class="footer-links">
              <li><a href="#">Privacy Policy</a></li>
              <li><a href="#">Terms of Service</a></li>
              <li><a href="#">Contact</a></li>
            </ul>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from './stores/auth'

export default {
  name: 'App',
  
  setup() {
    const router = useRouter()
    const authStore = useAuthStore()
    
    // Computed properties
    const isAuthenticated = computed(() => authStore.isAuthenticated)
    const user = computed(() => authStore.user)
    
    // Methods
    const logout = async () => {
      try {
        await authStore.logout()
        router.push('/login')
      } catch (error) {
        console.error('Logout error:', error)
      }
    }
    
    // Initialize the app
    onMounted(() => {
      // Auto login from saved token if available
      if (localStorage.getItem('token') && !authStore.isAuthenticated) {
        authStore.autoLogin()
      }
      
      // Initialize Bootstrap tooltips and popovers
      document.addEventListener('DOMContentLoaded', () => {
        const tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'))
        tooltipTriggerList.map(function (tooltipTriggerEl) {
          return new bootstrap.Tooltip(tooltipTriggerEl)
        })
      })
    })
    
    return {
      isAuthenticated,
      user,
      logout
    }
  }
}
</script>

<style>
/* Global Styles */
:root {
  --primary-color: #4a6bef;
  --secondary-color: #32325d;
  --accent-color: #7e57c2;
  --light-color: #f8f9fa;
  --dark-color: #212529;
  --success-color: #10b981;
  --danger-color: #ef4444;
  --warning-color: #f59e0b;
}

body {
  font-family: 'Inter', 'Segoe UI', Roboto, sans-serif;
  color: var(--dark-color);
  background-color: #f5f7fa;
  min-height: 100vh;
}

.app-container {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

/* Header Styles */
.main-header {
  background: linear-gradient(135deg, var(--primary-color), var(--accent-color));
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.navbar-brand {
  font-weight: 700;
  letter-spacing: 0.5px;
}

.navbar-dark .navbar-nav .nav-link {
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
  padding: 0.75rem 1rem;
}

.navbar-dark .navbar-nav .nav-link:hover {
  color: #fff;
}

.navbar-dark .navbar-nav .router-link-active {
  color: #fff;
  position: relative;
}

.navbar-dark .navbar-nav .router-link-active:after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 3px;
  background-color: white;
  border-radius: 3px;
}

/* Main Content Styles */
.main-content {
  flex: 1;
  padding: 2rem 0;
}

/* Footer Styles */
.main-footer {
  background-color: var(--secondary-color);
  color: rgba(255, 255, 255, 0.8);
  padding: 2rem 0;
  margin-top: auto;
}

.footer-links {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  justify-content: flex-end;
}

.footer-links li {
  margin-left: 1.5rem;
}

.footer-links a {
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  transition: color 0.2s;
}

.footer-links a:hover {
  color: #fff;
}

/* Button Styles */
.btn-primary {
  background-color: var(--primary-color);
  border-color: var(--primary-color);
}

.btn-primary:hover,
.btn-primary:focus {
  background-color: #3b5be0;
  border-color: #3b5be0;
  box-shadow: 0 2px 8px rgba(74, 107, 239, 0.25);
}

/* Form controls */
.form-control:focus {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 0.25rem rgba(74, 107, 239, 0.25);
}

/* Animation styles */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .footer-links {
    justify-content: center;
    margin-top: 1rem;
  }
  
  .col-md-6.text-end {
    text-align: center !important;
  }
  
  .col-md-6:first-child p {
    text-align: center;
    margin-bottom: 1rem;
  }
}
</style> 