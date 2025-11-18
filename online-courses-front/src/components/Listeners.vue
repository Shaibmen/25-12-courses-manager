<template>
  <Header title="Слушатели" />

  <div style="padding: 100px 20px 20px 20px; width: 100%;">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <input
        v-model="filter"
        placeholder="Поиск по фамилии"
        style="padding: 8px 12px; min-width: 250px; border-radius: 6px; border: 1px solid #ccc;"
      />
      <div style="display: flex; flex-wrap: wrap; gap: 10px;">
        <button class="btn btn-primary" @click="$router.push('/listeners/create')">Добавить слушателя</button>
        <button class="btn btn-primary" @click="toggleSort">
          Сортировать по дате рождения {{ sortOrder === 'asc' ? '↑' : '↓' }}
        </button>
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 1000px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Фамилия</th>
            <th>Имя</th>
            <th>Отчество</th>
            <th>Email</th>
            <th>Телефон</th>
            <th>Дата рождения</th>
            <th>СНИЛС</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="listener in sortedListeners" :key="listener.id_listener">
            <td :title="listener.second_name">{{ listener.second_name }}</td>
            <td :title="listener.first_name">{{ listener.first_name }}</td>
            <td :title="listener.middle_name || '-'">{{ listener.middle_name || '-' }}</td>
            <td :title="listener.email">{{ listener.email }}</td>
            <td :title="listener.contact_phone">{{ listener.contact_phone }}</td>
            <td>{{ formatDate(listener.date_of_birth) }}</td>
            <td>{{ listener.snils }}</td>
            <td style="white-space: nowrap;">
              <div style="display: flex; gap: 8px;">
                <button class="btn btn-success btn-sm" @click="viewListener(listener.id_listener)">Подробнее</button>
                <button class="btn btn-warning btn-sm" @click="editListener(listener.id_listener)">Изменить</button>
                <button class="btn btn-danger btn-sm" @click="deleteListener(listener.id_listener)">Удалить</button>
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
    title="Удаление слушателя"
    message="Вы точно хотите удалить этого слушателя?"
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
const listeners = ref([])
const loading = ref(false)
const page = ref(1)
const filter = ref('')
const sortOrder = ref('asc')
const hasMore = ref(false)
const token = localStorage.getItem('access_token')
let debounceTimer = null

const confirmModal = ref(null) 

const formatDate = (d) => {
  if (!d) return '-'
  try { return new Date(d).toLocaleDateString() } catch { return '-' }
}

const loadListeners = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/listener/?page=${page.value}&filter=${encodeURIComponent(filter.value)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    listeners.value = Array.isArray(data.data) ? data.data : []
    hasMore.value = data.data?.length === 25
  } catch (err) {
    toast.error(err.message || 'Ошибка при загрузке данных')
  } finally {
    loading.value = false
  }
}

const deleteListener = (id) => {
  if (!confirmModal.value || typeof confirmModal.value.open !== 'function') {
    const confirmed = confirm('Вы точно хотите удалить этого слушателя?')
    if (!confirmed) return
    return performDelete(id)
  }

  confirmModal.value.open(async () => {
    await performDelete(id)
  })
}

const performDelete = async (id) => {
  try {
    const res = await fetch(`${API_URL_CORE}/listener/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) {
      let text = `Ошибка удаления, возможно эта запись используется (${res.status})`
      try {
        const errData = await res.json()
        if (errData?.message) text = errData.message
      } catch {}
      throw new Error(text)
    }

    toast.success('Слушатель удалён')
    loadListeners()
  } catch (err) {
    toast.error(err.message || 'Ошибка при удалении, возможно эта запись используется')
  }
}

const editListener = (uuid) => router.push(`/listeners/edit/${uuid}`)
const viewListener = (uuid) => router.push(`/listeners/${uuid}`)

watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    loadListeners()
  }, 400)
})

const sortedListeners = computed(() => {
  return [...listeners.value].sort((a, b) => {
    const A = new Date(a.date_of_birth)
    const B = new Date(b.date_of_birth)
    return sortOrder.value === 'asc' ? A - B : B - A
  })
})

const toggleSort = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
}

const prevPage = () => {
  if (page.value > 1) {
    page.value--
    loadListeners()
  }
}
const nextPage = () => {
  page.value++
  loadListeners()
}

const goBack = () => router.push('/dashboard/worker')

onMounted(loadListeners)
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
