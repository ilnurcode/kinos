<template>
  <div class="admin-manufacturers py-4">
    <div class="container">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление производителями</h1>
        <button class="btn btn-primary" @click="openCreateModal">
          <i class="bi bi-plus-lg"></i> Добавить производителя
        </button>
      </div>

      <ErrorAlert v-if="error" :message="error" @dismiss="error = null" />

      <div v-if="loading" class="text-center py-5">
        <LoadingSpinner text="Загрузка производителей..." />
      </div>

      <div v-else-if="manufacturers.length === 0" class="alert alert-info">
        Производители не найдены
      </div>

      <div v-else class="card">
        <div class="card-body">
          <table class="table table-hover">
            <thead class="table-light">
              <tr>
                <th>ID</th>
                <th>Название</th>
                <th>Действия</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="manufacturer in manufacturers" :key="manufacturer.id">
                <td>{{ manufacturer.id }}</td>
                <td><strong>{{ manufacturer.name }}</strong></td>
                <td>
                  <button class="btn btn-sm btn-outline-primary me-2" @click="openEditModal(manufacturer)">
                    ✏️ Редактировать
                  </button>
                  <button class="btn btn-sm btn-outline-danger" @click="confirmDelete(manufacturer)">
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
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ isEditing ? 'Редактировать производителя' : 'Новый производитель' }}</h5>
          <button type="button" class="btn-close" @click="closeModal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveManufacturer">
            <div class="mb-3">
              <label class="form-label">Название производителя *</label>
              <input
                type="text"
                class="form-control"
                v-model="form.name"
                required
                placeholder="Например: Samsung"
              />
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal">Отмена</button>
          <button type="button" class="btn btn-primary" @click="saveManufacturer" :disabled="saving">
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

const manufacturers = ref([])
const loading = ref(false)
const error = ref(null)
const showModal = ref(false)
const isEditing = ref(false)
const saving = ref(false)

const form = ref({
  id: null,
  name: ''
})

const modalStyle = computed(() => ({
  display: showModal.value ? 'block' : 'none',
  backgroundColor: 'rgba(0,0,0,0.5)'
}))

onMounted(() => {
  loadManufacturers()
})

async function loadManufacturers() {
  loading.value = true
  error.value = null

  try {
    const response = await adminApi.getManufacturers({ limit: 100, offset: 0 })
    manufacturers.value = response.manufacturer || []
  } catch (err) {
    error.value = formatError(err, 'Не удалось загрузить производителей')
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  form.value = { id: null, name: '' }
  showModal.value = true
}

function openEditModal(manufacturer) {
  isEditing.value = true
  form.value = { id: manufacturer.id, name: manufacturer.name }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveManufacturer() {
  if (!form.value.name.trim()) {
    error.value = 'Введите название производителя'
    return
  }

  saving.value = true
  error.value = null

  try {
    if (isEditing.value) {
      await adminApi.updateManufacturer(form.value.id, form.value.name)
    } else {
      await adminApi.createManufacturer(form.value.name)
    }
    closeModal()
    await loadManufacturers()
  } catch (err) {
    error.value = formatError(err, 'Не удалось сохранить производителя')
  } finally {
    saving.value = false
  }
}

async function confirmDelete(manufacturer) {
  if (!confirm(`Удалить производителя "${manufacturer.name}"?`)) return

  try {
    await adminApi.deleteManufacturer(manufacturer.id)
    await loadManufacturers()
  } catch (err) {
    error.value = formatError(err, 'Не удалось удалить производителя')
  }
}
</script>

<style scoped>
.modal {
  background-color: rgba(0,0,0,0.5);
}
</style>
