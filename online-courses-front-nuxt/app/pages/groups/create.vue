<script setup lang="ts">
import type { GroupPayload } from '../../types/group'
import GroupForm from '../../components/features/groups/GroupForm.vue'

const router = useRouter()
const notifications = useNotifications()
const saving = ref(false)

const submit = async (payload: GroupPayload) => {
  saving.value = true

  try {
    await createGroupRequest(payload)
    notifications.success('Группа создана.', 'Группы')
    await router.push('/groups')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось создать группу',
      'Группы'
    )
  } finally {
    saving.value = false
  }
}
</script>

<template>
  <section class="stack content-shell">
    <GroupForm
      title="Создать группу"
      submit-label="Создать группу"
      :loading="saving"
      @submit="submit"
      @cancel="router.push('/groups')"
    />
  </section>
</template>
