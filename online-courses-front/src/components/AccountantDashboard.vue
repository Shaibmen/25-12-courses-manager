<template>
  <Header title="Панель бухгалтера" />

  <div style="display: flex; flex-wrap: wrap; gap: 20px; padding: 0 20px;">
    <div class="card p-3 shadow-sm" style="flex: 1 1 250px; min-width: 200px;">
      <h5>Всего слушателей</h5>
      <p>{{ dashboardData ? dashboardData.total_listener : '-' }}</p>
    </div>
    <div class="card p-3 shadow-sm" style="flex: 1 1 250px; min-width: 200px;">
      <h5>Всего программ</h5>
      <p>{{ dashboardData ? dashboardData.total_program : '-' }}</p>
    </div>
    <div class="card p-3 shadow-sm" style="flex: 1 1 250px; min-width: 200px;">
      <h5>Активные зачисления</h5>
      <p>{{ dashboardData ? dashboardData.active_enrollments : '-' }}</p>
    </div>
  </div>

  <div style="padding: 20px;">
    <div class="card p-3 w-100 shadow-lg rounded-3 my-3">
      <h5 class="mb-3">Записи на курсы, которые скоро заканчиваются</h5>
      <div v-if="!dashboardData || dashboardData.program_ending_soon.length === 0">
        <p class="text-muted">Нет записей с ближайшим окончанием</p>
      </div>
      <div v-else>
        <div
          v-for="(program, index) in dashboardData.program_ending_soon"
          :key="index"
          class="card p-3 col-4 shadow-sm rounded-3"
        >
          <p class="mb-1 fw-bold">{{ program.name_prof_education }}</p>
          <p class="mb-1">Дата окончания: {{ new Date(program.end_date).toLocaleDateString() }}</p>
          <p class="mb-0 text-muted">Слушателей: {{ program.total_listeners }}</p>
        </div>
      </div>
    </div>

    <div class="card p-3 shadow-lg rounded-3 my-3">
      <h5 class="mb-3">Отчёт по периодам</h5>
      <div class="d-flex gap-3">
        <input type="date" v-model="startDatePeriod" class="form-control" />
        <input type="date" v-model="endDatePeriod" class="form-control" />
        <button class="btn btn-success" @click="downloadReport('period')">Скачать</button>
      </div>
    </div>

    <div class="card p-3 shadow-lg rounded-3 my-3">
      <h5 class="mb-3">Отчёт по затратам</h5>
      <div class="d-flex gap-3">
        <input type="date" v-model="startDateExpensive" class="form-control" />
        <input type="date" v-model="endDateExpensive" class="form-control" />
        <button class="btn btn-warning" @click="downloadReport('expensive')">Скачать</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from "vue"
import { API_URL_CORE } from "../config"
import Header from "./Header.vue"
import { toast } from "vue3-toastify"

const dashboardData = ref(null)
const startDatePeriod = ref("")
const endDatePeriod = ref("")
const startDateExpensive = ref("")
const endDateExpensive = ref("")

onMounted(async () => {
  try {
    const token = localStorage.getItem("access_token")
    const res = await fetch(`${API_URL_CORE}/dashboard/user`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error("Ошибка загрузки данных")
    dashboardData.value = await res.json()
  } catch (err) {
    toast.error(err.message)
  }
})

const downloadReport = async (type) => {
  let start = type === "period" ? startDatePeriod.value : startDateExpensive.value
  let end = type === "period" ? endDatePeriod.value : endDateExpensive.value

  if (!start || !end) {
    toast.error("Выберите даты")
    return
  }

  const url =
    type === "period"
      ? `${API_URL_CORE}/report/period?start=${start}&end=${end}`
      : `${API_URL_CORE}/report/expensive?start=${start}&end=${end}`

  const fileName =
    type === "period"
      ? `отчёт_за_${start}_${end}.xlsx`
      : `отчёт_программ_за_${start}_${end}.xlsx`

  try {
    const token = localStorage.getItem("access_token")
    const res = await fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    if (!res.ok) throw new Error("Ошибка скачивания отчёта")

    const blob = await res.blob()
    const link = document.createElement("a")
    link.href = URL.createObjectURL(blob)
    link.download = fileName
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)

    toast.success("Отчёт успешно скачан")
  } catch (err) {
    toast.error(err.message)
  }
}
</script>
