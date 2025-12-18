<template>
  <Header title="Запись на курс" />

  <div style="padding: 100px 20px 20px 20px; max-width: 1200px; margin: 0 auto;">
    <div v-if="loading" class="text-center py-5">Загрузка данных...</div>
    <div v-else>
      <p><strong>Слушатель:</strong> {{ listener.second_name }} {{ listener.first_name }} {{ listener.middle_name || '' }}</p>

      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">Заказчик</h5>
        <div class="d-flex align-items-center mb-3">
          <input
            type="checkbox"
            id="hasContractor"
            v-model="hasContractor"
            @change="onHasContractorChange"
            class="form-check-input me-2"
          />
          <label for="hasContractor" class="fw-bold">Есть заказчик</label>
        </div>

        <div v-if="contractor" class="mb-3">
          <h6 class="mb-2">Текущий заказчик</h6>
          <p><strong>ФИО:</strong> {{ contractor.contractor.second_name }} {{ contractor.contractor.first_name }} {{ contractor.contractor.middle_name }}</p>
          <p><strong>Телефон:</strong> {{ contractor.contractor.contact_phone }}</p>
          <p><strong>Email:</strong> {{ contractor.contractor.email }}</p>
        </div>

        <div class="d-flex gap-2 flex-wrap">
          <button type="button" @click="openContractorModal" class="btn btn-primary">
            {{ contractor ? 'Изменить заказчика' : 'Добавить заказчика' }}
          </button>
          <button
            v-if="contractor"
            type="button"
            class="btn btn-danger"
            @click="confirmDeleteContractor"
          >
            Удалить заказчика
          </button>
        </div>
      </div>

      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">Договор</h5>

        <div class="d-flex flex-column gap-3">
          <select
            v-model="selectedContractId"
            class="form-select"
            :disabled="!contracts.length || loadingContracts"
          >
            <option value="">— Выберите договор —</option>
            <option
              v-for="c in filteredContracts"
              :key="c.id_contract"
              :value="c.id_contract"
            >
              {{ c.name }} ({{ c.type === 'bilateral' ? 'двусторонний' : 'трёхсторонний' }})
            </option>
          </select>

          <select v-model="studyLoadOption" class="form-select">
            <option value="">— Недельная учебная нагрузка —</option>
            <option value="3">
              Недельная учебная нагрузка по настоящему договору составляет 3 академических часа в неделю, включая 2 академических часа взаимодействия 
              с преподавателем и 1 академический час самостоятельной работы; общая продолжительность освоения — 81 неделя.
            </option>
            <option value="6">
              Недельная учебная нагрузка по настоящему договору составляет 6 академических часов в неделю, включая 4 академических часа взаимодействия 
              с преподавателем и 2 академических часа самостоятельной работы; общая продолжительность освоения — 41 неделя.
            </option>
            <option value="12">
              Недельная учебная нагрузка по настоящему договору составляет 12 академических часов в неделю, включая 8 академических часов взаимодействия 
              с преподавателем и 4 академических часа самостоятельной работы; общая продолжительность освоения — 21 неделя.
            </option>
            <option value="15">
              Недельная учебная нагрузка по настоящему договору составляет 15 академических часов в неделю, включая 10 академических часов взаимодействия 
              с преподавателем и 5 академических часов самостоятельной работы; общая продолжительность освоения — 17 недель.
            </option>
            <option value="30">
              Недельная учебная нагрузка по настоящему договору составляет 30 академических часов в неделю, включая 20 академических часов взаимодействия 
              с преподавателем и 10 академических часов самостоятельной работы; общая продолжительность освоения — 9 недель.
            </option>
            <option value="32">
              Недельная учебная нагрузка по настоящему договору составляет 32 академических часа в неделю, включая 20 академических часов взаимодействия 
              с преподавателем и 12 академических часов самостоятельной работы; общая продолжительность освоения — 8 недель.
            </option>
          </select>

          <select v-model="paymentOption" class="form-select">
            <option value="">— Порядок оплаты —</option>
            <option value="full">
              Оплата осуществляется в следующем порядке: 100% предоплата до начала обучения.
            </option>
            <option value="split">
              Оплата осуществляется в следующем порядке: аванс 50 % — предоплата до начала обучения, оставшиеся 50 % — в установленный срок.
            </option>
          </select>

          <input
            v-if="paymentOption === 'split'"
            type="date"
            v-model="secondPaymentDate"
            class="form-control"
            placeholder="Срок оплаты оставшихся 50 %"
          />

          <select v-model="loadVariant" class="form-select">
            <option value="">— Вариант учебной нагрузки —</option>
            <option
              v-for="(label, key) in activeLoadVariants"
              :key="key"
              :value="key"
            >
              {{ label }}
            </option>
          </select>

          <select v-model="ageCategory" class="form-select">
            <option value="">— Возрастная категория —</option>
            <option
              v-for="(label, key) in ageCategories"
              :key="key"
              :value="key"
            >
              {{ label }}
            </option>
          </select>
        </div>
      </div>

      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">Исполнитель с стороны 25-12</h5>
        <div class="d-flex flex-column gap-2">
          <div
            v-for="e in executors"
            :key="e.id_executor"
            class="form-check"
          >
            <input
              class="form-check-input"
              type="radio"
              :value="e.id_executor"
              v-model="selectedExecutorId"
              :id="e.id_executor"
            />
            <label class="form-check-label" :for="e.id_executor">
              {{ e.second_name }} {{ e.first_name }} {{ e.middle_name }}
            </label>
          </div>
          <div v-if="!executors.length" class="text-muted">Исполнители не найдены</div>
        </div>
      </div>

      <div class="card p-3 shadow-sm card-block mt-4">
        <h5 class="card-title mb-3">Программа обучения</h5>
        <div class="d-flex flex-column gap-3">
          <select v-model="selectedProgramId" @change="onProgramChange" class="form-select">
            <option v-for="p in programs" :key="p.id_program_education" :value="p.id_program_education">
              {{ p.name_prof_education }}
            </option>
          </select>

          <div class="d-flex align-items-center gap-3 flex-wrap">
            <div>
              <label class="form-label">Дата начала</label>
              <input type="date" v-model="startDate" class="form-control" required />
            </div>
            <div>
              <label class="form-label">Дата окончания</label>
              <input type="date" v-model="endDate" class="form-control" required />
            </div>
            <div>
              <label class="form-label">Цена</label>
              <select v-model="currentPrice" class="form-select" required>
                <option v-if="price.individual_price !== undefined" :value="price.individual_price">Индивидуальное: {{ price.individual_price }} ₽</option>
                <option v-if="price.group_price !== undefined" :value="price.group_price">Групповое: {{ price.group_price }} ₽</option>
                <option v-if="price.campus_price !== undefined" :value="price.campus_price">Кампус: {{ price.campus_price }} ₽</option>
              </select>
            </div>
          </div>

          <div class="row g-2">
            <div class="col-md-6">
              <label class="form-label">Группа</label>
              <input type="text" v-model="group" class="form-control" placeholder="Например: 0" />
            </div>
            <div class="col-md-6">
              <label class="form-label">Тип обучения</label>
              <select v-model="typeOfRetraining" class="form-select">
                <option value="Повышение квалификации">Повышение квалификации</option>
                <option value="Профессиональная переподготовка">Профессиональная переподготовка</option>
                <option value="Дополнительное образование">Дополнительное образование</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-4 d-flex gap-2">
        <button class="btn btn-success px-4" @click="createEnrollment" :disabled="!isFormValid || saving">
          {{ saving ? 'Создание...' : 'Создать запись' }}
        </button>
        <button class="btn btn-secondary" @click="goBack">Назад</button>
      </div>
    </div>
  </div>

  <div class="modal fade show" v-if="contractorModal" style="display: block; background: rgba(0,0,0,.5); z-index: 1050;">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">{{ contractor ? 'Изменить заказчика' : 'Добавить заказчика' }}</h5>
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
                :class="['form-control', { 'is-invalid': errors.second_name }]"
              />
              <div v-if="errors.second_name" class="invalid-feedback d-block">Только буквы</div>
            </div>
            <div class="col-md-4">
              <label class="form-label">Имя</label>
              <input
                v-model="contractorForm.contractor.first_name"
                @input="onlyLetters(contractorForm.contractor, 'first_name')"
                :class="['form-control', { 'is-invalid': errors.first_name }]"
              />
              <div v-if="errors.first_name" class="invalid-feedback d-block">Только буквы</div>
            </div>
            <div class="col-md-4">
              <label class="form-label">Отчество</label>
              <input
                v-model="contractorForm.contractor.middle_name"
                @input="onlyLetters(contractorForm.contractor, 'middle_name')"
                class="form-control"
              />
            </div>

            <div class="col-md-6">
              <label class="form-label">Телефон</label>
              <input
                v-model="contractorForm.contractor.contact_phone"
                @input="formatPhone"
                :class="['form-control', { 'is-invalid': errors.phone }]"
              />
              <div v-if="errors.phone" class="invalid-feedback d-block">Формат: +7XXXXXXXXXX</div>
            </div>
            <div class="col-md-6">
              <label class="form-label">Email</label>
              <input
                v-model="contractorForm.contractor.email"
                @input="validateEmail"
                :class="['form-control', { 'is-invalid': errors.email }]"
              />
              <div v-if="errors.email" class="invalid-feedback d-block">Некорректный Email</div>
            </div>
          </div>

          <hr class="my-4">
          <h5 class="mb-3">Паспорт заказчика</h5>
          <div class="row g-3">
            <div class="col-md-6">
              <input v-model="contractorForm.passport.place_birth" class="form-control" placeholder="Место рождения" />
            </div>
            <div class="col-md-6">
              <input v-model="contractorForm.passport.citizenship" class="form-control" placeholder="Гражданство" />
            </div>
            <div class="col-md-4">
              <select
                v-model="contractorForm.passport.gender"
                :class="['form-select', { 'is-invalid': errors.gender }]"
              >
                <option value="">Пол</option>
                <option value="Мужской">Мужской</option>
                <option value="Женский">Женский</option>
              </select>
              <div v-if="errors.gender" class="invalid-feedback d-block">Выберите пол</div>
            </div>
            <div class="col-md-4">
              <input
                v-model="contractorForm.passport.seria"
                @input="digitsLimit(contractorForm.passport, 'seria', 4)"
                :class="['form-control', { 'is-invalid': errors.seria }]"
                placeholder="Серия (4 цифры)"
              />
              <div v-if="errors.seria" class="invalid-feedback d-block">4 цифры</div>
            </div>
            <div class="col-md-4">
              <input
                v-model="contractorForm.passport.number"
                @input="digitsLimit(contractorForm.passport, 'number', 6)"
                :class="['form-control', { 'is-invalid': errors.number }]"
                placeholder="Номер (6 цифр)"
              />
              <div v-if="errors.number" class="invalid-feedback d-block">6 цифр</div>
            </div>
            <div class="col-md-12">
              <input v-model="contractorForm.passport.passport_given" class="form-control" placeholder="Кем выдан" />
            </div>
            <div class="col-md-6">
              <input type="date" v-model="contractorForm.passport.date_given" class="form-control" />
            </div>
            <div class="col-md-6">
              <input
                v-model="contractorForm.passport.code"
                @input="formatCode"
                :class="['form-control', { 'is-invalid': errors.code }]"
                placeholder="Код подразделения (000-000)"
              />
              <div v-if="errors.code" class="invalid-feedback d-block">Формат: 000-000</div>
            </div>
          </div>

          <hr class="my-4">
          <h5 class="mb-3">Адрес регистрации</h5>
          <div class="row g-3">
            <div class="col-md-4">
              <input
                v-model="contractorForm.reg_address.mail_index"
                @input="digitsLimit(contractorForm.reg_address, 'mail_index', 6)"
                :class="['form-control', { 'is-invalid': errors.index }]"
                placeholder="Индекс (6 цифр)"
              />
              <div v-if="errors.index" class="invalid-feedback d-block">6 цифр</div>
            </div>
            <div class="col-md-4">
              <input v-model="contractorForm.reg_address.region" class="form-control" placeholder="Регион" />
            </div>
            <div class="col-md-4">
              <input v-model="contractorForm.reg_address.city" class="form-control" placeholder="Город" />
            </div>
            <div class="col-md-6">
              <input v-model="contractorForm.reg_address.street" class="form-control" placeholder="Улица" />
            </div>
            <div class="col-md-3">
              <input v-model="contractorForm.reg_address.house" class="form-control" placeholder="Дом" />
            </div>
            <div class="col-md-3">
              <input v-model="contractorForm.reg_address.building" class="form-control" placeholder="Корпус" />
            </div>
            <div class="col-md-3">
              <input v-model="contractorForm.reg_address.apartment" class="form-control" placeholder="Квартира" />
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button @click="closeContractorModal" class="btn btn-secondary">Отмена</button>
          <button @click="saveContractor" :disabled="!formValid" class="btn btn-success">
            {{ contractor ? 'Сохранить изменения' : 'Добавить' }}
          </button>
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

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { toast } from 'vue3-toastify'
import Header from './Header.vue'
import ConfirmModal from './ConfirmModal.vue'
import { API_URL_CORE } from '../config'

