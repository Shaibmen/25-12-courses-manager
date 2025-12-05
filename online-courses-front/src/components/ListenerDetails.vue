<template>
  <Header title="Информация о слушателе" />

  <div style="padding: 100px 20px 20px 20px;">
    <div v-if="loading" class="text-center py-5">Загрузка данных...</div>

    <form v-else class="row g-4">
      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Личные данные</h5>
          <p><strong>Фамилия:</strong> {{ listener.second_name }}</p>
          <p><strong>Имя:</strong> {{ listener.first_name }}</p>
          <p><strong>Отчество:</strong> {{ listener.middle_name || '—' }}</p>
          <p><strong>Дата рождения:</strong> {{ formatDate(listener.date_of_birth) }}</p>
          <p><strong>СНИЛС:</strong> {{ listener.snils }}</p>
          <p><strong>Телефон:</strong> {{ listener.contact_phone }}</p>
          <p><strong>Email:</strong> {{ listener.email }}</p>
        </div>
      </div>

      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Паспорт</h5>
          <p><strong>Место рождения:</strong> {{ passport.place_birth }}</p>
          <p><strong>Гражданство:</strong> {{ passport.citizenship }}</p>
          <p><strong>Пол:</strong> {{ passport.gender }}</p>
          <p><strong>Серия и номер:</strong> {{ passport.seria }} {{ passport.number }}</p>
          <p><strong>Кем выдан:</strong> {{ passport.passport_given }}</p>
          <p><strong>Дата выдачи:</strong> {{ formatDate(passport.date_given) }}</p>
          <p><strong>Код подразделения:</strong> {{ passport.code }}</p>
        </div>
      </div>

      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Адрес регистрации</h5>
          <p><strong>Почтовый индекс:</strong> {{ regaddress.mail_index }}</p>
          <p><strong>Регион:</strong> {{ regaddress.region }}</p>
          <p><strong>Город:</strong> {{ regaddress.city }}</p>
          <p><strong>Улица:</strong> {{ regaddress.street }}</p>
          <p><strong>Дом:</strong> {{ regaddress.house }}</p>
          <p><strong>Корпус:</strong> {{ regaddress.building }}</p>
          <p><strong>Квартира:</strong> {{ regaddress.apartment }}</p>
        </div>
      </div>

      <div class="col-md-6">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Образование</h5>
          <div v-if="listener.looting_education">
            <p><strong>Получает образование сейчас:</strong> Да</p>
          </div>
          <div v-else>
            <p><strong>Учреждение:</strong> {{ education.educational_institution }}</p>
            <p><strong>Город:</strong> {{ education.city }}</p>
            <p><strong>Регион:</strong> {{ education.region }}</p>
            <p><strong>Специальность:</strong> {{ education.speciality }}</p>
            <p><strong>Уровень образования:</strong> {{ education.level_education }}</p>
            <p><strong>Диплом:</strong> {{ education.diplom_seria }} {{ education.diplom_number }}</p>
            <p><strong>Дата выдачи:</strong> {{ formatDate(education.date_given) }}</p>
          </div>
        </div>
      </div>

      <div class="col-md-6">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Место работы</h5>
          <p><strong>Компания:</strong> {{ placework?.name_company || '—' }}</p>
          <p><strong>Должность:</strong> {{ placework?.job_title || '—' }}</p>
          <p><strong>Общий стаж (лет):</strong> {{ placework?.all_experience || '—' }}</p>
          <p><strong>Стаж по должности (лет):</strong> {{ placework?.job_title_experience || '—' }}</p>
        </div>
      </div>

      <div class="col-12">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Курсы слушателя</h5>
          <div v-if="coursesLoading">Загрузка курсов...</div>
          <div v-else-if="enrollments.length === 0">Слушатель не записан на курсы</div>
          <ul v-else class="list-unstyled">
            <li v-for="c in enrollments" :key="c.id_enrollment">
              <strong>{{ c.name_prof_education }}</strong> — {{ c.current_price }} ₽,
              {{ formatDate(c.start_date) }} — {{ formatDate(c.end_date) }}
            </li>
          </ul>
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between">
        <button type="button" @click.prevent="goToCreateEnrollment" class="btn btn-success">
          Записать на курс
        </button>
        <button @click="goBack" class="btn btn-secondary">Назад</button>
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
const token = localStorage.getItem('access_token')

const listener = ref({})
const passport = ref({})
const regaddress = ref({})
const education = ref({})
const placework = ref({})
const enrollments = ref([])

const loading = ref(true)
const coursesLoading = ref(false)

const formatDate = (date) => {
  if (!date) return '—'
  const d = new Date(date)
  return isNaN(d) ? '—' : d.toLocaleDateString('ru-RU')
}

const loadDetails = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/listener/details/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки (${res.status})`)
    const data = await res.json()

    listener.value = data.data.listener || {}
    passport.value = data.data.passport || {}
    regaddress.value = data.data.regaddress || {}
    education.value = data.data.education_listener || {}
    placework.value = data.data.placework || {}

    loadEnrollments()
  } catch (err) {
    toast.error(err.message || 'Ошибка загрузки данных')
  } finally {
    loading.value = false
  }
}

const loadEnrollments = async () => {
  coursesLoading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/enrollment/details/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    enrollments.value = res.ok ? (await res.json()).data || [] : []
  } finally {
    coursesLoading.value = false
  }
}

const goToCreateEnrollment = () => router.push(`/enrollment/create/${id}`)
const goBack = () => router.push('/listeners')

onMounted(loadDetails)
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
.card p {
  margin: 4px 0;
}
.card p strong {
  min-width: 160px;
  display: inline-block;
}
</style>