<template>
  <div class="product-filters card">
    <div class="card-body">
      <h5 class="card-title mb-3">
        Фильтры
      </h5>
      
      <!-- Поиск -->
      <div class="mb-3">
        <label class="form-label">Поиск</label>
        <input 
          v-model="localFilters.search" 
          type="text" 
          class="form-control"
          placeholder="Название товара..."
          @input="applyFilters"
        >
      </div>
      
      <!-- Категория -->
      <div class="mb-3">
        <label class="form-label">Категория</label>
        <select 
          v-model="localFilters.category_id" 
          class="form-select"
          @change="applyFilters"
        >
          <option value="">
            Все категории
          </option>
          <option 
            v-for="category in categories" 
            :key="category.id" 
            :value="category.id"
          >
            {{ category.name }}
          </option>
        </select>
      </div>
      
      <!-- Производитель -->
      <div class="mb-3">
        <label class="form-label">Производитель</label>
        <select 
          v-model="localFilters.manufacturer_id" 
          class="form-select"
          @change="applyFilters"
        >
          <option value="">
            Все производители
          </option>
          <option 
            v-for="manufacturer in manufacturers" 
            :key="manufacturer.id" 
            :value="manufacturer.id"
          >
            {{ manufacturer.name }}
          </option>
        </select>
      </div>
      
      <!-- Цена от -->
      <div class="mb-3">
        <label class="form-label">Цена от</label>
        <input 
          v-model="localFilters.price_min" 
          type="number" 
          class="form-control"
          placeholder="0"
          min="0"
          @input="applyFilters"
        >
      </div>
      
      <!-- Цена до -->
      <div class="mb-3">
        <label class="form-label">Цена до</label>
        <input 
          v-model="localFilters.price_max" 
          type="number" 
          class="form-control"
          placeholder="∞"
          min="0"
          @input="applyFilters"
        >
      </div>
      
      <!-- Сбросить фильтры -->
      <button 
        class="btn btn-outline-secondary w-100" 
        @click="resetFilters"
      >
        Сбросить фильтры
      </button>
    </div>
  </div>
</template>

<script setup>
import { reactive, watch } from 'vue'

const props = defineProps({
  categories: {
    type: Array,
    default: () => []
  },
  manufacturers: {
    type: Array,
    default: () => []
  },
  filters: {
    type: Object,
    default: () => ({})
  }
})

const emit = defineEmits(['update:filters', 'apply'])

const localFilters = reactive({
  search: props.filters.search || '',
  category_id: props.filters.category_id || '',
  manufacturer_id: props.filters.manufacturer_id || '',
  price_min: props.filters.price_min || '',
  price_max: props.filters.price_max || ''
})

watch(() => props.filters, (newFilters) => {
  Object.assign(localFilters, newFilters)
}, { deep: true })

function applyFilters() {
  emit('update:filters', { ...localFilters })
  emit('apply', { ...localFilters })
}

function resetFilters() {
  Object.keys(localFilters).forEach(key => {
    localFilters[key] = ''
  })
  applyFilters()
}
</script>
