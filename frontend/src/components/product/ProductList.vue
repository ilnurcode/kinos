<template>
  <div class="product-list">
    <div
      v-if="loading"
      class="text-center py-5"
    >
      <LoadingSpinner text="Загрузка товаров..." />
    </div>

    <div
      v-else-if="error"
      class="alert alert-danger"
    >
      {{ error }}
    </div>

    <div
      v-else-if="products.length === 0"
      class="text-center py-5"
    >
      <div class="display-1">
        📦
      </div>
      <h3>Товары не найдены</h3>
      <p class="text-muted">
        Попробуйте изменить параметры фильтров
      </p>
    </div>

    <div
      v-else
      class="row g-4"
    >
      <div
        v-for="product in products"
        :key="product.id"
        class="col-md-6 col-lg-4"
      >
        <ProductCard
          :product="product"
          :adding="addingId === product.id"
          @add-to-cart="$emit('add-to-cart', $event)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import ProductCard from "./ProductCard.vue";
import LoadingSpinner from "@/components/common/LoadingSpinner.vue";

defineProps({
  products: {
    type: Array,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
  error: {
    type: String,
    default: null,
  },
  addingId: {
    type: Number,
    default: null,
  },
});

defineEmits(["add-to-cart"]);
</script>
