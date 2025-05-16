import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Import views
import HomePage from '@/views/HomePage.vue'
import LoginPage from '@/views/LoginPage.vue'
import RegisterPage from '@/views/RegisterPage.vue'
import ProfilePage from '@/views/ProfilePage.vue'
import MySurveysPage from '@/views/MySurveysPage.vue'
import SurveyEditorPage from '@/views/SurveyEditorPage.vue'
import SurveyTakePage from '@/views/SurveyTakePage.vue'
import SurveyResultsPage from '@/views/SurveyResultsPage.vue'
import SurveyResponsesPage from '@/views/SurveyResponsesPage.vue'
import NotFoundPage from '@/views/NotFoundPage.vue'

const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Home',
    component: HomePage,
    meta: { title: 'Home' }
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginPage,
    meta: { title: 'Login', requiresGuest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: RegisterPage,
    meta: { title: 'Register', requiresGuest: true }
  },
  {
    path: '/profile',
    name: 'Profile',
    component: ProfilePage,
    meta: { title: 'Profile', requiresAuth: true }
  },
  {
    path: '/my-surveys',
    name: 'MySurveys',
    component: MySurveysPage,
    meta: { title: 'My Surveys', requiresAuth: true }
  },
  {
    path: '/surveys/new',
    name: 'NewSurvey',
    component: SurveyEditorPage,
    meta: { title: 'Create Survey', requiresAuth: true }
  },
  {
    path: '/surveys/:id/edit',
    name: 'EditSurvey',
    component: SurveyEditorPage,
    meta: { title: 'Edit Survey', requiresAuth: true }
  },
  {
    path: '/surveys/:id/take',
    name: 'TakeSurvey',
    component: SurveyTakePage,
    meta: { title: 'Take Survey' }
  },
  {
    path: '/surveys/:id/results',
    name: 'SurveyResults',
    component: SurveyResultsPage,
    meta: { title: 'Survey Results', requiresAuth: true }
  },
  {
    path: '/surveys/:id/responses',
    name: 'SurveyResponses',
    component: SurveyResponsesPage,
    meta: { title: 'Survey Responses', requiresAuth: true }
  },
  {
    path: '/:catchAll(.*)',
    name: 'NotFound',
    component: NotFoundPage,
    meta: { title: 'Page Not Found' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  linkActiveClass: 'router-link-active'
})

// Navigation guards
router.beforeEach((to, from, next) => {
  // Update document title
  document.title = `${to.meta.title || 'Survey Platform'}`

  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest)

  if (requiresAuth && !authStore.isAuthenticated) {
    // If route requires authentication and user is not authenticated
    next('/login')
  } else if (requiresGuest && authStore.isAuthenticated) {
    // If route is for guests only and user is authenticated
    next('/')
  } else {
    // Proceed as normal
    next()
  }
})

export default router 