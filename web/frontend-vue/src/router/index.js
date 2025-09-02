import { createRouter, createWebHistory } from 'vue-router'
import DashboardView from '../views/DashboardView.vue'
import LoginView from '../views/LoginView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: DashboardView,
      meta: { requiresAuth: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    }
  ]
})

// Navigation Guard
router.beforeEach((to, from, next) => {
  const loggedIn = localStorage.getItem('jwt_token');

  if (to.meta.requiresAuth && !loggedIn) {
    next('/login');
  } else {
    next(); // Lanjutkan ke rute yang dituju
  }
});

export default router