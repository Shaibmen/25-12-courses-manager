<template>
  <Header title="Редактирование программы обучения" />

  <div style="padding: 100px 20px 20px 20px; max-width: 800px; width: 500px; margin: 0 auto;">
    <div class="card p-3 shadow-sm card-block">
      <form @submit.prevent="updateProgram" class="d-flex flex-column gap-3">
        <div>
          <label>Название программы:</label>
          <input v-model="name" required placeholder="Название программы" class="form-control" />
        </div>

        <div>
          <label>Длительность (часы):</label>
          <input v-model.number="time" type="number" min="1" required class="form-control" />
        </div>

        <div>
          <label>Индивид. цена (₽):</label>
          <input v-model.number="priceIndividual" type="number" min="0" required class="form-control" />
        </div>

        <div>
          <label>Групповая цена (₽):</label>
          <input v-model.number="priceGroup" type="number" min="0" required class="form-control" />
        </div>

        <div>
          <label>Кампусная цена (₽):</label>
          <input v-model.number="priceCampus" type="number" min="0" required class="form-control" />
        </div>

        <div>
          <label>Тип обучения:</label>
          <select v-model="selectedType" required class="form-select">
            <option value="">Выберите тип обучения</option>
            <option v-for="t in educationTypes" :key="t.id_educationType" :value="t.id_educationType">
              {{ t.typeName }}
            </option>
          </select>
        </div>

        <div>
          <label>Подразделение:</label>
          <select v-model="selectedDivision" required class="form-select">
            <option value="">Выберите подразделение</option>
            <option v-for="d in divisions" :key="d.id_divisionsEducation" :value="d.id_divisionsEducation">
              {{ d.divisions }}
            </option>
          </select>
        </div>

        <div class="d-flex justify-content-between mt-2">
          <button type="button" class="btn btn-secondary" @click="goBack" :disabled="loading">Назад</button>
          <button type="submit" class="btn btn-success" :disabled="loading">
            {{ loading ? 'Сохранение...' : 'Обновить' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { useRoute, useRouter } from "vue-router"
import { API_URL_CORE } from "../config"
import { toast } from "vue3-toastify"
import Header from './Header.vue'

const router = useRouter()
const route = useRoute()
const id = route.params.id

const name = ref("")
const time = ref(1)
const priceIndividual = ref(0)
const priceGroup = ref(0)
const priceCampus = ref(0)
const selectedType = ref("")
const selectedDivision = ref("")

const educationTypes = ref([])
const divisions = ref([])

const loading = ref(false)
const token = localStorage.getItem("access_token")

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

const loadProgram = async () => {
  loading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/programeducation/${id}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    const resp = await res.json()
    const program = resp?.data

    name.value = program.name_prof_education
    time.value = program.time_education
    priceIndividual.value = program.individual_price
    priceGroup.value = program.group_price
    priceCampus.value = program.campus_price
    selectedType.value = program.id_education_type
    selectedDivision.value = program.id_divisions_education
  } catch {
    toast.error('Ошибка загрузки программы')
  } finally {
    loading.value = false
  }
}

const updateProgram = async () => {
  loading.value = true
  try {
    const body = {
      name_prof_education: name.value,
      time_education: time.value,
      individual_price: priceIndividual.value,
      group_price: priceGroup.value,
      campus_price: priceCampus.value,
      ID_EducationType: selectedType.value,
      ID_DivisionsEducation: selectedDivision.value
    }

    const res = await fetch(`${API_URL_CORE}/programeducation/${id}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {})
      },
      body: JSON.stringify(body)
    })

    if (!res.ok) throw new Error()
    toast.success('Программа успешно обновлена!')
    setTimeout(() => router.push("/programs"), 800)
  } catch {
    toast.error('Ошибка обновления программы')
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push("/programs")

onMounted(async () => {
  await Promise.all([loadEducationTypes(), loadDivisions(), loadProgram()])
})
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
</style>
