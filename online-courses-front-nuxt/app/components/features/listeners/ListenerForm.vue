<script setup lang="ts">
import type {
  LevelEducationItem,
  ListenerFormPayload,
  ListenerFormState
} from '../../../types/listener'
import { createEmptyListenerFormState } from '../../../types/listener'
import { normalizeApiDate } from '../../../utils/date'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppCheckbox from '../../ui/AppCheckbox.vue'
import AppInput from '../../ui/AppInput.vue'
import AppSelect from '../../ui/AppSelect.vue'

const props = withDefaults(defineProps<{
  title: string
  submitLabel: string
  loading?: boolean
  initialState?: ListenerFormState | null
  mode?: 'create' | 'edit'
}>(), {
  loading: false,
  initialState: null,
  mode: 'create'
})

const emit = defineEmits<{
  submit: [payload: ListenerFormPayload]
  cancel: []
}>()

const state = reactive<ListenerFormState>(createEmptyListenerFormState())
const levels = ref<LevelEducationItem[]>([])
const levelsLoading = ref(false)
const firstNameError = ref('')
const secondNameError = ref('')
const middleNameError = ref('')
const emailError = ref('')
const dateOfBirthError = ref('')
const snilsError = ref('')
const phoneError = ref('')
const passportSeriaError = ref('')
const passportNumberError = ref('')
const passportCodeError = ref('')
const registrationIndexError = ref('')
const registrationRegionError = ref('')
const registrationCityError = ref('')
const registrationStreetError = ref('')
const registrationHouseError = ref('')
const registrationBuildingError = ref('')
const registrationApartmentError = ref('')
const educationCityError = ref('')
const educationRegionError = ref('')
const diplomaSeriaError = ref('')
const diplomaNumberError = ref('')
const workExperienceError = ref('')
const workRoleExperienceError = ref('')
const submitError = ref('')
const notifications = useNotifications()

const normalizeDateValue = (value?: string | null) => normalizeApiDate(value)

