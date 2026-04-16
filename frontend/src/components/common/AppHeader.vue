<template>
  <header class="navbar navbar-expand-lg navbar-dark bg-dark">
    <div class="container">
      <router-link
        class="navbar-brand fw-bold"
        to="/"
      >
        Kinos
      </router-link>

      <button
        class="navbar-toggler"
        type="button"
        data-bs-toggle="collapse"
        data-bs-target="#navbarNav"
      >
        <span class="navbar-toggler-icon" />
      </button>

      <div
        id="navbarNav"
        class="collapse navbar-collapse"
      >
        <ul class="navbar-nav me-auto">
          <li class="nav-item">
            <router-link
              class="nav-link"
              to="/catalog"
            >
              Каталог
            </router-link>
          </li>
          <li
            v-if="isAuthenticated"
            class="nav-item"
          >
            <router-link
              class="nav-link d-inline-flex align-items-center gap-1"
              to="/cart"
            >
              Корзина
              <span
                v-if="cartCount > 0"
                class="badge text-bg-light"
              >
                {{ cartCount }}
              </span>
            </router-link>
          </li>
          <li
            v-if="isAdmin"
            class="nav-item"
          >
            <router-link
              class="nav-link text-warning"
              to="/admin/dashboard"
            >
              Админка
            </router-link>
          </li>
        </ul>

        <ul class="navbar-nav">
          <template v-if="!isAuthenticated">
            <li class="nav-item">
              <router-link
                class="nav-link"
                to="/login"
              >
                Вход
              </router-link>
            </li>
            <li class="nav-item">
              <router-link
                class="nav-link"
                to="/register"
              >
                Регистрация
              </router-link>
            </li>
          </template>

          <template v-else>
            <li class="nav-item dropdown">
              <a
                id="navbarDropdown"
                class="nav-link dropdown-toggle d-flex align-items-center gap-1"
                href="javascript:void(0)"
                role="button"
                data-bs-toggle="dropdown"
                aria-expanded="false"
              >
                <span>{{ user?.username || "Профиль" }}</span>
                <span
                  v-if="isAdmin"
                  class="badge bg-warning text-dark"
                >
                  Admin
                </span>
              </a>
              <ul
                class="dropdown-menu dropdown-menu-end"
                aria-labelledby="navbarDropdown"
              >
                <li>
                  <router-link
                    class="dropdown-item"
                    to="/profile"
                    @click="closeDropdown"
                  >
                    Профиль
                  </router-link>
                </li>
                <li>
                  <router-link
                    class="dropdown-item"
                    to="/profile/edit"
                    @click="closeDropdown"
                  >
                    Настройки
                  </router-link>
                </li>
                <li><hr class="dropdown-divider"></li>
                <li>
                  <a
                    class="dropdown-item"
                    href="javascript:void(0)"
                    @click.prevent="handleLogout"
                  >
                    Выйти
                  </a>
                </li>
              </ul>
            </li>
          </template>
        </ul>
      </div>
    </div>
  </header>
</template>

<script setup>
import { computed, watch } from "vue";
import { useRouter } from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useCartStore } from "@/stores/cart";

const router = useRouter();
const authStore = useAuthStore();
const cartStore = useCartStore();

const isAuthenticated = computed(() => authStore.isAuthenticated);
const isAdmin = computed(() => authStore.isAdmin);
const user = computed(() => authStore.user);
const cartCount = computed(() => cartStore.count || cartStore.itemsCount);

watch(
  () => authStore.isAuthenticated,
  (authenticated) => {
    if (authenticated) {
      cartStore.fetchItemsCount();
      return;
    }

    cartStore.count = 0;
    cartStore.items = [];
    cartStore.total = 0;
  },
  { immediate: true }
);

const handleLogout = () => {
  authStore.logout();
  router.push("/");
};

const closeDropdown = () => {
  const dropdownElement = document.querySelector('.dropdown-toggle[aria-expanded="true"]');
  if (dropdownElement) {
    dropdownElement.click();
  }
};
</script>

<style scoped>
.navbar-brand {
  font-size: 1.5rem;
}

.nav-link.dropdown-toggle {
  display: inline-flex !important;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
}

.nav-link.dropdown-toggle::after {
  margin-left: 0.25rem;
  border-top: 0.3em solid;
  border-right: 0.3em solid transparent;
  border-bottom: 0;
  border-left: 0.3em solid transparent;
  vertical-align: middle;
}
</style>
