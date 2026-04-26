<script setup lang="ts">
import type { ContractorPayload } from '../../../types/enrollment'
import { normalizeApiDate } from '../../../utils/date'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'
import AppSelect from '../../ui/AppSelect.vue'

type ErrorState = {
  second_name: string
  first_name: string
  middle_name: string
  contact_phone: string
  email: string
  place_birth: string
  citizenship: string
  gender: string
  seria: string
  number: string
  passport_given: string
  date_given: string
  code: string
  mail_index: string
  region: string
  city: string
  street: string
  house: string
  building: string
  apartment: string
}

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

const createEmptyErrors = (): ErrorState => ({
  second_name: '',
  first_name: '',
  middle_name: '',
  contact_phone: '',
  email: '',
  place_birth: '',
  citizenship: '',
  gender: '',
  seria: '',
  number: '',
  passport_given: '',
  date_given: '',
  code: '',
  mail_index: '',
  region: '',
  city: '',
  street: '',
  house: '',
  building: '',
  apartment: ''
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
const errors = reactive<ErrorState>(createEmptyErrors())

const namePattern = /[^А-Яа-яЁёA-Za-z.\s-]/g
const emailPattern = /[^A-Za-z0-9@._+-]/g

const setError = (key: keyof ErrorState, message: string) => {
  errors[key] = message
  return !message
}

const resetErrors = () => {
  Object.assign(errors, createEmptyErrors())
}

const trimValue = (value: string) => value.replace(/\s{2,}/g, ' ').trim()

const syncState = (nextState?: ContractorPayload | null) => {
  const emptyState = createEmptyContractor()
  const source = nextState || emptyState

  Object.assign(state.contractor, emptyState.contractor, source.contractor || {})
  Object.assign(state.passport, emptyState.passport, source.passport || {}, {
    date_given: normalizeApiDate(source.passport?.date_given)
  })
  Object.assign(state.reg_address, emptyState.reg_address, source.reg_address || {})

  resetErrors()
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
  required = true
) => {
  state.contractor[key] = state.contractor[key].replace(namePattern, '').replace(/\s{2,}/g, ' ')
  const value = trimValue(state.contractor[key])
  state.contractor[key] = value

  const labels = {
    first_name: 'Укажите имя',
    second_name: 'Укажите фамилию',
    middle_name: ''
  } as const

  return setError(key, !value && required ? labels[key] : '')
}

const validateTextField = (
  section: 'passport' | 'reg_address',
  key: string,
  errorKey: keyof ErrorState,
  label: string
) => {
  const bucket = state[section] as Record<string, string>
  bucket[key] = trimValue(bucket[key] || '')
  return setError(errorKey, bucket[key] ? '' : `Укажите ${label}`)
}

const formatPhone = () => {
  let digits = state.contractor.contact_phone.replace(/\D/g, '')

  if (!digits) {
    state.contractor.contact_phone = ''
    return setError('contact_phone', 'Укажите телефон в формате +7XXXXXXXXXX')
  }

  if (digits.startsWith('8')) {
    digits = `7${digits.slice(1)}`
  } else if (!digits.startsWith('7')) {
    digits = `7${digits}`
  }

  digits = digits.slice(0, 11)
  state.contractor.contact_phone = `+${digits}`

  return setError(
    'contact_phone',
    digits.length === 11 ? '' : 'Телефон должен быть в формате +7XXXXXXXXXX'
  )
}

const validateEmail = () => {
  state.contractor.email = state.contractor.email.replace(emailPattern, '').trim()

  if (!state.contractor.email) {
    return setError('email', 'Укажите email')
  }

  const valid = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(state.contractor.email)
  return setError('email', valid ? '' : 'Проверьте формат email')
}

const onlyDigits = (
  section: 'passport' | 'reg_address',
  key: string,
  maxLength: number,
  errorKey?: keyof ErrorState
) => {
  const bucket = state[section] as Record<string, string>
  bucket[key] = String(bucket[key] || '').replace(/\D/g, '').slice(0, maxLength)

  if (errorKey) {
    errors[errorKey] = ''
  }
}

const validateExactDigits = (value: string, length: number, errorKey: keyof ErrorState, label: string) =>
  setError(errorKey, value.length === length ? '' : `${label}: ${length} цифр`)

const formatPassportCode = () => {
  let digits = state.passport.code.replace(/\D/g, '').slice(0, 6)

  if (digits.length >= 3) {
    digits = `${digits.slice(0, 3)}-${digits.slice(3)}`
  }

  state.passport.code = digits
  return setError('code', /^\d{3}-\d{3}$/.test(digits) ? '' : 'Код подразделения: 123-456')
}

const validateDateGiven = () => setError('date_given', state.passport.date_given ? '' : 'Укажите дату выдачи')

const validateGender = () => setError('gender', state.passport.gender ? '' : 'Выберите пол')

const validateForm = () => {
  const isValid = [
    validateNameField('second_name'),
    validateNameField('first_name'),
    validateNameField('middle_name', false),
    formatPhone(),
    validateEmail(),
    validateTextField('passport', 'place_birth', 'place_birth', 'место рождения'),
    validateTextField('passport', 'citizenship', 'citizenship', 'гражданство'),
    validateGender(),
    validateExactDigits(state.passport.seria, 4, 'seria', 'Серия паспорта'),
    validateExactDigits(state.passport.number, 6, 'number', 'Номер паспорта'),
    validateTextField('passport', 'passport_given', 'passport_given', 'кем выдан паспорт'),
    validateDateGiven(),
    formatPassportCode(),
    validateExactDigits(state.reg_address.mail_index, 6, 'mail_index', 'Индекс'),
    validateTextField('reg_address', 'region', 'region', 'регион'),
    validateTextField('reg_address', 'city', 'city', 'город'),
    validateTextField('reg_address', 'street', 'street', 'улицу'),
    validateTextField('reg_address', 'house', 'house', 'дом'),
    validateTextField('reg_address', 'building', 'building', 'корпус или строение'),
    validateTextField('reg_address', 'apartment', 'apartment', 'квартиру')
  ].every(Boolean)

  if (!isValid) {
    notifications.error('Проверьте поля заказчика: ошибки показаны прямо в форме.', 'Запись на курс')
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
      <section class="contractor-section">
        <div class="contractor-section__head">
          <h3>Контактные данные</h3>
          <p>Заполните ФИО и контакты человека, который заключает договор.</p>
        </div>

        <div class="contractor-form__grid contractor-form__grid--contacts">
          <AppInput
            v-model="state.contractor.second_name"
            label="Фамилия"
            placeholder="Например: Иванов"
            maxlength="255"
            :error="errors.second_name"
            @update:model-value="validateNameField('second_name')"
          />
          <AppInput
            v-model="state.contractor.first_name"
            label="Имя"
            placeholder="Например: Иван"
            maxlength="255"
            :error="errors.first_name"
            @update:model-value="validateNameField('first_name')"
          />
          <AppInput
            v-model="state.contractor.middle_name"
            label="Отчество"
            placeholder="Например: Иванович"
            maxlength="255"
            :error="errors.middle_name"
            @update:model-value="validateNameField('middle_name', false)"
          />
          <AppInput
            v-model="state.contractor.contact_phone"
            label="Телефон"
            placeholder="+79991234567"
            maxlength="12"
            inputmode="tel"
            :error="errors.contact_phone"
            @update:model-value="formatPhone"
          />
          <AppInput
            v-model="state.contractor.email"
            label="Email"
            placeholder="example@mail.ru"
            maxlength="255"
            :error="errors.email"
            @update:model-value="validateEmail"
          />
        </div>
      </section>

      <section class="contractor-section">
        <div class="contractor-section__head">
          <h3>Паспорт</h3>
          <p>Нужны данные паспорта заказчика в том виде, как они указаны в документе.</p>
        </div>

        <div class="contractor-form__grid contractor-form__grid--passport">
          <AppInput
            v-model="state.passport.place_birth"
            label="Место рождения"
            placeholder="Например: г. Москва"
            maxlength="255"
            :error="errors.place_birth"
            @update:model-value="validateTextField('passport', 'place_birth', 'place_birth', 'место рождения')"
          />
          <AppInput
            v-model="state.passport.citizenship"
            label="Гражданство"
            placeholder="Например: Российская Федерация"
            maxlength="255"
            :error="errors.citizenship"
            @update:model-value="validateTextField('passport', 'citizenship', 'citizenship', 'гражданство')"
          />
          <AppSelect
            v-model="state.passport.gender"
            label="Пол"
            placeholder="Выберите пол"
            :error="errors.gender"
            @update:model-value="validateGender"
          >
            <option value="Мужской">Мужской</option>
            <option value="Женский">Женский</option>
          </AppSelect>
          <AppInput
            v-model="state.passport.seria"
            label="Серия паспорта"
            placeholder="1234"
            maxlength="4"
            inputmode="numeric"
            :error="errors.seria"
            @update:model-value="onlyDigits('passport', 'seria', 4, 'seria')"
          />
          <AppInput
            v-model="state.passport.number"
            label="Номер паспорта"
            placeholder="123456"
            maxlength="6"
            inputmode="numeric"
            :error="errors.number"
            @update:model-value="onlyDigits('passport', 'number', 6, 'number')"
          />
          <AppInput
            v-model="state.passport.code"
            label="Код подразделения"
            placeholder="123-456"
            maxlength="7"
            inputmode="numeric"
            :error="errors.code"
            @update:model-value="formatPassportCode"
          />
          <AppInput
            v-model="state.passport.passport_given"
            label="Кем выдан"
            placeholder="Например: ОВД Тверского района г. Москвы"
            maxlength="255"
            :error="errors.passport_given"
            @update:model-value="validateTextField('passport', 'passport_given', 'passport_given', 'кем выдан паспорт')"
          />
          <AppInput
            v-model="state.passport.date_given"
            label="Дата выдачи"
            type="date"
            :error="errors.date_given"
            @update:model-value="validateDateGiven"
          />
        </div>
      </section>

      <section class="contractor-section">
        <div class="contractor-section__head">
          <h3>Адрес регистрации</h3>
          <p>Заполните адрес по паспорту. У каждого поля есть ожидаемый формат через подпись и плейсхолдер.</p>
        </div>

        <div class="contractor-form__grid contractor-form__grid--address">
          <AppInput
            v-model="state.reg_address.mail_index"
            label="Индекс"
            placeholder="101000"
            maxlength="6"
            inputmode="numeric"
            :error="errors.mail_index"
            @update:model-value="onlyDigits('reg_address', 'mail_index', 6, 'mail_index')"
          />
          <AppInput
            v-model="state.reg_address.region"
            label="Регион"
            placeholder="Например: Московская область"
            maxlength="255"
            :error="errors.region"
            @update:model-value="validateTextField('reg_address', 'region', 'region', 'регион')"
          />
          <AppInput
            v-model="state.reg_address.city"
            label="Город"
            placeholder="Например: Москва"
            maxlength="255"
            :error="errors.city"
            @update:model-value="validateTextField('reg_address', 'city', 'city', 'город')"
          />
          <AppInput
            v-model="state.reg_address.street"
            label="Улица"
            placeholder="Например: Тверская"
            maxlength="255"
            :error="errors.street"
            @update:model-value="validateTextField('reg_address', 'street', 'street', 'улицу')"
          />
          <AppInput
            v-model="state.reg_address.house"
            label="Дом"
            placeholder="Например: 10"
            maxlength="50"
            :error="errors.house"
            @update:model-value="validateTextField('reg_address', 'house', 'house', 'дом')"
          />
          <AppInput
            v-model="state.reg_address.building"
            label="Корпус / строение"
            placeholder="Например: 2"
            maxlength="50"
            :error="errors.building"
            @update:model-value="validateTextField('reg_address', 'building', 'building', 'корпус или строение')"
          />
          <AppInput
            v-model="state.reg_address.apartment"
            label="Квартира"
            placeholder="Например: 15"
            maxlength="50"
            :error="errors.apartment"
            @update:model-value="validateTextField('reg_address', 'apartment', 'apartment', 'квартиру')"
          />
        </div>
      </section>

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
  gap: 1.5rem;
}

.contractor-section {
  display: grid;
  gap: 1rem;
}

.contractor-section__head {
  display: grid;
  gap: 0.35rem;
}

.contractor-section__head h3 {
  margin: 0;
  font-size: 1rem;
  color: #0f172a;
}

.contractor-section__head p {
  margin: 0;
  color: #64748b;
  font-size: 0.92rem;
}

.contractor-form__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.contractor-form__grid--contacts {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.contractor-form__grid--address {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.contractor-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

@media (max-width: 1100px) {
  .contractor-form__grid,
  .contractor-form__grid--contacts,
  .contractor-form__grid--address {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 720px) {
  .contractor-form__actions {
    flex-direction: column;
  }
}
</style>
