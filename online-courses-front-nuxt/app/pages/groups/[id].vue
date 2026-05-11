<script setup lang="ts">
import type { GroupItem } from '../../types/group'
import { exportGroupSchedule, getGroupDetails } from '../../services/api/groups'
import GroupDetailsView from '../../components/features/groups/GroupDetailsView.vue'
import AppCard from '../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))

const group = ref<GroupItem | null>(null)
const loading = ref(true)
const exportLoading = ref(false)
const loadError = ref('')

const getScheduleExportFileName = (groupName: string) =>
  `Расписание - ${groupName.trim() || 'Группа'}`
    .replace(/[\\/:*?"<>|]/g, '-')
    .replace(/\s{2,}/g, ' ')
    .concat('.xlsx')

const load = async () => {
  loading.value = true
  loadError.value = ''

  try {
    group.value = await getGroupDetails(id.value)
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : 'Не удалось загрузить группу'
  } finally {
    loading.value = false
  }
}

const exportSchedule = async () => {
  if (!group.value || exportLoading.value) {
    return
  }

  exportLoading.value = true

  try {
    const blob = await exportGroupSchedule(group.value.group)
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = getScheduleExportFileName(group.value.name_group)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.setTimeout(() => URL.revokeObjectURL(url), 0)
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось скачать расписание группы',
      'Группы'
    )
  } finally {
    exportLoading.value = false
  }
}

onMounted(() => {
  void load()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю данные группы.</p>
    </AppCard>

    <AppCard v-else-if="loadError" title="Ошибка">
      <p>{{ loadError }}</p>
    </AppCard>

    <GroupDetailsView
      v-else-if="group"
      :group="group"
      :export-loading="exportLoading"
      @back="router.push('/groups')"
      @edit="router.push(`/groups/edit/${id}`)"
      @export="exportSchedule"
    />
  </section>
</template>
