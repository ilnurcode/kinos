import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { productsApi } from '@/api/products'

export const useProductsStore = defineStore('products', () => {
  const products = ref([])
  const loading = ref(false)
  const error = ref(null)
  const totalCount = ref(0)
  
  // Filters
  const filters = ref({
    category_id: '',
    manufacturer_id: '',
    price_min: '',
    price_max: '',
    search: ''
  })
  
  const categories = ref([])
  const manufacturers = ref([])
  
  async function fetchProducts(params = {}) {
    loading.value = true
    error.value = null
    
    try {
      const response = await productsApi.getList(params)
      products.value = response.product || []
      totalCount.value = response.total || products.value.length
    } catch (err) {
      error.value = err.message || 'Ошибка загрузки товаров'
    } finally {
      loading.value = false
    }
  }
  
  async function fetchCategories() {
    try {
      const response = await productsApi.getCategories()
      categories.value = response.category || []
    } catch (err) {
      console.error('Ошибка загрузки категорий:', err)
    }
  }
  
  async function fetchManufacturers() {
    try {
      const response = await productsApi.getManufacturers()
      manufacturers.value = response.manufacturer || []
    } catch (err) {
      console.error('Ошибка загрузки производителей:', err)
    }
  }
  
  function setFilter(key, value) {
    filters.value[key] = value
  }
  
  function resetFilters() {
    filters.value = {
      category_id: '',
      manufacturer_id: '',
      price_min: '',
      price_max: '',
      search: ''
    }
  }
  
  const filteredProducts = computed(() => {
    let result = products.value
    
    if (filters.value.category_id) {
      result = result.filter(p => p.category_id == filters.value.category_id)
    }
    
    if (filters.value.manufacturer_id) {
      result = result.filter(p => p.manufacturer_id == filters.value.manufacturer_id)
    }
    
    if (filters.value.price_min) {
      result = result.filter(p => p.price >= parseFloat(filters.value.price_min))
    }
    
    if (filters.value.price_max) {
      result = result.filter(p => p.price <= parseFloat(filters.value.price_max))
    }
    
    if (filters.value.search) {
      const search = filters.value.search.toLowerCase()
      result = result.filter(p => p.name.toLowerCase().includes(search))
    }
    
    return result
  })
  
  return {
    products,
    loading,
    error,
    totalCount,
    filters,
    categories,
    manufacturers,
    filteredProducts,
    fetchProducts,
    fetchCategories,
    fetchManufacturers,
    setFilter,
    resetFilters
  }
})
