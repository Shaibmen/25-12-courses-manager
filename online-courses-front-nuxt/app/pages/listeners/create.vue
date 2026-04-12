<script setup lang="ts">
import type { ListenerFormPayload } from '../../types/listener'
import ListenerForm from '../../components/features/listeners/ListenerForm.vue'
import AppCard from '../../components/ui/AppCard.vue'

const router = useRouter()
const route = useRoute()
const notifications = useNotifications()
const loading = ref(false)

const submit = async (payload: ListenerFormPayload) => {
  loading.value = true

  try {
    const legalEntityId = typeof route.query.id_legalentity === 'string'
      ? route.query.id_legalentity
      : ''

    if (legalEntityId) {
      payload.listener.id_legalentity = legalEntityId
    }

    await createListenerRequest(payload)
    notifications.success('Слушатель создан.', 'Слушатели')
    await router.push('/listeners')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать слушателя',
      'Слушатели'
    )
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <section class="stack">
    <AppCard title="Создать нового слушателя">
      <p>Заполните основные данные, документы, адрес и дополнительные сведения.</p>
    </AppCard>
    <ListenerForm
      title="Создать нового слушателя"
      submit-label="Создать"
      :loading="loading"
      mode="create"
      @submit="submit"
      @cancel="router.push('/listeners')"
    />
  </section>
</template>
