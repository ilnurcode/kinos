<template>
  <div class="login-view py-5">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-md-6 col-lg-4">
          <div class="card shadow-sm">
            <div class="card-header bg-primary text-white">
              <h4 class="mb-0">
                Вход
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
                    for="email"
                    class="form-label"
                  >Email</label>
                  <input 
                    id="email" 
                    v-model="form.email" 
                    type="email"
                    class="form-control"
                    :class="{ 'is-invalid': errors.email }"
                    placeholder="example@mail.ru"
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
                    for="password"
                    class="form-label"
                  >Пароль</label>
                  <input 
                    id="password" 
                    v-model="form.password" 
                    type="password"
                    class="form-control"
                    :class="{ 'is-invalid': errors.password }"
                    placeholder="••••••••"
                    required
                  >
                  <div
                    v-if="errors.password"
                    class="invalid-feedback"
                  >
                    {{ errors.password }}
                  </div>
                </div>
                
                <button 
                  type="submit" 
                  class="btn btn-primary w-100"
                  :disabled="loading"
                >
                  <span
                    v-if="loading"
                    class="spinner-border spinner-border-sm me-2"
                  />
                  {{ loading ? 'Вход...' : 'Войти' }}
                </button>
              </form>
            </div>
            <div class="card-footer text-center">
              <small>Нет аккаунта? <router-link to="/register">Зарегистрироваться</router-link></small>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import ErrorAlert from '@/components/common/ErrorAlert.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const form = reactive({
  email: '',
  password: ''
})

const errors = reactive({
  email: null,
  password: null
})

const loading = ref(false)
const error = ref(null)

const validateForm = () => {
  let isValid = true
  errors.email = null
  errors.password = null
  
  if (!form.email) {
    errors.email = 'Введите email'
    isValid = false
  } else if (!/\S+@\S+\.\S+/.test(form.email)) {
    errors.email = 'Введите корректный email'
    isValid = false
  }
  
  if (!form.password) {
    errors.password = 'Введите пароль'
    isValid = false
  } else if (form.password.length < 6) {
    errors.password = 'Пароль должен быть не менее 6 символов'
    isValid = false
  }
  
  return isValid
}

const handleSubmit = async () => {
  if (!validateForm()) return
  
  loading.value = true
  error.value = null
  
  try {
    await authStore.login(form.email, form.password)
    
    const redirect = route.query.redirect || '/'
    router.push(redirect)
  } catch (err) {
    error.value = err.message || 'Ошибка входа'
  } finally {
    loading.value = false
  }
}
</script>
