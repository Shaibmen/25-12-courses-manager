<script setup lang="ts">
import type { ExecuterFormState, ExecuterPayload } from '../../../types/executer'
import { createEmptyExecuterFormState } from '../../../types/executer'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'

const props = withDefaults(defineProps<{
  title: string
  submitLabel: string
  loading?: boolean
  initialState?: ExecuterFormState | null
}>(), {
  loading: false,
  initialState: null
})

const emit = defineEmits<{
  submit: [payload: ExecuterPayload]
  cancel: []
}>()

const notifications = useNotifications()
const state = reactive<ExecuterFormState>(createEmptyExecuterFormState())
const firstNameError = ref('')
const secondNameError = ref('')
const middleNameError = ref('')
const statusError = ref('')

const namePattern = /[^А-Яа-яЁёA-Za-z.\s-]/g

const validateRequiredText = (value: string, label: string, required = true) => {
  if (!value.trim()) {
    return required ? `Укажите ${label.toLowerCase()}` : ''
  }

  return ''
}

const validateNameField = (
  key: keyof Pick<ExecuterFormState, 'first_name' | 'second_name' | 'middle_name'>,
  label: string,
  errorRef: Ref<string>,
  required = true
) => {
  const sanitized = state[key].replace(namePattern, '').replace(/\s{2,}/g, ' ')

  state[key] = sanitized
  errorRef.value = validateRequiredText(sanitized, label, required)
  return !errorRef.value
}

const validateStatus = () => {
  statusError.value = validateRequiredText(state.status, 'должность')
  return !statusError.value
}

const syncState = (nextState?: ExecuterFormState | null) => {
  Object.assign(state, createEmptyExecuterFormState(), nextState || {})
  firstNameError.value = ''
  secondNameError.value = ''
  middleNameError.value = ''
  statusError.value = ''
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
    validateNameField('first_name', 'имя', firstNameError),
    validateNameField('second_name', 'фамилия', secondNameError),
    validateNameField('middle_name', 'отчество', middleNameError, false),
    validateStatus()
  ].every(Boolean)

  if (!isValid) {
    notifications.error('Проверьте обязательные поля исполнителя.', 'Исполнители')
    return
  }

  emit('submit', {
    first_name: state.first_name.trim(),
    second_name: state.second_name.trim(),
    middle_name: state.middle_name.trim() || undefined,
    status: state.status.trim(),
    doverenost: state.doverenost.trim() || undefined
  })
}
</script>

<template>
  <AppCard :title="title">
    <form class="executer-form" @submit.prevent="handleSubmit">
      <div class="executer-form__grid">
        <AppInput
          v-model="state.second_name"
          label="Фамилия"
          placeholder="Введите фамилию"
          :error="secondNameError"
          maxlength="100"
          @update:model-value="validateNameField('second_name', 'фамилия', secondNameError)"
        />
        <AppInput
          v-model="state.first_name"
          label="Имя"
          placeholder="Введите имя"
          :error="firstNameError"
          maxlength="100"
          @update:model-value="validateNameField('first_name', 'имя', firstNameError)"
        />
      </div>

      <div class="executer-form__grid">
        <AppInput
          v-model="state.middle_name"
          label="Отчество"
          placeholder="Введите отчество"
          :error="middleNameError"
          maxlength="100"
          @update:model-value="validateNameField('middle_name', 'отчество', middleNameError, false)"
        />
        <AppInput
          v-model="state.status"
          label="Должность / статус"
          placeholder="Например: Преподаватель"
          :error="statusError"
          maxlength="150"
          @update:model-value="validateStatus"
        />
      </div>

      <AppInput
        v-model="state.doverenost"
        label="Доверенность"
        placeholder="Например: № 17/24"
        maxlength="150"
      />

      <div class="executer-form__actions">
        <AppButton type="button" variant="ghost" @click="$emit('cancel')">Назад</AppButton>
        <AppButton type="submit" :disabled="loading">
          {{ loading ? 'Сохраняем...' : submitLabel }}
        </AppButton>
      </div>
    </form>
  </AppCard>
</template>

<style scoped>
.executer-form {
  display: grid;
  gap: 1rem;
}

.executer-form__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.executer-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

@media (max-width: 780px) {
  .executer-form__grid {
    grid-template-columns: 1fr;
  }

  .executer-form__actions {
    flex-direction: column;
  }
}
</style>
