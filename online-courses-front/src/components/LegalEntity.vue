<template>
  <Header title="Юридические лица" />

  <div style="">
    <div style="display: flex; flex-wrap: wrap; justify-content: space-between; align-items: center; margin-bottom: 20px; gap: 10px;">
      <input
        v-model="filter"
        placeholder="Поиск по названию компании"
        style="padding: 8px 12px; min-width: 250px; border-radius: 6px; border: 1px solid #ccc;"
      />
      <div style="display: flex; flex-wrap: wrap; gap: 10px;">
        <button class="btn btn-primary" @click="$router.push('/legalentity/create')">Добавить юридическое лицо</button>
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>

    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 1200px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Компания</th>
            <th>ИНН</th>
            <th>КПП</th>
            <th>ОГРН</th>
            <th>Телефон</th>
            <th>Email</th>
            <th>Контактное лицо</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="l in legalList" :key="l.ID_Legalentity">
            <td :title="l.NameCompany">{{ l.NameCompany }}</td>
            <td>{{ l.Inn }}</td>
            <td>{{ l.Kpp }}</td>
            <td>{{ l.Ogrn }}</td>
            <td>{{ l.Phone }}</td>
            <td>{{ l.Email }}</td>
            <td>
              {{ l.SecondName }} {{ l.FirstName }} {{ l.MiddleName || '' }}
            </td>
            <td style="white-space: nowrap;">
              <div style="display: flex; gap: 8px;">
                <button class="btn btn-success btn-sm" @click="viewEntity(l.ID_Legalentity)">Подробнее</button>
                <button class="btn btn-warning btn-sm" @click="editEntity(l.ID_Legalentity)">Изменить</button>
                <button class="btn btn-danger btn-sm" @click="deleteEntity(l.ID_Legalentity)">Удалить</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>

      <div style="margin-top: 20px; display: flex; justify-content: center; gap: 10px;">
        <button class="btn btn-primary" :disabled="page <= 1" @click="prevPage">Назад</button>
        <span>Страница {{ page }}</span>
        <button class="btn btn-primary" :disabled="!hasMore" @click="nextPage">Вперёд</button>
      </div>
    </div>
  </div>

  <ConfirmModal
    ref="confirmModal"
    title="Удаление юридического лица"
    message="Вы точно хотите удалить эту организацию?"
  />
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_CORE } from '../config'
import Header from './Header.vue'
import ConfirmModal from '../components/ConfirmModal.vue'

const router = useRouter()
const token = localStorage.getItem('access_token')

const legalList = ref([])
const loading = ref(false)
const page = ref(1)
const filter = ref('')
const hasMore = ref(false)
let debounceTimer = null

const confirmModal = ref(null)

const loadData = async () => {
  loading.value = true

  try {
    const res = await fetch(
      `${API_URL_CORE}/legalentity/?page=${page.value}&filter=${encodeURIComponent(filter.value)}`,
      { headers: { Authorization: `Bearer ${token}` } }
    )

    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)

    const data = await res.json()
    legalList.value = Array.isArray(data.data) ? data.data : []
    hasMore.value = data.data?.length === 25
  } catch (err) {
    toast.error(err.message || 'Ошибка загрузки данных')
  } finally {
    loading.value = false
  }
}

const deleteEntity = (id) => {
  if (!confirmModal.value) return

  confirmModal.value.open(async () => {
    try {
      const res = await fetch(`${API_URL_CORE}/legalentity/${id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      })

      if (!res.ok) {
        let text = `Ошибка удаления (${res.status})`
        try {
          const body = await res.json()
          if (body?.message) text = body.message
        } catch {}
        throw new Error(text)
      }

      toast.success('Организация удалена')
      loadData()
    } catch (err) {
      toast.error(err.message || 'Ошибка при удалении')
    }
  })
}

const editEntity = (id) => router.push(`/legalentity/edit/${id}`)
const viewEntity = (id) => router.push(`/legalentity/${id}`)

watch(filter, () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    page.value = 1
    loadData()
  }, 400)
})

const prevPage = () => {
  if (page.value > 1) {
    page.value--
    loadData()
  }
}

const nextPage = () => {
  page.value++
  loadData()
}

const goBack = () => router.push('/dashboard/worker')

onMounted(loadData)
</script>

<style scoped>
.table-hover tbody tr:hover {
  background-color: #e2f0d9;
  cursor: pointer;
}
th, td {
  vertical-align: middle;
}
</style>