const route = useRoute()
const router = useRouter()
const listenerId = route.params.listenerId
const token = localStorage.getItem('access_token')

const listener = ref({})
const contractor = ref(null)
const hasContractor = ref(false)
const loading = ref(true)

const ageCategory = ref('')

const documentType = ref('')
const loadVariant = ref('')

const loadVariantsDO = {
  1: 'с пониженной недельной учебной нагрузкой (1 акад. час в неделю)',
  2: 'с умеренной недельной учебной нагрузкой (2 акад. часа в неделю)',
  3: 'со стандартной недельной учебной нагрузкой (3 акад. часа в неделю)',
  4: 'с высокой недельной учебной нагрузкой (4 акад. часа в неделю)',
  5: 'с повышенной недельной учебной нагрузкой (6 акад. часов в неделю)'
}

const loadVariantsNotDO = {
  1: 'с пониженной недельной учебной нагрузкой (3 акад. часа в неделю)',
  2: 'с умеренной недельной учебной нагрузкой (6 акад. часов в неделю)',
  3: 'со стандартной недельной учебной нагрузкой (12 акад. часов в неделю)',
  4: 'с высокой недельной учебной нагрузкой (15 акад. часов в неделю)',
  5: 'с повышенной недельной учебной нагрузкой (30 акад. часов в неделю)',
  6: 'с интенсивной недельной учебной нагрузкой (36 акад. часов в неделю)'
}

