import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../store/auth'

// Lazy-loaded components
const Home = () => import('../views/Home.vue')
const Login = () => import('../views/Login.vue')
const Register = () => import('../views/Register.vue')
const Dashboard = () => import('../views/Dashboard.vue')
const CreateSurvey = () => import('../views/CreateSurvey.vue')
const EditSurvey = () => import('../views/EditSurvey.vue')
const SurveyResponses = () => import('../views/SurveyResponses.vue')
const TakeSurvey = () => import('../views/TakeSurvey.vue')
const NotFound = () => import('../views/NotFound.vue')
const SurveySuccess = () => import('../views/SurveySuccess.vue')

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: { requiresGuest: true }
  },
  {
    path: '/register',
    name: 'Register',
    component: Register,
    meta: { requiresGuest: true }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: Dashboard,
    meta: { requiresAuth: true }
  },
  {
    path: '/surveys/create',
    name: 'CreateSurvey',
    component: CreateSurvey,
    meta: { requiresAuth: true }
  },
  {
    path: '/surveys/:id/edit',
    name: 'EditSurvey',
    component: EditSurvey,
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/surveys/:id/responses',
    name: 'SurveyResponses',
    component: SurveyResponses,
    props: true,
    meta: { requiresAuth: true }
  },
  {
    path: '/surveys/:id',
    name: 'TakeSurvey',
    component: TakeSurvey,
    props: true
  },
  {
    path: '/surveys/:id/success',
    name: 'SurveySuccess',
    component: SurveySuccess,
    props: true,
    meta: { requiresAuth: false }
  },
  {
    path: '/surveys/:id/analytics',
    name: 'SurveyAnalytics',
    component: () => import('../views/SurveyAnalytics.vue'),
    meta: { requiresAuth: true }
  },
  // {
  //   path: '/profile',
  //   name: 'Profile',
  //   component: () => import('../views/Profile.vue'),
  //   meta: { requiresAuth: true }
  // },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: NotFound
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// Navigation guards
router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  const requiresAuth = to.matched.some(record => record.meta.requiresAuth)
  const requiresGuest = to.matched.some(record => record.meta.requiresGuest)

  if (requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'Login', query: { redirect: to.fullPath } })
  } else if (requiresGuest && authStore.isAuthenticated) {
    next({ name: 'Dashboard' })
  } else {
    next()
  }
})

export default router 