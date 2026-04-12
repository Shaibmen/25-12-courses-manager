<script setup lang="ts">
import type { LegalEntityPayload } from '../../types/legalentity'
import LegalEntityForm from '../../components/features/legalentities/LegalEntityForm.vue'

const router = useRouter()
const notifications = useNotifications()
const saving = ref(false)

const submit = async (payload: LegalEntityPayload) => {
  saving.value = true

  try {
    await createLegalEntityRequest(payload)
    notifications.success('Юридическое лицо создано.', 'Юридические лица')
    await router.push('/legalentities')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать юридическое лицо',
      'Юридические лица'
    )
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="stack">
    <LegalEntityForm
      title="Добавить юридическое лицо"
      submit-label="Создать"
      :loading="saving"
      @submit="submit"
      @cancel="router.push('/legalentities')"
    />
  </section>
</template>
