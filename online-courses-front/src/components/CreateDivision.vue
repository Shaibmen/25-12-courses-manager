<template>
  <Header title="Добавить подразделение" />

  <div style="padding: 100px 20px 20px 20px; max-width: 600px; margin: 0 auto;">
    <form @submit.prevent="createDivision" class="row g-4">

      <div class="col-12">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Название подразделения</h5>
          <input
            v-model="divisionName"
            type="text"
            placeholder="Введите название"
            required
            class="form-control"
          />
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="goBack" :disabled="loading">Назад</button>
        <button type="submit" class="btn btn-success px-4" :disabled="loading">
          {{ loading ? 'Сохранение...' : 'Создать' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import { API_URL_CORE } from '../config'

const router = useRouter()
const divisionName = ref('')
const loading = ref(false)
const token = localStorage.getItem('access_token')

const createDivision = async () => {
  if (!divisionName.value.trim()) {
    toast.error('Введите название подразделения')
    return
  }

  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/divisions/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ divisions: divisionName.value.trim() })
    })
    if (!res.ok) throw new Error(`Ошибка (${res.status})`)

    toast.success('Подразделение успешно добавлено')
    divisionName.value = ''
    setTimeout(() => router.push('/divisions'), 1000)
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push('/divisions')
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
</style>