const educationDisabled = computed(() => state.listener.looting_education)
const namePattern = /[^А-Яа-яЁёA-Za-z.\s]/g
const emailPattern = /[^A-Za-z0-9@._-]/g
const addressPattern = /[^0-9А-Яа-яЁёA-Za-z.,\-\s]/g
const addressNumberPattern = /[^0-9А-Яа-яЁёA-Za-z.,\-№#\s]/g

const validateRequiredText = (value: string, label: string, required = true) => {
  if (!value.trim()) {
    return required ? `Укажите ${label.toLowerCase()}` : ''
  }

  return ''
}

const validateEmail = () => {
  if (!state.listener.email) {
    emailError.value = 'Укажите email'
    return false
  }

  const valid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(state.listener.email)
  emailError.value = valid ? '' : 'Некорректный email'
  return valid
}

const validateDateOfBirth = () => {
  dateOfBirthError.value = state.listener.date_of_birth ? '' : 'Выберите дату: 02.01.2006'
  return !dateOfBirthError.value
}

const validateSnils = () => {
  if (!state.listener.snils) {
    snilsError.value = 'Укажите СНИЛС'
    return false
  }

  const valid = /^\d{3}-\d{3}-\d{3} \d{2}$/.test(state.listener.snils)
  snilsError.value = valid ? '' : 'СНИЛС должен быть в формате 123-456-789 00'
  return valid
}

const validatePhone = () => {
  const digits = state.listener.contact_phone.replace(/\D/g, '')

  if (!digits || digits === '7') {
    phoneError.value = 'Укажите телефон'
    return false
  }

  const valid = digits.length === 11
  phoneError.value = valid ? '' : 'Телефон должен быть в формате +7XXXXXXXXXX'
  return valid
}

const validatePassportSeria = () => {
  if (!state.passport.seria) {
    passportSeriaError.value = hasAnyValue(state.passport) ? 'Серия паспорта обязательна' : ''
    return !passportSeriaError.value
  }

  const valid = /^\d{4}$/.test(state.passport.seria)
  passportSeriaError.value = valid ? '' : 'Серия паспорта: 4 цифры'
  return valid
}

const validatePassportNumber = () => {
  if (!state.passport.number) {
    passportNumberError.value = hasAnyValue(state.passport) ? 'Номер паспорта обязателен' : ''
    return !passportNumberError.value
  }

  const valid = /^\d{6}$/.test(state.passport.number)
  passportNumberError.value = valid ? '' : 'Номер паспорта: 6 цифр'
  return valid
}

const validatePassportCode = () => {
  if (!state.passport.code) {
    passportCodeError.value = hasAnyValue(state.passport) ? 'Код подразделения обязателен' : ''
    return !passportCodeError.value
  }

  const valid = /^\d{3}-\d{3}$/.test(state.passport.code)
  passportCodeError.value = valid ? '' : 'Код подразделения: 123-456'
  return valid
}

const validateRegistrationIndex = () => {
  if (!state.registrationAddress.mail_index) {
    registrationIndexError.value = 'Укажите индекс'
    return false
  }

  const valid = /^\d{6}$/.test(state.registrationAddress.mail_index)
  registrationIndexError.value = valid ? '' : 'Индекс должен содержать 6 цифр'
  return valid
}

const validateDiplomaSeria = () => {
  if (!state.education.diplom_seria) {
    diplomaSeriaError.value = hasAnyValue(state.education) ? 'Серия диплома обязательна' : ''
    return !diplomaSeriaError.value
  }

  const valid = /^\d{6}$/.test(state.education.diplom_seria)
  diplomaSeriaError.value = valid ? '' : 'Серия диплома: 6 цифр'
  return valid
}

const validateDiplomaNumber = () => {
  if (!state.education.diplom_number) {
    diplomaNumberError.value = hasAnyValue(state.education) ? 'Номер диплома обязателен' : ''
    return !diplomaNumberError.value
  }

  const valid = /^\d{7}$/.test(state.education.diplom_number)
  diplomaNumberError.value = valid ? '' : 'Номер диплома: 7 цифр'
  return valid
}

const validateWorkExperience = () => {
  if (!state.placeWork.all_experience) {
    workExperienceError.value = hasAnyValue(state.placeWork) ? 'Укажите общий стаж' : ''
    return !workExperienceError.value
  }

  const valid = /^\d{1,2}$/.test(String(state.placeWork.all_experience))
  workExperienceError.value = valid ? '' : 'Общий стаж: только числа, до 2 цифр'
  return valid
}

const validateWorkRoleExperience = () => {
  if (!state.placeWork.job_title_expirience) {
    workRoleExperienceError.value = hasAnyValue(state.placeWork) ? 'Укажите стаж по должности' : ''
    return !workRoleExperienceError.value
  }

  const valid = /^\d{1,2}$/.test(String(state.placeWork.job_title_expirience))
  workRoleExperienceError.value = valid ? '' : 'Стаж по должности: только числа, до 2 цифр'
  return valid
}

const validateListenerName = (
  field: 'first_name' | 'second_name' | 'middle_name',
  label: string,
  required = true
) => {
  const raw = state.listener[field]
  const sanitized = raw.replace(namePattern, '')
  const hadInvalidChars = raw !== sanitized

  state.listener[field] = sanitized

  const targetError =
    field === 'first_name'
      ? firstNameError
      : field === 'second_name'
        ? secondNameError
        : middleNameError

  if (hadInvalidChars) {
    targetError.value = 'Допустимы только буквы и точка'
    return false
  }

  targetError.value = validateRequiredText(sanitized, label, required)
  return !targetError.value
}

const validateAddressField = (
  field: keyof ListenerFormState['registrationAddress'],
  label: string,
  allowNumbers = false
) => {
  const raw = state.registrationAddress[field]
  const pattern = allowNumbers ? addressNumberPattern : addressPattern
  const sanitized = raw.replace(pattern, '')
  const hadInvalidChars = raw !== sanitized

  state.registrationAddress[field] = sanitized

  const targetError =
    field === 'region'
      ? registrationRegionError
      : field === 'city'
        ? registrationCityError
        : field === 'street'
          ? registrationStreetError
          : field === 'house'
            ? registrationHouseError
            : field === 'building'
              ? registrationBuildingError
              : registrationApartmentError

  if (hadInvalidChars) {
    targetError.value = allowNumbers
      ? 'Допустимы буквы, цифры, точки, запятые, тире, № и #'
      : 'Допустимы буквы, цифры, точки, запятые и тире'
    return false
  }

  targetError.value = validateRequiredText(sanitized, label)
  return !targetError.value
}

const validateEducationAddressField = (field: 'region' | 'city', label: string) => {
  const raw = state.education[field]
  const sanitized = raw.replace(addressPattern, '')
  const hadInvalidChars = raw !== sanitized

  state.education[field] = sanitized
  const targetError = field === 'region' ? educationRegionError : educationCityError

  if (hadInvalidChars) {
    targetError.value = 'Допустимы буквы, цифры, точки, запятые и тире'
    return false
  }

  targetError.value = hasAnyValue(state.education)
    ? validateRequiredText(sanitized, label)
    : ''
  return !targetError.value
}

const numbersOnly = (
  section: 'passport' | 'education' | 'registrationAddress',
  field: string,
  max: number
) => {
  const target = state[section] as Record<string, string>
  target[field] = (target[field] || '').replace(/\D/g, '').slice(0, max)
}

const handlePassportDigitsInput = (field: 'seria' | 'number') => {
  numbersOnly('passport', field, field === 'seria' ? 4 : 6)

  if (field === 'seria') {
    validatePassportSeria()
    return
  }

  validatePassportNumber()
}

const handleEducationDigitsInput = (field: 'diplom_seria' | 'diplom_number') => {
  numbersOnly('education', field, field === 'diplom_seria' ? 6 : 7)

  if (field === 'diplom_seria') {
    validateDiplomaSeria()
    return
  }

  validateDiplomaNumber()
}

const handleRegistrationIndexInput = () => {
  numbersOnly('registrationAddress', 'mail_index', 6)
  validateRegistrationIndex()
}

const formatSnils = () => {
  const raw = state.listener.snils
  let value = raw.replace(/\D/g, '').slice(0, 11)

  if (value.length > 9) {
    value = `${value.slice(0, 3)}-${value.slice(3, 6)}-${value.slice(6, 9)} ${value.slice(9)}`
  } else if (value.length > 6) {
    value = `${value.slice(0, 3)}-${value.slice(3, 6)}-${value.slice(6)}`
  } else if (value.length > 3) {
    value = `${value.slice(0, 3)}-${value.slice(3)}`
  }

  state.listener.snils = value
  if (/[^\d\s-]/.test(raw)) {
    snilsError.value = 'Допустимы только цифры'
    return
  }

  validateSnils()
}

const formatPhone = () => {
  const raw = state.listener.contact_phone
  let value = raw.replace(/\D/g, '')

  if (!value.startsWith('7')) {
    value = `7${value}`
  }

  state.listener.contact_phone = `+${value.slice(0, 11)}`

  if (/[^\d+]/.test(raw)) {
    phoneError.value = 'Допустимы только цифры'
    return
  }

  validatePhone()
}

const formatPassportCode = () => {
  const raw = state.passport.code
  let value = raw.replace(/\D/g, '').slice(0, 6)

  if (value.length > 3) {
    value = `${value.slice(0, 3)}-${value.slice(3)}`
  }

  state.passport.code = value

  if (/[^\d-]/.test(raw)) {
    passportCodeError.value = 'Допустимы только цифры'
    return
  }

  validatePassportCode()
}

const resetEducation = () => {
  Object.assign(state.education, createEmptyListenerFormState().education)
}

const numbersOnlyWork = (field: 'all_experience' | 'job_title_expirience') => {
  const raw = state.placeWork[field]

  state.placeWork[field] = raw.replace(/\D/g, '').slice(0, 2)

  if (field === 'all_experience') {
    if (/\D/.test(raw)) {
      workExperienceError.value = 'Допустимы только цифры'
      return
    }

    validateWorkExperience()
    return
  }

  if (/\D/.test(raw)) {
    workRoleExperienceError.value = 'Допустимы только цифры'
    return
  }

  validateWorkRoleExperience()
}

const createEmptyEducationPayload = () => ({
  diplom_seria: '',
  diplom_number: '',
  date_given: '',
  city: '',
  region: '',
  educational_institution: '',
  speciality: '',
  level_education: ''
})

watch(
  () => state.listener.looting_education,
  (value) => {
    if (value) {
      resetEducation()
    }
  }
)

watch(levels, () => {
  syncEducationLevelWithOptions()
})

const hasAnyValue = (values: Record<string, string | number>) =>
  Object.values(values).some((value) => String(value || '').trim().length > 0)

const syncEducationLevelWithOptions = () => {
  if (!state.education.level_education) {
    return
  }

  const directMatch = levels.value.find(
    (level) => String(level.id_level_education) === String(state.education.level_education)
  )

  if (directMatch) {
    state.education.level_education = directMatch.id_level_education
    return
  }

  const byName = levels.value.find(
    (level) => level.education === state.education.level_education
  )

  if (byName) {
    state.education.level_education = byName.id_level_education
  }
}

const syncState = (nextState?: ListenerFormState | null) => {
  const source = nextState || createEmptyListenerFormState()

  Object.assign(state.listener, source.listener, {
    date_of_birth: normalizeDateValue(source.listener.date_of_birth)
  })
  Object.assign(state.passport, source.passport, {
    date_given: normalizeDateValue(source.passport.date_given)
  })
  Object.assign(state.registrationAddress, source.registrationAddress)
  Object.assign(state.education, source.education, {
    date_given: normalizeDateValue(source.education.date_given)
  })
  Object.assign(state.placeWork, source.placeWork, {
    all_experience: String(source.placeWork.all_experience || ''),
    job_title_expirience: String(source.placeWork.job_title_expirience || '')
  })
}

watch(
  () => props.initialState,
  (nextState) => {
    syncState(nextState)
    syncEducationLevelWithOptions()
  },
  { immediate: true, deep: true }
)

const validateRequired = () => {
  const isMainInfoValid = [
    validateListenerName('first_name', 'Имя'),
    validateListenerName('second_name', 'Фамилия'),
    validateListenerName('middle_name', 'Отчество', false),
    validateDateOfBirth(),
    validateSnils(),
    validatePhone(),
    validateEmail()
  ].every(Boolean)

  const isRegistrationValid = [
    validateRegistrationIndex(),
    validateAddressField('region', 'Регион'),
    validateAddressField('city', 'Город'),
    validateAddressField('street', 'Улица'),
    validateAddressField('house', 'Дом', true),
    validateAddressField('building', 'Корпус', true),
    validateAddressField('apartment', 'Квартира', true)
  ].every(Boolean)

  if (!isMainInfoValid || !isRegistrationValid) {
    return false
  }

  if (
    !validatePassportSeria() ||
    !validatePassportNumber() ||
    !validatePassportCode() ||
    !validateWorkExperience() ||
    !validateWorkRoleExperience()
  ) {
    return false
  }

  const requiredListenerFields = [
    state.listener.first_name,
    state.listener.second_name,
    state.listener.date_of_birth,
    state.listener.snils,
    state.listener.contact_phone,
    state.listener.email
  ]

  const requiredAddressFields = [
    state.registrationAddress.mail_index,
    state.registrationAddress.region,
    state.registrationAddress.city,
    state.registrationAddress.street,
    state.registrationAddress.house,
    state.registrationAddress.building,
    state.registrationAddress.apartment
  ]

  if (requiredListenerFields.some((value) => !String(value).trim())) {
    submitError.value = 'Заполните обязательные поля слушателя'
    return false
  }

  if (requiredAddressFields.some((value) => !String(value).trim())) {
    submitError.value = 'Заполните адрес регистрации'
    return false
  }

  if (hasAnyValue(state.passport) && Object.values(state.passport).some((value) => !String(value).trim())) {
    submitError.value = 'Если заполняете паспорт, заполните его полностью'
    return false
  }

  if (hasAnyValue(state.passport) && (!validatePassportSeria() || !validatePassportNumber() || !validatePassportCode())) {
    submitError.value = 'Проверьте формат паспортных данных'
    return false
  }

  if (!state.listener.looting_education) {
    const educationValues = Object.values(state.education)
    const hasPartialEducation = educationValues.some((value) => String(value).trim())

    if (hasPartialEducation && (!validateDiplomaSeria() || !validateDiplomaNumber())) {
      submitError.value = 'Проверьте серию и номер диплома'
      return false
    }

    if (hasPartialEducation && educationValues.some((value) => !String(value).trim())) {
      submitError.value = 'Если заполняете образование, заполните блок полностью'
      return false
    }
  }

  if (hasAnyValue(state.placeWork) && Object.values(state.placeWork).some((value) => !String(value).trim())) {
    submitError.value = 'Если заполняете место работы, заполните блок полностью'
    return false
  }

  if (hasAnyValue(state.placeWork) && (!validateWorkExperience() || !validateWorkRoleExperience())) {
    submitError.value = 'Стаж работы нужно указывать только цифрами'
    return false
  }

  submitError.value = ''
  return true
}

const buildPayload = (): ListenerFormPayload => {
  const placeWorkPayload = hasAnyValue(state.placeWork)
    ? {
        ...state.placeWork,
        all_experience: parseInt(String(state.placeWork.all_experience), 10),
        job_title_expirience: parseInt(String(state.placeWork.job_title_expirience), 10)
      }
    : undefined

  const payload: ListenerFormPayload = {
    listener: {
      ...state.listener,
      looting_education: Boolean(state.listener.looting_education)
    },
    registration_address: {
      ...state.registrationAddress,
      mail_index: String(state.registrationAddress.mail_index || '')
    }
  }

  if (props.mode === 'edit' || hasAnyValue(state.passport)) {
    payload.passport = {
      ...state.passport,
      seria: String(state.passport.seria || ''),
      number: String(state.passport.number || '')
    }
  }

  if (state.listener.looting_education) {
    if (props.mode === 'edit') {
      payload.education = createEmptyEducationPayload()
    }
  } else if (props.mode === 'edit' || hasAnyValue(state.education)) {
    payload.education = {
      ...state.education,
      diplom_seria: String(state.education.diplom_seria || ''),
      diplom_number: String(state.education.diplom_number || ''),
      level_education: state.education.level_education || ''
    }
  }

  if (placeWorkPayload) {
    payload.placeWork = placeWorkPayload
  }

  return payload
}

const submit = async () => {
  if (!validateRequired()) {
    return
  }

  emit('submit', buildPayload())
}

onMounted(async () => {
  levelsLoading.value = true

  try {
    levels.value = await getEducationLevels()
    syncEducationLevelWithOptions()
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось загрузить уровни образования',
      'Слушатели'
    )
  } finally {
    levelsLoading.value = false
  }
})
</script>

