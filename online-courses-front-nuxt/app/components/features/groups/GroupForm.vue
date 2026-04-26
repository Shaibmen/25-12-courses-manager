<script setup lang="ts">
import type {
  GroupFormScheduleItem,
  GroupFormState,
  GroupPayload
} from '../../../types/group'
import {
  createEmptyGroupFormState,
  createEmptyGroupScheduleItem
} from '../../../types/group'
import { formatApiDate, getDateSortValue, normalizeApiDate } from '../../../utils/date'
import { getRussianPublicHolidayDates, isRussianPublicHoliday } from '../../../utils/russian-holidays'
import AppButton from '../../ui/AppButton.vue'
import AppCard from '../../ui/AppCard.vue'
import AppInput from '../../ui/AppInput.vue'

const props = withDefaults(defineProps<{
  title: string
  submitLabel: string
  loading?: boolean
  initialState?: GroupFormState | null
}>(), {
  loading: false,
  initialState: null
})

const emit = defineEmits<{
  submit: [payload: GroupPayload]
  cancel: []
}>()

const notifications = useNotifications()
const state = reactive<GroupFormState>(createEmptyGroupFormState())
const nameError = ref('')
const scheduleErrors = ref<Record<string, { date: string; theme: string }>>({})

const blockedDates = getRussianPublicHolidayDates()

const schedulePreview = computed(() =>
  [...state.schedule]
    .map((item) => ({
      ...item,
      date: normalizeApiDate(item.date)
    }))
    .filter((item) => item.date || item.theme.trim())
    .sort((left, right) => getDateSortValue(left.date) - getDateSortValue(right.date))
)

const textPattern = /[^0-9А-Яа-яЁёA-Za-z\-./\s]/g

const ensureScheduleError = (id: string) => {
  if (!scheduleErrors.value[id]) {
    scheduleErrors.value[id] = { date: '', theme: '' }
  }

  return scheduleErrors.value[id]
}

const validateName = () => {
  const sanitized = state.name_group.replace(textPattern, '').replace(/\s{2,}/g, ' ')

  if (sanitized !== state.name_group) {
    state.name_group = sanitized
    nameError.value = 'Удалены недопустимые символы'
    return false
  }

  if (!sanitized.trim()) {
    nameError.value = 'Укажите название группы'
    return false
  }

  nameError.value = ''
  return true
}

const validateScheduleItem = (item: GroupFormScheduleItem) => {
  const errors = ensureScheduleError(item.id)
  const normalizedDate = normalizeApiDate(item.date)
  const trimmedTheme = item.theme.trim().replace(/\s{2,}/g, ' ')

  if (trimmedTheme !== item.theme) {
    item.theme = trimmedTheme
  }

  errors.date = ''
  errors.theme = ''

  if (!normalizedDate) {
    errors.date = 'Укажите дату занятия'
  } else if (isRussianPublicHoliday(normalizedDate)) {
    errors.date = 'На государственный праздник РФ занятие поставить нельзя'
    item.date = ''
  } else {
    item.date = normalizedDate
  }

  if (!trimmedTheme) {
    errors.theme = 'Укажите тему занятия'
  }

  return !errors.date && !errors.theme
}

const validateSchedule = () => {
  if (!state.schedule.length) {
    notifications.error('Добавьте хотя бы одно занятие в расписание.', 'Группы')
    return false
  }

  const isValid = state.schedule.map(validateScheduleItem).every(Boolean)

  if (!isValid) {
    return false
  }

  const duplicateDates = new Set<string>()
  const seenDates = new Set<string>()

  for (const item of state.schedule) {
    const normalizedDate = normalizeApiDate(item.date)

    if (seenDates.has(normalizedDate)) {
      duplicateDates.add(normalizedDate)
    }

    seenDates.add(normalizedDate)
  }

  if (duplicateDates.size) {
    for (const item of state.schedule) {
      if (duplicateDates.has(normalizeApiDate(item.date))) {
        ensureScheduleError(item.id).date = 'Для одной группы нельзя дублировать дату занятия'
      }
    }

    return false
  }

  return true
}

const addScheduleItem = () => {
  state.schedule.push(createEmptyGroupScheduleItem())
}

const removeScheduleItem = (id: string) => {
  if (state.schedule.length === 1) {
    state.schedule = [createEmptyGroupScheduleItem()]
    scheduleErrors.value = {}
    return
  }

  state.schedule = state.schedule.filter((item) => item.id !== id)
  delete scheduleErrors.value[id]
}

const handleDateChange = (item: GroupFormScheduleItem) => {
  const normalizedDate = normalizeApiDate(item.date)

  if (!normalizedDate) {
    item.date = ''
    validateScheduleItem(item)
    return
  }

  if (isRussianPublicHoliday(normalizedDate)) {
    item.date = ''
    ensureScheduleError(item.id).date = 'Эта дата относится к государственным праздникам РФ'
    notifications.info('Праздничную дату выбрать нельзя. Укажите другой день.', 'Группы')
    return
  }

  item.date = normalizedDate
  validateScheduleItem(item)
}

const syncState = (nextState?: GroupFormState | null) => {
  const defaultState = createEmptyGroupFormState()
  state.name_group = nextState?.name_group || defaultState.name_group
  state.schedule = (nextState?.schedule?.length
    ? nextState.schedule
    : defaultState.schedule
  ).map((item) =>
    createEmptyGroupScheduleItem({
      date: normalizeApiDate(item.date),
      theme: item.theme
    })
  )

  nameError.value = ''
  scheduleErrors.value = {}
}

