<script setup>
import { ref, onMounted, defineProps } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL_AUTH } from '../config'

const props = defineProps({
  title: {
    type: String,
    default: 'Панель работника'
  }
})

const router = useRouter()
const user = ref({ username: '', role: '' })
const rolesMap = {
  admin: 'Администратор',
  worker: 'Работник',
  accountant: 'Бухгалтер'
}

const logout = async () => {
  const token = localStorage.getItem('access_token')
  try {
    await fetch(`${API_URL_AUTH}/protected/logout`, {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` }
    })
  } catch (_) {}
  localStorage.clear()
  router.push('/login')
}

onMounted(() => {
  const userName = localStorage.getItem('username')
  const userRole = localStorage.getItem('role')
  user.value = { username: userName, role: userRole }
})
</script>

<template>
  <div class="header">
    <div class="d-flex justify-content-between align-items-center p-3 bg-white shadow-sm">
      <h2 class="mb-0">{{ props.title }}</h2>

      <div class="d-flex align-items-center gap-3">
        <div class="text-end">
          <p class="mb-0">Сотрудник: <b>{{ user.username }}</b></p>
          <p class="mb-0">Роль: <b>{{ rolesMap[user.role] || 'Неизвестно' }}</b></p>
        </div>
        <button class="btn btn-danger btn-sm" @click="logout">Выйти</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  z-index: 1000;
  border-radius: 0;
}

body, #app {
  padding-top: 80px;
}
</style>