<template>
  <form class="listener-form stack" @submit.prevent="submit">
    <div class="listener-form__grid listener-form__grid--triple">
      <AppCard title="Основная информация">
        <div class="form-stack">
          <AppInput
            v-model="state.listener.first_name"
            label="Имя"
            placeholder="Введите имя"
            :error="firstNameError"
            @update:model-value="validateListenerName('first_name', 'Имя')"
          />
          <AppInput
            v-model="state.listener.second_name"
            label="Фамилия"
            placeholder="Введите фамилию"
            :error="secondNameError"
            @update:model-value="validateListenerName('second_name', 'Фамилия')"
          />
          <AppInput
            v-model="state.listener.middle_name"
            label="Отчество"
            placeholder="Введите отчество"
            :error="middleNameError"
            @update:model-value="validateListenerName('middle_name', 'Отчество', false)"
          />
          <AppInput
            v-model="state.listener.date_of_birth"
            label="Дата рождения"
            type="date"
            :error="dateOfBirthError"
            @update:model-value="validateDateOfBirth"
          />
          <AppInput
            v-model="state.listener.snils"
            label="СНИЛС"
            placeholder="xxx-xxx-xxx xx"
            maxlength="14"
            :error="snilsError"
            @update:model-value="formatSnils"
          />
          <AppInput
            v-model="state.listener.contact_phone"
            label="Телефон"
            placeholder="+7XXXXXXXXXX"
            maxlength="12"
            :error="phoneError"
            @update:model-value="formatPhone"
          />
          <AppInput
            v-model="state.listener.email"
            label="Email"
            placeholder="example@mail.ru"
            :error="emailError"
            @update:model-value="validateEmail"
          />
        </div>
      </AppCard>

      <AppCard title="Паспорт">
        <div class="form-stack">
          <AppInput v-model="state.passport.place_birth" label="Место рождения" placeholder="Укажите место рождения" />
          <AppInput v-model="state.passport.citizenship" label="Гражданство" placeholder="Укажите гражданство" />

          <AppSelect
            v-model="state.passport.gender"
            label="Пол"
            placeholder="Выберите пол"
          >
            <option value="Мужской">Мужской</option>
            <option value="Женский">Женский</option>
          </AppSelect>

          <AppInput
            v-model="state.passport.seria"
            label="Серия"
            placeholder="4 цифры"
            maxlength="4"
            inputmode="numeric"
            :error="passportSeriaError"
            @update:model-value="handlePassportDigitsInput('seria')"
          />
          <AppInput
            v-model="state.passport.number"
            label="Номер"
            placeholder="6 цифр"
            maxlength="6"
            inputmode="numeric"
            :error="passportNumberError"
            @update:model-value="handlePassportDigitsInput('number')"
          />
          <AppInput v-model="state.passport.passport_given" label="Кем выдан" placeholder="Укажите орган выдачи" />
          <AppInput
            v-model="state.passport.date_given"
            label="Дата выдачи"
            type="date"
          />
          <AppInput
            v-model="state.passport.code"
            label="Код подразделения"
            placeholder="xxx-xxx"
            maxlength="7"
            inputmode="numeric"
            :error="passportCodeError"
            @update:model-value="formatPassportCode"
          />
        </div>
      </AppCard>

      <AppCard title="Адрес регистрации">
        <div class="form-stack">
          <AppInput
            v-model="state.registrationAddress.mail_index"
            label="Индекс"
            placeholder="6 цифр"
            maxlength="6"
            inputmode="numeric"
            :error="registrationIndexError"
            @update:model-value="handleRegistrationIndexInput"
          />
          <AppInput
            v-model="state.registrationAddress.region"
            label="Регион"
            placeholder="Укажите регион"
            :error="registrationRegionError"
            @update:model-value="validateAddressField('region', 'Регион')"
          />
          <AppInput
            v-model="state.registrationAddress.city"
            label="Город"
            placeholder="Укажите город"
            :error="registrationCityError"
            @update:model-value="validateAddressField('city', 'Город')"
          />
          <AppInput
            v-model="state.registrationAddress.street"
            label="Улица"
            placeholder="Укажите улицу"
            :error="registrationStreetError"
            @update:model-value="validateAddressField('street', 'Улица')"
          />
          <AppInput
            v-model="state.registrationAddress.house"
            label="Дом"
            placeholder="Номер дома"
            :error="registrationHouseError"
            @update:model-value="validateAddressField('house', 'Дом', true)"
          />
          <AppInput
            v-model="state.registrationAddress.building"
            label="Корпус"
            placeholder="Корпус"
            :error="registrationBuildingError"
            @update:model-value="validateAddressField('building', 'Корпус', true)"
          />
          <AppInput
            v-model="state.registrationAddress.apartment"
            label="Квартира"
            placeholder="Квартира"
            :error="registrationApartmentError"
            @update:model-value="validateAddressField('apartment', 'Квартира', true)"
          />
        </div>
      </AppCard>
    </div>

    <div class="listener-form__grid listener-form__grid--double">
      <AppCard title="Образование">
        <div class="form-stack">
          <div class="toggle-field">
            <AppCheckbox
              v-model="state.listener.looting_education"
              label="Получает образование сейчас"
            />
          </div>

          <AppInput
            v-model="state.education.diplom_seria"
            label="Серия диплома"
            placeholder="6 цифр"
            maxlength="6"
            inputmode="numeric"
            :disabled="educationDisabled"
            :error="diplomaSeriaError"
            @update:model-value="handleEducationDigitsInput('diplom_seria')"
          />
          <AppInput
            v-model="state.education.diplom_number"
            label="Номер диплома"
            placeholder="7 цифр"
            maxlength="7"
            inputmode="numeric"
            :disabled="educationDisabled"
            :error="diplomaNumberError"
            @update:model-value="handleEducationDigitsInput('diplom_number')"
          />
          <AppInput
            v-model="state.education.date_given"
            label="Дата выдачи"
            type="date"
            :disabled="educationDisabled"
          />
          <AppInput
            v-model="state.education.city"
            label="Город"
            placeholder="Город"
            :disabled="educationDisabled"
            :error="educationCityError"
            @update:model-value="validateEducationAddressField('city', 'Город')"
          />
          <AppInput
            v-model="state.education.region"
            label="Регион"
            placeholder="Регион"
            :disabled="educationDisabled"
            :error="educationRegionError"
            @update:model-value="validateEducationAddressField('region', 'Регион')"
          />
          <AppInput
            v-model="state.education.educational_institution"
            label="Учебное заведение"
            placeholder="Название учебного заведения"
            :disabled="educationDisabled"
          />
          <AppInput
            v-model="state.education.speciality"
            label="Специальность"
            placeholder="Специальность"
            :disabled="educationDisabled"
          />

          <AppSelect
            v-model="state.education.level_education"
            label="Уровень образования"
            placeholder="Выберите уровень"
            :disabled="educationDisabled || levelsLoading"
          >
            <option
              v-for="level in levels"
              :key="level.id_level_education"
              :value="level.id_level_education"
            >
              {{ level.education }}
            </option>
          </AppSelect>
        </div>
      </AppCard>

      <AppCard title="Место работы">
        <div class="form-stack">
          <AppInput v-model="state.placeWork.name_company" label="Компания" placeholder="Название компании" />
          <AppInput v-model="state.placeWork.job_title" label="Должность" placeholder="Текущая должность" />
          <AppInput
            v-model="state.placeWork.all_experience"
            label="Общий стаж"
            placeholder="Количество лет"
            inputmode="numeric"
            :error="workExperienceError"
            @update:model-value="numbersOnlyWork('all_experience')"
          />
          <AppInput
            v-model="state.placeWork.job_title_expirience"
            label="Стаж по должности"
            placeholder="Количество лет"
            inputmode="numeric"
            :error="workRoleExperienceError"
            @update:model-value="numbersOnlyWork('job_title_expirience')"
          />
        </div>
      </AppCard>
    </div>

    <p v-if="submitError" class="listener-form__message listener-form__message--error">
      {{ submitError }}
    </p>

    <div class="listener-form__actions">
      <AppButton type="button" variant="ghost" @click="emit('cancel')">
        Назад
      </AppButton>

      <AppButton type="submit" :disabled="loading">
        {{ loading ? 'Сохраняем...' : submitLabel }}
      </AppButton>
    </div>
  </form>
</template>

<style scoped>
.listener-form__grid {
  display: grid;
  gap: 1rem;
}

.listener-form__grid--triple {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.listener-form__grid--double {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.form-stack {
  display: grid;
  gap: 0.9rem;
}

.listener-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.listener-form__message {
  margin: 0;
  padding: 0.9rem 1rem;
  border-radius: 1rem;
}

.listener-form__message--error {
  background: rgba(254, 226, 226, 0.9);
  color: #991b1b;
}

.toggle-field {
  min-height: 3rem;
  padding: 0.85rem 1rem;
  border-radius: 1rem;
  background: rgba(255, 255, 255, 0.72);
  border: 1px solid rgba(15, 23, 42, 0.08);
  color: #0f172a;
}

@media (max-width: 1100px) {
  .listener-form__grid--triple,
  .listener-form__grid--double {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .listener-form__actions {
    flex-direction: column-reverse;
  }
}
</style>
