<template>
  <Header title="Точные записи на курс" />

  <div style="padding: 100px 20px 20px 20px;">
    <div style="display:flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
      <div style="display:flex; gap:10px; align-items:center;">
        <input v-model="filter" placeholder="Фильтр по курсу" class="form-control" style="min-width: 260px;" />
        <button class="btn btn-secondary" @click="refresh">Обновить</button>
      </div>
      <button class="btn btn-secondary" @click="goBack">Назад</button>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else style="overflow-x:auto;">
      <table class="table table-striped table-hover w-100" style="min-width: 900px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Курс (NameProfEducation)</th>
            <th>Часы (TimeEducation)</th>
            <th>Инд. цена</th>
            <th>Групповая цена</th>
            <th>Кампус цена</th>
            <th>Тип обучения</th>
            <th>Подразделение</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="item in filtered" :key="item.ID_Listener + item.NameProfEducation">
            <td :title="item.NameProfEducation">{{ item.NameProfEducation }}</td>
            <td>{{ item.TimeEducation }}</td>
            <td>{{ formatPrice(item.IndividualPrice) }}</td>
            <td>{{ formatPrice(item.GroupPrice) }}</td>
            <td>{{ formatPrice(item.CampusPrice) }}</td>
            <td>{{ item.EducationType }}</td>
            <td>{{ item.Division }}</td>
          </tr>
        </tbody>
      </table>

      <div v-if="!accurateList.length" class="text-center text-muted" style="padding: 20px;">Данные отсутствуют</div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'

const router = useRouter()
const accurateList = ref([])
const loading = ref(false)
const filter = ref('')
const token = localStorage.getItem('access_token')

const fetchAccurate = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/enrollment/accurate`, {
      method: 'GET',
      headers: {
        Authorization: `Bearer ${token}`,
        'Accept': 'application/json'
      }
    })
    if (!res.ok) {
      if (res.status === 401) throw new Error('401 — неавторизован. Проверьте токен.')
      throw new Error(`Ошибка загрузки (${res.status})`)
    }
    const data = await res.json()
    accurateList.value = Array.isArray(data) ? data : (data.data || [])
  } catch (err) {
    toast.error(err.message || 'Ошибка при загрузке точных записей')
  } finally {
    loading.value = false
  }
}

const filtered = computed(() => {
  if (!filter.value) return accurateList.value
  const f = filter.value.toLowerCase().trim()
  return accurateList.value.filter(i => (i.NameProfEducation || '').toLowerCase().includes(f))
})

const formatPrice = p => p == null || p === 0 ? '—' : `${Number(p).toLocaleString('ru-RU')} ₽`

const refresh = () => fetchAccurate()
const goBack = () => router.push('/enrollments')

onMounted(fetchAccurate)

watch(filter, () => {
}, { debounce: 300 })
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
