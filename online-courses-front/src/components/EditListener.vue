<template>
  <Header title="Редактирование слушателя" />

  <div style="padding: 100px 20px 20px 20px;">
    <form @submit.prevent="updateListener" class="row g-4">

      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Личные данные</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="listener.first_name" @input="onlyLetters('first_name')" class="form-control" placeholder="Имя" required />
            <input v-model="listener.second_name" @input="onlyLetters('second_name')" class="form-control" placeholder="Фамилия" required />
            <input v-model="listener.middle_name" @input="onlyLetters('middle_name')" class="form-control" placeholder="Отчество" />
            <input type="date" v-model="listener.date_of_birth" class="form-control" required />
            <input v-model="listener.snils" @input="formatSnils" class="form-control" placeholder="СНИЛС (xxx-xxx-xxx xx)" maxlength="14" />
            <input v-model="listener.contact_phone" @input="formatPhone" class="form-control" placeholder="Телефон (+7XXXXXXXXXX)" maxlength="12" />
            <input v-model="listener.email" @input="validateEmail"
                   :class="['form-control', { 'is-invalid': !emailValid && listener.email }]"
                   placeholder="Email" />
            <div v-if="!emailValid && listener.email" class="invalid-feedback d-block">
              Некорректный email
            </div>
          </div>
        </div>
      </div>

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
            <input v-model="passport.seria" @input="numbersOnly('passport','seria',4)" class="form-control" placeholder="Серия (4 цифры)">
            <input v-model="passport.number" @input="numbersOnly('passport','number',6)" class="form-control" placeholder="Номер (6 цифр)">
            <input v-model="passport.passport_given" class="form-control" placeholder="Кем выдан">
            <input type="date" v-model="passport.date_given" class="form-control">
            <input v-model="passport.code" @input="formatPassportCode" class="form-control" placeholder="Код (xxx-xxx)" maxlength="7">
          </div>
        </div>
      </div>

      <div class="col-md-4">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Адрес регистрации</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="registrationAddress.mail_index" @input="formatIndex" class="form-control" placeholder="Индекс (6 цифр)">
            <input v-model="registrationAddress.region" @input="onlyLettersAddress('region')" class="form-control" placeholder="Регион">
            <input v-model="registrationAddress.city" @input="onlyLettersAddress('city')" class="form-control" placeholder="Город">
            <input v-model="registrationAddress.street" class="form-control" placeholder="Улица">
            <input v-model="registrationAddress.house" class="form-control" placeholder="Дом">
            <input v-model="registrationAddress.building" class="form-control" placeholder="Корпус">
            <input v-model="registrationAddress.apartment" class="form-control" placeholder="Квартира">
          </div>
        </div>
      </div>

      <div class="col-md-6">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Образование</h5>
          <div class="d-flex flex-column gap-2">
            <div class="form-check mb-2">
              <input class="form-check-input" type="checkbox" v-model="listener.looting_education" id="lootingEduEdit" @change="onLootingChange">
              <label class="form-check-label" for="lootingEduEdit">Получает образование</label>
            </div>
            <input v-model="education.diplom_seria" @input="numbersOnly('education','diplom_seria',6)" class="form-control" placeholder="Серия диплома" :disabled="listener.looting_education">
            <input v-model="education.diplom_number" @input="numbersOnly('education','diplom_number',7)" class="form-control" placeholder="Номер диплома" :disabled="listener.looting_education">
            <input type="date" v-model="education.date_given" class="form-control" :disabled="listener.looting_education">
            <input v-model="education.city" @input="onlyLettersAddress('city')" class="form-control" placeholder="Город" :disabled="listener.looting_education">
            <input v-model="education.region" @input="onlyLettersAddress('region')" class="form-control" placeholder="Регион" :disabled="listener.looting_education">
            <input v-model="education.educational_institution" class="form-control" placeholder="Учебное заведение" :disabled="listener.looting_education">
            <input v-model="education.speciality" class="form-control" placeholder="Специальность" :disabled="listener.looting_education">
            <select v-model="education.level_education" class="form-select" :disabled="listener.looting_education">
              <option value="">Выберите уровень образования</option>
              <option v-for="lvl in educationLevels" :key="lvl.id_level_education" :value="lvl.id_level_education">
                {{ lvl.education }}
              </option>
            </select>
          </div>
        </div>
      </div>

      <div class="col-md-6">
        <div class="card p-3 shadow-sm card-block">
          <h5 class="card-title mb-3">Место работы</h5>
          <div class="d-flex flex-column gap-2">
            <input v-model="placeWork.name_company" class="form-control" placeholder="Компания">
            <input v-model="placeWork.job_title" class="form-control" placeholder="Должность">
            <input type="number" min="0" v-model.number="placeWork.all_experience" class="form-control" placeholder="Общий стаж (лет)">
            <input type="number" min="0" v-model.number="placeWork.job_title_expirience" class="form-control" placeholder="Стаж по должности (лет)">
          </div>
        </div>
      </div>

      <div class="col-12 d-flex justify-content-between">
        <button type="button" class="btn btn-secondary" @click="goBack" :disabled="loading">Отмена</button>
        <button type="submit" class="btn btn-success px-4" :disabled="loading">
          {{ loading ? 'Сохранение...' : 'Сохранить изменения' }}
        </button>
      </div>

    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { API_URL_CORE } from '../config'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'

const router = useRouter()
const route = useRoute()
const id = route.params.id
const token = localStorage.getItem('access_token')

const loading = ref(false)
const emailValid = ref(true)

const listener = ref({})
const passport = ref({})
const registrationAddress = ref({})
const education = ref({})
const placeWork = ref({})
const educationLevels = ref([])

const onlyLetters = field => {
  listener.value[field] = (listener.value[field] || '').replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}
