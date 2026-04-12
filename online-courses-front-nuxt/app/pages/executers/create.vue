<script setup lang="ts">
import type { ExecuterPayload } from '../../types/executer'
import ExecuterForm from '../../components/features/executers/ExecuterForm.vue'

const router = useRouter()
const notifications = useNotifications()
const saving = ref(false)

const submit = async (payload: ExecuterPayload) => {
  saving.value = true

  try {
    await createExecuterRequest(payload)
    notifications.success('Исполнитель создан.', 'Исполнители')
    await router.push('/executers')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать исполнителя',
      'Исполнители'
    )
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="stack">
    <ExecuterForm
      title="Добавить исполнителя"
      submit-label="Создать"
      :loading="saving"
      @submit="submit"
      @cancel="router.push('/executers')"
    />
  </section>
</template>
