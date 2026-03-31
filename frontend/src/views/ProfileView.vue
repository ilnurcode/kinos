<template>
  <div class="profile-view py-5">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-md-8 col-lg-6">
          <div class="card shadow-sm">
            <div class="card-header bg-primary text-white">
              <h4 class="mb-0">Профиль пользователя</h4>
            </div>
            <div class="card-body">
              <LoadingSpinner v-if="loading" text="Загрузка профиля..." />

              <div v-else-if="error" class="alert alert-danger">
                {{ error }}
              </div>

              <div v-else>
                <div class="mb-3">
                  <label class="form-label text-muted">Имя пользователя</label>
                  <p class="h5">{{ user?.username }}</p>
                </div>

                <div class="mb-3">
                  <label class="form-label text-muted">Email</label>
                  <p class="h5">{{ user?.email }}</p>
                </div>

                <div class="mb-3">
                  <label class="form-label text-muted">Телефон</label>
                  <p class="h5">{{ user?.phone || 'Не указан' }}</p>
                </div>

                <div class="mb-3">
                  <label class="form-label text-muted">Роль</label>
                  <p>
                    <span :class="badgeClass">
                      {{ user?.role === 'admin' ? 'Администратор' : 'Пользователь' }}
                    </span>
                  </p>
                </div>

                <div class="d-grid gap-2">
                  <router-link to="/profile/edit" class="btn btn-primary">
                    Редактировать профиль
                  </router-link>
                  <router-link to="/catalog" class="btn btn-outline-secondary">
                    В каталог
                  </router-link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useAuth } from '@/composables/useAuth'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'

const { user, loading, error, loadProfile } = useAuth()

const badgeClass = computed(() => {
  return user.value?.role === 'admin'
    ? 'badge bg-danger'
    : 'badge bg-secondary'
})

onMounted(async () => {
  await loadProfile()
})
</script>
