<template>
  <Header title="Записи на курсы" />

  <div style="padding: 100px 20px 20px 20px;">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <input v-model="filter" placeholder="Поиск по фамилии" class="form-control" style="min-width: 250px;" />

      <div style="display: flex; flex-wrap: wrap; gap: 10px;">
        <button class="btn btn-primary" @click="toggleSortPrice">
          Цена {{ sortOrder === 'asc' ? '↑' : sortOrder === 'desc' ? '↓' : '' }}
        </button>
        <button class="btn btn-success" @click="goCourseSearch">Поиск по курсу</button>
        <button class="btn btn-info" @click="goAccurateEnrollments">Точные записи на курс</button>
      </div>

      <button class="btn btn-secondary" @click="goBack">Назад</button>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 800px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Фамилия</th>
            <th>Имя</th>
            <th>Курс</th>
            <th>Цена</th>
            <th>Начало</th>
            <th>Окончание</th>
            <th>Группа</th>
            <th>Тип обучения</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="enr in enrollments" :key="enr.id_listener + enr.name_prof_education">
            <td :title="enr.second_name">{{ enr.second_name }}</td>
            <td :title="enr.first_name">{{ enr.first_name }}</td>
            <td :title="enr.name_prof_education">{{ enr.name_prof_education }}</td>
            <td>{{ enr.current_price }} ₽</td>
            <td>{{ formatDate(enr.start_date) }}</td>
            <td>{{ formatDate(enr.end_date) }}</td>
            <td>{{ enr.group || '—' }}</td>
            <td>{{ enr.type_of_retraining || '—' }}</td>
            <td style="white-space: nowrap;">
              <button class="btn btn-success btn-sm" @click="openDetails(enr.id_listener)">Подробнее</button>
            </td>
          </tr>
        </tbody>
      </table>

      <div style="margin-top: 15px; display: flex; justify-content: center; gap: 10px;">
        <button class="btn btn-primary" :disabled="page <= 1" @click="prevPage">Назад</button>
        <span>Страница {{ page }}</span>
        <button class="btn btn-primary" :disabled="!hasMore" @click="nextPage">Вперёд</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'

const router = useRouter()
const enrollments = ref([])
const originalEnrollments = ref([])
const loading = ref(false)
const page = ref(1)
const filter = ref('')
const hasMore = ref(false)
const sortOrder = ref(null)
let debounceTimer = null
const token = localStorage.getItem('access_token')

const fetchEnrollments = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/enrollment/?page=${page.value}&filter=${encodeURIComponent(filter.value)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    enrollments.value = data.data || []
    originalEnrollments.value = [...enrollments.value]
    applySort()
    hasMore.value = enrollments.value.length === 25
  } catch (err) {
    toast.error(err.message || 'Ошибка при загрузке данных')
  } finally {
    loading.value = false
  }
}

const applySort = () => {
  if (sortOrder.value === 'asc') enrollments.value.sort((a,b) => a.current_price - b.current_price)
  else if (sortOrder.value === 'desc') enrollments.value.sort((a,b) => b.current_price - a.current_price)
  else enrollments.value = [...originalEnrollments.value]
}

const toggleSortPrice = () => {
  if (sortOrder.value === null) sortOrder.value = 'asc'
  else if (sortOrder.value === 'asc') sortOrder.value = 'desc'
  else sortOrder.value = null
  applySort()
}

const prevPage = () => { if (page.value > 1) { page.value--; fetchEnrollments() } }
const nextPage = () => { page.value++; fetchEnrollments() }

watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    fetchEnrollments()
  }, 500)
})

const openDetails = (id) => router.push(`/enrollment/details/${id}`)
const goBack = () => router.push('/dashboard/worker')
const goCourseSearch = () => router.push('/enrollment/by-course')
const goAccurateEnrollments = () => router.push('/enrollment/accurate')

const formatDate = d => d?.split(' ')[0] || ''

onMounted(fetchEnrollments)
</script>

<style scoped>
.table-hover tbody tr:hover {
  background-color: #e2f0d9;
  cursor: pointer;
}
th, td {
  vertical-align: middle;
}
</style>