const ageCategories = {
  belowEighteen: 'Меньше восемнадцати',
  belowFourteen: 'Меньше четырнадцати',
  eighteen: 'Восемнадцать'
}

const contracts = ref([])
const loadingContracts = ref(false)
const selectedContractId = ref('')

const studyLoadOption = ref('')
const paymentOption = ref('')
const secondPaymentDate = ref('')

const executors = ref([])
const selectedExecutorId = ref('')

const programs = ref([])
const selectedProgramId = ref('')
const price = ref({ individual_price: 0, group_price: 0, campus_price: 0 })
const currentPrice = ref(0)
const startDate = ref('')
const endDate = ref('')
const group = ref('')
const typeOfRetraining = ref('')
const saving = ref(false)

const contractorModal = ref(false)
const confirmModal = ref(null)

const contractorForm = ref({
  contractor: { first_name: "", second_name: "", middle_name: "", contact_phone: "", email: "" },
  passport: { place_birth: "", citizenship: "", gender: "", seria: "", number: "", passport_given: "", date_given: "", code: "" },
  reg_address: { mail_index: "", region: "", city: "", street: "", house: "", building: "", apartment: "" }
})

const errors = ref({
  second_name: false, first_name: false, phone: false, email: false,
  gender: false, seria: false, number: false, code: false, index: false
})

