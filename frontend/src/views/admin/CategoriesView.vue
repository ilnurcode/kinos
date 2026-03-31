<template>
  <div class="admin-categories py-4">
    <div class="container">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление категориями</h1>
        <button class="btn btn-primary" @click="openCreateModal">
          <i class="bi bi-plus-lg"></i> Добавить категорию
        </button>
      </div>

      <ErrorAlert v-if="error" :message="error" @dismiss="error = null" />

      <div v-if="loading" class="text-center py-5">
        <LoadingSpinner text="Загрузка категорий..." />
      </div>

      <div v-else-if="categories.length === 0" class="alert alert-info">
        Категории не найдены
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
              <tr v-for="category in categories" :key="category.id">
                <td>{{ category.id }}</td>
                <td><strong>{{ category.name }}</strong></td>
                <td>
                  <button class="btn btn-sm btn-outline-primary me-2" @click="openEditModal(category)">
                    ✏️ Редактировать
                  </button>
                  <button class="btn btn-sm btn-outline-danger" @click="confirmDelete(category)">
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
          <h5 class="modal-title">{{ isEditing ? 'Редактировать категорию' : 'Новая категория' }}</h5>
          <button type="button" class="btn-close" @click="closeModal"></button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveCategory">
            <div class="mb-3">
              <label class="form-label">Название категории *</label>
              <input
                type="text"
                class="form-control"
                v-model="form.name"
                required
                placeholder="Например: Электроника"
              />
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="closeModal">Отмена</button>
          <button type="button" class="btn btn-primary" @click="saveCategory" :disabled="saving">
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
import { formatError } from '@/utils/errorHandler'
import { computed, onMounted, ref } from 'vue'

const categories = ref([])
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
  loadCategories()
})

async function loadCategories() {
  loading.value = true
  error.value = null

  try {
    const response = await adminApi.getCategories({ limit: 100, offset: 0 })
    categories.value = response.category || []
  } catch (err) {
    error.value = formatError(err, 'Не удалось загрузить категории')
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  form.value = { id: null, name: '' }
  showModal.value = true
}

function openEditModal(category) {
  isEditing.value = true
  form.value = { id: category.id, name: category.name }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function saveCategory() {
  if (!form.value.name.trim()) {
    error.value = 'Введите название категории'
    return
  }

  saving.value = true
  error.value = null

  try {
    if (isEditing.value) {
      await adminApi.updateCategory(form.value.id, form.value.name)
    } else {
      await adminApi.createCategory(form.value.name)
    }
    closeModal()
    await loadCategories()
  } catch (err) {
    error.value = formatError(err, 'Не удалось сохранить категорию')
  } finally {
    saving.value = false
  }
}

async function confirmDelete(category) {
  if (!confirm(`Удалить категорию "${category.name}"?`)) return

  try {
    await adminApi.deleteCategory(category.id)
    await loadCategories()
  } catch (err) {
    error.value = formatError(err, 'Не удалось удалить категорию')
  }
}
</script>

<style scoped>
.modal {
  background-color: rgba(0,0,0,0.5);
}
</style>
