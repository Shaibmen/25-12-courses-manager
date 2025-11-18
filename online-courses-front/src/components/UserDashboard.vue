<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import Header from './Header.vue'

const router = useRouter()
const dashboardData = ref(null)
const error = ref(null)

onMounted(async () => {
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch(`${API_URL_CORE}/dashboard/user`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка загрузки данных')
    dashboardData.value = await res.json()
  } catch (err) {
    error.value = err.message
  }
})
</script>

<template>
    <Header title="Панель работника" />
    <div style="display: flex; flex-wrap: wrap; gap: 20px; padding: 20px; justify-content: center;">
      <button class="btn btn-primary" @click="$router.push('/listeners')">Слушатели</button>
      <button class="btn btn-primary" @click="$router.push('/divisions')">Подразделения обучения</button>
      <button class="btn btn-primary" @click="$router.push('/types')">Типы обучения</button>
      <button class="btn btn-primary" @click="$router.push('/levels')">Уровни обучения</button>
      <button class="btn btn-primary" @click="$router.push('/programs')">Программы обучения</button>
      <button class="btn btn-primary" @click="$router.push('/enrollments')">Запись на курс</button>
    </div>

    <div style="display: flex; flex-wrap: wrap; gap: 20px; padding: 0 20px;">
      <div class="card p-3 shadow-sm" style="flex: 1 1 250px; min-width: 200px;">
        <h5>Всего слушателей</h5>
        <p>{{ dashboardData ? dashboardData.total_listener : '-' }}</p>
      </div>
      <div class="card p-3 shadow-sm" style="flex: 1 1 250px; min-width: 200px;">
        <h5>Всего программ</h5>
        <p>{{ dashboardData ? dashboardData.total_program : '-' }}</p>
      </div>
      <div class="card p-3 shadow-sm" style="flex: 1 1 250px; min-width: 200px;">
        <h5>Активные зачисления</h5>
        <p>{{ dashboardData ? dashboardData.active_enrollments : '-' }}</p>
      </div>
    </div>

    <div style="padding: 20px;">
     <div class="card p-3 w-100 shadow-lg rounded-3 my-3">
  <h5 class="mb-3">Записи на курсы, которые скоро заканчиваются</h5>

  <div v-if="!dashboardData || dashboardData.program_ending_soon.length === 0">
    <p class="text-muted">Нет записей с ближайшим окончанием</p>
  </div>

  <div v-else>
    <div
      v-for="(program, index) in dashboardData.program_ending_soon"
      :key="index"
      class="card p-3 col-4 shadow-sm rounded-3"
    >
      <p class="mb-1 fw-bold">{{ program.name_prof_education }}</p>
      <p class="mb-1">Дата окончания: {{ new Date(program.end_date).toLocaleDateString() }}</p>
      <p class="mb-0 text-muted">Слушателей: {{ program.total_listeners }}</p>
    </div>
  </div>
</div>

    </div>
</template>
