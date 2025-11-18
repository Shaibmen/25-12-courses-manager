<template>
  <Header title="Редактировать подразделение" />

  <div style="padding: 100px 20px 20px 20px; max-width: 600px; margin: 0 auto;">
    <form @submit.prevent="updateDivision" class="row g-4">

      <div v-if="loading" class="col-12 text-center">Загрузка данных...</div>

      <div v-else class="col-12">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Название подразделения</h5>
          <input
            v-model="divisionName"
            type="text"
            placeholder="Введите название подразделения"
            required
            class="form-control"
          />
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between" v-if="!loading">
        <button type="button" class="btn btn-secondary" @click="goBack" :disabled="saving">Назад</button>
        <button type="submit" class="btn btn-success px-4" :disabled="saving">
          {{ saving ? 'Сохранение...' : 'Сохранить изменения' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import { API_URL_CORE } from '../config'

const route = useRoute()
const router = useRouter()
const id = route.params.id
const token = localStorage.getItem('access_token')

const divisionName = ref('')
const loading = ref(true)
const saving = ref(false)

const loadDivision = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/divisions/?page=1&filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)

    const data = await res.json()
    const found = data.data.find(d => d.id_divisionsEducation === id)
    if (!found) throw new Error('Подразделение не найдено')

    divisionName.value = found.divisions
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const updateDivision = async () => {
  if (!divisionName.value.trim()) {
    toast.error('Введите название подразделения')
    return
  }

  saving.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/divisions/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ divisions: divisionName.value.trim() })
    })
    if (!res.ok) throw new Error(`Ошибка обновления (${res.status})`)

    toast.success('Подразделение успешно обновлено')
    setTimeout(() => router.push('/divisions'), 1000)
  } catch (err) {
    toast.error(err.message)
  } finally {
    saving.value = false
  }
}

const goBack = () => router.push('/divisions')

onMounted(loadDivision)
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
.form-control {
  border-radius: 6px;
}
</style>
