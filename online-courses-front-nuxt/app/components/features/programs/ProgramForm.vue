<script setup lang="ts">
import type { DivisionItem, EducationTypeItem } from '../../../types/catalogs'
import type {
  ProgramFormState,
  ProgramPayload,
  ProgramUpdatePayload
} from '../../../types/program'
import { createEmptyProgramFormState } from '../../../types/program'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'

const props = withDefaults(defineProps<{
  title: string
  submitLabel: string
  loading?: boolean
  initialState?: ProgramFormState | null
  educationTypes?: EducationTypeItem[]
  divisions?: DivisionItem[]
  mode?: 'create' | 'edit'
}>(), {
  loading: false,
  initialState: null,
  educationTypes: () => [],
  divisions: () => [],
  mode: 'create'
})

const emit = defineEmits<{
  submit: [payload: ProgramPayload | ProgramUpdatePayload]
  cancel: []
}>()

const notifications = useNotifications()
const state = reactive<ProgramFormState>(createEmptyProgramFormState())
const nameError = ref('')
const timeError = ref('')
const individualPriceError = ref('')
const groupPriceError = ref('')
const campusPriceError = ref('')
const typeError = ref('')
const divisionError = ref('')

const textPattern = /[^0-9А-Яа-яЁёA-Za-z.,()\-/"«»\s]/g
const numberFields: Array<keyof Pick<ProgramFormState, 'time_education' | 'individual_price' | 'group_price' | 'campus_price'>> = [
  'time_education',
  'individual_price',
  'group_price',
  'campus_price'
]

const selectFieldClass = (hasPlaceholder: boolean, hasError: boolean) => ({
  'form-select-field--placeholder': hasPlaceholder,
  'form-select-field--error': hasError
})

const validateRequiredText = (value: string, label: string) =>
  value.trim() ? '' : `Укажите ${label.toLowerCase()}`

const validateName = () => {
  const sanitized = state.name_prof_education.replace(textPattern, '').replace(/\s{2,}/g, ' ')

  if (sanitized !== state.name_prof_education) {
    state.name_prof_education = sanitized
    nameError.value = 'Удалены недопустимые символы'
    return false
  }

  nameError.value = validateRequiredText(sanitized, 'название программы')
  return !nameError.value
}

const sanitizeDigits = (value: string, maxLength = 6) => value.replace(/\D/g, '').slice(0, maxLength)

const handleTimeInput = () => handleNumericInput('time_education', 'Длительность', timeError, 5, 1)
const handleIndividualPriceInput = () =>
  handleNumericInput('individual_price', 'Индивидуальная цена', individualPriceError, 9)
const handleGroupPriceInput = () =>
  handleNumericInput('group_price', 'Групповая цена', groupPriceError, 9)
const handleCampusPriceInput = () =>
  handleNumericInput('campus_price', 'Кампусная цена', campusPriceError, 9)

const handleNumericInput = (
  key: keyof Pick<ProgramFormState, 'time_education' | 'individual_price' | 'group_price' | 'campus_price'>,
  label: string,
  errorRef: Ref<string>,
  maxLength: number,
  minValue = 0
) => {
  const sanitized = sanitizeDigits(state[key], maxLength)

  if (sanitized !== state[key]) {
    state[key] = sanitized
    errorRef.value = 'Допустимы только цифры'
    return
  }

  void validateNumericField(key, label, errorRef, minValue, maxLength)
}

const validateNumericField = (
  key: keyof Pick<ProgramFormState, 'time_education' | 'individual_price' | 'group_price' | 'campus_price'>,
  label: string,
  errorRef: Ref<string>,
  minValue = 0,
  maxLength = 6
) => {
  const sanitized = sanitizeDigits(state[key], maxLength)

  if (sanitized !== state[key]) {
    state[key] = sanitized
    errorRef.value = 'Допустимы только цифры'
    return false
  }

  if (!sanitized) {
    errorRef.value = `Укажите ${label.toLowerCase()}`
    return false
  }

  const numericValue = Number(sanitized)

  if (numericValue < minValue) {
    errorRef.value = `${label} не может быть меньше ${minValue}`
    return false
  }

  errorRef.value = ''
  return true
}

const validateType = () => {
  typeError.value = state.id_education_type ? '' : 'Выберите тип обучения'
  return !typeError.value
}

const validateDivision = () => {
  divisionError.value = state.id_divisions_education ? '' : 'Выберите подразделение'
  return !divisionError.value
}

const syncState = (nextState?: ProgramFormState | null) => {
  Object.assign(state, createEmptyProgramFormState(), nextState || {})
  numberFields.forEach((key) => {
    state[key] = sanitizeDigits(
      state[key],
      key === 'time_education' ? 5 : 9
    )
  })
  nameError.value = ''
  timeError.value = ''
  individualPriceError.value = ''
  groupPriceError.value = ''
  campusPriceError.value = ''
  typeError.value = ''
  divisionError.value = ''
}

watch(
  () => props.initialState,
  (value) => {
    syncState(value)
  },
  { immediate: true, deep: true }
)

const validateForm = () => {
  const isValid = [
    validateName(),
    validateNumericField('time_education', 'Длительность', timeError, 1, 5),
    validateNumericField('individual_price', 'Индивидуальная цена', individualPriceError, 0, 9),
    validateNumericField('group_price', 'Групповая цена', groupPriceError, 0, 9),
    validateNumericField('campus_price', 'Кампусная цена', campusPriceError, 0, 9),
    validateType(),
    validateDivision()
  ].every(Boolean)

  if (!isValid) {
    notifications.error('Проверьте обязательные поля программы.', 'Программы обучения')
  }

  return isValid
}

const handleSubmit = () => {
  if (!validateForm()) {
    return
  }

  if (props.mode === 'edit') {
    emit('submit', {
      name_prof_education: state.name_prof_education.trim(),
      time_education: Number(state.time_education),
      individual_price: Number(state.individual_price),
      group_price: Number(state.group_price),
      campus_price: Number(state.campus_price),
      ID_EducationType: state.id_education_type,
      ID_DivisionsEducation: state.id_divisions_education
    })

    return
  }

  emit('submit', {
    name_prof_education: state.name_prof_education.trim(),
    time_education: Number(state.time_education),
    individual_price: Number(state.individual_price),
    group_price: Number(state.group_price),
    campus_price: Number(state.campus_price),
    id_educationtype: state.id_education_type,
    id_divisionseducation: state.id_divisions_education
  })
}
</script>

<template>
  <AppCard :title="title">
    <form class="program-form" @submit.prevent="handleSubmit">
      <AppInput
        v-model="state.name_prof_education"
        label="Название программы"
        placeholder="Например: Повышение квалификации по охране труда"
        :error="nameError"
        maxlength="255"
        @update:model-value="validateName"
      />

      <div class="program-form__grid">
        <AppInput
          v-model="state.time_education"
          label="Длительность (часы)"
          placeholder="72"
          :error="timeError"
          inputmode="numeric"
          maxlength="5"
          @update:model-value="handleTimeInput"
        />

        <AppInput
          v-model="state.individual_price"
          label="Индивидуальная цена (₽)"
          placeholder="15000"
          :error="individualPriceError"
          inputmode="numeric"
          maxlength="9"
          @update:model-value="handleIndividualPriceInput"
        />

        <AppInput
          v-model="state.group_price"
          label="Групповая цена (₽)"
          placeholder="12000"
          :error="groupPriceError"
          inputmode="numeric"
          maxlength="9"
          @update:model-value="handleGroupPriceInput"
        />

        <AppInput
          v-model="state.campus_price"
          label="Кампусная цена (₽)"
          placeholder="8000"
          :error="campusPriceError"
          inputmode="numeric"
          maxlength="9"
          @update:model-value="handleCampusPriceInput"
        />
      </div>

      <div class="program-form__grid">
        <label class="app-select">
          <span class="form-select-field__label">Тип обучения</span>

          <div class="app-select__field" :class="selectFieldClass(!state.id_education_type, Boolean(typeError))">
            <select
              v-model="state.id_education_type"
              class="app-select__control"
              @change="validateType"
            >
              <option value="">Выберите тип обучения</option>
              <option
                v-for="type in educationTypes"
                :key="type.id_educationType"
                :value="type.id_educationType"
              >
                {{ type.typeName }}
              </option>
            </select>

            <span class="app-select__icon" aria-hidden="true">
              <svg viewBox="0 0 20 20" fill="none">
                <path
                  d="M5 7.5L10 12.5L15 7.5"
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.8"
                />
              </svg>
            </span>
          </div>
          <span v-if="typeError" class="form-select-field__message">{{ typeError }}</span>
        </label>

        <label class="app-select">
          <span class="form-select-field__label">Подразделение</span>

          <div class="app-select__field" :class="selectFieldClass(!state.id_divisions_education, Boolean(divisionError))">
            <select
              v-model="state.id_divisions_education"
              class="app-select__control"
              @change="validateDivision"
            >
              <option value="">Выберите подразделение</option>
              <option
                v-for="division in divisions"
                :key="division.id_divisionsEducation"
                :value="division.id_divisionsEducation"
              >
                {{ division.divisions }}
              </option>
            </select>

            <span class="app-select__icon" aria-hidden="true">
              <svg viewBox="0 0 20 20" fill="none">
                <path
                  d="M5 7.5L10 12.5L15 7.5"
                  stroke="currentColor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="1.8"
                />
              </svg>
            </span>
          </div>
          <span v-if="divisionError" class="form-select-field__message">{{ divisionError }}</span>
        </label>
      </div>

      <div class="program-form__actions">
        <AppButton type="button" variant="ghost" @click="$emit('cancel')">Назад</AppButton>
        <AppButton type="submit" :disabled="loading">
          {{ loading ? 'Сохраняем...' : submitLabel }}
        </AppButton>
      </div>
    </form>
  </AppCard>
</template>

<style scoped>
.program-form {
  display: grid;
  gap: 1rem;
}

.program-form__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.program-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.app-select {
  display: grid;
  gap: 0.45rem;
}

.form-select-field__label {
  font-size: 0.9rem;
  font-weight: 600;
  color: #0f172a;
}

.app-select__field {
  position: relative;
  border-radius: 1rem;
  border: 1px solid rgba(15, 23, 42, 0.12);
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.97) 0%, rgba(248, 250, 252, 0.94) 100%);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.7);
  transition:
    border-color 180ms ease,
    box-shadow 180ms ease,
    transform 180ms ease;
}

