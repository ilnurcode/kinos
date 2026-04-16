<template>
  <div class="admin-dashboard py-4">
    <div class="container">
      <h1 class="mb-4">
        Админ-панель
      </h1>
      
      <div class="row g-4">
        <!-- Пользователи -->
        <div class="col-md-6 col-lg-3">
          <router-link
            to="/admin/users"
            class="text-decoration-none"
          >
            <div class="card h-100 shadow-sm hover-card">
              <div class="card-body text-center">
                <div class="display-4 mb-3">
                  👥
                </div>
                <h5 class="card-title">
                  Пользователи
                </h5>
                <p class="card-text text-muted">
                  Управление пользователями и ролями
                </p>
              </div>
            </div>
          </router-link>
        </div>

        <!-- Категории -->
        <div class="col-md-6 col-lg-3">
          <router-link
            to="/admin/categories"
            class="text-decoration-none"
          >
            <div class="card h-100 shadow-sm hover-card">
              <div class="card-body text-center">
                <div class="display-4 mb-3">
                  📁
                </div>
                <h5 class="card-title">
                  Категории
                </h5>
                <p class="card-text text-muted">
                  Управление категориями товаров
                </p>
              </div>
            </div>
          </router-link>
        </div>

        <!-- Производители -->
        <div class="col-md-6 col-lg-3">
          <router-link
            to="/admin/manufacturers"
            class="text-decoration-none"
          >
            <div class="card h-100 shadow-sm hover-card">
              <div class="card-body text-center">
                <div class="display-4 mb-3">
                  🏭
                </div>
                <h5 class="card-title">
                  Производители
                </h5>
                <p class="card-text text-muted">
                  Управление производителями
                </p>
              </div>
            </div>
          </router-link>
        </div>

        <!-- Товары -->
        <div class="col-md-6 col-lg-3">
          <router-link
            to="/admin/products"
            class="text-decoration-none"
          >
            <div class="card h-100 shadow-sm hover-card">
              <div class="card-body text-center">
                <div class="display-4 mb-3">
                  📦
                </div>
                <h5 class="card-title">
                  Товары
                </h5>
                <p class="card-text text-muted">
                  Управление товарами
                </p>
              </div>
            </div>
          </router-link>
        </div>

        <!-- Запасы -->
        <div class="col-md-6 col-lg-3">
          <router-link
            to="/admin/inventory"
            class="text-decoration-none"
          >
            <div class="card h-100 shadow-sm hover-card">
              <div class="card-body text-center">
                <div class="display-4 mb-3">
                  📊
                </div>
                <h5 class="card-title">
                  Запасы
                </h5>
                <p class="card-text text-muted">
                  Управление запасами товаров
                </p>
              </div>
            </div>
          </router-link>
        </div>

        <!-- Склады -->
        <div class="col-md-6 col-lg-3">
          <router-link
            to="/admin/warehouses"
            class="text-decoration-none"
          >
            <div class="card h-100 shadow-sm hover-card">
              <div class="card-body text-center">
                <div class="display-4 mb-3">
                  🏢
                </div>
                <h5 class="card-title">
                  Склады
                </h5>
                <p class="card-text text-muted">
                  Управление складами
                </p>
              </div>
            </div>
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
</script>

<style scoped>
.hover-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.hover-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