const onlyLetters = (obj, field) => {
  obj[field] = obj[field].replace(/[^А-Яа-яЁёA-Za-z\s-]/g, "")
  errors.value[field] = obj[field].trim() === ""
}
const formatPhone = () => {
  let v = contractorForm.value.contractor.contact_phone.replace(/\D/g, "")
  if (!v.startsWith("7")) v = "7" + v
  v = v.slice(0, 11)
  contractorForm.value.contractor.contact_phone = v ? `+${v}` : ""
  errors.value.phone = v.length !== 11
}
const validateEmail = () => {
  const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  errors.value.email = !re.test(contractorForm.value.contractor.email)
}
const formatCode = () => {
  let v = contractorForm.value.passport.code.replace(/\D/g, "").slice(0, 6)
  if (v.length >= 3) v = v.slice(0, 3) + "-" + v.slice(3)
  contractorForm.value.passport.code = v
  errors.value.code = v.length !== 7 && v.length !== 0
}
const digitsLimit = (obj, field, len) => {
  obj[field] = obj[field].replace(/\D/g, "").slice(0, len)
  errors.value[field === "mail_index" ? "index" : field] = obj[field].length !== len && obj[field] !== ""
}
const validateGender = () => {
  errors.value.gender = contractorForm.value.passport.gender === ""
}
const formValid = computed(() => {
  validateGender()
  return !Object.values(errors.value).some(Boolean)
})

