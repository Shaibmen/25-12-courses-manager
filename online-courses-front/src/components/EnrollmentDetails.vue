<template>
  <Header title="Курсы слушателя" />
    
  <div style="padding: 10px 20px 20px 20px;">
    <div v-if="listener" class="card p-3 mb-3 shadow-sm card-block" style="display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap;">
    <div>
      <strong>{{ listener.second_name }} {{ listener.first_name }}</strong>
      <br />
      <span>{{ listener.contact_phone }} | {{ listener.email }}</span>
    </div>
    <div>
      <button class="btn btn-primary btn-md" style="margin-top: 10px;" @click="viewListener(listener.id_listener)">Просмотр профиля</button>
    </div>
  </div>


    <div style="display: flex; flex-wrap: wrap; gap: 10px; margin-bottom: 15px; margin-right: 20px; justify-content: flex-end;">
      <button class="btn btn-secondary" @click="goBack">Назад</button>
    </div>

    <div v-if="loading" class="text-center w-100">Загрузка...</div>
    <div v-else-if="error">{{ error }}</div>
    <div v-else style="overflow-x: auto;">
      <table class="table table-striped table-hover w-100" style="border-collapse: separate; border-spacing: 0; min-width: 800px;">
        <thead class="table-light sticky-top" style="top: 0; z-index: 2;">
          <tr>
            <th>Курс</th>
            <th>Тип</th>
            <th>Подразделение</th>
            <th>Группа</th>
            <th>Тип обучения</th>
            <th>Начало</th>
            <th>Окончание</th>
            <th>Цена</th>
            <th>Действия</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="item in enrollments" :key="item.id_program_education">
             <td>{{ item.name_prof_education }}</td>
            <td>{{ item.education_type }}</td>
            <td>{{ item.division_education }}</td>

            <td>{{ item.group || '—' }}</td>
            <td>{{ item.type_of_retraining || '—' }}</td>

            <td>{{ formatDate(item.start_date) }}</td>
            <td>{{ formatDate(item.end_date) }}</td>
            <td>{{ item.current_price }} ₽</td>
            <td style="white-space: nowrap;">
              <div style="display: flex; gap: 8px;">
                <button class="btn btn-warning btn-sm" @click="editRecord(item)">Изменить</button>
                <button class="btn btn-danger btn-sm" @click="deleteRecord(item)">Удалить</button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <ConfirmModal
      ref="confirmModal"
      title="Удаление записи"
      message="Вы точно хотите удалить эту запись?"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import ConfirmModal from '../components/ConfirmModal.vue'

const route = useRoute()
const router = useRouter()
const listenerId = route.params.listenerId

const loading = ref(false)
const error = ref(null)
const enrollments = ref([])
const listener = ref(null)
const token = localStorage.getItem('access_token')
const confirmModal = ref(null)

const loadEnrollments = async () => {
  const res = await fetch(`${API_URL_CORE}/enrollment/details/${listenerId}`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  const data = await res.json()
  enrollments.value = data.data || []
}

const loadListenerInfo = async () => {
  const res = await fetch(`${API_URL_CORE}/listener/details/${listenerId}`, {
    headers: { Authorization: `Bearer ${token}` }
  })
  const data = await res.json()
  listener.value = data.data.listener
}

const loadAll = async () => {
  try {
    loading.value = true
    await Promise.all([loadEnrollments(), loadListenerInfo()])
  } catch (err) {
    toast.error(err.message)
    error.value = err.message
  } finally {
    loading.value = false
  }
}

onMounted(loadAll)

const editRecord = (item) => router.push(`/enrollment/edit/${listenerId}/${item.id_program_education}`)

const deleteRecord = async (item) => {
  if (!confirmModal.value) return
  confirmModal.value.open(async () => {
    try {
      const res = await fetch(`${API_URL_CORE}/enrollment/${listenerId}/${item.id_program_education}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` }
      })
      if (!res.ok) throw new Error('Не удалось удалить запись')
      toast.success('Запись успешно удалена')
      await loadAll()
    } catch (err) {
      toast.error(err.message)
    }
  })
}

const formatDate = (str) => str?.split(" ")[0] ?? ''
const goBack = () => router.push(`/enrollments`)

const viewListener = (uuid) => router.push(`/listeners/${uuid}`)


</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
</style>