.app-select__field:focus-within {
  border-color: rgba(59, 130, 246, 0.5);
  box-shadow: 0 0 0 5px rgba(59, 130, 246, 0.12);
  transform: translateY(-1px);
}

.app-select__field--placeholder .app-select__control {
  color: #64748b;
}

.app-select__field--error {
  border-color: rgba(220, 38, 38, 0.35);
  box-shadow: 0 0 0 4px rgba(248, 113, 113, 0.08);
}

.app-select__control {
  appearance: none;
  -webkit-appearance: none;
  position: relative;
  z-index: 1;
  width: 100%;
  min-height: 3rem;
  padding: 0.8rem 2.75rem 0.8rem 1rem;
  border: none;
  border-radius: 1rem;
  background: transparent;
  color: #0f172a;
  font: inherit;
  outline: none;
  cursor: pointer;
}

.app-select__icon {
  position: absolute;
  top: 50%;
  right: 0.95rem;
  z-index: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 1rem;
  height: 1rem;
  color: #475569;
  pointer-events: none;
  transform: translateY(-50%);
}

.app-select__icon svg {
  width: 100%;
  height: 100%;
}

.form-select-field__message {
  font-size: 0.84rem;
  color: #b91c1c;
}

@media (max-width: 780px) {
  .program-form__grid {
    grid-template-columns: 1fr;
  }

  .program-form__actions {
    flex-direction: column;
  }
}
</style>
