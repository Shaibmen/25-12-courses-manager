<template>
  <Header title="Поиск слушателей по курсу" />

  <div style="padding: 100px 20px 20px 20px; width: 100%;">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <div style="display: flex; align-items: center; gap: 10px; flex-wrap: wrap;">
        <label>Курс:</label>
        <select v-model="selectedCourseId" @change="onCourseChange" class="form-select">
          <option v-for="c in courses" :key="c.id_program_education" :value="c.id_program_education">
            {{ c.name_prof_education }}
          </option>
        </select>

        <button class="btn btn-secondary btn-md" :disabled="coursePage <= 1" @click="prevCoursePage">« Предыдущая</button>
        <span>Страница {{ coursePage }}</span>
        <button class="btn btn-secondary btn-md" @click="nextCoursePage">Следующая »</button>
      </div>

      <div>
        <button class="btn btn-secondary btn-md" style="margin-bottom: 15px;" @click="goBack">Назад</button>
      </div>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else>
      <div v-if="enrollments.length === 0" style="margin-top: 20px; font-weight: bold;">
        На этом курсе слушателей нет
      </div>

      <div v-else style="overflow-x: auto;">
        <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 900px;">
          <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
            <tr>
              <th>Фамилия</th>
              <th>Имя</th>
              <th>Курс</th>
              <th>
                <button class="btn btn-link p-0 " @click="toggleSortPrice">
                  Цена 
                  <span v-if="sortOrder === 'asc'">↑</span>
                  <span v-else-if="sortOrder === 'desc'">↓</span>
                </button>
              </th>
              <th>Начало</th>
              <th>Окончание</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="enr in enrollments" :key="enr.id_listener + enr.name_prof_education">
              <td>{{ enr.second_name }}</td>
              <td>{{ enr.first_name }}</td>
              <td>{{ enr.name_prof_education }}</td>
              <td>{{ enr.current_price }} ₽</td>
              <td>{{ formatDate(enr.start_date) }}</td>
              <td>{{ formatDate(enr.end_date) }}</td>
              <td style="white-space: nowrap;">
                <button class="btn btn-success btn-sm" @click="openDetails(enr.id_listener)">Подробнее</button>
              </td>
            </tr>
          </tbody>
        </table>

        <div style="margin-top: 15px; display: flex; justify-content: center; gap: 10px; flex-wrap: wrap;">
          <button class="btn btn-primary btn-md" :disabled="page <= 1" @click="prevPage">Назад</button>
          <span>Страница {{ page }}</span>
          <button class="btn btn-primary btn-md" :disabled="!hasMore" @click="nextPage">Вперёд</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_CORE } from '../config'
import Header from './Header.vue'

const router = useRouter()
const token = localStorage.getItem('access_token')

const enrollments = ref([])
const originalEnrollments = ref([])
const loading = ref(false)
const page = ref(1)
const hasMore = ref(false)
const sortOrder = ref(null)

const courses = ref([])
const selectedCourseId = ref('')
const coursePage = ref(1)

const loadCourses = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/programeducation/?page=${coursePage.value}&filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки курсов (${res.status})`)
    const data = await res.json()
    courses.value = data.data || []
    if (!selectedCourseId.value && courses.value.length) {
      selectedCourseId.value = courses.value[0].id_program_education
      loadEnrollments()
    }
  } catch (err) {
    toast.error(err.message)
  }
}

const loadEnrollments = async () => {
  if (!selectedCourseId.value) return
  try {
    loading.value = true
    const res = await fetch(`${API_URL_CORE}/enrollment/${selectedCourseId.value}?page=${page.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) {
      if (res.status === 404 || res.status === 204) {
        enrollments.value = []
        originalEnrollments.value = []
        hasMore.value = false
        return
      }
      throw new Error(`Ошибка загрузки слушателей (${res.status})`)
    }
    const data = await res.json()
    enrollments.value = data.data || []
    originalEnrollments.value = [...enrollments.value]
    applySort()
    hasMore.value = enrollments.value.length === 25
  } catch (err) {
    toast.error(err.message)
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

const prevPage = () => { if (page.value > 1) { page.value--; loadEnrollments() } }
const nextPage = () => { page.value++; loadEnrollments() }

const prevCoursePage = () => { if (coursePage.value > 1) { coursePage.value--; loadCourses() } }
const nextCoursePage = () => { coursePage.value++; loadCourses() }

const onCourseChange = () => { page.value = 1; loadEnrollments() }

const openDetails = id => router.push(`/enrollment/details/${id}`)
const goBack = () => router.push('/enrollments')

const formatDate = d => d?.split(' ')[0] || ''

onMounted(() => {
  loadCourses()
})
</script>
