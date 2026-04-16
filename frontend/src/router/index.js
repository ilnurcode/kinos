import { useAuthStore } from "@/stores/auth";
import { createRouter, createWebHistory } from "vue-router";

const routes = [
  {
    path: "/",
    name: "home",
    component: () => import("@/views/HomeView.vue"),
  },
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/LoginView.vue"),
    meta: { guest: true },
  },
  {
    path: "/register",
    name: "register",
    component: () => import("@/views/RegisterView.vue"),
    meta: { guest: true },
  },
  {
    path: "/profile",
    name: "profile",
    component: () => import("@/views/ProfileView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/profile/edit",
    name: "profile-edit",
    component: () => import("@/views/ProfileEditView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/catalog",
    name: "catalog",
    component: () => import("@/views/CatalogView.vue"),
  },
  {
    path: "/cart",
    name: "cart",
    component: () => import("@/views/CartView.vue"),
    meta: { requiresAuth: true },
  },
  {
    path: "/admin",
    name: "admin",
    redirect: "/admin/dashboard",
    meta: { requiresAuth: true },
    children: [
      {
        path: "dashboard",
        name: "admin-dashboard",
        component: () => import("@/views/admin/DashboardView.vue"),
        meta: { requiresAdmin: true },
      },
      {
        path: "users",
        name: "admin-users",
        component: () => import("@/views/admin/UsersView.vue"),
        meta: { requiresAdmin: true },
      },
      {
        path: "categories",
        name: "admin-categories",
        component: () => import("@/views/admin/CategoriesView.vue"),
        meta: { requiresAdmin: true },
      },
      {
        path: "manufacturers",
        name: "admin-manufacturers",
        component: () => import("@/views/admin/ManufacturersView.vue"),
        meta: { requiresAdmin: true },
      },
      {
        path: "products",
        name: "admin-products",
        component: () => import("@/views/admin/ProductsView.vue"),
        meta: { requiresAdmin: true },
      },
      {
        path: "inventory",
        name: "admin-inventory",
        component: () => import("@/views/admin/InventoryView.vue"),
        meta: { requiresAdmin: true },
      },
      {
        path: "warehouses",
        name: "admin-warehouses",
        component: () => import("@/views/admin/WarehousesView.vue"),
        meta: { requiresAdmin: true },
      },
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition;
    } else {
      return { top: 0 };
    }
  },
});

// Navigation guards
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore();
  await authStore.initializeAuth();

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: "login", query: { redirect: to.fullPath } });
    return;
  }

  if (to.meta.requiresAdmin && !authStore.isAdmin) {
    next({ name: "home" });
    return;
  }

  if (to.meta.guest && authStore.isAuthenticated) {
    next({ name: "home" });
    return;
  }

  next();
});

export default router;
