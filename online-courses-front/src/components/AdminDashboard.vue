<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL_AUTH, API_URL_CORE } from '../config'
import Header from './Header.vue'
import { toast } from 'vue3-toastify'

const router = useRouter()
const dashboardData = ref(null)
const error = ref(null)

const username = ref('')
const password = ref('')
const selectedRole = ref('')
const roles = ref([]) 

const backups = ref([])
const selectedBackup = ref('')

onMounted(async () => {
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch(`${API_URL_CORE}/dashboard/admin`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка загрузки данных')
    const data = await res.json()
    dashboardData.value = data
    roles.value = data.role
  } catch (err) {
    toast.error(err.message)
  }

  await loadBackups()
})

const registerUser = async () => {
  if (!username.value || !password.value || !selectedRole.value) {
    return toast.error('Заполните все поля')
  }
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch(`${API_URL_AUTH}/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({
        username: username.value,
        password: password.value,
        role: selectedRole.value
      })
    })
    const data = await res.json().catch(() => ({})) 
    if (!res.ok) throw new Error(data.message || 'Ошибка регистрации')
    toast.success('Пользователь успешно зарегистрирован')
    username.value = ''
    password.value = ''
    selectedRole.value = ''
  } catch (err) {
    toast.error(err.message || 'Ошибка регистрации')
  }
}

const loadBackups = async () => {
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch(`${API_URL_CORE}/backup/all`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка загрузки списка бэкапов')
    const data = await res.json()
    backups.value = data.data || []
  } catch (err) {
    toast.error(err.message)
  }
}

const createBackup = async () => {
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch(`${API_URL_CORE}/backup/`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка при создании бэкапа')
    toast.success('Бэкап успешно создан')
    await loadBackups()
  } catch (err) {
    toast.error(err.message)
  }
}

const restoreBackup = async () => {
  if (!selectedBackup.value) return toast.error('Выберите файл для восстановления')
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch(`${API_URL_CORE}/backup/restore?file=${selectedBackup.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка при восстановлении бэкапа')
    toast.success(`Бэкап ${selectedBackup.value} успешно восстановлен`)
  } catch (err) {
    toast.error(err.message)
  }
}
</script>

<template>
  <Header title="Панель администратора" />

  <div class="container py-4" style="margin-top: 50px;">

    <div class="row g-3 mb-4">
      <div class="col-md-4" v-for="(stat, index) in [
        { label: 'База данных', value: dashboardData?.admin_stat.database_name },
        { label: 'Размер базы', value: dashboardData?.admin_stat.total_size },
        { label: 'Активные соединения', value: dashboardData?.admin_stat.active_connection },
        { label: 'Коммит транзакций', value: dashboardData?.admin_stat.committed_tx },
        { label: 'Откат транзакций', value: dashboardData?.admin_stat.rolledback_tx },
        { label: 'Чтение блоков диска', value: dashboardData?.admin_stat.disk_block_read },
        { label: 'Буферные попадания', value: dashboardData?.admin_stat.buffer_hits },
      ]" :key="index">
        <div class="card shadow-sm p-3 h-100">
          <h6>{{ stat.label }}</h6>
          <p class="mb-0">{{ stat.value ?? '-' }}</p>
        </div>
      </div>
    </div>

    <div class="card shadow-lg p-3 mb-4">
      <h5 class="mb-3">Процессы PostgreSQL</h5>
      <div v-if="!dashboardData || !dashboardData.pg_stat?.length">
        <p class="text-muted">Нет активных процессов</p>
      </div>
      <div class="row g-3" v-else>
        <div class="col-md-4" v-for="(proc, index) in dashboardData.pg_stat" :key="index">
          <div class="card p-3 shadow-sm h-100">
            <p class="mb-1">PID: {{ proc.pid }}</p>
            <p class="mb-1">Состояние: {{ proc.state }}</p>
            <p class="mb-1">Клиент: {{ proc.client_addr }}:{{ proc.client_port }}</p>
            <p class="mb-0 text-muted">Начало запроса: {{ new Date(proc.query_start).toLocaleString() }}</p>
          </div>
        </div>
      </div>
    </div>
    <div class="card shadow-lg p-3 mb-4">
      <h5 class="mb-3">Регистрация нового пользователя</h5>
      <div class="row g-3">
        <div class="col-md-4">
          <input v-model="username" class="form-control" type="text" placeholder="Имя пользователя" />
        </div>
        <div class="col-md-4">
          <input v-model="password" class="form-control" type="password" placeholder="Пароль" />
        </div>
        <div class="col-md-4">
          <select v-model="selectedRole" class="form-select">
            <option value="">Выберите роль</option>
            <option v-for="role in roles" :key="role.id" :value="role.id">{{ role.role }}</option>
          </select>
        </div>
      </div>
      <div class="mt-3">
        <button class="btn btn-success" @click="registerUser">Зарегистрировать</button>
      </div>
    </div>

    <div class="card shadow-lg p-3 mb-4">
      <h5 class="mb-3">Бэкапы базы данных</h5>
      <div class="d-flex align-items-center gap-3">
        <button class="btn btn-primary" @click="createBackup">Создать бэкап</button>
        <select v-model="selectedBackup" class="form-select w-auto">
          <option value="">Выберите бэкап для восстановления</option>
          <option v-for="file in backups" :key="file" :value="file">{{ file }}</option>
        </select>
        <button class="btn btn-warning" @click="restoreBackup">Восстановить бэкап</button>
      </div>
    </div>

  </div>
</template>

<style scoped>
.container {
  max-width: 1200px;
  margin: 0 auto;
}
</style>
