<script setup lang="ts">
import type { GroupItem } from '../../types/group'
import GroupDetailsView from '../../components/features/groups/GroupDetailsView.vue'
import AppCard from '../../components/ui/AppCard.vue'

const route = useRoute()
const router = useRouter()
const notifications = useNotifications()
const id = computed(() => String(route.params.id || ''))

const group = ref<GroupItem | null>(null)
const loading = ref(true)
const loadError = ref('')

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

const exportStub = () => {
  notifications.info('Экспорт пока работает как заглушка. Следующим шагом подключим действие.', 'Группы')
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
      @back="router.push('/groups')"
      @edit="router.push(`/groups/edit/${id}`)"
      @export="exportStub"
    />
  </section>
</template>
