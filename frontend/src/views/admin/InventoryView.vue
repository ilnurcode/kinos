<template>
  <div class="admin-inventory py-4">
    <div class="container">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление запасами</h1>
        <button class="btn btn-primary" @click="openCreateModal">
          <i class="bi bi-plus-lg"></i> Добавить запас
        </button>
      </div>

      <ErrorAlert v-if="error" :message="error" @dismiss="error = null" />

      <div v-if="loading" class="text-center py-5">
        <LoadingSpinner text="Загрузка запасов..." />
      </div>

      <div v-else-if="inventory.length === 0" class="alert alert-info">
        Запасы не найдены
      </div>

      <div v-else class="card">
        <div class="card-body">
          <table class="table table-hover">
            <thead class="table-light">
              <tr>
                <th>ID</th>
                <th>Товар</th>
                <th>Всего</th>
                <th>Зарезервировано</th>
                <th>Доступно</th>
                <th>Склад</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="item in inventory" :key="item.id">
                <td>{{ item.id }}</td>
                <td>{{ getProductName(item.product_id) }}</td>
                <td>{{ item.quantity }}</td>
                <td>{{ item.reserved_quantity }}</td>
                <td>
                  <span :class="getQuantityBadge(item.available_quantity)">
                    {{ item.available_quantity }}
                  </span>
                </td>
                <td>{{ item.warehouse_location }}</td>
                <td>
                  <button class="btn btn-sm btn-outline-primary me-2" @click="openEditModal(item)">
                    ✏️
                  </button>
                  <button class="btn btn-sm btn-outline-danger" @click="confirmDelete(item)">
                    🗑️
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
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Редактировать запас' : 'Добавить запас' }}</h5>
          <button type="button" class="btn-close" @click="closeModal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveInventory">
            <div class="mb-3">
              <label class="form-label">Товар *</label>
              <select class="form-select" v-model="form.product_id" required>
                <option value="">Выберите товар</option>
                <option v-for="product in products" :key="product.id" :value="product.id">
                  {{ product.name }}
                </option>
              </select>
            </div>
            <div class="mb-3">
              <label class="form-label">Количество *</label>
              <input
                type="number"
                class="form-control"
                v-model.number="form.quantity"
                required
                min="0"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Склад *</label>
              <select class="form-select" v-model="form.warehouse_location" required>
                <option value="">Выберите склад</option>
                <option v-for="wh in warehouses" :key="wh.id" :value="wh.name">
                  {{ wh.name }}{{ wh.city ? ` (${wh.city})` : '' }}
                </option>
              </select>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal">Отмена</button>
          <button type="button" class="btn btn-primary" @click="saveInventory" :disabled="saving">
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

const inventory = ref([])
const products = ref([])
const warehouses = ref([])
const loading = ref(false)
const error = ref(null)
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const form = ref({
  id: null,
  product_id: null,
  quantity: 0,
  warehouse_location: ''
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
  try {
    const [invRes, prodRes, whRes] = await Promise.all([
      adminApi.getInventory({ limit: 100, offset: 0 }),
      adminApi.getProducts({ limit: 100, offset: 0 }),
      adminApi.getWarehouses({ limit: 100, offset: 0 })
    ])
    inventory.value = invRes.inventory || []
    products.value = prodRes.product || []
    warehouses.value = whRes.warehouses || []
  } catch (err) {
    error.value = formatError(err, 'Не удалось загрузить данные')
  } finally {
    loading.value = false
  }
}

function getProductName(id) {
  const p = products.value.find(p => p.id === id)
  return p ? p.name : `ID: ${id}`
}

function getQuantityBadge(qty) {
  if (qty > 10) return 'badge bg-success'
  if (qty > 0) return 'badge bg-warning'
  return 'badge bg-danger'
}

function openCreateModal() {
  isEditing.value = false
  form.value = { id: null, product_id: null, quantity: 0, warehouse_location: '' }
  showModal.value = true
}

function openEditModal(item) {
  isEditing.value = true
  form.value = {
    id: item.id,
    product_id: item.product_id,
    quantity: item.quantity,
    warehouse_location: item.warehouse_location
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveInventory() {
  saving.value = true
  try {
    const data = {
      product_id: form.value.product_id,
      quantity: form.value.quantity,
      warehouse_location: form.value.warehouse_location
    }
    if (isEditing.value) {
      await adminApi.updateInventory(form.value.id, data)
    } else {
      await adminApi.createInventory(data)
    }
    closeModal()
    await loadAllData()
  } catch (err) {
    error.value = formatError(err, 'Не удалось сохранить запас')
  } finally {
    saving.value = false
  }
}

async function confirmDelete(item) {
  if (!confirm('Удалить этот запас?')) return
  try {
    await adminApi.deleteInventory(item.id)
    await loadAllData()
  } catch (err) {
    error.value = formatError(err, 'Не удалось удалить запас')
  }
}
</script>

<style scoped>
.modal { background-color: rgba(0,0,0,0.5); }
</style>
