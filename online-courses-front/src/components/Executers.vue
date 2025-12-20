<template>
  <Header title="Исполнители" />

  <div style="padding: 100px 20px 20px 20px; width: 100%;">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <input
        v-model="filter"
        placeholder="Поиск по фамилии"
        style="padding: 8px 12px; min-width: 250px; border-radius: 6px; border: 1px solid #ccc;"
      />
      <div style="display: flex; flex-wrap: wrap; gap: 10px;">
        <button class="btn btn-primary" @click="$router.push('/executers/create')">Добавить исполнителя</button>
        <button class="btn btn-primary" @click="toggleSort">
          Сортировать по фамилии {{ sortOrder === 'asc' ? '↑' : '↓' }}
        </button>
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 800px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Фамилия</th>
            <th>Имя</th>
            <th>Отчество</th>
            <th>Должность</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="executer in sortedExecuters" :key="executer.id_executor">
            <td :title="executer.second_name">{{ executer.second_name }}</td>
            <td :title="executer.first_name">{{ executer.first_name }}</td>
            <td :title="executer.middle_name || '-'">{{ executer.middle_name || '-' }}</td>
            <td :title="executer.status">{{ executer.status }}</td>
            <td style="white-space: nowrap;">
              <div style=" gap: 8px;">
                <button class="btn btn-danger btn-sm" @click="deleteExecuter(executer.id_executor)">Удалить</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div style="margin-top: 20px; display: flex; justify-content: center; gap: 10px;">
        <button class="btn btn-primary" :disabled="page <= 1" @click="prevPage">Назад</button>
        <span>Страница {{ page }}</span>
        <button class="btn btn-primary" :disabled="!hasMore" @click="nextPage">Вперёд</button>
      </div>
    </div>
  </div>

  <ConfirmModal
    ref="confirmModal"
    title="Удаление исполнителя"
    message="Вы точно хотите удалить этого исполнителя?"
  />
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_CORE } from '../config'
import Header from './Header.vue'
import ConfirmModal from '../components/ConfirmModal.vue'

const router = useRouter()
const executers = ref([])
const loading = ref(false)
const page = ref(1)
const filter = ref('')
const sortOrder = ref('asc')
const hasMore = ref(false)
const token = localStorage.getItem('access_token')
let debounceTimer = null

const confirmModal = ref(null)

const loadExecuters = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/executer/?page=${page.value}&filter=${encodeURIComponent(filter.value)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    executers.value = Array.isArray(data.data) ? data.data : []
    hasMore.value = data.data?.length === 25
  } catch (err) {
    toast.error(err.message || 'Ошибка при загрузке данных')
  } finally {
    loading.value = false
  }
}

const deleteExecuter = (id) => {
  confirmModal.value?.open(async () => {
    await performDelete(id)
  })
}

const performDelete = async (id) => {
  try {
    const res = await fetch(`${API_URL_CORE}/executer/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) {
      let text = `Ошибка удаления. Возможно, запись используется (${res.status})`
      try {
        const errData = await res.json()
        if (errData?.message) text = errData.message
      } catch {}
      throw new Error(text)
    }

    toast.success('Исполнитель удалён')
    loadExecuters()
  } catch (err) {
    toast.error(err.message || 'Не удалось удалить исполнителя')
  }
}


watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    loadExecuters()
  }, 400)
})

const sortedExecuters = computed(() => {
  return [...executers.value].sort((a, b) => {
    const aName = (a.second_name || '').toLowerCase()
    const bName = (b.second_name || '').toLowerCase()
    return sortOrder.value === 'asc' ? aName.localeCompare(bName) : bName.localeCompare(aName)
  })
})

const toggleSort = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  toast.info(`Сортировка: ${sortOrder.value === 'asc' ? 'По возрастанию' : 'По убыванию'}`)
}

const prevPage = () => {
  if (page.value > 1) {
    page.value--
    loadExecuters()
  }
}

const nextPage = () => {
  page.value++
  loadExecuters()
}

const goBack = () => router.push('/dashboard/worker')

onMounted(loadExecuters)
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