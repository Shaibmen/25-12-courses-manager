<template>
  <Header title="Добавить тип обучения" />

  <div style="padding: 100px 20px 20px 20px; max-width: 500px; margin: 0 auto;">
    <form @submit.prevent="createType" class="row g-4">

      <div class="col-12">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Название типа обучения</h5>
          <input
            v-model="typeName"
            placeholder="Например: Индивидуально с преподавателем"
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
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import { API_URL_CORE } from '../config'

const router = useRouter()
const typeName = ref('')
const loading = ref(false)
const token = localStorage.getItem('access_token')

const createType = async () => {
  if (!typeName.value.trim()) {
    toast.error('Введите название типа обучения')
    return
  }

  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/educationtype/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify({ type_name: typeName.value.trim() })
    })
    if (!res.ok) throw new Error(`Ошибка (${res.status}) при создании`)

    toast.success('Тип обучения успешно добавлен')
    typeName.value = ''
    setTimeout(() => router.push('/types'), 2000)
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push('/types')
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
.form-control {
  border-radius: 6px;
}
</style>
