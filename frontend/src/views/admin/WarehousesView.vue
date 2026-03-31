<template>
  <div class="admin-warehouses py-4">
    <div class="container">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление складами</h1>
        <button class="btn btn-primary" @click="openCreateModal">
          <i class="bi bi-plus-lg"></i> Добавить склад
        </button>
      </div>

      <ErrorAlert v-if="error" :message="error" @dismiss="error = null" />

      <div v-if="loading" class="text-center py-5">
        <LoadingSpinner text="Загрузка складов..." />
      </div>

      <div v-else-if="warehouses.length === 0" class="alert alert-info">
        Склады не найдены
      </div>

      <div v-else class="card">
        <div class="card-body">
          <table class="table table-hover">
            <thead class="table-light">
              <tr>
                <th>ID</th>
                <th>Название</th>
                <th>Город</th>
                <th>Адрес</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="wh in warehouses" :key="wh.id">
                <td>{{ wh.id }}</td>
                <td><strong>{{ wh.name }}</strong></td>
                <td>{{ wh.city || '—' }}</td>
                <td>{{ wh.address || '—' }}</td>
                <td>
                  <button class="btn btn-sm btn-outline-primary me-2" @click="openEditModal(wh)">
                    ✏️
                  </button>
                  <button class="btn btn-sm btn-outline-danger" @click="confirmDelete(wh)">
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
          <h5 class="modal-title">{{ isEditing ? 'Редактировать склад' : 'Новый склад' }}</h5>
          <button type="button" class="btn-close" @click="closeModal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveWarehouse">
            <div class="mb-3">
              <label class="form-label">Название склада *</label>
              <input
                type="text"
                class="form-control"
                v-model="form.name"
                required
                placeholder="Основной"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Город *</label>
              <input
                type="text"
                class="form-control"
                v-model="form.city"
                required
                placeholder="Москва"
              />
            </div>
            <div class="mb-3">
              <label class="form-label">Улица *</label>
              <input
                type="text"
                class="form-control"
                v-model="form.street"
                required
                placeholder="ул. Складская"
              />
            </div>
            <div class="row">
              <div class="col-md-6 mb-3">
                <label class="form-label">Дом</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.building"
                  placeholder="1"
                />
              </div>
              <div class="col-md-6 mb-3">
                <label class="form-label">Строение</label>
                <input
                  type="text"
                  class="form-control"
                  v-model="form.building2"
                  placeholder="1А"
                />
              </div>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal">Отмена</button>
          <button type="button" class="btn btn-primary" @click="saveWarehouse" :disabled="saving">
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

const warehouses = ref([])
const loading = ref(false)
const error = ref(null)
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const form = ref({
  id: null,
  name: '',
  city: '',
  street: '',
  building: '',
  building2: ''
})

const modalStyle = computed(() => ({
  display: showModal.value ? 'block' : 'none',
  backgroundColor: 'rgba(0,0,0,0.5)'
}))

onMounted(() => {
  loadWarehouses()
})

async function loadWarehouses() {
  loading.value = true
  try {
    const response = await adminApi.getWarehouses({ limit: 100, offset: 0 })
    warehouses.value = response.warehouses || []
  } catch (err) {
    error.value = formatError(err, 'Не удалось загрузить склады')
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  form.value = { id: null, name: '', city: '', street: '', building: '', building2: '' }
  showModal.value = true
}

function openEditModal(wh) {
  isEditing.value = true
  form.value = {
    id: wh.id,
    name: wh.name,
    city: wh.city || '',
    street: wh.street || '',
    building: wh.building || '',
    building2: wh.building2 || ''
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveWarehouse() {
  if (!form.value.name.trim()) {
    error.value = 'Введите название склада'
    return
  }
  saving.value = true
  try {
    const data = {
      name: form.value.name,
      city: form.value.city,
      street: form.value.street,
      building: form.value.building,
      building2: form.value.building2
    }
    if (isEditing.value) {
      await adminApi.updateWarehouse(form.value.id, data)
    } else {
      await adminApi.createWarehouse(data)
    }
    closeModal()
    await loadWarehouses()
  } catch (err) {
    error.value = formatError(err, 'Не удалось сохранить склад')
  } finally {
    saving.value = false
  }
}

async function confirmDelete(wh) {
  if (!confirm(`Удалить склад "${wh.name}"?`)) return
  try {
    await adminApi.deleteWarehouse(wh.id)
    await loadWarehouses()
  } catch (err) {
    error.value = formatError(err, 'Не удалось удалить склад')
  }
}
</script>

<style scoped>
.modal { background-color: rgba(0,0,0,0.5); }
</style>
