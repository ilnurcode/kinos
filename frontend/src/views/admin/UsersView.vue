<template>
  <div class="admin-users-view py-4">
    <div class="container">
      <div class="d-flex justify-content-between align-items-center mb-4">
        <h1>Управление пользователями</h1>
        <button
          class="btn btn-primary"
          @click="loadUsers"
        >
          <span
            v-if="loading"
            class="spinner-border spinner-border-sm"
          />
          {{ loading ? 'Загрузка...' : 'Обновить' }}
        </button>
      </div>

      <ErrorAlert
        v-if="error"
        :message="error"
        @dismiss="error = null"
      />

      <div
        v-if="loading"
        class="text-center py-5"
      >
        <LoadingSpinner text="Загрузка пользователей..." />
      </div>

      <div
        v-else-if="users.length === 0"
        class="alert alert-info"
      >
        Пользователи не найдены
      </div>

      <div
        v-else
        class="card"
      >
        <div class="card-body">
          <div class="table-responsive">
            <table class="table table-hover align-middle">
              <thead class="table-light">
                <tr>
                  <th>ID</th>
                  <th>Имя</th>
                  <th>Email</th>
                  <th>Телефон</th>
                  <th>Роль</th>
                  <th>Действия</th>
                </tr>
              </thead>
              <tbody>
                <tr
                  v-for="user in users"
                  :key="user.id"
                >
                  <td>{{ user.id }}</td>
                  <td>
                    <strong>{{ user.username }}</strong>
                    <span
                      v-if="user.id === currentUserId"
                      class="badge bg-info ms-1"
                    >Вы</span>
                  </td>
                  <td>{{ user.email }}</td>
                  <td>{{ user.phone || '—' }}</td>
                  <td>
                    <span
                      v-if="user.id === currentUserId"
                      class="badge"
                      :class="getRoleBadge(user.role)"
                    >
                      {{ getRoleLabel(user.role) }}
                    </span>
                    <select
                      v-else
                      :value="user.role"
                      class="form-select form-select-sm"
                      :class="getRoleBadgeClass(user.role)"
                      @change="changeRole(user.id, user.username, $event.target.value)"
                    >
                      <option value="user">
                        Пользователь
                      </option>
                      <option value="admin">
                        Администратор
                      </option>
                    </select>
                  </td>
                  <td>
                    <button
                      v-if="user.id !== currentUserId"
                      class="btn btn-sm btn-outline-danger"
                      :disabled="deletingUserId === user.id"
                      @click="confirmDelete(user)"
                    >
                      <span
                        v-if="deletingUserId === user.id"
                        class="spinner-border spinner-border-sm"
                      />
                      {{ deletingUserId === user.id ? '...' : 'Удалить' }}
                    </button>
                    <span
                      v-else
                      class="text-muted"
                    >—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { adminApi } from '@/api/admin'
import ErrorAlert from '@/components/common/ErrorAlert.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { useAuthStore } from '@/stores/auth'
import { formatError } from '@/utils/errorHandler'
import { computed, onMounted, ref } from 'vue'

const authStore = useAuthStore()
const currentUserId = computed(() => authStore.user?.id)

const users = ref([])
const loading = ref(false)
const error = ref(null)
const deletingUserId = ref(null)

onMounted(() => {
  loadUsers()
})

async function loadUsers() {
  loading.value = true
  error.value = null

  try {
    const response = await adminApi.getUsers({ limit: 100, offset: 0 })
    users.value = response.users || []
  } catch (err) {
    error.value = formatError(err, 'Не удалось загрузить пользователей')
  } finally {
    loading.value = false
  }
}

async function changeRole(userId, username, newRole) {
  if (!confirm(`Изменить роль пользователя ${username} на "${getRoleLabel(newRole)}"?`)) {
    loadUsers() // Вернуть старое значение
    return
  }

  try {
    await adminApi.updateUserRole(userId, newRole)
    // Обновить локально
    const user = users.value.find(u => u.id === userId)
    if (user) {
      user.role = newRole
    }
  } catch (err) {
    error.value = formatError(err, 'Не удалось изменить роль')
    loadUsers() // Перезагрузить список
  }
}

async function confirmDelete(user) {
  if (!confirm(`Вы уверены что хотите удалить пользователя ${user.username}?`)) {
    return
  }

  deletingUserId.value = user.id
  try {
    await adminApi.deleteUser(user.id)
    users.value = users.value.filter(u => u.id !== user.id)
  } catch (err) {
    error.value = formatError(err, 'Не удалось удалить пользователя')
  } finally {
    deletingUserId.value = null
  }
}

function getRoleBadge(role) {
  return role === 'admin' ? 'bg-warning text-dark' : 'bg-secondary'
}

function getRoleBadgeClass(role) {
  return role === 'admin' ? 'border-warning' : 'border-secondary'
}

function getRoleLabel(role) {
  return role === 'admin' ? 'Администратор' : 'Пользователь'
}
</script>

<style scoped>
.admin-users-view {
  min-height: calc(100vh - 200px);
}

.form-select {
  min-width: 150px;
  cursor: pointer;
}
</style>
