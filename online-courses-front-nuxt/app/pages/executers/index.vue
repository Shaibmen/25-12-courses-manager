<script setup lang="ts">
import type { ExecuterItem } from '../../types/executer'
import AppButton from '../../components/ui/AppButton.vue'
import AppCard from '../../components/ui/AppCard.vue'
import AppConfirmDialog from '../../components/ui/AppConfirmDialog.vue'
import AppInput from '../../components/ui/AppInput.vue'

const router = useRouter()
const notifications = useNotifications()
const executers = ref<ExecuterItem[]>([])
const page = ref(1)
const filter = ref('')
const sortOrder = ref<'asc' | 'desc'>('asc')
const hasMore = ref(false)
const loading = ref(false)
const deleteLoading = ref(false)
const errorMessage = ref('')
const executerToDelete = ref<ExecuterItem | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

const sortedExecuters = computed(() =>
  [...executers.value].sort((left, right) => {
    const leftName = (left.second_name || '').toLowerCase()
    const rightName = (right.second_name || '').toLowerCase()

    return sortOrder.value === 'asc'
      ? leftName.localeCompare(rightName)
      : rightName.localeCompare(leftName)
  })
)

const loadExecutersList = async () => {
  loading.value = true
  errorMessage.value = ''

  try {
    const data = await getExecuters(page.value, filter.value)
    executers.value = data
    hasMore.value = data.length === 25
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось загрузить список исполнителей'
  } finally {
    loading.value = false
  }
}

const confirmDelete = async () => {
  if (!executerToDelete.value) {
    return
  }

  deleteLoading.value = true

  try {
    await deleteExecuterRequest(executerToDelete.value.id_executor)
    executerToDelete.value = null
    await loadExecutersList()
    notifications.success('Исполнитель удалён. Список обновлён.', 'Исполнители')
  } catch (error) {
    notifications.error(
      error instanceof Error ? error.message : 'Не удалось удалить исполнителя',
      'Исполнители'
    )
  } finally {
    deleteLoading.value = false
  }
}

const toggleSort = () => {
  sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
}

watch(filter, () => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }

  debounceTimer = setTimeout(() => {
    page.value = 1
    void loadExecutersList()
  }, 350)
})

onBeforeUnmount(() => {
  if (debounceTimer) {
    clearTimeout(debounceTimer)
  }
})

onMounted(() => {
  void loadExecutersList()
})
</script>

<template>
  <section class="stack content-shell">
    <AppCard title="Исполнители">
      <div class="toolbar">
        <div class="toolbar__search">
          <AppInput v-model="filter" label="Поиск" placeholder="Поиск по фамилии" />
        </div>

        <div class="toolbar__actions">
          <AppButton to="/executers/create">Добавить исполнителя</AppButton>
          <AppButton variant="ghost" @click="toggleSort">
            Сортировать по фамилии {{ sortOrder === 'asc' ? '↑' : '↓' }}
          </AppButton>
          <AppButton variant="secondary" to="/dashboard">Назад</AppButton>
        </div>
      </div>
    </AppCard>

    <AppCard v-if="loading" title="Загрузка">
      <p>Подтягиваю список исполнителей.</p>
    </AppCard>

    <AppCard v-else-if="errorMessage" title="Ошибка">
      <p>{{ errorMessage }}</p>
      <AppButton variant="secondary" @click="loadExecutersList">Повторить запрос</AppButton>
    </AppCard>

    <AppCard v-else title="Список">
      <div class="table-wrap">
        <table class="executers-table">
          <thead>
            <tr>
              <th>Фамилия</th>
              <th>Имя</th>
              <th>Отчество</th>
              <th>Должность</th>
              <th>Доверенность</th>
              <th>Действия</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="!sortedExecuters.length">
              <td colspan="6" class="executers-table__empty">Исполнители пока не найдены.</td>
            </tr>
            <tr v-for="executer in sortedExecuters" :key="executer.id_executor">
              <td class="executers-table__text">{{ executer.second_name }}</td>
              <td class="executers-table__text">{{ executer.first_name }}</td>
              <td class="executers-table__text">{{ executer.middle_name || '—' }}</td>
              <td class="executers-table__text">{{ executer.status }}</td>
              <td class="executers-table__text">{{ executer.doverenost || '—' }}</td>
              <td>
                <AppButton variant="ghost" @click="executerToDelete = executer">
                  Удалить
                </AppButton>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="pagination">
        <AppButton variant="secondary" :disabled="page <= 1" @click="page > 1 && (page--, loadExecutersList())">
          Назад
        </AppButton>
        <span>Страница {{ page }}</span>
        <AppButton variant="secondary" :disabled="!hasMore" @click="page++, loadExecutersList()">
          Вперёд
        </AppButton>
      </div>
    </AppCard>

    <AppConfirmDialog
      :open="Boolean(executerToDelete)"
      title="Удаление исполнителя"
      :message="executerToDelete ? `Удалить исполнителя ${executerToDelete.second_name} ${executerToDelete.first_name}?` : ''"
      :loading="deleteLoading"
      confirm-label="Удалить"
      @cancel="executerToDelete = null"
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

.toolbar__actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
}

.table-wrap {
  overflow-x: auto;
}

.executers-table {
  width: 100%;
  min-width: 1040px;
  border-collapse: collapse;
}

.executers-table th,
.executers-table td {
  padding: 1rem 0.95rem;
  border-bottom: 1px solid rgba(15, 23, 42, 0.08);
  text-align: left;
  vertical-align: middle;
}

.executers-table th {
  color: #475569;
  font-size: 0.82rem;
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.executers-table tbody tr:hover {
  background: rgba(219, 234, 254, 0.26);
}

.executers-table__text {
  max-width: 14rem;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.executers-table__empty {
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
