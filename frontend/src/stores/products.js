import { defineStore } from "pinia";
import { ref } from "vue";
import { productsApi } from "@/api/products";

export const useProductsStore = defineStore("products", () => {
  const products = ref([]);
  const loading = ref(false);
  const error = ref(null);
  const totalCount = ref(0);

  // Filters
  const filters = ref({
    category_id: "",
    manufacturer_id: "",
    price_min: "",
    price_max: "",
    search: "",
  });

  const categories = ref([]);
  const manufacturers = ref([]);

  async function fetchProducts(params = {}) {
    loading.value = true;
    error.value = null;

    try {
      const response = await productsApi.getList(params);
      products.value = response.product || [];
      totalCount.value = response.total || products.value.length;
    } catch (err) {
      error.value = err.message || "Ошибка загрузки товаров";
    } finally {
      loading.value = false;
    }
  }

  async function fetchCategories() {
    try {
      const response = await productsApi.getCategories();
      categories.value = response.category || [];
    } catch (err) {
      // Тихая ошибка - категории не загрузились
    }
  }

  async function fetchManufacturers() {
    try {
      const response = await productsApi.getManufacturers();
      manufacturers.value = response.manufacturer || [];
    } catch (err) {
      // Тихая ошибка - производители не загрузились
    }
  }

  function setFilter(key, value) {
    filters.value[key] = value;
  }

  function resetFilters() {
    filters.value = {
      category_id: "",
      manufacturer_id: "",
      price_min: "",
      price_max: "",
      search: "",
    };
  }

  return {
    products,
    loading,
    error,
    totalCount,
    filters,
    categories,
    manufacturers,
    fetchProducts,
    fetchCategories,
    fetchManufacturers,
    setFilter,
    resetFilters,
  };
});