const onlyLettersAddress = field => {
  registrationAddress.value[field] = (registrationAddress.value[field] || '').replace(/[^А-Яа-яЁёA-Za-z\s-]/g, '')
}
const numbersOnly = (group, field, max) => {
  const target = group === 'passport' ? passport : education
  target.value[field] = (target.value[field] || '').toString().replace(/\D/g, '').slice(0, max)
}
const validateEmail = () => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  emailValid.value = re.test((listener.value.email || '').toString())
}
const formatSnils = () => {
  let v = (listener.value.snils || '').toString().replace(/\D/g, '').slice(0, 11)
  if (v.length > 9) v = v.slice(0, 3) + '-' + v.slice(3, 6) + '-' + v.slice(6, 9) + ' ' + v.slice(9)
  else if (v.length > 6) v = v.slice(0, 3) + '-' + v.slice(3, 6) + '-' + v.slice(6)
  else if (v.length > 3) v = v.slice(0, 3) + '-' + v.slice(3)
  listener.value.snils = v
}
const formatPhone = () => {
  let v = (listener.value.contact_phone || '').toString().replace(/\D/g, '')
  if (!v.startsWith('7')) v = '7' + v
  listener.value.contact_phone = '+' + v.slice(0, 11)
}
const formatPassportCode = () => {
  let v = (passport.value.code || '').toString().replace(/\D/g, '').slice(0, 6)
  if (v.length > 3) v = v.slice(0, 3) + '-' + v.slice(3)
  passport.value.code = v
}
const formatIndex = () => {
  registrationAddress.value.mail_index = (registrationAddress.value.mail_index || '').toString().replace(/\D/g, '').slice(0, 6)
}
const onLootingChange = () => {
  if (listener.value.looting_education) {
    education.value = { diplom_seria: '', diplom_number: '', date_given: '', city: '', region: '', educational_institution: '', speciality: '', level_education: '' }
  }
}

const normalizeToDateInput = dateStr => {
  if (!dateStr) return ''
  const first = dateStr.split(' ')[0]
  if (first.includes('T')) return first.split('T')[0]
  if (/^\d{4}-\d{2}-\d{2}$/.test(first)) return first
  const m = first.match(/(\d{4}-\d{2}-\d{2})/)
  return m ? m[1] : ''
}

const loadEducationLevels = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/leveleducation/?filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка загрузки уровней (${res.status})`)
    const result = await res.json()
    educationLevels.value = result.data || []
  } catch (err) {
    toast.error('Ошибка загрузки уровней образования')
    // don't throw further — let page try to load details (they may still show)
  }
}

const loadListenerDetails = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/listener/details/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Ошибка запроса (${res.status})`)
    const data = await res.json()

    listener.value = data.data.listener || {}
    passport.value = data.data.passport || {}
    registrationAddress.value = data.data.regaddress || {}
    education.value = { ...(data.data.education_listener || {}) }
    placeWork.value = data.data.placework || {}

    listener.value.date_of_birth = normalizeToDateInput(listener.value.date_of_birth)
    passport.value.date_given = normalizeToDateInput(passport.value.date_given)
    education.value.date_given = normalizeToDateInput(education.value.date_given)

    const current = education.value.id_level_education || education.value.level_education || ''
    if (!current) {
      education.value.level_education = ''
    } else {
      const uuidRegex = /^[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}$/
      if (uuidRegex.test(current)) {
        education.value.level_education = current
      } else {
        const foundById = educationLevels.value.find(l => String(l.id_level_education) === String(current))
        if (foundById) {
          education.value.level_education = foundById.id_level_education
        } else {
          const foundByName = educationLevels.value.find(l => l.education === current || l.education === (education.value.level_education))
          education.value.level_education = foundByName ? foundByName.id_level_education : ''
        }
      }
    }
  } catch (err) {
    toast.error('Ошибка загрузки данных слушателя')
  }
}

const updateListener = async () => {
  if (!emailValid.value) {
    toast.error('Проверьте правильность Email')
    return
  }

  loading.value = true
  try {
    let educationPayload = {
      ...education.value,
      diplom_seria: String(education.value.diplom_seria || ''),
      diplom_number: String(education.value.diplom_number || ''),
      level_education: education.value.level_education || null
    }

    if (listener.value.looting_education) {
      educationPayload = { diplom_seria: '', diplom_number: '', date_given: '', city: '', region: '', educational_institution: '', speciality: '', level_education: '' }
    }

    const payload = {
      listener: { ...listener.value, looting_education: !!listener.value.looting_education },
      passport: {
        ...passport.value,
        seria: String(passport.value.seria || ''),
        number: String(passport.value.number || '')
      },
      registration_address: {
        ...registrationAddress.value,
        mail_index: String(registrationAddress.value.mail_index || '')
      },
      education: educationPayload,
      placeWork: { ...placeWork.value }
    }

    const res = await fetch(`${API_URL_CORE}/listener/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(payload)
    })

    if (!res.ok) {
      let text = `Ошибка сохранения (${res.status})`
      try {
        const body = await res.text()
        if (body) text = body
      } catch {}
      throw new Error(text)
    }

    toast.success('Данные слушателя успешно обновлены')
    setTimeout(() => router.push('/listeners'), 2000)
  } catch (err) {
    toast.error(err.message || 'Ошибка при сохранении')
  } finally {
    loading.value = false
  }
}

const goBack = () => router.push('/listeners')

onMounted(async () => {
  loading.value = true
  try {
    await loadEducationLevels()
    await loadListenerDetails()
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.invalid-feedback {
  font-size: 13px;
}
.card-block {
  border-radius: 8px;
}
.is-invalid {
  border-color: red !important;
}
</style>
