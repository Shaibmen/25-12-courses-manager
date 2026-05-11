<script setup lang="ts">
import type { GroupFormState, GroupPayload } from '../../../types/group'
import { createEmptyGroupFormState, createEmptyGroupScheduleItem } from '../../../types/group'
import GroupForm from '../../../components/features/groups/GroupForm.vue'
import AppCard from '../../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))

const initialState = ref<GroupFormState | null>(null)
const loading = ref(true)
const saving = ref(false)
const loadError = ref('')

const load = async () => {
  loading.value = true
  loadError.value = ''

  try {
    const group = await getGroupDetails(id.value)
    initialState.value = {
      ...createEmptyGroupFormState(),
      name_group: group.name_group || '',
      schedule: (group.rapspisanie || []).length
        ? (group.rapspisanie || []).map((item) =>
            createEmptyGroupScheduleItem({
              date: item.date,
              theme: item.theme
            })
          )
        : createEmptyGroupFormState().schedule
    }
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : 'Не удалось загрузить группу'
  } finally {
    loading.value = false
  }
}

const submit = async (payload: GroupPayload) => {
  saving.value = true

  try {
    await updateGroupRequest(id.value, payload)
    notifications.success('Группа обновлена.', 'Группы')
    await router.push('/groups')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось обновить группу',
      'Группы'
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
  <section class="stack content-shell">
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные группы для редактирования.</p>
    </AppCard>

    <AppCard v-else-if="loadError" title="Ошибка">
      <p>{{ loadError }}</p>
    </AppCard>

    <GroupForm
      v-else
      title="Редактировать группу"
      submit-label="Сохранить изменения"
      :loading="saving"
      :initial-state="initialState"
      @submit="submit"
      @cancel="router.push('/groups')"
    />
  </section>
</template>
