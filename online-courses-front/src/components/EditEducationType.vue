<template>
  <Header title="Редактировать тип обучения" />

  <div style="padding: 100px 20px 20px 20px; max-width: 500px; margin: 0 auto;">
    <div v-if="loading" class="text-center w-100">Загрузка данных...</div>

    <form v-else @submit.prevent="updateType" class="row g-4">
      <div class="col-12">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Название типа обучения</h5>
          <input
            v-model="typeName"
            placeholder="Введите новое название"
            required
            class="form-control"
          />
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="goBack" :disabled="loading">Назад</button>
        <button type="submit" class="btn btn-success px-4" :disabled="loading">
          {{ loading ? 'Сохранение...' : 'Сохранить' }}
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
const typeName = ref('')
const loading = ref(true)
const token = localStorage.getItem('access_token')

const loadType = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/educationtype/?page=1&filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка (${res.status}) при загрузке типов`)

    const data = await res.json()
    const found = data.data.find(t => t.id_educationType === id)
    if (found) typeName.value = found.typeName
    else throw new Error('Тип обучения не найден')
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const updateType = async () => {
  if (!typeName.value.trim()) {
    toast.error('Введите название типа обучения')
    return
  }

  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/educationtype/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ type_name: typeName.value.trim() })
    })
    if (!res.ok) throw new Error(`Ошибка (${res.status}) при обновлении`)

    toast.success('Тип обучения успешно обновлён')
    setTimeout(() => router.push('/types'), 1000)
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push('/types')

onMounted(loadType)
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
.form-control {
  border-radius: 6px;
}
</style>
