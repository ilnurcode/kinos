<template>
  <div class="catalog-view py-4">
    <div class="container">
      <h1 class="mb-4">Каталог товаров</h1>

      <div class="row">
        <!-- Filters Sidebar -->
        <div class="col-lg-3 mb-4">
          <ProductFilters
            :categories="categories"
            :manufacturers="manufacturers"
            :filters="filters"
            @update:filters="updateFilters"
            @apply="applyFilters"
          />
        </div>

        <!-- Products Grid -->
        <div class="col-lg-9">
          <div class="d-flex justify-content-between align-items-center mb-4">
            <p class="mb-0">Найдено товаров: <strong>{{ totalCount }}</strong></p>
          </div>

          <ProductList
            :products="products"
            :loading="loading"
            :error="error"
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
import { ref, reactive, onMounted, computed } from 'vue'
import { useProducts } from '@/composables/useProducts'
import ProductFilters from '@/components/product/ProductFilters.vue'
import ProductList from '@/components/product/ProductList.vue'
import Pagination from '@/components/common/Pagination.vue'

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
  setFilter
} = useProducts()

const filters = reactive({
  category_id: '',
  manufacturer_id: '',
  price_min: '',
  price_max: '',
  search: ''
})

const page = ref(1)
const limit = 12

const totalPages = computed(() => Math.ceil(totalCount.value / limit))

const applyFilters = async () => {
  page.value = 1
  await loadProducts({
    limit,
    offset: 0,
    ...filters
  })
}

const updateFilters = (newFilters) => {
  Object.assign(filters, newFilters)
}

const changePage = (newPage) => {
  page.value = newPage
  loadProducts({
    limit,
    offset: (newPage - 1) * limit,
    ...filters
  })
}

onMounted(async () => {
  await Promise.all([
    loadProducts({ limit, offset: 0 }),
    loadCategories(),
    loadManufacturers()
  ])
})
</script>
