<template>
  <Header title="Подразделения обучения" />

  <div style="padding: 100px 20px 20px 20px;">
    <div style="display: flex; justify-content: space-between; margin-bottom: 20px;">
      <input
        v-model="filter"
        placeholder="Поиск по названию"
        style="padding: 6px 10px; width: 250px; border-radius: 6px; border: 1px solid #ccc;"
      />
      <div>
        <button class="btn btn-primary" @click="$router.push('/divisions/create')">Добавить подразделение</button>
        <button class="btn btn-secondary" @click="goBack" style="margin-left: 10px;">Назад</button>
      </div>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 600px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Название подразделения</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="division in divisions" :key="division.id_divisionsEducation">
            <td>{{ division.divisions }}</td>
            <td style="white-space: nowrap;">
              <button class="btn btn-warning btn-sm" @click="editDivision(division.id_divisionsEducation)">Изменить</button>
              <button class="btn btn-danger btn-sm" @click="deleteDivision(division.id_divisionsEducation)" style="margin-left: 8px;">
                Удалить
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>

  <ConfirmModal
    ref="confirmModal"
    title="Удаление подразделения"
    message="Вы точно хотите удалить это подразделение?"
  />
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_CORE } from '../config'
import Header from './Header.vue'
import ConfirmModal from '../components/ConfirmModal.vue'

const router = useRouter()
const divisions = ref([])
const loading = ref(false)
const filter = ref('')
const token = localStorage.getItem('access_token')
let debounceTimer = null

const confirmModal = ref(null)

const loadDivisions = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/divisions/?filter=${encodeURIComponent(filter.value)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    divisions.value = Array.isArray(data.data) ? data.data : []
  } catch (err) {
    toast.error(err.message || 'Ошибка загрузки данных')
  } finally {
    loading.value = false
  }
}

const deleteDivision = (id) => {
  if (!confirmModal.value || typeof confirmModal.value.open !== 'function') {
    const confirmed = confirm('Удалить это подразделение?')
    if (!confirmed) return
    return performDelete(id)
  }

  confirmModal.value.open(async () => {
    await performDelete(id)
  })
}

const performDelete = async (id) => {
  try {
    const res = await fetch(`${API_URL_CORE}/divisions/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка удаления, возможно эта запись используется (${res.status})`)
    divisions.value = divisions.value.filter(d => d.id_divisionsEducation !== id)
    toast.success('Подразделение удалено')
  } catch (err) {
    toast.error(err.message || 'Ошибка удаления, возможно эта запись используетсяы')
  }
}

const editDivision = (id) => router.push(`/divisions/edit/${id}`)
const goBack = () => router.push('/dashboard/worker')

watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(loadDivisions, 400)
})

onMounted(loadDivisions)
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
