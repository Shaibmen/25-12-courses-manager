<template>
  <Header title="Типы обучения" />

  <div style="padding: 100px 20px 20px 20px; width: 100%;">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <input
        v-model="filter"
        placeholder="Поиск по названию типа"
        style="padding: 8px 12px; min-width: 250px; border-radius: 6px; border: 1px solid #ccc;"
      />
      <div style="display: flex; flex-wrap: wrap; gap: 10px;">
        <button class="btn btn-primary" @click="$router.push('/types/create')">Добавить тип обучения</button>
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>
    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 600px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Название типа обучения</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="type in types" :key="type.id_educationType">
            <td>{{ type.typeName }}</td>
            <td style="white-space: nowrap;">
              <div style="display: flex; gap: 8px;">
                <button class="btn btn-warning btn-sm" @click="editType(type.id_educationType)">Изменить</button>
                <button class="btn btn-danger btn-sm" @click="openConfirm(type)">Удалить</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <ConfirmModal
      ref="confirmModal"
      title="Удаление типа обучения"
      message="Вы точно хотите удалить этот тип обучения?"
    />
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import ConfirmModal from '../components/ConfirmModal.vue'
import { API_URL_CORE } from '../config'

const router = useRouter()
const types = ref([])
const loading = ref(false)
const filter = ref('')
const token = localStorage.getItem('access_token')
let debounceTimer = null
const confirmModal = ref(null)

const loadTypes = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/educationtype/?filter=${encodeURIComponent(filter.value)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    types.value = Array.isArray(data.data) ? data.data : []
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const openConfirm = (type) => {
  if (!confirmModal.value || typeof confirmModal.value.open !== 'function') {
    if (confirm(`Вы точно хотите удалить тип обучения "${type.typeName}"?`)) {
      deleteType(type.id_educationType)
    }
    return
  }
  confirmModal.value.open(() => deleteType(type.id_educationType))
}

const deleteType = async (id) => {
  try {
    const res = await fetch(`${API_URL_CORE}/educationtype/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка удаления, возможно эта запись используется (${res.status})`)
    types.value = types.value.filter(t => t.id_educationType !== id)
    toast.success('Тип обучения успешно удалён')
  } catch (err) {
    toast.error(err.message)
  }
}

const editType = (id) => router.push(`/types/edit/${id}`)

watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(loadTypes, 400)
})

const goBack = () => router.push('/dashboard/worker')

onMounted(loadTypes)
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
