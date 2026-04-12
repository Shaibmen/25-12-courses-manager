<script setup lang="ts">
import type { ContractorPayload } from '../../../types/enrollment'
import { normalizeApiDate } from '../../../utils/date'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'

const createEmptyContractor = (): ContractorPayload => ({
  contractor: {
    first_name: '',
    second_name: '',
    middle_name: '',
    contact_phone: '',
    email: ''
  },
  passport: {
    place_birth: '',
    citizenship: '',
    gender: '',
    seria: '',
    number: '',
    passport_given: '',
    date_given: '',
    code: ''
  },
  reg_address: {
    mail_index: '',
    region: '',
    city: '',
    street: '',
    house: '',
    building: '',
    apartment: ''
  }
})

const props = withDefaults(defineProps<{
  initialState?: ContractorPayload | null
  loading?: boolean
}>(), {
  initialState: null,
  loading: false
})

const emit = defineEmits<{
  submit: [payload: ContractorPayload]
  cancel: []
}>()

const notifications = useNotifications()
const state = reactive<ContractorPayload>(createEmptyContractor())
const firstNameError = ref('')
const secondNameError = ref('')
const middleNameError = ref('')
const phoneError = ref('')
const emailError = ref('')
const passportSeriaError = ref('')
const passportNumberError = ref('')
const passportCodeError = ref('')
const registrationIndexError = ref('')

const namePattern = /[^А-Яа-яЁёA-Za-z.\s-]/g
const emailPattern = /[^A-Za-z0-9@._-]/g

const syncState = (nextState?: ContractorPayload | null) => {
  const source = nextState || createEmptyContractor()

  Object.assign(state.contractor, source.contractor || createEmptyContractor().contractor)
  Object.assign(state.passport, source.passport || createEmptyContractor().passport, {
    date_given: normalizeApiDate(source.passport?.date_given)
  })
  Object.assign(state.reg_address, source.reg_address || createEmptyContractor().reg_address)

  firstNameError.value = ''
  secondNameError.value = ''
  middleNameError.value = ''
  phoneError.value = ''
  emailError.value = ''
  passportSeriaError.value = ''
  passportNumberError.value = ''
  passportCodeError.value = ''
  registrationIndexError.value = ''
}

watch(
  () => props.initialState,
  (value) => {
    syncState(value)
  },
  { immediate: true, deep: true }
)

const validateNameField = (
  key: 'first_name' | 'second_name' | 'middle_name',
  errorRef: Ref<string>,
  required = true
) => {
  state.contractor[key] = state.contractor[key].replace(namePattern, '').replace(/\s{2,}/g, ' ')
  errorRef.value = !state.contractor[key].trim() && required ? 'Поле обязательно' : ''
  return !errorRef.value
}

const formatPhone = () => {
  let digits = state.contractor.contact_phone.replace(/\D/g, '')

  if (!digits) {
    phoneError.value = 'Укажите телефон'
    state.contractor.contact_phone = ''
    return false
  }

  if (!digits.startsWith('7')) {
    digits = `7${digits}`
  }

  digits = digits.slice(0, 11)
  state.contractor.contact_phone = `+${digits}`
  phoneError.value = digits.length === 11 ? '' : 'Телефон должен быть в формате +7XXXXXXXXXX'
  return !phoneError.value
}

const validateEmail = () => {
  state.contractor.email = state.contractor.email.replace(emailPattern, '')

  if (!state.contractor.email) {
    emailError.value = 'Укажите email'
    return false
  }

  const valid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(state.contractor.email)
  emailError.value = valid ? '' : 'Некорректный email'
  return valid
}

const onlyDigits = (section: 'passport' | 'reg_address', key: string, maxLength: number, errorRef: Ref<string>) => {
  const bucket = state[section] as Record<string, string>
  bucket[key] = String(bucket[key] || '').replace(/\D/g, '').slice(0, maxLength)
  errorRef.value = ''
}

const formatPassportCode = () => {
  let digits = state.passport.code.replace(/\D/g, '').slice(0, 6)

  if (digits.length >= 3) {
    digits = `${digits.slice(0, 3)}-${digits.slice(3)}`
  }

  state.passport.code = digits
  passportCodeError.value = digits.length === 7 ? '' : 'Код подразделения: 123-456'
  return !passportCodeError.value
}

