<script setup lang="ts">
import type { LegalEntityFormState, LegalEntityPayload } from '../../../types/legalentity'
import { createEmptyLegalEntityFormState } from '../../../types/legalentity'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'

const props = withDefaults(defineProps<{
  title: string
  submitLabel: string
  loading?: boolean
  initialState?: LegalEntityFormState | null
}>(), {
  loading: false,
  initialState: null
})

const emit = defineEmits<{
  submit: [payload: LegalEntityPayload]
  cancel: []
}>()

const notifications = useNotifications()
const state = reactive<LegalEntityFormState>(createEmptyLegalEntityFormState())

const companyNameError = ref('')
const innError = ref('')
const kppError = ref('')
const ogrnError = ref('')
const phoneError = ref('')
const emailError = ref('')
const firstNameError = ref('')
const secondNameError = ref('')
const middleNameError = ref('')
const statusError = ref('')
const registrationIndexError = ref('')

const namePattern = /[^А-Яа-яЁёA-Za-z.\s-]/g
const addressPattern = /[^0-9А-Яа-яЁёA-Za-z.,\-\s]/g
const addressNumberPattern = /[^0-9А-Яа-яЁёA-Za-z.,\-№#\s]/g
const companyPattern = /[^0-9А-Яа-яЁёA-Za-z.,()\-/"«»\s]/g
const emailPattern = /[^A-Za-z0-9@._-]/g

const validateRequiredText = (value: string, label: string, required = true) => {
  if (!value.trim()) {
    return required ? `Укажите ${label.toLowerCase()}` : ''
  }

  return ''
}

const validateCompanyName = () => {
  state.legal.name_company = state.legal.name_company.replace(companyPattern, '').replace(/\s{2,}/g, ' ')
  companyNameError.value = validateRequiredText(state.legal.name_company, 'название компании')
  return !companyNameError.value
}

const validateNameField = (
  key: keyof Pick<LegalEntityFormState['legal'], 'first_name' | 'second_name' | 'middle_name'>,
  label: string,
  errorRef: Ref<string>,
  required = true
) => {
  state.legal[key] = state.legal[key].replace(namePattern, '').replace(/\s{2,}/g, ' ')
  errorRef.value = validateRequiredText(state.legal[key], label, required)
  return !errorRef.value
}

const handleDigitsInput = (
  section: 'legal' | 'regAddress',
  key: string,
  errorRef: Ref<string>,
  label: string,
  maxLength: number,
  required = false
) => {
  const bucket = state[section] as Record<string, string>
  const sanitized = String(bucket[key] || '').replace(/\D/g, '').slice(0, maxLength)
  bucket[key] = sanitized

  if (!sanitized && required) {
    errorRef.value = `Укажите ${label.toLowerCase()}`
    return false
  }

  errorRef.value = ''
  return true
}

const formatPhone = () => {
  let digits = state.legal.phone.replace(/\D/g, '')

  if (!digits) {
    state.legal.phone = ''
    phoneError.value = 'Укажите телефон'
    return false
  }

  if (!digits.startsWith('7')) {
    digits = `7${digits}`
  }

  digits = digits.slice(0, 11)
  state.legal.phone = `+${digits}`
  phoneError.value = digits.length === 11 ? '' : 'Телефон должен быть в формате +7XXXXXXXXXX'
  return !phoneError.value
}

const validateEmail = () => {
  state.legal.email = state.legal.email.replace(emailPattern, '')

  if (!state.legal.email) {
    emailError.value = 'Укажите email'
    return false
  }

  const valid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(state.legal.email)
  emailError.value = valid ? '' : 'Некорректный email'
  return valid
}

const validateStatus = () => {
  state.legal.status = state.legal.status.replace(namePattern, '').replace(/\s{2,}/g, ' ')
  statusError.value = validateRequiredText(state.legal.status, 'должность')
  return !statusError.value
}

const validateAddressField = (
  key: keyof LegalEntityFormState['regAddress'],
  label: string,
  allowNumbers = false
) => {
  const pattern = allowNumbers ? addressNumberPattern : addressPattern
  state.regAddress[key] = state.regAddress[key].replace(pattern, '').replace(/\s{2,}/g, ' ')
  return true
}

const syncState = (nextState?: LegalEntityFormState | null) => {
  Object.assign(state, createEmptyLegalEntityFormState(), nextState || {})
  companyNameError.value = ''
  innError.value = ''
  kppError.value = ''
  ogrnError.value = ''
  phoneError.value = ''
  emailError.value = ''
  firstNameError.value = ''
  secondNameError.value = ''
  middleNameError.value = ''
  statusError.value = ''
  registrationIndexError.value = ''
}

watch(
  () => props.initialState,
  (value) => {
    syncState(value)
  },
  { immediate: true, deep: true }
)

const handleSubmit = () => {
  const isValid = [
    validateCompanyName(),
    handleDigitsInput('legal', 'inn', innError, 'ИНН', 12, true),
    handleDigitsInput('legal', 'kpp', kppError, 'КПП', 9),
    handleDigitsInput('legal', 'ogrn', ogrnError, 'ОГРН', 13),
    formatPhone(),
    validateEmail(),
    validateNameField('first_name', 'имя', firstNameError),
    validateNameField('second_name', 'фамилия', secondNameError),
    validateNameField('middle_name', 'отчество', middleNameError, false),
    validateStatus(),
    handleDigitsInput('regAddress', 'mail_index', registrationIndexError, 'индекс', 6)
  ].every(Boolean)

  if (!isValid) {
    notifications.error('Проверьте обязательные поля юридического лица.', 'Юридические лица')
    return
  }

  emit('submit', {
    legal_entity: {
      ...state.legal,
      name_company: state.legal.name_company.trim(),
      first_name: state.legal.first_name.trim(),
      second_name: state.legal.second_name.trim(),
      middle_name: state.legal.middle_name.trim(),
      status: state.legal.status.trim()
    },
    reg_address: {
      ...state.regAddress
    }
  })
}
</script>

<template>
  <form class="legalentity-form" @submit.prevent="handleSubmit">
    <div class="legalentity-form__grid">
      <AppCard title="Юридическое лицо">
        <div class="legalentity-form__stack">
          <AppInput
            v-model="state.legal.name_company"
            label="Название компании"
            placeholder="Введите название компании"
            :error="companyNameError"
            maxlength="255"
            @update:model-value="validateCompanyName"
          />
          <AppInput
            v-model="state.legal.inn"
            label="ИНН"
            placeholder="10 или 12 цифр"
            :error="innError"
            inputmode="numeric"
            maxlength="12"
            @update:model-value="handleDigitsInput('legal', 'inn', innError, 'ИНН', 12, true)"
          />
          <AppInput
            v-model="state.legal.kpp"
            label="КПП"
            placeholder="9 цифр"
            :error="kppError"
            inputmode="numeric"
            maxlength="9"
            @update:model-value="handleDigitsInput('legal', 'kpp', kppError, 'КПП', 9)"
          />
          <AppInput
            v-model="state.legal.ogrn"
            label="ОГРН"
            placeholder="13 цифр"
            :error="ogrnError"
            inputmode="numeric"
            maxlength="13"
            @update:model-value="handleDigitsInput('legal', 'ogrn', ogrnError, 'ОГРН', 13)"
          />
          <AppInput
            v-model="state.legal.phone"
            label="Телефон"
            placeholder="+79991234567"
            :error="phoneError"
            inputmode="tel"
            maxlength="12"
            @update:model-value="formatPhone"
          />
          <AppInput
            v-model="state.legal.email"
            label="Email"
            placeholder="company@example.ru"
            :error="emailError"
            maxlength="150"
            @update:model-value="validateEmail"
          />
        </div>
      </AppCard>

      <AppCard title="Представитель">
        <div class="legalentity-form__stack">
          <AppInput
            v-model="state.legal.second_name"
            label="Фамилия"
            placeholder="Введите фамилию"
            :error="secondNameError"
            maxlength="100"
            @update:model-value="validateNameField('second_name', 'фамилия', secondNameError)"
          />
          <AppInput
            v-model="state.legal.first_name"
            label="Имя"
            placeholder="Введите имя"
            :error="firstNameError"
            maxlength="100"
            @update:model-value="validateNameField('first_name', 'имя', firstNameError)"
          />
          <AppInput
            v-model="state.legal.middle_name"
            label="Отчество"
            placeholder="Введите отчество"
            :error="middleNameError"
            maxlength="100"
            @update:model-value="validateNameField('middle_name', 'отчество', middleNameError, false)"
          />
          <AppInput
            v-model="state.legal.status"
            label="Должность"
            placeholder="Например: Генеральный директор"
            :error="statusError"
            maxlength="150"
            @update:model-value="validateStatus"
          />
        </div>
      </AppCard>

      <AppCard title="Адрес регистрации">
        <div class="legalentity-form__stack">
          <AppInput
            v-model="state.regAddress.mail_index"
            label="Индекс"
            placeholder="6 цифр"
            :error="registrationIndexError"
            inputmode="numeric"
            maxlength="6"
            @update:model-value="handleDigitsInput('regAddress', 'mail_index', registrationIndexError, 'индекс', 6)"
          />
          <AppInput
            v-model="state.regAddress.region"
            label="Регион"
            placeholder="Укажите регион"
            maxlength="150"
            @update:model-value="validateAddressField('region', 'Регион')"
          />
          <AppInput
            v-model="state.regAddress.city"
            label="Город"
            placeholder="Укажите город"
            maxlength="150"
            @update:model-value="validateAddressField('city', 'Город')"
          />
          <AppInput
            v-model="state.regAddress.street"
            label="Улица"
            placeholder="Укажите улицу"
            maxlength="150"
            @update:model-value="validateAddressField('street', 'Улица')"
          />
          <AppInput
            v-model="state.regAddress.house"
            label="Дом"
            placeholder="Номер дома"
            maxlength="30"
            @update:model-value="validateAddressField('house', 'Дом', true)"
          />
          <AppInput
            v-model="state.regAddress.building"
            label="Корпус"
            placeholder="Корпус или строение"
            maxlength="30"
            @update:model-value="validateAddressField('building', 'Корпус', true)"
          />
          <AppInput
            v-model="state.regAddress.apartment"
            label="Квартира"
            placeholder="Номер квартиры"
            maxlength="30"
            @update:model-value="validateAddressField('apartment', 'Квартира', true)"
          />
        </div>
      </AppCard>
    </div>

    <div class="legalentity-form__actions">
      <AppButton type="button" variant="ghost" @click="$emit('cancel')">Назад</AppButton>
      <AppButton type="submit" :disabled="loading">
        {{ loading ? 'Сохраняем...' : submitLabel }}
      </AppButton>
    </div>
  </form>
</template>

<style scoped>
.legalentity-form {
  display: grid;
  gap: 1rem;
}

.legalentity-form__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 1rem;
}

.legalentity-form__stack {
  display: grid;
  gap: 0.85rem;
}

.legalentity-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

@media (max-width: 1100px) {
  .legalentity-form__grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .legalentity-form__actions {
    flex-direction: column;
  }
}
</style>
