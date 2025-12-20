<template>
  <Header title="Создать нового слушателя" />

  <div style="padding: 100px 20px 20px 20px;">
    <form @submit.prevent="createListener" class="row g-4">

      <!-- Основная информация -->
      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Основная информация</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="listener.first_name" @input="onlyLetters('first_name')" class="form-control"
                   placeholder="Имя" required>
            <input v-model="listener.second_name" @input="onlyLetters('second_name')" class="form-control"
                   placeholder="Фамилия" required>
            <input v-model="listener.middle_name" @input="onlyLetters('middle_name')" class="form-control"
                   placeholder="Отчество">
            <input type="date" v-model="listener.date_of_birth" class="form-control" required>
            <input v-model="listener.snils" @input="formatSnils" class="form-control"
                   placeholder="СНИЛС (xxx-xxx-xxx xx)" maxlength="14">
            <input v-model="listener.contact_phone" @input="formatPhone" class="form-control"
                   placeholder="Телефон (+7XXXXXXXXXX)" maxlength="12">
            <input v-model="listener.email" @input="validateEmail"
                   :class="['form-control', { 'is-invalid': !emailValid && listener.email }]"
                   placeholder="Email">
            <div v-if="!emailValid && listener.email" class="invalid-feedback d-block">
              Некорректный email
            </div>
          </div>
        </div>
      </div>

      <!-- Паспорт -->
      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Паспорт</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="passport.place_birth" class="form-control" placeholder="Место рождения">
            <input v-model="passport.citizenship" class="form-control" placeholder="Гражданство">
            <select v-model="passport.gender" class="form-select">
              <option>Мужской</option>
              <option>Женский</option>
            </select>
            <input v-model="passport.seria" @input="numbersOnly('passport','seria',4)" class="form-control"
                   placeholder="Серия (4 цифры)">
            <input v-model="passport.number" @input="numbersOnly('passport','number',6)" class="form-control"
                   placeholder="Номер (6 цифр)">
            <input v-model="passport.passport_given" class="form-control" placeholder="Кем выдан">
            <input type="date" v-model="passport.date_given" class="form-control">
            <input v-model="passport.code" @input="formatPassportCode" class="form-control"
                   placeholder="Код (xxx-xxx)" maxlength="7">
          </div>
        </div>
      </div>

      <!-- Адрес регистрации -->
      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Адрес регистрации</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="registration_address.mail_index" @input="formatIndex" class="form-control"
                   placeholder="Индекс (6 цифр)">
            <input v-model="registration_address.region" @input="onlyLettersAddress('region')" class="form-control"
                   placeholder="Регион">
            <input v-model="registration_address.city" @input="onlyLettersAddress('city')" class="form-control"
                   placeholder="Город">
            <input v-model="registration_address.street" class="form-control" placeholder="Улица">
            <input v-model="registration_address.house" class="form-control" placeholder="Дом">
            <input v-model="registration_address.building" class="form-control" placeholder="Корпус">
            <input v-model="registration_address.apartment" class="form-control" placeholder="Квартира">
          </div>
        </div>
      </div>

      <!-- Образование -->
      <div class="col-md-6">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Образование</h5>
          <div class="d-flex flex-column gap-2">
            <div class="form-check mb-2">
              <input class="form-check-input" type="checkbox" v-model="listener.looting_education" id="lootingEdu" @change="onLootingChange">
              <label class="form-check-label" for="lootingEdu">Получает образование</label>
            </div>
            <input v-model="education.diplom_seria" @input="numbersOnly('education','diplom_seria',6)" class="form-control"
                   placeholder="Серия диплома" :disabled="listener.looting_education">
            <input v-model="education.diplom_number" @input="numbersOnly('education','diplom_number',7)" class="form-control"
                   placeholder="Номер диплома" :disabled="listener.looting_education">
            <input type="date" v-model="education.date_given" class="form-control" :disabled="listener.looting_education">
            <input v-model="education.city" @input="onlyLettersAddress('city')" class="form-control"
                   placeholder="Город" :disabled="listener.looting_education">
            <input v-model="education.region" @input="onlyLettersAddress('region')" class="form-control"
                   placeholder="Регион" :disabled="listener.looting_education">
            <input v-model="education.educational_institution" class="form-control" placeholder="Учебное заведение" :disabled="listener.looting_education">
            <input v-model="education.speciality" class="form-control" placeholder="Специальность" :disabled="listener.looting_education">
            <select v-model="education.level_education" class="form-select" :disabled="listener.looting_education">
              <option value="">Выберите уровень</option>
              <option v-for="lvl in levels" :key="lvl.id_level_education" :value="lvl.id_level_education">
                {{ lvl.education }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <!-- Место работы -->
      <div class="col-md-6">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Место работы</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="placeWork.name_company" class="form-control" placeholder="Компания">
            <input v-model="placeWork.job_title" class="form-control" placeholder="Должность">
            <input type="number" min="0" v-model.number="placeWork.all_experience" class="form-control"
                   placeholder="Общий стаж (лет)">
            <input type="number" min="0" v-model.number="placeWork.job_title_expirience" class="form-control"
                   placeholder="Стаж по должности (лет)">
          </div>
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="goBack">Назад</button>
        <button type="submit" class="btn btn-success px-4" :disabled="loading">
          {{ loading ? 'Сохранение...' : 'Создать' }}
        </button>
      </div>

    </form>
  </div>
</template>


<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import { useRoute } from 'vue-router'
const route = useRoute()
const idLegalEntity = route.query.id_legalentity || ''


const router = useRouter()
const loading = ref(false)

const listener = ref({ first_name: '', second_name: '', middle_name: '', date_of_birth: '', snils: '', contact_phone: '', email: '', looting_education: false })
const passport = ref({ place_birth: '', citizenship: '', gender: 'Мужской', seria: '', number: '', passport_given: '', date_given: '', code: '' })
const registration_address = ref({ mail_index: '', region: '', city: '', street: '', house: '', building: '', apartment: '' })
const education = ref({ diplom_seria: '', diplom_number: '', date_given: '', city: '', region: '', educational_institution: '', speciality: '', level_education: '' })
const placeWork = ref({ name_company: '', job_title: '', all_experience: '', job_title_expirience: '' })

const levels = ref([])
const emailValid = ref(true)

onMounted(async () => {
  try {
    const token = localStorage.getItem('access_token')
    const res = await fetch(`${API_URL_CORE}/leveleducation/?filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    const data = await res.json()
    levels.value = data.data || []
  } catch {
    toast.error("Ошибка загрузки уровней образования")
  }
})

const goBack = () => router.push('/listeners')

const onlyLetters = field => {
  listener.value[field] = listener.value[field].replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}
const onlyLettersAddress = field => {
  registration_address.value[field] = registration_address.value[field].replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}
const numbersOnly = (group, field, max) => {
  const target = group === 'passport' ? passport : education
  target.value[field] = target.value[field].replace(/\D/g, '').slice(0, max)
}
const validateEmail = () => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  emailValid.value = re.test(listener.value.email)
}
const formatSnils = () => {
  let v = listener.value.snils.replace(/\D/g, '').slice(0, 11)
  if (v.length > 9) v = v.slice(0, 3) + '-' + v.slice(3, 6) + '-' + v.slice(6, 9) + ' ' + v.slice(9)
  else if (v.length > 6) v = v.slice(0, 3) + '-' + v.slice(3, 6) + '-' + v.slice(6)
  else if (v.length > 3) v = v.slice(0, 3) + '-' + v.slice(3)
  listener.value.snils = v
}
const formatPhone = () => {
  let v = listener.value.contact_phone.replace(/\D/g, '')
  if (!v.startsWith('7')) v = '7' + v
  listener.value.contact_phone = '+' + v.slice(0, 11)
}
const formatPassportCode = () => {
  let v = passport.value.code.replace(/\D/g, '').slice(0, 6)
  if (v.length > 3) v = v.slice(0, 3) + '-' + v.slice(3)
  passport.value.code = v
}
const formatIndex = () => {
  registration_address.value.mail_index = registration_address.value.mail_index.replace(/\D/g, '').slice(0, 6)
}
const onLootingChange = () => {
  if (listener.value.looting_education) {
    education.value = { diplom_seria: '', diplom_number: '', date_given: '', city: '', region: '', educational_institution: '', speciality: '', level_education: '' }
  }
}
const createListener = async () => {
  if (!emailValid.value) return toast.error("Проверьте правильность Email")

  loading.value = true
  const token = localStorage.getItem('access_token')
  const listenerPayload = { ...listener.value }
  // Явно указываем флаг, чтобы не зависеть от пропусков сериализации
  listenerPayload.looting_education = !!listener.value.looting_education
  if (idLegalEntity) listenerPayload.id_legalentity = idLegalEntity

  const payload = {
    listener: listenerPayload,
    passport: passport.value,
    registration_address: registration_address.value
  }

  if (!listener.value.looting_education && Object.values(education.value).some(v => v)) payload.education = education.value
  if (Object.values(placeWork.value).some(v => v)) payload.placeWork = placeWork.value

  try {
    const res = await fetch(`${API_URL_CORE}/listener/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload)
    })

    if (!res.ok) throw new Error("Ошибка создания")

    toast.success("Слушатель успешно создан")
    setTimeout(() => {
      router.push('/listeners')
    }, 2000)
  } catch (err) {
    toast.error(err.message || "Ошибка при создании")
  } finally {
    loading.value = false
  }
}

</script>

<style scoped>
.invalid {
  border-color: red !important;
}
</style>
