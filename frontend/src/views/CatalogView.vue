<template>
  <div class="catalog-view py-4">
    <div class="container">
      <h1 class="mb-4">
        Каталог товаров
      </h1>

      <div class="row">
        <div class="col-lg-3 mb-4">
          <ProductFilters
            :categories="categories"
            :manufacturers="manufacturers"
            :filters="filters"
            @update:filters="updateFilters"
            @apply="applyFilters"
          />
        </div>

        <div class="col-lg-9">
          <div class="d-flex justify-content-between align-items-center mb-4">
            <p class="mb-0">
              Найдено товаров: <strong>{{ totalCount }}</strong>
            </p>
          </div>

          <ProductList
            :products="products"
            :loading="loading"
            :error="error"
            :adding-id="addingProductId"
            @add-to-cart="handleAddToCart"
          />

          <Pagination
            v-if="totalPages > 1"
            :current-page="page"
            :total-pages="totalPages"
            @page-change="changePage"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import ProductFilters from "@/components/product/ProductFilters.vue";
import ProductList from "@/components/product/ProductList.vue";
import Pagination from "@/components/common/Pagination.vue";
import { useProducts } from "@/composables/useProducts";
import { useCartStore } from "@/stores/cart";
import { useAuthStore } from "@/stores/auth";

const router = useRouter();
const cartStore = useCartStore();
const authStore = useAuthStore();

const {
  products,
  loading,
  error,
  totalCount,
  categories,
  manufacturers,
  loadProducts,
  loadCategories,
  loadManufacturers,
} = useProducts();

const filters = reactive({
  category_id: "",
  manufacturer_id: "",
  price_min: "",
  price_max: "",
  search: "",
});

const page = ref(1);
const limit = 12;
const addingProductId = ref(null);

const totalPages = computed(() => Math.ceil(totalCount.value / limit));

const applyFilters = async () => {
  page.value = 1;
  await loadProducts({
    limit,
    offset: 0,
    ...filters,
  });
};

const updateFilters = (newFilters) => {
  Object.assign(filters, newFilters);
};

const changePage = (newPage) => {
  page.value = newPage;
  loadProducts({
    limit,
    offset: (newPage - 1) * limit,
    ...filters,
  });
};

const handleAddToCart = async (product) => {
  if (!authStore.isAuthenticated) {
    router.push({ name: "login", query: { redirect: router.currentRoute.value.fullPath } });
    return;
  }

  addingProductId.value = product.id;

  try {
    await cartStore.addItem(product.id, 1);
  } finally {
    addingProductId.value = null;
  }
};

onMounted(async () => {
  await Promise.all([
    loadProducts({ limit, offset: 0 }),
    loadCategories(),
    loadManufacturers(),
  ]);
});
</script>
