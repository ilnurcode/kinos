<template>
  <div class="home-view">
    <!-- Hero Section -->
    <section class="hero bg-light py-5">
      <div class="container">
        <div class="row align-items-center py-5">
          <div class="col-lg-6">
            <h1 class="display-4 fw-bold mb-4">Добро пожаловать в Kinos!</h1>
            <p class="lead mb-4">Ваш любимый интернет-магазин электроники и товаров для дома</p>
            <div class="d-flex gap-3">
              <router-link to="/catalog" class="btn btn-primary btn-lg">
                Перейти в каталог
              </router-link>
              <router-link v-if="!isAuthenticated" to="/register" class="btn btn-outline-success btn-lg">
                Регистрация
              </router-link>
            </div>
          </div>
          <div class="col-lg-6 text-center">
            <div class="display-1">🛍️</div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- Features Section -->
    <section class="features py-5">
      <div class="container">
        <h2 class="text-center mb-5">Почему выбирают нас</h2>
        <div class="row g-4">
          <div class="col-md-4">
            <div class="card h-100 border-0 shadow-sm">
              <div class="card-body text-center">
                <div class="display-4 mb-3">🚀</div>
                <h3 class="h5">Быстрая доставка</h3>
                <p class="text-muted">Доставим ваш заказ в кратчайшие сроки по всей стране</p>
              </div>
            </div>
          </div>
          <div class="col-md-4">
            <div class="card h-100 border-0 shadow-sm">
              <div class="card-body text-center">
                <div class="display-4 mb-3">🛡️</div>
                <h3 class="h5">Гарантия качества</h3>
                <p class="text-muted">Только оригинальная продукция с гарантией</p>
              </div>
            </div>
          </div>
          <div class="col-md-4">
            <div class="card h-100 border-0 shadow-sm">
              <div class="card-body text-center">
                <div class="display-4 mb-3">💳</div>
                <h3 class="h5">Удобная оплата</h3>
                <p class="text-muted">Различные способы оплаты: карта, наличные, рассрочка</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
    
    <!-- CTA Section -->
    <section class="cta bg-primary text-white py-5">
      <div class="container text-center">
        <h2 class="mb-4">Готовы сделать покупку?</h2>
        <p class="lead mb-4">Перейдите в каталог и выберите товары</p>
        <router-link to="/catalog" class="btn btn-light btn-lg">
          Смотреть каталог
        </router-link>
      </div>
    </section>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const isAuthenticated = computed(() => authStore.isAuthenticated)
</script>

<style scoped>
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.hero .btn-primary {
  background-color: white;
  border-color: white;
  color: #667eea;
}

.hero .btn-outline-success {
  border-color: white;
  color: white;
}

.hero .btn-outline-success:hover {
  background-color: white;
  color: #667eea;
}

.features {
  background-color: #f8f9fa;
}

.cta {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
</style>
