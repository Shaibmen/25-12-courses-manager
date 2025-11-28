<template>
  <div>
    <Header title="Информация о слушателе" />

    <div class="container py-4 mt-5">
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
            <p><strong>Учреждение:</strong> {{ education.educational_institution }}</p>
            <p><strong>Город:</strong> {{ education.city }}</p>
            <p><strong>Регион:</strong> {{ education.region }}</p>
            <p><strong>Специальность:</strong> {{ education.speciality }}</p>
            <p><strong>Уровень образования:</strong> {{ education.level_education }}</p>
            <p><strong>Диплом:</strong> {{ education.diplom_seria }} {{ education.diplom_number }}</p>
            <p><strong>Дата выдачи:</strong> {{ formatDate(education.date_given) }}</p>
          </div>
        </div>

        <div class="col-md-6">
          <div class="card p-3 shadow-sm card-block">
            <h5 class="card-title mb-3">Место работы</h5>
            <p><strong>Компания:</strong> {{ placework.name_company }}</p>
            <p><strong>Должность:</strong> {{ placework.job_title }}</p>
            <p><strong>Общий стаж (лет):</strong> {{ placework.all_experience }}</p>
            <p><strong>Стаж по должности (лет):</strong> {{ placework.job_title_expirience }}</p>
          </div>
        </div>

        <div class="col-12">
          <div class="card p-3 shadow-sm card-block">
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
          <div class="card p-3 shadow-sm card-block">
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

        <div class="col-12">
          <div class="card p-3 shadow-sm card-block">
            <div class="d-flex align-items-center mb-3">
              <input type="checkbox" v-model="hasContractor" class="form-check-input me-2" />
              <label class="fw-bold">Есть заказчик</label>
            </div>

            <div v-if="contractor">
              <h5 class="card-title mb-2">Текущий заказчик</h5>
              <p><strong>ФИО:</strong> {{ contractor.contractor.second_name }} {{ contractor.contractor.first_name }} {{ contractor.contractor.middle_name }}</p>
              <p><strong>Телефон:</strong> {{ contractor.contractor.contact_phone }}</p>
              <p><strong>Email:</strong> {{ contractor.contractor.email }}</p>
            </div>

            <button type="button" @click="openContractorModal" class="btn btn-primary">
              Добавить заказчика
            </button>

            <button
              type="button"
              class="btn btn-danger mt-2"
              @click="confirmDeleteContractor(contractor.Contractor.id_contractor)"
            >
              Удалить заказчика
            </button>
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

    <div class="modal fade show" v-if="contractorModal" style="display:block; background: rgba(0,0,0,.5);">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title">Добавление заказчика</h5>
            <button type="button" class="btn-close" @click="closeContractorModal"></button>

            
          </div>

          <div class="modal-body">
            <h5 class="mb-3">Данные заказчика</h5>

            <div class="row g-3">

              <div class="col-md-4">
                <label class="form-label">Фамилия</label>
                <input
                  v-model="contractorForm.contractor.second_name"
                  @input="onlyLetters(contractorForm.contractor, 'second_name')"
                  :class="{'form-control': true, 'is-invalid': errors.second_name}"
                >
                <div class="invalid-feedback" v-if="errors.second_name">Только буквы</div>
              </div>

              <div class="col-md-4">
                <label class="form-label">Имя</label>
                <input
                  v-model="contractorForm.contractor.first_name"
                  @input="onlyLetters(contractorForm.contractor, 'first_name')"
                  :class="{'form-control': true, 'is-invalid': errors.first_name}"
                >
                <div class="invalid-feedback" v-if="errors.first_name">Только буквы</div>
              </div>

              <div class="col-md-4">
                <label class="form-label">Отчество</label>
                <input
                  v-model="contractorForm.contractor.middle_name"
                  @input="onlyLetters(contractorForm.contractor, 'middle_name')"
                  class="form-control"
                >
              </div>

              <div class="col-md-6">
                <label class="form-label">Телефон</label>
                <input
                  v-model="contractorForm.contractor.contact_phone"
                  @input="formatPhone"
                  :class="{'form-control': true, 'is-invalid': errors.phone}"
                >
                <div class="invalid-feedback" v-if="errors.phone">Формат: +7XXXXXXXXXX</div>
              </div>

              <div class="col-md-6">
                <label class="form-label">Email</label>
                <input
                  v-model="contractorForm.contractor.email"
                  @input="validateEmail"
                  :class="{'form-control': true, 'is-invalid': errors.email}"
                >
                <div class="invalid-feedback" v-if="errors.email">Некорректный Email</div>
              </div>
            </div>

            <hr class="my-4">

            <h5 class="mb-3">Паспорт заказчика</h5>
            <div class="row g-3">

              <div class="col-md-6">
                <label class="form-label">Место рождения</label>
                <input v-model="contractorForm.passport.place_birth" class="form-control">
              </div>

              <div class="col-md-6">
                <label class="form-label">Гражданство</label>
                <input v-model="contractorForm.passport.citizenship" class="form-control">
              </div>

              <div class="col-md-4">
                <label class="form-label">Пол</label>
                <select
                  v-model="contractorForm.passport.gender"
                  :class="{'form-select': true, 'is-invalid': errors.gender}"
                >
                  <option value="">Выберите...</option>
                  <option value="Мужской">Мужской</option>
                  <option value="Женский">Женский</option>
                </select>
                <div class="invalid-feedback" v-if="errors.gender">Выберите пол</div>
              </div>

              <div class="col-md-4">
                <label class="form-label">Серия</label>
                <input
                  v-model="contractorForm.passport.seria"
                  @input="digitsLimit(contractorForm.passport, 'seria', 4)"
                  :class="{'form-control': true, 'is-invalid': errors.seria}"
                >
                <div class="invalid-feedback" v-if="errors.seria">4 цифры</div>
              </div>

              <div class="col-md-4">
                <label class="form-label">Номер</label>
                <input
                  v-model="contractorForm.passport.number"
                  @input="digitsLimit(contractorForm.passport, 'number', 6)"
                  :class="{'form-control': true, 'is-invalid': errors.number}"
                >
                <div class="invalid-feedback" v-if="errors.number">6 цифр</div>
              </div>

              <div class="col-md-12">
                <label class="form-label">Кем выдан</label>
                <input v-model="contractorForm.passport.passport_given" class="form-control">
              </div>

              <div class="col-md-6">
                <label class="form-label">Дата выдачи</label>
                <input type="date" v-model="contractorForm.passport.date_given" class="form-control">
              </div>

              <div class="col-md-6">
                <label class="form-label">Код подразделения</label>
                <input
                  v-model="contractorForm.passport.code"
                  @input="formatCode"
                  :class="{'form-control': true, 'is-invalid': errors.code}"
                >
                <div class="invalid-feedback" v-if="errors.code">Формат: 000-000</div>
              </div>
            </div>

            <hr class="my-4">

            <h5 class="mb-3">Адрес регистрации заказчика</h5>

            <div class="row g-3">

              <div class="col-md-4">
                <label class="form-label">Индекс</label>
                <input
                  v-model="contractorForm.reg_address.mail_index"
                  @input="digitsLimit(contractorForm.reg_address, 'mail_index', 6)"
                  :class="{'form-control': true, 'is-invalid': errors.index}"
                >
                <div class="invalid-feedback" v-if="errors.index">6 цифр</div>
              </div>

              <div class="col-md-4">
                <label class="form-label">Регион</label>
                <input v-model="contractorForm.reg_address.region" class="form-control">
              </div>

              <div class="col-md-4">
                <label class="form-label">Город</label>
                <input v-model="contractorForm.reg_address.city" class="form-control">
              </div>

              <div class="col-md-6">
                <label class="form-label">Улица</label>
                <input v-model="contractorForm.reg_address.street" class="form-control">
              </div>

              <div class="col-md-3">
                <label class="form-label">Дом</label>
                <input v-model="contractorForm.reg_address.house" class="form-control">
              </div>

              <div class="col-md-3">
                <label class="form-label">Корпус</label>
                <input v-model="contractorForm.reg_address.building" class="form-control">
              </div>

              <div class="col-md-3">
                <label class="form-label">Квартира</label>
                <input v-model="contractorForm.reg_address.apartment" class="form-control">
              </div>

            </div>

          </div>

          <div class="modal-footer">
            <button @click="closeContractorModal" class="btn btn-secondary">Отмена</button>
            <button @click="saveContractor" :disabled="!formValid" class="btn btn-success">Сохранить</button>
          </div>

        </div>
      </div>
    </div>

  </div>

  <ConfirmModal
  ref="confirmModal"
  title="Удаление заказчика"
  message="Вы точно хотите удалить этого заказчика?"
