<template>
  <div class="profile-edit-view py-5">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-md-8 col-lg-6">
          <div class="card shadow-sm">
            <div class="card-header bg-success text-white">
              <h4 class="mb-0">
                Редактирование профиля
              </h4>
            </div>
            <div class="card-body">
              <ErrorAlert
                v-if="error"
                :message="error"
                @dismiss="error = null"
              />

              <form @submit.prevent="handleSubmit">
                <div class="mb-3">
                  <label
                    for="username"
                    class="form-label"
                  >Имя пользователя</label>
                  <input
                    id="username"
                    v-model="form.username"
                    type="text"
                    class="form-control"
                    :class="{ 'is-invalid': errors.username }"
                    required
                  >
                  <div
                    v-if="errors.username"
                    class="invalid-feedback"
                  >
                    {{ errors.username }}
                  </div>
                </div>

                <div class="mb-3">
                  <label
                    for="email"
                    class="form-label"
                  >Email</label>
                  <input
                    id="email"
                    v-model="form.email"
                    type="email"
                    class="form-control"
                    :class="{ 'is-invalid': errors.email }"
                    required
                  >
                  <div
                    v-if="errors.email"
                    class="invalid-feedback"
                  >
                    {{ errors.email }}
                  </div>
                </div>

                <div class="mb-3">
                  <label
                    for="phone"
                    class="form-label"
                  >Телефон</label>
                  <div class="input-group">
                    <span class="input-group-text">+7</span>
                    <input
                      id="phone"
                      v-model="form.phone"
                      type="tel"
                      class="form-control"
                      :class="{ 'is-invalid': errors.phone }"
                      placeholder="999 123 45 67"
                      maxlength="14"
                    >
                  </div>
                  <div
                    v-if="errors.phone"
                    class="invalid-feedback"
                  >
                    {{ errors.phone }}
                  </div>
                  <small class="text-muted">Введите 10 цифр номера</small>
                </div>

                <div class="d-grid gap-2">
                  <button
                    type="submit"
                    class="btn btn-success"
                    :disabled="loading"
                  >
                    <span
                      v-if="loading"
                      class="spinner-border spinner-border-sm me-2"
                    />
                    {{ loading ? 'Сохранение...' : 'Сохранить' }}
                  </button>
                  <router-link
                    to="/profile"
                    class="btn btn-outline-secondary"
                  >
                    Отмена
                  </router-link>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '@/api/auth'
import { useAuth } from '@/composables/useAuth'
import ErrorAlert from '@/components/common/ErrorAlert.vue'

const router = useRouter()
const { loadProfile } = useAuth()

const form = reactive({
  username: '',
  email: '',
  phone: ''
})

const errors = reactive({
  username: null,
  email: null,
  phone: null
})

const loading = ref(false)
const error = ref(null)

const formatPhone = (value) => {
  let digits = value.replace(/\D/g, '')

  if (digits.length > 10) {
    digits = digits.substring(0, 10)
  }

  if (digits.length > 0) {
    let formatted = ''
    if (digits.length >= 3) {
      formatted += digits.substring(0, 3) + ' '
    } else {
      formatted += digits
    }
    if (digits.length >= 6) {
      formatted += digits.substring(3, 6) + ' '
    } else if (digits.length > 3) {
      formatted += digits.substring(3)
    }
    if (digits.length >= 8) {
      formatted += digits.substring(6, 8) + ' '
    } else if (digits.length > 6) {
      formatted += digits.substring(6)
    }
    if (digits.length >= 10) {
      formatted += digits.substring(8, 10)
    } else if (digits.length > 8) {
      formatted += digits.substring(8)
    }
    return formatted
  }

  return ''
}

const validateForm = () => {
  let isValid = true
  errors.username = null
  errors.email = null
  errors.phone = null

  if (!form.username) {
    errors.username = 'Введите имя пользователя'
    isValid = false
  } else if (form.username.length < 3) {
    errors.username = 'Имя должно быть не менее 3 символов'
    isValid = false
  }

  if (!form.email) {
    errors.email = 'Введите email'
    isValid = false
  } else if (!/\S+@\S+\.\S+/.test(form.email)) {
    errors.email = 'Введите корректный email'
    isValid = false
  }

  if (form.phone && form.phone.replace(/\s/g, '').length !== 10) {
    errors.phone = 'Введите 10 цифр номера'
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  if (!validateForm()) return

  loading.value = true
  error.value = null

  try {
    const phoneDigits = form.phone.replace(/\s/g, '')
    const phoneE164 = phoneDigits ? '+7' + phoneDigits : ''

    await authApi.updateProfile({
      username: form.username,
      email: form.email,
      phone: phoneE164
    })

    await loadProfile()
    router.push('/profile')
  } catch (err) {
    error.value = err.message || 'Ошибка сохранения'
  } finally {
    loading.value = false
  }
}

onMounted(async () => {
  const profile = await loadProfile()
  if (profile) {
    form.username = profile.username || ''
    form.email = profile.email || ''

    if (profile.phone) {
      const phoneDigits = profile.phone.replace(/\D/g, '').substring(1)
      form.phone = formatPhone(phoneDigits)
    }
  }
})
</script>
