<template>
  <Header title="Уровни обучения" />

  <div style="padding: 100px 20px 20px 20px; width: 100%;">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <input
        v-model="filter"
        placeholder="Поиск по названию уровня"
        style="padding: 8px 12px; min-width: 250px; border-radius: 6px; border: 1px solid #ccc;"
      />
      <div style="display: flex; flex-wrap: wrap; gap: 10px;">
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>
    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 400px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Название уровня образования</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="level in levels" :key="level.id_level_education">
            <td>{{ level.education }}</td>
          </tr>
        </tbody>
      </table>

      <div v-if="levels.length === 0" style="margin-top: 10px; text-align: center; color: #555;">
        Нет данных
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import { API_URL_CORE } from '../config'

const router = useRouter()
const levels = ref([])
const loading = ref(false)
const filter = ref('')
const token = localStorage.getItem('access_token')
let debounceTimer = null

const loadLevels = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/leveleducation/?filter=${encodeURIComponent(filter.value)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()
    levels.value = Array.isArray(data.data) ? data.data : []
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(loadLevels, 400)
})

const goBack = () => router.push('/dashboard/worker')

onMounted(loadLevels)
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