const validateForm = () => {
  const isValid = [
    validateNameField('second_name', secondNameError),
    validateNameField('first_name', firstNameError),
    validateNameField('middle_name', middleNameError, false),
    formatPhone(),
    validateEmail(),
    /^\d{4}$/.test(state.passport.seria) || !state.passport.seria ? !(passportSeriaError.value = state.passport.seria ? 'Серия: 4 цифры' : '') : false,
    /^\d{6}$/.test(state.passport.number) || !state.passport.number ? !(passportNumberError.value = state.passport.number ? 'Номер: 6 цифр' : '') : false,
    formatPassportCode(),
    /^\d{6}$/.test(state.reg_address.mail_index) || !state.reg_address.mail_index ? !(registrationIndexError.value = state.reg_address.mail_index ? 'Индекс: 6 цифр' : '') : false
  ].every(Boolean)

  if (!isValid) {
    notifications.error('Проверьте поля заказчика.', 'Запись на курс')
  }

  return isValid
}

const handleSubmit = () => {
  if (!validateForm()) {
    return
  }

  emit('submit', JSON.parse(JSON.stringify(state)) as ContractorPayload)
}
</script>

<template>
  <AppCard title="Заказчик">
    <form class="contractor-form" @submit.prevent="handleSubmit">
      <div class="contractor-form__grid">
        <AppInput v-model="state.contractor.second_name" label="Фамилия" :error="secondNameError" @update:model-value="validateNameField('second_name', secondNameError)" />
        <AppInput v-model="state.contractor.first_name" label="Имя" :error="firstNameError" @update:model-value="validateNameField('first_name', firstNameError)" />
        <AppInput v-model="state.contractor.middle_name" label="Отчество" :error="middleNameError" @update:model-value="validateNameField('middle_name', middleNameError, false)" />
        <AppInput v-model="state.contractor.contact_phone" label="Телефон" :error="phoneError" inputmode="tel" @update:model-value="formatPhone" />
        <AppInput v-model="state.contractor.email" label="Email" :error="emailError" @update:model-value="validateEmail" />
      </div>

      <div class="contractor-form__grid">
        <AppInput v-model="state.passport.place_birth" label="Место рождения" />
        <AppInput v-model="state.passport.citizenship" label="Гражданство" />
        <AppInput v-model="state.passport.seria" label="Серия паспорта" :error="passportSeriaError" maxlength="4" inputmode="numeric" @update:model-value="onlyDigits('passport', 'seria', 4, passportSeriaError)" />
        <AppInput v-model="state.passport.number" label="Номер паспорта" :error="passportNumberError" maxlength="6" inputmode="numeric" @update:model-value="onlyDigits('passport', 'number', 6, passportNumberError)" />
        <AppInput v-model="state.passport.passport_given" label="Кем выдан" />
        <AppInput v-model="state.passport.date_given" label="Дата выдачи" type="date" />
        <AppInput v-model="state.passport.code" label="Код подразделения" :error="passportCodeError" maxlength="7" inputmode="numeric" @update:model-value="formatPassportCode" />
      </div>

      <div class="contractor-form__grid">
        <AppInput v-model="state.reg_address.mail_index" label="Индекс" :error="registrationIndexError" maxlength="6" inputmode="numeric" @update:model-value="onlyDigits('reg_address', 'mail_index', 6, registrationIndexError)" />
        <AppInput v-model="state.reg_address.region" label="Регион" />
        <AppInput v-model="state.reg_address.city" label="Город" />
        <AppInput v-model="state.reg_address.street" label="Улица" />
        <AppInput v-model="state.reg_address.house" label="Дом" />
        <AppInput v-model="state.reg_address.building" label="Корпус" />
        <AppInput v-model="state.reg_address.apartment" label="Квартира" />
      </div>

      <div class="contractor-form__actions">
        <AppButton type="button" variant="ghost" @click="$emit('cancel')">Отмена</AppButton>
        <AppButton type="submit" :disabled="loading">{{ loading ? 'Сохраняем...' : 'Сохранить заказчика' }}</AppButton>
      </div>
    </form>
  </AppCard>
</template>

<style scoped>
.contractor-form {
  display: grid;
  gap: 1rem;
}

.contractor-form__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.contractor-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

@media (max-width: 1100px) {
  .contractor-form__grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .contractor-form__actions {
    flex-direction: column;
  }
}
</style>