/>

</template>

<style scoped>
.card {
  border-radius: 12px;
  padding: 16px;
}

.card p {
  display: flex;
  margin: 4px 0;
}
.card p strong {
  width: 200px;
  text-align: left;
  margin-right: 8px;
}
</style>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import { API_URL_CORE, API_URL_DOC } from '../config'
import ConfirmModal from '../components/ConfirmModal.vue'  
const confirmModal = ref(null) 


const confirmDeleteContractor = () => {
  contractor.value = data.data.contractor
  ? {
      Contractor: data.data.contractor.contractor,
      Passport: data.data.contractor.passport,
      RegAddress: data.data.contractor.reg_address
    }
  : null

  const id_contractor = contractor.value.Contractor.id_contractor

  if (!confirmModal.value || typeof confirmModal.value.open !== 'function') {
    const confirmed = confirm('Вы точно хотите удалить этого заказчика?')
    if (!confirmed) return
    return deleteContractor(id_contractor)
  }

  confirmModal.value.open(async () => {
    await deleteContractor(id_contractor)
  })
}

const deleteContractor = async (id_contractor) => {
  try {
    const res = await fetch(`${API_URL_CORE}/contractor/${id_contractor}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка удаления заказчика')

    toast.success('Заказчик успешно удалён!')
    contractor.value = null
    hasContractor.value = false
  } catch (e) {
    toast.error(e.message)
  }
}



const route = useRoute()
const router = useRouter()

const id = ref(route.params.id)
const token = localStorage.getItem('access_token')

const listener = ref({})
const passport = ref({
  place_birth: "",
  citizenship: "",
  gender: "",
  seria: "",
  number: "",
  passport_given: "",
  date_given: "",
  code: ""
})
const regaddress = ref({
  mail_index: "",
  region: "",
  city: "",
  street: "",
  house: "",
  building: "",
  apartment: ""
})
const education = ref({
  educational_institution: "",
  city: "",
  region: "",
  speciality: "",
  level_education: "",
  diplom_seria: "",
  diplom_number: "",
  date_given: ""
})
const placework = ref({
  name_company: "",
  job_title: "",
  all_experience: "",
  job_title_expirience: ""
})
const contractor = ref(null)

const loading = ref(true)
const docs = ref([])
const docsLoading = ref(false)
const enrollments = ref([])
const coursesLoading = ref(false)

const hasContractor = ref(false)
const contractorModal = ref(false)

const contractorForm = ref({
  contractor: {
    first_name: "",
    second_name: "",
    middle_name: "",
    contact_phone: "",
    email: ""
  },
  passport: {
    place_birth: "",
    citizenship: "",
    gender: "",
    seria: "",
    number: "",
    passport_given: "",
    date_given: "",
    code: ""
  },
  reg_address: {
    mail_index: "",
    region: "",
    city: "",
    street: "",
    house: "",
    building: "",
    apartment: ""
  }
})



const errors = ref({
  second_name: false,
  first_name: false,
  phone: false,
  email: false,
  gender: false,
  seria: false,
  number: false,
  code: false,
  index: false
})

const onlyLetters = (obj, field) => {
  obj[field] = obj[field].replace(/[^А-Яа-яЁё\s-]/g, "")
  errors.value[field] = obj[field].length === 0
}

const formatPhone = () => {
  let digits = contractorForm.value.contractor.contact_phone.replace(/\D/g, "")
  if (!digits.startsWith("7")) digits = "7" + digits
  digits = digits.substring(0, 11)
  contractorForm.value.contractor.contact_phone = "+7" + digits.substring(1)
  errors.value.phone = digits.length !== 11
}

const validateEmail = () => {
  const email = contractorForm.value.contractor.email
  errors.value.email = !/^[\w.-]+@[\w.-]+\.\w+$/.test(email)
}

const formatCode = () => {
  let digits = contractorForm.value.passport.code.replace(/\D/g, "").substring(0, 6)
  if (digits.length >= 3) digits = digits.slice(0, 3) + "-" + digits.slice(3)
  contractorForm.value.passport.code = digits
  errors.value.code = digits.length !== 7
}

const digitsLimit = (obj, field, length) => {
  obj[field] = obj[field].replace(/\D/g, "").substring(0, length)
  errors.value[field === "mail_index" ? "index" : field] = obj[field].length !== length
}

const validateGender = () => {
  errors.value.gender = contractorForm.value.passport.gender === ""
}

const formValid = computed(() => {
  validateGender()
  return !Object.values(errors.value).includes(true)
})



const formatDate = (date) => {
  if (!date) return "—"
  const m = date.match(/^(\d{4})-(\d{2})-(\d{2})/)
  return m ? `${m[3]}.${m[2]}.${m[1]}` : date
}

const shortenName = (n) => (n?.length > 40 ? n.slice(0, 37) + "..." : n)

const loadDetails = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/listener/details/${id.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error("Ошибка загрузки данных")

    const data = await res.json()

    listener.value = data.data.listener || {}
    passport.value = data.data.passport || passport.value
    regaddress.value = data.data.regaddress || regaddress.value
    education.value = data.data.education_listener || education.value
    placework.value = data.data.placework || placework.value
    contractor.value = data.data.contractor || null

    if (contractor.value) hasContractor.value = true

    if (listener.value.snils) loadDocs(listener.value.snils)
    loadEnrollments()

  } catch (e) {
    toast.error(e.message)
  } finally {
    loading.value = false
  }
}

// const loadDocs = async (snils) => {
//   docsLoading.value = true
//   try {
//     const res = await fetch(`http://localhost:8082/v1/doc/exists?card-name=${snils}`, {
//       headers: { Authorization: `Bearer ${token}` }
//     })
//     docs.value = res.ok ? await res.json() : []
//   } finally {
//     docsLoading.value = false
//   }
// }

const downloadDoc = async (name) => {
  try {
    const res = await fetch(`${API_URL_DOC}/download?card-name=${name}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error()
    const blob = await res.blob()
    const url = URL.createObjectURL(blob)
    const a = document.createElement("a")
    a.href = url
    a.download = name
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    toast.error("Ошибка скачивания документа")
  }
}

const loadEnrollments = async () => {
  coursesLoading.value = true
  try {
    const res = await fetch(`${API_URL_CORE}/enrollment/details/${id.value}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    enrollments.value = res.ok ? (await res.json()).data : []
  } finally {
    coursesLoading.value = false
  }
}

const goBack = () => router.push("/listeners")
const goToCreateEnrollment = () => router.push(`/enrollment/create/${id.value}`)



const openContractorModal = () => contractorModal.value = true
const closeContractorModal = () => contractorModal.value = false

const saveContractor = async () => {
  if (!formValid.value) {
    toast.error("Исправьте ошибки перед сохранением")
    return
  }

  try {
    const res = await fetch(`${API_URL_CORE}/contractor/${id.value}`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(contractorForm.value)
    })

    if (!res.ok) throw new Error("Ошибка добавления заказчика")
    toast.success("Заказчик успешно добавлен!")
    contractorModal.value = false
    location.reload()
  } catch (e) {
    toast.error(e.message)
  }
}

onMounted(loadDetails)
</script>
