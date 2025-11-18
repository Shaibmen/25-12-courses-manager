<template>
   <div class="d-flex justify-content-center align-items-center" style="background-color: #f8f9fa;">
    <div class="card shadow p-4" style="width: 350px; border: none;">
      <h2 class="text-center text-black mb-4">Авторизация</h2>
      <form @submit.prevent="login" class="d-flex flex-column align-items-center">
        <div class="mb-3 w-100">
          <input 
            v-model="username" 
            type="text" 
            id="username" 
            class="form-control mx-auto" 
            placeholder="Введите логин"
          />
        </div>
        <div class="mb-4 w-100">
          <input 
            v-model="password" 
            type="password" 
            id="password" 
            class="form-control mx-auto" 
            placeholder="Введите пароль"
          />
        </div>
        <button type="submit" class="btn btn-success w-100">Войти</button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_AUTH } from '../config'

const router = useRouter()
const username = ref('')
const password = ref('')

const login = async () => {
  try {
    const res = await fetch(`${API_URL_AUTH}/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: username.value,
        password: password.value
      })
    })

    if (!res.ok) throw new Error('Неверный логин или пароль')

    const data = await res.json()

    if (!data.token || !data.token.user) {
      throw new Error('Некорректный ответ сервера')
    }

    localStorage.setItem('access_token', data.token.access_token)
    localStorage.setItem('refresh_token', data.token.refresh_token)
    localStorage.setItem('access_token_expire_at', data.token.access_token_expire_at)
    localStorage.setItem('refresh_token_expire_at', data.token.refresh_token_expire_at)
    localStorage.setItem('username', data.token.user.username)
    localStorage.setItem('role', data.token.user.role)

    const role = data.token.user.role

    switch (role) {
      case 'admin':
        toast.success('Вход выполнен успешно!')
        router.push('/dashboard/admin')
        break
      case 'worker':
        toast.success('Вход выполнен успешно!')
        router.push('/dashboard/worker')
        break
      case 'accountant':
        toast.success('Вход выполнен успешно!')
        router.push('/dashboard/accountant')
        break
      default:
        toast.error('Неизвестная роль')
        break
    }
  } catch (e) {
    toast.error(e.message)
  }
}
</script>
