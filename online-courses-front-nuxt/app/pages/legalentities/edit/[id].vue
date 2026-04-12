<script setup lang="ts">
import type { LegalEntityFormState, LegalEntityPayload } from '../../../types/legalentity'
import { createEmptyLegalEntityFormState } from '../../../types/legalentity'
import LegalEntityForm from '../../../components/features/legalentities/LegalEntityForm.vue'
import AppCard from '../../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))
const initialState = ref<LegalEntityFormState | null>(null)
const loading = ref(true)
const saving = ref(false)
const loadError = ref('')

const load = async () => {
  loading.value = true
  loadError.value = ''

  try {
    const details = await getLegalEntityDetails(id.value)
    const state = createEmptyLegalEntityFormState()

    Object.assign(state.legal, details.legal_entity || {})
    Object.assign(state.regAddress, details.reg_address || {})
    initialState.value = state
  } catch (error) {
    loadError.value =
      error instanceof Error ? error.message : 'Не удалось загрузить юридическое лицо'
  } finally {
    loading.value = false
  }
}

const submit = async (payload: LegalEntityPayload) => {
  saving.value = true

  try {
    await updateLegalEntityRequest(id.value, payload)
    notifications.success('Юридическое лицо обновлено.', 'Юридические лица')
    await router.push('/legalentities')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось обновить юридическое лицо',
      'Юридические лица'
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
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю юридическое лицо.</p>
    </AppCard>

    <AppCard v-else-if="loadError" title="Ошибка">
      <p>{{ loadError }}</p>
    </AppCard>

    <LegalEntityForm
      v-else
      title="Редактировать юридическое лицо"
      submit-label="Сохранить изменения"
      :loading="saving"
      :initial-state="initialState"
      @submit="submit"
      @cancel="router.push('/legalentities')"
    />
  </section>
</template>
