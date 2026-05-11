<script setup lang="ts">
import type { GroupItem } from '../../types/group'
import { formatApiDate, getDateSortValue } from '../../utils/date'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../components/ui/AppInput.vue'

const router = useRouter()
const groups = ref<GroupItem[]>([])
const page = ref(1)
const filter = ref('')
const hasMore = ref(false)
const loading = ref(false)
const deleteLoading = ref(false)
const errorMessage = ref('')
const groupToDelete = ref<GroupItem | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const getScheduleCountLabel = (count: number) => {
  const mod10 = count % 10
  const mod100 = count % 100

  if (mod10 === 1 && mod100 !== 11) {
    return `${count} занятие`
  }

  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) {
    return `${count} занятия`
  }

  return `${count} занятий`
}

const getNextLesson = (group: GroupItem) =>
  [...(group.rapspisanie || [])]
    .sort((left, right) => getDateSortValue(left.date) - getDateSortValue(right.date))[0]

const loadGroupsList = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const data = await getGroups(page.value, filter.value)
    groups.value = data
    hasMore.value = data.length === 25
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось загрузить группы'
  } finally {
    loading.value = false
  }
}

const notifications = useNotifications()

const confirmDelete = async () => {
  if (!groupToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteGroupRequest(groupToDelete.value.group)
    groupToDelete.value = null
    await loadGroupsList()
    notifications.success('Группа удалена. Список обновлён.', 'Группы')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить группу',
      'Группы'
    )
  } finally {
    deleteLoading.value = false
  }
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    page.value = 1
    void loadGroupsList()
  }, 350)
})

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadGroupsList()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Группы">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по названию группы" />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/groups/create">Создать группу</AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю список групп и расписание занятий.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadGroupsList">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Список групп">
      <div class="table-wrap">
        <table class="groups-table">
          <thead>
            <tr>
              <th>Группа</th>
              <th>Ближайшее занятие</th>
              <th>Расписание</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!groups.length">
              <td colspan="4" class="groups-table__empty">Группы пока не найдены.</td>
            </tr>
            <tr v-for="group in groups" :key="group.group">
              <td>
                <div class="groups-table__main">
                  <strong>{{ group.name_group }}</strong>
                </div>
              </td>
              <td>
                <div class="groups-table__main">
                  <strong>{{ getNextLesson(group)?.theme || 'Занятий пока нет' }}</strong>
                  <span>
                    {{ getNextLesson(group)?.date ? formatApiDate(getNextLesson(group)?.date) : 'Дата не назначена' }}
                  </span>
                </div>
              </td>
              <td>{{ getScheduleCountLabel(group.rapspisanie?.length || 0) }}</td>
              <td>
                <div class="groups-table__actions">
                  <AppButton variant="ghost" @click="router.push(`/groups/${group.group}`)">
                    Подробнее
                  </AppButton>
                  <AppButton variant="secondary" @click="router.push(`/groups/edit/${group.group}`)">
                    Изменить
                  </AppButton>
                  <AppButton variant="ghost" @click="groupToDelete = group">
                    Удалить
                  </AppButton>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <AppButton variant="secondary" :disabled="page <= 1" @click="page > 1 && (page--, loadGroupsList())">
          Назад
        </AppButton>
        <span>Страница {{ page }}</span>
        <AppButton variant="secondary" :disabled="!hasMore" @click="page++, loadGroupsList()">
          Вперёд
        </AppButton>
      </div>
    </AppCard>

    <AppConfirmDialog
      :open="Boolean(groupToDelete)"
      title="Удаление группы"
      :message="groupToDelete ? `Удалить группу «${groupToDelete.name_group}»?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="groupToDelete = null"
      @confirm="confirmDelete"
    />
  </section>
</template>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  align-items: end;
  flex-wrap: wrap;
}

.toolbar__search {
  min-width: min(100%, 24rem);
  flex: 1 1 24rem;
}

.toolbar__actions,
.groups-table__actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.table-wrap {
  overflow-x: auto;
}

.groups-table {
  width: 100%;
  min-width: 1040px;
  border-collapse: collapse;
}

.groups-table th,
.groups-table td {
  padding: 1rem 0.95rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
  vertical-align: middle;
}

.groups-table th {
  color: #475569;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.groups-table tbody tr:hover {
  background: rgba(219, 234, 254, 0.26);
}

.groups-table__main {
  display: grid;
  gap: 0.3rem;
}

.groups-table__main span {
  color: #64748b;
  font-size: 0.86rem;
  overflow-wrap: anywhere;
}

.groups-table__empty {
  text-align: center;
  color: #64748b;
}

.pagination {
  margin-top: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 1rem;
}
</style>
