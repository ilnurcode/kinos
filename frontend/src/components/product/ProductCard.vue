<template>
  <div class="product-card card h-100">
    <div class="card-body">
      <h5 class="card-title">
        {{ product.name }}
      </h5>
      <p class="card-text text-muted">
        {{ product.category_name || "Категория не указана" }}
      </p>
      <div class="d-flex justify-content-between align-items-center gap-2">
        <span class="h4 mb-0 text-primary">{{ formatPrice(product.price) }}</span>
        <button
          class="btn btn-primary btn-sm"
          :disabled="adding"
          @click="$emit('add-to-cart', product)"
        >
          {{ adding ? "Добавление..." : "В корзину" }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { formatPrice } from "@/utils/format";

defineProps({
  product: {
    type: Object,
    required: true,
  },
  adding: {
    type: Boolean,
    default: false,
  },
});

defineEmits(["add-to-cart"]);
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