watch(
  () => props.initialState,
  (value) => {
    syncState(value)
  },
  { immediate: true, deep: true }
)

const handleSubmit = () => {
  const isValid = validateName() && validateSchedule()

  if (!isValid) {
    notifications.error('Проверьте название группы и расписание занятий.', 'Группы')
    return
  }

  emit('submit', {
    name_group: state.name_group.trim(),
    raspisanie: [...state.schedule]
      .map((item) => ({
        date: normalizeApiDate(item.date),
        theme: item.theme.trim()
      }))
      .sort((left, right) => getDateSortValue(left.date) - getDateSortValue(right.date))
  })
}
</script>

<template>
  <div class="page-grid">
    <AppCard :title="title">
      <form class="group-form" @submit.prevent="handleSubmit">
        <AppInput
          v-model="state.name_group"
          label="Название группы"
          placeholder="Например: P50-8-22"
          maxlength="120"
          :error="nameError"
          @update:model-value="validateName"
        />

        <div class="schedule-header">
          <div>
            <strong>Расписание группы</strong>
            <p>Добавляйте занятия по одному: дата и тема. Праздничные даты РФ для 2026 года недоступны.</p>
          </div>

          <AppButton type="button" variant="secondary" @click="addScheduleItem">
            Добавить занятие
          </AppButton>
        </div>

        <div class="schedule-list">
          <article
            v-for="(item, index) in state.schedule"
            :key="item.id"
            class="schedule-item"
          >
            <div class="schedule-item__top">
              <strong>Занятие {{ index + 1 }}</strong>
              <AppButton type="button" variant="ghost" @click="removeScheduleItem(item.id)">
                Удалить
              </AppButton>
            </div>

            <div class="schedule-item__grid">
              <AppInput
                v-model="item.date"
                type="date"
                label="Дата занятия"
                :error="scheduleErrors[item.id]?.date || ''"
                @update:model-value="handleDateChange(item)"
              />

              <AppInput
                v-model="item.theme"
                label="Тема занятия"
                placeholder="Например: Введение в программирование"
                maxlength="255"
                :error="scheduleErrors[item.id]?.theme || ''"
                @update:model-value="validateScheduleItem(item)"
              />
            </div>
          </article>
        </div>

        <div class="group-form__actions">
          <AppButton type="button" variant="ghost" @click="$emit('cancel')">Назад</AppButton>
          <AppButton type="submit" :disabled="loading">
            {{ loading ? 'Сохраняем...' : submitLabel }}
          </AppButton>
        </div>
      </form>
    </AppCard>

    <div class="stack">
      <AppCard title="Предпросмотр расписания">
        <p v-if="!schedulePreview.length" class="muted">
          Пока нет заполненных занятий. Как только вы укажете дату и тему, они появятся здесь по порядку.
        </p>

        <div v-else class="preview-list">
          <article v-for="item in schedulePreview" :key="item.id" class="preview-card">
            <strong>{{ item.theme || 'Без темы' }}</strong>
            <span>{{ item.date ? formatApiDate(item.date) : 'Дата не выбрана' }}</span>
          </article>
        </div>
      </AppCard>

      <AppCard title="Недоступные даты">
        <p class="muted">Для расписания заблокированы государственные праздники РФ в 2026 году:</p>
        <div class="holiday-list">
          <span v-for="date in blockedDates" :key="date" class="holiday-chip">
            {{ formatApiDate(date) }}
          </span>
        </div>
      </AppCard>
    </div>
  </div>
</template>

<style scoped>
.group-form,
.schedule-list,
.preview-list {
  display: grid;
  gap: 1rem;
}

.schedule-header,
.schedule-item__top,
.group-form__actions {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: center;
  flex-wrap: wrap;
}

.schedule-header p,
.muted {
  margin: 0.35rem 0 0;
  color: #64748b;
}

.schedule-item {
  padding: 1rem;
  border-radius: 1.15rem;
  background: rgba(248, 250, 252, 0.88);
  border: 1px solid rgba(15, 23, 42, 0.08);
}

.schedule-item__grid {
  display: grid;
  grid-template-columns: minmax(0, 0.85fr) minmax(0, 1.15fr);
  gap: 1rem;
  align-items: start;
  margin-top: 0.9rem;
}

.preview-card {
  display: grid;
  gap: 0.3rem;
  padding: 0.95rem 1rem;
  border-radius: 1rem;
  background: rgba(219, 234, 254, 0.38);
  border: 1px solid rgba(59, 130, 246, 0.14);
}

.preview-card span {
  color: #475569;
}

.holiday-list {
  display: flex;
  gap: 0.65rem;
  flex-wrap: wrap;
}

.holiday-chip {
  display: inline-flex;
  align-items: center;
  min-height: 2.2rem;
  padding: 0.4rem 0.8rem;
  border-radius: 999px;
  background: rgba(254, 226, 226, 0.8);
  color: #991b1b;
  border: 1px solid rgba(248, 113, 113, 0.24);
}

@media (max-width: 900px) {
  .schedule-item__grid {
    grid-template-columns: 1fr;
  }

  .group-form__actions {
    flex-direction: column;
    align-items: stretch;
  }
}
</style>
