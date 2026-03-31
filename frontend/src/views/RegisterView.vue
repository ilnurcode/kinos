<template>
  <div class="register-view py-5">
    <div class="container">
      <div class="row justify-content-center">
        <div class="col-md-6 col-lg-5">
          <div class="card shadow-sm">
            <div class="card-header bg-success text-white">
              <h4 class="mb-0">Регистрация</h4>
            </div>
            <div class="card-body">
              <ErrorAlert 
                v-if="error" 
                :message="error" 
                @dismiss="error = null" 
              />
              
              <form @submit.prevent="handleSubmit">
                <div class="mb-3">
                  <label for="username" class="form-label">Имя пользователя</label>
                  <input 
                    type="text" 
                    class="form-control" 
                    id="username"
                    v-model="form.username"
                    :class="{ 'is-invalid': errors.username }"
                    placeholder="Ivan"
                    required
                  />
                  <div class="invalid-feedback" v-if="errors.username">
                    {{ errors.username }}
                  </div>
                </div>
                
                <div class="mb-3">
                  <label for="email" class="form-label">Email</label>
                  <input 
                    type="email" 
                    class="form-control" 
                    id="email"
                    v-model="form.email"
                    :class="{ 'is-invalid': errors.email }"
                    placeholder="example@mail.ru"
                    required
                  />
                  <div class="invalid-feedback" v-if="errors.email">
                    {{ errors.email }}
                  </div>
                </div>
                
                <div class="mb-3">
                  <label for="password" class="form-label">Пароль</label>
                  <input 
                    type="password" 
                    class="form-control" 
                    id="password"
                    v-model="form.password"
                    :class="{ 'is-invalid': errors.password }"
                    placeholder="••••••••"
                    minlength="8"
                    required
                  />
                  <div class="invalid-feedback" v-if="errors.password">
                    {{ errors.password }}
                  </div>
                </div>
                
                <div class="mb-3">
                  <label for="phone" class="form-label">Телефон</label>
                  <div class="input-group">
                    <span class="input-group-text">+7</span>
                    <input 
                      type="tel" 
                      class="form-control" 
                      id="phone"
                      v-model="form.phone"
                      :class="{ 'is-invalid': errors.phone }"
                      placeholder="999 123 45 67"
                      maxlength="14"
                      required
                    />
                  </div>
                  <div class="invalid-feedback" v-if="errors.phone">
                    {{ errors.phone }}
                  </div>
                  <small class="text-muted">Введите 10 цифр номера (например: 999 123 45 67)</small>
                </div>
                
                <button 
                  type="submit" 
                  class="btn btn-success w-100"
                  :disabled="loading"
                >
                  <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
                  {{ loading ? 'Регистрация...' : 'Зарегистрироваться' }}
                </button>
              </form>
            </div>
            <div class="card-footer text-center">
              <small>Уже есть аккаунт? <router-link to="/login">Войти</router-link></small>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import ErrorAlert from '@/components/common/ErrorAlert.vue'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  username: '',
  email: '',
  password: '',
  phone: ''
})

const errors = reactive({
  username: null,
  email: null,
  password: null,
  phone: null
})

const loading = ref(false)
const error = ref(null)

// Маска телефона
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
  errors.password = null
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
  
  if (!form.password) {
    errors.password = 'Введите пароль'
    isValid = false
  } else if (form.password.length < 8) {
    errors.password = 'Пароль должен быть не менее 8 символов'
    isValid = false
  }
  
  if (!form.phone) {
    errors.phone = 'Введите телефон'
    isValid = false
  } else if (form.phone.replace(/\s/g, '').length !== 10) {
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
    // Форматируем телефон в E.164
    const phoneDigits = form.phone.replace(/\s/g, '')
    const phoneE164 = '+7' + phoneDigits
    
    await authStore.register({
      username: form.username,
      email: form.email,
      password: form.password,
      phone: phoneE164
    })
    
    // После успешной регистрации - логин
    await authStore.login(form.email, form.password)
    router.push('/')
  } catch (err) {
    error.value = err.message || 'Ошибка регистрации'
  } finally {
    loading.value = false
  }
}

// Обработчик ввода телефона
const handlePhoneInput = (event) => {
  form.phone = formatPhone(event.target.value)
}
</script>
