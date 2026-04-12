<script setup lang="ts">
import type { ListenerFormPayload, ListenerFormState } from '../../../types/listener'
import { createEmptyListenerFormState } from '../../../types/listener'
import { normalizeApiDate } from '../../../utils/date'
import ListenerForm from '../../../components/features/listeners/ListenerForm.vue'
import AppCard from '../../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))

const initialState = ref<ListenerFormState | null>(null)
const loading = ref(true)
const saving = ref(false)
const loadError = ref('')

const load = async () => {
  loading.value = true
  loadError.value = ''

  try {
    const details = await getListenerDetails(id.value)
    const state = createEmptyListenerFormState()

    Object.assign(state.listener, details.listener, {
      date_of_birth: normalizeApiDate(details.listener.date_of_birth)
    })
    Object.assign(state.passport, details.passport || {}, {
      date_given: normalizeApiDate(details.passport?.date_given)
    })
    Object.assign(state.registrationAddress, details.regaddress || {})
    Object.assign(state.education, details.education_listener || {}, {
      date_given: normalizeApiDate(details.education_listener?.date_given),
      level_education: details.education_listener?.level_education || ''
    })
    Object.assign(state.placeWork, details.placework || {}, {
      all_experience: String(details.placework?.all_experience ?? ''),
      job_title_expirience: String(
        details.placework?.job_title_expirience ?? details.placework?.job_title_experience ?? ''
      )
    })

    initialState.value = state
  } catch (error) {
    loadError.value =
      error instanceof Error ? error.message : 'Не удалось загрузить данные слушателя'
  } finally {
    loading.value = false
  }
}

const submit = async (payload: ListenerFormPayload) => {
  saving.value = true

  try {
    await updateListenerRequest(id.value, payload)
    notifications.success('Данные слушателя обновлены.', 'Слушатели')
    await router.push('/listeners')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось обновить данные слушателя',
      'Слушатели'
    )
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack">
    <AppCard title="Редактирование слушателя">
      <p>Обновите личные данные, документы и связанные сведения по слушателю.</p>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные слушателя для редактирования.</p>
    </AppCard>

    <AppCard v-else-if="loadError" title="Ошибка">
      <p>{{ loadError }}</p>
    </AppCard>

    <ListenerForm
      v-else
      title="Редактирование слушателя"
      submit-label="Сохранить изменения"
      :loading="saving"
      :initial-state="initialState"
      mode="edit"
      @submit="submit"
      @cancel="router.push('/listeners')"
    />
  </section>
</template>
