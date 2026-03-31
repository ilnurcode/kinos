import { useAuthStore } from '@/stores/auth'
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'home',
    component: () => import('@/views/HomeView.vue')
  },
  {
    path: '/login',
    name: 'login',
    component: () => import('@/views/LoginView.vue'),
    meta: { guest: true }
  },
  {
    path: '/register',
    name: 'register',
    component: () => import('@/views/RegisterView.vue'),
    meta: { guest: true }
  },
  {
    path: '/profile',
    name: 'profile',
    component: () => import('@/views/ProfileView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/profile/edit',
    name: 'profile-edit',
    component: () => import('@/views/ProfileEditView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/catalog',
    name: 'catalog',
    component: () => import('@/views/CatalogView.vue')
  },
  {
    path: '/admin',
    name: 'admin',
    redirect: '/admin/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: 'dashboard',
        name: 'admin-dashboard',
        component: () => import('@/views/admin/DashboardView.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'users',
        name: 'admin-users',
        component: () => import('@/views/admin/UsersView.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'categories',
        name: 'admin-categories',
        component: () => import('@/views/admin/CategoriesView.vue'),
        meta: { requiresAdmin: true }
      },
      {
        path: 'manufacturers',
        name: 'admin-manufacturers',
        component: () => import('@/views/admin/ManufacturersView.vue')
      },
      {
        path: 'products',
        name: 'admin-products',
        component: () => import('@/views/admin/ProductsView.vue')
      },
      {
        path: 'inventory',
        name: 'admin-inventory',
        component: () => import('@/views/admin/InventoryView.vue')
      },
      {
        path: 'warehouses',
        name: 'admin-warehouses',
        component: () => import('@/views/admin/WarehousesView.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // Проверяем токен если это первый заход
  const token = localStorage.getItem('access_token')
  if (token && !authStore.isAuthenticated) {
    try {
      await authStore.fetchProfile()
    } catch (err) {
      console.error('Failed to restore session:', err)
      authStore.logout()
    }
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: 'home' })
    return
  }

  if (to.meta.guest && authStore.isAuthenticated) {
    next({ name: 'home' })
    return
  }

  next()
})

export default router
