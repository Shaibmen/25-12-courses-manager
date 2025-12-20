<template>
  <Header title="Программы обучения" />

  <div style="padding: 100px 20px 20px 20px; width: 100%;">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <input
        v-model="filter"
        placeholder="Поиск по названию программы"
        style="padding: 8px 12px; min-width: 250px; border-radius: 6px; border: 1px solid #ccc;"
      />
      <div style="display: flex; flex-wrap: wrap; gap: 10px;">
        <button class="btn btn-primary" @click="$router.push('/programs/create')">Добавить программу</button>
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>

    <div style="display: flex; flex-wrap: wrap; gap: 10px; margin-bottom: 15px; justify-content: flex-start;">
      <button class="btn btn-primary" @click="toggleSort('individual')">
        Индивидуальная цена <span v-if="sortField === 'individual'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
      </button>
      <button class="btn btn-primary" @click="toggleSort('group')">
        Групповая цена <span v-if="sortField === 'group'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
      </button>
      <button class="btn btn-primary" @click="toggleSort('campus')">
        Кампусная цена <span v-if="sortField === 'campus'">{{ sortOrder === 'asc' ? '↑' : '↓' }}</span>
      </button>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>
    <div v-else-if="error" class="text-danger mb-3">{{ error }}</div>
    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 800px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Название программы</th>
            <th>Длительность (часы)</th>
            <th>Индивидуально (₽)</th>
            <th>Групповое (₽)</th>
            <th>На кампусе (₽)</th>
            <th>Тип обучения</th>
            <th>Подразделение</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="program in sortedPrograms" :key="program.id_program_education">
            <td>{{ program.name_prof_education }}</td>
            <td>{{ program.time_education }}</td>
            <td>{{ program.individual_price }}</td>
            <td>{{ program.group_price }}</td>
            <td>{{ program.campus_price }}</td>
            <td>{{ getTypeName(program.id_education_type) }}</td>
            <td>{{ getDivisionName(program.id_divisions_education) }}</td>
            <td style="white-space: nowrap;">
              <div style="display: flex; gap: 8px;">
                <button class="btn btn-warning btn-sm" @click="editProgram(program.id_program_education)">Изменить</button>
                <button class="btn btn-danger btn-sm" @click="confirmModalRef.open(() => deleteProgram(program.id_program_education))">Удалить</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div style="margin-top: 15px; display: flex; justify-content: center; gap: 10px;">
        <button class="btn btn-secondary" :disabled="page <= 1" @click="prevPage">Назад</button>
        <span>Страница {{ page }}</span>
        <button class="btn btn-secondary" :disabled="!hasMore" @click="nextPage">Вперёд</button>
      </div>

      <div v-if="sortedPrograms.length === 0" style="margin-top: 10px; text-align: center; color: #555;">
        Нет данных
      </div>
    </div>

    <ConfirmModal
      ref="confirmModalRef"
      title="Удаление программы"
      message="Вы точно хотите удалить эту программу?"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import ConfirmModal from '../components/ConfirmModal.vue'
import { API_URL_CORE } from '../config'

const router = useRouter()
const programs = ref([])
const educationTypes = ref([])
const divisions = ref([])
const loading = ref(false)
const error = ref(null)
const page = ref(1)
const filter = ref('')
const hasMore = ref(false)
const token = localStorage.getItem('access_token')

const sortField = ref(null)
const sortOrder = ref(null)
let debounceTimer = null

const confirmModalRef = ref(null)

const loadPrograms = async () => {
  loading.value = true
  error.value = null
  try {
    const res = await fetch(`${API_URL_CORE}/programeducation/?page=${page.value}&filter=${encodeURIComponent(filter.value)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    programs.value = Array.isArray(data.data) ? data.data : []
    hasMore.value = data.data?.length === 25
  } catch (err) {
    toast.error(err.message)
    error.value = err.message
  } finally {
    loading.value = false
  }
}

const loadEducationTypes = async () => {
  const res = await fetch(`${API_URL_CORE}/educationtype/?page=1&filter=`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  const data = await res.json()
  educationTypes.value = data.data || []
}

const loadDivisions = async () => {
  const res = await fetch(`${API_URL_CORE}/divisions/?page=1&filter=`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  const data = await res.json()
  divisions.value = data.data || []
}

const getTypeName = (id) => {
  const t = educationTypes.value.find(e => e.id_educationType === id)
  return t ? t.typeName : '—'
}

const getDivisionName = (id) => {
  const d = divisions.value.find(e => e.id_divisionsEducation === id)
  return d ? d.divisions : '—'
}

const editProgram = (id) => router.push(`/programs/edit/${id}`)

const deleteProgram = async (id) => {
  try {
    const res = await fetch(`${API_URL_CORE}/programeducation/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка удаления: возможно на курс записан слушатель`)
    programs.value = programs.value.filter(p => p.id_program_education !== id)
    toast.success('Программа успешно удалена')
  } catch (err) {
    toast.error(err.message)
  }
}

const toggleSort = (field) => {
  if (sortField.value !== field) {
    sortField.value = field
    sortOrder.value = 'asc'
  } else {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : sortOrder.value === 'desc' ? null : 'asc'
    if (!sortOrder.value) sortField.value = null
  }
}

const sortedPrograms = computed(() => {
  let list = [...programs.value]
  if (filter.value.trim()) {
    list = list.filter(p => p.name_prof_education.toLowerCase().includes(filter.value.toLowerCase()))
  }
  if (sortField.value && sortOrder.value) {
    const keyMap = { individual: 'individual_price', group: 'group_price', campus: 'campus_price' }
    const key = keyMap[sortField.value]
    list.sort((a, b) => (sortOrder.value === 'asc' ? a[key] - b[key] : b[key] - a[key]))
  }
  return list
})

watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    loadPrograms()
  }, 500)
})

const prevPage = () => { if (page.value > 1) { page.value--; loadPrograms() } }
const nextPage = () => { page.value++; loadPrograms() }
const goBack = () => router.push('/dashboard/worker')

onMounted(async () => {
  await Promise.all([loadPrograms(), loadEducationTypes(), loadDivisions()])
})
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
