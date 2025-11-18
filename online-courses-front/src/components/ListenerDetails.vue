<template>
  <div>
    <Header title="Информация о слушателе" />

    <div class="container py-4 mt-5">
      <div v-if="loading" class="text-center py-5">Загрузка данных...</div>

      <form v-else class="row g-4">
        <div class="col-md-4">
          <div class="card p-3 shadow-sm">
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
          <div class="card p-3 shadow-sm">
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
          <div class="card p-3 shadow-sm">
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
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Образование</h5>
            <p><strong>Учреждение:</strong> {{ education.educational_institution }}</p>
            <p><strong>Город:</strong> {{ education.city }}</p>
            <p><strong>Регион:</strong> {{ education.region }}</p>
            <p><strong>Специальность:</strong> {{ education.speciality }}</p>
            <p><strong>Уровень образования:</strong> {{ education.level_education }}</p>
            <p><strong>Серия и номер диплома:</strong> {{ education.diplom_seria }} {{ education.diplom_number }}</p>
            <p><strong>Дата выдачи:</strong> {{ formatDate(education.date_given) }}</p>
          </div>
        </div>

        <div class="col-md-6">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Место работы</h5>
            <p><strong>Компания:</strong> {{ placework.name_company }}</p>
            <p><strong>Должность:</strong> {{ placework.job_title }}</p>
            <p><strong>Общий стаж(лет):</strong> {{ placework.all_experience }}</p>
            <p><strong>Стаж по должности(лет):</strong> {{ placework.job_title_expirience }}</p>
          </div>
        </div>

        <div class="col-12">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Личные дела</h5>
            <div v-if="docsLoading">Загрузка личных дел...</div>
            <div v-else-if="docs.length === 0">Личных дел не найдено</div>
            <ul v-else class="list-unstyled">
              <li v-for="doc in docs" :key="doc" class="d-flex align-items-center mb-2">
                {{ shortenName(doc) }}
               <button type="button" @click="downloadDoc(doc)" class="btn btn-sm btn-success ms-2">Скачать</button>
              </li>
            </ul>
          </div>
        </div>

        <div class="col-12">
          <div class="card p-3 shadow-sm">
            <h5 class="card-title mb-3">Курсы слушателя</h5>
            <div v-if="coursesLoading">Загрузка курсов...</div>
            <div v-else-if="enrollments.length === 0">Слушатель не записан на курсы</div>
            <ul v-else class="list-unstyled">
              <li v-for="c in enrollments" :key="c.id_program_education">
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
  </div>
</template>

<style scoped>
.card {
  border-radius: 12px;
  padding: 16px;
}

.card p {
  display: flex;
  justify-content: flex-start;
  margin: 4px 0;
}

.card p strong {
  width: 200px; 
  text-align: left;
  margin-right: 8px;
}

h5.card-title {
  margin-top: 0;
}
</style>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import { API_URL_CORE, API_URL_DOC } from '../config'
import Header from './Header.vue'

const route = useRoute()
const router = useRouter()
const id = ref(route.params.id)
const token = localStorage.getItem('access_token')

const listener = ref({})
const passport = ref({})
const regaddress = ref({})
const education = ref({})
const placework = ref({})
const loading = ref(true)

const docs = ref([])
const docsLoading = ref(false)

const enrollments = ref([])
const coursesLoading = ref(false)

const formatDate = date => {
  if (!date) return '—'
  const match = date.match(/^(\d{4})-(\d{2})-(\d{2})/)
  return match ? `${match[3]}.${match[2]}.${match[1]}` : date
}

const loadDetails = async () => {
  try {
    if (!id.value) throw new Error('ID слушателя не найден')
    const res = await fetch(`${API_URL_CORE}/listener/details/${id.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка запроса (${res.status})`)
    const data = await res.json()
    listener.value = data.data?.listener || {}
    passport.value = data.data?.passport || {}
    regaddress.value = data.data?.regaddress || {}
    education.value = data.data?.education_listener || {}
    placework.value = data.data?.placework || {}

    if (listener.value.snils) loadDocs(listener.value.snils)
    loadEnrollments()
  } catch (err) {
    toast.error(err.message)
  } finally {
    loading.value = false
  }
}

const loadDocs = async (snils) => {
  docsLoading.value = true
  try {
    const res = await fetch(`http://localhost:8082/v1/doc/exists?card-name=${encodeURIComponent(snils)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) {
      docs.value = []
      return
    }
    const data = await res.json()
    docs.value = data || []
  } catch (err) {
    docs.value = []
    toast.error('Ошибка загрузки личных дел')
  } finally {
    docsLoading.value = false
  }
}

const shortenName = (name) => {
  return name.length > 40 ? name.slice(0, 37) + '...' : name
}

const downloadDoc = async (name) => {
  try {
    const res = await fetch(`${API_URL_DOC}/download?card-name=${encodeURIComponent(name)}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки документа (${res.status})`)
    const blob = await res.blob()
    const link = document.createElement('a')
    link.href = window.URL.createObjectURL(blob)
    link.download = name
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
  } catch (err) {
    toast.error('Не удалось скачать документ')
  }
}

const loadEnrollments = async () => {
  coursesLoading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/enrollment/details/${id.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) {
      enrollments.value = []
      return
    }
    const data = await res.json()
    enrollments.value = data.data || []
  } catch (err) {
    enrollments.value = []
    toast.error('Ошибка загрузки курсов')
  } finally {
    coursesLoading.value = false
  }
}

const goBack = () => router.push('/listeners')

const goToCreateEnrollment = () => {
  if (!id.value) {
    toast.error('ID слушателя не загружен')
    return
  }
  router.push(`/enrollment/create/${id.value}`)
}


onMounted(loadDetails)
</script>


