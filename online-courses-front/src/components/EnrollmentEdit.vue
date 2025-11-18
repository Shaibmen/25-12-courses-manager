<template>
  <Header title="Редактирование обучения" />

  <div style="padding: 100px 20px 20px 20px; max-width: 800px; margin: 0 auto;">
    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else>
      <div class="mb-3">
        <label>Слушатель: <b><span>{{ listener?.second_name }} {{ listener?.first_name }}</span></b> </label> 
      </div>

      <div class="mb-3">
        <label>Программа обучения: <b><span>{{ program.name_prof_education }}</span></b></label>
      </div>

      <div class="d-flex gap-3 align-items-center mb-3">
        <div>
          <label>Дата начала:</label>
          <input type="date" v-model="startDate" class="form-control" />
        </div>
        <div>
          <label>Дата окончания:</label>
          <input type="date" v-model="endDate" class="form-control" />
        </div>
        <div>
          <label>Цена:</label>
          <select v-model="currentPrice" class="form-select">
            <option :value="program.individual_price">Индивидуальное: {{ program.individual_price }}</option>
            <option :value="program.group_price">Групповое: {{ program.group_price }}</option>
            <option :value="program.campus_price">Кампус: {{ program.campus_price }}</option>
          </select>
        </div>
      </div>

      <div class="d-flex gap-3">
        <button class="btn btn-success" @click="save" :disabled="loading">
          {{ loading ? 'Сохранение...' : 'Сохранить' }}
        </button>
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'

const route = useRoute()
const router = useRouter()

const idStudent = route.params.idStudent
const idProgram = route.params.idProgram

const loading = ref(false)
const listener = ref(null)
const program = ref({ individual_price: 0, group_price: 0, campus_price: 0, name_prof_education: '' })
const currentPrice = ref(0)
const startDate = ref('')
const endDate = ref('')

const token = localStorage.getItem('access_token')

const loadData = async () => {
  try {
    loading.value = true
    const [listenerRes, programRes, enrollmentRes] = await Promise.all([
      fetch(`${API_URL_CORE}/listener/details/${idStudent}`, { headers: { Authorization: `Bearer ${token}` } }),
      fetch(`${API_URL_CORE}/programeducation/${idProgram}`, { headers: { Authorization: `Bearer ${token}` } }),
      fetch(`${API_URL_CORE}/enrollment/details/${idStudent}`, { headers: { Authorization: `Bearer ${token}` } })
    ])
    const listenerData = await listenerRes.json()
    const programData = await programRes.json()
    const enrollmentData = await enrollmentRes.json()

    listener.value = listenerData.data.listener
    program.value = programData.data

    const enrollment = enrollmentData.data.find(x => x.id_program_education === idProgram)
    if (!enrollment) throw new Error('Данные обучения не найдены')

    startDate.value = enrollment.start_date.substring(0, 10)
    endDate.value = enrollment.end_date.substring(0, 10)
    currentPrice.value = enrollment.current_price || program.value.individual_price
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const save = async () => {
  loading.value = true
  try {
    const body = {
      id_program: idProgram,
      start_date: startDate.value,
      end_date: endDate.value,
      current_price: Number(currentPrice.value)
    }

    const res = await fetch(`${API_URL_CORE}/enrollment/${idStudent}/${idProgram}`, {
      method: 'PUT',
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body)
    })

    if (!res.ok) throw new Error('Ошибка сохранения')
    toast.success('Обучение успешно обновлено')
    setTimeout(goBack, 500)
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push(`/enrollment/details/${idStudent}`)

onMounted(loadData)
</script>

<style scoped>
.form-control, .form-select {
  min-width: 180px;
}
</style>
