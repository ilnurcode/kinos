<template>
  <div class="admin-products py-4">
    <div class="container">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление товарами</h1>
        <button class="btn btn-primary" @click="openCreateModal">
          <i class="bi bi-plus-lg"></i> Добавить товар
        </button>
      </div>

      <ErrorAlert v-if="error" :message="error" @dismiss="error = null" />

      <div v-if="loading" class="text-center py-5">
        <LoadingSpinner text="Загрузка товаров..." />
      </div>

      <div v-else-if="products.length === 0" class="alert alert-info">
        Товары не найдены
      </div>

      <div v-else class="card">
        <div class="card-body">
          <table class="table table-hover">
            <thead class="table-light">
              <tr>
                <th>ID</th>
                <th>Название</th>
                <th>Категория</th>
                <th>Производитель</th>
                <th>Цена</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="product in products" :key="product.id">
                <td>{{ product.id }}</td>
                <td><strong>{{ product.name }}</strong></td>
                <td>{{ getCategoryName(product.category_id) }}</td>
                <td>{{ getManufacturerName(product.manufacturer_id) }}</td>
                <td>{{ formatPrice(product.price) }}</td>
                <td>
                  <button class="btn btn-sm btn-outline-primary me-2" @click="openEditModal(product)">
                    ✏️ Редактировать
                  </button>
                  <button class="btn btn-sm btn-outline-danger" @click="confirmDelete(product)">
                    🗑️ Удалить
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>

  <!-- Модальное окно -->
  <div class="modal fade" :class="{ show: showModal }" :style="modalStyle" tabindex="-1" @click.self="closeModal">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Редактировать товар' : 'Новый товар' }}</h5>
          <button type="button" class="btn-close" @click="closeModal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveProduct">
            <div class="row">
              <div class="col-md-12 mb-3">
                <label class="form-label">Название товара *</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.name"
                  required
                  placeholder="Например: Смартфон Samsung Galaxy S24"
                />
              </div>
              <div class="col-md-4 mb-3">
                <label class="form-label">Категория *</label>
                <select class="form-select" v-model="form.category_id" required>
                  <option value="">Выберите категорию</option>
                  <option v-for="cat in categories" :key="cat.id" :value="cat.id">
                    {{ cat.name }}
                  </option>
                </select>
              </div>
              <div class="col-md-4 mb-3">
                <label class="form-label">Производитель *</label>
                <select class="form-select" v-model="form.manufacturer_id" required>
                  <option value="">Выберите производителя</option>
                  <option v-for="man in manufacturers" :key="man.id" :value="man.id">
                    {{ man.name }}
                  </option>
                </select>
              </div>
              <div class="col-md-4 mb-3">
                <label class="form-label">Цена (₽) *</label>
                <input
                  type="number"
                  class="form-control"
                  v-model.number="form.price"
                  required
                  min="0"
                  step="1"
                  placeholder="99990"
                />
              </div>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal">Отмена</button>
          <button type="button" class="btn btn-primary" @click="saveProduct" :disabled="saving">
            <span v-if="saving" class="spinner-border spinner-border-sm"></span>
            {{ saving ? 'Сохранение...' : 'Сохранить' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { adminApi } from '@/api/admin'
import ErrorAlert from '@/components/common/ErrorAlert.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { computed, onMounted, ref } from 'vue'

const products = ref([])
const categories = ref([])
const manufacturers = ref([])
const loading = ref(false)
const error = ref(null)
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const form = ref({
  id: null,
  name: '',
  category_id: null,
  manufacturer_id: null,
  price: 0
})

const modalStyle = computed(() => ({
  display: showModal.value ? 'block' : 'none',
  backgroundColor: 'rgba(0,0,0,0.5)'
}))

onMounted(() => {
  loadAllData()
})

async function loadAllData() {
  loading.value = true
  error.value = null

  try {
    const [productsRes, categoriesRes, manufacturersRes] = await Promise.all([
      adminApi.getProducts({ limit: 100, offset: 0 }),
      adminApi.getCategories({ limit: 100, offset: 0 }),
      adminApi.getManufacturers({ limit: 100, offset: 0 })
    ])

    products.value = productsRes.product || []
    categories.value = categoriesRes.category || []
    manufacturers.value = manufacturersRes.manufacturer || []
  } catch (err) {
    error.value = formatError(err, 'Не удалось загрузить данные')
  } finally {
    loading.value = false
  }
}

function getCategoryName(id) {
  const cat = categories.value.find(c => c.id === id)
  return cat ? cat.name : `ID: ${id}`
}

function getManufacturerName(id) {
  const man = manufacturers.value.find(m => m.id === id)
  return man ? man.name : `ID: ${id}`
}

function formatPrice(price) {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency: 'RUB',
    minimumFractionDigits: 0
  }).format(price)
}

function openCreateModal() {
  isEditing.value = false
  form.value = { id: null, name: '', category_id: null, manufacturer_id: null, price: 0 }
  showModal.value = true
}

function openEditModal(product) {
  isEditing.value = true
  form.value = {
    id: product.id,
    name: product.name,
    category_id: product.category_id,
    manufacturer_id: product.manufacturer_id,
    price: product.price
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveProduct() {
  if (!form.value.name.trim() || !form.value.category_id || !form.value.manufacturer_id || !form.value.price) {
    error.value = 'Заполните все обязательные поля'
    return
  }

  saving.value = true
  error.value = null

  try {
    const productData = {
      name: form.value.name,
      category_id: form.value.category_id,
      manufacturer_id: form.value.manufacturer_id,
      price: form.value.price
    }

    if (isEditing.value) {
      await adminApi.updateProduct(form.value.id, productData)
    } else {
      await adminApi.createProduct(productData)
    }
    closeModal()
    await loadAllData()
  } catch (err) {
    error.value = formatError(err, 'Не удалось сохранить товар')
  } finally {
    saving.value = false
  }
}

async function confirmDelete(product) {
  if (!confirm(`Удалить товар "${product.name}"?`)) return

  try {
    await adminApi.deleteProduct(product.id)
    await loadAllData()
  } catch (err) {
    error.value = formatError(err, 'Не удалось удалить товар')
  }
}
</script>

<style scoped>
.modal {
  background-color: rgba(0,0,0,0.5);
}
</style>
