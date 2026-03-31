import { ref, computed } from 'vue'
import { useProductsStore } from '@/stores/products'

export function useProducts() {
  const productsStore = useProductsStore()
  
  const products = computed(() => productsStore.products)
  const loading = computed(() => productsStore.loading)
  const error = computed(() => productsStore.error)
  const totalCount = computed(() => productsStore.totalCount)
  const categories = computed(() => productsStore.categories)
  const manufacturers = computed(() => productsStore.manufacturers)
  const filters = computed(() => productsStore.filters)
  
  async function loadProducts(params = {}) {
    await productsStore.fetchProducts(params)
  }
  
  async function loadCategories() {
    await productsStore.fetchCategories()
  }
  
  async function loadManufacturers() {
    await productsStore.fetchManufacturers()
  }
  
  function setFilter(key, value) {
    productsStore.setFilter(key, value)
  }
  
  function resetFilters() {
    productsStore.resetFilters()
  }
  
  return {
    products,
    loading,
    error,
    totalCount,
    categories,
    manufacturers,
    filters,
    loadProducts,
    loadCategories,
    loadManufacturers,
    setFilter,
    resetFilters
  }
}
