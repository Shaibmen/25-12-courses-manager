<template>
  <Header title="Добавить программу обучения" />

  <div style="padding: 100px 20px 20px 20px; max-width: 600px; width: 500px; margin: 0 auto;">
    <form @submit.prevent="createProgram" class="d-flex flex-column gap-3">

      <div class="card p-3 shadow-sm card-block">
        <h5 class="card-title mb-3">Данные программы</h5>
        <div class="d-flex flex-column gap-2">
          <input v-model="name" placeholder="Название программы" class="form-control" required />
          <input v-model.number="time" type="number" min="1" placeholder="Длительность (мес.)" class="form-control" required />
          <input v-model.number="priceIndividual" type="number" min="0" placeholder="Индивидуальная цена (₽)" class="form-control" required />
          <input v-model.number="priceGroup" type="number" min="0" placeholder="Групповая цена (₽)" class="form-control" required />
          <input v-model.number="priceCampus" type="number" min="0" placeholder="Кампусная цена (₽)" class="form-control" required />
          <select v-model="selectedType" class="form-select" required>
            <option value="">Выберите тип обучения</option>
            <option v-for="t in educationTypes" :key="t.id_educationType" :value="t.id_educationType">{{ t.typeName }}</option>
          </select>
          <select v-model="selectedDivision" class="form-select" required>
            <option value="">Выберите подразделение</option>
            <option v-for="d in divisions" :key="d.id_divisionsEducation" :value="d.id_divisionsEducation">{{ d.divisions }}</option>
          </select>
        </div>
      </div>

      <div class="d-flex justify-content-between mt-3">
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
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_CORE } from '../config'
import Header from './Header.vue'

const router = useRouter()
const name = ref('')
const time = ref()
const priceIndividual = ref()
const priceGroup = ref()
const priceCampus = ref()
const selectedType = ref('')
const selectedDivision = ref('')
const educationTypes = ref([])
const divisions = ref([])
const loading = ref(false)
const token = localStorage.getItem('access_token')

const loadEducationTypes = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/educationtype/`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    const data = await res.json()
    educationTypes.value = data.data || []
  } catch {
    toast.error('Ошибка загрузки типов обучения')
  }
}

const loadDivisions = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/divisions/`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    const data = await res.json()
    divisions.value = data.data || []
  } catch {
    toast.error('Ошибка загрузки подразделений')
  }
}

const createProgram = async () => {
  if (!name.value.trim() || !selectedType.value || !selectedDivision.value) {
    toast.error('Заполните все поля')
    return
  }

  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/programeducation/`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      body: JSON.stringify({
        name_prof_education: name.value,
        time_education: time.value,
        individual_price: priceIndividual.value,
        group_price: priceGroup.value,
        campus_price: priceCampus.value,
        id_educationtype: selectedType.value,
        id_divisionseducation: selectedDivision.value
      })
    })

    if (!res.ok) throw new Error(`Ошибка (${res.status}) при создании`)
    toast.success('Программа успешно добавлена')
    setTimeout(() => router.push('/programs'), 1000)
  } catch (err) {
    toast.error(err.message || 'Ошибка при создании программы')
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push('/programs')

onMounted(async () => {
  await Promise.all([loadEducationTypes(), loadDivisions()])
})
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
</style>