const selectedContract = computed(() =>
  contracts.value.find(c => c.id_contract === selectedContractId.value)
)

const filteredContracts = computed(() => {
  if (hasContractor.value) {
    return contracts.value.filter(c => c.type === 'trilateral')
  } else {
    return contracts.value.filter(c => c.type === 'bilateral')
  }
})

const isDoFiz = computed(() => {
  const c = contracts.value.find(c => c.id_contract === selectedContractId.value)
  return c?.name === 'ДО с оплатой физическим лицом'
})

const activeLoadVariants = computed(() => {
  return isDoFiz.value ? loadVariantsDO : loadVariantsNotDO
})

const isFormValid = computed(() => {
  return (
    selectedProgramId.value &&
    startDate.value &&
    endDate.value &&
    currentPrice.value != null &&
    (!hasContractor.value || selectedContractId.value)
  )
})

const loadListenerAndContractor = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/listener/details/${listenerId}`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error(`Слушатель не найден (${res.status})`)
    const data = await res.json()
    listener.value = data.data.listener || {}
    contractor.value = data.data.contractor || null
    hasContractor.value = !!contractor.value
  } catch (err) {
    toast.error(err.message || 'Ошибка загрузки слушателя')
  }
}

const loadPrograms = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/programeducation/?page=1&filter=`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    const data = await res.json()
    programs.value = data.data || []
    if (programs.value.length && !selectedProgramId.value) {
      selectedProgramId.value = programs.value[0].id_program_education
      onProgramChange()
    }
  } catch (err) {
    toast.error('Не удалось загрузить программы')
  }
}

const loadContracts = async () => {
  loadingContracts.value = true
  contracts.value = [
    { id_contract: 'DO-FIZ-3', name: 'ДО с оплатой физическим лицом', type: 'trilateral' },
    { id_contract: 'PK-FIZ-2', name: 'ПК с оплатой физическим лицом', type: 'bilateral' },
    { id_contract: 'PK-FIZ-3', name: 'ПК с оплатой физическим лицом', type: 'trilateral' },
    { id_contract: 'PK-YUR-3', name: 'ПК с оплатой юридическим лицом', type: 'trilateral' },
    { id_contract: 'PP-FIZ-2', name: 'ПП с оплатой физическим лицом', type: 'bilateral' },
    { id_contract: 'PP-FIZ-3', name: 'ПП с оплатой физическим лицом', type: 'trilateral' },
    { id_contract: 'PP-YUR-3', name: 'ПП с оплатой юридическим лицом', type: 'trilateral' }
  ]
  loadingContracts.value = false
}

