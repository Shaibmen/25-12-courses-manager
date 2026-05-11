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
import AppSelect from '../../ui/AppSelect.vue'

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
const priceError = ref('')
const typeError = ref('')
const divisionError = ref('')

const textPattern = /[^0-9А-Яа-яЁёA-Za-z.,()\-/"«»\s]/g

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

const handleTimeInput = () => {
  const sanitized = sanitizeDigits(state.time_education, 5)
  if (sanitized !== state.time_education) {
    state.time_education = sanitized
    timeError.value = 'Допустимы только цифры'
    return
  }
  validateTime()
}

const validateTime = () => {
  const sanitized = sanitizeDigits(state.time_education, 5)
  if (!sanitized) {
    timeError.value = 'Укажите длительность'
    return false
  }

  if (Number(sanitized) < 1) {
    timeError.value = 'Длительность не может быть меньше 1'
    return false
  }

  timeError.value = ''
  return true
}

const handlePriceInput = () => {
  const sanitized = sanitizeDigits(state.price, 9)
  if (sanitized !== state.price) {
    state.price = sanitized
    priceError.value = 'Допустимы только цифры'
    return
  }
  validatePrice()
}

const validatePrice = () => {
  const sanitized = sanitizeDigits(state.price, 9)
  if (!sanitized) {
    priceError.value = 'Укажите цену'
    return false
  }

  if (Number(sanitized) < 0) {
    priceError.value = 'Цена не может быть меньше 0'
    return false
  }

  priceError.value = ''
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
  state.time_education = sanitizeDigits(state.time_education, 5)
  state.price = sanitizeDigits(state.price, 9)
  nameError.value = ''
  timeError.value = ''
  priceError.value = ''
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
    validateTime(),
    validatePrice(),
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
      price: Number(state.price),
      id_education_type: state.id_education_type,
      id_divisions_education: state.id_divisions_education
    })

    return
  }

  emit('submit', {
    name_prof_education: state.name_prof_education.trim(),
    time_education: Number(state.time_education),
    price: Number(state.price),
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
          v-model="state.price"
          label="Цена (₽)"
          placeholder="15000"
          :error="priceError"
          inputmode="numeric"
          maxlength="9"
          @update:model-value="handlePriceInput"
        />
      </div>

      <div class="program-form__grid">
        <AppSelect
          v-model="state.id_education_type"
          label="Тип обучения"
          placeholder="Выберите тип обучения"
          :error="typeError"
          @update:model-value="validateType"
        >
          <option
            v-for="type in educationTypes"
            :key="type.id_educationType"
            :value="type.id_educationType"
          >
            {{ type.typeName }}
          </option>
        </AppSelect>

        <AppSelect
          v-model="state.id_divisions_education"
          label="Подразделение"
          placeholder="Выберите подразделение"
          :error="divisionError"
          @update:model-value="validateDivision"
        >
          <option
            v-for="division in divisions"
            :key="division.id_divisionsEducation"
            :value="division.id_divisionsEducation"
          >
            {{ division.divisions }}
          </option>
        </AppSelect>
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

@media (max-width: 780px) {
  .program-form__grid {
    grid-template-columns: 1fr;
  }

  .program-form__actions {
    flex-direction: column;
  }
}
</style>
