<template>
  <div class="product-card card h-100">
    <div class="card-body">
      <h5 class="card-title">{{ product.name }}</h5>
      <p class="card-text text-muted">
        {{ product.category_name || 'Категория не указана' }}
      </p>
      <div class="d-flex justify-content-between align-items-center">
        <span class="h4 mb-0 text-primary">{{ formatPrice(product.price) }}</span>
        <button class="btn btn-outline-primary btn-sm">
          Подробнее
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  product: {
    type: Object,
    required: true
  }
})

function formatPrice(price) {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 0
  }).format(price)
}
</script>

<style scoped>
.product-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.product-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}
</style>
