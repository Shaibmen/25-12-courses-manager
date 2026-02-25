<template>
  <Header title="Создать исполнителя" />

  <div style="padding: 100px 20px 20px 20px;">
    <form @submit.prevent="createExecuter" class="row g-4">
      <div class="col-md-6">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Основные данные</h5>
          <div class="d-flex flex-column gap-2">
            <input
              v-model="executer.second_name"
              @input="onlyLetters('second_name')"
              class="form-control"
              placeholder="Фамилия"
              required
            />
            <input
              v-model="executer.first_name"
              @input="onlyLetters('first_name')"
              class="form-control"
              placeholder="Имя"
              required
            />
            <input
              v-model="executer.middle_name"
              @input="onlyLetters('middle_name')"
              class="form-control"
              placeholder="Отчество"
            />
            <input
              v-model="executer.status"
              class="form-control"
              placeholder="Должность / Статус"
              required
            />
             <input
              v-model="executer.doverenost"
              class="form-control"
              placeholder="№ Доверенности"
            />
          </div>
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="goBack" :disabled="loading">
          Назад
        </button>
        <button type="submit" class="btn btn-success px-4" :disabled="loading">
          {{ loading ? 'Создание...' : 'Создать' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_CORE } from '../config'
import Header from './Header.vue'

const router = useRouter()
const loading = ref(false)

const executer = ref({
  first_name: '',
  second_name: '',
  middle_name: '',
  status: '',
  doverenost:''
})

const onlyLetters = (field) => {
  executer.value[field] = (executer.value[field] || '').replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}

const goBack = () => router.push('/executers')

const createExecuter = async () => {
  const token = localStorage.getItem('access_token')
  if (!token) {
    toast.error('Не авторизован')
    return
  }

  const payload = { ...executer.value }
  if (!payload.middle_name?.trim()) delete payload.middle_name

  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/executer/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(payload)
    })

    if (!res.ok) {
      let msg = 'Ошибка при создании'
      try {
        const err = await res.json()
        if (err?.message) msg = err.message
      } catch {}
      throw new Error(msg)
    }

    toast.success('Исполнитель успешно создан')
    setTimeout(() => router.push('/executers'), 2000)
  } catch (err) {
    toast.error(err.message || 'Не удалось создать исполнителя')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
</style>