const loadExecutors = async () => {
  try {
    const res = await fetch(`${API_URL_CORE}/executer/`, {
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка загрузки исполнителей')
    const data = await res.json()
    executors.value = data.data || []
  } catch (err) {
    toast.error(err.message || 'Не удалось загрузить исполнителей')
  }
}

const onProgramChange = () => {
  const p = programs.value.find(prog => prog.id_program_education === selectedProgramId.value)
  if (p) {
    price.value = {
      individual_price: p.individual_price,
      group_price: p.group_price,
      campus_price: p.campus_price
    }
    currentPrice.value = p.individual_price
  }
}

const openContractorModal = () => {
  if (contractor.value) {
    contractorForm.value = {
      contractor: { ...contractor.value.contractor },
      passport: { ...contractor.value.passport },
      reg_address: { ...contractor.value.reg_address }
    }
  } else {
    contractorForm.value = {
      contractor: { first_name: "", second_name: "", middle_name: "", contact_phone: "", email: "" },
      passport: { place_birth: "", citizenship: "", gender: "", seria: "", number: "", passport_given: "", date_given: "", code: "" },
      reg_address: { mail_index: "", region: "", city: "", street: "", house: "", building: "", apartment: "" }
    }
  }
  Object.keys(errors.value).forEach(k => errors.value[k] = false)
  contractorModal.value = true
}

const closeContractorModal = () => {
  contractorModal.value = false
}

const onHasContractorChange = () => {
  if (hasContractor.value) {
    if (!contractor.value) {
      openContractorModal()
    }
  } else {
    if (contractor.value) {
      selectedContractId.value = ''
      confirmDeleteContractor()
    } else {
      selectedContractId.value = ''
    }
  }
}

const saveContractor = async () => {
  if (!formValid.value) return toast.error("Исправьте ошибки в форме")

  try {
    saving.value = true

    const url = `${API_URL_CORE}/contractor/${listenerId}`

    const res = await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`
      },
      body: JSON.stringify(contractorForm.value)
    })

    if (!res.ok) {
      const e = await res.json().catch(() => ({}))
      throw new Error(e.message || "Ошибка добавления заказчика")
    }

    toast.success("Заказчик успешно добавлен")

    await loadListenerAndContractor()
    closeContractorModal()

  } catch (err) {
    toast.error(err.message || "Не удалось добавить заказчика")
  } finally {
    saving.value = false
  }
}

const confirmDeleteContractor = () => {
  confirmModal.value?.open(async () => {
    await deleteContractor()
  })
}

const deleteContractor = async () => {
  try {
    saving.value = true
    const id = contractor.value.contractor.id_contractor
    const res = await fetch(`${API_URL_CORE}/contractor/${id}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` }
    })
    if (!res.ok) throw new Error('Ошибка удаления заказчика(заказчик игнорируется)')
    toast.success('Заказчик удалён')
    await loadListenerAndContractor()
  } catch (err) {
    toast.info(err.message || 'Не удалось удалить заказчика')
  } finally {
    saving.value = false
  }
}

const createEnrollment = async () => {
  if (!isFormValid.value) return toast.warn('Заполните обязательные поля')

  try {
    saving.value = true

    const body = {
      id_listener: listenerId,
      id_program: selectedProgramId.value,
      start_date: startDate.value,
      end_date: endDate.value,
      current_price: Number(currentPrice.value),
      group: group.value || null,
      type_of_retraining: typeOfRetraining.value || null,
      is_active: true
    }

    const res = await fetch(`${API_URL_CORE}/enrollment/`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
      body: JSON.stringify(body)
    })

    if (!res.ok) {
      const e = await res.json().catch(() => ({}))
      throw new Error(e.message || `Ошибка создания записи, возможно такая запись на курс уже существует (${res.status})`)
    }

    toast.success('Запись успешно создана')

    const frontData = {
      variant: loadVariant.value ? Number(loadVariant.value) : null,
      dogovor_type: selectedContractId.value || null,
      opion_nagruz: studyLoadOption.value || null,
      opt_document: ageCategory.value || null,
      opt_price: currentPrice.value != null ? String(currentPrice.value) : null
    }

    if (paymentOption.value === 'split' && secondPaymentDate.value) {
      frontData.second_payment_date = secondPaymentDate.value
    }

    const documentPayload = {
      id_listener: listenerId,
      id_program: selectedProgramId.value,
      id_executor: selectedExecutorId.value || null,
      FrontData: frontData
    }

    try {
      const docRes = await fetch(`${API_URL_CORE}/document/`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify(documentPayload)
      })
      if (!docRes.ok) {
        const e = await docRes.json().catch(() => ({}))
        toast.error(e.message || 'Не удалось создать документ')
      } else {
        toast.success('Документ создан')
      }
    } catch (e) {
      toast.error('Ошибка при создании документа')
    }

    router.push(`/listeners/${listenerId}`)
  } catch (err) {
    toast.error(err.message || 'Не удалось создать запись, возможно такая запись на курс уже существует')
  } finally {
    saving.value = false
  }
}

const goBack = () => router.push(`/listeners/${listenerId}`)

onMounted(async () => {
  await loadListenerAndContractor()
  await loadPrograms()
  await loadContracts()
  await loadExecutors()
  loading.value = false
})
</script>

<style scoped>
.card-block {
  border-radius: 8px;
}
.form-label {
  font-weight: 500;
}
.invalid-feedback {
  font-size: 0.875em;
}
.table-hover tbody tr:hover {
  background-color: #e2f0d9;
  cursor: pointer;
}
th, td {
  vertical-align: middle;
}
</style>
