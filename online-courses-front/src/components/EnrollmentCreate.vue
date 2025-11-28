<template>
  <div style="padding: 20px; max-width: 900px; margin: 0 auto;">
    <h2>Запись на курс</h2>

    <div v-if="loading">Загрузка программ...</div>
    <div v-else-if="error" style="color: red;">{{ error }}</div>
    <div v-else>
      <p><strong>Слушатель:</strong> {{ listener.second_name }} {{ listener.first_name }}</p>

      <div style="margin-top: 10px; ">
        <label style="margin-bottom: 10px;">Программа обучения:</label>
        <select v-model="selectedProgramId" @change="onProgramChange" class="form-select">
          <option v-for="p in programs" :key="p.id_program_education" :value="p.id_program_education">
            {{ p.name_prof_education }}
          </option>
        </select>
      </div>

      <div style="margin-top: 20px; display: flex; gap: 5px; justify-content: flex-end;">
        <button class="btn btn-secondary btn-md" :disabled="page <= 1" @click="prevPage">« Предыдущая</button>
        <span style="margin-top: 5px;">Страница {{ page }}</span>
        <button class="btn btn-secondary btn-md" @click="nextPage">Следующая »</button>
      </div>

      <div style="display: flex; align-items: center; gap: 15px; margin-top: 15px; flex-wrap: wrap;">
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
            <option :value="price.individual_price">Индивидуальное: {{ price.individual_price }}</option>
            <option :value="price.group_price">Групповое: {{ price.group_price }}</option>
            <option :value="price.campus_price">Кампус: {{ price.campus_price }}</option>
          </select>
        </div>
        <div>
        <label>Группа:</label>
        <input type="text" v-model="group" class="form-control" placeholder="Например: 0" />
      </div>

      <div>
        <label>Тип обучения:</label>
        <select v-model="typeOfRetraining" class="form-select">
          <option value="Повышение квалификации">Повышение квалификации</option>
          <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
          <option value="Дополительное образование">Дополительное образование</option>
        </select>
      </div>

      </div>

      <div style="margin-top: 20px; display: flex; gap: 10px; flex-wrap: wrap;">
        <button class="btn btn-success btn-md" @click="createEnrollment">Создать запись</button>
        <button class="btn btn-secondary btn-md" @click="goBack">Назад</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import 'vue3-toastify/dist/index.css'
import { API_URL_CORE } from '../config'

const route = useRoute()
const router = useRouter()
const listenerId = ref(route.params.listenerId)

const token = localStorage.getItem('access_token')

const listener = ref({})
const programs = ref([])
const selectedProgramId = ref('')
const price = ref({ individual_price: 0, group_price: 0, campus_price: 0 })
const currentPrice = ref(0)
const startDate = ref('')
const endDate = ref('')
const page = ref(1)
const loading = ref(true)
const error = ref(null)
const group = ref('')
const typeOfRetraining = ref('')

const loadListener = async () => {
  if (!listenerId) {
    toast.error('ID слушателя не найден')
    router.push('/listeners')
    return
  }

  try {
    const res = await fetch(`${API_URL_CORE}/listener/details/${listenerId.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки слушателя (${res.status})`)
    const data = await res.json()
    listener.value = data.data?.listener || {}
  } catch (err) {
    error.value = err.message
    toast.error(err.message)
  }
}

const loadPrograms = async () => {
  try {
    loading.value = true
    const res = await fetch(`${API_URL_CORE}/programeducation/?page=${page.value}&filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки программ (${res.status})`)
    const data = await res.json()
    programs.value = data.data || []
    if (!selectedProgramId.value && programs.value.length) {
      selectedProgramId.value = programs.value[0].id_program_education
      onProgramChange()
    }
  } catch (err) {
    error.value = err.message
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const onProgramChange = () => {
  const sel = programs.value.find(p => p.id_program_education === selectedProgramId.value)
  if (!sel) return
  price.value = { individual_price: sel.individual_price, group_price: sel.group_price, campus_price: sel.campus_price }
  currentPrice.value = price.value.individual_price
}
const createEnrollment = async () => {
  if (!selectedProgramId.value) return toast.warn('Выберите программу')
  if (!listenerId.value) return toast.error('ID слушателя не найден')

  try {
    const body = {
  id_listener: listenerId.value,
  id_program: selectedProgramId.value,
  start_date: startDate.value,
  end_date: endDate.value,
  current_price: Number(currentPrice.value),

  group: group.value,
  type_of_retraining: typeOfRetraining.value,

  is_active: true
}


    const res = await fetch(`${API_URL_CORE}/enrollment/`, {
      method: 'POST',
      headers: { 
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(body)
    })

    if (!res.ok) throw new Error(`Ошибка создания записи (${res.status})`)
    toast.success('Запись успешно создана')

    setTimeout(() => {
      router.push(`/enrollment/details/${listenerId.value}`)
    }, 2000)

  } catch (err) {
    toast.error(err.message)
  }
}


const nextPage = () => { page.value++; loadPrograms() }
const prevPage = () => { if (page.value > 1) { page.value--; loadPrograms() } }
const goBack = () => router.push(`/enrollment/details/${listenerId.value}`)

onMounted(() => {
  loadListener()
  loadPrograms()
})
</script>
